package mandrill

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/common/auth"
)

type Webhook struct {
	Path string
	acc  telegraf.Accumulator
	log  telegraf.Logger
	auth.BasicAuth
}

// Register registers the webhook with the provided router
func (md *Webhook) Register(router *mux.Router, acc telegraf.Accumulator, log telegraf.Logger) {
	router.HandleFunc(md.Path, returnOK).Methods("HEAD")
	router.HandleFunc(md.Path, md.eventHandler).Methods("POST")

	md.log = log
	md.log.Infof("Started the webhooks_mandrill on %s", md.Path)
	md.acc = acc
}

func returnOK(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (md *Webhook) eventHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !md.Verify(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := url.ParseQuery(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var events []mandrillEvent
	err = json.Unmarshal([]byte(data.Get("mandrill_events")), &events)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, event := range events {
		md.acc.AddFields("mandrill_webhooks", event.fields(), event.tags(), time.Unix(event.TimeStamp, 0))
	}

	w.WriteHeader(http.StatusOK)
}
