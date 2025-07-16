package repository

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	basePath        = "avatars"
	presignedExpiry = 15 * time.Minute
)

type AvatarStorage interface {
	PresignedUploadURL(ctx context.Context, uid uint, filename string) (string, error)
	PresignedViewURL(ctx context.Context, uid uint, filename string) (string, error)
	SaveToLocal(ctx context.Context, uid uint, filename string, reader io.Reader) (string, error)
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

func (r *avatarStorage) objectKey(uid uint, filename string) string {
	return fmt.Sprintf("%s/%d/%s", basePath, uid, filename)
}

func (r *avatarStorage) bucketExists(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exists, err := r.m.client.BucketExists(ctx, r.m.bucket)
	if err != nil {
		return false, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	return exists, nil
}

func (r *avatarStorage) PresignedUploadURL(ctx context.Context, uid uint, filename string) (string, error) {
	// 检查桶是否存在
	if exist, _ := r.bucketExists(ctx); !exist {
		return "", fmt.Errorf("bucket does not exist")
	}

	// 生成预签名上传链接
	objectKey := r.objectKey(uid, filename)
	url, err := r.m.client.PresignedPutObject(ctx, r.m.bucket, objectKey, presignedExpiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return url.String(), nil
}

func (r *avatarStorage) PresignedViewURL(ctx context.Context, uid uint, filename string) (string, error) {
	// 检查桶是否存在
	if exist, _ := r.bucketExists(ctx); !exist {
		return "", fmt.Errorf("bucket does not exist")
	}

	// 生成预签名下载链接
	objectKey := r.objectKey(uid, filename)
	url, err := r.m.client.PresignedGetObject(ctx, r.m.bucket, objectKey, presignedExpiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	return url.String(), nil
}

func (r *avatarStorage) SaveToLocal(ctx context.Context, uid uint, filename string, reader io.Reader) (string, error) {
	// 创建用户目录
	userDir := filepath.Join(basePath, fmt.Sprintf("%d", uid))
	if err := os.MkdirAll(userDir, 0744); err != nil {
		return "", fmt.Errorf("failed to create user directory: %w", err)
	}

	// 构建文件路径
	filePath := filepath.Join(userDir, filename)

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

	return filePath, nil
}
