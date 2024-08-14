package main

import (
	"github.com/minio/minio-go/v7"
	"time"
)

const (
	bucket = "images"
)

type ImageStorageModel struct {
	client *minio.Client
}

func (s *ImageStorageModel) PutImage(content string, bucket string, path string) error {
	// TODO
	return nil
}
func (s *ImageStorageModel) GetSignedUrl(bucket string, path string, expiration time.Time) (string, error) {
	// TODO
	return "", nil
}
func (s *ImageStorageModel) DeleteImage(bucket string, path string) error {
	// TODO
	return nil
}
