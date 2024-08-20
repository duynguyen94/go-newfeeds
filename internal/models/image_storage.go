package models

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

const (
	DefaultBucket = "images"
)

type ImagePostStorageModel struct {
	Client *minio.Client
	Bucket string
}

func (s *ImagePostStorageModel) BucketExists() error {
	ctx := context.Background()
	exists, errBucketExists := s.Client.BucketExists(ctx, s.Bucket)

	if errBucketExists != nil {
		return errBucketExists
	}

	if exists == false {
		err := errors.New("Bucket is not existed")
		return err
	}

	return nil
}

func (s *ImagePostStorageModel) getKey(postId int, filename string) string {
	return "post/" + strconv.Itoa(postId) + "/" + filename
}

func (s *ImagePostStorageModel) PutImage(reader *os.File, postId int, filename string, fileSize int64) (string, error) {
	objectKey := s.getKey(postId, filename)

	info, err := s.Client.PutObject(
		context.Background(), s.Bucket, objectKey,
		reader, fileSize, minio.PutObjectOptions{},
	)

	if err != nil {
		return "", err
	}

	return info.Key, nil
}

func (s *ImagePostStorageModel) GetSignedUrl(path string, expiration time.Duration) (string, error) {
	reqParams := make(url.Values)

	presignedURL, err := s.Client.PresignedGetObject(
		context.Background(), s.Bucket, path,
		expiration, reqParams,
	)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return presignedURL.String(), nil
}
func (s *ImagePostStorageModel) DeleteImage(path string) error {
	// TODO
	return nil
}
