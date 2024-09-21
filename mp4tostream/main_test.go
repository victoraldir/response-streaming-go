package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestRenderPage(t *testing.T) {

	// Test case: Test the renderPage function
	t.Run("Should return a valid response", func(t *testing.T) {

		// Arrange
		request := events.LambdaFunctionURLRequest{
			QueryStringParameters: map[string]string{
				"url": "https://example.com/video.mp4",
			},
		}

		// Act
		response, err := lambdaHandler(&request)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, response)

	})

}
