import { forwardRef } from 'react';
import type { TextareaHTMLAttributes } from 'react';
import type { InputStatus } from '../../types';
import './Textarea.css';

export interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  status?: InputStatus;
  helperText?: string;
  resize?: 'none' | 'vertical' | 'horizontal' | 'both';
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  (
    {
      label,
      status = 'default',
      helperText,
      resize = 'vertical',
      className = '',
      id,
      disabled,
      ...props
    },
    ref
  ) => {
    const textareaId = id || `textarea-${Math.random().toString(36).substring(2, 11)}`;
    const baseClass = 'ui-textarea';
    const statusClass = status !== 'default' ? `ui-textarea--${status}` : '';
    const disabledClass = disabled ? 'ui-textarea--disabled' : '';
    const resizeClass = `ui-textarea--resize-${resize}`;

    return (
      <div className={`ui-textarea-wrapper ${className}`}>
        {label && (
          <label htmlFor={textareaId} className="ui-textarea-label">
            {label}
          </label>
        )}
        <textarea
          ref={ref}
          id={textareaId}
          className={[baseClass, statusClass, disabledClass, resizeClass]
            .filter(Boolean)
            .join(' ')}
          disabled={disabled}
          {...props}
        />
        {helperText && (
          <p className={`ui-textarea-helper ${status !== 'default' ? `ui-textarea-helper--${status}` : ''}`}>
            {helperText}
          </p>
        )}
      </div>
    );
  }
);

Textarea.displayName = 'Textarea';
