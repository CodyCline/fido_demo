import * as React from 'react';
import { axiosAuth } from '../utils/axios';
import { Credential } from '../components/credential/credential';

export const Credentials = () => {
    const [state, setState] = React.useState<any>({
        credentials: [],
        loaded: false,
    });
    React.useEffect(() => {
        async function getCreds() {
            try {
                const req = await axiosAuth.get("/api/credentials")
                setState({
                    ...state,
                    credentials: req.data.credentials,
                    loaded: true,
                });
            } catch (error) {
                setState({
                    ...state,
                    errors: true,
                })
            }
        }
        getCreds();
    }, []);
    return (
        <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center" }}>
            <div style={{ flex: "0 1 700px", margin: "5px" }}>
                <h1>Credentials</h1>
                {state.loaded ?
                    state.credentials.map((cred: any, inc: number) => {
                        return (
                            <Credential key={inc}/>
                        )
                    })
                    : <p>{state.errors ? "Error getting credentials" : "Loading ..."}</p>
                }
            </div>
        </div>
    )
}