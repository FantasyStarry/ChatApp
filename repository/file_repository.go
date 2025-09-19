package repository

import (
	"chatapp/models"

	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

// Create 创建文件记录
func (r *FileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

// GetByID 根据ID获取文件
func (r *FileRepository) GetByID(id uint) (*models.File, error) {
	var file models.File
	err := r.db.Preload("ChatRoom").Preload("Uploader").First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// GetByChatRoomID 根据聊天室ID获取文件列表
func (r *FileRepository) GetByChatRoomID(chatRoomID uint) ([]models.File, error) {
	var files []models.File
	err := r.db.Where("chatroom_id = ?", chatRoomID).
		Preload("Uploader").
		Order("uploaded_at DESC").
		Find(&files).Error
	return files, err
}

// GetByUserID 根据用户ID获取用户上传的文件列表
func (r *FileRepository) GetByUserID(userID uint) ([]models.File, error) {
	var files []models.File
	err := r.db.Where("uploader_id = ?", userID).
		Preload("ChatRoom").
		Order("uploaded_at DESC").
		Find(&files).Error
	return files, err
}

// Delete 软删除文件记录
func (r *FileRepository) Delete(id uint) error {
	return r.db.Delete(&models.File{}, id).Error
}

// Update 更新文件信息
func (r *FileRepository) Update(file *models.File) error {
	return r.db.Save(file).Error
}

// GetByChatRoomIDWithPagination 分页获取聊天室文件列表
func (r *FileRepository) GetByChatRoomIDWithPagination(chatRoomID uint, page, pageSize int) ([]models.File, int64, error) {
	var files []models.File
	var total int64
	
	// 获取总数
	err := r.db.Model(&models.File{}).Where("chatroom_id = ?", chatRoomID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.Where("chatroom_id = ?", chatRoomID).
		Preload("Uploader").
		Order("uploaded_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&files).Error
	
	return files, total, err
}

// GetByFilePath 根据文件路径获取文件记录
func (r *FileRepository) GetByFilePath(filePath string) (*models.File, error) {
	var file models.File
	err := r.db.Where("file_path = ?", filePath).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}