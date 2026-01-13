import type { ReactNode, CSSProperties } from 'react';
import './Container.css';

export interface ContainerProps {
  children: ReactNode;
  className?: string;
  maxWidth?: 'small' | 'medium' | 'large' | 'full';
  padding?: boolean;
  center?: boolean;
  style?: CSSProperties;
}

export function Container({
  children,
  className = '',
  maxWidth = 'medium',
  padding = true,
  center = true,
  style
}: ContainerProps) {
  const maxWidthClass = maxWidth !== 'full' ? `ui-container--${maxWidth}` : '';
  const paddingClass = padding ? 'ui-container--padding' : '';
  const centerClass = center ? 'ui-container--center' : '';

  return (
    <div
      className={`ui-container ${maxWidthClass} ${paddingClass} ${centerClass} ${className}`}
      style={style}
    >
      {children}
    </div>
  );
}
