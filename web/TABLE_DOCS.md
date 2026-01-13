# Table Component 文档

功能完整的数据表格组件，基于 @tanstack/react-table 构建，支持排序、分页、行选择等高级功能。

## 主要特性

### ✅ 核心功能
- **排序 (Sorting)**: 点击表头进行升序/降序排序
- **分页 (Pagination)**: 自动分页，可自定义每页显示数量
- **行选择 (Row Selection)**: 支持单选/多选行
- **加载状态 (Loading)**: 优雅的 loading 动画
- **空状态 (Empty State)**: 无数据时的友好提示

### 🎨 样式功能
- **斑马条纹 (Striped)**: 隔行变色，提高可读性
- **悬停高亮 (Hoverable)**: 鼠标悬停时高亮行
- **边框模式 (Bordered)**: 可选的表格边框
- **固定表头 (Sticky Header)**: 滚动时保持表头可见
- **自定义高度**: 可设置最大高度并启用滚动

### 🔧 高级功能
- **操作列**: 内置操作按钮样式
- **状态徽章**: 预设的状态标签样式
- **点击行事件**: 可配置行点击回调
- **响应式设计**: 移动端友好
- **TypeScript 支持**: 完整的类型定义

## 基础用法

```tsx
import { Table } from './components/ui';
import type { ColumnDef } from '@tanstack/react-table';

interface User {
  id: number;
  name: string;
  email: string;
}

function MyTable() {
  const data: User[] = [
    { id: 1, name: 'John', email: 'john@example.com' },
    { id: 2, name: 'Jane', email: 'jane@example.com' }
  ];

  const columns: ColumnDef<User>[] = [
    {
      accessorKey: 'id',
      header: 'ID'
    },
    {
      accessorKey: 'name',
      header: 'Name'
    },
    {
      accessorKey: 'email',
      header: 'Email'
    }
  ];

  return <Table data={data} columns={columns} />;
}
```

## Props 详解

```typescript
interface TableProps<TData> {
  // 必需属性
  data: TData[];                          // 表格数据
  columns: ColumnDef<TData, any>[];       // 列定义

  // 样式属性
  striped?: boolean;                      // 斑马条纹 (默认: false)
  hoverable?: boolean;                    // 悬停高亮 (默认: true)
  bordered?: boolean;                     // 显示边框 (默认: false)
  stickyHeader?: boolean;                 // 固定表头 (默认: false)
  maxHeight?: string;                     // 最大高度，如 '400px'
  className?: string;                     // 自定义类名

  // 功能属性
  loading?: boolean;                      // 加载状态 (默认: false)
  enableSorting?: boolean;                // 启用排序 (默认: true)
  enablePagination?: boolean;             // 启用分页 (默认: true)
  enableRowSelection?: boolean;           // 启用行选择 (默认: false)

  // 分页属性
  pageSize?: number;                      // 每页数量 (默认: 10)
  pageSizeOptions?: number[];             // 分页选项 (默认: [10, 20, 50, 100])

  // 回调函数
  onRowClick?: (row: Row<TData>) => void; // 行点击事件
  onRowSelectionChange?: (selectedRows: TData[]) => void; // 选择变化

  // 自定义文本
  emptyMessage?: string;                  // 空数据提示
}
```

## 完整示例

### 1. 带所有功能的表格

```tsx
import { useState } from 'react';
import { Table } from './components/ui';
import type { ColumnDef } from '@tanstack/react-table';

interface User {
  id: number;
  name: string;
  email: string;
  status: 'active' | 'inactive';
}

function FullFeaturedTable() {
  const [selectedRows, setSelectedRows] = useState<User[]>([]);

  const columns: ColumnDef<User>[] = [
    {
      accessorKey: 'id',
      header: 'ID',
      size: 60
    },
    {
      accessorKey: 'name',
      header: 'Name',
      size: 150
    },
    {
      accessorKey: 'email',
      header: 'Email',
      size: 200
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: (info) => (
        <span className={`ui-table__badge ui-table__badge--${
          info.getValue() === 'active' ? 'success' : 'danger'
        }`}>
          {info.getValue()}
        </span>
      )
    }
  ];

  return (
    <Table
      data={data}
      columns={columns}
      striped={true}
      hoverable={true}
      enableSorting={true}
      enablePagination={true}
      enableRowSelection={true}
      pageSize={10}
      onRowClick={(row) => console.log('Clicked:', row.original)}
      onRowSelectionChange={setSelectedRows}
    />
  );
}
```

### 2. 带操作列的表格

