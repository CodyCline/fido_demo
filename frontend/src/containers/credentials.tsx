import * as React from 'react';
import { axiosAuth } from '../utils/axios';


export const Credentials = () => {
    const [state, setState] = React.useState<any>();
    React.useEffect(() => {
        async function getCreds() {
            try {
                const req = await axiosAuth.get("/credentials")
                setState(req.data.credentials);
            } catch (error) {

            }   
        }
        getCreds();
    }, []);
    return (
        <div>
            Your Credentials
            {state.credentials && state.credentials.map((cred: any, inc:number) => {
                return (
                    <div style={{border: "5px solid green"}} key={inc}>
                    <p>Nickname: {cred.name}</p>
                    <p>Last used: {cred.last_used}</p>
                    <p>Times used: {cred.counter}</p>
                </div>
                )
            })}
        </div>
    )
}