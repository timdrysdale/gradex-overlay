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
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattetti/filebuffer"
	"github.com/unidoc/unipdf/v3/annotator"
	unicommon "github.com/unidoc/unipdf/v3/common"
	creator "github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	pdf "github.com/unidoc/unipdf/v3/model"
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

	//basename := strings.TrimSuffix(inputPath, suffix)
	//outputPath := basename + "-mark" + suffix

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

	form := model.NewPdfAcroForm()

	AddImagePage("./jpg/edited5-covered0005.jpg", "page5", c, form)

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

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pdfWriter := pdf.NewPdfWriter()

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		err = pdfWriter.SetForms(createForm(page))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	of, err := os.Create("./testForm.pdf")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer of.Close()

	pdfWriter.Write(of)
}

// textFieldsDef is a list of text fields to add to the form. The Rect field specifies the coordinates of the
// field.
var textFieldsDef = []struct {
	Name string
	Rect []float64
}{
	{Name: "full_name", Rect: []float64{123.97, 619.02, 343.99, 633.6}},
	{Name: "address_line_1", Rect: []float64{142.86, 596.82, 347.3, 611.4}},
	{Name: "address_line_2", Rect: []float64{143.52, 574.28, 347.96, 588.86}},
	{Name: "age", Rect: []float64{95.15, 551.75, 125.3, 566.33}},
	{Name: "city", Rect: []float64{96.47, 506.35, 168.37, 520.93}},
	{Name: "country", Rect: []float64{114.69, 483.82, 186.59, 498.4}},
}

// checkboxFieldDefs is a list of checkboxes to add to the form.
var checkboxFieldDefs = []struct {
	Name    string
	Rect    []float64
	Checked bool
}{
	{Name: "male", Rect: []float64{113.7, 525.57, 125.96, 540.15}, Checked: true},
	{Name: "female", Rect: []float64{157.44, 525.24, 169.7, 539.82}, Checked: false},
}

// choiceFieldDefs is a list of comboboxes to add to the form with specified options.
var choiceFieldDefs = []struct {
	Name    string
	Rect    []float64
	Options []string
}{
	{
		Name:    "fav_color",
		Rect:    []float64{144.52, 461.61, 243.92, 476.19},
		Options: []string{"Black", "Blue", "Green", "Orange", "Red", "White", "Yellow"},
	},
}

// createForm creates the form and fields to be placed on the `page`.
func createForm(page *model.PdfPage) *model.PdfAcroForm {
	form := model.NewPdfAcroForm()

	// Add ZapfDingbats font.
	zapfdb := model.NewStandard14FontMustCompile(model.ZapfDingbatsName)
	form.DR = model.NewPdfPageResources()
	form.DR.SetFontByName(`ZaDb`, zapfdb.ToPdfObject())

	for _, fdef := range textFieldsDef {
		opt := annotator.TextFieldOptions{}
		textf, err := annotator.NewTextField(page, fdef.Name, fdef.Rect, opt)
		if err != nil {
			panic(err)
		}

		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
	}

	for _, cbdef := range checkboxFieldDefs {
		opt := annotator.CheckboxFieldOptions{}
		checkboxf, err := annotator.NewCheckboxField(page, cbdef.Name, cbdef.Rect, opt)
		if err != nil {
			panic(err)
		}

		*form.Fields = append(*form.Fields, checkboxf.PdfField)
		page.AddAnnotation(checkboxf.Annotations[0].PdfAnnotation)
	}

	for _, chdef := range choiceFieldDefs {
		opt := annotator.ComboboxFieldOptions{Choices: chdef.Options}
		comboboxf, err := annotator.NewComboboxField(page, chdef.Name, chdef.Rect, opt)
		if err != nil {
			panic(err)
		}

		*form.Fields = append(*form.Fields, comboboxf.PdfField)
		page.AddAnnotation(comboboxf.Annotations[0].PdfAnnotation)
	}

	return form
}
