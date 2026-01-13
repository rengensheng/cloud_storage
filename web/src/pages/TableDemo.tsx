import { useState, useMemo } from 'react';
import type { ColumnDef } from '@tanstack/react-table';
import { Table, Card, Container, Button } from '../components/ui';
import { Edit, Trash2, Eye } from 'lucide-react';
import '../styles/globals.css';

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  status: 'active' | 'inactive' | 'pending';
  joinDate: string;
}

// 模拟数据
const generateMockData = (count: number): User[] => {
  const roles = ['Admin', 'Editor', 'Viewer', 'Contributor'];
  const statuses: ('active' | 'inactive' | 'pending')[] = ['active', 'inactive', 'pending'];
  const names = [
    'John Doe', 'Jane Smith', 'Bob Johnson', 'Alice Williams',
    'Charlie Brown', 'Diana Prince', 'Edward Norton', 'Fiona Apple',
    'George Martin', 'Hannah Montana', 'Ian McKellen', 'Julia Roberts'
  ];

  return Array.from({ length: count }, (_, i) => ({
    id: i + 1,
    name: names[i % names.length] + ` ${Math.floor(i / names.length) + 1}`,
    email: `user${i + 1}@example.com`,
    role: roles[Math.floor(Math.random() * roles.length)],
    status: statuses[Math.floor(Math.random() * statuses.length)],
    joinDate: new Date(2020 + Math.floor(Math.random() * 4), Math.floor(Math.random() * 12), Math.floor(Math.random() * 28) + 1).toLocaleDateString()
  }));
};

