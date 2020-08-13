import * as React from 'react';

export const Input = ({value, onChange, type, validationText} : any) => {
    return (
        <div className="input-block">
            <input className="input" type={type} value={value} onChange={onChange}/>
            {validationText && <span>{validationText}</span>}
        </div>
    )
}