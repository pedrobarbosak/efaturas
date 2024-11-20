package main

import (
	"context"
	"encoding/base64"
	"log"

	"efaturas-xtreme/pkg/efaturas"
)

func main() {
	uname, _ := base64.StdEncoding.DecodeString("MjQzNTY0MjEw")
	pword, _ := base64.StdEncoding.DecodeString("MjQzNTY0MjEwMTA=")

	ctx := context.Background()

	s := efaturas.New()
	cookies, err := s.Login(ctx, string(uname), string(pword))
	if err != nil {
		log.Panicln(err)
	}

	log.Println("🍪 cookies 🍪")
	for k, v := range cookies {
		log.Println(k, "-", v)
	}

	response, _, err := s.GetInvoices(ctx, cookies)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(response)

	invoice := response[0]

	isValid, err := s.CheckInvoice(ctx, cookies, invoice, invoice.Activity.Category)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("isValid:", isValid)
	return
}
