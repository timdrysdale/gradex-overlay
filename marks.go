package main

import (
	"bytes"
	"fmt"
	"math"
	"os"

	"github.com/mattetti/filebuffer"
	"github.com/unidoc/unipdf/v3/annotator"
	creator "github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	pdf "github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

type markOpt struct {
	left             bool
	right            bool
	barwidth         float64
	pageWidth        float64
	pageHeight       float64
	marksEvery       float64
	markHeight       float64
	markWidth        float64
	markMargin       float64
	markBottomMargin float64
}

func createMarks(page *model.PdfPage, opt markOpt, formID string) *model.PdfAcroForm {

	form := model.NewPdfAcroForm()

	// mirror each other, close to the page
	xright := opt.pageWidth - opt.barwidth + opt.markMargin
	xleft := opt.barwidth - opt.markWidth - opt.markMargin

	numMarks := math.Floor(opt.pageHeight / opt.marksEvery)

	yTop := 0.5 * (opt.pageHeight + opt.marksEvery*(numMarks-1) - opt.markHeight)

	if opt.left {
		for idx := 0; idx < int(numMarks); idx = idx + 1 {
			yPos := yTop - (float64(idx) * opt.marksEvery)
			tfopt := annotator.TextFieldOptions{}
			name := fmt.Sprintf("%s-left-%02d", formID, idx)
			rect := []float64{xleft, yPos, xleft + opt.markWidth, yPos + opt.markHeight}
			textf, err := annotator.NewTextField(page, name, rect, tfopt)
			if err != nil {
				panic(err)
			}
			*form.Fields = append(*form.Fields, textf.PdfField)
			page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
		}

	}

	// right
	if opt.right {
		for idx := 0; idx < int(numMarks); idx = idx + 1 {
			yPos := yTop - (float64(idx) * opt.marksEvery)
			tfopt := annotator.TextFieldOptions{}
			name := fmt.Sprintf("%s-right-%02d", formID, idx)
			rect := []float64{xright, yPos, xright + opt.markWidth, yPos + opt.markHeight}
			textf, err := annotator.NewTextField(page, name, rect, tfopt)
			if err != nil {
				panic(err)
			}
			*form.Fields = append(*form.Fields, textf.PdfField)
			page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
		}
	}

	return form
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

	pdfWriter.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    80,
		ImageUpperPPI:                   100,
	}))

	pdfWriter.Write(of)
}
