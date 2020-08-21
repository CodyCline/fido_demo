import * as React from 'react';
import { Redirect } from 'react-router-dom';
import Cookies from 'universal-cookie';



export const AuthRequired = ({ isAuthenticated, children }: any) => {
    const cookies = new Cookies();
    const [isAuth, setAuth] = React.useState(cookies.get("token"));
    React.useEffect(() => {
        if( cookies.get("token") ) {
            // setAuth(true);
        }
        
        console.log("HELLO")
    });

    if (isAuth) {
        return (
            <React.Fragment>
                {children}
            </React.Fragment>
        )
    }
    return (
        <Redirect to="/login" />
    )
}
