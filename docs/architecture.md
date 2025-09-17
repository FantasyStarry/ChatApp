# Repository-Service-Controller Architecture

## 🏗️ 架构改进总结

### 改进前的问题
- 控制器直接操作数据库
- 业务逻辑与数据访问紧耦合
- 代码重复和难以测试
- 违反单一职责原则

### 改进后的分层架构

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

## 📁 新的目录结构

```
ChatApp/
├── repository/                 # 数据访问层
│   ├── user_repository.go
│   ├── chatroom_repository.go
│   └── message_repository.go
├── service/                    # 业务逻辑层
│   ├── auth_service.go
│   ├── chatroom_service.go
│   └── message_service.go
├── controllers/                # 控制器层
│   ├── auth_controller.go
│   └── chatroom_controller.go
└── ... (其他现有目录)
```

## 🔧 各层职责

### Repository Layer (数据访问层)
- **职责**: 封装数据库操作
- **特点**: 
  - 提供接口抽象
  - 隐藏SQL细节
  - 支持多种数据源切换
  - 便于单元测试

```go
type UserRepository interface {
    Create(user *models.User) error
    GetByID(id uint) (*models.User, error)
    GetByUsername(username string) (*models.User, error)
    // ... 其他方法
}
```

### Service Layer (业务逻辑层)
- **职责**: 处理业务逻辑和规则
- **特点**:
  - 实现业务规则验证
  - 协调多个Repository
  - 处理事务管理
  - 提供可复用的业务操作

```go
type AuthService interface {
    Login(username, password string) (*models.User, string, error)
    CreateUser(user *models.User) error
    // ... 其他方法
}
```

### Controller Layer (控制器层)
- **职责**: 处理HTTP请求和响应
- **特点**:
  - 参数绑定和验证
  - 调用Service层方法
  - 格式化响应数据
  - 错误处理

```go
func (ctrl *AuthController) Login(c *gin.Context) {
    // 1. 绑定请求参数
    // 2. 调用Service层
    // 3. 处理响应
}
```

## ✅ 架构优势

### 1. **关注点分离**
- 每层专注于特定职责
- 代码更清晰、易维护

### 2. **可测试性**
- 每层可独立测试
- 支持依赖注入和Mock

### 3. **可扩展性**
- 易于添加新功能
- 支持水平和垂直扩展

### 4. **代码复用**
- Service层可被多个Controller使用
- Repository层可被多个Service使用

### 5. **技术栈无关**
- Repository层可切换数据库
- Service层业务逻辑不依赖具体技术

## 🚀 依赖注入

### 初始化顺序
```go
// 1. 初始化Repository
userRepo := repository.NewUserRepository(db)
chatRoomRepo := repository.NewChatRoomRepository(db)

// 2. 初始化Service (注入Repository)
authService := service.NewAuthService(userRepo)
chatRoomService := service.NewChatRoomService(chatRoomRepo, userRepo)

// 3. 初始化Controller (注入Service)
authController := controllers.NewAuthController(authService)
chatRoomController := controllers.NewChatRoomController(chatRoomService)
```

## 🔍 使用示例

### 用户登录流程
```go
// 1. Controller 接收请求
func (ctrl *AuthController) Login(c *gin.Context) {
    var req LoginRequest
    c.ShouldBindJSON(&req)
    
    // 2. 调用 Service 层
    user, token, err := ctrl.authService.Login(req.Username, req.Password)
    
    // 3. 返回响应
    c.JSON(http.StatusOK, LoginResponse{Token: token, User: *user})
}

// 4. Service 处理业务逻辑
func (s *authService) Login(username, password string) (*models.User, string, error) {
    // 验证用户名密码
    user, err := s.userRepo.GetByUsername(username)
    // 生成JWT Token
    token, err := utils.GenerateToken(user.ID, user.Username)
    return user, token, nil
}

// 5. Repository 执行数据库操作
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
    var user models.User
    err := r.db.Where("username = ?", username).First(&user).Error
    return &user, err
}
```

## 📊 性能和维护性改进

### Before (原架构)
```go
// 控制器直接操作数据库
func Login(c *gin.Context) {
    var req LoginRequest
    c.ShouldBindJSON(&req)
    
    var user models.User
    config.DB.Where("username = ?", req.Username).First(&user)  // 直接数据库操作
    
    if !utils.CheckPasswordHash(req.Password, user.Password) { // 业务逻辑混在控制器
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }
    // ... 更多代码
}
```

### After (新架构)
```go
// 清晰的分层架构
func (ctrl *AuthController) Login(c *gin.Context) {
    var req LoginRequest
    c.ShouldBindJSON(&req)
    
    user, token, err := ctrl.authService.Login(req.Username, req.Password) // 委托给Service
    if err != nil {
        c.JSON(401, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, LoginResponse{Token: token, User: *user})
}
```

## 🧪 测试改进

### Repository测试
```go
func TestUserRepository_GetByUsername(t *testing.T) {
    // 可以使用内存数据库或Mock
    repo := repository.NewUserRepository(testDB)
    user, err := repo.GetByUsername("testuser")
    assert.NoError(t, err)
    assert.Equal(t, "testuser", user.Username)
}
```

### Service测试
```go
func TestAuthService_Login(t *testing.T) {
    // Mock Repository
    mockRepo := &MockUserRepository{}
    service := service.NewAuthService(mockRepo)
    
    user, token, err := service.Login("testuser", "password")
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}
```

## 🎯 结论

这个架构改进带来了：

1. **更好的代码组织** - 清晰的职责分离
2. **更高的可测试性** - 每层可独立测试
3. **更强的可维护性** - 易于修改和扩展
4. **更好的代码复用** - 减少重复代码
5. **更灵活的技术选择** - 易于替换技术栈组件

这是现代Go Web应用的最佳实践架构，符合SOLID原则和清洁架构理念。