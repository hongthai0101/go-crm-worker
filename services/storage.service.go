package services

import (
	"cloud.google.com/go/storage"
	"context"
	"crm-worker-go/config"
	"crm-worker-go/utils"
	"io"
	"log"
	"time"
)

type StorageService interface {
	UploadFile(file io.Reader, object string) (*storage.Writer, error)
}

type clientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

func NewStorageService(uploadPath string) StorageService {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		utils.Logger.Error(err)
		log.Fatalf("Failed to create client: %v", err)
	}
	return &clientUploader{
		client,
		config.GetConfig().GCSConfig.ProjectId,
		config.GetConfig().GCSConfig.Buket,
		uploadPath,
	}
}

// UploadFile uploads an object
func (c *clientUploader) UploadFile(file io.Reader, object string) (*storage.Writer, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + "/" + object).NewWriter(ctx)
	wc.ObjectAttrs.ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if _, err := io.Copy(wc, file); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	utils.Debug(wc.MediaLink, wc.Size)

	if err := wc.Close(); err != nil {
		log.Fatalf("Writer.Close: %v", err)
	}
	return wc, nil
}
