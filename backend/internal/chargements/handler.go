package chargements

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"

	chargementsapi "github.com/KochKevin/effective-spoon-v2/internal/chargements/generated"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Service interface {
	CreateChargementIntent(ctx context.Context, userId uuid.UUID, amount float32) (ChargementIntent, error)
	CancelCurrentChargementIntent(ctx context.Context) error
}

type Api struct {
	Service Service
}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/chargements/generated.ServerInterface

// PostAddBalance Add new balance to the currently logged in user, using stripe
// (POST /add-balance)
/*
func (a *Api) PostAddBalance(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Fehler beim Lesen des Bodys: %v", err)
	}
	defer r.Body.Close()

	addBalance := chargementsapi.PostAddBalanceJSONRequestBody{}
	err = json.Unmarshal(bodyBytes, &addBalance)
	if err != nil {
		log.Fatalf("Fehler beim Unmarshal: %v", err)
	}
	/*
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			err := errors.New("cannot get user_id from context")
			slog.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
}



	url, err := a.Service.CreatePaymentLink()
	if err != nil {
		log.Fatalf("error calling stripe: %v", err)
	}

	render.JSON(w, r, chargementsapi.AddBalanceResponse{PaymentLink: url})
}
*/

// PostChargementsCurrent Create a new chargement and set it as the current one
// (POST /chargements/current)
func (a *Api) PostChargementsCurrent(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading body: %v", err)
	}
	defer r.Body.Close()

	createChargement := chargementsapi.CreateChargement{}
	err = json.Unmarshal(bodyBytes, &createChargement)
	if err != nil {
		log.Fatalf("error unmarshaling json to struct: %v", err)
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	chargementIntent, err := a.Service.CreateChargementIntent(r.Context(), userID, createChargement.AmountToAdd)
	if err != nil {
		slog.Error("error creating chargement intent", "err", err)
		render.Status(r, http.StatusInternalServerError)
	}

	render.JSON(w, r, chargementsapi.Chargement{PaymentLink: chargementIntent.PaymentLink, Amount: chargementIntent.Amount.GetAsEuro()})

}

// PostChargementsCurrentCancel Cancel current Chargement
// (POST /chargements/current/cancel)
func (a *Api) PostChargementsCurrentCancel(w http.ResponseWriter, r *http.Request) {

	err := a.Service.CancelCurrentChargementIntent(r.Context())
	if err != nil {
		slog.Error("error on cancelling current chargement intent", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
