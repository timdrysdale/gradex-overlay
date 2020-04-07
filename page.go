package main

import (
	"math"

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
	barWidth := 30 * creator.PPMM
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

	// these are tweaked - see vspace hack
	// TODO make this make sense

	opt := &markOpt{
		left:             isLandscape,
		right:            true,
		barwidth:         barWidth,
		pageWidth:        pageWidth,
		pageHeight:       pageHeight,
		marksEvery:       26.25 * creator.PPMM,
		markHeight:       18 * creator.PPMM,
		markWidth:        20 * creator.PPMM,
		markMargin:       5 * creator.PPMM,
		markBottomMargin: 0 * creator.PPMM,
	}

	// coloured box for the marks

	boxX := pageWidth - barWidth
	boxY := 0.0

	rect := c.NewRectangle(boxX, boxY, barWidth, pageHeight)
	rect.SetBorderColor(creator.ColorRed)
	rect.SetFillColor(creator.ColorRGBFromHex("#FFCCCB"))
	c.Draw(rect)

	if isLandscape {
		boxX = 0.0
		rect = c.NewRectangle(boxX, boxY, barWidth, pageHeight)
		rect.SetBorderColor(creator.ColorRed)
		rect.SetFillColor(creator.ColorRGBFromHex("#FFCCCB"))
		c.Draw(rect)

	}

	xright := opt.pageWidth - opt.barwidth + opt.markMargin
	xleft := opt.barwidth - opt.markWidth - opt.markMargin

	numMarks := math.Floor(opt.pageHeight / opt.marksEvery)

	yTop := 0.5 * (pageHeight + opt.marksEvery*(numMarks-1) - opt.markHeight)

	for idx := 0; idx < int(numMarks); idx = idx + 1 {
		yPos := yTop - (float64(idx) * opt.marksEvery)
		if opt.left {
			rect = c.NewRectangle(xleft, yPos, opt.markWidth, opt.markHeight)
			rect.SetBorderColor(creator.ColorRed)
			rect.SetFillColor(creator.ColorRGBFromHex("#FFFFFF"))
			c.Draw(rect)
		}
		if opt.right {

			rect = c.NewRectangle(xright, yPos, opt.markWidth, opt.markHeight)
			rect.SetBorderColor(creator.ColorRed)
			rect.SetFillColor(creator.ColorRGBFromHex("#FFFFFF"))
			c.Draw(rect)
		}

	}

	return opt, nil
}
