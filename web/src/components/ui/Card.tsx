import type { ReactNode, CSSProperties } from 'react';
import './Card.css';

export interface CardProps {
  children: ReactNode;
  className?: string;
  padding?: 'none' | 'small' | 'medium' | 'large';
  shadow?: 'none' | 'small' | 'medium' | 'large';
  hover?: boolean;
  style?: CSSProperties;
}

export function Card({
  children,
  className = '',
  padding = 'medium',
  shadow = 'small',
  hover = false,
  style
}: CardProps) {
  const paddingClass = `ui-card--padding-${padding}`;
  const shadowClass = shadow !== 'none' ? `ui-card--shadow-${shadow}` : '';
  const hoverClass = hover ? 'ui-card--hover' : '';

  return (
    <div className={`ui-card ${paddingClass} ${shadowClass} ${hoverClass} ${className}`} style={style}>
      {children}
    </div>
  );
}
