package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(conf *viper.Viper) *LocalStorage {
	basePath := conf.GetString("storage.local.path")
	if basePath == "" {
		panic("local storage base path is required")
	}

	// 确保目录存在
	if err := os.MkdirAll(basePath, 0744); err != nil {
		panic(fmt.Sprintf("failed to create local storage directory: %v", err))
	}

	return &LocalStorage{
		basePath: basePath,
	}
}

func (s *LocalStorage) UploadAvatar(ctx context.Context, uid string, filename string, reader io.Reader) (string, error) {
	// 创建用户目录
	userDir := filepath.Join(s.basePath, uid)
	if err := os.MkdirAll(userDir, 0744); err != nil {
		return "", fmt.Errorf("failed to create user directory: %v", err)
	}

	// 构建文件路径
	filePath := filepath.Join(userDir, filename)

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// 写入文件内容
	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file content: %v", err)
	}

	return filePath, nil
}

func (s *LocalStorage) GetAvatarURL(ctx context.Context, uid string, filename string) (string, error) {
	// 构建文件路径
	filePath := filepath.Join(s.basePath, uid, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("avatar file does not exist: %v", err)
	}

	// 返回文件路径
	return filePath, nil
}
