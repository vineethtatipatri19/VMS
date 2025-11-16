import React from 'react';
import './Badge.css';

/**
 * Badge Component - Status indicators
 */
const Badge = ({ 
  children, 
  variant = 'default',
  size = 'md',
  rounded = false,
  dot = false,
  className = '',
  ...props 
}) => {
  const badgeClasses = [
    'badge',
    `badge-${variant}`,
    `badge-${size}`,
    rounded && 'badge-rounded',
    dot && 'badge-dot',
    className
  ].filter(Boolean).join(' ');

  return (
    <span className={badgeClasses} {...props}>
      {dot && <span className="badge-dot-indicator" />}
      {children}
    </span>
  );
};

export default Badge;
