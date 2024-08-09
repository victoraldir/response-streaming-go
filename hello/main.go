package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func lambdaHandler(request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {

	// Get url from query parameter
	mp4Url := request.QueryStringParameters["url"]

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

	// Save it into /tmp/myfile.mp4
	file, err := os.Create("/tmp/myfile.mp4")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	log.Printf("Downloading %s\n", mp4Url)
	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return nil, err
	}

	// Pipe the file to the response
	go func() {
		defer w.Close()

		// Open the file
		file, err := os.Open("/tmp/myfile.mp4")

		if err != nil {
			log.Printf("Error opening file: %v\n", err)
			return
		}

		io.Copy(w, file)
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
