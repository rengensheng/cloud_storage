# Cloud Storage 服务部署指南

## 项目概述

基于Go语言和Gin框架构建的网盘文件存储服务，提供完备的RESTful API接口。

## 系统要求

### 开发环境
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Git

### 生产环境
- Docker 20.10+
- Docker Compose 2.0+
- 至少2GB RAM
- 至少10GB存储空间

## 快速开始

### 1. 克隆项目
```bash
git clone <repository-url>
cd cloud-storage
```

### 2. 配置环境变量
```bash
cp .env.example .env
# 编辑.env文件，配置数据库连接等信息
```

### 3. 使用Docker Compose运行
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f app
```

### 4. 访问服务
- API服务: http://localhost:8080
- 健康检查: http://localhost:8080/health
- API文档: http://localhost:8080/swagger/index.html (如果启用了Swagger)
- 数据库管理: http://localhost:8081 (Adminer)
- Redis管理: http://localhost:8082 (Redis Commander)

## 手动部署

### 1. 安装依赖
```bash
go mod download
```

### 2. 配置数据库
```bash
# 创建数据库
createdb cloud_storage

# 或使用PostgreSQL客户端
psql -U postgres -c "CREATE DATABASE cloud_storage;"
```

### 3. 运行数据库迁移
```bash
# 需要先实现迁移脚本
go run ./cmd/migrate
```

### 4. 构建和运行
```bash
# 构建
go build -o cloud-storage ./cmd/server

# 运行
./cloud-storage
```

## 配置说明

### 主要配置项

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| SERVER_PORT | 8080 | 服务器端口 |
| SERVER_HOST | 0.0.0.0 | 服务器主机 |
| DB_HOST | localhost | 数据库主机 |
| DB_PORT | 5432 | 数据库端口 |
| DB_NAME | cloud_storage | 数据库名称 |
| DB_USER | postgres | 数据库用户 |
| DB_PASSWORD | password | 数据库密码 |
| REDIS_HOST | localhost | Redis主机 |
| REDIS_PORT | 6379 | Redis端口 |
| JWT_SECRET | - | JWT密钥（必须修改） |
| STORAGE_PATH | ./storage/uploads | 文件存储路径 |
| MAX_UPLOAD_SIZE | 104857600 | 最大上传大小（100MB） |

### 安全配置建议

1. **修改JWT密钥**: 生产环境必须修改JWT_SECRET
2. **使用强密码**: 数据库和Redis使用强密码
3. **启用HTTPS**: 生产环境必须启用HTTPS
4. **配置防火墙**: 只开放必要端口
5. **定期备份**: 定期备份数据库和文件

## 数据库管理

### 备份数据库
```bash
# 使用pg_dump
docker-compose exec postgres pg_dump -U postgres cloud_storage > backup.sql

# 定时备份（crontab）
0 2 * * * docker-compose exec -T postgres pg_dump -U postgres cloud_storage > /backups/cloud_storage_$(date +\%Y\%m\%d).sql
```

### 恢复数据库
```bash
cat backup.sql | docker-compose exec -T postgres psql -U postgres cloud_storage
```

### 监控数据库
```bash
# 查看数据库大小
docker-compose exec postgres psql -U postgres -c "SELECT pg_size_pretty(pg_database_size('cloud_storage'));"

# 查看连接数
docker-compose exec postgres psql -U postgres -c "SELECT count(*) FROM pg_stat_activity;"
```

## 文件存储管理

### 存储路径结构
```
storage/
├── uploads/          # 用户上传的文件
│   └── {user_id}/    # 按用户ID分组
├── temp/             # 临时文件
└── versions/         # 文件版本（如果启用）
```

### 清理临时文件
```bash
# 清理7天前的临时文件
find ./storage/temp -type f -mtime +7 -delete

# 清理空目录
find ./storage -type d -empty -delete
```

### 存储配额监控
```bash
# 查看存储使用情况
du -sh ./storage/uploads

# 按用户查看存储使用
# 需要实现相应的脚本或查询
```

## 性能优化

### 数据库优化
1. **添加索引**: 在频繁查询的字段上添加索引
2. **分区表**: 大数据量时考虑表分区
3. **连接池**: 配置合适的数据库连接池大小

### Redis优化
1. **内存优化**: 配置最大内存和淘汰策略
2. **持久化**: 根据需求配置RDB或AOF
3. **集群**: 高并发时考虑Redis集群

### 应用优化
1. **GOMAXPROCS**: 设置合适的GOMAXPROCS
2. **连接复用**: 启用HTTP Keep-Alive
3. **缓存策略**: 合理使用缓存

## 监控和日志

### 日志配置
- 日志级别通过LOG_LEVEL环境变量控制
- 日志文件路径通过LOG_FILE环境变量配置
- 支持日志轮转（需要外部工具）

### 健康检查
```bash
curl http://localhost:8080/health
```

### 监控指标
建议监控以下指标:
- 系统资源（CPU、内存、磁盘）
- 数据库连接数
- API响应时间
- 错误率
- 存储使用率

## 高可用部署

### 架构设计
```
负载均衡器 (Nginx/Haproxy)
        ↓
    [应用服务器集群]
        ↓
    共享存储 (S3/NFS)
        ↓
