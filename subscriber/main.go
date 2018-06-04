package main

import (
	"context"

	"pack.ag/amqp"

	"github.com/fabric8-services/fabric8-amqp/client"
	"github.com/fabric8-services/fabric8-amqp/configuration"
	"github.com/fabric8-services/fabric8-amqp/log"
)

func main() {

	// loads the configuration
	config := configuration.New()

	ctx := context.Background()
	session, err := client.NewAMQPSession(config)
	if err != nil {
		log.Fatalf("failed to connect to '%s': %v", config.GetBrokerURL(), err.Error())
	}
	defer session.Close(ctx)

	receiver, err := session.NewReceiver(amqp.LinkAddress(config.GetTargetAddresses()[0]))
	if err != nil {
		log.Fatalf("failed to create receiver: %v", err)
	}

	for {
		msg, err := receiver.Receive(ctx)
		if err != nil {
			log.Errorf("failed to receive msg: %v", err.Error())
		} else {
			log.Infof("received msg: %v", string(msg.GetData()))
		}

	}
}
