// package msg receives messages from Google Pub/Sub
package msg

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/siuyin/dflt"
	"google.golang.org/api/iterator"
)

var projectID = dflt.EnvString("PROJECT_ID", "someID")
var ctx = context.Background()
var client *pubsub.Client
var Sub *pubsub.Subscription

// err := Sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
	// fmt.Printf("ger: %s", m.Data)
// })
// if err != nil {
	// log.Printf("msg.Get: %v", err)
// }

const subName = "backend"

func init() {
	cl, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	client = cl
	Sub = client.Subscription(subName)
}

func Close() {
	if err := client.Close(); err != nil {
		log.Printf("pubsub.Close: %v", err)
	}
}


func SetTopic(name string) error {
	return fmt.Errorf("not implemented")
}

func ListTopics() []*pubsub.Topic {
	topics := []*pubsub.Topic{}
	iter := client.Topics(ctx)
	for {
		t, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("msg.ListTopics: %v", err)
		}
		topics = append(topics, t)
	}
	return topics
}
