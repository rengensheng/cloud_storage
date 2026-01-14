import { useState, useEffect } from 'react';
import { Copy } from 'lucide-react';
import { Dialog, Button, Input } from './ui';

interface RenameDialogProps {
  open: boolean;
  onClose: () => void;
  onRename: (newName: string) => Promise<void>;
  currentName: string;
  itemType: 'file' | 'directory';
}

export default function RenameDialog({
  open,
  onClose,
  onRename,
  currentName,
  itemType,
}: RenameDialogProps) {
  const [name, setName] = useState(currentName);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    setName(currentName);
    setError('');
  }, [currentName, open]);

  const handleSubmit = async () => {
    if (!name.trim()) {
      setError(`${itemType === 'directory' ? 'Folder' : 'File'} name is required`);
      return;
    }

    if (name === currentName) {
      onClose();
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      await onRename(name.trim());
      onClose();
    } catch (err: any) {
      setError(err.message || `Failed to rename ${itemType}`);
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
    <Dialog open={open} onClose={onClose} title={`Rename ${itemType === 'directory' ? 'Folder' : 'File'}`} size="small">
      <div className="space-y-4">
        <div className="flex items-center gap-3 p-4 bg-blue-50 rounded-lg">
          <Copy className="w-8 h-8 text-blue-600" />
          <div className="text-sm text-blue-900">
            <p className="font-medium">Rename {itemType === 'directory' ? 'folder' : 'file'}</p>
            <p className="text-blue-700">Current: {currentName}</p>
          </div>
        </div>

        <div>
          <label htmlFor="newName" className="block text-sm font-medium text-gray-700 mb-2">
            New Name
          </label>
          <Input
            id="newName"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Enter new name"
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
            Rename
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
