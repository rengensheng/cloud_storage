import { useState } from 'react';
import { FolderPlus } from 'lucide-react';
import { Dialog, Button, Input } from './ui';

interface CreateFolderDialogProps {
  open: boolean;
  onClose: () => void;
  onCreate: (name: string) => Promise<void>;
  parentId?: string;
}

export default function CreateFolderDialog({
  open,
  onClose,
  onCreate,
  parentId,
}: CreateFolderDialogProps) {
  const [name, setName] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    if (!name.trim()) {
      setError('Folder name is required');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      await onCreate(name.trim());
      setName('');
      onClose();
    } catch (err: any) {
      setError(err.message || 'Failed to create folder');
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !isLoading) {
      handleSubmit();
    }
  };

  return (
    <Dialog open={open} onClose={onClose} title="Create Folder" size="small">
      <div className="space-y-4">
        <div className="flex items-center gap-3 p-4 bg-blue-50 rounded-lg">
          <FolderPlus className="w-8 h-8 text-blue-600" />
          <div className="text-sm text-blue-900">
            <p className="font-medium">Create a new folder</p>
            <p className="text-blue-700">Enter the name for your new folder</p>
          </div>
        </div>

        <div>
          <label htmlFor="folderName" className="block text-sm font-medium text-gray-700 mb-2">
            Folder Name
          </label>
          <Input
            id="folderName"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="My Folder"
            status={error ? 'error' : 'default'}
            autoFocus
          />
          {error && (
            <p className="text-sm text-red-600 mt-1">{error}</p>
          )}
        </div>

        <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
          <Button variant="secondary" onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button variant="primary" onClick={handleSubmit} loading={isLoading}>
            Create Folder
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
