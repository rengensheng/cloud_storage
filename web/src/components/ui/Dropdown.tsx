import { Menu, MenuItem, MenuButton, MenuItems } from '@headlessui/react';
import { motion, AnimatePresence } from 'framer-motion';
import type { ReactNode } from 'react';
import './Dropdown.css';

export interface DropdownItemProps {
  children: ReactNode;
  onClick?: () => void;
  disabled?: boolean;
  icon?: ReactNode;
}

export function DropdownItem({ children, onClick, disabled, icon }: DropdownItemProps) {
  return (
    <MenuItem disabled={disabled}>
      {({ focus }) => (
        <button
          onClick={onClick}
          className={`ui-dropdown-item ${focus ? 'ui-dropdown-item--active' : ''} ${disabled ? 'ui-dropdown-item--disabled' : ''}`}
        >
          {icon && <span className="ui-dropdown-item__icon">{icon}</span>}
          <span className="ui-dropdown-item__text">{children}</span>
        </button>
      )}
    </MenuItem>
  );
}

export interface DropdownProps {
  trigger: ReactNode;
  children: ReactNode;
  className?: string;
}

export function Dropdown({ trigger, children, className = '' }: DropdownProps) {
  return (
    <Menu as="div" className={`ui-dropdown ${className}`}>
      <MenuButton className="ui-dropdown-trigger">
        {trigger}
      </MenuButton>

      <AnimatePresence>
        <MenuItems className="ui-dropdown-items-wrapper">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.15 }}
            className="ui-dropdown-items"
          >
            {children}
          </motion.div>
        </MenuItems>
      </AnimatePresence>
    </Menu>
  );
}

Dropdown.Item = DropdownItem;
