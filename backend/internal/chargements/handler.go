package chargements

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	chargementsapi "github.com/KochKevin/effective-spoon-v2/internal/chargements/generated"
	"github.com/go-chi/render"
)

type Service interface {
	CreatePaymentLink() (string, error)
}

type Api struct {
	Service Service
}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/chargements/generated.ServerInterface

// PostAddBalance Add new balance to the currently logged in user, using stripe
// (POST /add-balance)
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
	*/
	

	url, err := a.Service.CreatePaymentLink()
	if err != nil {
		log.Fatalf("error calling stripe: %v", err)
	}

	render.JSON(w, r, chargementsapi.AddBalanceResponse{PaymentLink: url})
}
