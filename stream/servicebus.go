package stream

import (
	"log"
	"os"

	"github.com/michaelbironneau/asbclient"
)

var (
	namespace string = os.Getenv("SERVICE_BUS_NAMESPACE")
	keyValue  string = os.Getenv("SERVICE_BUS_SHARED_ACCESS_KEY_VALUE")
)

func SendToServiceBus(topic string, message string) {
	serviceBusClient := asbclient.New(asbclient.Topic, namespace, "RootManageSharedAccessKey", keyValue)
	err := serviceBusClient.Send(topic, &asbclient.Message{
		Body: []byte(message),
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Printf(`SERVICEBUS	: "%s" ---> "%s"`, message, topic)
	}
}
