package main

import (
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func lambdaHandler(request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {

	// Get url from query parameter
	mp4Url := request.QueryStringParameters["url"]

	log.Printf("MP4 URL to be streamed: %s\n", mp4Url)

	// We will use a pipe to stream the response from the client to the response
	r, w := io.Pipe()

	// Prepare request
	req, err := http.NewRequest(http.MethodGet, mp4Url, nil)
	if err != nil {
		return nil, err
	}

	// Perform request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Pipe the file to the response
	go func() {
		defer w.Close()
		defer resp.Body.Close()

		if err != nil {
			log.Printf("Error opening file: %v\n", err)
			return
		}

		log.Printf("Copying response body to pipe\n")
		io.Copy(w, resp.Body)

	}()

	return &events.LambdaFunctionURLStreamingResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":        "video/mp4",
			"Content-Disposition": "attachment; filename=myfile.mp4",
		},
		Body: r,
	}, nil
}

func main() {
	lambda.Start(lambdaHandler)
}
