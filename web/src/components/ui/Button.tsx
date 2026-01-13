import { forwardRef, type MouseEvent } from 'react';
import type { ButtonHTMLAttributes } from 'react';
import type { ButtonVariant, ButtonSize } from '../../types';
import './Button.css';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  fullWidth?: boolean;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  loading?: boolean;
  disableRipple?: boolean;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      variant = 'primary',
      size = 'medium',
      fullWidth = false,
      leftIcon,
      rightIcon,
      loading = false,
      disabled,
      disableRipple = false,
      className = '',
      children,
      onClick,
      ...props
    },
    ref
  ) => {
    const baseClass = 'ui-button';
    const variantClass = `ui-button--${variant}`;
    const sizeClass = `ui-button--${size}`;
    const fullWidthClass = fullWidth ? 'ui-button--full-width' : '';
    const loadingClass = loading ? 'ui-button--loading' : '';

    const classes = [
      baseClass,
      variantClass,
      sizeClass,
      fullWidthClass,
      loadingClass,
      className
    ].filter(Boolean).join(' ');

    const createRipple = (event: MouseEvent<HTMLButtonElement>) => {
      const button = event.currentTarget;

      // Remove existing ripples
      const existingRipples = button.querySelectorAll('.ui-button__ripple');
      existingRipples.forEach(ripple => ripple.remove());

      // Create ripple element
      const ripple = document.createElement('span');
      ripple.classList.add('ui-button__ripple');

      // Calculate ripple size and position
      const rect = button.getBoundingClientRect();
      const size = Math.max(rect.width, rect.height);
      const x = event.clientX - rect.left - size / 2;
      const y = event.clientY - rect.top - size / 2;

      ripple.style.width = ripple.style.height = `${size}px`;
      ripple.style.left = `${x}px`;
      ripple.style.top = `${y}px`;

      button.appendChild(ripple);

      // Remove ripple after animation
      setTimeout(() => {
        ripple.remove();
      }, 600);
    };

    const handleClick = (event: MouseEvent<HTMLButtonElement>) => {
      if (!disableRipple && !disabled && !loading) {
        createRipple(event);
      }
      onClick?.(event);
    };

    return (
      <button
        ref={ref}
        className={classes}
        disabled={disabled || loading}
        onClick={handleClick}
        {...props}
      >
        {loading && (
          <span className="ui-button__spinner" aria-hidden="true">
            <svg className="ui-button__spinner-icon" viewBox="0 0 24 24">
              <circle
                className="ui-button__spinner-circle"
                cx="12"
                cy="12"
                r="10"
                fill="none"
                strokeWidth="3"
              />
            </svg>
          </span>
        )}
        {!loading && leftIcon && (
          <span className="ui-button__icon ui-button__icon--left">
            {leftIcon}
          </span>
        )}
        <span className="ui-button__text">{children}</span>
        {!loading && rightIcon && (
          <span className="ui-button__icon ui-button__icon--right">
            {rightIcon}
          </span>
        )}
      </button>
    );
  }
);

Button.displayName = 'Button';
