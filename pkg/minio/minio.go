package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"log"
	"os"
)

type Storage struct {
	c          *minio.Client
	bucketName string
}

func New(c *minio.Client, bucketName string) Storage {
	return Storage{
		c:          c,
		bucketName: bucketName,
	}
}

func (s Storage) Store(ctx context.Context, file *os.File) (string, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("failed to get information about file %v", err)
		return "", err
	}

	_, err = s.c.PutObject(ctx, s.bucketName, fileInfo.Name(), file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		log.Printf("failed to load PDF file in MinIO: %v", err)
		return "", err
	}

	url := "http://localhost:" + s.c.EndpointURL().Port() + "/" + s.bucketName + "/" + fileInfo.Name()

	return url, nil
}
