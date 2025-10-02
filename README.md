# ChatApp - Go WebSocket Chat Application

A real-time chat application built with Go, Gin, GORM, and PostgreSQL, featuring multi-storage file management and WebSocket-based real-time communication.

## 🚀 Features

- **Real-time Messaging**: WebSocket-based instant messaging with room-specific broadcasting
- **Multi-User Authentication**: JWT-based secure authentication system
- **Multiple Chat Rooms**: Create and join different chat rooms
- **File Management**: Upload, download, and manage files with multi-storage support (Minio/Qiniu)
- **Message Persistence**: All messages are stored in PostgreSQL database
- **RESTful API**: Comprehensive API for chat rooms, users, and file operations
- **CORS Support**: Configurable cross-origin resource sharing
- **Production Ready**: Structured logging, configuration management, and security best practices

## 🏗️ Architecture Overview

This application follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controllers   │◄───│     Services    │◄───│   Repositories  │
│   (HTTP Layer)  │    │  (Business Logic)│    │   (Data Access) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   WebSocket     │    │     Models      │    │   PostgreSQL    │
│   Handlers      │    │  (Data Schema)  │    │    Database     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │
         ▼
┌─────────────────┐
│   Storage       │
│   (Minio/Qiniu) │
└─────────────────┘
```

## 📁 Project Structure

```
ChatApp/
├── cmd/                          # Command-line tools
│   ├── seed/                     # Database seeding
│   └── test/                     # API testing utilities
├── config/                       # Configuration management
│   ├── config.go                 # Main configuration structure
│   └── database.go               # Database connection and migrations
├── controllers/                  # HTTP request handlers
│   ├── auth_controller.go        # Authentication endpoints
│   ├── chatroom_controller.go    # Chat room management
│   └── file_controller.go        # File operations
├── docs/                         # Documentation
├── examples/                     # Usage examples
├── handlers/                     # WebSocket handlers
│   └── websocket.go              # Real-time communication hub
├── middleware/                   # HTTP middleware
│   └── auth.go                   # JWT authentication middleware
├── models/                       # Data models
│   ├── user.go                   # User entity
│   ├── chatroom.go               # Chat room entity
│   ├── message.go                # Message entity
│   └── file.go                   # File entity
├── repository/                   # Data access layer
│   ├── user_repository.go        # User data operations
│   ├── chatroom_repository.go    # Chat room data operations
│   ├── message_repository.go     # Message data operations
│   └── file_repository.go        # File metadata operations
├── service/                      # Business logic layer
│   ├── auth_service.go           # Authentication logic
│   ├── chatroom_service.go       # Chat room operations
│   ├── message_service.go        # Message handling
│   └── file_service.go           # File management
├── storage/                      # Multi-storage abstraction
│   ├── factory.go                # Storage factory pattern
│   ├── storage.go                # Storage interface
│   ├── minio.go                  # Minio storage implementation
│   └── qiniu.go                  # Qiniu storage implementation
├── utils/                        # Utility functions
│   ├── jwt.go                    # JWT token handling
│   ├── password.go               # Password hashing
│   └── response.go               # Standardized API responses
├── main.go                       # Application entry point
├── go.mod                        # Go module dependencies
└── README.md                     # This file
```

## 🛠️ Technology Stack

- **Backend Framework**: Gin v1.9.1
- **ORM**: GORM v1.25.4
- **Database**: PostgreSQL
- **File Storage**: Minio (default) or Qiniu
- **Real-time Communication**: Gorilla WebSocket v1.5.0
- **Authentication**: JWT v5.0.0
- **Password Hashing**: bcrypt
- **Configuration**: Viper for YAML/ENV configuration
- **Logging**: Structured logging with configurable levels

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Minio server (for file storage) OR Qiniu account
- Git

## ⚙️ Configuration

The application uses a comprehensive YAML-based configuration system. Copy the example configuration:

```bash
cp config.example.yaml config.yaml
```

### Configuration Structure

```yaml
# Application Settings
app:
  name: "ChatApp"
  version: "1.0.0"
  debug: false
  environment: "development"

# Server Configuration
server:
  port: ":8080"
  host: "localhost"
  read_timeout: "30s"
  write_timeout: "30s"

# Database Configuration
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

# JWT Configuration
jwt:
  secret: "your-jwt-secret-key"
  expire_hours: 24
  issuer: "chatapp"

# WebSocket Configuration
websocket:
  read_buffer_size: 1024
  write_buffer_size: 1024
  read_deadline: "60s"
  write_deadline: "10s"
  ping_period: "54s"

# CORS Configuration
cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowed_headers: ["Origin", "Content-Type", "Authorization"]

