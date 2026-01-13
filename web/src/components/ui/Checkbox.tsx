import { Checkbox as HeadlessCheckbox } from '@headlessui/react';
import { Check } from 'lucide-react';
import './Checkbox.css';

export interface CheckboxProps {
  checked?: boolean;
  onChange?: (checked: boolean) => void;
  label?: string;
  disabled?: boolean;
  className?: string;
}

export function Checkbox({
  checked = false,
  onChange,
  label,
  disabled = false,
  className = ''
}: CheckboxProps) {
  return (
    <HeadlessCheckbox
      checked={checked}
      onChange={onChange}
      disabled={disabled}
      className={`ui-checkbox-wrapper ${className}`}
    >
      <span className={`ui-checkbox ${checked ? 'ui-checkbox--checked' : ''} ${disabled ? 'ui-checkbox--disabled' : ''}`}>
        {checked && <Check className="ui-checkbox__icon" strokeWidth={3} />}
      </span>
      {label && <span className="ui-checkbox-label">{label}</span>}
    </HeadlessCheckbox>
  );
}
