package main

import (
	"log"

	"efaturas-xtreme/server"

	"efaturas-xtreme/pkg/errors"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load .env:", err)
	}

	var config server.Config
	if _, err = env.UnmarshalFromEnviron(&config); err != nil {
		log.Panicln(errors.New(err))
	}

	s, err := server.New(config)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("exit:", s.Run())
}
