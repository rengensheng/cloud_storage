# 云存储服务 (Cloud Storage Service)

基于Go语言和Gin框架构建的网盘文件存储服务，提供完备的RESTful API接口。

## 功能特性

### 核心功能
- ✅ 文件上传/下载
- ✅ 文件管理（列表、重命名、移动、复制、删除）
- ✅ 文件夹管理（创建、删除、重命名）
- ✅ 文件分享（生成分享链接、密码保护、访问控制、过期设置）
- ✅ 回收站功能（软删除、恢复、永久删除）

### 用户管理
- ✅ 用户注册、登录、注销
- ✅ JWT身份验证和令牌刷新
- ✅ 角色权限管理（管理员、普通用户）
- ✅ 用户配额管理
- ✅ 用户状态管理（激活/停用）

### 高级功能
- ✅ 文件版本控制（自动创建版本、恢复历史版本）
- ✅ 文件搜索和过滤（按名称搜索）
- ✅ 操作日志记录（完整的操作审计）
- ✅ 系统管理（用户管理、系统统计）
- ✅ 数据库迁移工具（自动创建表结构）

## 技术栈

### 后端
- **编程语言**: Go 1.21+
- **Web框架**: Gin
- **数据库**: PostgreSQL（主数据库） + Redis（缓存）
- **对象存储**: 本地文件系统（可扩展为S3/MinIO）
- **认证授权**: JWT

### 开发工具
- **依赖管理**: Go Modules
- **API文档**: Swagger/OpenAPI
- **测试框架**: Go标准测试包 + testify
- **代码质量**: golangci-lint

## 项目结构

```
cloud-storage/
├── cmd/
│   ├── server/
│   │   └── main.go              # 应用入口
│   └── migrate/
│       └── main.go              # 数据库迁移工具
├── internal/
│   ├── config/                  # 配置管理
│   ├── database/               # 数据库连接
│   ├── models/                 # 数据模型
│   │   ├── user.go
│   │   ├── file.go
│   │   ├── file_version.go
│   │   ├── share.go
│   │   ├── operation_log.go
│   │   └── upload.go
│   ├── repositories/           # 数据访问层
│   │   ├── user_repository.go
│   │   ├── file_repository.go
│   │   ├── file_version_repository.go
│   │   ├── share_repository.go
│   │   └── operation_log_repository.go
│   ├── services/              # 业务逻辑层
│   │   ├── file_service.go
│   │   ├── share_service.go
│   │   └── operation_log_service.go
│   ├── handlers/              # HTTP处理器
│   │   ├── auth_handler.go
│   │   ├── file_handler.go
│   │   ├── share_handler.go
│   │   └── admin_handler.go
│   ├── middleware/            # 中间件
│   │   └── auth_middleware.go
│   └── pkg/                   # 可复用包
│       └── storage/           # 存储抽象层
├── migrations/               # SQL迁移文件
│   ├── 001_create_users_table.sql
│   ├── 002_create_files_table.sql
│   ├── 003_create_file_versions_table.sql
│   ├── 004_create_shares_table.sql
│   └── 005_create_operation_logs_table.sql
├── storage/                   # 文件存储目录
│   ├── uploads/              # 上传文件
│   └── temp/                 # 临时文件
├── .env.example              # 环境变量示例
├── go.mod                    # Go模块定义
├── go.sum                    # 依赖校验
├── docker-compose.yml        # Docker编排
├── Dockerfile                # Docker构建文件
├── Makefile                  # 构建脚本
└── README.md                 # 项目文档
```

## 数据模型设计

### 用户表 (users)
```sql
id, username, email, password_hash, role, storage_quota, used_storage,
created_at, updated_at, last_login_at, is_active
```

### 文件表 (files)
```sql
id, user_id, parent_id, name, path, size, mime_type, hash,
type, is_public, share_token, version, deleted_at,
created_at, updated_at
```

### 文件版本表 (file_versions)
```sql
id, file_id, version_number, file_size, file_hash, storage_path,
mime_type, created_by, created_at
```

