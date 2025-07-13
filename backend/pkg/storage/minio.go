package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type MinIOConfig struct {
	Endpoint  string // 服务端地址
	AccessKey string // 密钥
	SecretKey string // 密钥
	Bucket    string // 桶
	Region    string // 区域
	Secure    bool   // 是否使用 HTTPS
}

type MinIOStorage struct {
	client *minio.Client // 客户端
	bucket string        // 桶
}

func NewMinIOStorage(conf *viper.Viper) Storage {
	cfg := &MinIOConfig{
		Endpoint:  conf.GetString("storage.minio.endpoint"),
		AccessKey: conf.GetString("storage.minio.access_key"),
		SecretKey: conf.GetString("storage.minio.secret_key"),
		Bucket:    conf.GetString("storage.minio.bucket"),
		Region:    conf.GetString("storage.minio.region"),
		Secure:    conf.GetBool("storage.minio.secure"),
	}

	// 初始化 MinIO 客户端
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Region: cfg.Region,
		Secure: cfg.Secure,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to initialize MinIO client: %v", err))
	}

	// 检查桶是否存在，不存在则创建
	exists, err := client.BucketExists(context.Background(), cfg.Bucket)
	if err != nil {
		panic(fmt.Sprintf("failed to check bucket existence: %v", err))
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			panic(fmt.Sprintf("failed to create bucket: %v", err))
		}
	}

	return &MinIOStorage{
		client: client,
		bucket: cfg.Bucket,
	}
}

func (s *MinIOStorage) UploadAvatar(ctx context.Context, uid string, filename string, reader io.Reader) (string, error) {
	// 构建对象键
	objectKey := fmt.Sprintf("avatars/%s/%s", uid, filename)

	// 获取文件大小
	fileInfo, ok := reader.(interface{ Stat() (os.FileInfo, error) })
	if !ok {
		return "", fmt.Errorf("failed to get file size")
	}
	fileSize, err := fileInfo.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file size: %v", err)
	}

	// 上传文件
	_, err = s.client.PutObject(ctx, s.bucket, objectKey, reader, fileSize.Size(), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %v", err)
	}

	// 返回对象键
	return objectKey, nil
}

// GetAvatarURL 获取MinIO上的头像URL
func (s *MinIOStorage) GetAvatarURL(ctx context.Context, uid string, filename string) (string, error) {
	// 构建对象键（类似文件路径）
	objectKey := fmt.Sprintf("avatars/%s/%s", uid, filename)

	// 返回访问URL（根据实际情况修改）
	return fmt.Sprintf("/minio/%s/%s", s.bucket, objectKey), nil
}
