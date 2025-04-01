package quote

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func GenerateQuote(ctx workflow.Context, input GenerateTextInput) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var textResult GenerateTextOutput
	err := workflow.ExecuteActivity(ctx, GenerateText, input).Get(ctx, &textResult)
	if err != nil {
		return err
	}

	var imageResult GenerateImageOutput
	err = workflow.ExecuteActivity(ctx, GenerateImage, textResult).Get(ctx, &imageResult)
	if err != nil {
		return err
	}

	var publicationResult string
	err = workflow.ExecuteActivity(ctx, PublishInstagramPost, imageResult).Get(ctx, &publicationResult)
	if err != nil {
		return err
	}

	return err
}
