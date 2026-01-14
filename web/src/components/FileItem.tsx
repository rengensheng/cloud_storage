import { useState } from 'react';
import {
  File,
  Folder,
  MoreVertical,
  Download,
  Share2,
  Trash2,
  Copy,
  Scissors,
  Eye,
  History,
} from 'lucide-react';
import { Menu } from './ui/Menu';
import type { File as FileType } from '../types/api';

interface FileItemProps {
  file: FileType;
  view: 'grid' | 'list';
  isSelected: boolean;
  onSelect: () => void;
  onOpen: () => void;
  onDownload: () => void;
  onShare: () => void;
  onDelete: () => void;
  onRename: () => void;
  onCopy: () => void;
  onMove: () => void;
  onVersions: () => void;
}

export default function FileItem({
  file,
  view,
  isSelected,
  onSelect,
  onOpen,
  onDownload,
  onShare,
  onDelete,
  onRename,
  onCopy,
  onMove,
  onVersions,
}: FileItemProps) {
  const [showMenu, setShowMenu] = useState(false);

  const getFileIcon = () => {
    if (file.type === 'directory') {
      return <Folder className="w-6 h-6 text-blue-500" />;
    }

    const ext = file.name.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf':
        return <File className="w-6 h-6 text-red-500" />;
      case 'doc':
      case 'docx':
        return <File className="w-6 h-6 text-blue-600" />;
      case 'xls':
      case 'xlsx':
        return <File className="w-6 h-6 text-green-600" />;
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
        return <File className="w-6 h-6 text-purple-500" />;
      case 'mp4':
      case 'mov':
      case 'avi':
        return <File className="w-6 h-6 text-pink-500" />;
      case 'zip':
      case 'rar':
        return <File className="w-6 h-6 text-yellow-600" />;
      default:
        return <File className="w-6 h-6 text-gray-500" />;
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
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 0) {
      const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
      if (diffHours === 0) {
        const diffMinutes = Math.floor(diffTime / (1000 * 60));
        return diffMinutes <= 1 ? 'Just now' : `${diffMinutes}m ago`;
      }
      return `${diffHours}h ago`;
    } else if (diffDays === 1) {
      return 'Yesterday';
    } else if (diffDays < 7) {
      return `${diffDays} days ago`;
    } else {
      return date.toLocaleDateString();
    }
  };

  const menuItems = [
    { key: 'open', label: 'Open', icon: <Eye className="w-4 h-4" />, onClick: onOpen },
    { key: 'download', label: 'Download', icon: <Download className="w-4 h-4" />, onClick: onDownload },
    { key: 'share', label: 'Share', icon: <Share2 className="w-4 h-4" />, onClick: onShare },
    { key: 'rename', label: 'Rename', icon: <Copy className="w-4 h-4" />, onClick: onRename },
    { key: 'copy', label: 'Copy', icon: <Copy className="w-4 h-4" />, onClick: onCopy },
    { key: 'move', label: 'Move', icon: <Scissors className="w-4 h-4" />, onClick: onMove },
    ...(file.type === 'file' ? [{ key: 'versions', label: 'Versions', icon: <History className="w-4 h-4" />, onClick: onVersions }] : []),
    { key: 'delete', label: 'Delete', icon: <Trash2 className="w-4 h-4" />, onClick: onDelete, danger: true },
  ];

  if (view === 'grid') {
    return (
      <div
        className={`relative bg-white rounded-xl border-2 transition-all duration-200 cursor-pointer hover:shadow-lg group ${
          isSelected ? 'border-blue-500 shadow-lg' : 'border-gray-200 hover:border-blue-300'
        }`}
        onClick={onOpen}
        onContextMenu={(e) => {
          e.preventDefault();
          setShowMenu(true);
        }}
      >
        <div
          className="absolute top-3 left-3 z-10"
          onClick={(e) => {
            e.stopPropagation();
            onSelect();
          }}
        >
          <input
            type="checkbox"
            checked={isSelected}
            onChange={(e) => e.stopPropagation()}
            className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
        </div>

        <div className="p-4">
          <div className="flex justify-center mb-3">
            {getFileIcon()}
          </div>

          <div className="text-center">
            <p className="font-medium text-gray-900 text-sm truncate" title={file.name}>
              {file.name}
            </p>
            <p className="text-xs text-gray-500 mt-1">
              {file.type === 'directory' ? 'Folder' : formatFileSize(file.size)}
            </p>
          </div>
        </div>

        <div className="absolute top-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity" onClick={(e) => e.stopPropagation()}>
          <Menu
            trigger={
              <button
                className="p-1.5 hover:bg-gray-100 rounded-lg"
              >
                <MoreVertical className="w-4 h-4 text-gray-600" />
              </button>
            }
            items={menuItems}
            position="bottom-end"
          />
        </div>
      </div>
    );
  }

  return (
    <div
      className={`flex items-center gap-3 p-3 rounded-lg border transition-all duration-200 cursor-pointer hover:bg-gray-50 group ${
        isSelected ? 'bg-blue-50 border-blue-500' : 'bg-white border-gray-200 hover:border-blue-300'
      }`}
      onClick={onOpen}
    >
      <input
        type="checkbox"
        checked={isSelected}
        onChange={(e) => {
          e.stopPropagation();
          onSelect();
        }}
        className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
      />

      <div className="flex-shrink-0">{getFileIcon()}</div>

      <div className="flex-1 min-w-0">
        <p className="font-medium text-gray-900 text-sm truncate" title={file.name}>
          {file.name}
        </p>
      </div>

      <div className="hidden md:block flex-shrink-0 w-24 text-sm text-gray-600">
        {file.type === 'directory' ? 'Folder' : formatFileSize(file.size)}
      </div>

      <div className="hidden lg:block flex-shrink-0 w-32 text-sm text-gray-600">
        {formatDate(file.updated_at)}
      </div>

      <div className="flex-shrink-0" onClick={(e) => e.stopPropagation()}>
        <Menu
          trigger={
            <button
              className="p-1.5 hover:bg-gray-100 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity"
            >
              <MoreVertical className="w-4 h-4 text-gray-600" />
            </button>
          }
          items={menuItems}
          position="bottom-end"
        />
      </div>
    </div>
  );
}
