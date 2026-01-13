import { Switch as HeadlessSwitch } from '@headlessui/react';
import './Switch.css';

export interface SwitchProps {
  checked?: boolean;
  onChange?: (checked: boolean) => void;
  label?: string;
  disabled?: boolean;
  className?: string;
}

export function Switch({
  checked = false,
  onChange,
  label,
  disabled = false,
  className = ''
}: SwitchProps) {
  return (
    <HeadlessSwitch.Group as="div" className={`ui-switch-group ${className}`}>
      <HeadlessSwitch
        checked={checked}
        onChange={onChange}
        disabled={disabled}
        className={`ui-switch ${checked ? 'ui-switch--checked' : ''} ${disabled ? 'ui-switch--disabled' : ''}`}
      >
        <span className="ui-switch__thumb" />
      </HeadlessSwitch>
      {label && (
        <HeadlessSwitch.Label className="ui-switch-label">
          {label}
        </HeadlessSwitch.Label>
      )}
    </HeadlessSwitch.Group>
  );
}
