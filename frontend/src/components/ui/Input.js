import React, { forwardRef } from 'react';
import './Input.css';

/**
 * Input Component - Form input with label and error handling
 */
const Input = forwardRef(({ 
  label,
  error,
  helperText,
  leftIcon,
  rightIcon,
  size = 'md',
  fullWidth = false,
  className = '',
  containerClassName = '',
  ...props 
}, ref) => {
  const inputClasses = [
    'input',
    `input-${size}`,
    error && 'input-error',
    leftIcon && 'input-with-left-icon',
    rightIcon && 'input-with-right-icon',
    className
  ].filter(Boolean).join(' ');

  const containerClasses = [
    'input-container',
    fullWidth && 'input-full-width',
    containerClassName
  ].filter(Boolean).join(' ');

  return (
    <div className={containerClasses}>
      {label && (
        <label className="input-label">
          {label}
          {props.required && <span className="input-required">*</span>}
        </label>
      )}
      
      <div className="input-wrapper">
        {leftIcon && <span className="input-icon-left">{leftIcon}</span>}
        
        <input
          ref={ref}
          className={inputClasses}
          {...props}
        />
        
        {rightIcon && <span className="input-icon-right">{rightIcon}</span>}
      </div>
      
      {(error || helperText) && (
        <span className={error ? 'input-error-text' : 'input-helper-text'}>
          {error || helperText}
        </span>
      )}
    </div>
  );
});

Input.displayName = 'Input';

export default Input;
