# ChatApp - Go WebSocket 聊天应用

一个使用 Go、Gin、GORM 和 PostgreSQL 构建的实时聊天应用，具有多存储文件管理和基于 WebSocket 的实时通信功能。

## 🚀 功能特性

- **实时消息传递**: 基于 WebSocket 的即时消息传递，支持房间特定广播
- **多用户认证**: JWT 安全认证系统
- **多聊天室**: 创建和加入不同的聊天室
- **文件管理**: 上传、下载和管理文件，支持多存储（Minio/Qiniu）
- **消息持久化**: 所有消息都存储在 PostgreSQL 数据库中
- **RESTful API**: 完整的聊天室、用户和文件操作 API
- **CORS 支持**: 可配置的跨域资源共享
- **生产就绪**: 结构化日志、配置管理和安全最佳实践

## 🏗️ 架构概览

本应用采用清晰的架构模式，关注点分离明确：

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   控制器层      │◄───│    服务层       │◄───│   数据访问层    │
│  (HTTP 层)      │    │  (业务逻辑)     │    │  (数据操作)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   WebSocket     │    │     数据模型    │    │   PostgreSQL    │
│    处理器       │    │   (数据结构)    │    │     数据库      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │
         ▼
┌─────────────────┐
│    存储层       │
│  (Minio/Qiniu)  │
└─────────────────┘
```

## 📁 项目结构

```
ChatApp/
├── cmd/                          # 命令行工具
│   ├── seed/                     # 数据库种子数据
│   └── test/                     # API 测试工具
├── config/                       # 配置管理
│   ├── config.go                 # 主配置结构
│   └── database.go               # 数据库连接和迁移
├── controllers/                  # HTTP 请求处理器
│   ├── auth_controller.go        # 认证端点
│   ├── chatroom_controller.go    # 聊天室管理
│   └── file_controller.go        # 文件操作
├── docs/                         # 文档
├── examples/                     # 使用示例
├── handlers/                     # WebSocket 处理器
│   └── websocket.go              # 实时通信中心
├── middleware/                   # HTTP 中间件
│   └── auth.go                   # JWT 认证中间件
├── models/                       # 数据模型
│   ├── user.go                   # 用户实体
│   ├── chatroom.go               # 聊天室实体
│   ├── message.go                # 消息实体
│   └── file.go                   # 文件实体
├── repository/                   # 数据访问层
│   ├── user_repository.go        # 用户数据操作
│   ├── chatroom_repository.go    # 聊天室数据操作
│   ├── message_repository.go     # 消息数据操作
│   └── file_repository.go        # 文件元数据操作
├── service/                      # 业务逻辑层
│   ├── auth_service.go           # 认证逻辑
│   ├── chatroom_service.go       # 聊天室操作
│   ├── message_service.go        # 消息处理
│   └── file_service.go           # 文件管理
├── storage/                      # 多存储抽象
│   ├── factory.go                # 存储工厂模式
│   ├── storage.go                # 存储接口
│   ├── minio.go                  # Minio 存储实现
│   └── qiniu.go                  # 七牛云存储实现
├── utils/                        # 工具函数
│   ├── jwt.go                    # JWT 令牌处理
│   ├── password.go               # 密码哈希
│   └── response.go               # 标准化 API 响应
├── main.go                       # 应用入口点
├── go.mod                        # Go 模块依赖
└── README.md                     # 本文档
```

## 🛠️ 技术栈

- **后端框架**: Gin v1.9.1
- **ORM**: GORM v1.25.4
- **数据库**: PostgreSQL
- **文件存储**: Minio（默认）或 七牛云
- **实时通信**: Gorilla WebSocket v1.5.0
- **认证**: JWT v5.0.0
- **密码哈希**: bcrypt
- **配置管理**: Viper 支持 YAML/ENV 配置
- **日志**: 结构化日志，可配置级别

## 📋 前置要求

- Go 1.21 或更高版本
- PostgreSQL 数据库
- Minio 服务器（用于文件存储）或 七牛云账户
- Git

## ⚙️ 配置

应用使用基于 YAML 的全面配置系统。复制示例配置：

```bash
cp config.example.yaml config.yaml
```

### 配置结构

```yaml
# 应用设置
app:
  name: "ChatApp"
  version: "1.0.0"
  debug: false
  environment: "development"

