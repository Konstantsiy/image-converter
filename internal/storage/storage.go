// Package storage implements the functionality of file storage ()
package storage

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	uuid "github.com/satori/go.uuid"
)

const (
	URLTimeout       = 10 * time.Minute
	LocationTemplate = "https://%s.amazonaws.com/%s/%s"
)

//S3Error represent storage-related errors.
type S3Error struct {
	message string
}

func (s3err *S3Error) Error() string {
	return fmt.Sprintf("storage error: %s", s3err.message)
}

// S3Config used to configure the session, create a bucket,
// and connect to the SDK's service client.
type S3Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

//Storage implements the functionality of file storage (Amazon S3).
type Storage struct {
	svc        *s3.S3
	s3conf     *S3Config
	bucketName string
}

//NewStorage creates new file storage with the given S3 configs and bucket name.
func NewStorage(s3conf S3Config, bucketName string) *Storage {
	return &Storage{
		svc:        &s3.S3{},
		s3conf:     &s3conf,
		bucketName: bucketName,
	}
}

//InitS3ServiceClient initializes SDK's service client.
func (s *Storage) InitS3ServiceClient() error {
	s3session, err := s.createSession()
	if err != nil {
		return &S3Error{message: err.Error()}
	}

	s.svc = s3.New(s3session)
	err = s.createBucket()
	if err != nil {
		return err
	}

	return nil
}

//createSession creates and returns a new session.
func (s *Storage) createSession() (*session.Session, error) {
	s3session, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.s3conf.Region),
		Credentials: credentials.NewStaticCredentials(s.s3conf.AccessKeyID, s.s3conf.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, &S3Error{message: err.Error()}
	}
	return s3session, nil
}

//createBucket
func (s *Storage) createBucket() error {
	_, err := s.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s.bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(s.s3conf.Region),
		},
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				return &S3Error{message: "bucket name already in use"}
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				return &S3Error{message: "bucket exists and is owned by you"}
			default:
				return &S3Error{message: err.Error()}
			}
		}
	}
	return nil
}

//UploadFile uploads the given file to the bucket.
func (s *Storage) UploadFile(file multipart.File) (string, string, error) {
	fileUUID := uuid.NewV4()
	fileUUIDStr := fileUUID.String()
	location := fmt.Sprintf(LocationTemplate, s.s3conf.Region, s.bucketName, fileUUID)

	_, err := s.svc.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(location),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})
	if err != nil {
		return "", "", &S3Error{message: err.Error()}
	}

	return fileUUIDStr, location, nil
}

//GetDownloadURL returns URL to download а file from the bucket by the given file location.
func (s *Storage) GetDownloadURL(location string) (string, error) {
	req, _ := s.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(location),
	})

	url, err := req.Presign(URLTimeout)
	if err != nil {
		return "", &S3Error{message: fmt.Sprintf("can't create requets's presigned URL, %s", err.Error())}
	}

	return url, err
}
