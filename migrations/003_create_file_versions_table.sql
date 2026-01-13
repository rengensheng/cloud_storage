-- 003_create_file_versions_table.sql
-- 创建文件版本表

CREATE TABLE IF NOT EXISTS file_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID NOT NULL,
    version_number INTEGER NOT NULL,
    file_size BIGINT NOT NULL,
    file_hash VARCHAR(255),
    storage_path VARCHAR(512) NOT NULL,
    mime_type VARCHAR(100),
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_file_versions_file FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    CONSTRAINT fk_file_versions_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_file_version UNIQUE (file_id, version_number)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_file_versions_file_id ON file_versions(file_id);
CREATE INDEX IF NOT EXISTS idx_file_versions_version_number ON file_versions(file_id, version_number);
CREATE INDEX IF NOT EXISTS idx_file_versions_created_by ON file_versions(created_by);
