# 🎨 Modern UI Component Library

一个完整的、可切换主题的现代化 UI 组件库，基于 React、TypeScript、Framer Motion、HeadlessUI 和 TanStack Table 构建。

## ✨ 特性

- 🎨 **4 个精美主题** - Modern Blue、Warm Sunset、Neo Mint、Slate Dark
- 🔄 **完全可主题化** - 基于 CSS 变量的设计系统
- ♿ **无障碍访问** - WCAG AA 标准合规
- ⚡ **流畅动画** - Framer Motion 驱动的动效
- 📱 **响应式设计** - 移动端友好
- 🎯 **TypeScript** - 完整的类型支持
- 🎛️ **自定义类名** - 所有组件支持 className
- 📊 **功能完整的表格** - 基于 TanStack Table

## 📦 组件列表

### 基础表单组件
- ✅ **Button** - 多种变体、尺寸、加载状态、图标支持
- ✅ **Input** - 标签、图标、错误状态、帮助文本
- ✅ **Textarea** - 可调整大小、状态提示

### 选择组件
- ✅ **Checkbox** - 复选框
- ✅ **Radio Group** - 单选按钮组
- ✅ **Switch** - 开关切换
- ✅ **Select** - 下拉选择器

### 布局组件
- ✅ **Card** - 卡片容器，支持阴影、内边距、悬停效果
- ✅ **Container** - 响应式容器

### 复杂组件
- ✅ **Dialog** - 模态对话框
- ✅ **Dropdown** - 下拉菜单
- ✅ **Tabs** - 标签页

### 数据展示
- ✅ **Table** - 功能完整的数据表格
  - 排序
  - 分页
  - 行选择
  - 斑马条纹
  - 固定表头
  - 加载状态
  - 操作列
  - 状态徽章
- ✅ **Pagination** - 分页组件

## 🚀 快速开始

### 安装依赖

```bash
pnpm install
```

### 启动开发服务器

```bash
pnpm dev
```

### 构建生产版本

```bash
pnpm build
```

## 📖 使用文档

### 主题切换

```tsx
import { ThemeProvider, useTheme } from './contexts';

function App() {
  return (
    <ThemeProvider defaultTheme="modern-blue">
      <YourApp />
    </ThemeProvider>
  );
}

function ThemeSwitcher() {
  const { theme, setTheme } = useTheme();

  return (
    <Button onClick={() => setTheme('warm-sunset')}>
      切换主题
    </Button>
  );
}
```

### 基础组件使用

```tsx
import { Button, Input, Card } from './components/ui';
import { Mail } from 'lucide-react';

function MyComponent() {
  return (
    <Card padding="large" shadow="medium">
      <Input
        label="邮箱"
        type="email"
        leftIcon={<Mail />}
        placeholder="you@example.com"
      />
      <Button variant="primary" fullWidth>
        提交
      </Button>
    </Card>
  );
}
```

### 表格组件使用

```tsx
import { Table } from './components/ui';
import type { ColumnDef } from '@tanstack/react-table';

interface User {
  id: number;
  name: string;
  email: string;
}

const columns: ColumnDef<User>[] = [
  { accessorKey: 'id', header: 'ID' },
  { accessorKey: 'name', header: '姓名' },
  { accessorKey: 'email', header: '邮箱' }
];

function UserTable() {
  const [users, setUsers] = useState<User[]>([]);

  return (
    <Table
      data={users}
      columns={columns}
      striped={true}
      enablePagination={true}
      enableRowSelection={true}
      pageSize={10}
    />
  );
}
```

## 🎨 主题系统

### 可用主题

1. **Modern Blue** - 专业、干净的蓝色主题
2. **Warm Sunset** - 温暖、友好的橙红主题
3. **Neo Mint** - 清新、现代的青绿主题
4. **Slate Dark** - 深色模式，蓝色点缀

### 自定义主题

在 `src/styles/globals.css` 中添加：

```css
.theme-custom {
  --color-primary: #your-color;
  --color-primary-hover: #your-hover-color;
  --color-bg: #your-bg-color;
  /* ... 更多颜色变量 */
}
```

## 🎯 设计系统

### 色彩系统
- Primary - 主色
- Secondary - 辅助色
- Success - 成功状态
- Warning - 警告状态
- Danger - 危险状态
- Neutrals - 中性色

### 间距系统
基于 8pt 网格：
- 4px, 8px, 12px, 16px, 20px, 24px, 32px, 40px, 48px

### 字体层级
- H1: 40px, Bold
- H2: 28px, Semi-bold
- H3: 22px, Semi-bold
- Body L: 18px
- Body M: 16px
- Caption: 12px

### 圆角
- Small: 6px
- Medium: 10px
- Large: 12px
- XL: 16px

### 阴影
- Small: 0 1px 3px
- Medium: 0 4px 12px
- Large: 0 8px 24px
- XL: 0 16px 48px

## 📚 详细文档

- [UI 组件库完整文档](./UI_LIBRARY.md)
- [表格组件详细文档](./TABLE_DOCS.md)

## 🛠️ 技术栈

- **React 19** - UI 框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Framer Motion** - 动画库
- **HeadlessUI** - 无样式组件基础
- **TanStack Table** - 表格解决方案
- **Lucide React** - 图标库

## 🎨 组件演示

项目包含两个演示页面：

1. **组件展示页** - 展示所有基础 UI 组件
2. **表格演示页** - 展示表格的所有功能

运行项目后，点击右上角的切换按钮可以在两个页面之间切换。

## 📁 项目结构

```
src/
├── components/
│   └── ui/              # UI 组件库
│       ├── Button.tsx
│       ├── Input.tsx
│       ├── Table.tsx
│       └── ...
├── contexts/            # React Context
│   └── ThemeContext.tsx
├── pages/              # 页面
│   ├── ComponentShowcase.tsx
│   └── TableDemo.tsx
├── styles/             # 全局样式
│   └── globals.css
├── types/              # TypeScript 类型
│   └── index.ts
└── App.tsx
```

## ♿ 无障碍访问

所有组件都遵循 WCAG AA 标准：

- ✅ 适当的颜色对比度 (4.5:1)
- ✅ 键盘导航支持
- ✅ 屏幕阅读器兼容
- ✅ ARIA 标签
- ✅ 焦点指示器
- ✅ 最小触摸目标 (44x44px)

## 🌐 浏览器支持

- Chrome (最新)
- Firefox (最新)
- Safari (最新)
- Edge (最新)

## 📄 许可证

MIT

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 💡 设计理念

本组件库遵循以下设计原则：

1. **简约** - 去除不必要的装饰，以内容为核心
2. **一致性** - 统一的颜色、间距、动效
3. **可用性** - 易用、直观的交互
4. **可访问性** - 对所有用户友好
5. **性能** - 流畅的动画和渲染

## 🎓 参考资料

- Material Design
- Ant Design
- Chakra UI
- TailwindCSS

---

Built with ❤️ using modern web technologies
