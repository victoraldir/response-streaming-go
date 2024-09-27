package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func lambdaHandler(request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {

	log.Printf("Request: %v", request)

	// Get url from query parameter
	m3u8Url := request.QueryStringParameters["url"]

	// Create a temporary .mp4 file to store the stream
	log.Printf("M3U8 URL to be streamed test: %s\n", m3u8Url)

	videoPath := fmt.Sprintf("/tmp/%s.mp4", time.Now().Format("20060102150405"))

	log.Printf("Temp file created: %s\n", videoPath)

	var errBug, outBuf strings.Builder

	// execute cmd ffmpeg -i {url} -c copy -bsf:a aac_adtstoasc {file}
	cmd := exec.Command("ffmpeg", "-y", "-i", m3u8Url, "-c", "copy", "-bsf:a", "aac_adtstoasc", videoPath)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBug
	err := cmd.Run()
	if err != nil {
		log.Printf("Error executing ffmpeg: %v\n", err)
		return nil, err
	}

	log.Printf("ffmpeg output: %s\n", outBuf.String())
	log.Printf("ffmpeg error: %s\n", errBug.String())

	// We will use a pipe to stream the file to the response
	r, w := io.Pipe()

	// Pipe the file to the response
	go func() {
		file, err := os.OpenFile(videoPath, os.O_RDONLY, 0644)
		defer w.Close()

		if err != nil {
			log.Printf("Error opening file: %v\n", err)
			return
		}

		io.Copy(w, file)
	}()

	return &events.LambdaFunctionURLStreamingResponse{
		StatusCode: 200,
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