export default function TableDemo() {
  const [data] = useState<User[]>(() => generateMockData(50));
  const [loading, setLoading] = useState(false);
  const [selectedRows, setSelectedRows] = useState<User[]>([]);

  // 模拟加载
  const handleRefresh = () => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
    }, 1500);
  };

  // 定义列
  const columns = useMemo<ColumnDef<User, any>[]>(
    () => [
      {
        accessorKey: 'id',
        header: 'ID',
        size: 60,
        cell: (info) => <span style={{ fontWeight: 600 }}>#{info.getValue()}</span>
      },
      {
        accessorKey: 'name',
        header: 'Name',
        size: 150,
        cell: (info) => (
          <div style={{ fontWeight: 500 }}>
            {info.getValue()}
          </div>
        )
      },
      {
        accessorKey: 'email',
        header: 'Email',
        size: 200
      },
      {
        accessorKey: 'role',
        header: 'Role',
        size: 100,
        cell: (info) => (
          <span className="ui-table__badge ui-table__badge--info">
            {info.getValue()}
          </span>
        )
      },
      {
        accessorKey: 'status',
        header: 'Status',
        size: 100,
        cell: (info) => {
          const status = info.getValue();
          const badgeClass =
            status === 'active' ? 'ui-table__badge--success' :
            status === 'inactive' ? 'ui-table__badge--danger' :
            'ui-table__badge--warning';

          return (
            <span className={`ui-table__badge ${badgeClass}`}>
              {status.charAt(0).toUpperCase() + status.slice(1)}
            </span>
          );
        }
      },
      {
        accessorKey: 'joinDate',
        header: 'Join Date',
        size: 120
      },
      {
        id: 'actions',
        header: 'Actions',
        size: 120,
        cell: ({ row }) => (
          <div className="ui-table__actions">
            <button
              className="ui-table__action-btn"
              onClick={(e) => {
                e.stopPropagation();
                alert(`View user: ${row.original.name}`);
              }}
              title="View"
            >
              <Eye />
            </button>
            <button
              className="ui-table__action-btn"
              onClick={(e) => {
                e.stopPropagation();
                alert(`Edit user: ${row.original.name}`);
              }}
              title="Edit"
            >
              <Edit />
            </button>
            <button
              className="ui-table__action-btn ui-table__action-btn--danger"
              onClick={(e) => {
                e.stopPropagation();
                alert(`Delete user: ${row.original.name}`);
              }}
              title="Delete"
            >
              <Trash2 />
            </button>
          </div>
        )
      }
    ],
    []
  );

  return (
    <div style={{ minHeight: '100vh', paddingTop: '2rem', paddingBottom: '2rem', backgroundColor: 'var(--color-bg)' }}>
      <Container maxWidth="large">
        {/* Header */}
        <div style={{ marginBottom: '2rem' }}>
          <h1 className="text-h1" style={{ marginBottom: '0.5rem' }}>
            Table Component Demo
          </h1>
          <p className="text-body-l" style={{ color: 'var(--color-text-muted)' }}>
            完整功能的数据表格组件，基于 TanStack Table 构建
          </p>
        </div>

        {/* Features Info */}
        <Card padding="medium" shadow="small" style={{ marginBottom: '2rem' }}>
          <h3 className="text-h3" style={{ marginBottom: '1rem' }}>功能特性</h3>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem' }}>
            <div>
              <strong>✅ 排序</strong>
              <p className="text-caption">点击表头排序</p>
            </div>
            <div>
              <strong>✅ 分页</strong>
              <p className="text-caption">自动分页导航</p>
            </div>
            <div>
              <strong>✅ 行选择</strong>
              <p className="text-caption">多选/单选支持</p>
            </div>
            <div>
              <strong>✅ 斑马条纹</strong>
              <p className="text-caption">隔行变色</p>
            </div>
            <div>
              <strong>✅ 悬停高亮</strong>
              <p className="text-caption">鼠标悬停效果</p>
            </div>
            <div>
              <strong>✅ 固定表头</strong>
              <p className="text-caption">滚动时表头固定</p>
            </div>
            <div>
              <strong>✅ 加载状态</strong>
              <p className="text-caption">Loading 动画</p>
            </div>
            <div>
              <strong>✅ 操作列</strong>
              <p className="text-caption">自定义操作按钮</p>
            </div>
          </div>
        </Card>

        {/* Selected Rows Info */}
        {selectedRows.length > 0 && (
          <Card padding="medium" shadow="small" style={{ marginBottom: '1rem', backgroundColor: 'var(--color-primary-light)' }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <span className="text-body-m">
                <strong>{selectedRows.length}</strong> 行已选中
              </span>
              <Button size="small" onClick={() => setSelectedRows([])}>
                清除选择
              </Button>
            </div>
          </Card>
        )}

        {/* Controls */}
        <Card padding="medium" shadow="small" style={{ marginBottom: '1rem' }}>
          <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
            <Button onClick={handleRefresh} loading={loading}>
              刷新数据
            </Button>
            <Button variant="secondary" onClick={() => alert('导出功能演示')}>
              导出数据
            </Button>
            <Button variant="tertiary" onClick={() => alert('添加用户')}>
              添加用户
            </Button>
          </div>
        </Card>

        {/* Table 1: Full Featured */}
        <Card padding="none" shadow="medium" style={{ marginBottom: '2rem' }}>
          <div style={{ padding: 'var(--spacing-4)', borderBottom: '1px solid var(--color-border)' }}>
            <h3 className="text-h3" style={{ margin: 0 }}>用户列表 - 完整功能</h3>
            <p className="text-caption" style={{ marginTop: '0.5rem' }}>
              包含排序、分页、行选择、斑马条纹、悬停高亮等所有功能
            </p>
          </div>
          <Table
            data={data}
            columns={columns}
            loading={loading}
            striped={true}
            hoverable={true}
            stickyActions={true}
            stickySelection={true}
            enableSorting={true}
            enablePagination={true}
            enableRowSelection={true}
            pageSize={10}
            pageSizeOptions={[5, 10, 20, 50]}
            onRowClick={(row) => {
              console.log('Row clicked:', row.original);
            }}
            onRowSelectionChange={setSelectedRows}
          />
        </Card>

        {/* Table 2: With Fixed Header */}
        <Card padding="none" shadow="medium" style={{ marginBottom: '2rem' }}>
          <div style={{ padding: 'var(--spacing-4)', borderBottom: '1px solid var(--color-border)' }}>
            <h3 className="text-h3" style={{ margin: 0 }}>固定表头示例</h3>
            <p className="text-caption" style={{ marginTop: '0.5rem' }}>
              表格滚动时表头保持固定
            </p>
          </div>
          <Table
            data={data}
            columns={columns}
            loading={loading}
            stickyHeader={true}
            bordered={true}
            enablePagination={false}
            maxHeight="400px"
          />
        </Card>

        {/* Table 3: Simple Table */}
        <Card padding="none" shadow="medium">
          <div style={{ padding: 'var(--spacing-4)', borderBottom: '1px solid var(--color-border)' }}>
            <h3 className="text-h3" style={{ margin: 0 }}>简单表格</h3>
            <p className="text-caption" style={{ marginTop: '0.5rem' }}>
              最小配置的表格示例
            </p>
          </div>
          <Table
            data={data.slice(0, 5)}
            columns={columns.slice(0, 4)}
            enablePagination={false}
            enableSorting={false}
          />
        </Card>

        {/* Footer */}
        <div style={{ marginTop: '3rem', textAlign: 'center', padding: '2rem 0', borderTop: '1px solid var(--color-border)' }}>
          <p className="text-caption">
            表格组件基于 @tanstack/react-table v8 构建
          </p>
        </div>
      </Container>
    </div>
  );
}