# 服务器配置
server:
  port: ":8080"
  host: "localhost"
  read_timeout: "30s"
  write_timeout: "30s"

# 数据库配置
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "your-password"
  dbname: "chatapp"
  sslmode: "disable"
  timezone: "UTC"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: "3600s"

# JWT 配置
jwt:
  secret: "your-jwt-secret-key"
  expire_hours: 24
  issuer: "chatapp"

# WebSocket 配置
websocket:
  read_buffer_size: 1024
  write_buffer_size: 1024
  read_deadline: "60s"
  write_deadline: "10s"
  ping_period: "54s"

# CORS 配置
cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowed_headers: ["Origin", "Content-Type", "Authorization"]

# 存储配置
storage:
  type: "minio"  # "minio" 或 "qiniu"

# Minio 配置（如果使用 Minio）
minio:
  endpoint: "127.0.0.1:9000"
  access_key: "minioadmin"
  secret_key: "minioadmin"
  bucket_name: "chatapp"
  use_ssl: false
  region: ""

# 七牛云配置（如果使用七牛云）
qiniu:
  access_key: "your-qiniu-access-key"
  secret_key: "your-qiniu-secret-key"
  bucket: "your-bucket-name"
  domain: "your-domain.com"
  region: "south-china"
  use_https: true
```

### 环境变量

所有配置值都可以通过环境变量覆盖：

```bash
export DATABASE_HOST=localhost
export DATABASE_USER=postgres
export DATABASE_PASSWORD=your-password
export JWT_SECRET=your-jwt-secret
export MINIO_ENDPOINT=127.0.0.1:9000
```

## 🚀 安装与设置

### 1. 克隆并安装依赖

```bash
git clone <repository-url>
cd ChatApp
go mod tidy
```

### 2. 数据库设置

创建 PostgreSQL 数据库：

```sql
CREATE DATABASE chatapp;
```

### 3. 配置

```bash
cp config.example.yaml config.yaml
# 编辑 config.yaml 文件，填入你的设置
```

### 4. 数据库种子数据

创建测试用户和聊天室：

```bash
go run cmd/seed/main.go
```

### 5. 启动应用

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动

## 👥 测试用户

种子脚本创建以下测试用户（所有用户密码均为 `password123`）：

- **admin** - 管理员用户
- **user1** - 普通用户
- **user2** - 普通用户  
- **user3** - 普通用户

## 🔐 认证流程

### 1. 登录

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

响应：
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com"
    }
  }
}
```

### 2. 在请求中使用令牌

在 Authorization 头中包含 JWT 令牌：

```bash
curl -X GET http://localhost:8080/api/chatrooms \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 💬 聊天室管理

### 获取所有聊天室

```bash
curl -X GET http://localhost:8080/api/chatrooms \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 创建聊天室

```bash
curl -X POST http://localhost:8080/api/chatrooms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "技术讨论",
    "description": "讨论编程和技术"
  }'
```

### 获取聊天室消息

```bash
curl -X GET http://localhost:8080/api/chatrooms/1/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📁 文件管理

### 上传文件

```bash
curl -X POST http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/file.jpg" \
  -F "chatroom_id=1"
```

### 获取聊天室中的文件

```bash
curl -X GET http://localhost:8080/api/files/chatroom/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 下载文件

```bash
curl -X GET http://localhost:8080/api/files/download/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -o downloaded-file.jpg
```

## 🔌 WebSocket 实时通信

### 连接 URL

```
ws://localhost:8080/api/ws/{chatroom_id}
```

### 认证

连接后立即发送认证消息：

```json
{
  "type": "auth",
  "token": "YOUR_JWT_TOKEN",
  "chatroom_id": 1
}
```

### 发送消息

```json
{
  "type": "message",
  "content": "你好，世界！",
  "chatroom_id": 1
}
```

### 消息类型

- `auth` - 认证消息
- `message` - 常规聊天消息
- `auth_success` - 认证确认
- `system` - 系统通知

### WebSocket 中心架构

WebSocket 实现使用中心模式：

