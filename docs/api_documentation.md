# ChatApp API Documentation

## 概述

ChatApp 是一个基于 Go 语言开发的实时聊天应用，提供了 RESTful API 和 WebSocket 接口。该文档详细描述了所有可用的 API 端点和 WebSocket 通信协议。

## 技术栈

- **后端框架**: Gin v1.9.1
- **ORM**: GORM v1.25.4
- **数据库**: PostgreSQL
- **WebSocket**: Gorilla WebSocket v1.5.0
- **认证**: JWT v5.0.0
- **密码加密**: bcrypt

## 认证机制

大多数 API 端点都需要通过 JWT Token 进行认证。获取 Token 的方式是通过登录接口。

### 获取 Token

使用登录接口获取 JWT Token：

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

返回示例：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 使用 Token

在需要认证的接口中，在请求头中添加 Authorization 字段：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## API 端点

### 认证相关

#### 用户登录
- **URL**: `POST /api/login`
- **描述**: 用户登录并获取 JWT Token
- **请求参数**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **响应**:
  ```json
  {
    "token": "string",
    "user": {
      "id": "integer",
      "username": "string",
      "email": "string",
      "created_at": "datetime"
    }
  }
  ```

#### 获取用户信息
- **URL**: `GET /api/profile`
- **描述**: 获取当前登录用户的信息
- **认证**: 需要 Bearer Token
- **响应**:
  ```json
  {
    "id": "integer",
    "username": "string",
    "email": "string",
    "created_at": "datetime"
  }
  ```

#### 用户登出
- **URL**: `POST /api/logout`
- **描述**: 用户登出
- **认证**: 需要 Bearer Token
- **响应**:
  ```json
  {
    "message": "Logout successful",
    "user_id": "integer"
  }
  ```

### 聊天室相关

#### 获取所有聊天室
- **URL**: `GET /api/chatrooms`
- **描述**: 获取所有聊天室列表
- **认证**: 需要 Bearer Token
- **响应**:
  ```json
  [
    {
      "id": "integer",
      "name": "string",
      "description": "string",
      "created_by": "integer",
      "creator": {
        "id": "integer",
        "username": "string",
        "email": "string"
      },
      "created_at": "datetime"
    }
  ]
  ```

#### 创建聊天室
- **URL**: `POST /api/chatrooms`
- **描述**: 创建新的聊天室
- **认证**: 需要 Bearer Token
- **请求参数**:
  ```json
  {
    "name": "string",
    "description": "string"
  }
  ```
- **响应**:
  ```json
  {
    "id": "integer",
    "name": "string",
    "description": "string",
    "created_by": "integer",
    "creator": {
      "id": "integer",
      "username": "string",
      "email": "string"
    },
    "created_at": "datetime"
  }
  ```

#### 获取特定聊天室
- **URL**: `GET /api/chatrooms/{id}`
- **描述**: 获取特定聊天室的详细信息
- **认证**: 需要 Bearer Token
- **路径参数**:
  - `id`: 聊天室ID
- **响应**:
  ```json
  {
    "id": "integer",
    "name": "string",
    "description": "string",
    "created_by": "integer",
    "creator": {
      "id": "integer",
      "username": "string",
      "email": "string"
    },
    "created_at": "datetime",
    "messages": [
      {
        "id": "integer",
        "content": "string",
        "user_id": "integer",
        "user": {
          "id": "integer",
          "username": "string"
        },
        "chatroom_id": "integer",
        "created_at": "datetime"
      }
    ]
  }
  ```

#### 获取聊天室消息
- **URL**: `GET /api/chatrooms/{id}/messages`
- **描述**: 获取特定聊天室的消息列表
- **认证**: 需要 Bearer Token
- **路径参数**:
  - `id`: 聊天室ID
- **查询参数**:
  - `limit`: 每页消息数量（默认50）
  - `offset`: 偏移量（默认0）
- **响应**:
  ```json
  [
    {
      "id": "integer",
      "content": "string",
      "user_id": "integer",
      "user": {
        "id": "integer",
        "username": "string"
      },
      "chatroom_id": "integer",
      "created_at": "datetime"
    }
  ]
  ```

## WebSocket 接口

### 连接地址
```
ws://localhost:8080/api/ws/{chatroom_id}
```

### 认证流程

WebSocket 的认证机制已更新，采用消息认证方式：

1. 建立 WebSocket 连接（无需认证）
2. 发送认证消息
3. 等待认证成功响应
4. 开始发送和接收消息

### 认证消息格式

连接建立后，客户端需要立即发送认证消息：

```json
{
  "type": "auth",
  "token": "your-jwt-token-here",
  "chatroomId": 1
}
```

### 认证响应

认证成功：
```json
{
  "type": "auth_success",
  "content": "Authentication successful",
  "timestamp": "2023-12-18T10:30:00Z"
}
```

### 消息格式

#### 发送消息
```json
{
  "type": "message",
  "content": "Hello, world!"
}
```

#### 接收消息
```json
{
  "type": "message",
  "content": "Hello, world!",
  "user_id": 1,
  "username": "admin",
  "chatroom_id": 1,
  "timestamp": "2023-12-18T10:30:00Z"
}
```

### 客户端实现示例

```javascript
// 1. 建立连接
const ws = new WebSocket('ws://localhost:8080/api/ws/1');

// 2. 连接打开后发送认证消息
ws.onopen = function(event) {
  const authMessage = {
    type: 'auth',
    token: 'your-jwt-token-here',
    chatroomId: 1
  };
  ws.send(JSON.stringify(authMessage));
};

// 3. 处理消息
ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  
  if (data.type === 'auth_success') {
    console.log('认证成功');
  } else if (data.type === 'message') {
    console.log(`[${data.username}]: ${data.content}`);
  }
};

// 4. 发送消息
function sendMessage(content) {
  const message = {
    type: 'message',
    content: content
  };
  ws.send(JSON.stringify(message));
}
```

## 错误响应格式

所有错误响应都遵循以下格式：

```json
{
  "error": "错误描述信息"
}
```

常见的 HTTP 状态码：
- `400`: 请求参数错误
- `401`: 未认证或认证失败
- `404`: 资源未找到
- `500`: 服务器内部错误

## 测试用户

应用提供了以下测试用户用于开发和测试：

- **用户名**: `admin`, **密码**: `password123`
- **用户名**: `user1`, **密码**: `password123`
- **用户名**: `user2`, **密码**: `password123`
- **用户名**: `user3`, **密码**: `password123`

## 默认聊天室

应用默认创建了以下聊天室：

1. **General** - 通用讨论室
2. **Tech Talk** - 技术讨论和编程
3. **Random** - 随意聊天和有趣内容

## 安全注意事项

1. 不要在版本控制系统中提交包含真实凭证的配置文件
2. 生产环境中使用强 JWT 密钥（至少64个字符）
3. 启用数据库 SSL 连接
4. 设置特定的 CORS 允许来源，不要使用通配符
5. 定期轮换密钥和密码