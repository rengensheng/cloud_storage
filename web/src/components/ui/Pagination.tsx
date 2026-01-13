import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-react';
import type { TablePaginationProps } from '../../types';
import { Button } from './Button';
import { Select } from './Select';
import './Pagination.css';

export function Pagination({
  currentPage,
  totalPages,
  pageSize,
  totalItems,
  pageSizeOptions = [10, 20, 50, 100],
  onPageChange,
  onPageSizeChange
}: TablePaginationProps) {
  const startItem = totalItems === 0 ? 0 : (currentPage - 1) * pageSize + 1;
  const endItem = Math.min(currentPage * pageSize, totalItems);

  const canGoPrevious = currentPage > 1;
  const canGoNext = currentPage < totalPages;

  // Generate page numbers to display
  const getPageNumbers = () => {
    const pages: (number | string)[] = [];
    const maxVisible = 7; // Maximum number of page buttons to show

    if (totalPages <= maxVisible) {
      // Show all pages if total is less than max
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      // Always show first page
      pages.push(1);

      if (currentPage > 3) {
        pages.push('...');
      }

      // Show pages around current page
      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (currentPage < totalPages - 2) {
        pages.push('...');
      }

      // Always show last page
      pages.push(totalPages);
    }

    return pages;
  };

  const pageSizeSelectOptions = pageSizeOptions.map(size => ({
    value: size.toString(),
    label: `${size} / page`
  }));

  return (
    <div className="ui-pagination">
      {/* Items info */}
      <div className="ui-pagination__info">
        <span className="text-body-s">
          Showing <strong>{startItem}</strong> to <strong>{endItem}</strong> of{' '}
          <strong>{totalItems}</strong> items
        </span>
      </div>

      {/* Page size selector */}
      <div className="ui-pagination__page-size">
        <Select
          value={pageSize.toString()}
          onChange={(value) => onPageSizeChange(Number(value))}
          options={pageSizeSelectOptions}
        />
      </div>

      {/* Page navigation */}
      <div className="ui-pagination__navigation">
        {/* First page */}
        <Button
          variant="tertiary"
          size="small"
          onClick={() => onPageChange(1)}
          disabled={!canGoPrevious}
          aria-label="Go to first page"
        >
          <ChevronsLeft size={16} />
        </Button>

        {/* Previous page */}
        <Button
          variant="tertiary"
          size="small"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={!canGoPrevious}
          aria-label="Go to previous page"
        >
          <ChevronLeft size={16} />
        </Button>

        {/* Page numbers */}
        <div className="ui-pagination__pages">
          {getPageNumbers().map((page, index) => {
            if (page === '...') {
              return (
                <span key={`ellipsis-${index}`} className="ui-pagination__ellipsis">
                  ...
                </span>
              );
            }

            return (
              <button
                key={page}
                onClick={() => onPageChange(page as number)}
                className={`ui-pagination__page ${
                  currentPage === page ? 'ui-pagination__page--active' : ''
                }`}
                aria-label={`Go to page ${page}`}
                aria-current={currentPage === page ? 'page' : undefined}
              >
                {page}
              </button>
            );
          })}
        </div>

        {/* Next page */}
        <Button
          variant="tertiary"
          size="small"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={!canGoNext}
          aria-label="Go to next page"
        >
          <ChevronRight size={16} />
        </Button>

        {/* Last page */}
        <Button
          variant="tertiary"
          size="small"
          onClick={() => onPageChange(totalPages)}
          disabled={!canGoNext}
          aria-label="Go to last page"
        >
          <ChevronsRight size={16} />
        </Button>
      </div>
    </div>
  );
}
