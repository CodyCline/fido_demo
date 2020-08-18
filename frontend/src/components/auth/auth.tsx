import * as React from 'react';
import { Redirect } from 'react-router-dom';



export const AuthRequired = ({isAuthenticated, children} : any) => {
    if (isAuthenticated) {
        return (
            <React.Fragment>
                {children}
            </React.Fragment>
        )
    }
    return (
        <Redirect to="/login"/>
    )
}
