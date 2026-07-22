package pushservice

const UserLoginEvent string = "user.login"
const ShoppingCartUpdateEvent string = "shoppingcart.update"

type PushService struct {
	eventChannel chan string
}

func New() *PushService {
	return &PushService{
		eventChannel: make(chan string),
	}
}

func (p *PushService) PushUserLogin() {
	p.eventChannel <- UserLoginEvent
}

func (p *PushService) PushShoppingCartUpdate() {
	p.eventChannel <- ShoppingCartUpdateEvent
}

func (p *PushService) GetEventChannel() <-chan string {
	return p.eventChannel
}
