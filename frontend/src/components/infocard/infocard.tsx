import * as React from 'react';

export const InfoCard = ({children, header, subText} :any) => {
    return (
        <div className="infocard">
            <h2 className="info-header">{header}</h2>
            <p className="info-text">{subText}</p>
            {children}
        </div>
    )
}