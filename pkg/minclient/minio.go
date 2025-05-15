package minio

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type MinioClient struct {
	*minio.Client
}

func NewMinioClient() (*MinioClient, error) {
	minioClient := &MinioClient{}
	var err error

	minioClient.Client, err = minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: false,
	})

	return minioClient, err
}

func (m *MinioClient) NewBucket(ctx context.Context, bucketName string) {
	bucketsExists, err := m.BucketExists(ctx, bucketName)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	if !bucketsExists {
		err := m.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			logrus.Fatal(err.Error())
		}
	}
}

func (m *MinioClient) GenerateStorageUrl(bucketName, fileName string) string {
	storageUrl := fmt.Sprintf("http://%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), bucketName, fileName)
	return storageUrl
}

func (m *MinioClient) UploadFile(
	ctx context.Context,
	bucketName, fileName string,
	reader io.Reader,
	fileSize int64,
	contentType string) error {

	_, err := m.PutObject(ctx,
		bucketName, fileName, reader, fileSize,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	return nil
}
