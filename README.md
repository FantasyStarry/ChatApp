# ChatApp - Go WebSocket Chat Application

A real-time chat application built with Go, Gin, GORM, and PostgreSQL.

## Features

- JWT-based user authentication
- Multiple chat rooms support
- Real-time messaging via WebSocket
- Message history persistence
- File upload and download with Minio
- File management per chat room
- RESTful API endpoints

## Technology Stack

- **Backend Framework**: Gin v1.9.1
- **ORM**: GORM v1.25.4
- **Database**: PostgreSQL
- **File Storage**: Minio
- **WebSocket**: Gorilla WebSocket v1.5.0
- **Authentication**: JWT v5.0.0
- **Password Hashing**: bcrypt

## Project Structure

```
ChatApp/
├── config/
│   ├── config.go
│   └── database.go
├── controllers/
│   ├── auth_controller.go
│   ├── chatroom_controller.go
│   └── file_controller.go
├── handlers/
│   └── websocket.go
├── middleware/
│   └── auth.go
├── models/
│   ├── chatroom.go
│   ├── message.go
│   ├── user.go
│   └── file.go
├── repository/
│   ├── chatroom_repository.go
│   ├── message_repository.go
│   ├── user_repository.go
│   └── file_repository.go
├── service/
│   ├── auth_service.go
│   ├── chatroom_service.go
│   ├── message_service.go
│   └── file_service.go
├── utils/
│   ├── jwt.go
│   ├── password.go
│   └── response.go
├── main.go
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Minio server (for file storage)

## Configuration

The application uses a YAML-based configuration system. Copy and modify `config.yaml` for your environment:

```yaml
server:
  port: ":8080"
database:
  host: "your-db-host"
  port: 5432
  user: "your-username"
  password: "your-password"
  dbname: "ChatApp"
jwt:
  secret: "your-jwt-secret"
minio:
  endpoint: "127.0.0.1:9000"
  access_key: "minioadmin"
  secret_key: "minioadmin"
  bucket_name: "chatapp"
  use_ssl: false
```

See [Configuration Guide](docs/configuration.md) for detailed configuration options.

## Installation

1. Install Go dependencies:

```bash
go mod tidy
```

2. **Security Setup**: Copy the example configuration and add your credentials:

```bash
cp config.example.yaml config.yaml
# Edit config.yaml with your database credentials and JWT secret
```

**⚠️ Important**: Never commit `config.yaml` with real credentials to version control!

3. Configure the application by editing `config.yaml` or setting environment variables:

```bash
# Using environment variables
export DATABASE_HOST=your-db-host
export DATABASE_USER=your-username
export DATABASE_PASSWORD=your-password
export JWT_SECRET=your-jwt-secret
```

3. Run database seeding (creates test users and chat rooms):

```bash
go run seed.go
```

4. Start the server:

```bash
go run main.go
```

The server will start on port 8080.

## Test Users

The following test users are created by the seed script:

- **Username**: `admin`, **Password**: `password123`
- **Username**: `user1`, **Password**: `password123`
- **Username**: `user2`, **Password**: `password123`
- **Username**: `user3`, **Password**: `password123`

## API Endpoints

### Authentication

- `POST /api/login` - User login

### User

- `GET /api/profile` - Get user profile (requires auth)

### Chat Rooms

- `GET /api/chatrooms` - Get all chat rooms (requires auth)
- `POST /api/chatrooms` - Create new chat room (requires auth)
- `GET /api/chatrooms/:id` - Get specific chat room (requires auth)
- `GET /api/chatrooms/:id/messages` - Get chat room messages (requires auth)

### File Management

- `POST /api/files/upload` - Upload file to chat room (requires auth)
- `GET /api/files/download/:id` - Get file download link (requires auth)
- `GET /api/files/chatroom/:chatroom_id` - Get files in chat room (requires auth)
- `GET /api/files/my` - Get user's uploaded files (requires auth)
- `DELETE /api/files/:id` - Delete file (requires auth, uploader only)
- `GET /api/files/:id` - Get file information (requires auth)
- `GET /api/files/upload-url` - Get presigned upload URL (requires auth)

### WebSocket

- `GET /api/ws/:chatroom_id` - WebSocket connection for chat room (requires auth)

## Usage Examples

### 1. Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

### 2. Get Chat Rooms

```bash
curl -X GET http://localhost:8080/api/chatrooms \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. Create Chat Room

```bash
curl -X POST http://localhost:8080/api/chatrooms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name": "New Room", "description": "A new chat room"}'
```

### 4. WebSocket Connection

Connect to WebSocket at:

```
ws://localhost:8080/api/ws/1?token=YOUR_JWT_TOKEN
```

Send messages in JSON format:

```json
{
  "type": "message",
  "content": "Hello, world!",
  "chatroom_id": 1
}
```

## Default Chat Rooms

The seed script creates these default chat rooms:

1. **General** - General discussion room
2. **Tech Talk** - Technical discussions and programming
3. **Random** - Random conversations and fun stuff

## Configuration

You can modify the configuration in `config.yaml` or use environment variables:

### Key Configuration Options:

- **Server**: Port, timeouts, host settings
- **Database**: Connection settings, pool configuration
- **JWT**: Secret key, expiration, issuer
- **WebSocket**: Buffer sizes, timeouts
- **CORS**: Allowed origins, methods, headers
- **Logging**: Level, format, output

For detailed configuration options, see [Configuration Guide](docs/configuration.md).

**🔒 Security**: See [Security Guide](docs/security.md) for production deployment and security best practices.

## WebSocket Message Format

Messages sent and received via WebSocket follow this format:

```json
{
  "type": "message",
  "content": "Message content",
  "user_id": 1,
  "username": "admin",
  "chatroom_id": 1,
  "timestamp": "2023-01-01T12:00:00Z"
}
```

## Security Features

- JWT-based authentication
- Password hashing with bcrypt
- CORS support
- Request validation

## Build

To build the application:

```bash
go build -o chatapp main.go
```

Then run:

```bash
./chatapp
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
