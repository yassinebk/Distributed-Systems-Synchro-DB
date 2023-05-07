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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "seed",
				Value: false,
				Usage: "Seed database or not",
			},
			&cli.StringFlag{
				Name:     "whoami",
				Value:    "OH",
				Usage:    "Specify which branch you are OH,BH1,BH2",
				Required: true,
			},
		},

		Action: func(ctx *cli.Context) error {

			if ctx.Bool("seed") {
				err := db.SeedDB("test-db.sqlite")
				if err != nil {
					log.Panicln("[-] Error seeding database", err)
				}
			}

			dbConnection, err := db.ConnectToDb("test-db.sqlite")

			if err != nil {
				log.Panicln("[-] Error connecting to database", err)
			}

			productsRepo := db.NewProductSalesRepo(dbConnection)

			products := productsRepo.FindAll()

			ui.CreateApp(ctx.String("whoami"), &products, &productsRepo)
			return nil
		},
	}

}
