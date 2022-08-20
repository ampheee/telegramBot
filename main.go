package main

import (
	"flag"
	client "github.com/ampheee/telegramBot/v2/clients/telegram"
	consEvents "github.com/ampheee/telegramBot/v2/consumer/cons-events"
	"github.com/ampheee/telegramBot/v2/events/telegram"
	"github.com/ampheee/telegramBot/v2/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "userpages-storage"
	batchSize   = 150
)

func main() {
	eventsProcessor := telegram.New(
		client.NewClient(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := consEvents.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"token output",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("token is required")
	}
	return *token
}
