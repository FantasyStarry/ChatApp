# ChatApp 多存储后端支持文档

## 概述

ChatApp 现在支持多种文件存储后端，你可以通过配置文件轻松切换存储方式：
- **MinIO**: 本地或私有云对象存储
- **七牛云**: 公有云对象存储服务

## 配置方式

### 1. 基本配置结构

在 `config.yaml` 文件中添加以下配置：

```yaml
storage:
  type: "qiniu"  # 选择存储类型："minio" 或 "qiniu"

# MinIO 配置（当 storage.type = "minio" 时使用）
minio:
  endpoint: "127.0.0.1:9000"
  access_key: "minioadmin"
  secret_key: "minioadmin"
  bucket_name: "chatapp"
  use_ssl: false
  region: "us-east-1"

# 七牛云配置（当 storage.type = "qiniu" 时使用）
qiniu:
  access_key: "你的七牛云AccessKey"
  secret_key: "你的七牛云SecretKey"
  bucket: "你的存储空间名称"
  domain: "你的自定义域名"
  region: "south-china"  # 存储区域
  use_https: true
```

### 2. 七牛云存储区域

支持的区域选项：
- `east-china`: 华东-浙江
- `north-china`: 华北-河北
- `south-china`: 华南-广东
- `north-america`: 北美
- `southeast-asia`: 东南亚

### 3. 你当前的七牛云配置

```yaml
storage:
  type: "qiniu"

qiniu:
  access_key: "LYc71kKaC3uOHOF3SUoIm9GayQCdd3_RaPw5JWZ4"
  secret_key: "p48XPt-eJ2WHCiuVf3lfCLWQliY8iTCHKPzbO60W"
  bucket: "spongzi"
  domain: "t2u2huu55.hn-bkt.clouddn.com"
  region: "south-china"
  use_https: true
```

## API 使用

所有文件相关的API保持不变，存储后端的切换对应用层完全透明：

### 文件上传
```bash
POST /api/files/upload
Content-Type: multipart/form-data

# Form data:
# file: [文件]
# chatroom_id: [聊天室ID]
```

### 文件下载
```bash
GET /api/files/download/{file_id}
```

### 获取聊天室文件列表
```bash
GET /api/files/chatroom/{chatroom_id}?page=1&page_size=20
```

### 获取用户文件列表
```bash
GET /api/files/my
```

### 删除文件
```bash
DELETE /api/files/{file_id}
```

### 获取文件信息
```bash
GET /api/files/{file_id}
```

### 获取上传预签名URL
```bash
GET /api/files/upload-url?filename=test.jpg&chatroom_id=1
```

## 存储后端特性对比

| 特性 | MinIO | 七牛云 |
|------|-------|--------|
| 部署方式 | 本地/私有云 | 公有云服务 |
| 访问控制 | 私有访问 | 支持公有/私有访问 |
| CDN加速 | 需自行配置 | 内置CDN |
| 成本 | 服务器成本 | 按使用量计费 |
| 可用性 | 依赖本地环境 | 99.9%+ |
| 预签名URL | 支持 | 支持 |
| 区域支持 | 单一部署区域 | 多区域支持 |

## 切换存储后端

如果你想从一种存储后端切换到另一种：

1. **修改配置文件**：更改 `storage.type` 值
2. **重启应用**：新的存储配置将在下次启动时生效
3. **数据迁移**（可选）：如果需要迁移已有文件，需要编写专门的迁移脚本

## 扩展新的存储后端

如果你想添加新的存储后端（如阿里云OSS、AWS S3等），只需要：

1. 在 `storage/` 目录下创建新的适配器文件
2. 实现 `Storage` 接口
3. 在 `storage/factory.go` 中添加新的存储类型
4. 在 `config/config.go` 中添加相应配置结构

## 安全注意事项

1. **配置文件安全**：
   - 不要将包含真实密钥的 `config.yaml` 提交到版本控制
   - 使用 `config.example.yaml` 作为模板
   - 生产环境中使用环境变量或安全的配置管理

2. **访问控制**：
   - 七牛云建议使用私有空间配合预签名URL
   - 定期轮换访问密钥

3. **网络安全**：
   - 生产环境建议启用HTTPS
   - 配置适当的CORS策略

## 故障排除

### 常见问题

1. **七牛云上传失败**：
   - 检查AccessKey和SecretKey是否正确
   - 确认存储空间名称和区域设置
   - 验证域名配置是否正确

2. **MinIO连接失败**：
   - 检查MinIO服务是否运行
   - 验证网络连接和端口
   - 确认bucket是否存在

3. **文件访问权限**：
   - 检查存储空间的访问权限设置
   - 确认预签名URL的有效期

### 日志调试

应用启动时会显示当前使用的存储类型：
```
Configuration loaded successfully!
Storage type: qiniu
```

启用debug模式可以查看更详细的存储操作日志。