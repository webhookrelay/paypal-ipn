package ipn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Listener struct {
	debug bool
}

func New(debug bool) *Listener {
	return &Listener{
		debug: debug,
	}
}

// Listen for webhooks
func (l *Listener) WebhooksHandler(cb func(err error, n *PaypalNotification)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			cb(fmt.Errorf("failed to read body: %s", err), nil)
			return
		}

		var notification PaypalNotification
		err = json.Unmarshal(body, &notification)
		if err != nil {
			cb(fmt.Errorf("failed to decode request body: %s", err), nil)
			return
		}

		if l.debug {
			fmt.Printf("paypal: body: %s, parsed: %+v\n", body, notification)
		}

		w.WriteHeader(http.StatusOK)
		cb(nil, &notification)
	}
}

type PaypalNotification struct {
	ID           string    `json:"id"`
	CreateTime   time.Time `json:"create_time"`
	ResourceType string    `json:"resource_type"`
	EventType    string    `json:"event_type"`
	Summary      string    `json:"summary"`
	Resource     struct {
		ParentPayment string    `json:"parent_payment"`
		UpdateTime    time.Time `json:"update_time"`
		Amount        struct {
			Total    string `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
		CreateTime time.Time `json:"create_time"`
		Links      []struct {
			Href   string `json:"href"`
			Rel    string `json:"rel"`
			Method string `json:"method"`
		} `json:"links"`
		ID    string `json:"id"`
		State string `json:"state"`
	} `json:"resource"`
	Links []struct {
		Href    string `json:"href"`
		Rel     string `json:"rel"`
		Method  string `json:"method"`
		EncType string `json:"encType"`
	} `json:"links"`
	EventVersion string `json:"event_version"`
}
