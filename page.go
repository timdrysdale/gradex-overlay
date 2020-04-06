package main

import (
	creator "github.com/unidoc/unipdf/v3/creator"
)

// see https://github.com/unidoc/unipdf-examples/blob/master/image/pdf_add_image_to_page.go
// xPos and yPos define the upper left corner of the image location, and iwidth
// is the width of the image in PDF document dimensions (height/width ratio is maintained).

func AddImagePage(imgPath string, c *creator.Creator) (*markOpt, error) {

	// load image
	img, err := c.NewImageFromFile(imgPath)
	if err != nil {
		return &markOpt{}, err
	}

	// start out as A4 portrait, swap to landscape if need be
	barWidth := 25 * creator.PPMM
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

	// create new page with image
	c.SetPageSize(creator.PageSize{pageWidth, pageHeight})
	c.NewPage()
	c.Draw(img)

	opts := &markOpt{
		left:             isLandscape,
		right:            true,
		barwidth:         barWidth,
		pageWidth:        pageWidth,
		pageHeight:       pageHeight,
		marksEvery:       25 * creator.PPMM,
		markHeight:       12 * creator.PPMM,
		markWidth:        20 * creator.PPMM,
		markMargin:       0 * creator.PPMM,
		markBottomMargin: 5 * creator.PPMM,
	}

	return opts, nil
}
