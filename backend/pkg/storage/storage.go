package storage

import (
	"context"
	"io"

	"github.com/spf13/viper"
)

type StorageType string

const (
	Local = "local"
	MinIO = "minio"
)

type Storage interface {
	// Avatar
	UploadAvatar(ctx context.Context, uid string, filename string, reader io.Reader) (string, error)
	GetAvatarURL(ctx context.Context, uid string, filename string) (string, error)
}

func NewStorage(conf *viper.Viper) *Storage {
	var s Storage
	storageType := conf.GetString("storage.type")
	switch storageType {
	case Local: // 本地存储
		s = NewLocalStorage(conf)
	case MinIO: // MinIO 存储
		s = NewMinIOStorage(conf)
	default:
		panic("unsupported storage type")
	}
	return &s
}
