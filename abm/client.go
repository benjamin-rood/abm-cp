package abm

import "time"

// Client is the wrapper interface between the different comm protocols with the 'user' / 'caller' and the Model.
type Client interface {
	Monitor(chan struct{})
	Dead() bool
	TimeStamp() time.Time
}
