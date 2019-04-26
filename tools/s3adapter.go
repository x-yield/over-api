package tools

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"

	//"github.com/x-yield/over-api/internal/config"
)

const (
	ammoBucket      = "ammo"
	artifactsBucket = "artifacts"
)

// S3Uploader describes s3 uploader
type S3Uploader interface {
	UploadWithContext(ctx aws.Context, input *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type S3Downloader interface {
	DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error)
}

// S3service service
type S3service struct {
	Client          *s3.S3
	Uploader        S3Uploader
	Downloader      S3Downloader
	ammoBucket      string
	artifactsBucket string
}

func (s *S3service) SetAmmoBucket(bucket string) {
	s.ammoBucket = bucket
}

func (s *S3service) GetAmmoBucket() string {
	return s.ammoBucket
}

func (s *S3service) SetArtifactsBucket(bucket string) {
	s.artifactsBucket = bucket
}

func (s *S3service) GetArtifactsBucket() string {
	return s.artifactsBucket
}

// New create new s3 service
func NewS3Service() *S3service {

	var (
		s3Endpoint = config.GetValue(context.Background(), config.S3Endpoint).String()
		s3Access   = config.GetValue(context.Background(), config.S3Access).String()
		s3Secret   = config.GetValue(context.Background(), config.S3Secret).String()
	)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(s3Access, s3Secret, ""),
		Endpoint:         aws.String(s3Endpoint),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Println(err)
	}

	Client := s3.New(newSession)
	Client.Handlers.Validate.RemoveByName("core.ValidateEndpointHandler")
	Uploader := s3manager.NewUploaderWithClient(Client, func(u *s3manager.Uploader) {})
	Downloader := s3manager.NewDownloaderWithClient(Client, func(d *s3manager.Downloader) {})

	return &S3service{
		Client:          Client,
		Uploader:        Uploader,
		Downloader:      Downloader,
		ammoBucket:      ammoBucket,
		artifactsBucket: artifactsBucket,
	}
}

// Close - does nothing. made for overall consistency
func (s *S3service) Close() error {
	return nil
}

// Upload - upload file to s3 into any bucket
func (s *S3service) Upload(ctx context.Context, bucket string, filename string, reader io.Reader, acl string) (*s3manager.UploadOutput, error) {

	resp, err := s.Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   reader,
		ACL:    &acl,
	})

	if err != nil {
		log.Printf("Failed to upload file %v", err)
		return nil, errors.Wrap(err, "cannot upload the file")
	}
	return resp, nil
}

// UploadAmmo - upload file into "ammo" bucket
func (s *S3service) UploadAmmo(ctx context.Context, filename string, reader io.Reader) (*s3manager.UploadOutput, error) {
	return s.Upload(ctx, s.ammoBucket, filename, reader, s3.ObjectCannedACLPublicRead)
}

// UploadArtifact - upload file into "artifact" bucket
func (s *S3service) UploadArtifact(ctx context.Context, filename string, reader io.Reader) (*s3manager.UploadOutput, error) {
	return s.Upload(ctx, s.artifactsBucket, filename, reader, s3.ObjectCannedACLPublicRead)
}

// Download - download file from any s3 bucket
func (s *S3service) Download(ctx context.Context, bucket string, filename string, writer io.WriterAt) (int64, error) {
	// Create a downloader with the s3 client and custom options
	resp, err := s.Downloader.DownloadWithContext(ctx, writer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		log.Printf("Failed to download file %v", err)
		return 0, errors.Wrap(err, "Failed to download file ")
	}
	return resp, nil
}

// DownloadAmmo - download file from "ammo" bucket
func (s *S3service) DownloadAmmo(ctx context.Context, filename string, writer io.WriterAt) (int64, error) {
	return s.Download(ctx, s.ammoBucket, filename, writer)
}

// ListObjects - get list of files in a bucket
// TODO: pagination
func (s *S3service) ListObjects(ctx context.Context, bucket *string, prefix *string) ([]*s3.Object, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: bucket,
		Prefix: prefix,
	}
	output, err := s.Client.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		log.Printf("Failed to list objects for bucket %v", err)
		return []*s3.Object{}, errors.Wrap(err, "Failed to list objects")
	}
	return output.Contents, nil
}

// ListAmmo - get list of ammofiles (files in "ammo" bucket)
func (s *S3service) ListAmmo(ctx context.Context) ([]*s3.Object, error) {
	return s.ListObjects(ctx, &s.ammoBucket, nil)
}

// ListArtifacts - get list of ammofiles (files in "ammo" bucket)
func (s *S3service) ListArtifacts(ctx context.Context, prefix string) ([]*s3.Object, error) {
	return s.ListObjects(ctx, &s.artifactsBucket, &prefix)
}

// Delete - deletes file from any bucket
func (s *S3service) Delete(key string, bucket string) error {
	object := &s3.DeleteObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(bucket),
	}
	_, err := s.Client.DeleteObject(object)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAmmo - deletes file from ammo bucket
func (s *S3service) DeleteAmmo(ctx context.Context, key string) error {
	return s.Delete(key, ammoBucket)
}
