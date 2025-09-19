package controllers

import (
	"chatapp/service"
	"chatapp/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService *service.FileService
}

func NewFileController(fileService *service.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传文件到指定聊天室
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param chatroom_id formData int true "聊天室ID"
// @Param file formData file true "要上传的文件"
// @Success 200 {object} utils.Response{data=models.File}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/files/upload [post]
func (c *FileController) UploadFile(ctx *gin.Context) {
	// 获取用户ID（从JWT中间件获取）
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "用户未认证")
		return
	}

	// 获取聊天室ID
	chatRoomIDStr := ctx.PostForm("chatroom_id")
	if chatRoomIDStr == "" {
		utils.BadRequestResponse(ctx, "聊天室ID不能为空")
		return
	}

	chatRoomID, err := strconv.ParseUint(chatRoomIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "无效的聊天室ID")
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.BadRequestResponse(ctx, "获取文件失败: "+err.Error())
		return
	}

	// 检查文件大小（限制50MB）
	const maxFileSize = 50 * 1024 * 1024 // 50MB
	if file.Size > maxFileSize {
		utils.BadRequestResponse(ctx, "文件大小不能超过50MB")
		return
	}

	// 上传文件
	fileRecord, err := c.fileService.UploadFile(file, uint(chatRoomID), userID.(uint))
	if err != nil {
		utils.InternalErrorResponse(ctx, "文件上传失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(ctx, "文件上传成功", fileRecord)
}

// DownloadFile 下载文件
// @Summary 下载文件
// @Description 获取文件下载链接
// @Tags files
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} utils.Response{data=map[string]interface{}}
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/files/download/{id} [get]
func (c *FileController) DownloadFile(ctx *gin.Context) {
	// 获取文件ID
	fileIDStr := ctx.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "无效的文件ID")
		return
	}

	// 获取下载链接
	downloadURL, fileInfo, err := c.fileService.DownloadFile(uint(fileID))
	if err != nil {
		utils.NotFoundResponse(ctx, "文件不存在或获取失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(ctx, "获取下载链接成功", gin.H{
		"download_url": downloadURL,
		"file_info":    fileInfo,
	})
}

// GetFilesByRoom 获取聊天室文件列表
// @Summary 获取聊天室文件列表
// @Description 获取指定聊天室的文件列表（支持分页）
// @Tags files
// @Produce json
// @Param chatroom_id path int true "聊天室ID"
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认20，最大100）"
// @Success 200 {object} utils.Response{data=map[string]interface{}}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/files/chatroom/{chatroom_id} [get]
func (c *FileController) GetFilesByRoom(ctx *gin.Context) {
	// 获取聊天室ID
	chatRoomIDStr := ctx.Param("chatroom_id")
	chatRoomID, err := strconv.ParseUint(chatRoomIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "无效的聊天室ID")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	// 获取文件列表
	files, total, err := c.fileService.GetFilesByRoomWithPagination(uint(chatRoomID), page, pageSize)
	if err != nil {
		utils.InternalErrorResponse(ctx, "获取文件列表失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(ctx, "获取文件列表成功", gin.H{
		"files":       files,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetFilesByUser 获取用户上传的文件列表
// @Summary 获取用户文件列表
// @Description 获取当前用户上传的所有文件列表
// @Tags files
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.File}
// @Failure 500 {object} utils.Response
// @Router /api/files/my [get]
func (c *FileController) GetFilesByUser(ctx *gin.Context) {
	// 获取用户ID（从JWT中间件获取）
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "用户未认证")
		return
	}

	// 获取用户文件列表
	files, err := c.fileService.GetFilesByUser(userID.(uint))
	if err != nil {
		utils.InternalErrorResponse(ctx, "获取文件列表失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(ctx, "获取用户文件列表成功", files)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除指定文件（只有上传者可以删除）
// @Tags files
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/files/{id} [delete]
func (c *FileController) DeleteFile(ctx *gin.Context) {
	// 获取用户ID（从JWT中间件获取）
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "用户未认证")
		return
	}

	// 获取文件ID
	fileIDStr := ctx.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "无效的文件ID")
		return
	}

	// 删除文件
	err = c.fileService.DeleteFile(uint(fileID), userID.(uint))
	if err != nil {
		if err.Error() == "permission denied: only uploader can delete the file" {
			utils.ForbiddenResponse(ctx, "权限不足：只有上传者可以删除文件")
			return
		}
		utils.InternalErrorResponse(ctx, "删除文件失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(ctx, "文件删除成功", nil)
}

// GetFileInfo 获取文件信息
// @Summary 获取文件信息
// @Description 获取指定文件的详细信息
// @Tags files
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} utils.Response{data=models.File}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/files/{id} [get]
func (c *FileController) GetFileInfo(ctx *gin.Context) {
	// 获取文件ID
	fileIDStr := ctx.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "无效的文件ID")
		return
	}

	// 获取文件信息
	fileInfo, err := c.fileService.GetFileInfo(uint(fileID))
	if err != nil {
		utils.NotFoundResponse(ctx, "文件不存在: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(ctx, "获取文件信息成功", fileInfo)
}

// GetUploadURL 获取文件上传预签名URL（可选功能）
// @Summary 获取上传预签名URL
// @Description 获取文件上传的预签名URL，用于前端直接上传到Minio
// @Tags files
// @Produce json
// @Param filename query string true "文件名"
// @Param chatroom_id query int true "聊天室ID"
// @Success 200 {object} utils.Response{data=map[string]string}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/files/upload-url [get]
func (c *FileController) GetUploadURL(ctx *gin.Context) {
	// 获取参数
	fileName := ctx.Query("filename")
	chatRoomIDStr := ctx.Query("chatroom_id")

	if fileName == "" || chatRoomIDStr == "" {
		utils.BadRequestResponse(ctx, "文件名和聊天室ID不能为空")
		return
	}

	chatRoomID, err := strconv.ParseUint(chatRoomIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "无效的聊天室ID")
		return
	}

	// 获取预签名上传URL
	uploadURL, objectPath, err := c.fileService.GetUploadURL(fileName, uint(chatRoomID))
	if err != nil {
		utils.InternalErrorResponse(ctx, "获取上传URL失败: "+err.Error())
		return
	}

	utils.SuccessResponseWithMessage(ctx, "获取上传URL成功", gin.H{
		"upload_url":  uploadURL,
		"object_path": objectPath,
	})
}