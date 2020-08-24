import * as React from 'react';
import './button.css';

export const Button = ({disabled, onClick, style, children} :any) => {
    return (
        <button 
            style={style}
            onClick={onClick} 
            className="button" 
            disabled={disabled}
        >
            {children}
        </button>
    )
}