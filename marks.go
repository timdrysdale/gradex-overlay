package main

import (
	"fmt"
	"math"

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

	numMarks := math.Floor((opt.pageHeight - opt.markBottomMargin) / opt.marksEvery)
	ytop := opt.markBottomMargin + ((numMarks - 1) * opt.marksEvery)

	//do this way to get tab order correct (left column, then right column, top to bottom)
	//for ypos := opt.markBottomMargin; ypos < opt.pageHeight-opt.marksEvery; ypos = ypos + opt.marksEvery {

	if opt.left {
		for idx := 0; idx < int(numMarks); idx = idx + 1 {

			ypos := ytop - (float64(idx) * opt.marksEvery)

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

	}

	// right
	if opt.right {
		for idx := 0; idx < int(numMarks); idx = idx + 1 {
			ypos := ytop - (float64(idx) * opt.marksEvery)
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
	}

	return form
}
