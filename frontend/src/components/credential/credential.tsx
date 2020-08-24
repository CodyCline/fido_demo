import * as React from 'react';
import { ReactComponent as Key } from '../../assets/key-solid.svg';
import { ReactComponent as Trash } from '../../assets/trash-solid.svg';
import './credential.css';


export const Credential = ({name, lastUsed, useCount} :any) => {
    return (
        <div className="credential">
            <div className="side-bar">
                <Key style={{height: "40px", color: "#CCC"}}/>
            </div>
            <div className="credential-info">
                <h3>{name || "Default Credential"}</h3>
                <p>Last used: {lastUsed || "Never"}</p>
                <p>Times used: {useCount || 0}</p>
            </div>
            <div className="delete-section">
                <Trash className="delete-btn"/>
            </div>
        </div>
    )
}