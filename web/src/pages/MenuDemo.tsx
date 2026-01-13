import { useState } from 'react';
import { MenuBar, Button, Card, Container } from '../components/ui';
import type { MenuBarMenu } from '../components/ui';
import {
  FileText,
  FolderOpen,
  Save,
  Upload,
  Download,
  Printer,
  X,
  Copy,
  Scissors,
  Clipboard,
  Search,
  RotateCcw,
  RotateCw,
  ZoomIn,
  ZoomOut,
  Maximize,
  Settings,
  HelpCircle,
  BookOpen,
  Info,
  Keyboard,
  Package,
  LogOut,
  User,
  Bell,
  Mail,
  Home
} from 'lucide-react';
import '../styles/globals.css';

export default function MenuDemo() {
  const [autoSave, setAutoSave] = useState(false);
  const [showToolbar, setShowToolbar] = useState(true);
  const [showSidebar, setShowSidebar] = useState(true);

  // File Menu
  const fileMenu: MenuBarMenu = {
    key: 'file',
    label: 'File',
    items: [
      {
        key: 'new',
        label: 'New',
        icon: <FileText />,
        shortcut: '⌘N',
        onClick: () => alert('New File')
      },
      {
        key: 'open',
        label: 'Open',
        icon: <FolderOpen />,
        shortcut: '⌘O',
        children: [
          {
            key: 'open-file',
            label: 'Open File...',
            shortcut: '⌘O',
            onClick: () => alert('Open File')
          },
          {
            key: 'open-folder',
            label: 'Open Folder...',
            shortcut: '⌘⇧O',
            onClick: () => alert('Open Folder')
          },
          {
            key: 'open-recent',
            label: 'Open Recent',
            children: [
              {
                key: 'recent-1',
                label: 'project-1.txt',
                onClick: () => alert('Open project-1.txt')
              },
              {
                key: 'recent-2',
                label: 'document.pdf',
                onClick: () => alert('Open document.pdf')
              },
              {
                key: 'recent-3',
                label: 'notes.md',
                onClick: () => alert('Open notes.md')
              }
            ]
          }
        ]
      },
      { key: 'divider-1', divider: true },
      {
        key: 'save',
        label: 'Save',
        icon: <Save />,
        shortcut: '⌘S',
        onClick: () => alert('Save')
      },
      {
        key: 'save-as',
        label: 'Save As...',
        shortcut: '⌘⇧S',
        onClick: () => alert('Save As')
      },
      {
        key: 'save-all',
        label: 'Save All',
        disabled: true,
        onClick: () => alert('Save All')
      },
      { key: 'divider-2', divider: true },
      {
        key: 'export',
        label: 'Export',
        icon: <Download />,
        children: [
          {
            key: 'export-pdf',
            label: 'Export as PDF',
            onClick: () => alert('Export PDF')
          },
          {
            key: 'export-html',
            label: 'Export as HTML',
            onClick: () => alert('Export HTML')
          },
          {
            key: 'export-md',
            label: 'Export as Markdown',
            onClick: () => alert('Export Markdown')
          }
        ]
      },
      {
        key: 'import',
        label: 'Import',
        icon: <Upload />,
        onClick: () => alert('Import')
      },
      { key: 'divider-3', divider: true },
      {
        key: 'print',
        label: 'Print',
        icon: <Printer />,
        shortcut: '⌘P',
        onClick: () => alert('Print')
      },
      { key: 'divider-4', divider: true },
      {
        key: 'close',
        label: 'Close',
        icon: <X />,
        shortcut: '⌘W',
        onClick: () => alert('Close')
      }
    ]
  };

  // Edit Menu
  const editMenu: MenuBarMenu = {
    key: 'edit',
    label: 'Edit',
    items: [
      {
        key: 'undo',
        label: 'Undo',
        icon: <RotateCcw />,
        shortcut: '⌘Z',
        onClick: () => alert('Undo')
      },
      {
        key: 'redo',
        label: 'Redo',
        icon: <RotateCw />,
        shortcut: '⌘⇧Z',
        onClick: () => alert('Redo')
      },
      { key: 'divider-1', divider: true },
      {
        key: 'cut',
        label: 'Cut',
        icon: <Scissors />,
        shortcut: '⌘X',
        onClick: () => alert('Cut')
      },
      {
        key: 'copy',
        label: 'Copy',
        icon: <Copy />,
        shortcut: '⌘C',
        onClick: () => alert('Copy')
      },
      {
        key: 'paste',
        label: 'Paste',
        icon: <Clipboard />,
        shortcut: '⌘V',
        onClick: () => alert('Paste')
      },
      { key: 'divider-2', divider: true },
      {
        key: 'find',
        label: 'Find',
        icon: <Search />,
        shortcut: '⌘F',
        onClick: () => alert('Find')
      },
      {
        key: 'replace',
        label: 'Replace',
        shortcut: '⌘⌥F',
        onClick: () => alert('Replace')
      }
    ]
  };

  // View Menu
  const viewMenu: MenuBarMenu = {
    key: 'view',
    label: 'View',
    items: [
      {
        key: 'toolbar',
        label: 'Show Toolbar',
        checked: showToolbar,
        onClick: () => {
          setShowToolbar(!showToolbar);
          alert(`Toolbar ${!showToolbar ? 'shown' : 'hidden'}`);
        }
      },
      {
        key: 'sidebar',
        label: 'Show Sidebar',
        checked: showSidebar,
        onClick: () => {
          setShowSidebar(!showSidebar);
          alert(`Sidebar ${!showSidebar ? 'shown' : 'hidden'}`);
        }
      },
      { key: 'divider-1', divider: true },
      {
        key: 'zoom-in',
        label: 'Zoom In',
        icon: <ZoomIn />,
        shortcut: '⌘+',
        onClick: () => alert('Zoom In')
      },
      {
        key: 'zoom-out',
        label: 'Zoom Out',
        icon: <ZoomOut />,
        shortcut: '⌘-',
        onClick: () => alert('Zoom Out')
      },
      {
        key: 'reset-zoom',
        label: 'Reset Zoom',
        shortcut: '⌘0',
        onClick: () => alert('Reset Zoom')
      },
      { key: 'divider-2', divider: true },
      {
        key: 'fullscreen',
        label: 'Full Screen',
        icon: <Maximize />,
        shortcut: '⌘⌃F',
        onClick: () => alert('Toggle Fullscreen')
      }
    ]
  };

  // Help Menu
  const helpMenu: MenuBarMenu = {
    key: 'help',
    label: 'Help',
    items: [
      {
        key: 'docs',
        label: 'Documentation',
        icon: <BookOpen />,
        onClick: () => alert('Open Documentation')
      },
      {
        key: 'shortcuts',
        label: 'Keyboard Shortcuts',
        icon: <Keyboard />,
        shortcut: '⌘/',
        onClick: () => alert('Show Shortcuts')
      },
      {
        key: 'support',
        label: 'Get Support',
        icon: <HelpCircle />,
        onClick: () => alert('Get Support')
      },
      { key: 'divider-1', divider: true },
      {
        key: 'updates',
        label: 'Check for Updates',
        icon: <Package />,
        onClick: () => alert('Check Updates')
      },
      { key: 'divider-2', divider: true },
      {
        key: 'about',
        label: 'About',
        icon: <Info />,
        onClick: () => alert('About Application')
      }
    ]
  };

  // Preferences Menu
  const preferencesMenu: MenuBarMenu = {
    key: 'preferences',
    label: 'Preferences',
    items: [
      {
        key: 'settings',
        label: 'Settings',
        icon: <Settings />,
        shortcut: '⌘,',
        onClick: () => alert('Open Settings')
      },
      {
        key: 'auto-save',
        label: 'Auto Save',
        icon: <Save />,
        checked: autoSave,
        onClick: () => {
          setAutoSave(!autoSave);
          alert(`Auto Save ${!autoSave ? 'enabled' : 'disabled'}`);
        }
      },
      { key: 'divider-1', divider: true },
      {
        key: 'account',
        label: 'Account',
        icon: <User />,
        children: [
          {
            key: 'profile',
            label: 'My Profile',
            onClick: () => alert('View Profile')
          },
          {
            key: 'billing',
            label: 'Billing',
            onClick: () => alert('View Billing')
          },
          { key: 'divider-1', divider: true },
          {
            key: 'logout',
            label: 'Logout',
            icon: <LogOut />,
            danger: true,
            onClick: () => alert('Logout')
          }
        ]
      }
    ]
  };

  const menus = [fileMenu, editMenu, viewMenu, preferencesMenu, helpMenu];

  return (
    <div style={{ minHeight: '100vh' }}>
      {/* MenuBar Demo */}
      <MenuBar
        menus={menus}
        logo={
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <Home size={20} style={{ color: 'var(--color-primary)' }} />
            <span style={{ fontWeight: 700, fontSize: '16px' }}>MyApp</span>
          </div>
        }
        actions={
          <>
            <Button variant="tertiary" style={{ padding: '0.5rem' }}>
              <Bell size={18} />
            </Button>
            <Button variant="tertiary" style={{ padding: '0.5rem' }}>
              <Mail size={18} />
            </Button>
            <Button variant="tertiary" style={{ padding: '0.5rem' }}>
              <User size={18} />
            </Button>
          </>
        }
      />

      {/* Content */}
      <Container maxWidth="large" style={{ paddingTop: '2rem', paddingBottom: '2rem' }}>
        {/* Header */}
        <div style={{ marginBottom: '3rem', textAlign: 'center' }}>
          <h1 className="text-h1" style={{ marginBottom: '0.5rem' }}>
            MenuBar Component
          </h1>
          <p className="text-body-l" style={{ color: 'var(--color-text-muted)' }}>
            应用程序风格的水平菜单栏，类似于桌面应用的顶部菜单
          </p>
        </div>

        {/* Status Card */}
        <Card padding="large" shadow="medium" style={{ marginBottom: '2rem' }}>
          <h3 className="text-h3" style={{ marginBottom: '1rem' }}>当前状态</h3>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem' }}>
            <div>
              <strong className="text-body-m">Auto Save:</strong>
              <p className="text-body-s" style={{ color: 'var(--color-text-muted)', marginTop: '0.25rem' }}>
                {autoSave ? '✅ Enabled' : '❌ Disabled'}
              </p>
            </div>
            <div>
              <strong className="text-body-m">Toolbar:</strong>
              <p className="text-body-s" style={{ color: 'var(--color-text-muted)', marginTop: '0.25rem' }}>
                {showToolbar ? '✅ Visible' : '❌ Hidden'}
              </p>
            </div>
            <div>
              <strong className="text-body-m">Sidebar:</strong>
              <p className="text-body-s" style={{ color: 'var(--color-text-muted)', marginTop: '0.25rem' }}>
                {showSidebar ? '✅ Visible' : '❌ Hidden'}
              </p>
            </div>
          </div>
        </Card>

        {/* Features */}
        <Card padding="large" shadow="medium" style={{ marginBottom: '2rem' }}>
          <h3 className="text-h3" style={{ marginBottom: '1rem' }}>功能特性</h3>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1rem' }}>
            <div>
              <strong className="text-body-m">✅ 水平菜单栏</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                类似桌面应用的顶部菜单
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ Logo 区域</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                可自定义的应用 Logo
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ 多级菜单</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                支持无限层级嵌套
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ 图标支持</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                Lucide React 图标集成
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ 快捷键显示</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                显示键盘快捷键提示
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ 选中状态</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                复选标记显示状态
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ 分隔线</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                菜单项分组
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ 禁用状态</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                灰化不可用选项
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ 危险操作</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                红色高亮删除等操作
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ Actions 区域</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                右侧自定义操作按钮
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ 粘性定位</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                滚动时保持在顶部
              </p>
            </div>
            <div>
              <strong className="text-body-m">✅ 响应式设计</strong>
              <p className="text-caption" style={{ marginTop: '0.5rem' }}>
                移动端适配
              </p>
            </div>
          </div>
        </Card>

        {/* Usage Example */}
        <Card padding="large" shadow="medium">
          <h3 className="text-h3" style={{ marginBottom: '1rem' }}>使用示例</h3>
          <pre
            style={{
              padding: 'var(--spacing-4)',
              backgroundColor: 'var(--color-surface-hover)',
              borderRadius: 'var(--radius-md)',
              overflow: 'auto',
              fontSize: '13px',
              lineHeight: '1.6'
            }}
          >
            {`<MenuBar
  menus={[fileMenu, editMenu, viewMenu]}
  logo={
    <div>
      <Logo />
      <span>MyApp</span>
    </div>
  }
  actions={
    <>
      <Button><Bell /></Button>
      <Button><User /></Button>
    </>
  }
/>`}
          </pre>
        </Card>

        {/* Footer */}
        <div style={{ marginTop: '3rem', textAlign: 'center', padding: '2rem 0', borderTop: '1px solid var(--color-border)' }}>
          <p className="text-caption">
            MenuBar 组件基于 @headlessui/react 构建，完美适配所有主题
          </p>
        </div>
      </Container>
    </div>
  );
}
