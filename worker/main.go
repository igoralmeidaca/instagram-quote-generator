package main

import (
	"log"
	"os"

	"github.com/igoralmeidaca/quote"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the Temporal client with the specific namespace
	temporalHost := os.Getenv("TEMPORAL_HOST")
	c, err := client.Dial(client.Options{
		HostPort:  temporalHost,
		Namespace: quote.Namespace,
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// Create a worker for the task queue in the specified namespace
	w := worker.New(c, quote.GenerateTaskQueue, worker.Options{})

	// Register the workflow and activity functions
	w.RegisterWorkflow(quote.GenerateQuote)
	w.RegisterActivity(quote.GenerateText)
	w.RegisterActivity(quote.GenerateImage)
	w.RegisterActivity(quote.PublishInstagramPost)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
