package input

import (
	"net/http"

	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	inputapi "github.com/KochKevin/effective-spoon-v2/internal/input/generated"
	"github.com/go-chi/render"
)

type InputService interface {
	EnterBarcode(input string) error
	EnterRfid(input string) error
}

type Api struct {
	InputService InputService
	Txm          infrastructure.TxManager
}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/input/generated.ServerInterface

// Use Inputsystem with barcode. Mainly for testing purposes
// (POST /input/barcode)
func (a *Api) PostInputBarcode(w http.ResponseWriter, r *http.Request, params inputapi.PostInputBarcodeParams) {

	a.InputService.EnterBarcode(params.Input)

	render.Status(r, http.StatusOK)
}

// Use Inputsystem with rfid. Mainly for testing purposes
// (POST /input/rfid)
func (a *Api) PostInputRfid(w http.ResponseWriter, r *http.Request, params inputapi.PostInputRfidParams) {

	a.InputService.EnterRfid(params.Input)

	render.Status(r, http.StatusOK)
}
