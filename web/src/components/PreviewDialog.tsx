import { useState, useEffect } from 'react';
import { File as FileIcon, Download } from 'lucide-react';
import { Dialog, Button, Card } from './ui';
import type { File as FileType } from '../types/api';
import { fileApi } from '../services/api';

interface PreviewDialogProps {
  open: boolean;
  onClose: () => void;
  file: FileType;
}

export default function PreviewDialog({ open, onClose, file }: PreviewDialogProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [previewUrl, setPreviewUrl] = useState<string>('');
  const [error, setError] = useState<string>('');

  useEffect(() => {
    if (open && file) {
      loadPreview();
    }

    return () => {
      if (previewUrl) {
        URL.revokeObjectURL(previewUrl);
        setPreviewUrl('');
      }
    };
  }, [open, file.id]);

  const loadPreview = async () => {
    setError('');
    setIsLoading(true);
    try {
      const blob = await fileApi.downloadFile(file.id);
      const url = URL.createObjectURL(blob);
      setPreviewUrl(url);
    } catch (err: any) {
      console.error('Failed to load preview:', err);
      setError(err.message || 'Failed to load preview');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDownload = async () => {
    setIsLoading(true);
    try {
      const blob = await fileApi.downloadFile(file.id);
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = file.name;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      setTimeout(() => URL.revokeObjectURL(url), 100);
    } catch (error) {
      console.error('Download failed:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const renderPreview = () => {
    const ext = file.name.split('.').pop()?.toLowerCase();

    if (error) {
      return (
        <div className="flex flex-col items-center justify-center py-16">
          <FileIcon className="w-16 h-16 text-red-400 mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">Preview Failed</h3>
          <p className="text-sm text-gray-600 text-center max-w-md">{error}</p>
        </div>
      );
    }

    if (isLoading) {
      return (
        <div className="flex flex-col items-center justify-center py-16">
          <div className="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mb-4"></div>
          <p className="text-sm text-gray-600">Loading preview...</p>
        </div>
      );
    }

    if (!previewUrl) {
      return (
        <div className="flex flex-col items-center justify-center py-16">
          <p className="text-sm text-gray-600">No preview available</p>
        </div>
      );
    }

    // Image files
    if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(ext || '')) {
      return (
        <img
          src={previewUrl}
          alt={file.name}
          className="max-w-full max-h-[70vh] object-contain mx-auto rounded-lg"
        />
      );
    }

    // Video files
    if (['mp4', 'webm', 'ogg', 'mov', 'avi'].includes(ext || '')) {
      return (
        <video
          controls
          className="max-w-full max-h-[70vh] mx-auto rounded-lg"
          src={previewUrl}
        >
          Your browser does not support the video tag.
        </video>
      );
    }

    // PDF files
    if (ext === 'pdf') {
      return (
        <iframe
          src={previewUrl}
          className="w-full h-[70vh] rounded-lg"
          title={file.name}
        />
      );
    }

    // Audio files
    if (['mp3', 'wav', 'ogg', 'm4a'].includes(ext || '')) {
      return (
        <div className="flex flex-col items-center justify-center py-12">
          <FileIcon className="w-16 h-16 text-gray-400 mb-4" />
          <audio
            controls
            className="w-full max-w-md"
            src={previewUrl}
          >
            Your browser does not support the audio tag.
          </audio>
        </div>
      );
    }

    // Unsupported file types
    return (
      <div className="flex flex-col items-center justify-center py-16">
        <FileIcon className="w-20 h-20 text-gray-300 mb-4" />
        <h3 className="text-lg font-medium text-gray-900 mb-2">Preview Not Available</h3>
        <p className="text-sm text-gray-600 text-center max-w-md">
          This file type ({ext?.toUpperCase()}) cannot be previewed directly.
          Please download the file to view its contents.
        </p>
      </div>
    );
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <Dialog open={open} onClose={onClose} title={file.name} size="large">
      <div className="space-y-4">
        {/* File Info */}
        <Card padding="small" className="bg-gray-50">
          <div className="flex items-center justify-between text-sm">
            <div className="flex items-center gap-4">
              <span className="text-gray-600">Size: <span className="font-medium text-gray-900">{formatFileSize(file.size)}</span></span>
              <span className="text-gray-600">Type: <span className="font-medium text-gray-900">{file.mime_type || 'Unknown'}</span></span>
            </div>
            <span className="text-gray-600">Modified: <span className="font-medium text-gray-900">{new Date(file.updated_at).toLocaleDateString()}</span></span>
          </div>
        </Card>

        {/* Preview Area */}
        <div className="bg-gray-100 rounded-lg p-4 flex items-center justify-center min-h-[400px]">
          {renderPreview()}
        </div>

        {/* Actions */}
        <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
          <Button variant="secondary" onClick={onClose}>
            Close
          </Button>
          <Button
            variant="primary"
            leftIcon={<Download className="w-4 h-4" />}
            onClick={handleDownload}
            loading={isLoading}
          >
            Download File
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
