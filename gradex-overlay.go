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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	unicommon "github.com/unidoc/unipdf/v3/common"
	creator "github.com/unidoc/unipdf/v3/creator"
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

	basename := strings.TrimSuffix(inputPath, suffix)
	outputPath := basename + "-mark" + suffix

	jpegPath := "./jpg"

	err := ensureDir(jpegPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = convertPDFToJPEGs(inputPath, jpegPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing

	AddImagePage("./jpg/edited5-covered0005.jpg", c)

	err = c.WriteToFile(outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
