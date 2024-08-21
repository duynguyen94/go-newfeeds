package object_store

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"net/url"
	"os"
	"strconv"
	"time"
)

type ImageStorage interface {

	// BucketExists simple ping for check connection and bucket exised
	BucketExists(bucketName string) error

	// Put upload image to blob store
	Put(reader *os.File, postId int, filename string, fileSize int64) (string, error)

	// GetSignedUrl get downloadable url for given image
	GetSignedUrl(bucketName string, path string, expiration time.Duration) (string, error)
}

func NewImageStorage(client *minio.Client) ImageStorage {
	return &imagePostStorage{client: client}
}

type imagePostStorage struct {
	client *minio.Client
	Bucket string
}

func (s *imagePostStorage) BucketExists(bucketName string) error {
	ctx := context.Background()
	exists, errBucketExists := s.client.BucketExists(ctx, bucketName)

	if errBucketExists != nil {
		return errBucketExists
	}

	if exists == false {
		err := errors.New("Bucket is not existed")
		return err
	}

	return nil
}

func (s *imagePostStorage) getKey(postId int, filename string) string {
	return "post/" + strconv.Itoa(postId) + "/" + filename
}

func (s *imagePostStorage) Put(reader *os.File, postId int, filename string, fileSize int64) (string, error) {
	objectKey := s.getKey(postId, filename)

	info, err := s.client.PutObject(
		context.Background(), s.Bucket, objectKey,
		reader, fileSize, minio.PutObjectOptions{},
	)

	if err != nil {
		return "", err
	}

	return info.Key, nil
}

func (s *imagePostStorage) GetSignedUrl(bucketName string, path string, expiration time.Duration) (string, error) {
	reqParams := make(url.Values)

	presignedURL, err := s.client.PresignedGetObject(
		context.Background(), s.Bucket, path,
		expiration, reqParams,
	)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return presignedURL.String(), nil
}
