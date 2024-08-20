package conn

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
)

func CreateMinioClient() (*minio.Client, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT")+":"+os.Getenv("MINIO_API_PORT"), &minio.Options{
		Creds: credentials.NewStaticV4(
			os.Getenv("MINIO_ACCESS_KEY_ID"),
			os.Getenv("MINIO_SECRET_ACCESS_KEY"), "",
		),
		Secure: false,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
