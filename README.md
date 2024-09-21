
# What is that response-streaming-go?

This is a simple example of how to use [Lambda Response Streaming](https://aws.amazon.com/blogs/compute/introducing-aws-lambda-response-streaming/) in Go. 
In this example, the Lambda function receives a request with an URL of a MP4 video file, then the Go function makes a request to the URL and pipe the response to the Lambda response stream in a go routine. 

# Current limitations and considerations

- Response payload size is limited to 20MB. It's a soft limit and can be increased by contacting AWS support.
- Lambda function timeout is limited to 15 minutes.

## Requirements

- [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)
- [Go](https://golang.org/doc/install)
- [Taskfile](https://taskfile.dev/#/installation)

## How to run

Within the project folder, run the following command:

```bash
task deploy
```

## How to test

After the deployment, you can test the function using the following command:

```bash
curl -v "https://{DISTRIBUTION_ID}.cloudfront.net/?url=https://www.pexels.com/download/video/7230308/" -o sample-30s2.mp4
```
## Cleanup

To remove all resources created by the SAM template, run the following command:

```bash
task destroy
```
