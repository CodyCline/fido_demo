import * as React from 'react';
import './input.css';

export const Input = ({value, onChange, type, validationText, label, placeHolder, name} : any) => {
    return (
        <div className="input-block">
            <label htmlFor={name} className="input-label">{label}</label>
            <input placeholder={placeHolder} className="input" name={name} type={type} value={value} onChange={onChange}/>
            {validationText && <span className="validation-text">{validationText}</span>}
        </div>
    )
}