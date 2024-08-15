package main

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
	defaultBucket = "images"
)

type ImagePostStorageModel struct {
	client *minio.Client
	bucket string
}

func (s *ImagePostStorageModel) bucketExists() error {
	ctx := context.Background()
	exists, errBucketExists := s.client.BucketExists(ctx, s.bucket)

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

	info, err := s.client.PutObject(
		context.Background(), s.bucket, objectKey,
		reader, fileSize, minio.PutObjectOptions{},
	)

	if err != nil {
		return "", err
	}

	return info.Key, nil
}

func (s *ImagePostStorageModel) GetSignedUrl(path string, expiration time.Duration) (string, error) {
	reqParams := make(url.Values)

	presignedURL, err := s.client.PresignedGetObject(
		context.Background(), s.bucket, path,
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
