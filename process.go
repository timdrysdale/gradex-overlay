package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/timdrysdale/parsesvg"
	"github.com/timdrysdale/pool"
)

func actionExam(exam string, action string, pcChan chan int) error {

	root := fmt.Sprintf("./exams/%s/submitted/", exam)
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".pdf") {
			files = append(files, path)
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = doDir(files, action, exam, pcChan)

	return err
}

func doDir(inputPath []string, spreadName string, exam string, pcChan chan int) error {

	N := len(inputPath)

	tasks := []*pool.Task{}

	for i := 0; i < N; i++ {

		inputPDF := inputPath[i]
		spreadName := spreadName
		exam := exam
		newtask := pool.NewTask(func() error {
			pc, err := doOneDoc(inputPDF, spreadName, exam)
			pcChan <- pc
			return err
		})
		tasks = append(tasks, newtask)
	}

	p := pool.NewPool(tasks, runtime.GOMAXPROCS(-1))
	p.Run()

	var numErrors int
	for _, task := range p.Tasks {
		if task.Err != nil {
			fmt.Println(task.Err)
			numErrors++
		}
	}

	if numErrors > 0 {
		return errors.New("processing errors")
	}

	return nil

}

func doOneDoc(inputPath, spreadName, exam string) (int, error) {

	if strings.ToLower(filepath.Ext(inputPath)) != ".pdf" {
		return 0, errors.New(fmt.Sprintf("%s does not appear to be a pdf", inputPath))
	}

	// need page count to find the jpeg files again later
	numPages, err := countPages(inputPath)

	fmt.Printf("Starting: %s %d pages\n", inputPath, numPages)

	// render to image

	jpegPath := filepath.Join(filepath.Dir(inputPath), "/jpg")
	err = ensureDir(jpegPath)
	if err != nil {
		return 0, err
	}
	suffix := filepath.Ext(inputPath)

	basename := strings.TrimSuffix(filepath.Base(inputPath), suffix)
	jpegFileOption := fmt.Sprintf("%s/%s%%04d.jpg", jpegPath, basename)

	err = convertPDFToJPEGs(inputPath, jpegPath, jpegFileOption)
	if err != nil {
		return 0, err
	}

	// convert images to individual pdfs, with form overlay

	pagePath := "./pdf"
	err = ensureDir(pagePath)
	if err != nil {
		return 0, err
	}

	pageFileOption := fmt.Sprintf("%s/%s%%04d.pdf", pagePath, basename)

	mergePaths := []string{}

	// gs starts indexing at 1
	for imgIdx := 1; imgIdx <= numPages; imgIdx = imgIdx + 1 {

		// construct image name
		previousImagePath := fmt.Sprintf(jpegFileOption, imgIdx)
		pageFilename := fmt.Sprintf(pageFileOption, imgIdx)

		//TODO select Layout to suit landscape or portrait
		svgLayoutPath := "./test/layout-312pt-static-mark-dynamic-moderate-comment-static-check.svg"

		err := parsesvg.RenderSpread(svgLayoutPath, spreadName, previousImagePath, imgIdx, pageFilename)
		if err != nil {
			return 0, err
		}

		//save the pdf filename for the merge at the end
		mergePaths = append(mergePaths, pageFilename)
	}

	outPath := fmt.Sprintf("./exams/%s/%s", exam, spreadName)
	err = ensureDir(outPath)
	if err != nil {
		return 0, err
	}

	outputPath := fmt.Sprintf("%s/%s-%s.pdf", outPath, basename, spreadName)
	err = mergePdf(mergePaths, outputPath)
	if err != nil {
		return 0, err
	}

	return numPages, nil

}
