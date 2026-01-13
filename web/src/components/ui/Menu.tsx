import { Fragment } from 'react';
import { Menu as HeadlessMenu, MenuButton, MenuItems, MenuItem, Transition } from '@headlessui/react';
import type { ReactNode } from 'react';
import { ChevronRight, Check } from 'lucide-react';
import './Menu.css';

export interface MenuItemData {
  key: string;
  label?: string;
  icon?: ReactNode;
  shortcut?: string;
  disabled?: boolean;
  danger?: boolean;
  checked?: boolean;
  children?: MenuItemData[];
  onClick?: () => void;
  divider?: boolean;
}

export interface MenuProps {
  trigger: ReactNode;
  items: MenuItemData[];
  className?: string;
  position?: 'bottom-start' | 'bottom-end' | 'top-start' | 'top-end';
}

interface MenuItemGroupProps {
  items: MenuItemData[];
  onClose?: () => void;
}

function MenuItemGroup({ items, onClose }: MenuItemGroupProps) {
  return (
    <>
      {items.map((item, index) => {
        if (item.divider) {
          return <div key={`divider-${index}`} className="ui-menu__divider" />;
        }

        if (item.children && item.children.length > 0) {
          // Submenu
          return (
            <HeadlessMenu key={item.key} as="div" className="ui-menu__submenu-wrapper">
              {({ open }) => (
                <>
                  <MenuButton
                    className={`ui-menu__item ui-menu__item--has-submenu ${
                      item.disabled ? 'ui-menu__item--disabled' : ''
                    }`}
                    disabled={item.disabled}
                  >
                    {item.icon && <span className="ui-menu__item-icon">{item.icon}</span>}
                    <span className="ui-menu__item-label">{item.label}</span>
                    <ChevronRight className="ui-menu__submenu-arrow" size={16} />
                  </MenuButton>

                  <Transition
                    show={open}
                    as={Fragment}
                    enter="ui-menu__transition-enter"
                    enterFrom="ui-menu__transition-enter-from"
                    enterTo="ui-menu__transition-enter-to"
                    leave="ui-menu__transition-leave"
                    leaveFrom="ui-menu__transition-leave-from"
                    leaveTo="ui-menu__transition-leave-to"
                  >
                    <MenuItems className="ui-menu__items ui-menu__items--submenu">
                      <MenuItemGroup items={item.children || []} onClose={onClose} />
                    </MenuItems>
                  </Transition>
                </>
              )}
            </HeadlessMenu>
          );
        }

        // Regular menu item
        return (
          <MenuItem key={item.key} disabled={item.disabled}>
            {({ focus }) => (
              <button
                className={`ui-menu__item ${
                  focus ? 'ui-menu__item--focus' : ''
                } ${item.disabled ? 'ui-menu__item--disabled' : ''} ${
                  item.danger ? 'ui-menu__item--danger' : ''
                }`}
                onClick={() => {
                  if (!item.disabled && item.onClick) {
                    item.onClick();
                    onClose?.();
                  }
                }}
                disabled={item.disabled}
              >
                {item.icon && <span className="ui-menu__item-icon">{item.icon}</span>}
                <span className="ui-menu__item-label">{item.label}</span>
                {item.checked && (
                  <Check className="ui-menu__item-check" size={16} />
                )}
                {item.shortcut && (
                  <span className="ui-menu__item-shortcut">{item.shortcut}</span>
                )}
              </button>
            )}
          </MenuItem>
        );
      })}
    </>
  );
}

export function Menu({ trigger, items, className = '', position = 'bottom-start' }: MenuProps) {
  const positionClasses = {
    'bottom-start': 'ui-menu--bottom-start',
    'bottom-end': 'ui-menu--bottom-end',
    'top-start': 'ui-menu--top-start',
    'top-end': 'ui-menu--top-end'
  };

  return (
    <HeadlessMenu as="div" className={`ui-menu ${className}`}>
      {({ close }) => (
        <>
          <MenuButton as={Fragment}>{trigger}</MenuButton>

          <Transition
            as={Fragment}
            enter="ui-menu__transition-enter"
            enterFrom="ui-menu__transition-enter-from"
            enterTo="ui-menu__transition-enter-to"
            leave="ui-menu__transition-leave"
            leaveFrom="ui-menu__transition-leave-from"
            leaveTo="ui-menu__transition-leave-to"
          >
            <MenuItems className={`ui-menu__items ${positionClasses[position]}`}>
              <MenuItemGroup items={items} onClose={close} />
            </MenuItems>
          </Transition>
        </>
      )}
    </HeadlessMenu>
  );
}

// Compound components for flexible usage
Menu.Item = MenuItem;
Menu.Items = MenuItems;
Menu.Button = MenuButton;
