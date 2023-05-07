package ui

import (
	"fmt"
	"synchro-db/db"
)

func ConvertDataToDb(products *[]db.Product) [][]string {

	data := make([][]string, len(*products)+1)
	data[0] = []string{"Date", "Region", "Product", "Qty", "Cost", "Amt", "Tax", "Total"}
	for i, p := range *products {
		data[i+1] = []string{
			fmt.Sprintf("%d-%d-%d",p.Date.Day(),p.Date.Month(),p.Date.Year()),
			p.Region,
			p.Product,
			fmt.Sprintf("%d", p.Qty),
			fmt.Sprintf("%.2f", p.Cost),
			fmt.Sprintf("%.2f", p.Cost*float32(p.Qty)),
			fmt.Sprintf("%.2f", p.Tax),
			fmt.Sprintf("%.2f", p.Tax+p.Cost*float32(p.Qty)),
		}
	}

	return data
}