# Storage Configuration
storage:
  type: "minio"  # "minio" or "qiniu"

# Minio Configuration (if using Minio)
minio:
  endpoint: "127.0.0.1:9000"
  access_key: "minioadmin"
  secret_key: "minioadmin"
  bucket_name: "chatapp"
  use_ssl: false
  region: ""

# Qiniu Configuration (if using Qiniu)
qiniu:
  access_key: "your-qiniu-access-key"
  secret_key: "your-qiniu-secret-key"
  bucket: "your-bucket-name"
  domain: "your-domain.com"
  region: "south-china"
  use_https: true
```

### Environment Variables

All configuration values can be overridden with environment variables:

```bash
export DATABASE_HOST=localhost
export DATABASE_USER=postgres
export DATABASE_PASSWORD=your-password
export JWT_SECRET=your-jwt-secret
export MINIO_ENDPOINT=127.0.0.1:9000
```

## 🚀 Installation & Setup

### 1. Clone and Install Dependencies

```bash
git clone <repository-url>
cd ChatApp
go mod tidy
```

### 2. Database Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE chatapp;
```

### 3. Configuration

```bash
cp config.example.yaml config.yaml
# Edit config.yaml with your settings
```

### 4. Database Seeding

Create test users and chat rooms:

```bash
go run cmd/seed/main.go
```

### 5. Start the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## 👥 Test Users

The seed script creates these test users (all with password `password123`):

- **admin** - Administrator user
- **user1** - Regular user
- **user2** - Regular user  
- **user3** - Regular user

## 🔐 Authentication Flow

### 1. Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

Response:
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

### 2. Use Token in Requests

Include the JWT token in the Authorization header:

```bash
curl -X GET http://localhost:8080/api/chatrooms \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 💬 Chat Room Management

### Get All Chat Rooms

```bash
curl -X GET http://localhost:8080/api/chatrooms \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Create Chat Room

```bash
curl -X POST http://localhost:8080/api/chatrooms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Tech Discussions",
    "description": "Discuss programming and technology"
  }'
```

### Get Chat Room Messages

```bash
curl -X GET http://localhost:8080/api/chatrooms/1/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📁 File Management

### Upload File

```bash
curl -X POST http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/file.jpg" \
  -F "chatroom_id=1"
```

### Get Files in Chat Room

```bash
curl -X GET http://localhost:8080/api/files/chatroom/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Download File

```bash
curl -X GET http://localhost:8080/api/files/download/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -o downloaded-file.jpg
```

## 🔌 WebSocket Real-time Communication

### Connection URL

```
ws://localhost:8080/api/ws/{chatroom_id}
```

### Authentication

Send an authentication message immediately after connecting:

```json
{
  "type": "auth",
  "token": "YOUR_JWT_TOKEN",
  "chatroom_id": 1
}
```

### Send Message

```json
{
  "type": "message",
  "content": "Hello, world!",
  "chatroom_id": 1
}
```

### Message Types

- `auth` - Authentication message
- `message` - Regular chat message
- `auth_success` - Authentication confirmation
- `system` - System notifications

### WebSocket Hub Architecture

The WebSocket implementation uses a hub pattern:

- **Hub**: Manages all connected clients and room-specific broadcasting
- **Client**: Represents a single WebSocket connection with user context
- **Room-based Broadcasting**: Messages are broadcast only to clients in the same chat room
- **Authentication**: Token-based authentication over WebSocket
- **Message Persistence**: All messages are saved to database

## 🗃️ Data Models

### User Model
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

### ChatRoom Model
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

### Message Model
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

