import * as React from 'react';
import { axiosAuth } from '../utils/axios';


export const Credentials = () => {
    const [state, setState] = React.useState<any>({
        credentials: [],
        loaded: false,
    });
    React.useEffect(() => {
        async function getCreds() {
            try {
                const req = await axiosAuth.get("/api/credentials")
                console.log(req);
                setState({
                    ...state,
                    credentials: req.data.credentials,
                    loaded: true,
                });
            } catch (error) {

            }
        }
        getCreds();
    }, []);
    return (
        <div>
            Your Credentials
            {state.loaded ?
                state.credentials.map((cred: any, inc: number) => {
                    return (
                        <div style={{ border: "5px solid green" }} key={inc}>
                            <p>Nickname: {cred.name}</p>
                            <p>Last used: {cred.last_used}</p>
                            <p>Times used: {cred.counter}</p>
                        </div>
                    )
                })
                : <p>LOading</p>
            }
        </div>
    )
}