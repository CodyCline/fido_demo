import * as React from 'react';
import { Route, Redirect, useHistory } from 'react-router-dom';
import Cookies from 'universal-cookie'



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