- **中心（Hub）**: 管理所有连接的客户端和房间特定广播
- **客户端（Client）**: 代表具有用户上下文的单个 WebSocket 连接
- **基于房间的广播**: 消息仅广播到同一聊天室中的客户端
- **认证**: 通过 WebSocket 进行基于令牌的认证
- **消息持久化**: 所有消息都保存到数据库

## 🗃️ 数据模型

### 用户模型
```go
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"uniqueIndex;not null"`
    Password  string         `json:"-" gorm:"not null"`
    Email     string         `json:"email" gorm:"uniqueIndex"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### 聊天室模型
```go
type ChatRoom struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"not null"`
    Description string         `json:"description"`
    CreatedBy   uint           `json:"created_by"`
    Creator     User           `json:"creator" gorm:"foreignKey:CreatedBy"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    Messages    []Message      `json:"messages,omitempty" gorm:"foreignKey:ChatRoomID"`
}
```

### 消息模型
```go
type Message struct {
    ID         uint           `json:"id" gorm:"primaryKey"`
    Content    string         `json:"content" gorm:"not null;type:text"`
    UserID     uint           `json:"user_id"`
    User       User           `json:"user" gorm:"foreignKey:UserID"`
    ChatRoomID uint           `json:"chat_room_id"`
    Type       string         `json:"type" gorm:"type:varchar(20);default:'message';not null"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### 文件模型
```go
type File struct {
    ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    FileName    string    `json:"file_name" gorm:"not null;size:255"`
    FilePath    string    `json:"file_path" gorm:"not null;size:500"`
    FileSize    int64     `json:"file_size" gorm:"not null"`
    ContentType string    `json:"content_type" gorm:"size:100"`
    ChatRoomID  uint      `json:"chat_room_id" gorm:"not null;index"`
    UploaderID  uint      `json:"uploader_id" gorm:"not null;index"`
    UploadedAt  time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 关联关系
    ChatRoom ChatRoom `json:"chatroom,omitempty" gorm:"foreignKey:ChatRoomID"`
    Uploader User     `json:"uploader,omitempty" gorm:"foreignKey:UploaderID"`
}
```

## 🔧 API 端点

### 认证
- `POST /api/login` - 用户登录

### 用户管理
- `GET /api/profile` - 获取当前用户资料
- `POST /api/logout` - 用户登出

### 聊天室
- `GET /api/chatrooms` - 列出所有聊天室
- `POST /api/chatrooms` - 创建新聊天室
- `GET /api/chatrooms/:id` - 获取特定聊天室
- `GET /api/chatrooms/:id/messages` - 获取聊天室消息

### 文件管理
- `POST /api/files/upload` - 上传文件到聊天室
- `GET /api/files/download/:id` - 下载文件
- `GET /api/files/chatroom/:chatroom_id` - 获取聊天室中的文件
- `GET /api/files/my` - 获取用户上传的文件
- `DELETE /api/files/:id` - 删除文件（仅上传者）
- `GET /api/files/:id` - 获取文件信息
- `GET /api/files/upload-url` - 获取预签名上传 URL

### WebSocket
- `GET /api/ws/:chatroom_id` - 实时聊天的 WebSocket 连接

## 🛡️ 安全特性

- **JWT 认证**: 安全的基于令牌的认证
- **密码哈希**: 使用 bcrypt 安全存储密码
- **CORS 保护**: 可配置的跨域策略
- **输入验证**: 请求参数验证
- **SQL 注入防护**: GORM 参数化查询
- **文件类型验证**: 上传文件的 MIME 类型检查
- **访问控制**: 基于用户的文件和资源权限

## 📦 存储实现

### 多存储架构

应用通过工厂模式支持多个存储后端：

```go
type Storage interface {
    UploadFile(file []byte, fileName string) (string, error)
    DownloadFile(filePath string) ([]byte, error)
    DeleteFile(filePath string) error
    GetUploadURL(fileName string) (string, error)
}
```

### 支持的存储提供商

1. **Minio**（默认）: 自托管对象存储
2. **七牛云**: 云存储服务

### 存储配置

在 `config.yaml` 中切换存储提供商：

```yaml
storage:
  type: "minio"  # 或 "qiniu"
```

## 🧪 测试

### API 测试

使用包含的测试工具：

```bash
go run cmd/test/api_check.go
```

### WebSocket 测试

使用提供的 HTML 测试客户端：

```bash
open examples/websocket_test.html
```

## 🚀 部署

### 生产环境构建

```bash
go build -o chatapp main.go
```

### 生产环境运行

```bash
./chatapp
```

### Docker 部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o chatapp main.go

EXPOSE 8080

CMD ["./chatapp"]
```

## 📚 代码解释

### 应用入口点 (`main.go`)

主应用文件处理：
- 使用 Viper 加载配置
- 数据库连接和迁移
- 依赖注入（存储库 → 服务 → 控制器）
- WebSocket 中心初始化
- 使用 Gin 框架设置 HTTP 服务器

### 配置管理 (`config/`)

- **config.go**: 定义完整的配置结构，包含默认值
- **database.go**: 处理数据库连接、迁移和 DSN 生成
- 使用 Viper 进行 YAML/ENV 配置，支持环境变量覆盖

### 数据层 (`models/`, `repository/`)

- **模型**: 基于 GORM 的数据结构，具有正确的关系
- **存储库**: 实现存储库模式的数据访问层
- 支持软删除和适当的索引

### 业务逻辑 (`service/`)

- **服务**: 包含业务逻辑并协调存储库操作
- **认证服务**: 处理用户认证和 JWT 令牌生成
- **消息服务**: 管理消息创建和 WebSocket 集成
- **文件服务**: 与存储提供商协调文件操作

### HTTP 层 (`controllers/`, `middleware/`)

- **控制器**: 处理 HTTP 请求和响应
- **中间件**: JWT 认证和 CORS 处理
- 使用工具函数实现标准化响应格式

### 实时通信 (`handlers/websocket.go`)

- **中心模式**: 中央 WebSocket 中心管理所有连接
- **基于房间的广播**: 消息仅发送到同一聊天室的用户
- **认证**: 通过 WebSocket 协议进行基于令牌的认证
- **消息持久化**: 所有 WebSocket 消息都保存到数据库

### 存储抽象 (`storage/`)

- **工厂模式**: 支持多个存储提供商（Minio/七牛云）
- **基于接口**: 易于添加新的存储提供商
- **文件元数据**: 文件信息存储在数据库中，包含存储引用

## 🔄 工作流程示例

### 用户注册和登录
1. 用户通过 `/api/login` 端点登录
2. 服务器验证凭据并返回 JWT 令牌
3. 客户端存储令牌用于后续请求

### 实时聊天
1. 用户连接到带有聊天室 ID 的 WebSocket 端点
2. 发送带有 JWT 令牌的认证消息
3. 认证成功后，可以实时发送/接收消息
4. 所有消息都持久化到数据库

### 文件上传
1. 用户通过 `/api/files/upload` 端点上传文件
2. 文件存储在配置的存储提供商（Minio/七牛云）中
3. 文件元数据保存到数据库，包含上传者和聊天室信息
4. 聊天室中的其他用户可以下载该文件

## 🐛 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 PostgreSQL 是否运行
   - 验证配置中的数据库凭据
   - 确保数据库存在

2. **WebSocket 连接失败**
   - 检查 JWT 令牌是否有效
   - 验证聊天室 ID 是否存在
   - 检查 CORS 配置

3. **文件上传失败**
   - 验证 Minio/七牛云配置
   - 检查存储桶是否存在且可访问
   - 验证文件大小限制

### 日志

在配置中启用调试日志以进行详细故障排除：

```yaml
app:
  debug: true

logging:
  level: "debug"
```

## 🤝 贡献

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m '添加一些很棒的功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 有关详细信息，请参阅 LICENSE 文件。

## 🙏 致谢

- Gin 框架用于 HTTP 路由
- GORM 用于数据库操作
- Gorilla WebSocket 用于实时通信
- Minio 用于对象存储
- 七牛云用于云存储集成

---

**注意**: 这是一个仅包含后端的应用程序。要构建完整的聊天应用程序，您需要构建一个使用 REST API 和 WebSocket 端点的前端客户端。
