package iot

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type OeeAlertsConsumer struct {
	logger watermill.LoggerAdapter
}

func NewOeeAlertConsumer(logger watermill.LoggerAdapter) *OeeAlertsConsumer {
	return &OeeAlertsConsumer{
		logger: logger,
	}
}

func (c OeeAlertsConsumer) Handle(msg *message.Message) error {
	c.logger.Info(fmt.Sprintf("message received. %v", msg.Payload), watermill.LogFields{})
	return nil
}
