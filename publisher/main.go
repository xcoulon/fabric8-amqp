package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/fabric8-services/fabric8-amqp/client"
	"github.com/fabric8-services/fabric8-amqp/configuration"
	"github.com/fabric8-services/fabric8-amqp/log"
	"pack.ag/amqp"
)

func main() {

	// loads the configuration
	config := configuration.New()

	// message channel with enough capacity to handle a disconnection (hopefully...)
	msgChan := make(chan amqp.Message, 1000)

	ctx := context.Background()
	session, err := client.NewAMQPSession(config)
	if err != nil {
		log.Fatalf("failed to connect to '%s': %v", config.GetBrokerURL(), err.Error())
	}
	defer session.Close(ctx)

	// async function to publish the messages on the server
	go func() {
		// wait for a new message to arrive
		targetAddresses := config.GetTargetAddresses()
		for {
			msg, ok := <-msgChan
			if !ok {
				log.Warnf("msg channel closed. Stopping the publish routine.")
				runtime.Goexit()
			}
			log.Infof("publishing msg '%s'...", msg.GetData())
			for _, addr := range targetAddresses {
				sender, err := session.NewSender(amqp.LinkTargetAddress(addr))
				if err != nil {
					log.Fatalf("failed to create sender: %v", err)
				}
				err = sender.Send(ctx, &msg)
				if err != nil {
					log.Errorf("failed to publish msg '%s' to address '%s': %v", string(msg.GetData()), addr, err.Error())
				} else {
					log.Infof("published msg '%s' to address '%s'", string(msg.GetData()), addr)
				}
			}
		}
	}()

	// infinite loop of message publishing...
	count := 1
	for {
		// block for a few seconds...
		<-time.After(3 * time.Second)
		data := fmt.Sprintf("message #%d", count)
		log.Infof("preparing msg '%s", data)
		msg := amqp.NewMessage([]byte(data))
		msgChan <- *msg
		count++
	}
}
