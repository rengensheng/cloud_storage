import { forwardRef } from 'react';
import type { InputHTMLAttributes } from 'react';
import type { InputStatus } from '../../types';
import './Input.css';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  status?: InputStatus;
  helperText?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      label,
      status = 'default',
      helperText,
      leftIcon,
      rightIcon,
      className = '',
      id,
      disabled,
      ...props
    },
    ref
  ) => {
    const inputId = id || `input-${Math.random().toString(36).substring(2, 11)}`;
    const baseClass = 'ui-input';
    const statusClass = status !== 'default' ? `ui-input--${status}` : '';
    const disabledClass = disabled ? 'ui-input--disabled' : '';
    const hasIconClass = (leftIcon || rightIcon) ? 'ui-input--has-icon' : '';

    return (
      <div className={`ui-input-wrapper ${className}`}>
        {label && (
          <label htmlFor={inputId} className="ui-input-label">
            {label}
          </label>
        )}
        <div className="ui-input-container">
          {leftIcon && (
            <span className="ui-input__icon ui-input__icon--left">
              {leftIcon}
            </span>
          )}
          <input
            ref={ref}
            id={inputId}
            className={[baseClass, statusClass, disabledClass, hasIconClass]
              .filter(Boolean)
              .join(' ')}
            disabled={disabled}
            {...props}
          />
          {rightIcon && (
            <span className="ui-input__icon ui-input__icon--right">
              {rightIcon}
            </span>
          )}
        </div>
        {helperText && (
          <p className={`ui-input-helper ${status !== 'default' ? `ui-input-helper--${status}` : ''}`}>
            {helperText}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
