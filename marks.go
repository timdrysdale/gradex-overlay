package main

import (
	"fmt"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/model"
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

	idx := 0

	for ypos := opt.markBottomMargin; ypos < opt.pageHeight-opt.marksEvery; ypos = ypos + opt.marksEvery {

		// right
		if opt.right {
			tfopt := annotator.TextFieldOptions{}
			name := fmt.Sprintf("%s-right-%02d", formID, idx)
			rect := []float64{xright, ypos, xright + opt.markWidth, ypos + opt.markHeight}
			textf, err := annotator.NewTextField(page, name, rect, tfopt)
			if err != nil {
				panic(err)
			}
			*form.Fields = append(*form.Fields, textf.PdfField)
			page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
		}

		if opt.left {

			tfopt := annotator.TextFieldOptions{}
			name := fmt.Sprintf("%s-left-%02d", formID, idx)
			rect := []float64{xleft, ypos, xleft + opt.markWidth, ypos + opt.markHeight}
			textf, err := annotator.NewTextField(page, name, rect, tfopt)
			if err != nil {
				panic(err)
			}
			*form.Fields = append(*form.Fields, textf.PdfField)
			page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
		}

		idx = idx + 1
	}

	return form
}