### 分享表 (shares)
```sql
id, file_id, user_id, share_token, password_hash, access_type,
expires_at, max_downloads, download_count, is_active,
created_at, updated_at
```

### 操作日志表 (operation_logs)
```sql
id, user_id, operation, resource_type, resource_id,
result, details, error_message, ip_address, user_agent,
duration, created_at
```

### 上传会话表 (upload_sessions)
```sql
id, user_id, file_name, file_size, file_hash, parent_id,
chunk_size, total_chunks, uploaded_chunks, storage_path,
mime_type, status, error_message, created_at, updated_at, expires_at
```

## API接口设计

### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户注销
- `POST /api/v1/auth/refresh` - 刷新令牌
- `GET /api/v1/auth/profile` - 获取用户信息
- `PUT /api/v1/auth/profile` - 更新用户信息
- `PUT /api/v1/auth/password` - 修改密码

### 文件操作
- `GET /api/v1/files` - 获取文件列表
- `GET /api/v1/files/{id}` - 获取文件详情
- `POST /api/v1/files` - 创建文件/文件夹
- `PUT /api/v1/files/{id}` - 更新文件信息
- `DELETE /api/v1/files/{id}` - 删除文件
- `POST /api/v1/files/{id}/copy` - 复制文件
- `POST /api/v1/files/{id}/move` - 移动文件
- `GET /api/v1/files/{id}/download` - 下载文件
- `GET /api/v1/files/{id}/versions` - 获取文件版本列表
- `POST /api/v1/files/{id}/restore-version` - 恢复文件版本

### 文件上传
- `POST /api/v1/upload` - 文件上传
- `POST /api/v1/upload/chunk` - 分片上传

### 回收站操作
- `GET /api/v1/recycle` - 查看回收站文件
- `POST /api/v1/recycle/{id}/restore` - 恢复文件
- `DELETE /api/v1/recycle/cleanup` - 清理回收站

### 分享管理
- `POST /api/v1/shares` - 创建分享
- `GET /api/v1/shares` - 获取分享列表
- `GET /api/v1/shares/{id}` - 获取分享详情
- `PUT /api/v1/shares/{id}` - 更新分享
- `DELETE /api/v1/shares/{id}` - 删除分享
- `POST /api/v1/shares/batch-delete` - 批量删除分享
- `GET /api/v1/shares/stats` - 获取分享统计
- `GET /api/v1/s/{token}` - 访问分享（公开）
- `GET /api/v1/s/{token}/download` - 下载分享文件

### 搜索和统计
- `GET /api/v1/search` - 搜索文件
- `GET /api/v1/stats/storage` - 获取存储使用情况
- `GET /api/v1/stats/files` - 获取文件统计

### 系统管理
- `GET /api/v1/admin/stats` - 系统统计信息
- `GET /api/v1/admin/users` - 获取用户列表
- `GET /api/v1/admin/users/{id}` - 获取用户详情
- `PUT /api/v1/admin/users/{id}` - 更新用户信息
- `DELETE /api/v1/admin/users/{id}` - 删除用户
- `POST /api/v1/admin/users/{id}/activate` - 激活用户
- `POST /api/v1/admin/users/{id}/deactivate` - 停用用户

### 操作日志
- `GET /api/v1/logs` - 获取操作日志
- `GET /api/v1/logs/stats` - 获取日志统计
- `DELETE /api/v1/logs/cleanup` - 清理过期日志（管理员）

## 配置说明

### 环境变量
```env
# 服务器配置
APP_ENV=development
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cloud_storage
DB_USER=postgres
DB_PASSWORD=password

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

# 存储配置
STORAGE_PATH=./storage/uploads
MAX_UPLOAD_SIZE=104857600  # 100MB
ENABLE_CHUNK_UPLOAD=true
```

## 部署方式

