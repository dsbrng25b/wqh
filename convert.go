package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/otiai10/gosseract"
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

func convertTesseract(reader io.Reader) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetVariable("load_system_dawg", "false")
	client.SetVariable("load_freq_dawg", "false")
	client.SetVariable("user_words_suffix", "user-words")

	picture, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	client.SetImageFromBytes(picture)

	return client.Text()
}
