import { useState, useEffect } from 'react';
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import {
  FolderKanban,
  Trash2,
  Search,
  Settings,
  LogOut,
  Menu,
  X,
  Upload,
  User,
  HardDrive,
} from 'lucide-react';
import { Button, Card, Dialog } from '../components/ui';
import { Logo } from '../components/Logo';
import { useAuth } from '../contexts/AuthContext';
import type { StorageStats } from '../types/api';
import { statsApi } from '../services/api';

export default function MainLayout() {
  const { user, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [storageStats, setStorageStats] = useState<StorageStats | null>(null);
  const [showLogoutDialog, setShowLogoutDialog] = useState(false);

  const loadStorageStats = async () => {
    try {
      const stats = await statsApi.getStorageStats();
      setStorageStats(stats);
    } catch (error) {
      console.error('Failed to load storage stats:', error);
    }
  };

  useEffect(() => {
    loadStorageStats();
  }, []);

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const navigation = [
    { name: '我的文件', href: '/', icon: FolderKanban },
    { name: '搜索文件', href: '/search', icon: Search },
    { name: '回收站', href: '/recycle', icon: Trash2 },
    { name: '设置', href: '/settings', icon: Settings },
  ];

  return (
    <div className="min-h-screen bg-gray-50/50">
      <Dialog
        open={showLogoutDialog}
        onClose={() => setShowLogoutDialog(false)}
        title="退出登录"
        size="small"
      >
        <div className="space-y-4">
          <p className="text-gray-600">Are you sure you want to logout?</p>
          <div className="flex justify-end gap-3">
            <Button
              variant="secondary"
              onClick={() => setShowLogoutDialog(false)}
            >
              Cancel
            </Button>
            <Button variant="danger" onClick={handleLogout}>
              Logout
            </Button>
          </div>
        </div>
      </Dialog>

      <nav className="fixed top-0 left-0 right-0 h-16 bg-white/80 backdrop-blur-lg border-b border-gray-200 z-50">
        <div className="flex items-center justify-between h-full px-4">
          <div className="flex items-center gap-4">
            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="p-2 hover:bg-gray-100 rounded-lg transition-colors lg:hidden"
            >
              {sidebarOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
            <Link to="/" className="flex items-center gap-2">
              <Logo size={32} />
              <span className="font-bold text-lg bg-gradient-to-r from-blue-600 to-cyan-600 bg-clip-text text-transparent">
                Cloud Storage
              </span>
            </Link>
          </div>

          <div className="flex items-center gap-3">
            <Button
              variant="primary"
              size="small"
              leftIcon={<Upload className="w-4 h-4" />}
              onClick={() => navigate('/upload')}
            >
              <span className="hidden sm:inline">Upload</span>
            </Button>
            <div className="flex items-center gap-2 ml-4 pl-4 border-l border-gray-200">
              <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-cyan-500 rounded-full flex items-center justify-center">
                <User className="w-4 h-4 text-white" />
              </div>
              <span className="hidden md:block text-sm font-medium text-gray-700">{user?.username}</span>
            </div>
            <button
              onClick={() => setShowLogoutDialog(true)}
              className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              title="Logout"
            >
              <LogOut className="w-5 h-5 text-gray-500" />
            </button>
          </div>
        </div>
      </nav>

      <aside
        className={`fixed left-0 top-16 bottom-0 w-64 bg-white border-r border-gray-200 z-40 transition-transform duration-300 lg:translate-x-0 ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="p-4 space-y-6">
          <div>
            <Card padding="medium" shadow="small">
              <div className="flex items-center gap-3 mb-3">
                <HardDrive className="w-5 h-5 text-blue-600" />
                <span className="font-semibold text-gray-900">Storage</span>
              </div>
              {storageStats && (
                <>
                  <div className="space-y-2 mb-3">
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Used</span>
                      <span className="font-medium">{storageStats.usage_readable}</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                      <div
                        className="bg-gradient-to-r from-blue-500 to-cyan-500 h-full rounded-full transition-all duration-500"
                        style={{ width: `${Math.min(storageStats.usage_percent, 100)}%` }}
                      />
                    </div>
                  </div>
                  <p className="text-xs text-gray-500">
                    {storageStats.available > 0 ? `${(storageStats.available / (1024 * 1024 * 1024)).toFixed(2)} GB available` : 'Full'}
                  </p>
                </>
              )}
            </Card>
          </div>

          <nav className="space-y-1">
            {navigation.map((item) => {
              const isActive = location.pathname === item.href;
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-200 ${
                    isActive
                      ? 'bg-gradient-to-r from-blue-50 to-cyan-50 text-blue-700 font-medium'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`}
                >
                  <item.icon className={`w-5 h-5 ${isActive ? 'text-blue-600' : ''}`} />
                  <span>{item.name}</span>
                </Link>
              );
            })}
          </nav>
        </div>
      </aside>

      <main className="lg:ml-64 pt-16 min-h-screen">
        <div className="p-4 lg:p-6">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
