package sqsclient

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)


func sqsOperation(ctx context.Context){
    // Load the Shared AWS Configuration (~/.aws/config)
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Create an Amazon SQS service client
    client := sqs.NewFromConfig(cfg)
	

	getQueueAttributes := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/478867930449/inspirate-email-notification-queue"),
 	}

	result, err := client.GetQueueAttributes(ctx, getQueueAttributes)

    // Get the first page of results for ListObjectsV2 for a bucket
    // //output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
    //     Bucket: aws.String("amzn-s3-demo-bucket"),
    // })
    // if err != nil {
    //     log.Fatal(err)
    // }

	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(result)
}
