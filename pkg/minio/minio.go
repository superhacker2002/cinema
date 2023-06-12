package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"log"
	"os"
)

type Storage struct {
	c *minio.Client
}

func New(c *minio.Client) Storage {
	return Storage{
		c: c,
	}
}

func (s Storage) StoreTicket(ctx context.Context, file *os.File) (string, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("failed to get information about file %v", err)
		return "", err
	}

	_, err = s.c.PutObject(ctx, "tickets", fileInfo.Name(), file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		log.Printf("failed to load PDF file in MinIO: %v", err)
		return "", err
	}

	url := s.c.EndpointURL().String() + "/" + "tickets" + "/" + fileInfo.Name()

	return url, nil
}
