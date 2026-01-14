import { useState, useEffect } from 'react';
import {
  User,
  Mail,
  Lock,
  Shield,
  HardDrive,
} from 'lucide-react';
import { TabsControlled } from '../components/ui/TabsControlled';
import { Button, Input, Card, Dialog } from '../components/ui';
import { authApi, statsApi } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

export default function Settings() {
  const { user, logout, refreshUser } = useAuth();
  const [activeTab, setActiveTab] = useState('profile');
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  // Profile form
  const [profileForm, setProfileForm] = useState({
    username: user?.username || '',
    email: user?.email || '',
  });

  // Password form
  const [passwordForm, setPasswordForm] = useState({
    old_password: '',
    new_password: '',
    confirm_password: '',
  });

  const [showLogoutDialog, setShowLogoutDialog] = useState(false);
  const [storageStats, setStorageStats] = useState<{ used: number; quota: number } | null>(null);

  useEffect(() => {
    const loadStorageStats = async () => {
      try {
        const stats = await statsApi.getStorageStats();
        setStorageStats({ used: stats.used, quota: stats.quota });
      } catch (error) {
        console.error('Failed to load storage stats:', error);
      }
    };

    loadStorageStats();
  }, []);

  const handleProfileUpdate = async () => {
    setIsLoading(true);
    setMessage(null);

    try {
      await authApi.updateProfile(profileForm);
      await refreshUser();
      setMessage({ type: 'success', text: 'Profile updated successfully' });
    } catch (error: any) {
      setMessage({ type: 'error', text: error.message || 'Failed to update profile' });
    } finally {
      setIsLoading(false);
    }
  };

  const handlePasswordChange = async () => {
    if (passwordForm.new_password !== passwordForm.confirm_password) {
      setMessage({ type: 'error', text: 'Passwords do not match' });
      return;
    }

    if (passwordForm.new_password.length < 6) {
      setMessage({ type: 'error', text: 'Password must be at least 6 characters' });
      return;
    }

    setIsLoading(true);
    setMessage(null);

    try {
      await authApi.changePassword({
        old_password: passwordForm.old_password,
        new_password: passwordForm.new_password,
      });
      setPasswordForm({ old_password: '', new_password: '', confirm_password: '' });
      setMessage({ type: 'success', text: 'Password changed successfully' });
    } catch (error: any) {
      setMessage({ type: 'error', text: error.message || 'Failed to change password' });
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = async () => {
    await logout();
    window.location.href = '/login';
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleDateString();
  };

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <p className="text-gray-600">Loading...</p>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Settings</h1>
        <p className="text-sm text-gray-600">Manage your account settings</p>
      </div>

      {/* Tabs */}
      <TabsControlled
        tabs={[
          { id: 'profile', label: 'Profile', icon: <User className="w-4 h-4" /> },
          { id: 'security', label: 'Security', icon: <Shield className="w-4 h-4" /> },
          { id: 'storage', label: 'Storage', icon: <HardDrive className="w-4 h-4" /> },
        ]}
        activeTab={activeTab}
        onChange={setActiveTab}
      />

      {/* Message */}
      {message && (
        <Card
          padding="medium"
          className={message.type === 'success' ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'}
        >
          <p className={message.type === 'success' ? 'text-green-900' : 'text-red-900'}>
            {message.text}
          </p>
        </Card>
      )}

      {/* Profile Tab */}
      {activeTab === 'profile' && (
        <div className="space-y-6">
          <Card padding="large">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Profile Information</h2>

            <div className="flex items-center gap-4 mb-6 p-4 bg-gray-50 rounded-lg">
              <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-cyan-500 rounded-full flex items-center justify-center">
                <User className="w-8 h-8 text-white" />
              </div>
              <div>
                <h3 className="font-medium text-gray-900">{user?.username}</h3>
                <p className="text-sm text-gray-600">{user?.email}</p>
                <span className="inline-block mt-1 px-2 py-0.5 text-xs font-medium bg-blue-100 text-blue-800 rounded">
                  {user?.role}
                </span>
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Username
                </label>
                <Input
                  type="text"
                  value={profileForm.username}
                  onChange={(e) => setProfileForm({ ...profileForm, username: e.target.value })}
                  leftIcon={<User className="w-5 h-5 text-gray-400" />}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Email
                </label>
                <Input
                  type="email"
                  value={profileForm.email}
                  onChange={(e) => setProfileForm({ ...profileForm, email: e.target.value })}
                  leftIcon={<Mail className="w-5 h-5 text-gray-400" />}
                />
              </div>

              <div className="flex justify-end pt-4 border-t border-gray-200">
                <Button
                  variant="primary"
                  onClick={handleProfileUpdate}
                  loading={isLoading}
                  disabled={
                    profileForm.username === user?.username &&
                    profileForm.email === user?.email
                  }
                >
                  Save Changes
                </Button>
              </div>
            </div>
          </Card>

          <Card padding="large">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Account Info</h2>
            <div className="space-y-3">
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-sm text-gray-600">Member Since</span>
                <span className="text-sm font-medium text-gray-900">
                  {formatDate(user?.created_at)}
                </span>
              </div>
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-sm text-gray-600">Last Login</span>
                <span className="text-sm font-medium text-gray-900">
                  {formatDate(user?.last_login_at)}
                </span>
              </div>
              <div className="flex items-center justify-between py-2">
                <span className="text-sm text-gray-600">Account Status</span>
                <span className={`inline-flex items-center px-2 py-1 text-xs font-medium rounded ${
                  user?.is_active
                    ? 'bg-green-100 text-green-800'
                    : 'bg-red-100 text-red-800'
                }`}>
                  {user?.is_active ? 'Active' : 'Inactive'}
                </span>
              </div>
            </div>
          </Card>
        </div>
      )}

      {/* Security Tab */}
      {activeTab === 'security' && (
        <div className="space-y-6">
          <Card padding="large">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Change Password</h2>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Current Password
                </label>
                <Input
                  type="password"
                  value={passwordForm.old_password}
                  onChange={(e) => setPasswordForm({ ...passwordForm, old_password: e.target.value })}
                  leftIcon={<Lock className="w-5 h-5 text-gray-400" />}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  New Password
                </label>
                <Input
                  type="password"
                  value={passwordForm.new_password}
                  onChange={(e) => setPasswordForm({ ...passwordForm, new_password: e.target.value })}
                  leftIcon={<Lock className="w-5 h-5 text-gray-400" />}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Confirm New Password
                </label>
                <Input
                  type="password"
                  value={passwordForm.confirm_password}
                  onChange={(e) => setPasswordForm({ ...passwordForm, confirm_password: e.target.value })}
                  leftIcon={<Lock className="w-5 h-5 text-gray-400" />}
                />
              </div>

              <div className="flex justify-end pt-4 border-t border-gray-200">
                <Button
                  variant="primary"
                  onClick={handlePasswordChange}
                  loading={isLoading}
                  disabled={!passwordForm.old_password || !passwordForm.new_password}
                >
                  Change Password
                </Button>
              </div>
            </div>
          </Card>

          <Card padding="large">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Account Actions</h2>
            <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
              <div>
                <h3 className="font-medium text-gray-900">Sign Out</h3>
                <p className="text-sm text-gray-600">Sign out from your account on this device</p>
              </div>
              <Button variant="secondary" onClick={() => setShowLogoutDialog(true)}>
                Sign Out
              </Button>
            </div>
          </Card>
        </div>
      )}

      {/* Storage Tab */}
      {activeTab === 'storage' && (
        <div className="space-y-6">
          <Card padding="large">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Storage Usage</h2>

            <div className="mb-6">
              <div className="flex justify-between mb-2">
                <span className="text-sm font-medium text-gray-900">Used Space</span>
                <span className="text-sm text-gray-600">
                  {storageStats ? formatFileSize(storageStats.used) : 'Loading...'} /{' '}
                  {storageStats ? formatFileSize(storageStats.quota) : '...'}
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
                <div
                  className="bg-gradient-to-r from-blue-500 to-cyan-500 h-full rounded-full transition-all duration-500"
                  style={{
                    width: storageStats
                      ? `${Math.min((storageStats.used / storageStats.quota) * 100, 100)}%`
                      : '0%'
                  }}
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="p-4 bg-blue-50 rounded-lg">
                <p className="text-sm text-blue-600 font-medium">Used</p>
                <p className="text-2xl font-bold text-blue-900">
                  {storageStats ? formatFileSize(storageStats.used) : '-'}
                </p>
              </div>
              <div className="p-4 bg-green-50 rounded-lg">
                <p className="text-sm text-green-600 font-medium">Available</p>
                <p className="text-2xl font-bold text-green-900">
                  {storageStats ? formatFileSize(storageStats.quota - storageStats.used) : '-'}
                </p>
              </div>
            </div>
          </Card>

          <Card padding="large">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Storage Plan</h2>
            <div className="p-4 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-lg text-white">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-semibold text-lg">Free Plan</h3>
                  <p className="text-blue-100 text-sm">
                    {storageStats ? formatFileSize(storageStats.quota) : '...'} storage
                  </p>
                </div>
                <HardDrive className="w-10 h-10 text-white/80" />
              </div>
            </div>
            <p className="text-sm text-gray-600 mt-4">
              You are currently on the free plan. Upgrade to get more storage space.
            </p>
          </Card>
        </div>
      )}

      {/* Logout Dialog */}
      <Dialog
        open={showLogoutDialog}
        onClose={() => setShowLogoutDialog(false)}
        title="Sign Out"
        size="small"
      >
        <div className="space-y-4">
          <p className="text-gray-600">Are you sure you want to sign out?</p>
          <div className="flex justify-end gap-3">
            <Button variant="secondary" onClick={() => setShowLogoutDialog(false)}>
              Cancel
            </Button>
            <Button variant="danger" onClick={handleLogout}>
              Sign Out
            </Button>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
