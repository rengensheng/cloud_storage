import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  File,
  FileListResponse,
  CreateFileRequest,
  UpdateFileRequest,
  Share,
  ShareListResponse,
  CreateShareRequest,
  UpdateShareRequest,
  ShareStats,
  FileVersion,
  OperationLogResponse,
  StorageStats,
  FileStats,
  SystemStats,
  SearchParams,
  FileListParams,
  LogParams,
  AdminUserListResponse,
  AdminUser,
} from '../types/api';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

class ApiClient {
  private getAuthHeaders(): Record<string, string> {
    const token = localStorage.getItem('access_token');
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...this.getAuthHeaders(),
      ...options.headers,
    };

    const response = await fetch(url, { ...options, headers });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  private async requestWithFile<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = {
      ...this.getAuthHeaders(),
    };

    const response = await fetch(url, { ...options, headers });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  async uploadFile(
    file: any,
    parentId?: string,
    isPublic: boolean = false,
    override: boolean = false
  ): Promise<any> {
    const formData = new FormData();
    formData.append('file', file);
    if (parentId) formData.append('parent_id', parentId);
    formData.append('is_public', isPublic.toString());
    formData.append('override', override.toString());

    return this.requestWithFile('/upload', {
      method: 'POST',
      body: formData,
    });
  }

  async downloadFile(fileId: string): Promise<Blob> {
    const url = `${API_BASE_URL}/files/${fileId}/download`;
    const headers = this.getAuthHeaders();
    const response = await fetch(url, { headers });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || `HTTP error! status: ${response.status}`);
    }

    return response.blob();
  }

  async refreshToken(): Promise<AuthResponse> {
    const refreshToken = localStorage.getItem('refresh_token');
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${refreshToken}`,
      },
    });

    if (!response.ok) {
      throw new Error('Failed to refresh token');
    }

    const data = await response.json();
    localStorage.setItem('access_token', data.tokens.access_token);
    localStorage.setItem('refresh_token', data.tokens.refresh_token);
    return data;
  }
}

export const apiClient = new ApiClient();

export const authApi = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    localStorage.setItem('access_token', response.tokens.access_token);
    localStorage.setItem('refresh_token', response.tokens.refresh_token);
    localStorage.setItem('user_id', response.user.id);
    return response;
  },

  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await apiClient.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    localStorage.setItem('access_token', response.tokens.access_token);
    localStorage.setItem('refresh_token', response.tokens.refresh_token);
    localStorage.setItem('user_id', response.user.id);
    return response;
  },

  logout: async (): Promise<void> => {
    await apiClient.request<void>('/auth/logout', { method: 'POST' });
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user_id');
  },

  getProfile: async (): Promise<any> => {
    return apiClient.request<any>('/auth/profile');
  },

  updateProfile: async (data: any): Promise<any> => {
    return apiClient.request<any>('/auth/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  changePassword: async (data: { old_password: string; new_password: string }): Promise<any> => {
    return apiClient.request<any>('/auth/password', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },
};

export const fileApi = {
  getFiles: async (params?: FileListParams): Promise<FileListResponse> => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.page_size) searchParams.append('page_size', params.page_size.toString());
    if (params?.parent_id) searchParams.append('parent_id', params.parent_id);
    if (params?.type) searchParams.append('type', params.type);
    if (params?.sort_by) searchParams.append('sort_by', params.sort_by);
    if (params?.sort_order) searchParams.append('sort_order', params.sort_order);

    const response = await apiClient.request<any>(`/files?${searchParams.toString()}`);
    return {
      files: response.files || [],
      total: response.total || 0,
      page: response.page || 1,
      page_size: response.size || 20,
      total_pages: Math.ceil((response.total || 0) / (response.size || 20)),
    };
  },

  getFile: async (id: string): Promise<File> => {
    return apiClient.request<File>(`/files/${id}`);
  },

  createFile: async (data: CreateFileRequest): Promise<File> => {
    return apiClient.request<File>('/files', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  updateFile: async (id: string, data: UpdateFileRequest): Promise<File> => {
    return apiClient.request<File>(`/files/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  deleteFile: async (id: string, permanent: boolean = false): Promise<void> => {
    const params = permanent ? '?permanent=true' : '';
    await apiClient.request<void>(`/files/${id}${params}`, { method: 'DELETE' });
  },

  copyFile: async (id: string, targetParentId: string): Promise<File> => {
    return apiClient.request<File>(`/files/${id}/copy`, {
      method: 'POST',
      body: JSON.stringify({ target_parent_id: targetParentId }),
    });
  },

  moveFile: async (id: string, targetParentId: string): Promise<File> => {
    return apiClient.request<File>(`/files/${id}/move`, {
      method: 'POST',
      body: JSON.stringify({ target_parent_id: targetParentId }),
    });
  },

  downloadFile: async (id: string): Promise<Blob> => {
    return apiClient.downloadFile(id);
  },

  getVersions: async (id: string): Promise<FileVersion[]> => {
    const response = await apiClient.request<{ versions: FileVersion[] }>(`/files/${id}/versions`);
    return response.versions || [];
  },

  restoreVersion: async (id: string, versionNumber: number): Promise<File> => {
    return apiClient.request<File>(`/files/${id}/restore-version`, {
      method: 'POST',
      body: JSON.stringify({ version_number: versionNumber }),
    });
  },

  uploadFile: (
    file: File,
    parentId?: string,
    isPublic?: boolean,
    override?: boolean
  ): Promise<any> => {
    return apiClient.uploadFile(file, parentId, isPublic, override);
  },
};

