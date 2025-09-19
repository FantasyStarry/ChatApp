package models

import (
	"time"

	"gorm.io/gorm"
)

// File represents a file uploaded to the system
type File struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	FileName     string    `json:"file_name" gorm:"not null;size:255"`           // 原始文件名
	FilePath     string    `json:"file_path" gorm:"not null;size:500"`           // Minio中的对象路径
	FileSize     int64     `json:"file_size" gorm:"not null"`                    // 文件大小（字节）
	ContentType  string    `json:"content_type" gorm:"size:100"`                 // MIME类型
	ChatRoomID   uint      `json:"chat_room_id" gorm:"column:chat_room_id;not null;index"`    // 所属聊天室ID
	UploaderID   uint      `json:"uploader_id" gorm:"column:uploader_id;not null;index"`      // 上传用户ID
	UploadedAt   time.Time `json:"uploaded_at" gorm:"autoCreateTime"`            // 上传时间
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	ChatRoom ChatRoom `json:"chatroom,omitempty" gorm:"foreignKey:ChatRoomID"`
	Uploader User     `json:"uploader,omitempty" gorm:"foreignKey:UploaderID"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// BeforeCreate 创建前的钩子函数
func (f *File) BeforeCreate(tx *gorm.DB) error {
	f.UploadedAt = time.Now()
	return nil
}