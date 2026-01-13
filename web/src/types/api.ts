export interface User {
  id: string;
  username: string;
  email: string;
  role: 'user' | 'admin';
  storage_quota: number;
  used_storage: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  last_login_at?: string;
}

export interface Tokens {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
}

export interface AuthResponse {
  message: string;
  user: User;
  tokens: Tokens;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface File {
  id: string;
  user_id: string;
  parent_id: string | null;
  name: string;
  path: string;
  size: number;
  mime_type: string;
  hash: string;
  type: 'file' | 'directory';
  is_public: boolean;
  share_token: string | null;
  version: number;
  deleted_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface FileListResponse {
  files: File[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface CreateFileRequest {
  name: string;
  type: 'file' | 'directory';
  parent_id?: string;
}

export interface UpdateFileRequest {
  name?: string;
  is_public?: boolean;
}

export interface Share {
  id: string;
  file_id: string;
  user_id: string;
  share_token: string;
  password_hash: string | null;
  access_type: 'view' | 'download' | 'edit';
  expires_at: string | null;
  max_downloads: number | null;
  download_count: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  file?: File;
}

export interface CreateShareRequest {
  file_id: string;
  password?: string;
  access_type: 'view' | 'download' | 'edit';
  expires_at?: string;
  max_downloads?: number;
}

export interface UpdateShareRequest {
  password?: string;
  access_type?: 'view' | 'download' | 'edit';
  expires_at?: string;
  max_downloads?: number;
}

export interface ShareListResponse {
  shares: Share[];
  total: number;
  page: number;
  page_size: number;
}

export interface ShareStats {
  total_shares: number;
  active_shares: number;
  expired_shares: number;
  total_downloads: number;
}

export interface FileVersion {
  id: string;
  file_id: string;
  version_number: number;
  file_size: number;
  file_hash: string;
  storage_path: string;
  mime_type: string;
  created_by: string;
  created_at: string;
}

export interface OperationLog {
  id: string;
  user_id: string;
  operation: string;
  resource_type: string;
  resource_id: string;
  result: 'success' | 'failure';
  details: string;
  error_message: string | null;
  ip_address: string;
  user_agent: string;
  duration: number;
  created_at: string;
}

export interface OperationLogResponse {
  logs: OperationLog[];
  total: number;
  page: number;
  page_size: number;
}

export interface StorageStats {
  used: number;
  quota: number;
  available: number;
  usage_percent: number;
  usage_readable: string;
}

export interface FileStats {
  total_files: number;
  total_directories: number;
  total_size: number;
  public_files: number;
  shared_files: number;
  recent_files: File[];
}

export interface SystemStats {
  total_users: number;
  total_files: number;
  total_directories: number;
  total_storage_used: number;
  total_storage_quota: number;
  active_shares: number;
  operation_logs_today: number;
}

export interface PaginatedParams {
  page?: number;
  page_size?: number;
}

export interface SortParams {
  sort_by?: 'name' | 'size' | 'created_at' | 'updated_at';
  sort_order?: 'asc' | 'desc';
}

export interface FileListParams extends PaginatedParams, SortParams {
  parent_id?: string;
  type?: 'file' | 'directory';
}

export interface SearchParams extends PaginatedParams {
  q: string;
  search_in?: 'name' | 'path' | 'all';
}

export interface LogParams extends PaginatedParams {
  user_id?: string;
  operation?: string;
  result?: 'success' | 'failure';
  start_date?: string;
  end_date?: string;
}

export interface ApiResponse<T = any> {
  data?: T;
  error?: string;
  message?: string;
}

export interface AdminUser extends User {
  file_count: number;
  storage_used: number;
}

export interface AdminUserListResponse {
  users: AdminUser[];
  total: number;
  page: number;
  page_size: number;
}
