package main

import (
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// https://github.com/catherinelu/evangelist/blob/master/server.go

// allow at most 1 MB of form data to be passed to the server
const MAX_MULTIPART_FORM_BYTES = 1024 * 1024

// number of workers to run simultaneously to convert a PDF
const NUM_WORKERS_CONVERT = 2

// number of workers to run simultaneously to upload a PDF
const NUM_WORKERS_UPLOAD = 10

// possible alpha numeric characters
const ALPHA_NUMERIC = "abcdefghijklmnopqrstuvwxyz0123456789"

/* Returns the number of pages in the PDF specified by `pdfPath`. */
func getNumPages(pdfPath string) (int, error) {
	// ghostscript can retrieve us the number of pages
	cmd := exec.Command("gs", "-q", "-dNODISPLAY", "-c",
		fmt.Sprintf("(%s) (r) file runpdfbegin pdfpagecount = quit", pdfPath))
	numPagesBytes, err := cmd.Output()

	// convert []byte -> string -> int (painful, but necessary)
	if err != nil {
		return -1, err
	}
	numPagesStr := strings.Trim(string(numPagesBytes), "\n")
	numPagesInt64, err := strconv.ParseInt(numPagesStr, 10, 0)

	if err != nil {
		return -1, err
	}
	return int(numPagesInt64), nil
}

/* Resizes the JPEG at `jpegPath` to have a width at most `maxWidth` and
 * a height at most `maxHeight`. Maintains aspect ratio. Saves the resized
 * JPEG to `resizedJPEGPath`. */
func resizeAndSaveImage(jpegPath string, resizedJPEGPath string, maxWidth int,
	maxHeight int) error {
	dimension := fmt.Sprintf("%dx%d", maxWidth, maxHeight)
	cmd := exec.Command("convert", "-resize", dimension, jpegPath, resizedJPEGPath)
	return cmd.Run()
}

/* Converts the PDF at `pdfPath` to JPEGs. Outputs the JPEGs to the provided
 * `jpegPath` (note: '%d' in `jpegPath` will be replaced by the JPEG
 * number). Returns the path to the JPEGs (contains a %d that should be
 * replaced with the page number) and the number of pages in the PDF. */
func convertPDFToJPEGsParallel(pdfPath string, jpegPath string, smallJPEGPath string,
	largeJPEGPath string) (int, error) {
	numPages, err := getNumPages(pdfPath)
	if err != nil {
		return -1, err
	}

	// find number of pages to convert per worker
	numPagesPerWorkerFloat64 := float64(numPages) / float64(NUM_WORKERS_CONVERT)
	numPagesPerWorker := int(math.Ceil(numPagesPerWorkerFloat64))

	var wg sync.WaitGroup

	for firstPage := 1; firstPage <= numPages; firstPage = firstPage + numPagesPerWorker {
		// spawn workers, keeping track of them to wait until they're finished
		wg.Add(1)
		lastPage := firstPage + numPagesPerWorker - 1
		if lastPage > numPages {
			lastPage = numPages
		}

		go convertPagesToJPEGs(&wg, pdfPath, jpegPath, smallJPEGPath,
			largeJPEGPath, firstPage, lastPage)
	}

	wg.Wait()
	return numPages, err
}

/* Converts the PDF at `pdfPath` to JPEGs. Outputs the JPEGs to the provided
 * `jpegPath` (note: '%d' in `jpegPath` will be replaced by the JPEG
 * number). Converts pages within the range [`firstPage`, `lastPage`]. Calls
 * `wg.Done()` once finished. Returns an error on the given channel. */
func convertPagesToJPEGs(wg *sync.WaitGroup, pdfPath string, jpegPath string,
	smallJPEGPath string, largeJPEGPath string, firstPage int, lastPage int) {
	defer wg.Done()

	// use ghostscript for PDF -> JPEG conversion at 300 density
	for pageNum := firstPage; pageNum <= lastPage; pageNum = pageNum + 1 {
		// convert a single page at a time with the correct output JPEG path
		firstPageOption := fmt.Sprintf("-dFirstPage=%d", pageNum)
		lastPageOption := fmt.Sprintf("-dLastPage=%d", pageNum)

		// convert to two sizes: normal and large
		jpegPathForPage := fmt.Sprintf(jpegPath, pageNum)
		smallJPEGPathForPage := fmt.Sprintf(smallJPEGPath, pageNum)
		largeJPEGPathForPage := fmt.Sprintf(largeJPEGPath, pageNum)

		outputFileOption := fmt.Sprintf("-sOutputFile=%s", largeJPEGPathForPage)

		cmd := exec.Command("gs", "-dNOPAUSE", "-sDEVICE=jpeg", firstPageOption,
			lastPageOption, outputFileOption, "-dJPEGQ=90", "-r200", "-q", pdfPath,
			"-c", "quit")
		err := cmd.Run()

		if err != nil {
			fmt.Printf("gs command failed: %s\n", err.Error())
			return
		}

		resizeAndSaveImage(largeJPEGPathForPage, jpegPathForPage, 800, 800)
		if err != nil {
			fmt.Printf("Couldn't resize image: %s\n", err.Error())
			return
		}

		resizeAndSaveImage(jpegPathForPage, smallJPEGPathForPage, 300, 300)
		if err != nil {
			fmt.Printf("Couldn't resize image: %s\n", err.Error())
			return
		}
	}
}

func convertPDFToJPEGs(pdfPath string, jpegPath string, smallJPEGPath string,
	largeJPEGPath string) error {

	suffix := filepath.Ext(pdfPath)
	basename := strings.TrimSuffix(pdfPath, suffix)

	outputFileOption := fmt.Sprintf("-sOutputFile=%s%%03d.jpg", basename)

	cmd := exec.Command("gs", "-dNOPAUSE", "-sDEVICE=jpeg", outputFileOption, "-dJPEGQ=90", "-r200", "-q", pdfPath,
		"-c", "quit")

	err := cmd.Run()
	if err != nil {
		fmt.Printf("gs command failed: %s\n", err.Error())
		return err
	}

	return nil
}

// This worked
// gs -dNOPAUSE -sDEVICE=jpeg -sOutputFile=edited-%d.jpg -dJPEGQ=95 -r300 -q edited5-covered.pdf -c quit
