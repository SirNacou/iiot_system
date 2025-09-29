package iotalerts

import (
	"fmt"
)

type AlertType int

// Use a slice to define the string representations in order.
var alertTypeNames = []string{
	"VIBRATION_EXCEEDED_THRESHOLD",
	"TEMPERATURE_CRITICAL",
	"SENSOR_OFFLINE",
}

// Define the enum values using iota.
const (
	VIBRATION_EXCEEDED_THRESHOLD AlertType = iota
	TEMPERATURE_CRITICAL
	SENSOR_OFFLINE
)

// The reverse map is generated automatically.
var alertTypeValues map[string]AlertType

func init() {
	alertTypeValues = make(map[string]AlertType, len(alertTypeNames))
	for i, name := range alertTypeNames {
		alertTypeValues[name] = AlertType(i)
	}
}

// String provides the string representation of the alert type.
func (at AlertType) String() string {
	if at >= 0 && int(at) < len(alertTypeNames) {
		return alertTypeNames[at]
	}
	return "Unknown Alert Type"
}

// AlertTypeFromString converts a string to an AlertType.
func AlertTypeFromString(s string) (AlertType, error) {
	if at, ok := alertTypeValues[s]; ok {
		return at, nil
	}
	return -1, fmt.Errorf("unknown alert type: %s", s)
}
