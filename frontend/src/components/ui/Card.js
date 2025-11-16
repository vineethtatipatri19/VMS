import React from 'react';
import './Card.css';

/**
 * Card Component - Flexible container
 * Cross-platform compatible design
 */
const Card = ({ 
  children, 
  title,
  subtitle,
  headerAction,
  footer,
  variant = 'default',
  hoverable = false,
  padding = 'md',
  className = '',
  onClick,
  ...props 
}) => {
  const cardClasses = [
    'card',
    `card-${variant}`,
    `card-padding-${padding}`,
    hoverable && 'card-hoverable',
    onClick && 'card-clickable',
    className
  ].filter(Boolean).join(' ');

  return (
    <div className={cardClasses} onClick={onClick} {...props}>
      {(title || subtitle || headerAction) && (
        <div className="card-header">
          <div className="card-header-content">
            {title && <h3 className="card-title">{title}</h3>}
            {subtitle && <p className="card-subtitle">{subtitle}</p>}
          </div>
          {headerAction && <div className="card-header-action">{headerAction}</div>}
        </div>
      )}
      
      <div className="card-body">
        {children}
      </div>
      
      {footer && (
        <div className="card-footer">
          {footer}
        </div>
      )}
    </div>
  );
};

export default Card;
