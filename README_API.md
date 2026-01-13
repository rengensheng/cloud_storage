# Cloud Storage API 使用指南

本文档提供 Cloud Storage 服务的 API 使用示例和说明。

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: Bearer Token (JWT)
- **默认端口**: 8080

## 快速开始

### 1. 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

响应示例:
```json
{
  "message": "user registered successfully",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "testuser",
    "email": "test@example.com",
    "role": "user",
    "storage_quota": 10737418240,
    "used_storage": 0,
    "is_active": true,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. 使用访问令牌

将获取到的访问令牌添加到请求头中:

```bash
export ACCESS_TOKEN="your_access_token_here"

curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## 文件操作示例

### 1. 创建目录

```bash
curl -X POST http://localhost:8080/api/v1/files \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "documents",
    "type": "directory"
  }'
```

### 2. 上传文件

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -F "file=@/path/to/your/file.pdf" \
  -F "parent_id=123e4567-e89b-12d3-a456-426614174000" \
  -F "is_public=false" \
  -F "override=false"
```

### 3. 获取文件列表

```bash
curl -X GET "http://localhost:8080/api/v1/files?page=1&page_size=20&sort_by=name&sort_order=asc" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

### 4. 下载文件

```bash
curl -X GET http://localhost:8080/api/v1/files/{file_id}/download \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  --output downloaded_file.pdf
```

### 5. 更新文件信息

```bash
curl -X PUT http://localhost:8080/api/v1/files/{file_id} \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "renamed_file.pdf",
    "is_public": true
  }'
```

### 6. 删除文件

```bash
# 软删除（移动到回收站）
curl -X DELETE http://localhost:8080/api/v1/files/{file_id} \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 永久删除
curl -X DELETE "http://localhost:8080/api/v1/files/{file_id}?permanent=true" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## 回收站操作

### 1. 查看回收站文件

```bash
curl -X GET "http://localhost:8080/api/v1/recycle?page=1&page_size=20" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

### 2. 恢复文件

```bash
curl -X POST http://localhost:8080/api/v1/recycle/{file_id}/restore \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

### 3. 清理回收站

```bash
# 清理30天前的文件
curl -X DELETE "http://localhost:8080/api/v1/recycle/cleanup?days=30" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## 搜索文件

```bash
curl -X GET "http://localhost:8080/api/v1/search?q=document&search_in=name&page=1&page_size=20" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## 系统信息

### 1. 健康检查

```bash
curl -X GET http://localhost:8080/health
```

### 2. 获取存储使用情况

```bash
curl -X GET http://localhost:8080/api/v1/stats/storage \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

响应示例:
```json
{
  "used": 104857600,
  "quota": 10737418240,
  "available": 10632560640,
  "usage_percent": 0.98,
  "usage_readable": "100 MB / 10 GB"
}
```

### 3. 获取文件统计

```bash
curl -X GET http://localhost:8080/api/v1/stats/files \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## 错误处理

API 使用标准的 HTTP 状态码:

- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证或令牌无效
- `403 Forbidden`: 权限不足
- `404 Not Found`: 资源不存在
- `409 Conflict`: 资源冲突（如文件名重复）
- `413 Payload Too Large`: 文件太大
- `429 Too Many Requests`: 请求频率限制
- `500 Internal Server Error`: 服务器内部错误

错误响应格式:
```json
{
  "error": "Error message description"
}
```

## 速率限制

API 有默认的速率限制:
- 每个 IP 地址每分钟最多 100 次请求
- 文件上传大小限制: 100MB

超过限制会返回 `429 Too Many Requests` 状态码。

## 分页和排序

列表接口支持分页和排序:

### 分页参数
- `page`: 页码（从1开始）
- `page_size`: 每页大小（1-100，默认20）

### 排序参数
- `sort_by`: 排序字段（name, size, created_at, updated_at）
- `sort_order`: 排序顺序（asc, desc）

示例:
```
/files?page=1&page_size=20&sort_by=name&sort_order=asc
```

## 环境变量配置

服务支持以下环境变量配置:

```bash
# 服务器配置
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
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRE_HOURS=24

# 存储配置
STORAGE_PATH=./storage/uploads
MAX_UPLOAD_SIZE=104857600  # 100MB
```

## Docker 部署

使用 Docker Compose 快速部署:

```bash
# 复制环境变量文件
cp .env.example .env

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

## 开发环境

### 本地运行

```bash
# 安装依赖
go mod download

# 运行数据库迁移
go run ./cmd/migrate

# 启动服务
go run ./cmd/server

# 或使用 Makefile
make dev
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test ./internal/services -v

# 生成测试覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## API 文档

Swagger UI 文档（如果启用了）:
- URL: `http://localhost:8080/swagger/index.html`
- OpenAPI 规范: `http://localhost:8080/swagger/doc.json`

## 客户端 SDK

### Go 客户端示例

```go
package main

import (
    "context"
    "fmt"
    "net/http"

    "github.com/google/uuid"
)

type CloudStorageClient struct {
    baseURL    string
    token      string
    httpClient *http.Client
}

func NewCloudStorageClient(baseURL, token string) *CloudStorageClient {
    return &CloudStorageClient{
        baseURL:    baseURL,
        token:      token,
        httpClient: &http.Client{},
    }
}

func (c *CloudStorageClient) UploadFile(ctx context.Context, filePath string) error {
    // 实现文件上传逻辑
    return nil
}

func (c *CloudStorageClient) DownloadFile(ctx context.Context, fileID uuid.UUID, savePath string) error {
    // 实现文件下载逻辑
    return nil
}

// 更多客户端方法...
```

### Python 客户端示例

```python
import requests
import uuid

class CloudStorageClient:
    def __init__(self, base_url, token):
        self.base_url = base_url
        self.token = token
        self.headers = {
            "Authorization": f"Bearer {token}",
            "Content-Type": "application/json"
        }

    def upload_file(self, file_path, parent_id=None, is_public=False):
        url = f"{self.base_url}/upload"
        files = {"file": open(file_path, "rb")}
        data = {
            "parent_id": parent_id,
            "is_public": is_public
        }

        response = requests.post(url, headers=self.headers, files=files, data=data)
        return response.json()

    def download_file(self, file_id, save_path):
        url = f"{self.base_url}/files/{file_id}/download"
        response = requests.get(url, headers=self.headers, stream=True)

        with open(save_path, "wb") as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)

        return save_path

# 使用示例
client = CloudStorageClient("http://localhost:8080/api/v1", "your_access_token")
client.upload_file("/path/to/file.pdf")
```

## 常见问题

### Q: 如何重置密码？
A: 目前需要通过管理员或实现密码重置功能。

### Q: 支持哪些文件类型？
A: 支持所有文件类型，系统会根据文件扩展名自动识别MIME类型。

### Q: 最大文件大小是多少？
A: 默认100MB，可通过环境变量 `MAX_UPLOAD_SIZE` 配置。

### Q: 如何备份数据？
A: 使用数据库备份工具（如 pg_dump）备份 PostgreSQL 数据库，并备份存储目录。

### Q: 支持集群部署吗？
A: 是的，可以通过配置共享存储（如S3）和负载均衡实现集群部署。

## 联系支持

如有问题或建议，请通过以下方式联系:
- GitHub Issues: [项目地址]
- 邮箱: support@cloud-storage.example.com

---

*文档最后更新: 2023年12月*