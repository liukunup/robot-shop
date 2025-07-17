package repository

import (
	v1 "backend/api/v1"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

const (
	basePath    = "avatars"
	minioPrefix = "minio://"
	localPrefix = "local://"
	httpPrefix  = "http://"
	httpsPrefix = "https://"
)

type AvatarStorage interface {
	SaveToMinIO(ctx context.Context, req *v1.AvatarRequest, reader io.Reader) (string, error)
	SaveToLocal(ctx context.Context, req *v1.AvatarRequest, reader io.Reader) (string, error)
	GetURL(ctx context.Context, avatar string) (string, error)
}

func NewAvatarStorage(
	repository *Repository,
) AvatarStorage {
	return &avatarStorage{
		Repository: repository,
	}
}

type avatarStorage struct {
	*Repository
}

func (r *avatarStorage) objectName(uid uint, filename string) string {
	return fmt.Sprintf("%s/%d/%s", basePath, uid, filename)
}

func (r *avatarStorage) SaveToMinIO(ctx context.Context, req *v1.AvatarRequest, reader io.Reader) (string, error) {
	// 检查 Bucket 是否存在
	exists, err := r.m.client.BucketExists(ctx, r.m.bucket)
	if err != nil {
		return "", fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("bucket does not exist")
	}

	// 构建 Object 名称
	objectName := r.objectName(req.UserID, req.Filename)

	// 上传文件到 MinIO
	_, err = r.m.client.PutObject(ctx, r.m.bucket, objectName, reader, req.Size,
		minio.PutObjectOptions{
			ContentType: req.Type,
		})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	return fmt.Sprintf("%s%s/%s", minioPrefix, r.m.bucket, objectName), nil
}

func (r *avatarStorage) SaveToLocal(ctx context.Context, req *v1.AvatarRequest, reader io.Reader) (string, error) {
	// 创建用户目录
	userDir := filepath.Join(basePath, fmt.Sprintf("%d", req.UserID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create user directory: %w", err)
	}

	// 构建文件路径
	filePath := filepath.Join(userDir, req.Filename)

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 写入文件内容
	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file content: %w", err)
	}

	return fmt.Sprintf("%s%s", localPrefix, filePath), nil
}

func (r *avatarStorage) GetURL(ctx context.Context, avatar string) (string, error) {
	// MinIO
	if strings.HasPrefix(avatar, minioPrefix) {
		noPrefix := strings.TrimPrefix(avatar, minioPrefix)
		url := fmt.Sprintf("%s%s", r.m.client.EndpointURL(), noPrefix)
		return url, nil
	}

	// Local
	if strings.HasPrefix(avatar, localPrefix) {
		noPrefix := strings.TrimPrefix(avatar, localPrefix)
		return fmt.Sprintf("%s%s", "http://localhost:8000", noPrefix), nil
	}

	// Http & Https
	if strings.HasPrefix(avatar, httpPrefix) || strings.HasPrefix(avatar, httpsPrefix) {
		return avatar, nil
	}

	return avatar, fmt.Errorf("unknown avatar prefix")
}
