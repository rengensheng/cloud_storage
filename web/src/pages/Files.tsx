import { useState, useEffect, useCallback } from 'react';
import {
  Search,
  Grid3x3,
  List,
  FolderPlus,
  Upload,
  RefreshCw,
  ArrowLeft,
  SortAsc,
  SortDesc,
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { Button, Input, Select, Card } from '../components/ui';
import FileItem from '../components/FileItem';
import CreateFolderDialog from '../components/CreateFolderDialog';
import RenameDialog from '../components/RenameDialog';
import ShareDialog from '../components/ShareDialog';
import FileVersionsDialog from '../components/FileVersionsDialog';
import PreviewDialog from '../components/PreviewDialog';
import { fileApi, shareApi } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import type { File as FileType, Share } from '../types/api';

type ViewMode = 'grid' | 'list';
type SortBy = 'name' | 'size' | 'created_at' | 'updated_at';

export default function Files() {
  const navigate = useNavigate();
  const { refreshUser } = useAuth();
  const [files, setFiles] = useState<FileType[]>([]);
  const [currentPath, setCurrentPath] = useState<FileType[]>([]);
  const [parentId, setParentId] = useState<string | undefined>(undefined);
  const [viewMode, setViewMode] = useState<ViewMode>('grid');
  const [sortBy, setSortBy] = useState<SortBy>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [selectedFiles, setSelectedFiles] = useState<Set<string>>(new Set());
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');

  // Dialog states
  const [showCreateFolder, setShowCreateFolder] = useState(false);
  const [showRename, setShowRename] = useState(false);
  const [showShare, setShowShare] = useState(false);
  const [showVersions, setShowVersions] = useState(false);
  const [showPreview, setShowPreview] = useState(false);
  const [selectedFile, setSelectedFile] = useState<FileType | null>(null);
  const [existingShare, setExistingShare] = useState<Share | null>(null);

  const loadFiles = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await fileApi.getFiles({
        parent_id: parentId,
        sort_by: sortBy,
        sort_order: sortOrder,
      });
      setFiles(response.files || []);
    } catch (error) {
      console.error('Failed to load files:', error);
    } finally {
      setIsLoading(false);
    }
  }, [parentId, sortBy, sortOrder]);

  useEffect(() => {
    loadFiles();
  }, [loadFiles]);

  const handleFileClick = (file: FileType) => {
    if (file.type === 'directory') {
      setCurrentPath([...currentPath, file]);
      setParentId(file.id);
      setSelectedFiles(new Set());
    } else {
      handlePreview(file);
    }
  };

  const handlePreview = (file: FileType) => {
    setSelectedFile(file);
    setShowPreview(true);
  };

  const navigateToFolder = (index: number) => {
    const newPath = currentPath.slice(0, index);
    setCurrentPath(newPath);
    setParentId(newPath.length > 0 ? newPath[newPath.length - 1].id : undefined);
    setSelectedFiles(new Set());
  };

  const handleNavigateUp = () => {
    if (currentPath.length > 0) {
      const newPath = currentPath.slice(0, -1);
      setCurrentPath(newPath);
      setParentId(newPath.length > 0 ? newPath[newPath.length - 1].id : undefined);
      setSelectedFiles(new Set());
    }
  };

  const navigateToUpload = () => {
    const url = parentId ? `/upload?parent_id=${parentId}` : '/upload';
    navigate(url);
  };

  const handleCreateFolder = async (name: string) => {
    await fileApi.createFile({
      name,
      type: 'directory',
      parent_id: parentId,
    });
    await loadFiles();
    await refreshUser();
  };

  const handleRename = async (newName: string) => {
    if (selectedFile) {
      await fileApi.updateFile(selectedFile.id, { name: newName });
      await loadFiles();
    }
  };

  const handleDownload = async (file: FileType) => {
    try {
      const blob = await fileApi.downloadFile(file.id);
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = file.name;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Download failed:', error);
    }
  };

  const handleDelete = async (file: FileType) => {
    await fileApi.deleteFile(file.id);
    await loadFiles();
    await refreshUser();
  };

  const handleCopy = async (file: FileType) => {
    await fileApi.copyFile(file.id, parentId || 'root');
    await loadFiles();
    await refreshUser();
  };

  const handleMove = async (file: FileType) => {
    await fileApi.moveFile(file.id, parentId || 'root');
    await loadFiles();
  };

  const handleCreateShare = async (data: any) => {
    return await shareApi.createShare(data);
  };

  const handleUpdateShare = async (id: string, data: any) => {
    return await shareApi.updateShare(id, data);
  };

  const handleDeleteShare = async (id: string) => {
    await shareApi.deleteShare(id);
    setExistingShare(null);
  };

  const handleRestoreVersion = async (versionNumber: number) => {
    if (selectedFile) {
      await fileApi.restoreVersion(selectedFile.id, versionNumber);
      await loadFiles();
    }
  };

  const handleLoadVersions = async () => {
    if (selectedFile) {
      return await fileApi.getVersions(selectedFile.id);
    }
    return [];
  };

  const openRenameDialog = (file: FileType) => {
    setSelectedFile(file);
    setShowRename(true);
  };

  const openShareDialog = async (file: FileType) => {
    setSelectedFile(file);
    try {
      const shares = await shareApi.getShares(1, 100);
      const fileShare = shares.shares.find((s) => s.file_id === file.id);
      setExistingShare(fileShare || null);
    } catch (error) {
      console.error('Failed to load shares:', error);
    }
    setShowShare(true);
  };

  const openVersionsDialog = (file: FileType) => {
    setSelectedFile(file);
    setShowVersions(true);
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

  const handleBulkDelete = async () => {
    for (const fileId of selectedFiles) {
      await fileApi.deleteFile(fileId);
    }
    setSelectedFiles(new Set());
    await loadFiles();
    await refreshUser();
  };

  const filteredFiles = (files || []).filter((file) =>
    file.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">My Files</h1>
          <p className="text-sm text-gray-600">
            {files.length} {files.length === 1 ? 'item' : 'items'}
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="secondary"
            size="small"
            leftIcon={<RefreshCw className="w-4 h-4" />}
            onClick={loadFiles}
            disabled={isLoading}
          >
            Refresh
          </Button>
          <Button
            variant="primary"
            size="small"
            leftIcon={<FolderPlus className="w-4 h-4" />}
            onClick={() => setShowCreateFolder(true)}
          >
            New Folder
          </Button>
        </div>
      </div>

      {/* Breadcrumb */}
      {currentPath.length > 0 && (
        <div className="flex items-center gap-2 bg-white rounded-lg p-3 border border-gray-200">
          <Button
            variant="tertiary"
            size="small"
            leftIcon={<ArrowLeft className="w-4 h-4" />}
            onClick={handleNavigateUp}
            disabled={currentPath.length === 0}
          >
            Back
          </Button>
          <div className="flex items-center gap-2 flex-1 overflow-x-auto">
            <button
              onClick={() => navigateToFolder(0)}
              className="text-sm text-blue-600 hover:text-blue-700 whitespace-nowrap"
            >
              Home
            </button>
            {currentPath.map((folder, index) => (
              <div key={folder.id} className="flex items-center gap-2">
                <span className="text-gray-400">/</span>
                <button
                  onClick={() => navigateToFolder(index + 1)}
                  className="text-sm text-blue-600 hover:text-blue-700 whitespace-nowrap"
                >
                  {folder.name}
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Toolbar */}
      <Card padding="medium">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <Input
                type="text"
                placeholder="Search files..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Select
              value={sortBy}
              onChange={(value) => setSortBy(value as SortBy)}
              options={[
                { value: 'name', label: 'Name' },
                { value: 'size', label: 'Size' },
                { value: 'created_at', label: 'Created' },
                { value: 'updated_at', label: 'Modified' },
              ]}
            />

            <Button
              variant="secondary"
              size="medium"
              onClick={() => setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')}
            >
              {sortOrder === 'asc' ? <SortAsc className="w-5 h-5" /> : <SortDesc className="w-5 h-5" />}
            </Button>

            <div className="flex border border-gray-300 rounded-lg overflow-hidden">
              <Button
                variant={viewMode === 'grid' ? 'primary' : 'tertiary'}
                size="medium"
                onClick={() => setViewMode('grid')}
              >
                <Grid3x3 className="w-5 h-5" />
              </Button>
              <Button
                variant={viewMode === 'list' ? 'primary' : 'tertiary'}
                size="medium"
                onClick={() => setViewMode('list')}
              >
                <List className="w-5 h-5" />
              </Button>
            </div>
          </div>
        </div>

        {selectedFiles.size > 0 && (
          <div className="mt-4 pt-4 border-t border-gray-200 flex items-center justify-between">
            <span className="text-sm text-gray-600">
              {selectedFiles.size} {selectedFiles.size === 1 ? 'item' : 'items'} selected
            </span>
            <Button
              variant="danger"
              size="small"
              onClick={handleBulkDelete}
            >
              Delete Selected
            </Button>
          </div>
        )}
      </Card>

      {/* Files */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-600">Loading files...</div>
      ) : filteredFiles.length === 0 ? (
        <div className="text-center py-12">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-gray-100 rounded-full mb-4">
            <FolderPlus className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No files yet</h3>
          <p className="text-gray-600 mb-4">
            {searchQuery ? 'Try a different search term' : 'Upload files or create a folder to get started'}
          </p>
          {!searchQuery && (
            <Button
              variant="primary"
              leftIcon={<Upload className="w-4 h-4" />}
              onClick={navigateToUpload}
            >
              Upload Files
            </Button>
          )}
        </div>
      ) : (
        <div
          className={
            viewMode === 'grid'
              ? 'grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4'
              : 'space-y-2'
          }
        >
          {filteredFiles.map((file) => (
            <FileItem
              key={file.id}
              file={file}
              view={viewMode}
              isSelected={selectedFiles.has(file.id)}
              onSelect={() => toggleFileSelection(file.id)}
              onOpen={() => handleFileClick(file)}
              onDownload={() => handleDownload(file)}
              onShare={() => openShareDialog(file)}
              onDelete={() => handleDelete(file)}
              onRename={() => openRenameDialog(file)}
              onCopy={() => handleCopy(file)}
              onMove={() => handleMove(file)}
              onVersions={() => openVersionsDialog(file)}
            />
          ))}
        </div>
      )}

      {/* Dialogs */}
      <CreateFolderDialog
        open={showCreateFolder}
        onClose={() => setShowCreateFolder(false)}
        onCreate={handleCreateFolder}
        parentId={parentId}
      />

      {selectedFile && (
        <>
          <RenameDialog
            open={showRename}
            onClose={() => setShowRename(false)}
            onRename={handleRename}
            currentName={selectedFile.name}
            itemType={selectedFile.type}
          />

          <ShareDialog
            open={showShare}
            onClose={() => setShowShare(false)}
            fileId={selectedFile.id}
            fileName={selectedFile.name}
            existingShare={existingShare}
            onCreate={handleCreateShare}
            onUpdate={handleUpdateShare}
            onDelete={handleDeleteShare}
          />

          <FileVersionsDialog
            open={showVersions}
            onClose={() => setShowVersions(false)}
            fileId={selectedFile.id}
            fileName={selectedFile.name}
            onRestore={handleRestoreVersion}
            onLoadVersions={handleLoadVersions}
          />

          <PreviewDialog
            open={showPreview}
            onClose={() => setShowPreview(false)}
            file={selectedFile}
          />
        </>
      )}
    </div>
  );
}
