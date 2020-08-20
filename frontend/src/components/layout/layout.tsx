import * as React from 'react';
import './layout.css';

export const Layout = ({children, isAuth} : any) => {
    return (
        <div className="layout">
            <nav className="navbar">
                <h1 style={{marginBlockStart: 0, marginBlockEnd: 0}}>WebAuthn Demo</h1>
            </nav>
            {children}
            <div style={{height: "300px"}}/>
            <footer>
                footer
            </footer>
        </div>
    )
}