import { Listbox } from '@headlessui/react';
import { ChevronDown, Check } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import './Select.css';

export interface SelectOption {
  value: string;
  label: string;
  disabled?: boolean;
}

export interface SelectProps {
  value?: string;
  onChange?: (value: string) => void;
  options: SelectOption[];
  placeholder?: string;
  label?: string;
  disabled?: boolean;
  className?: string;
}

export function Select({
  value,
  onChange,
  options,
  placeholder = 'Select an option',
  label,
  disabled = false,
  className = ''
}: SelectProps) {
  const selectedOption = options.find(opt => opt.value === value);

  return (
    <Listbox value={value} onChange={onChange} disabled={disabled}>
      <div className={`ui-select-wrapper ${className}`}>
        {label && <Listbox.Label className="ui-select-label">{label}</Listbox.Label>}
        <Listbox.Button className={`ui-select-button ${disabled ? 'ui-select-button--disabled' : ''}`}>
          <span className="ui-select-button__text">
            {selectedOption ? selectedOption.label : placeholder}
          </span>
          <ChevronDown className="ui-select-button__icon" />
        </Listbox.Button>

        <AnimatePresence>
          <Listbox.Options as="div" className="ui-select-options-wrapper">
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              transition={{ duration: 0.15 }}
              className="ui-select-options"
            >
              {options.map((option) => (
                <Listbox.Option
                  key={option.value}
                  value={option.value}
                  disabled={option.disabled}
                  className="ui-select-option"
                >
                  {({ selected }) => (
                    <>
                      <span className={`ui-select-option__text ${selected ? 'ui-select-option__text--selected' : ''}`}>
                        {option.label}
                      </span>
                      {selected && <Check className="ui-select-option__check" />}
                    </>
                  )}
                </Listbox.Option>
              ))}
            </motion.div>
          </Listbox.Options>
        </AnimatePresence>
      </div>
    </Listbox>
  );
}
