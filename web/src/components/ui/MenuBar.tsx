import { Fragment } from 'react';
import { Menu as HeadlessMenu, MenuButton, MenuItems, MenuItem, Transition } from '@headlessui/react';
import type { ReactNode } from 'react';
import { ChevronRight, Check } from 'lucide-react';
import './MenuBar.css';

export interface MenuBarItemData {
  key: string;
  label?: string;
  icon?: ReactNode;
  shortcut?: string;
  disabled?: boolean;
  danger?: boolean;
  checked?: boolean;
  children?: MenuBarItemData[];
  onClick?: () => void;
  divider?: boolean;
}

export interface MenuBarMenu {
  key: string;
  label: string;
  items: MenuBarItemData[];
}

export interface MenuBarProps {
  menus: MenuBarMenu[];
  className?: string;
  logo?: ReactNode;
  actions?: ReactNode;
}

interface MenuItemGroupProps {
  items: MenuBarItemData[];
  onClose?: () => void;
}

function MenuItemGroup({ items, onClose }: MenuItemGroupProps) {
  return (
    <>
      {items.map((item, index) => {
        if (item.divider) {
          return <div key={`divider-${index}`} className="ui-menubar__divider" />;
        }

        if (item.children && item.children.length > 0) {
          // Submenu
          return (
            <HeadlessMenu key={item.key} as="div" className="ui-menubar__submenu-wrapper">
              {({ open }) => (
                <>
                  <MenuButton
                    className={`ui-menubar__item ui-menubar__item--has-submenu ${
                      item.disabled ? 'ui-menubar__item--disabled' : ''
                    }`}
                    disabled={item.disabled}
                  >
                    {item.icon && <span className="ui-menubar__item-icon">{item.icon}</span>}
                    <span className="ui-menubar__item-label">{item.label}</span>
                    <ChevronRight className="ui-menubar__submenu-arrow" size={16} />
                  </MenuButton>

                  <Transition
                    show={open}
                    as={Fragment}
                    enter="ui-menubar__transition-enter"
                    enterFrom="ui-menubar__transition-enter-from"
                    enterTo="ui-menubar__transition-enter-to"
                    leave="ui-menubar__transition-leave"
                    leaveFrom="ui-menubar__transition-leave-from"
                    leaveTo="ui-menubar__transition-leave-to"
                  >
                    <MenuItems className="ui-menubar__dropdown ui-menubar__dropdown--submenu">
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
                className={`ui-menubar__item ${
                  focus ? 'ui-menubar__item--focus' : ''
                } ${item.disabled ? 'ui-menubar__item--disabled' : ''} ${
                  item.danger ? 'ui-menubar__item--danger' : ''
                }`}
                onClick={() => {
                  if (!item.disabled && item.onClick) {
                    item.onClick();
                    onClose?.();
                  }
                }}
                disabled={item.disabled}
              >
                {item.icon && <span className="ui-menubar__item-icon">{item.icon}</span>}
                <span className="ui-menubar__item-label">{item.label}</span>
                {item.checked && (
                  <Check className="ui-menubar__item-check" size={16} />
                )}
                {item.shortcut && (
                  <span className="ui-menubar__item-shortcut">{item.shortcut}</span>
                )}
              </button>
            )}
          </MenuItem>
        );
      })}
    </>
  );
}

export function MenuBar({ menus, className = '', logo, actions }: MenuBarProps) {
  return (
    <nav className={`ui-menubar ${className}`}>
      {/* Logo Section */}
      {logo && <div className="ui-menubar__logo">{logo}</div>}

      {/* Menu Items */}
      <div className="ui-menubar__menus">
        {menus.map((menu) => (
          <HeadlessMenu key={menu.key} as="div" className="ui-menubar__menu">
            {({ open, close }) => (
              <>
                <MenuButton
                  className={`ui-menubar__trigger ${open ? 'ui-menubar__trigger--active' : ''}`}
                >
                  {menu.label}
                </MenuButton>

                <Transition
                  as={Fragment}
                  enter="ui-menubar__transition-enter"
                  enterFrom="ui-menubar__transition-enter-from"
                  enterTo="ui-menubar__transition-enter-to"
                  leave="ui-menubar__transition-leave"
                  leaveFrom="ui-menubar__transition-leave-from"
                  leaveTo="ui-menubar__transition-leave-to"
                >
                  <MenuItems className="ui-menubar__dropdown">
                    <MenuItemGroup items={menu.items} onClose={close} />
                  </MenuItems>
                </Transition>
              </>
            )}
          </HeadlessMenu>
        ))}
      </div>

      {/* Actions Section */}
      {actions && <div className="ui-menubar__actions">{actions}</div>}
    </nav>
  );
}
