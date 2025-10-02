package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuStorage 七牛云存储实现
type QiniuStorage struct {
	mac       *qbox.Mac
	cfg       *storage.Config
	bucket    string
	domain    string
	useHTTPS  bool
	bucketMgr *storage.BucketManager
	uploader  *storage.FormUploader
}

// QiniuStorageConfig 七牛云配置
type QiniuStorageConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
	Region    string // "south-china", "east-china", "north-china", "north-america", "southeast-asia"
	UseHTTPS  bool
}

// NewQiniuStorage 创建七牛云存储实例
func NewQiniuStorage(config QiniuStorageConfig) (*QiniuStorage, error) {
	if config.AccessKey == "" || config.SecretKey == "" {
		return nil, fmt.Errorf("qiniu access key and secret key are required")
	}
	if config.Bucket == "" {
		return nil, fmt.Errorf("qiniu bucket is required")
	}
	if config.Domain == "" {
		return nil, fmt.Errorf("qiniu domain is required")
	}

	mac := qbox.NewMac(config.AccessKey, config.SecretKey)

	cfg := &storage.Config{
		UseHTTPS:      config.UseHTTPS,
		UseCdnDomains: false,
	}

	// 根据区域设置
	switch config.Region {
	case "east-china":
		cfg.Region = &storage.ZoneHuadong
	case "north-china":
		cfg.Region = &storage.ZoneHuabei
	case "south-china":
		cfg.Region = &storage.ZoneHuanan
	case "north-america":
		cfg.Region = &storage.ZoneBeimei
	case "southeast-asia":
		cfg.Region = &storage.ZoneXinjiapo
	default:
		cfg.Region = &storage.ZoneHuanan // 默认华南
	}

	bucketMgr := storage.NewBucketManager(mac, cfg)
	uploader := storage.NewFormUploader(cfg)

	return &QiniuStorage{
		mac:       mac,
		cfg:       cfg,
		bucket:    config.Bucket,
		domain:    config.Domain,
		useHTTPS:  config.UseHTTPS,
		bucketMgr: bucketMgr,
		uploader:  uploader,
	}, nil
}

// Upload 上传文件
func (q *QiniuStorage) Upload(objectPath string, reader io.Reader, options UploadOptions) (*UploadResult, error) {
	// 生成上传凭证
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", q.bucket, objectPath),
	}
	upToken := putPolicy.UploadToken(q.mac)

	// 读取文件内容到buffer
	var buf bytes.Buffer
	size, err := io.Copy(&buf, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// 上传文件
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	if options.ContentType != "" {
		putExtra.MimeType = options.ContentType
	}

	err = q.uploader.Put(context.Background(), &ret, upToken, objectPath, bytes.NewReader(buf.Bytes()), size, &putExtra)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to qiniu: %w", err)
	}

	return &UploadResult{
		ObjectPath: objectPath,
		Size:       size,
		ETag:       ret.Hash,
	}, nil
}

// Download 获取文件下载URL
func (q *QiniuStorage) Download(objectPath string, expiry time.Duration) (string, error) {
	// 构建访问URL
	scheme := "http"
	if q.useHTTPS {
		scheme = "https"
	}

	publicAccessURL := storage.MakePublicURL(q.domain, objectPath)

	// 如果设置了过期时间，生成私有访问URL
	if expiry > 0 {
		deadline := time.Now().Add(expiry).Unix()
		privateAccessURL := storage.MakePrivateURL(q.mac, q.domain, objectPath, deadline)
		return privateAccessURL, nil
	}

	// 确保URL使用正确的协议
	if u, err := url.Parse(publicAccessURL); err == nil {
		u.Scheme = scheme
		return u.String(), nil
	}

	return publicAccessURL, nil
}

// Delete 删除文件
func (q *QiniuStorage) Delete(objectPath string) error {
	err := q.bucketMgr.Delete(q.bucket, objectPath)
	if err != nil {
		return fmt.Errorf("failed to delete file from qiniu: %w", err)
	}
	return nil
}

// GetFileInfo 获取文件信息
func (q *QiniuStorage) GetFileInfo(objectPath string) (*FileInfo, error) {
	fileInfo, err := q.bucketMgr.Stat(q.bucket, objectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info from qiniu: %w", err)
	}

	return &FileInfo{
		Key:          objectPath,
		Size:         fileInfo.Fsize,
		LastModified: time.Unix(fileInfo.PutTime/10000000, 0), // 七牛时间戳是纳秒，需要转换
		ContentType:  fileInfo.MimeType,
		ETag:         fileInfo.Hash,
	}, nil
}

// GetUploadURL 获取预签名上传URL
func (q *QiniuStorage) GetUploadURL(objectPath string, expiry time.Duration) (string, error) {
	// 七牛云的预签名上传通过上传凭证实现
	putPolicy := storage.PutPolicy{
		Scope:   fmt.Sprintf("%s:%s", q.bucket, objectPath),
		Expires: uint64(time.Now().Add(expiry).Unix()),
	}
	upToken := putPolicy.UploadToken(q.mac)

	// 返回上传地址和token
	// 这里返回的是上传凭证，前端需要使用这个凭证进行上传
	scheme := "http"
	if q.useHTTPS {
		scheme = "https"
	}

	uploadHost := fmt.Sprintf("%s://upload.qiniup.com", scheme)
	if q.cfg.Region != nil {
		uploadHost = fmt.Sprintf("%s://up.qiniup.com", scheme) // 使用通用上传地址
	}

	// 构建包含token的上传URL（这是一个简化实现，实际使用时前端需要处理token）
	uploadURL := fmt.Sprintf("%s?token=%s&key=%s", uploadHost, upToken, url.QueryEscape(objectPath))

	return uploadURL, nil
}

// Exists 检查文件是否存在
func (q *QiniuStorage) Exists(objectPath string) (bool, error) {
	_, err := q.bucketMgr.Stat(q.bucket, objectPath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") ||
			strings.Contains(err.Error(), "612") { // 七牛云文件不存在错误码
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}
	return true, nil
}

// GetStorageType 获取存储类型
func (q *QiniuStorage) GetStorageType() string {
	return "qiniu"
}
