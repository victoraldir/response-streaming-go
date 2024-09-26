package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	log.Printf("M3U8 URL to be streamed: %s\n", m3u8Url)
	// file, err := os.CreateTemp("", "playlist*.mp4")
	// if err != nil {
	// 	return nil, err
	// }

	// Create a client to probe the m3u8 file
	client := &http.Client{}

	// Probe the m3u8 file
	resp, err := client.Get(m3u8Url)
	if err != nil {
		log.Printf("Error probing m3u8 file: %v\n", err)
		return nil, err
	}

	// Check if the response is OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error probing m3u8 file: %v\n", err)
		return nil, fmt.Errorf("Error probing m3u8 file: %v", err)
	}

	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	fmt.Printf("Content-Length: %s\n", resp.Header.Get("Content-Length"))
	fmt.Printf("Content-Disposition: %s\n", resp.Header.Get("Content-Disposition"))
	fmt.Printf("Content-Encoding: %s\n", resp.Header.Get("Content-Encoding"))

	videoPath := fmt.Sprintf("/tmp/%s.mp4", time.Now().Format("20060102150405"))

	log.Printf("Temp file created: %s\n", videoPath)

	var errBug, outBuf strings.Builder

	// Print ffmpeg -version
	cmd := exec.Command("ffmpeg", "-version")
	cmd.Stdout = &outBuf
	err = cmd.Run()
	if err != nil {
		log.Printf("Error executing ffmpeg: %v\n", err)
		return nil, err
	}

	log.Printf("ffmpeg version: %s\n", outBuf.String())

	// Probe the m3u8 file with ffprobe
	cmd = exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", m3u8Url)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBug
	err = cmd.Run()
	if err != nil {
		log.Printf("Error executing ffprobe: %v\n", err)
		return nil, err
	}

	log.Printf("ffprobe output: %s\n", outBuf.String())

	// execute cmd ffmpeg -i {url} -c copy -bsf:a aac_adtstoasc {file}
	cmd = exec.Command("ffmpeg", "-y", "-i", m3u8Url, "-c", "copy", "-bsf:a", "aac_adtstoasc", videoPath)

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBug
	err = cmd.Run()
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