[数据库集群]  [Redis集群]
```

### 部署步骤
1. 配置共享存储（如Amazon S3、MinIO）
2. 部署数据库集群（PostgreSQL流复制）
3. 部署Redis集群
4. 部署多个应用实例
5. 配置负载均衡器
6. 配置监控和告警

### 配置共享存储
修改环境变量:
```env
STORAGE_TYPE=s3
S3_BUCKET=your-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
```

## 安全建议

### 1. 网络安全
- 使用VPC/私有网络
- 配置安全组/防火墙
- 启用DDoS防护

### 2. 应用安全
- 定期更新依赖
- 启用CORS白名单
- 实施速率限制
- 验证文件类型

### 3. 数据安全
- 加密敏感数据
- 定期备份
- 访问日志审计
- 实施数据保留策略

## 故障排除

### 常见问题

#### 1. 服务启动失败
```bash
# 检查端口占用
netstat -tlnp | grep 8080

# 检查依赖服务
docker-compose ps

# 查看日志
docker-compose logs app
```

#### 2. 数据库连接失败
```bash
# 测试数据库连接
psql -h localhost -U postgres -d cloud_storage

# 检查数据库状态
docker-compose exec postgres pg_isready
```

#### 3. 文件上传失败
- 检查存储路径权限
- 检查磁盘空间
- 检查文件大小限制

#### 4. 性能问题
- 检查数据库索引
- 检查Redis连接
- 监控系统资源

### 日志分析
```bash
# 查看错误日志
grep -i error ./logs/app.log

# 查看慢查询
docker-compose exec postgres psql -U postgres -c "SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;"

# 查看API访问日志
tail -f ./logs/access.log
```

## 备份和恢复

### 完整备份方案
```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/$DATE"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
docker-compose exec -T postgres pg_dump -U postgres cloud_storage > $BACKUP_DIR/database.sql

# 备份文件
rsync -av ./storage/uploads $BACKUP_DIR/

# 备份配置
cp .env $BACKUP_DIR/

# 打包备份
tar -czf /backups/cloud-storage-$DATE.tar.gz $BACKUP_DIR

# 清理旧备份（保留最近30天）
find /backups -name "*.tar.gz" -mtime +30 -delete

echo "备份完成: /backups/cloud-storage-$DATE.tar.gz"
```

### 恢复方案
```bash
#!/bin/bash
# restore.sh

BACKUP_FILE=$1
RESTORE_DIR="/tmp/restore_$(date +%s)"

# 解压备份
tar -xzf $BACKUP_FILE -C $RESTORE_DIR

# 恢复数据库
cat $RESTORE_DIR/database.sql | docker-compose exec -T postgres psql -U postgres cloud_storage

# 恢复文件
rsync -av $RESTORE_DIR/uploads/ ./storage/uploads/

# 恢复配置（谨慎操作）
# cp $RESTORE_DIR/.env ./

echo "恢复完成"
```

## 扩展和定制

### 添加新功能
1. 在`internal/models`中添加数据模型
2. 在`internal/repositories`中添加数据访问层
3. 在`internal/services`中添加业务逻辑
4. 在`internal/handlers`中添加HTTP处理器
5. 在`internal/middleware`中添加中间件（如果需要）

### 集成第三方服务
1. **邮件服务**: 用于用户注册、密码重置
2. **短信服务**: 用于双因素认证
3. **对象存储**: 替换本地存储为S3/MinIO
4. **CDN**: 加速文件下载
5. **监控服务**: Prometheus + Grafana

### 性能测试
```bash
# 使用ab进行压力测试
ab -n 1000 -c 100 http://localhost:8080/health

# 使用wrk进行性能测试
wrk -t12 -c400 -d30s http://localhost:8080/health
```

## 支持与维护

### 获取帮助
- 查看文档: `README.md` 和 `README_API.md`
- 检查日志: `./logs/` 目录
- 提交Issue: [项目Issue跟踪]

### 版本升级
1. 备份当前版本
2. 拉取新代码
3. 更新依赖: `go mod download`
4. 运行数据库迁移
5. 重启服务

### 贡献指南
1. Fork项目
2. 创建特性分支
3. 提交更改
4. 创建Pull Request

## 许可证

本项目采用MIT许可证。详见[LICENSE](LICENSE)文件。

---

*最后更新: 2023年12月*