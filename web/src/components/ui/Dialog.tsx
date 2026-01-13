import { Dialog as HeadlessDialog } from '@headlessui/react';
import { X } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import type { ReactNode } from 'react';
import './Dialog.css';

export interface DialogProps {
  open: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  className?: string;
  size?: 'small' | 'medium' | 'large';
}

export function Dialog({
  open,
  onClose,
  title,
  children,
  className = '',
  size = 'medium'
}: DialogProps) {
  return (
    <AnimatePresence>
      {open && (
        <HeadlessDialog
          static
          open={open}
          onClose={onClose}
          className="ui-dialog-container"
        >
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="ui-dialog-backdrop"
            aria-hidden="true"
          />

          {/* Dialog Panel */}
          <div className="ui-dialog-wrapper">
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              transition={{ duration: 0.2 }}
              className="ui-dialog-content-wrapper"
            >
              <HeadlessDialog.Panel className={`ui-dialog-panel ui-dialog-panel--${size} ${className}`}>
                {/* Header */}
                {title && (
                  <div className="ui-dialog-header">
                    <HeadlessDialog.Title className="ui-dialog-title">
                      {title}
                    </HeadlessDialog.Title>
                    <button
                      onClick={onClose}
                      className="ui-dialog-close"
                      aria-label="Close dialog"
                    >
                      <X />
                    </button>
                  </div>
                )}

                {/* Body */}
                <div className="ui-dialog-body">
                  {children}
                </div>
              </HeadlessDialog.Panel>
            </motion.div>
          </div>
        </HeadlessDialog>
      )}
    </AnimatePresence>
  );
}
