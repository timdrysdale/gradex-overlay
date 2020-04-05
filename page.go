package main

import (
	"fmt"

	creator "github.com/unidoc/unipdf/v3/creator"
)

// see https://github.com/unidoc/unipdf-examples/blob/master/image/pdf_add_image_to_page.go
// xPos and yPos define the upper left corner of the image location, and iwidth
// is the width of the image in PDF document dimensions (height/width ratio is maintained).

func AddImagePage(imgPath string, c *creator.Creator) error {

	// load image

	// decide on page size and orientation (A4+boxwidth only, for now) based on image

	// decide image scaling to suit page

	// place on page

	// Prepare the image.
	img, err := c.NewImageFromFile(imgPath)
	if err != nil {
		return err
	}

	fmt.Printf("%f %f\n", img.Width(), img.Height())

	// start out as A4 portrait, swap to landscape if needbe
	barWidth := 50 * creator.PPMM
	A4Width := 210 * creator.PPMM
	A4Height := 297 * creator.PPMM
	pageWidth := A4Width + barWidth
	pageHeight := A4Height
	imgLeft := 0.0

	isLandscape := img.Height() < img.Width()

	if isLandscape {
		pageWidth = pageWidth + barWidth
		pageHeight = A4Width
		imgLeft = imgLeft + barWidth
	}
	img.ScaleToHeight(pageHeight)
	img.SetPos(0, imgLeft) //top,left

	fmt.Printf("%f %f\n", img.Width(), img.Height())

	c.SetPageSize(creator.PageSize{pageWidth, A4Height})

	c.NewPage()

	c.Draw(img)

	return nil

}
