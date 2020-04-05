package main

import (
	"fmt"

	"github.com/unidoc/unipdf/v3/annotator"
	creator "github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

// see https://github.com/unidoc/unipdf-examples/blob/master/image/pdf_add_image_to_page.go
// xPos and yPos define the upper left corner of the image location, and iwidth
// is the width of the image in PDF document dimensions (height/width ratio is maintained).

func AddImagePage(imgPath string, markID string, c *creator.Creator, form *model.PdfAcroForm) error {

	// load image
	img, err := c.NewImageFromFile(imgPath)
	if err != nil {
		return err
	}

	fmt.Printf("%f %f\n", img.Width(), img.Height())

	// Choose page size
	// start out as A4 portrait, swap to landscape if need be
	barWidth := 50 * creator.PPMM
	A4Width := 210 * creator.PPMM
	A4Height := 297 * creator.PPMM
	pageWidth := A4Width + barWidth
	pageHeight := A4Height
	imgLeft := 0.0

	isLandscape := img.Height() < img.Width()

	if isLandscape {
		pageWidth = A4Height + (2 * barWidth)
		pageHeight = A4Width
		imgLeft = barWidth
	}

	// scale and position image
	img.ScaleToHeight(pageHeight)
	img.SetPos(imgLeft, 0) //left, top

	// create new page
	c.SetPageSize(creator.PageSize{pageWidth, pageHeight})
	page := c.NewPage()

	// add image
	c.Draw(img)

	// add mark box forms, prepend markID string
	opt := annotator.TextFieldOptions{}
	opt.Value = "THIS IS MY VALUE"
	textf, err := annotator.NewTextField(page, "markbox", []float64{0, 0, 200, 200}, opt)
	if err != nil {
		panic(err)
	}

	*form.Fields = append(*form.Fields, textf.PdfField)
	page.AddAnnotation(textf.Annotations[0].PdfAnnotation)

	return nil

}
