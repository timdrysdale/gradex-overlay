package main

/*
 * Add a cover page to a PDF file
 * Generates cover page then merges, including form field data (AcroForms).
 *
 * Run as: gradex-coverpage <barefile>.pdf
 *
 * outputs: <barefile>-covered.pdf (using internally generated cover page)
 *
 * Adapted from github.com/unidoc/unipdf-examples/pages/pdf_merge_advanced.go
 *
 *
 */

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattetti/filebuffer"
	unicommon "github.com/unidoc/unipdf/v3/common"
	creator "github.com/unidoc/unipdf/v3/creator"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Debug log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelInfo))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Requires one argument: input_path\n")
		fmt.Printf("Usage: gradex-coverpage.exe input.pdf\n")
		os.Exit(0)
	}

	inputPath := os.Args[1]

	suffix := filepath.Ext(inputPath)

	// sanity check
	if suffix != ".pdf" {
		fmt.Printf("Error: input path must be a .pdf\n")
		os.Exit(1)
	}

	// need page count to find the jpeg files again later
	numPages, err := countPages(inputPath)

	// render to images
	jpegPath := "./jpg"
	err = ensureDir(jpegPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	basename := strings.TrimSuffix(inputPath, suffix)
	jpegFileOption := fmt.Sprintf("%s/%s%%04d.jpg", jpegPath, basename)

	err = convertPDFToJPEGs(inputPath, jpegPath, jpegFileOption)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// convert images to individual pdfs, with form overlay

	pagePath := "./pdf"
	err = ensureDir(pagePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pageFileOption := fmt.Sprintf("%s/%s%%04d.pdf", pagePath, basename)
	formNameOption := fmt.Sprintf("%s%%04d", basename)

	// gs starts indexing at 1
	for imgIdx := 1; imgIdx <= numPages; imgIdx = imgIdx + 1 {

		// construct image name
		jpegFilename := fmt.Sprintf(jpegFileOption, imgIdx)
		pageFilename := fmt.Sprintf(pageFileOption, imgIdx)
		formID := fmt.Sprintf(formNameOption, imgIdx)

		// do the overlay
		convertJPEGToOverlaidPDF(jpegFilename, pageFilename, formID)

	}

}

func convertJPEGToOverlaidPDF(jpegFilename string, pageFilename string, formID string) {

	c := creator.New()

	c.SetPageMargins(0, 0, 0, 0) // we're not printing

	markOptions, err := AddImagePage(jpegFilename, c) //isLandscape
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// write to memory
	var buf bytes.Buffer

	err = c.Write(&buf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := pdf.NewPdfReader(fbuf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pdfWriter := pdf.NewPdfWriter()

	page, err := pdfReader.GetPage(1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = pdfWriter.SetForms(createMarks(page, *markOptions, formID))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = pdfWriter.AddPage(page)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	of, err := os.Create(pageFilename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer of.Close()

	pdfWriter.Write(of)
}
