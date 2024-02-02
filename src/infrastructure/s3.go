package infrastructure

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/justepl2/gopro_to_gpx_api/domain"
)

type S3FileStorage struct {
	s3     *s3.S3
	bucket string
	sess   *session.Session
}

func NewS3FileStorage() domain.FileStorage {
	bucket := os.Getenv("AWS_S3_BUCKET")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))

	return &S3FileStorage{
		s3:     s3.New(sess),
		bucket: bucket,
		sess:   sess,
	}
}

func (s *S3FileStorage) UploadFiles(path string, file []byte) error {
	// Upload file to S3
	_, err := s.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(file),
	})

	return err
}

func (s *S3FileStorage) GetFile(path string) ([]byte, error) {
	// TODO: Implement GetFile
	return nil, nil
}
