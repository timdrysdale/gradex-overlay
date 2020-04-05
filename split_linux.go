package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// https://github.com/catherinelu/evangelist/blob/master/server.go

// possible alpha numeric characters

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

func convertPDFToJPEGs(pdfPath string, jpegPath string) error {

	suffix := filepath.Ext(pdfPath)
	basename := strings.TrimSuffix(pdfPath, suffix)

	outputFileOption := fmt.Sprintf("-sOutputFile=./jpg/%s%%04d.jpg", basename)

	cmd := exec.Command("gs", "-dNOPAUSE", "-sDEVICE=jpeg", outputFileOption, "-dJPEGQ=95", "-r300", "-q", pdfPath,
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
