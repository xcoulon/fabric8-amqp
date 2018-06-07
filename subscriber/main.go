package main

import (
	"context"
	"strings"

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

	receiver, err := session.NewReceiver(amqp.LinkSourceAddress(config.GetTargetAddresses()[0]))
	if err != nil {
		log.Fatalf("failed to create receiver: %v", err)
	}

	for {
		msg, err := receiver.Receive(ctx)
		msgData := string(msg.GetData())
		if err != nil {
			log.Errorf("failed to receive msg: %v", err.Error())
		} else {
			// reject messages ending with `0` or `5`
			if strings.HasSuffix(msgData, "0") || strings.HasSuffix(msgData, "5") {
				log.Warnf("rejected message '%s'", msgData)
				msg.Modify(true, true, amqp.Annotations{})
			} else {
				log.Infof("accepted message '%s'", msgData)
				msg.Accept()
			}
		}
	}
}
