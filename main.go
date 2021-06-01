package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"strings"
)
func LambdaHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err!=nil {
		return err
	}
	for _, message := range sqsEvent.Records {
		d := decode(message.Body)
		fmt.Println(" decode message body: " + d)
		upload(cfg, "mybucket0106202586", message.MessageId, d)
		//fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}

	return nil
}

func decode(message string) string {
	b := []byte(message)
	sEnc := b64.StdEncoding.EncodeToString(b)
	return sEnc
}

func upload(cfg aws.Config, bucket string, filename string, data string) {
	awsClient := s3.NewFromConfig(cfg)
	if bucket == "" || filename == "" {
		fmt.Println("You must supply a bucket name   and file name ")
		return
	}

	// create a reader from data data in memory
	reader := strings.NewReader(data)

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &filename,
		Body:   reader,
	}

	_, err := awsClient.PutObject(context.TODO(), input)
	if err != nil {
		fmt.Println("Got error uploading file:")
		fmt.Println(err)
		return
	}

}
func main() {
	lambda.Start(LambdaHandler)
}
