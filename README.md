# ChatApp - 企业级聊天应用

一个现代化的企业级聊天应用，采用前后端分离架构，前端使用 Next.js + TailwindCSS + shadcn/ui，后端使用 Go。

## 项目结构

```
ChatApp/
├── frontend/                 # Next.js 前端应用
│   ├── src/
│   │   ├── app/             # App Router
│   │   ├── components/      # UI 组件
│   │   └── lib/             # 工具函数
│   ├── package.json
│   └── tailwind.config.ts
├── backend/                  # Go 后端服务
│   ├── cmd/                 # 命令行工具
│   ├── config/              # 配置文件
│   ├── controllers/         # 控制器层
│   ├── handlers/            # 请求处理器
│   ├── middleware/          # 中间件
│   ├── models/              # 数据模型
│   ├── repository/          # 数据访问层
│   ├── service/             # 业务逻辑层
│   ├── storage/             # 文件存储
│   ├── utils/               # 工具函数
│   ├── main.go
│   └── go.mod
├── docs/                    # 项目文档
└── examples/                # 示例文件
```

## 技术栈

### 前端
- **框架**: Next.js 15.5.4 (App Router)
- **语言**: TypeScript
- **样式**: TailwindCSS v4
- **UI 组件**: shadcn/ui
- **包管理**: pnpm
- **主题**: 蓝色主题，类似钉钉风格

### 后端
- **语言**: Go
- **Web 框架**: 标准库 + 自定义路由
- **数据库**: 支持多种数据库 (通过 repository 层)
- **认证**: JWT
- **文件存储**: 支持 MinIO、七牛云等
- **WebSocket**: 实时消息推送

## 快速开始

### 前端开发

1. 进入前端目录：
```bash
cd frontend
```

2. 安装依赖：
```bash
pnpm install
```

3. 启动开发服务器：
```bash
pnpm run dev
```

前端应用将在 http://localhost:3000 运行。

### 后端开发

1. 进入后端目录：
```bash
cd backend
```

2. 安装 Go 依赖：
```bash
go mod tidy
```

3. 配置环境变量：
```bash
cp config.example.yaml config.yaml
# 编辑 config.yaml 文件配置数据库等
```

4. 启动后端服务：
```bash
go run main.go
```

后端服务默认运行在 http://localhost:8080

## 功能特性

### 前端功能
- ✅ 现代化的聊天界面（类似钉钉风格）
- ✅ 响应式设计，适配不同屏幕
- ✅ 蓝色主题配色，统一视觉风格
- ✅ 实时消息发送和显示
- ✅ 联系人列表，支持在线状态显示
- ✅ 聊天室管理，支持未读消息计数
- ✅ 消息时间戳显示
- ✅ 搜索功能
- ✅ 加载状态处理
- ✅ 防Hydration错误处理

### 后端功能
- ✅ RESTful API
- ✅ WebSocket 实时通信
- ✅ JWT 认证
- ✅ 用户管理
- ✅ 聊天室管理
- ✅ 消息存储
- ✅ 文件上传
- ✅ 多存储支持

## 开发指南

### 添加新的 UI 组件

使用 shadcn/ui 添加新组件：
```bash
cd frontend
pnpm dlx shadcn@latest add [component-name]
```

### 自定义主题

编辑 `frontend/src/app/globals.css` 文件中的 CSS 变量来修改主题颜色。

### API 集成

前端通过 REST API 和 WebSocket 与后端通信：
- REST API: `http://localhost:8080/api/v1`
- WebSocket: `ws://localhost:8080/ws`

## 部署

### 前端部署

构建生产版本：
```bash
cd frontend
pnpm run build
pnpm run start
```

### 后端部署

构建可执行文件：
```bash
cd backend
go build -o chatapp main.go
./chatapp
```

## 许可证

MIT License
