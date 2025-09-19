# 前端文件管理功能集成指南

## 概述

本文档描述了如何在前端应用中集成 ChatApp 的文件管理功能。ChatApp 现在支持文件上传、下载、列表查看和删除等功能，所有文件都按聊天室进行组织存储。

## 文件存储架构

- **存储服务**: Minio
- **组织方式**: 按聊天室分组存储（路径格式：`chatroom-{id}/文件名`）
- **文件大小限制**: 50MB
- **文件类型**: 无限制
- **权限控制**: 只有上传者可以删除文件

## API 端点概览

| 操作           | 方法   | 端点                       | 描述                       |
| -------------- | ------ | -------------------------- | -------------------------- |
| 上传文件       | POST   | `/api/files/upload`        | 上传文件到指定聊天室       |
| 下载文件       | GET    | `/api/files/download/{id}` | 获取文件下载链接           |
| 获取聊天室文件 | GET    | `/api/files/chatroom/{id}` | 获取聊天室文件列表（分页） |
| 获取我的文件   | GET    | `/api/files/my`            | 获取当前用户上传的文件     |
| 删除文件       | DELETE | `/api/files/{id}`          | 删除文件（仅上传者）       |
| 获取文件信息   | GET    | `/api/files/{id}`          | 获取文件详细信息           |
| 获取上传 URL   | GET    | `/api/files/upload-url`    | 获取预签名上传 URL         |

## 前端实现示例

### 1. 文件上传组件

#### HTML 结构

```html
<div class="file-upload-component">
  <input type="file" id="fileInput" multiple />
  <button onclick="uploadFiles()">上传文件</button>
  <div id="uploadProgress" style="display: none;">
    <div class="progress-bar">
      <div class="progress-fill" style="width: 0%"></div>
    </div>
    <span class="progress-text">0%</span>
  </div>
</div>
```

#### JavaScript 实现

```javascript
async function uploadFiles() {
  const fileInput = document.getElementById("fileInput");
  const files = fileInput.files;
  const chatRoomId = getCurrentChatRoomId(); // 获取当前聊天室ID

  if (!files || files.length === 0) {
    alert("请选择文件");
    return;
  }

  for (const file of files) {
    await uploadSingleFile(file, chatRoomId);
  }
}

async function uploadSingleFile(file, chatRoomId) {
  // 检查文件大小
  const maxSize = 50 * 1024 * 1024; // 50MB
  if (file.size > maxSize) {
    alert(`文件 ${file.name} 超过50MB限制`);
    return;
  }

  const formData = new FormData();
  formData.append("file", file);
  formData.append("chatroom_id", chatRoomId);

  try {
    showUploadProgress(true);

    const response = await fetch("/api/files/upload", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${getAuthToken()}`,
      },
      body: formData,
    });

    const result = await response.json();

    if (result.code === 1000) {
      console.log("文件上传成功:", result.data);
      // 刷新文件列表
      await loadChatRoomFiles(chatRoomId);
      // 清空文件输入
      document.getElementById("fileInput").value = "";
    } else {
      alert("上传失败: " + result.messages);
    }
  } catch (error) {
    console.error("上传错误:", error);
    alert("上传失败，请重试");
  } finally {
    showUploadProgress(false);
  }
}

function showUploadProgress(show) {
  const progressDiv = document.getElementById("uploadProgress");
  progressDiv.style.display = show ? "block" : "none";
}
```

### 2. 文件列表组件

#### HTML 结构

```html
<div class="file-list-component">
  <div class="file-list-header">
    <h3>聊天室文件</h3>
    <button onclick="refreshFileList()">刷新</button>
  </div>
  <div id="fileList" class="file-list">
    <!-- 文件项将在这里动态加载 -->
  </div>
  <div class="pagination">
    <button onclick="previousPage()" id="prevBtn">上一页</button>
    <span id="pageInfo">第 1 页</span>
    <button onclick="nextPage()" id="nextBtn">下一页</button>
  </div>
</div>
```

#### JavaScript 实现

```javascript
let currentPage = 1;
const pageSize = 20;
let totalPages = 1;

