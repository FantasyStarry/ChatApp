# ChatApp 代码结构分析

## 项目概述

ChatApp 是一个基于 Go 语言开发的实时聊天应用，采用了分层架构设计，包括 Repository、Service 和 Controller 三层。项目使用 Gin 作为 Web 框架，GORM 作为 ORM，PostgreSQL 作为数据库，WebSocket 实现实时通信。

## 项目架构

项目采用了 Repository-Service-Controller 分层架构：

```
┌─────────────────────┐
│   Controller Layer  │  ← HTTP请求处理，参数验证，响应格式化
├─────────────────────┤
│   Service Layer     │  ← 业务逻辑，业务规则验证，事务管理
├─────────────────────┤
│   Repository Layer  │  ← 数据访问，数据库操作，数据映射
├─────────────────────┤
│   Model Layer       │  ← 数据模型定义
└─────────────────────┘
```

## 目录结构

```
ChatApp/
├── config/                 # 配置管理
│   ├── config.go          # 配置加载和管理
│   └── database.go        # 数据库连接
├── controllers/            # 控制器层
│   ├── auth_controller.go # 认证相关控制器
│   └── chatroom_controller.go # 聊天室相关控制器
├── handlers/               # WebSocket 处理器
│   └── websocket.go       # WebSocket 连接和消息处理
├── middleware/             # 中间件
│   └── auth.go            # 认证中间件
├── models/                 # 数据模型
│   ├── chatroom.go        # 聊天室模型
│   ├── message.go         # 消息模型
│   └── user.go            # 用户模型
├── repository/             # 数据访问层
│   ├── chatroom_repository.go # 聊天室数据访问
│   ├── message_repository.go  # 消息数据访问
│   └── user_repository.go     # 用户数据访问
├── service/                # 业务逻辑层
│   ├── auth_service.go        # 认证业务逻辑
│   ├── chatroom_service.go    # 聊天室业务逻辑
│   └── message_service.go     # 消息业务逻辑
├── utils/                  # 工具函数
│   ├── jwt.go             # JWT 认证工具
│   └── password.go        # 密码加密工具
├── docs/                   # 文档
│   ├── api_documentation.md      # API 文档
│   ├── architecture.md           # 架构文档
│   ├── configuration.md          # 配置文档
│   ├── security.md               # 安全文档
│   ├── websocket-refactoring.md  # WebSocket 重构文档
│   └── websocket_auth_update.md  # WebSocket 认证更新文档
├── examples/               # 示例文件
│   ├── api_examples.http  # API 调用示例
│   └── websocket_test.html # WebSocket 测试页面
├── cmd/                    # 命令行工具
│   └── seed/              # 数据库种子数据
│       └── main.go
├── main.go                 # 程序入口
├── seed.go                 # 数据库种子脚本
├── go.mod                  # Go 模块定义
├── go.sum                  # Go 模块校验和
├── README.md               # 项目说明
├── config.example.yaml     # 配置文件示例
└── .gitignore              # Git 忽略文件
```

## 核心功能模块

### 1. 用户认证模块

**功能**:
- 用户登录和登出
- JWT Token 生成和验证
- 用户信息获取

**关键文件**:
- `utils/jwt.go`: JWT Token 的生成和验证
- `utils/password.go`: 密码加密和验证
- `middleware/auth.go`: 认证中间件
- `controllers/auth_controller.go`: 认证控制器
- `service/auth_service.go`: 认证业务逻辑

**流程**:
1. 用户通过 `/api/login` 接口登录
2. 服务验证用户名和密码
3. 生成 JWT Token 并返回给客户端
4. 客户端在后续请求中通过 `Authorization: Bearer <token>` 头传递 Token
5. 服务通过中间件验证 Token 有效性

### 2. 聊天室管理模块

**功能**:
- 聊天室创建、查询
- 聊天室消息查询

**关键文件**:
- `models/chatroom.go`: 聊天室数据模型
- `repository/chatroom_repository.go`: 聊天室数据访问
- `service/chatroom_service.go`: 聊天室业务逻辑
- `controllers/chatroom_controller.go`: 聊天室控制器

