import { useState } from 'react';
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getPaginationRowModel,
  getFilteredRowModel,
  flexRender
} from '@tanstack/react-table';
import type { SortingState, ColumnDef, Row } from '@tanstack/react-table';
import { ArrowUp, ArrowDown, ChevronsUpDown, Loader2 } from 'lucide-react';
import { Checkbox } from './Checkbox';
import { Pagination } from './Pagination';
import './Table.css';

export interface TableProps<TData> {
  data: TData[];
  columns: ColumnDef<TData, any>[];
  loading?: boolean;
  striped?: boolean;
  hoverable?: boolean;
  bordered?: boolean;
  stickyHeader?: boolean;
  stickyActions?: boolean;
  stickySelection?: boolean;
  enableSorting?: boolean;
  enablePagination?: boolean;
  enableRowSelection?: boolean;
  pageSize?: number;
  pageSizeOptions?: number[];
  onRowClick?: (row: Row<TData>) => void;
  onRowSelectionChange?: (selectedRows: TData[]) => void;
  emptyMessage?: string;
  className?: string;
  maxHeight?: string;
}

export function Table<TData>({
  data,
  columns,
  loading = false,
  striped = false,
  hoverable = true,
  bordered = false,
  stickyHeader = false,
  stickyActions = false,
  stickySelection = false,
  enableSorting = true,
  enablePagination = true,
  enableRowSelection = false,
  pageSize: initialPageSize = 10,
  pageSizeOptions = [10, 20, 50, 100],
  onRowClick,
  onRowSelectionChange,
  emptyMessage = 'No data available',
  className = '',
  maxHeight
}: TableProps<TData>) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [rowSelection, setRowSelection] = useState({});
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: initialPageSize
  });

  // Build columns with selection if enabled
  const tableColumns = enableRowSelection
    ? [
        {
          id: 'select',
          header: ({ table }: any) => (
            <Checkbox
              checked={table.getIsAllRowsSelected()}
              onChange={(checked) => {
                table.toggleAllRowsSelected(checked);
              }}
            />
          ),
          cell: ({ row }: any) => (
            <Checkbox
              checked={row.getIsSelected()}
              onChange={(checked) => {
                row.toggleSelected(checked);
              }}
            />
          ),
          size: 40
        } as ColumnDef<TData, any>,
        ...columns
      ]
    : columns;

  const table = useReactTable({
    data,
    columns: tableColumns,
    state: {
      sorting,
      rowSelection,
      pagination
    },
    enableSorting,
    enableRowSelection,
    onSortingChange: setSorting,
    onRowSelectionChange: (updater) => {
      setRowSelection(updater);
      if (onRowSelectionChange) {
        const selection = typeof updater === 'function' ? updater(rowSelection) : updater;
        const selectedRows = Object.keys(selection)
          .filter(key => selection[key as keyof typeof selection])
          .map(key => data[parseInt(key)]);
        onRowSelectionChange(selectedRows);
      }
    },
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: enablePagination ? getPaginationRowModel() : undefined,
    getFilteredRowModel: getFilteredRowModel(),
    manualPagination: false
  });

  const containerClasses = [
    'ui-table-container',
    bordered && 'ui-table-container--bordered',
    stickyHeader && 'ui-table-container--sticky-header',
    className
  ]
    .filter(Boolean)
    .join(' ');

  const tableClasses = [
    'ui-table',
    striped && 'ui-table--striped',
    hoverable && 'ui-table--hoverable',
    loading && 'ui-table--loading'
  ]
    .filter(Boolean)
    .join(' ');

  const containerStyle = maxHeight ? { maxHeight } : undefined;

  return (
    <div className="ui-table-wrapper">
      <div className={containerClasses} style={containerStyle}>
        <table className={tableClasses}>
          <thead className="ui-table__thead">
            {table.getHeaderGroups().map(headerGroup => (
              <tr key={headerGroup.id} className="ui-table__tr">
                {headerGroup.headers.map((header, index) => {
                  const canSort = header.column.getCanSort();
                  const isSorted = header.column.getIsSorted();
                  const isLastColumn = index === headerGroup.headers.length - 1;
                  const isFirstColumn = index === 0;
                  const shouldStickRight = stickyActions && isLastColumn;
                  const shouldStickLeft = stickySelection && isFirstColumn && enableRowSelection;

                  return (
                    <th
                      key={header.id}
                      className={`ui-table__th ${shouldStickRight ? 'ui-table__th--sticky-right' : ''} ${shouldStickLeft ? 'ui-table__th--sticky-left' : ''}`}
                      style={{ width: header.getSize() !== 150 ? header.getSize() : undefined }}
                    >
                      {header.isPlaceholder ? null : (
                        <div
                          className={`ui-table__th-content ${
                            canSort ? 'ui-table__th-content--sortable' : ''
                          }`}
                          onClick={canSort ? header.column.getToggleSortingHandler() : undefined}
                        >
                          {flexRender(header.column.columnDef.header, header.getContext())}
                          {canSort && (
                            <span className="ui-table__sort-icon">
                              {isSorted === 'asc' ? (
                                <ArrowUp size={16} />
                              ) : isSorted === 'desc' ? (
                                <ArrowDown size={16} />
                              ) : (
                                <ChevronsUpDown size={16} />
                              )}
                            </span>
                          )}
                        </div>
                      )}
                    </th>
                  );
                })}
              </tr>
            ))}
          </thead>
          <tbody className="ui-table__tbody">
            {loading ? (
              <tr className="ui-table__tr">
                <td colSpan={tableColumns.length} className="ui-table__td ui-table__loading">
                  <div className="ui-table__loading-content">
                    <Loader2 className="ui-table__spinner" size={32} />
                    <p>Loading...</p>
                  </div>
                </td>
              </tr>
            ) : table.getRowModel().rows.length === 0 ? (
              <tr className="ui-table__tr">
                <td colSpan={tableColumns.length} className="ui-table__td ui-table__empty">
                  <div className="ui-table__empty-content">
                    <p>{emptyMessage}</p>
                  </div>
                </td>
              </tr>
            ) : (
              table.getRowModel().rows.map(row => (
                <tr
                  key={row.id}
                  className={`ui-table__tr ${
                    onRowClick ? 'ui-table__tr--clickable' : ''
                  } ${row.getIsSelected() ? 'ui-table__tr--selected' : ''}`}
                  onClick={() => onRowClick?.(row)}
                >
                  {row.getVisibleCells().map((cell, index) => {
                    const isLastColumn = index === row.getVisibleCells().length - 1;
                    const isFirstColumn = index === 0;
                    const shouldStickRight = stickyActions && isLastColumn;
                    const shouldStickLeft = stickySelection && isFirstColumn && enableRowSelection;

                    return (
                      <td key={cell.id} className={`ui-table__td ${shouldStickRight ? 'ui-table__td--sticky-right' : ''} ${shouldStickLeft ? 'ui-table__td--sticky-left' : ''}`}>
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </td>
                    );
                  })}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {enablePagination && !loading && data.length > 0 && (
        <Pagination
          currentPage={table.getState().pagination.pageIndex + 1}
          totalPages={table.getPageCount()}
          pageSize={table.getState().pagination.pageSize}
          totalItems={data.length}
          pageSizeOptions={pageSizeOptions}
          onPageChange={(page) => table.setPageIndex(page - 1)}
          onPageSizeChange={(size) => table.setPageSize(size)}
        />
      )}
    </div>
  );
}
