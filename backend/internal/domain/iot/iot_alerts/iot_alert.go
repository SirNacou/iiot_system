package iotalerts

import (
	"time"
)

type IotAlert struct {
	Time         time.Time
	DeviceID     string
	AlertType    string
	Severity     string
	Message      string
	CurrentValue int
}
