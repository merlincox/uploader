package main

import (
	"os"
	"flag"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3"
	"fmt"
)

type App struct {
	bucket   string
	region   string
	profile  string
	filePath string
	key      string
	list     bool
	upload   bool
}

var app App

func main() {

	readArgs()

	awsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(app.region)},
		Profile: app.profile,
	})

	if err != nil {
		log.Fatalln("Getting aws session threw error: ", err)
	}

	if app.upload {
		upload(awsSession, app.bucket, app.filePath, app.key)
	}

	if app.list {
		listBucket(awsSession, app.bucket)
	}
}

func listBucket(awsSession *session.Session, bucket string) {

	svc := s3.New(awsSession)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}

	resp, _ := svc.ListObjects(params)

	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
	}
}

func readArgs() {

	var noUpload bool

	flag.StringVar(&app.profile, "p", "default", "Defines AWS profile to use")
	flag.StringVar(&app.bucket, "b", "", "Defines S3 bucket to use")
	flag.StringVar(&app.region, "r", "", "Defines AWS region to use")
	flag.StringVar(&app.filePath, "f", "", "Defines path of local source file to upload")
	flag.StringVar(&app.key, "d", "", "Defines S3 destination file path to upload to: (default same as path of local source file)")
	flag.BoolVar(&app.list, "l", false, "Lists bucket contents after upload")
	flag.BoolVar(&noUpload, "n", false, "Only lists bucket contents (no upload)")

	flag.Parse()

	app.upload = !noUpload

	if noUpload {
		app.list = true
	}

	if app.bucket == "" {
		log.Fatalln("S3 bucket must be defined with -b")
	}

	if app.region == "" {
		log.Fatalln("AWS region must be defined with -r")
	}

	if app.upload && app.filePath == "" {
		log.Fatalln("Local source file path must be defined with -f")
	}

	if app.profile == "default" {
		fmt.Println("AWS profile set to 'default'. To change use -p")
	}

	if app.upload && app.key == "" {
		app.key = app.filePath
		log.Printf("S3 destination file path set to '%v'. To change use -d\n", app.filePath)
	}

	if app.list == false {
		fmt.Println("S3 bucket listing is OFF. To change use -l")
	}
}

func upload(awsSession *session.Session, bucket string, filePath string, key string) {

	uploader := s3manager.NewUploader(awsSession)

	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal("Error opening file: ", err)
	}

	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
		Body: file,
	})

	if err != nil {
		log.Fatal("Error uploading: ", err)
	}

	log.Printf("Uploaded %v to %v S3 bucket as %v\n", filePath, bucket, key)
}
