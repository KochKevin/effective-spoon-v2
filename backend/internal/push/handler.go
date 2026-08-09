package push

import (
	"log/slog"
	"net/http"
	"time"

	pushapi "github.com/KochKevin/effective-spoon-v2/internal/push/generated"
)

type PushService interface {
	GetEventChannel() <-chan string
}

type Api struct {
	pushapi.Unimplemented
	PushService PushService
}

const endOfLine = "\n"
const lineSeperator = "\n"

// This system is currently limited to one connection. It would break if multiple clients would be connected. But for its current purpose its enough
func (a *Api) GetPushes(w http.ResponseWriter, r *http.Request) {

	slog.Debug("GET on /pushes")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(30 * time.Second)
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

			w.Write([]byte("data: " + event + endOfLine + lineSeperator))

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

}
