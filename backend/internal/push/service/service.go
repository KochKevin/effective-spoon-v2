package pushservice

import (
	"log/slog"

	pushapi "github.com/KochKevin/effective-spoon-v2/internal/push/generated"
)

type PushService struct {
	eventChannel chan string
}

func New() *PushService {
	return &PushService{
		eventChannel: make(chan string, 16),
	}
}

func (p *PushService) PushUserLogin() {
	p.push(string(pushapi.UserLogin))
}

func (p *PushService) PushShoppingCartUpdate() {
	p.push(string(pushapi.ShoppingcartUpdate))
}

func (p *PushService) push(event string) {
	select {
	case p.eventChannel <- event:
	default:
		slog.Warn("dropping push event, no listner or channel, event will be lost", "event", event)
	}
}

func (p *PushService) GetEventChannel() <-chan string {
	return p.eventChannel
}
