package main

import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
)

func convert(ctx context.Context, reader io.Reader) (string, error) {

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return "", fmt.Errorf("Failed to create client: %v", err)
	}
	defer client.Close()

	image, err := vision.NewImageFromReader(reader)
	if err != nil {
		return "", fmt.Errorf("Failed to create image: %v", err)
	}

	text, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to detect text: %v", err)
	}

	return text.Text, nil
}
