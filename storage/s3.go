package storage

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
	"time"
)

type S3CloudStorage struct {
	session *session.Session
	s3      *s3.S3
	bucket  string
	prefix  string
}

func (cloud S3CloudStorage) Setup(bucket string, prefix string, test string, client string, secret string, endpointURL string, region string) error {
	config := aws.Config{}
	inf := log.Info()
	if region != "" {
		config.Region = aws.String(region)
		inf.Str("region", region)
	}
	if client != "" && secret != "" {
		log.Debug().Str("client", client).Msg("AWS client+secret set")
		config.Credentials = credentials.NewStaticCredentials(client, secret, "")
		inf.Bool("client_set", true)
		inf.Bool("secret_set", true)
	}
	if region != "" {
		config.Endpoint = aws.String(endpointURL)
		inf.Str("endpoint_url", endpointURL)
	}
	cloud.session = session.Must(session.NewSession(&config))
	cloud.s3 = s3.New(cloud.session)
	cloud.bucket = bucket
	cloud.prefix = strings.TrimRight(prefix, "/")
	inf.Msg("Establishing AWS S3 session...")
	return cloud.testConnection(test)
}

func (cloud S3CloudStorage) testConnection(test string) error {
	helloWorld := fmt.Sprintf("hello %s", time.Now())
	err := cloud.Put(test, strings.NewReader(helloWorld))
	if err != nil {
		return err
	}
	buf := &aws.WriteAtBuffer{}
	err = cloud.Get(test, buf)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(string(buf.Bytes()), "hello ") {
		return errors.New("wrong prefix in test file")
	}
	return nil
	//out, err := sts.New(cloud.session).GetCallerIdentity(&sts.GetCallerIdentityInput{})
	//if err != nil {
	//	msg := "Failed to establish AWS S3 session"
	//	if awsErr, ok := err.(awserr.Error); ok {
	//		switch awsErr.Code() {
	//		default:
	//			log.Error().Err(awsErr.OrigErr()).Str("code", awsErr.Code()).Str("message", awsErr.Message()).Msg(msg)
	//		}
	//	} else {
	//		log.Error().Err(err).Msg(msg)
	//	}
	//	return err
	//}
	//log.Info().Str("arn", aws.StringValue(out.Arn)).Str("user_id", aws.StringValue(out.UserId)).Str("account", aws.StringValue(out.Account)).Msg("AWS S3 session established")
}

func (cloud S3CloudStorage) GetKey(path string) string {
	return fmt.Sprintf("%s/%s", cloud.prefix, strings.TrimLeft(path, "/"))
}

func (cloud S3CloudStorage) Put(path string, reader io.Reader) error {
	key := cloud.GetKey(path)
	uploader := s3manager.NewUploaderWithClient(cloud.s3)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cloud.bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	return err
}

func (cloud S3CloudStorage) Get(path string, writer io.WriterAt) error {
	key := cloud.GetKey(path)
	downloader := s3manager.NewDownloaderWithClient(cloud.s3)
	_, err := downloader.Download(writer, &s3.GetObjectInput{
		Bucket: aws.String(cloud.bucket),
		Key:    aws.String(key),
	})
	return err
}
