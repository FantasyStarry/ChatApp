package service

import (
	"chatapp/config"
	"chatapp/models"
	"chatapp/repository"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileService struct {
	fileRepo    *repository.FileRepository
	minioClient *minio.Client
	bucketName  string
}

func NewFileService(fileRepo *repository.FileRepository) *FileService {
	// 初始化Minio客户端
	minioClient, err := minio.New(config.GlobalConfig.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.GlobalConfig.Minio.AccessKey, config.GlobalConfig.Minio.SecretKey, ""),
		Secure: config.GlobalConfig.Minio.UseSSL,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Minio client: %v", err))
	}

	// 确保bucket存在
	ctx := context.Background()
	bucketName := config.GlobalConfig.Minio.BucketName
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		panic(fmt.Sprintf("Failed to check bucket existence: %v", err))
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: config.GlobalConfig.Minio.Region})
		if err != nil {
			panic(fmt.Sprintf("Failed to create bucket: %v", err))
		}
	}

	return &FileService{
		fileRepo:    fileRepo,
		minioClient: minioClient,
		bucketName:  bucketName,
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

	// 上传到Minio
	ctx := context.Background()
	uploadInfo, err := s.minioClient.PutObject(ctx, s.bucketName, objectPath, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to minio: %w", err)
	}

	// 创建数据库记录
	fileRecord := &models.File{
		FileName:    file.Filename,
		FilePath:    objectPath,
		FileSize:    uploadInfo.Size,
		ContentType: file.Header.Get("Content-Type"),
		ChatRoomID:  chatRoomID,
		UploaderID:  uploaderID,
	}

	err = s.fileRepo.Create(fileRecord)
	if err != nil {
		// 如果数据库插入失败，尝试删除已上传的文件
		s.minioClient.RemoveObject(ctx, s.bucketName, objectPath, minio.RemoveObjectOptions{})
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

	// 生成预签名下载URL（有效期1小时）
	ctx := context.Background()
	presignedURL, err := s.minioClient.PresignedGetObject(ctx, s.bucketName, fileRecord.FilePath, time.Hour, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate download url: %w", err)
	}

	return presignedURL.String(), fileRecord, nil
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

	// 从Minio删除文件
	ctx := context.Background()
	err = s.minioClient.RemoveObject(ctx, s.bucketName, fileRecord.FilePath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from minio: %w", err)
	}

	// 从数据库删除记录
	err = s.fileRepo.Delete(fileID)
	if err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

// GetUploadURL 获取文件上传的预签名URL（可选功能，用于前端直接上传到Minio）
func (s *FileService) GetUploadURL(fileName string, chatRoomID uint) (string, string, error) {
	// 生成对象路径
	timestamp := time.Now().Unix()
	fileExt := filepath.Ext(fileName)
	baseFileName := strings.TrimSuffix(fileName, fileExt)
	objectPath := fmt.Sprintf("chatroom-%d/%d-%s%s", chatRoomID, timestamp, baseFileName, fileExt)

	// 生成预签名上传URL（有效期15分钟）
	ctx := context.Background()
	presignedURL, err := s.minioClient.PresignedPutObject(ctx, s.bucketName, objectPath, 15*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate upload url: %w", err)
	}

	return presignedURL.String(), objectPath, nil
}