**主要接口**:
- `GET /api/chatrooms`: 获取所有聊天室
- `POST /api/chatrooms`: 创建聊天室
- `GET /api/chatrooms/{id}`: 获取特定聊天室信息

### 3. 消息管理模块

**功能**:
- 消息发送和存储
- 消息查询（分页）

**关键文件**:
- `models/message.go`: 消息数据模型
- `repository/message_repository.go`: 消息数据访问
- `service/message_service.go`: 消息业务逻辑

**主要接口**:
- `GET /api/chatrooms/{id}/messages`: 获取聊天室消息（支持分页）

### 4. WebSocket 实时通信模块

**功能**:
- 实时消息推送
- 用户连接管理

**关键文件**:
- `handlers/websocket.go`: WebSocket 处理器
- `models/message.go`: WebSocket 消息模型

**特点**:
- 采用消息认证方式，连接建立后发送认证消息
- 支持多聊天室，客户端连接时指定聊天室ID
- 消息广播到同一聊天室的所有客户端

**认证流程**:
1. 客户端建立 WebSocket 连接
2. 发送包含 JWT Token 的认证消息
3. 服务端验证 Token 并注册客户端
4. 认证成功后开始正常消息通信

## 数据模型

### User (用户)
- `ID`: 用户ID
- `Username`: 用户名
- `Password`: 加密后的密码
- `Email`: 邮箱
- `CreatedAt`: 创建时间
- `UpdatedAt`: 更新时间

### ChatRoom (聊天室)
- `ID`: 聊天室ID
- `Name`: 聊天室名称
- `Description`: 聊天室描述
- `CreatedBy`: 创建者ID
- `Creator`: 创建者信息（关联User）
- `CreatedAt`: 创建时间
- `UpdatedAt`: 更新时间
- `Messages`: 聊天室消息列表（关联Message）

### Message (消息)
- `ID`: 消息ID
- `Content`: 消息内容
- `UserID`: 发送者ID
- `User`: 发送者信息（关联User）
- `ChatRoomID`: 聊天室ID
- `ChatRoom`: 聊天室信息（关联ChatRoom）
- `CreatedAt`: 创建时间
- `UpdatedAt`: 更新时间

### WSMessage (WebSocket消息)
- `Type`: 消息类型（auth, auth_success, message）
- `Content`: 消息内容
- `UserID`: 用户ID
- `Username`: 用户名
- `ChatRoomID`: 聊天室ID
- `Timestamp`: 时间戳
- `Token`: JWT Token（仅用于认证消息）

## 技术特点

### 1. 分层架构
- 清晰的职责分离
- 易于测试和维护
- 支持依赖注入

### 2. 安全性
- JWT Token 认证
- 密码加密存储
- CORS 支持
- 配置文件安全管理

### 3. 实时通信
- WebSocket 实现
- 消息广播机制
- 连接状态管理

### 4. 可扩展性
- 模块化设计
- 接口抽象
- 易于添加新功能

## 部署和配置

### 环境要求
- Go 1.21 或更高版本
- PostgreSQL 数据库

### 配置方式
1. 复制 `config.example.yaml` 为 `config.yaml`
2. 修改数据库连接信息和 JWT 密钥
3. 或使用环境变量配置

### 安全注意事项
- 不要在版本控制中提交真实的配置文件
- 生产环境使用强 JWT 密钥
- 启用数据库 SSL 连接
- 设置特定的 CORS 允许来源

## 测试和开发

### 测试用户
- admin/password123
- user1/password123
- user2/password123
- user3/password123

### API 测试
- 使用 `examples/api_examples.http` 文件进行 API 测试

### WebSocket 测试
- 使用 `examples/websocket_test.html` 文件进行 WebSocket 测试

## 总结

ChatApp 采用了现代化的 Go Web 应用架构，具有清晰的分层结构、完善的安全机制和实时通信功能。代码组织良好，易于维护和扩展，适合作为聊天应用的基础框架。