async function loadChatRoomFiles(chatRoomId, page = 1) {
  try {
    const response = await fetch(
      `/api/files/chatroom/${chatRoomId}?page=${page}&page_size=${pageSize}`,
      {
        headers: {
          Authorization: `Bearer ${getAuthToken()}`,
        },
      }
    );

    const result = await response.json();

    if (result.code === 1000) {
      const { files, total, page: currentPageNum, total_pages } = result.data;
      currentPage = currentPageNum;
      totalPages = total_pages;

      renderFileList(files);
      updatePagination();
    } else {
      console.error("加载文件列表失败:", result.messages);
    }
  } catch (error) {
    console.error("加载文件列表错误:", error);
  }
}

function renderFileList(files) {
  const fileList = document.getElementById("fileList");

  if (files.length === 0) {
    fileList.innerHTML = '<div class="no-files">暂无文件</div>';
    return;
  }

  const fileItems = files.map((file) => createFileItem(file)).join("");
  fileList.innerHTML = fileItems;
}

function createFileItem(file) {
  const fileSize = formatFileSize(file.file_size);
  const uploadTime = formatDate(file.uploaded_at);
  const canDelete = file.uploader_id === getCurrentUserId();

  return `
    <div class="file-item" data-file-id="${file.id}">
      <div class="file-icon">
        <i class="icon ${getFileIcon(file.content_type)}"></i>
      </div>
      <div class="file-info">
        <div class="file-name">${escapeHtml(file.file_name)}</div>
        <div class="file-meta">
          <span class="file-size">${fileSize}</span>
          <span class="uploader">上传者: ${escapeHtml(
            file.uploader.username
          )}</span>
          <span class="upload-time">${uploadTime}</span>
        </div>
      </div>
      <div class="file-actions">
        <button onclick="downloadFile(${
          file.id
        })" class="btn-download">下载</button>
        ${
          canDelete
            ? `<button onclick="deleteFile(${file.id})" class="btn-delete">删除</button>`
            : ""
        }
      </div>
    </div>
  `;
}

function updatePagination() {
  document.getElementById(
    "pageInfo"
  ).textContent = `第 ${currentPage} 页，共 ${totalPages} 页`;
  document.getElementById("prevBtn").disabled = currentPage <= 1;
  document.getElementById("nextBtn").disabled = currentPage >= totalPages;
}

function previousPage() {
  if (currentPage > 1) {
    loadChatRoomFiles(getCurrentChatRoomId(), currentPage - 1);
  }
}

function nextPage() {
  if (currentPage < totalPages) {
    loadChatRoomFiles(getCurrentChatRoomId(), currentPage + 1);
  }
}
```

### 3. 文件下载功能

```javascript
async function downloadFile(fileId) {
  try {
    const response = await fetch(`/api/files/download/${fileId}`, {
      headers: {
        Authorization: `Bearer ${getAuthToken()}`,
      },
    });

    const result = await response.json();

    if (result.code === 1000) {
      const { download_url, file_info } = result.data;

      // 创建隐藏的下载链接
      const link = document.createElement("a");
      link.href = download_url;
      link.download = file_info.file_name;
      link.style.display = "none";

      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } else {
      alert("下载失败: " + result.messages);
    }
  } catch (error) {
    console.error("下载错误:", error);
    alert("下载失败，请重试");
  }
}
```

### 4. 文件删除功能

```javascript
async function deleteFile(fileId) {
  if (!confirm("确定要删除这个文件吗？此操作无法撤销。")) {
    return;
  }

  try {
    const response = await fetch(`/api/files/${fileId}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${getAuthToken()}`,
      },
    });

    const result = await response.json();

    if (result.code === 1000) {
      alert("文件删除成功");
      // 刷新文件列表
      await loadChatRoomFiles(getCurrentChatRoomId(), currentPage);
    } else {
      alert("删除失败: " + result.messages);
    }
  } catch (error) {
    console.error("删除错误:", error);
    alert("删除失败，请重试");
  }
}
```

### 5. 工具函数

```javascript
function formatFileSize(bytes) {
  if (bytes === 0) return "0 Bytes";

  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}

function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleString("zh-CN");
}

