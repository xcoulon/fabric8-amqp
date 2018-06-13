package main

import (
	"github.com/fabric8-services/fabric8-amqp/configuration"
	"github.com/fabric8-services/fabric8-amqp/log"
	"qpid.apache.org/electron"
)

func main() {

	// loads the configuration
	config := configuration.New()

	// message channel with enough capacity to handle a disconnection (hopefully...)
	container := electron.NewContainer(config.GetPodName())
	c, err := container.Dial("tcp", config.GetBrokerURL())
	if err != nil {
		log.Fatalf("failed to connect to '%s': %v", config.GetBrokerURL(), err.Error())
	}
	defer c.Close(nil)
	r, err := c.Receiver(electron.Source(config.GetQueueName()))
	log.Infof("opened connection to `%s` and initialized receiver to `%s`", config.GetBrokerURL(), config.GetQueueName())
	if err != nil {
		log.Fatalf("failed to initialize the received to '%s' on : %v", config.GetQueueName(), config.GetBrokerURL(), err.Error())
	}

	for {
		msg, err := r.Receive()
		if err != nil {
			log.Fatalf("failed to receive msg: %v", err.Error())
		} else {
			log.Infof("received msg: %v", msg.Message.Body())
			msg.Accept()
		}
	}
}
