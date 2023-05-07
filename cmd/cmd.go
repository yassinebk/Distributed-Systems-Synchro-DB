package cmd

import (
	"fmt"
	"log"
	"os"
	"synchro-db/db"
	"synchro-db/ui"

	"github.com/urfave/cli/v2" // imports as package "cli"

	"github.com/fatih/color"
)

func printASCIIArt(*cli.Context) error {

	content, err := os.ReadFile("./ascii.txt")
	if err != nil {
		log.Panicln("[-] Error reading ascii art")
	}

	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println(cyan(string(content)))
	return nil
}

func Setup() cli.App {

	return cli.App{
		Usage:       "Welcome to DB Syhncronizer, use wisely to synchronize your db between HO and BO",
		Description: "Distributed Systems final project, @2023 3rd SE INSAT ",
		Authors: []*cli.Author{
			{Name: "Mahmoud Hedi Nefzi", Email: "mahmoudhnefzi@gmail.com"},
			{Name: "Yassine Belkhadem", Email: "yassine.belkhadem@insat.ucar.tn"},
		},
		Name:     "SyncDB-Super!",
		Before:   printASCIIArt,
		Commands: []*cli.Command{},
		Action: func(ctx *cli.Context) error {
			err := db.SeedDB("test-db.sqlite")
			if err != nil {
				return err
			}

			dbConnection, err := db.ConnectToDb("test-db.sqlite")

			productsRepo := db.NewProductSalesRepo(dbConnection)

			products := productsRepo.FindAll()

			ui.CreateApp("Default app", &products, &productsRepo)
			return nil
		},
	}

}
