package minio

import (
	"github.com/minio/minio-go/v7"
	"log"
	"os"
)

type Storage struct {
	c *minio.Client
}

func (s Storage) New(c *minio.Client) Storage {
	return Storage{
		c: c,
	}
}

func (s Storage) StoreTicket(objName string, file *os.File) error {
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("failed to get information about file %v", err)
	}

	_, err = s.c.PutObject(nil, "tickets", objName, file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		log.Fatalln("Failed to load PDF file in MinIO:", err)
	}
}
