import { RadioGroup as HeadlessRadioGroup } from '@headlessui/react';
import './Radio.css';

export interface RadioOption {
  value: string;
  label: string;
  disabled?: boolean;
}

export interface RadioGroupProps {
  value?: string;
  onChange?: (value: string) => void;
  options: RadioOption[];
  label?: string;
  className?: string;
}

export function RadioGroup({
  value,
  onChange,
  options,
  label,
  className = ''
}: RadioGroupProps) {
  return (
    <HeadlessRadioGroup value={value} onChange={onChange} className={`ui-radio-group ${className}`}>
      {label && <HeadlessRadioGroup.Label className="ui-radio-group-label">{label}</HeadlessRadioGroup.Label>}
      <div className="ui-radio-group-options">
        {options.map((option) => (
          <HeadlessRadioGroup.Option
            key={option.value}
            value={option.value}
            disabled={option.disabled}
            className="ui-radio-wrapper"
          >
            {({ checked, disabled }) => (
              <>
                <span className={`ui-radio ${checked ? 'ui-radio--checked' : ''} ${disabled ? 'ui-radio--disabled' : ''}`}>
                  {checked && <span className="ui-radio__dot" />}
                </span>
                <span className="ui-radio-label">{option.label}</span>
              </>
            )}
          </HeadlessRadioGroup.Option>
        ))}
      </div>
    </HeadlessRadioGroup>
  );
}
