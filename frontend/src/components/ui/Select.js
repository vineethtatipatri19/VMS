import React, { forwardRef } from 'react';
import { ChevronDown } from 'lucide-react';
import './Select.css';

/**
 * Select Component - Dropdown selector
 */
const Select = forwardRef(({ 
  label,
  error,
  helperText,
  options = [],
  placeholder = 'Select...',
  size = 'md',
  fullWidth = false,
  className = '',
  containerClassName = '',
  ...props 
}, ref) => {
  const selectClasses = [
    'select',
    `select-${size}`,
    error && 'select-error',
    className
  ].filter(Boolean).join(' ');

  const containerClasses = [
    'select-container',
    fullWidth && 'select-full-width',
    containerClassName
  ].filter(Boolean).join(' ');

  return (
    <div className={containerClasses}>
      {label && (
        <label className="select-label">
          {label}
          {props.required && <span className="select-required">*</span>}
        </label>
      )}
      
      <div className="select-wrapper">
        <select
          ref={ref}
          className={selectClasses}
          {...props}
        >
          {placeholder && (
            <option value="" disabled>
              {placeholder}
            </option>
          )}
          {options.map((option, index) => (
            <option 
              key={option.value || index} 
              value={option.value}
              disabled={option.disabled}
            >
              {option.label}
            </option>
          ))}
        </select>
        
        <ChevronDown className="select-icon" size={20} />
      </div>
      
      {(error || helperText) && (
        <span className={error ? 'select-error-text' : 'select-helper-text'}>
          {error || helperText}
        </span>
      )}
    </div>
  );
});

Select.displayName = 'Select';

export default Select;
