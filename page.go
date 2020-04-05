package main

import creator "github.com/unidoc/unipdf/v3/creator"

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
	img.ScaleToWidth(210 * creator.PPMM)
	img.SetPos(0, 0) //top, left for now

	c.SetPageSize(creator.PageSize{260 * creator.PPMM, 297 * creator.PPMM})

	c.NewPage()

	c.Draw(img)

	return nil

}
