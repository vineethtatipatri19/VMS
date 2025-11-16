import React from 'react';
import './Button.css';

/**
 * Button Component - Cross-platform compatible
 * Designed to work on web and easily adaptable to React Native
 */
const Button = ({ 
  children, 
  variant = 'primary', 
  size = 'md', 
  fullWidth = false,
  disabled = false,
  loading = false,
  leftIcon = null,
  rightIcon = null,
  onClick,
  type = 'button',
  className = '',
  ...props 
}) => {
  const buttonClasses = [
    'btn',
    `btn-${variant}`,
    `btn-${size}`,
    fullWidth && 'btn-full-width',
    disabled && 'btn-disabled',
    loading && 'btn-loading',
    className
  ].filter(Boolean).join(' ');

  return (
    <button
      type={type}
      className={buttonClasses}
      onClick={onClick}
      disabled={disabled || loading}
      {...props}
    >
      {loading && <span className="btn-spinner" />}
      {!loading && leftIcon && <span className="btn-icon-left">{leftIcon}</span>}
      <span className="btn-text">{children}</span>
      {!loading && rightIcon && <span className="btn-icon-right">{rightIcon}</span>}
    </button>
  );
};

export default Button;
