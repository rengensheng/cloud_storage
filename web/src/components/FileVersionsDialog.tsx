import { useState, useEffect } from 'react';
import { History, RotateCcw, Calendar, HardDrive, User as UserIcon } from 'lucide-react';
import { Dialog, Button, Card } from './ui';
import type { FileVersion } from '../types/api';

interface FileVersionsDialogProps {
  open: boolean;
  onClose: () => void;
  fileId: string;
  fileName: string;
  onRestore: (versionNumber: number) => Promise<void>;
  onLoadVersions: () => Promise<FileVersion[]>;
}

export default function FileVersionsDialog({
  open,
  onClose,
  fileId,
  fileName,
  onRestore,
  onLoadVersions,
}: FileVersionsDialogProps) {
  const [versions, setVersions] = useState<FileVersion[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isRestoring, setIsRestoring] = useState<number | null>(null);

  useEffect(() => {
    if (open) {
      loadVersions();
    }
  }, [open, fileId]);

  const loadVersions = async () => {
    setIsLoading(true);
    try {
      const data = await onLoadVersions();
      setVersions(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Failed to load versions:', error);
      setVersions([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRestore = async (versionNumber: number) => {
    setIsRestoring(versionNumber);
    try {
      await onRestore(versionNumber);
      await loadVersions();
    } catch (error) {
      console.error('Failed to restore version:', error);
    } finally {
      setIsRestoring(null);
    }
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
    return date.toLocaleString();
  };

  return (
    <Dialog open={open} onClose={onClose} title="File Versions" size="medium">
      <div className="space-y-4">
        <div className="flex items-center gap-2 text-sm text-gray-600">
          <History className="w-4 h-4" />
          <span>{fileName}</span>
        </div>

        {isLoading ? (
          <div className="text-center py-8 text-gray-600">Loading versions...</div>
        ) : versions.length === 0 ? (
          <div className="text-center py-8 text-gray-600">No versions available</div>
        ) : (
          <div className="space-y-3 max-h-96 overflow-y-auto">
            {versions.map((version, index) => (
              <Card key={version.id} padding="medium" className="hover:bg-gray-50">
                <div className="flex items-start justify-between gap-4">
                  <div className="flex items-start gap-3 flex-1">
                    <div className="flex-shrink-0 w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                      <span className="text-blue-600 font-bold text-sm">
                        v{version.version_number}
                      </span>
                    </div>

                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <Calendar className="w-4 h-4 text-gray-500" />
                        <span className="text-sm font-medium text-gray-900">
                          {formatDate(version.created_at)}
                        </span>
                      </div>

                      <div className="flex items-center gap-4 text-xs text-gray-600">
                        <div className="flex items-center gap-1">
                          <HardDrive className="w-3.5 h-3.5" />
                          <span>{formatFileSize(version.file_size)}</span>
                        </div>
                        <div className="flex items-center gap-1">
                          <UserIcon className="w-3.5 h-3.5" />
                          <span>{version.created_by}</span>
                        </div>
                      </div>

                      <div className="mt-1">
                        <span className="text-xs text-gray-500">
                          Hash: <span className="font-mono">{version.file_hash.slice(0, 12)}...</span>
                        </span>
                      </div>
                    </div>
                  </div>

                  {index > 0 && (
                    <Button
                      size="small"
                      variant="secondary"
                      leftIcon={<RotateCcw className="w-4 h-4" />}
                      onClick={() => handleRestore(version.version_number)}
                      loading={isRestoring === version.version_number}
                    >
                      Restore
                    </Button>
                  )}
                </div>
              </Card>
            ))}
          </div>
        )}

        <div className="flex justify-end pt-4 border-t border-gray-200">
          <Button variant="secondary" onClick={onClose}>
            Close
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
