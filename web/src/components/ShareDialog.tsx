import { useState, useEffect } from 'react';
import { Link2, Copy, Check, Lock, Unlock } from 'lucide-react';
import { Dialog, Button, Input, Select, Card } from './ui';
import type { Share, CreateShareRequest, UpdateShareRequest } from '../types/api';

interface ShareDialogProps {
  open: boolean;
  onClose: () => void;
  fileId: string;
  fileName: string;
  existingShare?: Share | null;
  onCreate: (data: CreateShareRequest) => Promise<Share>;
  onUpdate: (id: string, data: UpdateShareRequest) => Promise<Share>;
  onDelete: (id: string) => Promise<void>;
}

export default function ShareDialog({
  open,
  onClose,
  fileId,
  fileName,
  existingShare,
  onCreate,
  onUpdate,
  onDelete,
}: ShareDialogProps) {
  const [password, setPassword] = useState('');
  const [accessType, setAccessType] = useState<'view' | 'download' | 'edit'>('view');
  const [expiresIn, setExpiresIn] = useState<string>('never');
  const [maxDownloads, setMaxDownloads] = useState<string>('unlimited');
  const [isLoading, setIsLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const [shareLink, setShareLink] = useState('');
  const [successMessage, setSuccessMessage] = useState('');
  const [showSuccess, setShowSuccess] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [showError, setShowError] = useState(false);

  useEffect(() => {
    setShowSuccess(false);
    setShowError(false);
    setSuccessMessage('');
    setErrorMessage('');

    if (existingShare) {
      setAccessType(existingShare.access_type);
      setExpiresIn(existingShare.expires_at ? 'custom' : 'never');
      setMaxDownloads(existingShare.max_downloads ? String(existingShare.max_downloads) : 'unlimited');
      setShareLink(`${window.location.origin}/s/${existingShare.share_token}`);
    } else {
      setPassword('');
      setAccessType('view');
      setExpiresIn('never');
      setMaxDownloads('unlimited');
      setShareLink('');
    }
  }, [existingShare, open]);

  const handleSubmit = async () => {
    setIsLoading(true);
    try {
      const data: CreateShareRequest | UpdateShareRequest = {
        file_id: fileId,
        access_type: accessType,
      };

      if (password) {
        data.password = password;
      }

      if (expiresIn !== 'never') {
        const expiresAt = new Date();
        switch (expiresIn) {
          case '1h':
            expiresAt.setHours(expiresAt.getHours() + 1);
            break;
          case '1d':
            expiresAt.setDate(expiresAt.getDate() + 1);
            break;
          case '7d':
            expiresAt.setDate(expiresAt.getDate() + 7);
            break;
          case '30d':
            expiresAt.setDate(expiresAt.getDate() + 30);
            break;
        }
        data.expires_at = expiresAt.toISOString();
      }

      if (maxDownloads !== 'unlimited') {
        data.max_downloads = parseInt(maxDownloads);
      }

      let share: Share;
      if (existingShare) {
        share = await onUpdate(existingShare.id, data);
        setSuccessMessage('Share updated successfully');
      } else {
        share = await onCreate(data);
        setSuccessMessage('Share created successfully');
      }

      setShareLink(`${window.location.origin}/s/${share.share_token}`);
      setShowSuccess(true);
      setTimeout(() => setShowSuccess(false), 3000);
    } catch (error: any) {
      console.error('Failed to save share:', error);
      setErrorMessage(error.message || 'Failed to save share');
      setShowError(true);
      setTimeout(() => setShowError(false), 3000);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopyLink = () => {
    navigator.clipboard.writeText(shareLink);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleDelete = async () => {
    if (existingShare) {
      setIsLoading(true);
      try {
        await onDelete(existingShare.id);
        setSuccessMessage('Share deleted successfully');
        setShowSuccess(true);
        setTimeout(() => {
          setShowSuccess(false);
          onClose();
        }, 1500);
      } catch (error: any) {
        console.error('Failed to delete share:', error);
        setErrorMessage(error.message || 'Failed to delete share');
        setShowError(true);
        setTimeout(() => setShowError(false), 3000);
      } finally {
        setIsLoading(false);
      }
    }
  };

  const hasPassword = !!password || existingShare?.password_hash;

  return (
    <Dialog open={open} onClose={onClose} title="Share File" size="medium">
      <div className="space-y-5">
        {showSuccess && successMessage && (
          <div className="bg-green-50 border border-green-200 rounded-lg p-3">
            <p className="text-sm text-green-900">{successMessage}</p>
          </div>
        )}

        {showError && errorMessage && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-3">
            <p className="text-sm text-red-900">{errorMessage}</p>
          </div>
        )}

        <div>
          <p className="text-sm text-gray-600">File: <span className="font-medium text-gray-900">{fileName}</span></p>
        </div>

        {shareLink && (
          <Card padding="medium" className="bg-blue-50 border-blue-200">
            <div className="flex items-center gap-2">
              <Link2 className="w-5 h-5 text-blue-600 flex-shrink-0" />
              <div className="flex-1 min-w-0">
                <p className="text-sm font-mono text-blue-900 truncate">{shareLink}</p>
              </div>
              <Button
                size="small"
                variant="secondary"
                leftIcon={copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
                onClick={handleCopyLink}
              >
                {copied ? 'Copied' : 'Copy'}
              </Button>
            </div>
          </Card>
        )}

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Access Type
          </label>
          <Select
            value={accessType}
            onChange={(value) => setAccessType(value as any)}
            options={[
              { value: 'view', label: 'View only' },
              { value: 'download', label: 'View & Download' },
              { value: 'edit', label: 'Full access' },
            ]}
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            <div className="flex items-center gap-2">
              {hasPassword ? <Lock className="w-4 h-4" /> : <Unlock className="w-4 h-4" />}
              Password Protection
            </div>
          </label>
          <Input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder={existingShare?.password_hash ? 'Enter new password to change' : 'Leave empty for no password'}
          />
          {existingShare?.password_hash && !password && (
            <p className="text-xs text-gray-500 mt-1">Current password is active. Enter a new one to change it.</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Expiration
          </label>
          <Select
            value={expiresIn}
            onChange={(value) => setExpiresIn(value)}
            options={[
              { value: 'never', label: 'Never' },
              { value: '1h', label: '1 hour' },
              { value: '1d', label: '1 day' },
              { value: '7d', label: '7 days' },
              { value: '30d', label: '30 days' },
            ]}
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Download Limit
          </label>
          <Select
            value={maxDownloads}
            onChange={(value) => setMaxDownloads(value)}
            options={[
              { value: 'unlimited', label: 'Unlimited' },
              { value: '1', label: '1 download' },
              { value: '5', label: '5 downloads' },
              { value: '10', label: '10 downloads' },
              { value: '50', label: '50 downloads' },
            ]}
          />
        </div>

        <div className="flex gap-3 pt-4 border-t border-gray-200">
          {existingShare ? (
            <>
              <Button
                variant="secondary"
                onClick={onClose}
                className="flex-1"
              >
                Close
              </Button>
              <Button
                variant="danger"
                onClick={handleDelete}
                loading={isLoading}
              >
                Delete Share
              </Button>
              <Button
                variant="primary"
                onClick={handleSubmit}
                loading={isLoading}
              >
                Update
              </Button>
            </>
          ) : (
            <>
              <Button
                variant="secondary"
                onClick={onClose}
                className="flex-1"
              >
                Cancel
              </Button>
              <Button
                variant="primary"
                onClick={handleSubmit}
                loading={isLoading}
              >
                Create Share
              </Button>
            </>
          )}
        </div>
      </div>
    </Dialog>
  );
}
