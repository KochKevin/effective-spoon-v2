package input

import (
	"log/slog"
	"net/http"

	inputapi "github.com/KochKevin/effective-spoon-v2/internal/input/generated"
	"github.com/go-chi/render"
)

type InputService interface {
	EnterBarcode(input string) error
	EnterRfid(input string) error
}

type Api struct {
	InputService InputService
}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/input/generated.ServerInterface

//Barcode
//curl -X POST "https://effective-waddle-4j7xj7vp44pxh7xrp-8080.app.github.dev/api/input/barcode?input=1"
//RFID
//curl -X POST "https://effective-waddle-4j7xj7vp44pxh7xrp-8080.app.github.dev/api/input/rfid?input=1"

// Use Inputsystem with barcode. Mainly for testing purposes
// (POST /input/barcode)
func (a *Api) PostInputBarcode(w http.ResponseWriter, r *http.Request, params inputapi.PostInputBarcodeParams) {

	err := a.InputService.EnterBarcode(params.Input)
	if err != nil {
		slog.Error("error while trying enter input with barcode", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
}

// Use Inputsystem with rfid. Mainly for testing purposes
// (POST /input/rfid)
func (a *Api) PostInputRfid(w http.ResponseWriter, r *http.Request, params inputapi.PostInputRfidParams) {

	a.InputService.EnterRfid(params.Input)

	render.Status(r, http.StatusOK)
}
