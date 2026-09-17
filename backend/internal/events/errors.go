package events

import "errors"

var NoCurrentEventErr = errors.New("the is no current event. Create a new one")
