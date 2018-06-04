package client

import (
	"github.com/fabric8-services/fabric8-amqp/configuration"
	"github.com/fabric8-services/fabric8-amqp/log"
	"pack.ag/amqp"
)

// NewAMQPSession initializes a new AMQP session using the given configuration settings.
func NewAMQPSession(config configuration.Config) (*amqp.Session, error) {
	log.Infof("opening connection to broker on '%s'...", config.GetBrokerURL())
	client, err := amqp.Dial(config.GetBrokerURL(),
		// amqp.ConnSASLPlain(config.GetUsername(), config.GetPassword()),
		amqp.ConnSASLAnonymous(),
		// amqp.ConnProperty("key"+numStr, "value"+numStr),
		// amqp.ConnTLSConfig(&tls.Config{
		// InsecureSkipVerify: true,
		// }),
	)
	if err != nil {
		log.Fatalf("error while dialing the AMQP server at %s: %v", config.GetBrokerURL(), err)
	}

	log.Infof("connection established with server at '%s'", config.GetBrokerURL())
	// Open a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("error while creating an AMQP session: %v", err)
		return nil, err
	}
	log.Infof("obtained a new session at '%s'", config.GetBrokerURL())
	return session, nil
}
