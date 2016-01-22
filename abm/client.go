package abm

import "time"

// Client defines the owner of the Model.
type Client struct {
	UID string
	Model
	Active bool
	Stamp  time.Time
}
