package storage

import (
	"context"
	"io"

	"github.com/spf13/viper"
)

// 存储类型
type StorageType string

const (
	Local = "local"
	MinIO = "minio"
)

// 定义存储接口
type Storage interface {
	// Avatar
	UploadAvatar(ctx context.Context, uid string, filename string, reader io.Reader) (string, error)
	GetAvatarURL(ctx context.Context, uid string, filename string) (string, error)
}

func NewStorage(conf *viper.Viper) Storage {
	StorageType := conf.GetString("storage.type")
	switch StorageType {
	case Local:
		return NewLocalStorage(conf)
	case MinIO:
		return NewMinIOStorage(conf)
	default:
		panic("unsupported storage type")
	}
}
