import { useState, useEffect } from 'react';
import {
  Trash2,
  RefreshCw,
  RotateCcw,
  Trash,
  AlertOctagon,
} from 'lucide-react';
import { Button, Card, Dialog, Select } from '../components/ui';
import { recycleApi } from '../services/api';
import type { File as FileType } from '../types/api';

export default function Recycle() {
  const [files, setFiles] = useState<FileType[]>([]);
  const [selectedFiles, setSelectedFiles] = useState<Set<string>>(new Set());
  const [isLoading, setIsLoading] = useState(true);
  const [showRestoreDialog, setShowRestoreDialog] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showCleanupDialog, setShowCleanupDialog] = useState(false);
  const [cleanupDays, setCleanupDays] = useState('30');

  const loadRecycleFiles = async () => {
    setIsLoading(true);
    try {
      const response = await recycleApi.getRecycleFiles(1, 100);
      setFiles(response.files);
    } catch (error) {
      console.error('Failed to load recycle files:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadRecycleFiles();
  }, []);

  const handleRestore = async (fileId?: string) => {
    try {
      const filesToRestore = fileId ? [fileId] : Array.from(selectedFiles);

      for (const id of filesToRestore) {
        await recycleApi.restoreFile(id);
      }

      setSelectedFiles(new Set());
      setShowRestoreDialog(false);
      await loadRecycleFiles();
    } catch (error) {
      console.error('Failed to restore files:', error);
    }
  };

  const handlePermanentDelete = async (fileId?: string) => {
    try {
      const filesToDelete = fileId ? [fileId] : Array.from(selectedFiles);

      for (const id of filesToDelete) {
        await recycleApi.restoreFile(id);
      }

      setSelectedFiles(new Set());
      setShowDeleteDialog(false);
      await loadRecycleFiles();
    } catch (error) {
      console.error('Failed to delete files:', error);
    }
  };

  const handleCleanup = async () => {
    try {
      await recycleApi.cleanupRecycle(parseInt(cleanupDays));
      setShowCleanupDialog(false);
      await loadRecycleFiles();
    } catch (error) {
      console.error('Failed to cleanup recycle bin:', error);
    }
  };

  const toggleFileSelection = (fileId: string) => {
    const newSelection = new Set(selectedFiles);
    if (newSelection.has(fileId)) {
      newSelection.delete(fileId);
    } else {
      newSelection.add(fileId);
    }
    setSelectedFiles(newSelection);
  };

  const selectAll = () => {
    setSelectedFiles(new Set(files.map((f) => f.id)));
  };

  const clearSelection = () => {
    setSelectedFiles(new Set());
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 0) {
      return 'Today';
    } else if (diffDays === 1) {
      return 'Yesterday';
    } else if (diffDays < 7) {
      return `${diffDays} days ago`;
    } else {
      return date.toLocaleDateString();
    }
  };

  const getFileIcon = (file: FileType) => {
    if (file.type === 'directory') {
      return <Trash2 className="w-6 h-6 text-blue-500" />;
    }

    const ext = file.name.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf':
        return <Trash2 className="w-6 h-6 text-red-500" />;
      case 'doc':
      case 'docx':
        return <Trash2 className="w-6 h-6 text-blue-600" />;
      case 'xls':
      case 'xlsx':
        return <Trash2 className="w-6 h-6 text-green-600" />;
      default:
        return <Trash2 className="w-6 h-6 text-gray-500" />;
    }
  };

  return (
    <div className="max-w-5xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
            <Trash2 className="w-6 h-6 text-red-500" />
            Recycle Bin
          </h1>
          <p className="text-sm text-gray-600">
            {files.length} {files.length === 1 ? 'item' : 'items'}
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="secondary"
            size="small"
            leftIcon={<RefreshCw className="w-4 h-4" />}
            onClick={loadRecycleFiles}
            disabled={isLoading}
          >
            Refresh
          </Button>
          <Button
            variant="secondary"
            size="small"
            leftIcon={<Trash className="w-4 h-4" />}
            onClick={() => setShowCleanupDialog(true)}
            disabled={files.length === 0}
          >
            Empty Old
          </Button>
        </div>
      </div>

      {/* Info Banner */}
      <Card className="bg-amber-50 border-amber-200" padding="medium">
        <div className="flex items-start gap-3">
          <AlertOctagon className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
          <div className="text-sm text-amber-900">
            <p className="font-medium">About Recycle Bin</p>
            <p className="text-amber-700 mt-1">
              Items in the recycle bin will be permanently deleted after 30 days.
              You can restore items or delete them permanently at any time.
            </p>
          </div>
        </div>
      </Card>

      {/* Toolbar */}
      {files.length > 0 && (
        <Card padding="medium">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Button
                variant="tertiary"
                size="small"
                onClick={selectedFiles.size === files.length ? clearSelection : selectAll}
              >
                {selectedFiles.size === files.length ? 'Deselect All' : 'Select All'}
              </Button>
              {selectedFiles.size > 0 && (
                <span className="text-sm text-gray-600">
                  {selectedFiles.size} selected
                </span>
              )}
            </div>

            <div className="flex gap-2">
              <Button
                variant="secondary"
                size="small"
                leftIcon={<RotateCcw className="w-4 h-4" />}
                onClick={() => setShowRestoreDialog(true)}
                disabled={selectedFiles.size === 0}
              >
                Restore
              </Button>
              <Button
                variant="danger"
                size="small"
                leftIcon={<Trash className="w-4 h-4" />}
                onClick={() => setShowDeleteDialog(true)}
                disabled={selectedFiles.size === 0}
              >
                Delete Forever
              </Button>
            </div>
          </div>
        </Card>
      )}

      {/* File List */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-600">Loading...</div>
      ) : files.length === 0 ? (
        <div className="text-center py-12">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-gray-100 rounded-full mb-4">
            <Trash2 className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">Recycle bin is empty</h3>
          <p className="text-gray-600">
            Items you delete will appear here
          </p>
        </div>
      ) : (
        <Card padding="medium">
          <div className="space-y-2">
            {files.map((file) => (
              <div
                key={file.id}
                className="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <input
                  type="checkbox"
                  checked={selectedFiles.has(file.id)}
                  onChange={() => toggleFileSelection(file.id)}
                  className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />

                <div className="flex-shrink-0">{getFileIcon(file)}</div>

                <div className="flex-1 min-w-0">
                  <p className="font-medium text-gray-900 truncate">{file.name}</p>
                  <div className="flex items-center gap-3 mt-1 text-xs text-gray-600">
                    <span>{file.type === 'directory' ? 'Folder' : formatFileSize(file.size)}</span>
                    <span>Deleted {formatDate(file.deleted_at!)}</span>
                  </div>
                </div>

                <div className="flex gap-2">
                  <Button
                    variant="secondary"
                    size="small"
                    leftIcon={<RotateCcw className="w-4 h-4" />}
                    onClick={() => handleRestore(file.id)}
                  >
                    Restore
                  </Button>
                  <Button
                    variant="danger"
                    size="small"
                    leftIcon={<Trash className="w-4 h-4" />}
                    onClick={() => handlePermanentDelete(file.id)}
                  >
                    Delete
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}

      {/* Restore Dialog */}
      <Dialog
        open={showRestoreDialog}
        onClose={() => setShowRestoreDialog(false)}
        title="Restore Items"
        size="small"
      >
        <div className="space-y-4">
          <p className="text-gray-600">
            Are you sure you want to restore {selectedFiles.size} item{selectedFiles.size > 1 ? 's' : ''}?
          </p>
          <div className="flex justify-end gap-3">
            <Button variant="secondary" onClick={() => setShowRestoreDialog(false)}>
              Cancel
            </Button>
            <Button variant="primary" onClick={() => handleRestore()}>
              Restore
            </Button>
          </div>
        </div>
      </Dialog>

      {/* Delete Forever Dialog */}
      <Dialog
        open={showDeleteDialog}
        onClose={() => setShowDeleteDialog(false)}
        title="Delete Forever"
        size="small"
      >
        <div className="space-y-4">
          <div className="flex items-start gap-3 p-3 bg-red-50 rounded-lg">
            <AlertOctagon className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
            <div className="text-sm text-red-900">
              <p className="font-medium">Warning: This action cannot be undone</p>
              <p className="text-red-700 mt-1">
                Are you sure you want to permanently delete {selectedFiles.size} item{selectedFiles.size > 1 ? 's' : ''}?
              </p>
            </div>
          </div>
          <div className="flex justify-end gap-3">
            <Button variant="secondary" onClick={() => setShowDeleteDialog(false)}>
              Cancel
            </Button>
            <Button variant="danger" onClick={() => handlePermanentDelete()}>
              Delete Forever
            </Button>
          </div>
        </div>
      </Dialog>

      {/* Cleanup Dialog */}
      <Dialog
        open={showCleanupDialog}
        onClose={() => setShowCleanupDialog(false)}
        title="Empty Old Items"
        size="small"
      >
        <div className="space-y-4">
          <p className="text-gray-600">
            Permanently delete items that were deleted more than a certain number of days ago.
          </p>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Delete items older than
            </label>
            <Select
              value={cleanupDays}
              onChange={(value) => setCleanupDays(value)}
              options={[
                { value: '7', label: '7 days' },
                { value: '14', label: '14 days' },
                { value: '30', label: '30 days' },
                { value: '60', label: '60 days' },
                { value: '90', label: '90 days' },
              ]}
            />
          </div>

          <div className="flex justify-end gap-3">
            <Button variant="secondary" onClick={() => setShowCleanupDialog(false)}>
              Cancel
            </Button>
            <Button variant="danger" onClick={handleCleanup}>
              Empty
            </Button>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
