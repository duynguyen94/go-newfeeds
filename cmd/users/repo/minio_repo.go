package repo

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	defaultEndpoint        = "localhost"
	defaultAccessKeyID     = ""
	defaultSecretAccessKey = ""
	defaultUseSSL          = false
)

func CreateMinioClient() (*minio.Client, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(defaultEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(defaultAccessKeyID, defaultSecretAccessKey, ""),
		Secure: defaultUseSSL,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
