package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage MinIO存储实现
type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

// MinioConfig MinIO配置
type MinioConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
	Region     string
}

// NewMinioStorage 创建MinIO存储实例
func NewMinioStorage(config MinioConfig) (*MinioStorage, error) {
	// 初始化Minio客户端
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Minio client: %w", err)
	}

	// 确保bucket存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{Region: config.Region})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinioStorage{
		client:     client,
		bucketName: config.BucketName,
	}, nil
}

// Upload 上传文件
func (m *MinioStorage) Upload(objectPath string, reader io.Reader, options UploadOptions) (*UploadResult, error) {
	ctx := context.Background()
	uploadInfo, err := m.client.PutObject(ctx, m.bucketName, objectPath, reader, options.Size, minio.PutObjectOptions{
		ContentType: options.ContentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to minio: %w", err)
	}

	return &UploadResult{
		ObjectPath: objectPath,
		Size:       uploadInfo.Size,
		ETag:       uploadInfo.ETag,
	}, nil
}

// Download 获取文件下载URL
func (m *MinioStorage) Download(objectPath string, expiry time.Duration) (string, error) {
	ctx := context.Background()
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectPath, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate download url: %w", err)
	}
	return presignedURL.String(), nil
}

// Delete 删除文件
func (m *MinioStorage) Delete(objectPath string) error {
	ctx := context.Background()
	err := m.client.RemoveObject(ctx, m.bucketName, objectPath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from minio: %w", err)
	}
	return nil
}

// GetFileInfo 获取文件信息
func (m *MinioStorage) GetFileInfo(objectPath string) (*FileInfo, error) {
	ctx := context.Background()
	objInfo, err := m.client.StatObject(ctx, m.bucketName, objectPath, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info from minio: %w", err)
	}

	return &FileInfo{
		Key:          objInfo.Key,
		Size:         objInfo.Size,
		LastModified: objInfo.LastModified,
		ContentType:  objInfo.ContentType,
		ETag:         objInfo.ETag,
	}, nil
}

// GetUploadURL 获取预签名上传URL
func (m *MinioStorage) GetUploadURL(objectPath string, expiry time.Duration) (string, error) {
	ctx := context.Background()
	presignedURL, err := m.client.PresignedPutObject(ctx, m.bucketName, objectPath, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload url: %w", err)
	}
	return presignedURL.String(), nil
}

// Exists 检查文件是否存在
func (m *MinioStorage) Exists(objectPath string) (bool, error) {
	ctx := context.Background()
	_, err := m.client.StatObject(ctx, m.bucketName, objectPath, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}
	return true, nil
}

// GetStorageType 获取存储类型
func (m *MinioStorage) GetStorageType() string {
	return "minio"
}
