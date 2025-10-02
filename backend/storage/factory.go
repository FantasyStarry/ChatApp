package storage

import (
	"fmt"
)

// StorageType 存储类型常量
const (
	StorageTypeMinio = "minio"
	StorageTypeQiniu = "qiniu"
)

// StorageFactory 存储工厂
type StorageFactory struct{}

// NewStorageFactory 创建存储工厂实例
func NewStorageFactory() *StorageFactory {
	return &StorageFactory{}
}

// CreateStorage 根据配置创建存储实例
func (f *StorageFactory) CreateStorage(storageType string, config interface{}) (Storage, error) {
	switch storageType {
	case StorageTypeMinio:
		minioConfig, ok := config.(MinioConfig)
		if !ok {
			return nil, fmt.Errorf("invalid minio config type")
		}
		return NewMinioStorage(minioConfig)

	case StorageTypeQiniu:
		qiniuConfig, ok := config.(QiniuStorageConfig)
		if !ok {
			return nil, fmt.Errorf("invalid qiniu config type")
		}
		return NewQiniuStorage(qiniuConfig)

	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// GetSupportedStorageTypes 获取支持的存储类型列表
func (f *StorageFactory) GetSupportedStorageTypes() []string {
	return []string{StorageTypeMinio, StorageTypeQiniu}
}

// IsValidStorageType 检查存储类型是否有效
func (f *StorageFactory) IsValidStorageType(storageType string) bool {
	supportedTypes := f.GetSupportedStorageTypes()
	for _, t := range supportedTypes {
		if t == storageType {
			return true
		}
	}
	return false
}
