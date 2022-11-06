package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"io"
)

type S3CloudStorage struct {
	session *session.Session
	s3      *s3.S3
	bucket  string
	prefix  string
}

func (cloud S3CloudStorage) Setup(bucket string, prefix string, client string, secret string, endpointURL string, region string) error {
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
	cloud.prefix = prefix
	inf.Msg("Establishing AWS S3 session...")
	return cloud.testConnection()
}

func (cloud S3CloudStorage) testConnection() error {
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
	return nil
}

func (cloud S3CloudStorage) Put(path string, reader *io.ReadSeeker) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", cloud.prefix, path)
	_, err := cloud.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(cloud.bucket),
		Key:    aws.String(key),
		Body:   *reader,
	})
	return err
}

func (cloud S3CloudStorage) Get(path string, writer *io.Writer) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", cloud.prefix, path)

}
