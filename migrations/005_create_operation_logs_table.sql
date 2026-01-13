-- 005_create_operation_logs_table.sql
-- 创建操作日志表

CREATE TABLE IF NOT EXISTS operation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    operation_type VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    details JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'success',
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_operation_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_operation_type ON operation_logs(operation_type);
CREATE INDEX IF NOT EXISTS idx_operation_logs_resource_id ON operation_logs(resource_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_resource_type ON operation_logs(resource_type);
CREATE INDEX IF NOT EXISTS idx_operation_logs_status ON operation_logs(status);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_operation_logs_user_id_created_at ON operation_logs(user_id, created_at);

-- 添加注释
COMMENT ON TABLE operation_logs IS '操作日志表，记录用户的操作历史';
COMMENT ON COLUMN operation_logs.operation_type IS '操作类型：upload, download, delete, share, rename, move, copy等';
COMMENT ON COLUMN operation_logs.resource_type IS '资源类型：file, directory, user, share等';
COMMENT ON COLUMN operation_logs.details IS '操作详情，JSON格式存储';
COMMENT ON COLUMN operation_logs.status IS '操作状态：success, failed';