```tsx
import { Edit, Trash2, Eye } from 'lucide-react';

const columns: ColumnDef<User>[] = [
  // ... 其他列
  {
    id: 'actions',
    header: 'Actions',
    cell: ({ row }) => (
      <div className="ui-table__actions">
        <button
          className="ui-table__action-btn"
          onClick={(e) => {
            e.stopPropagation();
            handleView(row.original);
          }}
          title="View"
        >
          <Eye />
        </button>
        <button
          className="ui-table__action-btn"
          onClick={(e) => {
            e.stopPropagation();
            handleEdit(row.original);
          }}
          title="Edit"
        >
          <Edit />
        </button>
        <button
          className="ui-table__action-btn ui-table__action-btn--danger"
          onClick={(e) => {
            e.stopPropagation();
            handleDelete(row.original);
          }}
          title="Delete"
        >
          <Trash2 />
        </button>
      </div>
    )
  }
];
```

### 3. 固定表头的表格

```tsx
<Table
  data={data}
  columns={columns}
  stickyHeader={true}
  maxHeight="500px"
  bordered={true}
/>
```

### 4. 带状态徽章

```tsx
{
  accessorKey: 'status',
  header: 'Status',
  cell: (info) => {
    const status = info.getValue();
    const badgeClass =
      status === 'active' ? 'ui-table__badge--success' :
      status === 'pending' ? 'ui-table__badge--warning' :
      'ui-table__badge--danger';

    return (
      <span className={`ui-table__badge ${badgeClass}`}>
        {status}
      </span>
    );
  }
}
```

## 内置 CSS 类

### 表格状态徽章
```tsx
// 成功状态（绿色）
<span className="ui-table__badge ui-table__badge--success">Active</span>

// 警告状态（黄色）
<span className="ui-table__badge ui-table__badge--warning">Pending</span>

// 危险状态（红色）
<span className="ui-table__badge ui-table__badge--danger">Inactive</span>

// 信息状态（蓝色）
<span className="ui-table__badge ui-table__badge--info">Admin</span>
```

### 操作按钮
```tsx
// 普通操作按钮
<button className="ui-table__action-btn">
  <Icon />
</button>

// 危险操作按钮（删除等）
<button className="ui-table__action-btn ui-table__action-btn--danger">
  <Trash2 />
</button>
```

## 列配置 (Column Definition)

TanStack Table 支持的列配置：

```typescript
const columns: ColumnDef<DataType>[] = [
  {
    accessorKey: 'fieldName',    // 数据字段名
    header: 'Column Title',       // 列标题
    size: 150,                    // 列宽度（像素）
    enableSorting: true,          // 是否可排序
    cell: (info) => {             // 自定义单元格渲染
      return <div>{info.getValue()}</div>;
    }
  }
];
```

## Pagination 组件

分页组件会自动集成到 Table 中，但也可以单独使用：

```tsx
import { Pagination } from './components/ui';

<Pagination
  currentPage={1}
  totalPages={10}
  pageSize={20}
  totalItems={200}
  pageSizeOptions={[10, 20, 50]}
  onPageChange={(page) => console.log(page)}
  onPageSizeChange={(size) => console.log(size)}
/>
```

## 性能优化建议

1. **使用 useMemo 缓存列定义**
```tsx
const columns = useMemo(() => [...], []);
```

2. **大数据集考虑虚拟滚动**
```tsx
// 对于超过 1000 行的数据，建议使用虚拟滚动库
// 如 @tanstack/react-virtual
```

3. **异步加载数据**
```tsx
const [loading, setLoading] = useState(false);

useEffect(() => {
  setLoading(true);
  fetchData().then(data => {
    setData(data);
    setLoading(false);
  });
}, []);
```

## 样式自定义

### 通过 CSS 变量
```css
.ui-table {
  --table-border-color: var(--color-border);
  --table-hover-bg: var(--color-primary-light);
}
```

### 通过 className
```tsx
<Table
  className="my-custom-table"
  data={data}
  columns={columns}
/>
```

## 可访问性 (A11y)

表格组件已实现以下可访问性特性：

- ✅ 语义化 HTML 标签 (`<table>`, `<thead>`, `<tbody>`)
- ✅ ARIA 标签支持
- ✅ 键盘导航支持
- ✅ 屏幕阅读器友好
- ✅ 高对比度模式支持

## 浏览器支持

- Chrome (最新版)
- Firefox (最新版)
- Safari (最新版)
- Edge (最新版)

## 相关资源

- [TanStack Table 文档](https://tanstack.com/table/latest)
- [Lucide React Icons](https://lucide.dev)
