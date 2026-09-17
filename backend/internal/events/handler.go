package events

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"

	eventsapi "github.com/KochKevin/effective-spoon-v2/internal/events/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
)

type Service interface {
	CreateCurrentEvent(ctx context.Context, authorId uuid.UUID, amountPerPerson int, endDateTime time.Time) (Event, error)
	GetEventUsage(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (EventUsage, error)
	GetCurrentEvent(ctx context.Context) (Event, error)
}

type UserService interface {
	GetUser(ctx context.Context, userId uuid.UUID) (users.User, error)
}

type Api struct {
	Service     Service
	UserService UserService
}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/events/generated.ServerInterface

// // GetEventsCurrent Get the current event and the event usage of the current user
// (GET /events/current)
func (a *Api) GetEventsCurrent(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	event, err := a.Service.GetCurrentEvent(r.Context())

	if err != nil {
		slog.Error("getting event", "err", err)

		if errors.Is(err, NoCurrentEventErr) {
			render.Status(r, http.StatusInternalServerError)
			http.Error(w, "Not Found", http.StatusNotFound)
		} else {
			render.Status(r, http.StatusInternalServerError)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	author, err := a.UserService.GetUser(r.Context(), event.UserId)
	if err != nil {
		slog.Error("getting user:", "err", err)
		render.Status(r, http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	eventUsage, err := a.Service.GetEventUsage(r.Context(), userID, event.Id)
	if err != nil {
		slog.Error("getting event usage:", "err", err)
		render.Status(r, http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, eventsapi.Event{
		AuthorName:               author.Name,
		AvailableAmountPerPerson: event.AmountFreeProductsPerUser,
		EventEndDateTime:         event.EndTimestamp,
		EventUsage: eventsapi.EventUsage{
			AmountUsed: eventUsage.AmountUsed,
			UserId:     eventUsage.UserId.String(),
		}})

}

// PostEventsCurrent Create a new event and set it as the current one
// (POST /events/current)
func (a *Api) PostEventsCurrent(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("reading request body", "err", err)
	}
	defer r.Body.Close()

	createEvent := eventsapi.CreateEvent{}
	err = json.Unmarshal(bodyBytes, &createEvent)
	if err != nil {
		slog.Error("unmarshaling create event request body", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	event, err := a.Service.CreateCurrentEvent(r.Context(), userID, createEvent.AmountPerPerson, createEvent.EndDateTime)
	if err != nil {
		slog.Error("creating event", "err", err)
		render.Status(r, http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	author, err := a.UserService.GetUser(r.Context(), event.UserId)
	if err != nil {
		slog.Error("getting user:", "err", err)
		render.Status(r, http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	eventUsage, err := a.Service.GetEventUsage(r.Context(), userID, event.Id)
	if err != nil {
		slog.Error("getting event usage:", "err", err)
		render.Status(r, http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, eventsapi.Event{
		AuthorName:               author.Name,
		AvailableAmountPerPerson: event.AmountFreeProductsPerUser,
		EventEndDateTime:         event.EndTimestamp,
		EventUsage: eventsapi.EventUsage{
			AmountUsed: eventUsage.AmountUsed,
			UserId:     eventUsage.UserId.String(),
		}})

}

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
*/
