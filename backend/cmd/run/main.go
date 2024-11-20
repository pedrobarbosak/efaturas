package main

import (
	"context"
	"log"

	"efaturas-xtreme/internal/service"
	"efaturas-xtreme/internal/service/repository"
	"efaturas-xtreme/pkg/db"
	"efaturas-xtreme/pkg/efaturas"
	"efaturas-xtreme/pkg/errors"
	"efaturas-xtreme/pkg/sse"
	"efaturas-xtreme/server"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load .env:", err)
	}

	var cfg server.Config
	if _, err = env.UnmarshalFromEnviron(&cfg); err != nil {
		log.Panicln(errors.New(err))
	}

	db, err := db.New(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		log.Panicln("failed to create db:", err)
	}

	sse := sse.New()
	efaturas := efaturas.New()

	repo, err := repository.New(db)
	if err != nil {
		log.Panicln("failed to create repo:", err)
	}

	service := service.New(repo, efaturas, sse)

	//////////////////////////////////////////////////

	ctx := context.Background()
	uname := "227130782"
	pword := "PHTXUQLMTTGL"

	invoices, cookies, err := service.CreateOrUpdate(ctx, uname, pword)
	if err != nil {
		log.Panicln(err)
	}

	for _, inv := range invoices {
		log.Println(inv.ID, inv.Document.Number, inv.Document.Description)
	}

	//////////////////////////////////////////////////

	log.Println(service.ScanInvoices(ctx, cookies, invoices))
}
