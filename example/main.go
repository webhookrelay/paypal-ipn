package main

import (
	"log"
	"net/http"

	ipn "github.com/webhookrelay/paypal-ipn"
)

const (
	debug = false
)

func main() {
	mux := http.NewServeMux()

	listener := ipn.New(debug)

	mux.Handle("/v1/paypal-webhooks", listener.WebhooksHandler(func(err error, n *ipn.PaypalNotification) {
		if err != nil {
			log.Printf("IPN error: %v", err)
			return
		}

		log.Printf("event type: %s", n.EventType)
		log.Printf("event resource type: %s", n.ResourceType)
		log.Printf("summary: %s", n.Summary)
	}))
	log.Println("server starting on :8080")
	log.Fatalf("failed to run http server: %v", http.ListenAndServe(":8080", mux))
}
