package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/y-yagi/tsudura/utils"
)

type Client struct {
	S3  *s3.S3
	cfg *utils.Config
}

type Result struct {
	Key  string
	ETag string
}

func Init(cfg *utils.Config) (*Client, error) {
	return &Client{S3: buildClient(cfg), cfg: cfg}, nil
}

func buildClient(cfg *utils.Config) *s3.S3 {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.Secret, cfg.Token, ""),
		Endpoint:         aws.String(cfg.Endpoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)
	client := s3.New(newSession)

	return client
}

func (c *Client) Upload(path string) (*Result, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var output *s3manager.UploadOutput
	uploader := s3manager.NewUploaderWithClient(c.S3)
	key := c.sanitizeKey(strings.Split(path, c.cfg.Root)[1])

	err = retry.Do(
		func() error {
			output, err = uploader.Upload(&s3manager.UploadInput{
				Body:   strings.NewReader(string(content)),
				Bucket: aws.String(c.cfg.Bucket),
				Key:    aws.String(key),
			})

			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &Result{Key: key, ETag: *output.ETag}, nil
}

func (c *Client) Destroy(path string) (*Result, error) {
	key := c.sanitizeKey(strings.Split(path, c.cfg.Root)[1])

	err := retry.Do(
		func() error {
			_, err := c.S3.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(c.cfg.Bucket),
				Key:    aws.String(key),
			})

			if err != nil {
				return err
			}
			return nil
		},
	)

	return &Result{Key: key}, err
}

func (c *Client) Download(path string, etag string) (*Result, error) {
	key := c.sanitizeKey(strings.Split(path, c.cfg.Root)[1])
	downloader := s3manager.NewDownloaderWithClient(c.S3)
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %v", err)
	}

	err = retry.Do(
		func() error {
			_, err := downloader.Download(f, &s3.GetObjectInput{
				Bucket: aws.String(c.cfg.Bucket),
				Key:    aws.String(key),
			})

			if err != nil {
				return err
			}
			return nil
		},
	)

	return &Result{Key: key}, err
}

func (c *Client) sanitizeKey(key string) string {
	if os.PathSeparator == '\\' {
		key = strings.ReplaceAll(key, string(os.PathSeparator), "/")
	}

	return key
}