function getFileIcon(contentType) {
  if (contentType.startsWith("image/")) return "icon-image";
  if (contentType.startsWith("video/")) return "icon-video";
  if (contentType.startsWith("audio/")) return "icon-audio";
  if (contentType.includes("pdf")) return "icon-pdf";
  if (contentType.includes("word")) return "icon-word";
  if (contentType.includes("excel")) return "icon-excel";
  if (contentType.includes("powerpoint")) return "icon-powerpoint";
  return "icon-file";
}

function escapeHtml(text) {
  const div = document.createElement("div");
  div.textContent = text;
  return div.innerHTML;
}

function getAuthToken() {
  // 从 localStorage 或其他地方获取认证令牌
  return localStorage.getItem("authToken");
}

function getCurrentUserId() {
  // 获取当前用户ID
  const token = getAuthToken();
  if (!token) return null;

  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.user_id;
  } catch (e) {
    return null;
  }
}

function getCurrentChatRoomId() {
  // 获取当前聊天室ID，根据你的应用实现
  return window.currentChatRoomId || 1;
}
```

## CSS 样式示例

```css
.file-upload-component {
  margin: 20px 0;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.file-list-component {
  margin: 20px 0;
}

.file-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid #eee;
  transition: background-color 0.2s;
}

.file-item:hover {
  background-color: #f9f9f9;
}

.file-icon {
  margin-right: 12px;
  font-size: 24px;
  color: #666;
}

.file-info {
  flex: 1;
}

.file-name {
  font-weight: 500;
  margin-bottom: 4px;
}

.file-meta {
  font-size: 12px;
  color: #888;
}

.file-meta span {
  margin-right: 12px;
}

.file-actions {
  display: flex;
  gap: 8px;
}

.btn-download,
.btn-delete {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.btn-download {
  background-color: #007bff;
  color: white;
}

.btn-delete {
  background-color: #dc3545;
  color: white;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background-color: #f0f0f0;
  border-radius: 3px;
  overflow: hidden;
  margin: 10px 0;
}

.progress-fill {
  height: 100%;
  background-color: #007bff;
  transition: width 0.3s ease;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-top: 20px;
}

.pagination button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
  border-radius: 4px;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.no-files {
  text-align: center;
  color: #999;
  padding: 40px;
}
```

## 错误处理

在实现文件管理功能时，需要处理以下常见错误：

1. **文件大小超限** (code: 4000)

   ```javascript
   if (result.code === 4000 && result.messages.includes("50MB")) {
     alert("文件大小不能超过50MB，请选择较小的文件");
   }
   ```

2. **权限不足** (code: 4003)

   ```javascript
   if (result.code === 4003) {
     alert("权限不足：只有文件上传者才能删除文件");
   }
   ```

3. **文件不存在** (code: 4004)

   ```javascript
   if (result.code === 4004) {
     alert("文件不存在或已被删除");
     // 刷新文件列表
     loadChatRoomFiles(getCurrentChatRoomId());
   }
   ```

4. **网络错误**
   ```javascript
   try {
     // API 调用
   } catch (error) {
     if (error.name === "NetworkError") {
       alert("网络连接失败，请检查网络后重试");
     } else {
       alert("操作失败，请重试");
     }
   }
   ```

## 最佳实践

1. **文件类型检查**: 在前端进行基础的文件类型检查，提升用户体验
2. **进度显示**: 对于大文件上传，显示上传进度
3. **缓存策略**: 合理缓存文件列表，避免频繁请求
4. **错误重试**: 对于网络错误，提供重试机制
5. **用户反馈**: 及时给用户反馈操作结果
6. **权限检查**: 在 UI 中根据用户权限显示/隐藏操作按钮

## 安全注意事项

1. **始终验证 JWT Token**: 所有文件操作都需要有效的认证令牌
2. **文件名转义**: 显示文件名时进行 HTML 转义，防止 XSS 攻击
3. **下载链接时效**: 下载链接具有时效性（1 小时），需要及时使用
4. **文件大小限制**: 前端和后端都应检查文件大小限制
5. **Content-Type 检查**: 根据需要检查文件 MIME 类型

这个集成指南提供了完整的前端实现示例，你可以根据具体的前端框架（如 React、Vue 等）进行相应的调整和优化。