export const shareApi = {
  createShare: async (data: CreateShareRequest): Promise<Share> => {
    return apiClient.request<Share>('/shares', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getShares: async (page: number = 1, pageSize: number = 20): Promise<ShareListResponse> => {
    const response = await apiClient.request<any>(`/shares?page=${page}&page_size=${pageSize}`);
    return {
      shares: response.shares || [],
      total: response.total || 0,
      page: response.page || 1,
      page_size: response.size || 20,
    };
  },

  getShare: async (id: string): Promise<Share> => {
    return apiClient.request<Share>(`/shares/${id}`);
  },

  updateShare: async (id: string, data: UpdateShareRequest): Promise<Share> => {
    return apiClient.request<Share>(`/shares/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  deleteShare: async (id: string): Promise<void> => {
    await apiClient.request<void>(`/shares/${id}`, { method: 'DELETE' });
  },

  batchDeleteShares: async (ids: string[]): Promise<void> => {
    await apiClient.request<void>('/shares/batch-delete', {
      method: 'DELETE',
      body: JSON.stringify({ share_ids: ids }),
    });
  },

  getStats: async (): Promise<ShareStats> => {
    return apiClient.request<ShareStats>('/shares/stats');
  },

  accessShare: async (token: string, password?: string): Promise<Share> => {
    const body = password ? JSON.stringify({ password }) : undefined;
    const response = await apiClient.request<{ share: Share }>(`/s/${token}`, {
      method: 'POST',
      body,
    });
    return response.share;
  },

  downloadSharedFile: async (token: string): Promise<Blob> => {
    const url = `${API_BASE_URL}/s/${token}/download`;
    const headers = apiClient['getAuthHeaders']();
    const response = await fetch(url, { headers });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || `HTTP error! status: ${response.status}`);
    }

    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      const error = await response.json();
      throw new Error(error.error || 'Download failed');
    }

    return response.blob();
  },
};

export const recycleApi = {
  getRecycleFiles: async (page: number = 1, pageSize: number = 20): Promise<FileListResponse> => {
    const response = await apiClient.request<any>(`/recycle?page=${page}&page_size=${pageSize}`);
    return {
      files: response.files || [],
      total: response.total || 0,
      page: response.page || 1,
      page_size: response.size || 20,
      total_pages: Math.ceil((response.total || 0) / (response.size || 20)),
    };
  },

  restoreFile: async (id: string): Promise<File> => {
    return apiClient.request<File>(`/recycle/${id}/restore`, { method: 'POST' });
  },

  cleanupRecycle: async (days: number = 30): Promise<void> => {
    await apiClient.request<void>(`/recycle/cleanup?days=${days}`, { method: 'DELETE' });
  },
};

export const searchApi = {
  searchFiles: async (params: SearchParams): Promise<FileListResponse> => {
    const searchParams = new URLSearchParams();
    searchParams.append('q', params.q);
    if (params.search_in) searchParams.append('search_in', params.search_in);
    if (params.page) searchParams.append('page', params.page.toString());
    if (params.page_size) searchParams.append('page_size', params.page_size.toString());

    const response = await apiClient.request<any>(`/search?${searchParams.toString()}`);
    return {
      files: response.files || [],
      total: response.total || 0,
      page: response.page || 1,
      page_size: response.size || 20,
      total_pages: Math.ceil((response.total || 0) / (response.size || 20)),
    };
  },
};

export const statsApi = {
  getStorageStats: async (): Promise<StorageStats> => {
    return apiClient.request<StorageStats>('/stats/storage');
  },

  getFileStats: async (): Promise<FileStats> => {
    return apiClient.request<FileStats>('/stats/files');
  },
};

export const logApi = {
  getLogs: async (params: LogParams): Promise<OperationLogResponse> => {
    const searchParams = new URLSearchParams();
    if (params.page) searchParams.append('page', params.page.toString());
    if (params.page_size) searchParams.append('page_size', params.page_size.toString());
    if (params.user_id) searchParams.append('user_id', params.user_id);
    if (params.operation) searchParams.append('operation', params.operation);
    if (params.result) searchParams.append('result', params.result);
    if (params.start_date) searchParams.append('start_date', params.start_date);
    if (params.end_date) searchParams.append('end_date', params.end_date);

    const response = await apiClient.request<any>(`/logs?${searchParams.toString()}`);
    return {
      logs: response.logs || [],
      total: response.total || 0,
      page: response.page || 1,
      page_size: response.size || 50,
    };
  },

  getStats: async (): Promise<any> => {
    return apiClient.request<any>('/logs/stats');
  },

  cleanupLogs: async (): Promise<void> => {
    await apiClient.request<void>('/logs/cleanup', { method: 'DELETE' });
  },
};

export const adminApi = {
  getStats: async (): Promise<SystemStats> => {
    return apiClient.request<SystemStats>('/admin/stats');
  },

  getUsers: async (page: number = 1, pageSize: number = 20): Promise<AdminUserListResponse> => {
    const response = await apiClient.request<any>(`/admin/users?page=${page}&page_size=${pageSize}`);
    return {
      users: response.users || [],
      total: response.total || 0,
      page: response.page || 1,
      page_size: response.size || 20,
    };
  },

  getUser: async (id: string): Promise<AdminUser> => {
    return apiClient.request<AdminUser>(`/admin/users/${id}`);
  },

  updateUser: async (id: string, data: Partial<AdminUser>): Promise<AdminUser> => {
    return apiClient.request<AdminUser>(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  deleteUser: async (id: string): Promise<void> => {
    await apiClient.request<void>(`/admin/users/${id}`, { method: 'DELETE' });
  },

  activateUser: async (id: string): Promise<void> => {
    await apiClient.request<void>(`/admin/users/${id}/activate`, { method: 'POST' });
  },

  deactivateUser: async (id: string): Promise<void> => {
    await apiClient.request<void>(`/admin/users/${id}/deactivate`, { method: 'POST' });
  },
};
