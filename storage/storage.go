package storage

import (
	"io"
	"time"
)

// UploadOptions 上传选项
type UploadOptions struct {
	ContentType string
	Size        int64
}

// UploadResult 上传结果
type UploadResult struct {
	ObjectPath string
	Size       int64
	ETag       string
}

// FileInfo 文件信息
type FileInfo struct {
	Key          string
	Size         int64
	LastModified time.Time
	ContentType  string
	ETag         string
}

// Storage 存储接口，定义文件存储的统一抽象
type Storage interface {
	// Upload 上传文件
	Upload(objectPath string, reader io.Reader, options UploadOptions) (*UploadResult, error)

	// Download 获取文件下载URL
	Download(objectPath string, expiry time.Duration) (string, error)

	// Delete 删除文件
	Delete(objectPath string) error

	// GetFileInfo 获取文件信息
	GetFileInfo(objectPath string) (*FileInfo, error)

	// GetUploadURL 获取预签名上传URL（可选实现）
	GetUploadURL(objectPath string, expiry time.Duration) (string, error)

	// Exists 检查文件是否存在
	Exists(objectPath string) (bool, error)

	// GetStorageType 获取存储类型
	GetStorageType() string
}
