export type Theme =
  | 'modern-blue'
  | 'warm-sunset'
  | 'neo-mint'
  | 'slate-dark'
  | 'purple-dream'
  | 'ocean-breeze'
  | 'forest-green'
  | 'rose-gold'
  | 'midnight-purple'
  | 'sakura-pink'
  | 'cyber-neon';

export interface ThemeContextValue {
  theme: Theme;
  setTheme: (theme: Theme) => void;
}

export type ButtonVariant = 'primary' | 'secondary' | 'tertiary' | 'danger';
export type ButtonSize = 'small' | 'medium' | 'large';

export type InputStatus = 'default' | 'error' | 'success';

// Table Types
export interface PaginationState {
  pageIndex: number;
  pageSize: number;
}

export interface TablePaginationProps {
  currentPage: number;
  totalPages: number;
  pageSize: number;
  totalItems: number;
  pageSizeOptions?: number[];
  onPageChange: (page: number) => void;
  onPageSizeChange: (size: number) => void;
}
