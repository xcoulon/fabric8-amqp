package main

import (
	"fmt"
	"time"

	"github.com/fabric8-services/fabric8-amqp/configuration"
	"github.com/fabric8-services/fabric8-amqp/log"
	"qpid.apache.org/amqp"
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
	s, err := c.Sender(electron.Target(config.GetPublishAddress()))
	log.Infof("opened connection to `%s` and initialized sender to `%s`", config.GetBrokerURL(), config.GetPublishAddress())
	if err != nil {
		log.Fatalf("failed to initialize the sender to '%s': %v", config.GetPublishAddress(), err.Error())
	}

	// infinite loop of message publishing...
	count := 1
	for {
		// block for a few seconds...
		<-time.After(3 * time.Second)
		t := time.Now()
		data := fmt.Sprintf("message #%d (%v)", count, t.Format("2006-01-02 15:04:05"))
		msg := amqp.NewMessageWith(data)
		s.SendSync(msg)
		log.Infof("sent message `%s`", data)
		count++
	}
}
