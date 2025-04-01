package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/igoralmeidaca/quote"
)

func main() {
	// Create a Temporal namespace client
	temporalHost := os.Getenv("TEMPORAL_HOST")
	namespaceClient, err := client.NewNamespaceClient(client.Options{
		HostPort: temporalHost,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal namespace client:", err)
	}

	// Define the retention period (30 days) as *durationpb.Duration
	retention := durationpb.New(30 * 24 * time.Hour) // 30 days

	// Register the namespace
	err = namespaceClient.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        quote.Namespace,
		WorkflowExecutionRetentionPeriod: retention,
	})

	if err != nil {
		log.Print("Namespace already registered:", err)
	} else {
		log.Println("Namespace 'Instagram Quote Generator' registered successfully.")
	}

	// Close the namespace client
	namespaceClient.Close()

	// Create the Temporal client for workflows
	c, err := client.Dial(client.Options{
		Namespace: quote.Namespace, // Use the newly registered namespace
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	//scheduleWorkflow(c)
}

func scheduleWorkflow(c client.Client) {
	// Define the schedule ID
	scheduleID := "schedule_quote"

	// Create the schedule
	_, err := c.ScheduleClient().Create(context.Background(), client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			Intervals: []client.ScheduleIntervalSpec{
				{Every: 60 * time.Second},
			},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        uuid.New().String(),
			Workflow:  quote.GenerateQuote,
			TaskQueue: quote.GenerateTaskQueue,
			Args:      []any{quote.GenerateTextInput{}},
		},
	})

	if err != nil {
		log.Fatalln("Unable to create schedule:", err)
	} else {
		log.Println("Schedule created successfully!")
	}
}
