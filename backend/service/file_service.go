package service

import (
	"chatapp/config"
	"chatapp/models"
	"chatapp/repository"
	"chatapp/storage"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type FileService struct {
	fileRepo *repository.FileRepository
	storage  storage.Storage
}

func NewFileService(fileRepo *repository.FileRepository) *FileService {
	// 创建存储工厂
	factory := storage.NewStorageFactory()

	// 根据配置选择存储类型
	storageType := config.GlobalConfig.Storage.Type
	if storageType == "" {
		storageType = storage.StorageTypeMinio // 默认使用MinIO
	}

	// 验证存储类型
	if !factory.IsValidStorageType(storageType) {
		panic(fmt.Sprintf("Unsupported storage type: %s", storageType))
	}

	var storageInstance storage.Storage
	var err error

	// 根据类型创建存储实例
	switch storageType {
	case storage.StorageTypeMinio:
		minioConfig := storage.MinioConfig{
			Endpoint:   config.GlobalConfig.Minio.Endpoint,
			AccessKey:  config.GlobalConfig.Minio.AccessKey,
			SecretKey:  config.GlobalConfig.Minio.SecretKey,
			BucketName: config.GlobalConfig.Minio.BucketName,
			UseSSL:     config.GlobalConfig.Minio.UseSSL,
			Region:     config.GlobalConfig.Minio.Region,
		}
		storageInstance, err = factory.CreateStorage(storageType, minioConfig)

	case storage.StorageTypeQiniu:
		qiniuConfig := storage.QiniuStorageConfig{
			AccessKey: config.GlobalConfig.Qiniu.AccessKey,
			SecretKey: config.GlobalConfig.Qiniu.SecretKey,
			Bucket:    config.GlobalConfig.Qiniu.Bucket,
			Domain:    config.GlobalConfig.Qiniu.Domain,
			Region:    config.GlobalConfig.Qiniu.Region,
			UseHTTPS:  config.GlobalConfig.Qiniu.UseHTTPS,
		}
		storageInstance, err = factory.CreateStorage(storageType, qiniuConfig)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize storage: %v", err))
	}

	return &FileService{
		fileRepo: fileRepo,
		storage:  storageInstance,
	}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(file *multipart.FileHeader, chatRoomID, uploaderID uint) (*models.File, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// 生成唯一的文件路径，按聊天室分组
	timestamp := time.Now().Unix()
	fileExt := filepath.Ext(file.Filename)
	fileName := strings.TrimSuffix(file.Filename, fileExt)
	objectPath := fmt.Sprintf("chatroom-%d/%d-%s%s", chatRoomID, timestamp, fileName, fileExt)

	// 上传到存储
	uploadOptions := storage.UploadOptions{
		ContentType: file.Header.Get("Content-Type"),
		Size:        file.Size,
	}

	uploadResult, err := s.storage.Upload(objectPath, src, uploadOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// 创建数据库记录
	fileRecord := &models.File{
		FileName:    file.Filename,
		FilePath:    uploadResult.ObjectPath,
		FileSize:    uploadResult.Size,
		ContentType: file.Header.Get("Content-Type"),
		ChatRoomID:  chatRoomID,
		UploaderID:  uploaderID,
	}

	err = s.fileRepo.Create(fileRecord)
	if err != nil {
		// 如果数据库插入失败，尝试删除已上传的文件
		s.storage.Delete(objectPath)
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return fileRecord, nil
}

// GetFileInfo 获取文件信息
func (s *FileService) GetFileInfo(fileID uint) (*models.File, error) {
	return s.fileRepo.GetByID(fileID)
}

// DownloadFile 获取文件下载链接
func (s *FileService) DownloadFile(fileID uint) (string, *models.File, error) {
	// 获取文件记录
	fileRecord, err := s.fileRepo.GetByID(fileID)
	if err != nil {
		return "", nil, fmt.Errorf("file not found: %w", err)
	}

	// 生成下载URL（有效期1小时）
	downloadURL, err := s.storage.Download(fileRecord.FilePath, time.Hour)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate download url: %w", err)
	}

	return downloadURL, fileRecord, nil
}

// GetFilesByRoom 获取聊天室文件列表
func (s *FileService) GetFilesByRoom(chatRoomID uint) ([]models.File, error) {
	return s.fileRepo.GetByChatRoomID(chatRoomID)
}

// GetFilesByRoomWithPagination 分页获取聊天室文件列表
func (s *FileService) GetFilesByRoomWithPagination(chatRoomID uint, page, pageSize int) ([]models.File, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return s.fileRepo.GetByChatRoomIDWithPagination(chatRoomID, page, pageSize)
}

// GetFilesByUser 获取用户上传的文件列表
func (s *FileService) GetFilesByUser(userID uint) ([]models.File, error) {
	return s.fileRepo.GetByUserID(userID)
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(fileID, userID uint) error {
	// 获取文件记录
	fileRecord, err := s.fileRepo.GetByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// 检查权限（只有上传者可以删除）
	if fileRecord.UploaderID != userID {
		return fmt.Errorf("permission denied: only uploader can delete the file")
	}

	// 从存储删除文件
	err = s.storage.Delete(fileRecord.FilePath)
	if err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// 从数据库删除记录
	err = s.fileRepo.Delete(fileID)
	if err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

// GetUploadURL 获取文件上传的预签名URL（可选功能，用于前端直接上传）
func (s *FileService) GetUploadURL(fileName string, chatRoomID uint) (string, string, error) {
	// 生成对象路径
	timestamp := time.Now().Unix()
	fileExt := filepath.Ext(fileName)
	baseFileName := strings.TrimSuffix(fileName, fileExt)
	objectPath := fmt.Sprintf("chatroom-%d/%d-%s%s", chatRoomID, timestamp, baseFileName, fileExt)

	// 生成预签名上传URL（有效期15分钟）
	uploadURL, err := s.storage.GetUploadURL(objectPath, 15*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate upload url: %w", err)
	}

	return uploadURL, objectPath, nil
}
