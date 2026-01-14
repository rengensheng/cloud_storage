import { useState, useCallback, useRef, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import {
  Upload as UploadIcon,
  X,
  File,
  Check,
  AlertCircle,
  FolderOpen,
} from 'lucide-react';
import { Button, Card, Dialog } from '../components/ui';
import { fileApi, statsApi } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

interface UploadFile {
  id: string;
  file: File;
  progress: number;
  status: 'pending' | 'uploading' | 'success' | 'error';
  error?: string;
}

export default function Upload() {
  const navigate = useNavigate();
  const { user, refreshUser } = useAuth();
  const [searchParams] = useSearchParams();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [uploadFiles, setUploadFiles] = useState<UploadFile[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const [storageUsed, setStorageUsed] = useState<{ used: number; quota: number } | null>(null);

  // Get parent_id from URL query params
  const parentId = searchParams.get('parent_id');

  const loadStorageStats = async () => {
    try {
      const stats = await statsApi.getStorageStats();
      setStorageUsed({ used: stats.used, quota: stats.quota });
    } catch (error) {
      console.error('Failed to load storage stats:', error);
    }
  };

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const droppedFiles = Array.from(e.dataTransfer.files);
    addFiles(droppedFiles);
  }, []);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const selectedFiles = Array.from(e.target.files);
      addFiles(selectedFiles);
    }
  };

  const addFiles = (files: File[]) => {
    const newUploadFiles: UploadFile[] = files.map((file) => ({
      id: Math.random().toString(36).substring(7),
      file,
      progress: 0,
      status: 'pending',
    }));

    setUploadFiles((prev) => [...prev, ...newUploadFiles]);
  };

  const removeFile = (id: string) => {
    setUploadFiles((prev) => prev.filter((f) => f.id !== id));
  };

  const uploadFile = async (uploadFile: UploadFile) => {
    setUploadFiles((prev) =>
      prev.map((f) =>
        f.id === uploadFile.id ? { ...f, status: 'uploading' } : f
      )
    );

    try {
      const onProgress = (progress: number) => {
        setUploadFiles((prev) =>
          prev.map((f) =>
            f.id === uploadFile.id ? { ...f, progress } : f
          )
        );
      };

      await fileApi.uploadFile(uploadFile.file, parentId || undefined, false, false);

      setUploadFiles((prev) =>
        prev.map((f) =>
          f.id === uploadFile.id ? { ...f, progress: 100, status: 'success' } : f
        )
      );

      await refreshUser();
      await loadStorageStats();
    } catch (error: any) {
      setUploadFiles((prev) =>
        prev.map((f) =>
          f.id === uploadFile.id
            ? { ...f, status: 'error', error: error.message || 'Upload failed' }
            : f
        )
      );
    }
  };

  const uploadAllFiles = async () => {
    const pendingFiles = uploadFiles.filter((f) => f.status === 'pending');

    for (const file of pendingFiles) {
      await uploadFile(file);
    }
  };

  const clearCompleted = () => {
    setUploadFiles((prev) => prev.filter((f) => f.status !== 'success'));
  };

  const clearAll = () => {
    setUploadFiles([]);
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getFileIcon = (fileName: string) => {
    const ext = fileName.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf':
        return <File className="w-8 h-8 text-red-500" />;
      case 'doc':
      case 'docx':
        return <File className="w-8 h-8 text-blue-600" />;
      case 'xls':
      case 'xlsx':
        return <File className="w-8 h-8 text-green-600" />;
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
        return <File className="w-8 h-8 text-purple-500" />;
      case 'mp4':
      case 'mov':
      case 'avi':
        return <File className="w-8 h-8 text-pink-500" />;
      default:
        return <File className="w-8 h-8 text-gray-500" />;
    }
  };

  const totalSize = uploadFiles.reduce((acc, f) => acc + f.file.size, 0);
  const uploadedSize = uploadFiles
    .filter((f) => f.status === 'success')
    .reduce((acc, f) => acc + f.file.size, 0);
  const progress = uploadFiles.length > 0 ? (uploadedSize / totalSize) * 100 : 0;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Upload Files</h1>
        <p className="text-sm text-gray-600">
          {storageUsed
            ? `Storage: ${formatFileSize(storageUsed.used)} / ${formatFileSize(storageUsed.quota)} used`
            : 'Loading storage info...'}
        </p>
      </div>

      {/* Upload Destination */}
      {parentId ? (
        <Card padding="medium" className="bg-blue-50 border-blue-200">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <FolderOpen className="w-5 h-5 text-blue-600" />
              <div>
                <p className="text-sm font-medium text-blue-900">
                  Uploading to a subfolder
                </p>
                <p className="text-xs text-blue-700">
                  Files will be uploaded to the selected folder
                </p>
              </div>
            </div>
            <Button
              variant="secondary"
              size="small"
              onClick={() => navigate('/')}
            >
              Change Folder
            </Button>
          </div>
        </Card>
      ) : (
        <Card padding="medium" className="bg-gray-50 border-gray-200">
          <div className="flex items-center gap-3">
            <FolderOpen className="w-5 h-5 text-gray-600" />
            <div>
              <p className="text-sm font-medium text-gray-900">
                Uploading to root directory
              </p>
              <p className="text-xs text-gray-600">
                Files will be uploaded to your main folder
              </p>
            </div>
          </div>
        </Card>
      )}

      {/* Upload Area */}
      <Card
        padding="large"
        className={`border-2 border-dashed transition-colors ${
          isDragging
            ? 'border-blue-500 bg-blue-50'
            : 'border-gray-300 hover:border-blue-400'
        }`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        <div className="text-center py-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-4">
            <UploadIcon className="w-8 h-8 text-blue-600" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            Drop files here or click to browse
          </h3>
          <p className="text-sm text-gray-600 mb-4">
            Upload any file up to 100MB
          </p>
          <input
            type="file"
            ref={fileInputRef}
            multiple
            onChange={handleFileSelect}
            className="hidden"
          />
          <Button
            variant="primary"
            onClick={() => fileInputRef.current?.click()}
          >
            Browse Files
          </Button>
        </div>
      </Card>

      {/* File List */}
      {uploadFiles.length > 0 && (
        <Card padding="medium">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-medium text-gray-900">
              {uploadFiles.length} {uploadFiles.length === 1 ? 'file' : 'files'}
            </h3>
            <div className="flex gap-2">
              <Button variant="secondary" size="small" onClick={clearCompleted}>
                Clear Completed
              </Button>
              <Button variant="secondary" size="small" onClick={clearAll}>
                Clear All
              </Button>
            </div>
          </div>

          {/* Overall Progress */}
          {uploadFiles.some((f) => f.status === 'uploading') && (
            <div className="mb-4">
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-600">Uploading...</span>
                <span className="font-medium">{Math.round(progress)}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                <div
                  className="bg-gradient-to-r from-blue-500 to-cyan-500 h-full rounded-full transition-all duration-300"
                  style={{ width: `${progress}%` }}
                />
              </div>
            </div>
          )}

          <div className="space-y-3">
            {uploadFiles.map((uploadFile) => (
              <div
                key={uploadFile.id}
                className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg"
              >
                <div className="flex-shrink-0">{getFileIcon(uploadFile.file.name)}</div>

                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-gray-900 truncate">
                    {uploadFile.file.name}
                  </p>
                  <div className="flex items-center gap-2 mt-1">
                    <span className="text-xs text-gray-600">
                      {formatFileSize(uploadFile.file.size)}
                    </span>
                    {uploadFile.status === 'uploading' && (
                      <div className="flex-1 max-w-32">
                        <div className="w-full bg-gray-200 rounded-full h-1.5 overflow-hidden">
                          <div
                            className="bg-blue-500 h-full rounded-full transition-all duration-300"
                            style={{ width: `${uploadFile.progress}%` }}
                          />
                        </div>
                      </div>
                    )}
                  </div>
                  {uploadFile.status === 'error' && (
                    <p className="text-xs text-red-600 mt-1">{uploadFile.error}</p>
                  )}
                </div>

                <div className="flex items-center gap-2">
                  {uploadFile.status === 'pending' && (
                    <Button
                      variant="primary"
                      size="small"
                      onClick={() => uploadFile(uploadFile)}
                    >
                      Upload
                    </Button>
                  )}
                  {uploadFile.status === 'success' && (
                    <Check className="w-5 h-5 text-green-500" />
                  )}
                  {uploadFile.status === 'error' && (
                    <AlertCircle className="w-5 h-5 text-red-500" />
                  )}
                  <button
                    onClick={() => removeFile(uploadFile.id)}
                    className="p-1 hover:bg-gray-200 rounded"
                  >
                    <X className="w-4 h-4 text-gray-500" />
                  </button>
                </div>
              </div>
            ))}
          </div>

          <div className="flex justify-end gap-3 mt-4 pt-4 border-t border-gray-200">
            <Button
              variant="primary"
              onClick={uploadAllFiles}
              disabled={!uploadFiles.some((f) => f.status === 'pending')}
            >
              Upload All ({uploadFiles.filter((f) => f.status === 'pending').length})
            </Button>
            <Button
              variant="secondary"
              onClick={() => navigate('/')}
              disabled={uploadFiles.some((f) => f.status === 'uploading')}
            >
              View Files
            </Button>
          </div>
        </Card>
      )}

      {/* Quick Links */}
      <Card padding="medium">
        <h3 className="font-medium text-gray-900 mb-3">Quick Actions</h3>
        <div className="flex gap-3">
          <Button
            variant="secondary"
            leftIcon={<FolderOpen className="w-4 h-4" />}
            onClick={() => navigate('/')}
          >
            My Files
          </Button>
          <Button
            variant="secondary"
            leftIcon={<UploadIcon className="w-4 h-4" />}
            onClick={() => navigate('/recycle')}
          >
            Recycle Bin
          </Button>
        </div>
      </Card>
    </div>
  );
}
