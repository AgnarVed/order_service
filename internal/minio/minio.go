package minio

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"tests2/internal/config"
)

type MinIOClient struct {
	cfg *config.Config
	*minio.Client
}

func NewMinIOClient(cfg *config.Config) (*MinIOClient, error) {
	client, err := minio.New(cfg.MinIOURL, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.MinIOUser, cfg.MinIOPass, ""),
	})

	if err != nil {
		return nil, err
	}

	return &MinIOClient{Client: client, cfg: cfg}, nil
}

func (m *MinIOClient) InitDocBucket(ctx context.Context) error {
	exists, err := m.BucketExists(ctx, m.cfg.MinIODocBucket)
	if err != nil {
		return err
	}

	if !exists {
		if err := m.MakeBucket(ctx, m.cfg.MinIODocBucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (m *MinIOClient) DownloadDoc(ctx context.Context, name string) ([]byte, error) {

	obj, err := m.GetObject(ctx, m.cfg.MinIODocBucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return io.ReadAll(obj)
}

func (m *MinIOClient) UploadDoc(ctx context.Context, name string, buffer *bytes.Buffer) error {
	_, err := m.PutObject(ctx, m.cfg.MinIODocBucket, name, buffer, int64(buffer.Len()), minio.PutObjectOptions{})
	return err
}