### File Model
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
    
    // Relationships
    ChatRoom ChatRoom `json:"chatroom,omitempty" gorm:"foreignKey:ChatRoomID"`
    Uploader User     `json:"uploader,omitempty" gorm:"foreignKey:UploaderID"`
}
```

## 🔧 API Endpoints

### Authentication
- `POST /api/login` - User login

### User Management
- `GET /api/profile` - Get current user profile
- `POST /api/logout` - User logout

### Chat Rooms
- `GET /api/chatrooms` - List all chat rooms
- `POST /api/chatrooms` - Create new chat room
- `GET /api/chatrooms/:id` - Get specific chat room
- `GET /api/chatrooms/:id/messages` - Get chat room messages

### File Management
- `POST /api/files/upload` - Upload file to chat room
- `GET /api/files/download/:id` - Download file
- `GET /api/files/chatroom/:chatroom_id` - Get files in chat room
- `GET /api/files/my` - Get user's uploaded files
- `DELETE /api/files/:id` - Delete file (uploader only)
- `GET /api/files/:id` - Get file information
- `GET /api/files/upload-url` - Get presigned upload URL

### WebSocket
- `GET /api/ws/:chatroom_id` - WebSocket connection for real-time chat

## 🛡️ Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt for secure password storage
- **CORS Protection**: Configurable cross-origin policies
- **Input Validation**: Request parameter validation
- **SQL Injection Protection**: GORM parameterized queries
- **File Type Validation**: MIME type checking for uploads
- **Access Control**: User-based file and resource permissions

## 📦 Storage Implementation

### Multi-Storage Architecture

The application supports multiple storage backends through a factory pattern:

```go
type Storage interface {
    UploadFile(file []byte, fileName string) (string, error)
    DownloadFile(filePath string) ([]byte, error)
    DeleteFile(filePath string) error
    GetUploadURL(fileName string) (string, error)
}
```

### Supported Storage Providers

1. **Minio** (Default): Self-hosted object storage
2. **Qiniu**: Cloud storage service

### Storage Configuration

Switch between storage providers in `config.yaml`:

```yaml
storage:
  type: "minio"  # or "qiniu"
```

## 🧪 Testing

### API Testing

Use the included test utilities:

```bash
go run cmd/test/api_check.go
```

### WebSocket Testing

Use the provided HTML test client:

```bash
open examples/websocket_test.html
```

## 🚀 Deployment

### Build for Production

```bash
go build -o chatapp main.go
```

### Run in Production

```bash
./chatapp
```

### Docker Deployment

Create a `Dockerfile`:

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

## 📚 Code Explanation

### Application Entry Point (`main.go`)

The main application file handles:
- Configuration loading using Viper
- Database connection and migrations
- Dependency injection (repositories → services → controllers)
- WebSocket hub initialization
- HTTP server setup with Gin framework

### Configuration Management (`config/`)

- **config.go**: Defines the complete configuration structure with defaults
- **database.go**: Handles database connection, migrations, and DSN generation
- Uses Viper for YAML/ENV configuration with environment variable overrides

### Data Layer (`models/`, `repository/`)

- **Models**: GORM-based data structures with proper relationships
- **Repositories**: Data access layer implementing repository pattern
- Supports soft deletion and proper indexing

### Business Logic (`service/`)

- **Services**: Contain business logic and orchestrate repository operations
- **Auth Service**: Handles user authentication and JWT token generation
- **Message Service**: Manages message creation and WebSocket integration
- **File Service**: Coordinates file operations with storage providers

### HTTP Layer (`controllers/`, `middleware/`)

- **Controllers**: Handle HTTP requests and responses
- **Middleware**: JWT authentication and CORS handling
- Standardized response format using utility functions

### Real-time Communication (`handlers/websocket.go`)

- **Hub Pattern**: Central WebSocket hub managing all connections
- **Room-based Broadcasting**: Messages are sent only to users in the same chat room
- **Authentication**: Token-based authentication over WebSocket protocol
- **Message Persistence**: All WebSocket messages are saved to database

### Storage Abstraction (`storage/`)

- **Factory Pattern**: Supports multiple storage providers (Minio/Qiniu)
- **Interface-based**: Easy to add new storage providers
- **File Metadata**: File information stored in database with storage references

## 🔄 Workflow Examples

### User Registration & Login
1. User logs in via `/api/login` endpoint
2. Server validates credentials and returns JWT token
3. Client stores token for subsequent requests

### Real-time Chat
1. User connects to WebSocket endpoint with chat room ID
2. Sends authentication message with JWT token
3. After successful auth, can send/receive messages in real-time
4. All messages are persisted to database

### File Upload
1. User uploads file via `/api/files/upload` endpoint
2. File is stored in configured storage provider (Minio/Qiniu)
3. File metadata is saved to database with uploader and chat room info
4. Other users in the chat room can download the file

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check PostgreSQL is running
   - Verify database credentials in config
   - Ensure database exists

2. **WebSocket Connection Failed**
   - Check JWT token is valid
   - Verify chat room ID exists
   - Check CORS configuration

3. **File Upload Fails**
   - Verify Minio/Qiniu configuration
   - Check storage bucket exists and is accessible
   - Verify file size limits

### Logs

Enable debug logging in configuration for detailed troubleshooting:

```yaml
app:
  debug: true

logging:
  level: "debug"
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Gin framework for the HTTP routing
- GORM for database operations
- Gorilla WebSocket for real-time communication
- Minio for object storage
- Qiniu for cloud storage integration

---

**Note**: This is a backend-only application. For a complete chat application, you'll need to build a frontend client that consumes the REST API and WebSocket endpoints.
