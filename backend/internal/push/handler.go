package push

import (
	"log/slog"
	"net/http"
	"time"

	pushservice "github.com/KochKevin/effective-spoon-v2/internal/push/service"
	"github.com/go-chi/chi/v5"
)

type PushService interface {
	GetEventChannel() <-chan string
}

type Api struct {
	PushService PushService
}

const endOfLine = "\n"
const lineSeperator = "\n"


//This system is currently limited to one connection. It would break if multiple clients would be connected. But for its current purpose its enough

func (a *Api) GetEvents(r chi.Router) {

	r.Get("/events", func(w http.ResponseWriter, r *http.Request) {

		slog.Debug("GET on /events")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {

			//DO WORK
			//slog.Debug("Hallo!")

			w.Write([]byte(": heartbeat" + endOfLine + lineSeperator))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

			select {
			//rerun work
			case <-ticker.C:
				continue

			case event := <-a.PushService.GetEventChannel():

				if event == pushservice.ShoppingCartUpdateEvent {

					w.Write([]byte("data: " + event + endOfLine + lineSeperator))

				} else if event == pushservice.UserLoginEvent {
					w.Write([]byte("data: " + event + endOfLine + lineSeperator))
				}

				//Push new message
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
				continue

			//End Ticker when connection is closed
			case <-r.Context().Done():
				return
			}
		}

	})

}
