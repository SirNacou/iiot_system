package iotalerts

type IotAlertRepository interface {
	FindBySpec() []*IotAlert
}
