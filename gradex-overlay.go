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
	"fmt"
	"os"

	unicommon "github.com/unidoc/unipdf/v3/common"
)

func init() {
	// Debug log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelInfo))
}

func main() {

	c := make(chan int, 1000)

	closed := make(chan struct{})

	err := actionExam("ENGI01020", "mark", c)

	go func() {
		scripts := 0
		pages := 0
	LOOP:
		for {
			select {
			case <-closed:
				break LOOP
			case val := <-c:
				scripts++
				pages = pages + val
				fmt.Printf("%03d / %04d\n", scripts, pages)
			}

		}
	}()

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

}

/*
	a := app.New()

	//closed := make(chan struct{})
	//c := make(chan float64)
	rand.Seed(time.Now().UnixNano())

	//gridSize := fyne.NewSize(100+theme.Padding(), 100+theme.Padding())
	//cellSize := fyne.NewSize(50, 50)

	w := a.NewWindow("Hello")

	tc := widget.NewTabContainer()

	//entry := widget.NewEntry()

	examLabel := widget.NewLabel("ELEE09442 - Special Topics in Counting")

	submitCount := widget.NewEntry()
	submitCount.SetText("0")
	submitButton := widget.NewButton("Count submitted scripts", func() { submitCount.SetText("10") })

	markCount := widget.NewEntry()
	markCount.SetText("0/10")
	markButton := widget.NewButton("Add marking sidebar", func() { markCount.SetText("10/10") })

	markedCount := widget.NewEntry()
	markedCount.SetText("0/10")
	markedButton := widget.NewButton("Count marked scripts", func() { markedCount.SetText("10/10") })

	moderateCount := widget.NewEntry()
	moderateCount.SetText("0/10")
	moderateButton := widget.NewButton("Add moderating sidebar", func() { moderateCount.SetText("10/10") })

	checkCount := widget.NewEntry()
	checkCount.SetText("0/10")
	checkButton := widget.NewButton("Add checking sidebar", func() { checkCount.SetText("10/10") })

	moderatedCount := widget.NewEntry()
	moderatedCount.SetText("0/10")
	moderatedButton := widget.NewButton("Count moderated scripts", func() { moderatedCount.SetText("10/10") })

	checkedCount := widget.NewEntry()
	checkedCount.SetText("0/10")
	checkedButton := widget.NewButton("Count checked scripts", func() { checkedCount.SetText("10/10") })

	tab1 := widget.NewTabItem("ELEE09442", fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		examLabel,
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			submitButton,
			submitCount,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			markButton,
			markCount,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			markedButton,
			markedCount,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			moderateButton,
			moderateCount,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			moderatedButton,
			moderatedCount,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			checkButton,
			checkCount,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			checkedButton,
			checkedCount,
		),
	))
	examLabel2 := widget.NewLabel("ENGI12345 - Things to do with spanners")
	submitCount2 := widget.NewEntry()
	submitCount2.SetText("0")
	submitButton2 := widget.NewButton("Count submitted scripts", func() { submitCount2.SetText("10") })

	markCount2 := widget.NewEntry()
	markCount2.SetText("0/10")
	markButton2 := widget.NewButton("Add marking sidebar", func() { markCount2.SetText("10/10") })

	markedCount2 := widget.NewEntry()
	markedCount2.SetText("0/10")
	markedButton2 := widget.NewButton("Count marked scripts", func() { markedCount2.SetText("10/10") })

	moderateCount2 := widget.NewEntry()
	moderateCount2.SetText("0/10")
	moderateButton2 := widget.NewButton("Add moderating sidebar", func() { moderateCount2.SetText("10/10") })

	checkCount2 := widget.NewEntry()
	checkCount2.SetText("0/10")
	checkButton2 := widget.NewButton("Add checking sidebar", func() { checkCount2.SetText("10/10") })

	moderatedCount2 := widget.NewEntry()
	moderatedCount2.SetText("0/10")
	moderatedButton2 := widget.NewButton("Count moderated scripts", func() { moderatedCount2.SetText("10/10") })

	checkedCount2 := widget.NewEntry()
	checkedCount2.SetText("0/10")
	checkedButton2 := widget.NewButton("Count checked scripts", func() { checkedCount2.SetText("10/10") })
	tab2 := widget.NewTabItem("ENGI12345", fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		examLabel2,
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			submitButton2,
			submitCount2,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			markButton2,
			markCount2,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			markedButton2,
			markedCount2,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			moderateButton2,
			moderateCount2,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			moderatedButton2,
			moderatedCount2,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			checkButton2,
			checkCount2,
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			checkedButton2,
			checkedCount2,
		),
	))

	tc.Append(tab1)
	tc.Append(tab2)

	w.SetContent(tc)

	/*go func() {

		total := int(0)
	LOOP:
		for {
			select {
			case <-c:
				total++
				entry.SetText(fmt.Sprintf("%d", total))
				entry.Refresh()
			case <-closed:
				break LOOP
			}
		}

	}()
*/

/*
	w.ShowAndRun()
}
*/
