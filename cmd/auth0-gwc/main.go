package main

import (
	"context"
	"log"

	"github.com/auth0/go-auth0/management"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	flag "github.com/spf13/pflag"
)

func main() {
	var (
		eventType string
		target    string
		source    string
		userEmail string
	)
	flag.StringVarP(&eventType, "event-type", "e", "", "event type [post-login | post-change-password | pre-user-registration | post-user-registration | sync]")
	flag.StringVarP(&userEmail, "user-email", "u", "", "user email")
	flag.StringVarP(&target, "target", "t", "http://auth0-gw.eastus.cloudapp.azure.com:8383/events", "target url")
	flag.StringVarP(&source, "source", "s", "auth0.com", "event source")
	flag.Parse()

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	user := management.User{
		Email: &userEmail,
	}

	// Create an Event.
	event := cloudevents.NewEvent()
	event.SetSource(source)
	event.SetType(eventType)
	if err := event.SetData(cloudevents.ApplicationJSON, user); err != nil {
		log.Println(err)
	}

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.Background(), target)

	// Send that Event.
	if result := c.Send(ctx, event); cloudevents.IsUndelivered(result) {
		log.Fatalf("failed to send, %v", result)
	} else {
		log.Printf("sent: %v", event)
		log.Printf("result: %v", result)
	}
}
