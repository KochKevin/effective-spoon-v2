package pushservice

import pushapi "github.com/KochKevin/effective-spoon-v2/internal/push/generated"


type PushService struct {
	eventChannel chan string
}

func New() *PushService {
	return &PushService{
		eventChannel: make(chan string),
	}
}

func (p *PushService) PushUserLogin() {
	p.eventChannel <- string(pushapi.UserLogin)
}

func (p *PushService) PushShoppingCartUpdate() {
	p.eventChannel <- string(pushapi.ShoppingcartUpdate)
}

func (p *PushService) GetEventChannel() <-chan string {
	return p.eventChannel
}
