import * as React from 'react';
import { axiosAuth } from '../utils/axios';
import { Credential } from '../components/credential/credential';
import { Input } from '../components/input/input';
import { Button } from '../components/button/button';

export const Credentials = () => {
    const date = new Date();
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
            <div style={{ flex: "0 1 800px", margin: "5px" }}>
                <h1>Credentials</h1>
                <div>
                    <Input
                        label="Add another credential in case you lose one."
                        placeHolder="Name for credential (e.g. bluetooth key)"
                    />
                    <Button style={{width:"100px"}}>Add</Button>
                </div>
                {state.loaded ?
                    state.credentials.map((cred: any, inc: number) => {
                        const updated = new Date(cred.updated_at)
                        return (
                            <Credential 
                                key={cred.id} 
                                lastUsed={updated.toLocaleString()} 
                                useCount={cred.sign_count} 
                                //Onclick to delete
                            />
                        )
                    })
                    : <p>{state.errors ? "Error getting credentials" : "Loading ..."}</p>
                }
            </div>
        </div>
    )
}