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
	msgDataChan := make(chan []byte, 1000)

	ctx := context.Background()
	session, err := client.NewAMQPSession(config)
	if err != nil {
		log.Fatalf("failed to connect to '%s': %v", config.GetBrokerURL(), err.Error())
	}
	defer session.Close(ctx)

	// async function to publish the messages on the server
	go func() {
		// wait for a new message to arrive
		addr := config.GetPublishAddress()
		sender, err := session.NewSender(amqp.LinkTargetAddress(addr))
		if err != nil {
			log.Fatalf("failed to create sender: %v", err)
		}
		defer sender.Close(ctx)
		for {
			msgData, ok := <-msgDataChan
			if !ok {
				log.Warnf("msg channel closed. Stopping the publish routine.")
				runtime.Goexit()
			}
			msg := amqp.Message{
				Header: &amqp.MessageHeader{
					DeliveryCount: 2,
					FirstAcquirer: true,
				},
				Data: [][]byte{msgData},
			}
			log.Infof("publishing msg '%s' on '%s'...", msg.GetData(), addr)
			err = sender.Send(ctx, &msg)
			if err != nil {
				log.Errorf("failed to publish msg '%s' to address '%s': %v", string(msg.GetData()), sender.Address(), err.Error())
			} else {
				log.Infof("published msg '%s' to address '%s'", string(msg.GetData()), sender.Address())
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
		msgDataChan <- []byte(data)
		count++
	}
}
