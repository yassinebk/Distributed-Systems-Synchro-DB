package ui

import (
	"fmt"
	"strconv"
	"synchro-db/db"
	"synchro-db/shared"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fatih/color"
)

func triggerError(message string, myWindow fyne.Window) {

	dialog.ShowError(fmt.Errorf(message), myWindow)
}

func showForm(app *fyne.App, productRepo *db.ProductSalesRepo) {
	newWindow := (*app).NewWindow("Add product form")

	var (
		productEntry = widget.NewEntry()
		dateEntry    = widget.NewEntry()
		regionEntry  = widget.NewEntry()
		qtyEntry     = widget.NewEntry()
		costEntry    = widget.NewEntry()
		taxEntry     = widget.NewEntry()
	)

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Product", Widget: productEntry},
			{Text: "Date", Widget: dateEntry},
			{Text: "Region", Widget: regionEntry},
			{Text: "Quantity", Widget: qtyEntry},
			{Text: "Cost", Widget: costEntry},
			{Text: "Tax", Widget: taxEntry},
		},

		OnSubmit: func() { // optional, handle form submission
			var (
				qty  uint32
				cost float32
				tax  float32
				x1   uint64
				x2   float64
				x3   float64
				err  error
			)

			if x1, err = strconv.ParseUint(qtyEntry.Text, 10, 32); err != nil {
				triggerError("Quantity should be a number", newWindow)
				return
			}

			if x2, err = strconv.ParseFloat(costEntry.Text, 32); err != nil {

				triggerError("Cost should be a decimal number", newWindow)
				return
			}

			if x3, err = strconv.ParseFloat(costEntry.Text, 32); err != nil {
				triggerError("Cost should be a decimal number", newWindow)
				return
			}

			qty = uint32(x1)
			cost = float32(x2)
			tax = float32(x3)
			time, err := time.Parse(time.DateOnly, dateEntry.Text)
			if err != nil {
				triggerError("Time should be dd-mm-yyyy"+err.Error(), newWindow)
				return
			}
			newProduct := &db.Product{

				Date:    time,
				Product: productEntry.Text,
				Region:  regionEntry.Text,
				Qty:     qty,
				Cost:    cost,
				Tax:     tax,
			}

			_, err = (*productRepo).CreateProduct(*newProduct)
			if err != nil {
				triggerError("Error creating new Product", newWindow)
				return
			}

			newWindow.Close()
		},
	}

	newWindow.SetContent(form)
	newWindow.Resize(fyne.NewSize(1000, 1000))
	newWindow.Show()

}

func CreateApp(whoami string, tableData *[]db.Product, productRepo *db.ProductSalesRepo) {

	myApp := app.New()

	myWindow := myApp.NewWindow("DB Synchronizer " + whoami)

	data := ConvertDataToDb(tableData)

	// Create a table widget for the main body
	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			x := widget.NewLabel("xx")
			return x
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		},
	)
	for id := range data[0] {
		table.SetColumnWidth(id, float32(120))
	}

	// Create a label widget for the sidebar
	sidebarLabel := widget.NewLabel("Welcome to DB Synchronizer")

	// Create a button widget for the sidebar
	sidebarButton := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		notSentProducts := (*productRepo).FindAllNotSent()

		// log.Println("Not sent products", notSentProducts)

		for i := range notSentProducts {
			notSentProducts[i].Sent = true
			(*productRepo).UpdateProduct(notSentProducts[i])
		}

		if whoami != "ho" && len(notSentProducts) > 0 {
			shared.SendProductsToHO(notSentProducts, whoami)
		}

	})

	// Create a VBox container for the sidebar
	sidebar := container.NewVBox(
		sidebarLabel,
		sidebarButton,
	)

	// Create a label widget for the main body
	mainLabel := widget.NewLabel("Welcome to your workspace - " + whoami)

	if whoami == "ho" {
		go shared.RecvDataFromTheWire("bo1")
		go shared.RecvDataFromTheWire("bo2")
		go shared.OrderProductIntoHeap()
		go shared.PerformDbOp(func() {
			green := color.New(color.FgHiGreen).SprintFunc()
			fmt.Println(green("[!] Updating UI with new data"))
			tempData := (*productRepo).FindAll()
			data = ConvertDataToDb(&tempData)

			// Create a table widget for the main body
			myWindow.Canvas().Refresh(table)

		})

	}
	// Create an HBox container for the main body
	tableCtr := container.NewVScroll(table)
	tableCtr.SetMinSize(fyne.NewSize(600, 500))

	addEntryButton := widget.NewButtonWithIcon("Add product", theme.ContentAddIcon(), func() { showForm(&myApp, productRepo) })
	mainBody := container.NewVBox(mainLabel, tableCtr, addEntryButton)

	// Create a Split container to hold the sidebar and main body
	split := container.NewHBox(sidebar, mainBody)

	// Set the content of the window to the Split container
	myWindow.SetContent(split)

	// Show the window

	// Blocking
	myWindow.ShowAndRun()

}
