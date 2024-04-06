// package msg receives messages from Google Pub/Sub
package msg

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/siuyin/dflt"
)

var projectID = dflt.EnvString("PROJECT_ID", "someID")
var ctx = context.Background()
var client *pubsub.Client

func init() {
	cl, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	client = cl
}

func Close() {
	if err := client.Close(); err != nil {
		log.Printf("pubsub.Close: %v", err)
	}
}

func Get() string {
	return "hello"
}

func SetTopic(name string) error {

}