### 本地开发
```bash
# 安装依赖
go mod download

# 启动数据库
docker-compose up -d postgres redis

# 运行数据库迁移
go run cmd/migrate/main.go

# 设置默认管理员账户（可选）
export ADMIN_USERNAME=admin
export ADMIN_EMAIL=admin@cloud-storage.local
export ADMIN_PASSWORD=your_secure_password

# 重新运行迁移创建管理员
go run cmd/migrate/main.go

# 启动服务
go run cmd/server/main.go
```

### 数据库迁移
```bash
# 执行所有数据库迁移
go run cmd/migrate/main.go

# 查看回滚信息（仅显示，不执行）
go run cmd/migrate/main.go --rollback
```

### Docker部署
```bash
# 构建镜像
docker build -t cloud-storage .

# 运行容器
docker-compose up -d
```

### Kubernetes部署
```bash
# 部署到Kubernetes
kubectl apply -f k8s/
```

## 开发计划

### 第一阶段 (基础功能) ✅
- [x] 项目结构搭建
- [x] 用户认证系统
- [x] 文件上传下载
- [x] 基础文件管理

### 第二阶段 (核心功能) ✅
- [x] 文件夹管理
- [x] 文件分享功能
- [x] 回收站功能
- [x] 配额管理

### 第三阶段 (高级功能) ✅
- [x] 文件版本控制
- [x] 文件搜索和过滤
- [x] 操作日志
- [x] 系统管理功能
- [x] 数据库迁移工具

### 第四阶段 (优化和扩展) 🚧
- [ ] 分片上传完整实现
- [ ] 断点续传下载
- [ ] 文件预览（图片、文档、视频）
- [ ] 邮箱验证和密码重置
- [ ] 性能优化
- [ ] 监控和日志
- [ ] 容器化部署
- [ ] 压力测试

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 已实现功能详情

### 用户系统 ✅
- 用户注册/登录/注销
- JWT令牌认证和刷新
- 角色权限（user/admin）
- 用户配额管理
- 用户状态激活/停用
- 个人信息管理

### 文件管理 ✅
- 文件上传/下载
- 文件列表和分页
- 文件重命名/移动/复制
- 文件删除（软删除）
- 文件夹创建和管理
- 文件搜索（按名称）
- 存储使用统计
- 文件类型过滤

### 分享功能 ✅
- 创建分享链接
- 密码保护
- 访问类型控制（view/download/edit）
- 过期时间设置
- 下载次数限制
- 分享列表管理
- 批量删除分享
- 分享统计信息
- 公开访问分享

### 版本控制 ✅
- 自动创建文件版本
- 查看版本历史
- 恢复到指定版本
- 版本元数据管理
- 版本对比功能

### 回收站 ✅
- 软删除到回收站
- 查看回收站文件
- 恢复已删除文件
- 永久删除
- 批量清理旧文件

### 操作日志 ✅
- 记录所有用户操作
- 操作类型分类（上传、下载、删除等）
- 操作结果追踪
- IP地址和UserAgent记录
- 日志查询和过滤
- 用户操作统计
- 自动清理过期日志

### 系统管理 ✅
- 系统统计信息（用户数、文件数、存储使用等）
- 用户列表和详情
- 用户信息更新
- 用户激活/停用
- 用户配额管理
- 操作日志查询
- 日志清理功能

### 数据库工具 ✅
- 自动迁移表结构
- 索引创建
- 默认管理员账户创建
- 迁移工具命令行

## 快速开始指南

### 第一步：准备环境
```bash
# 安装Go 1.21+
# 安装PostgreSQL 12+
# 安装Redis 6+ (可选）
```

### 第二步：配置环境
```bash
cp .env.example .env
# 编辑.env文件，配置数据库、Redis等
```

### 第三步：数据库迁移
```bash
# 运行迁移
go run cmd/migrate/main.go
```

### 第四步：启动服务
```bash
# 启动PostgreSQL和Redis
docker-compose up -d

# 启动应用
go run cmd/server/main.go
```

### 第五步：访问服务
```bash
# 健康检查
curl http://localhost:8080/health

# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123"}'
```

## 联系方式

如有问题或建议，请通过以下方式联系：
- 提交 Issue
- 发送邮件至 support@cloud-storage.example.com