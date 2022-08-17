package main

import (
	"flag"
	"github.com/ampheee/telegramBot/v2/clients/telegram"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.NewClient(tgBotHost, mustToken())
	// fetcher = fetcher.NewFetcher()
	// processor = processor.NewProcessor() // will talk with telegramBot api
	//consumer.Start(fetcher, processor) // with telegramBOT api. Fetcher to fetch reqs, processor
	// to process reqs and send resps
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
