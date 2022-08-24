package service

import (
	"bytes"
	"context"
	"tests2/internal/minio"
	"tests2/internal/models"
)

type doc struct {
	minio *minio.MinIOClient
}

func (d *doc) DownloadDoc(ctx context.Context, filename string) ([]byte, error) {
	return d.minio.DownloadDoc(ctx, filename)
}

func (d *doc) UploadDoc(ctx context.Context, fileInfo *models.FileInfo) error {
	toMinIO := &bytes.Buffer{}

	minioName := fileInfo.ObjectName

	if err := d.minio.UploadDoc(ctx, minioName, toMinIO); err != nil {
		return err
	}
	return nil
}

func NewDocService(minio *minio.MinIOClient) Doc {
	return &doc{minio: minio}
}
