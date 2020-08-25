import * as React from 'react';
import { axiosAuth } from '../utils/axios';
import { Credential } from '../components/credential/credential';
import { Input } from '../components/input/input';
import { Button } from '../components/button/button';
import { bufferDecode, bufferEncode } from '../utils/webauthn';
import Cookies from 'universal-cookie';

export const Credentials = () => {
    const cookies = new Cookies();
    const [state, setState] = React.useState<any>({
        credentials: [],
        credName: "",
        loaded: false,
    });
    const addCredential = async() => {
        //Add...-+
        try {
            const req = await axiosAuth.get("/api/credential/start");

            //Todo remove, will make this jwt
            let sessJSON = JSON.stringify(req.data.session_data);
            let sessB64 = Buffer.from(sessJSON).toString("base64");
            cookies.set("cred-token", sessB64)

            //After receiving registration data, decode/mutate response
            let credentialCreationOptions = req.data.options;
            let { user, challenge, excludeCredentials } = credentialCreationOptions.publicKey;
            credentialCreationOptions.publicKey.challenge = bufferDecode(challenge);
            credentialCreationOptions.publicKey.user.id = bufferDecode(user.id);
            if (excludeCredentials) {
                excludeCredentials.map((cred: any) => {
                    cred.id = bufferDecode(cred.id)
                })
            }
            //Call browser to insert key
            const credential: any = await navigator.credentials.create({
                publicKey: credentialCreationOptions.publicKey
            })
            const sessionData = cookies.get("cred-token")
            console.log(credential)
            let { attestationObject, clientDataJSON, rawId } = credential.response;
            const newCred = await axiosAuth.post(`/api/credential/finish/${state.credName}/${sessionData}`, {
                id: credential.id,
                rawId: bufferEncode(rawId),
                type: credential.type,
                response: {
                    attestationObject: bufferEncode(attestationObject),
                    clientDataJSON: bufferEncode(clientDataJSON),
                },
            })
            console.log(newCred)
            setState({
                ...state,
                credentials: [
                    ...state.credentials,
                    newCred.data,
                ]
            })
        } catch (error) {
            console.log(error);
        }
        
    }

    const deleteCredential = async(id: string | number) => {
        const req = await axiosAuth.delete("/api/credential/" + id);
        const updatedCreds = state.credentials.filter((cred:any) => {
            return cred.id !== id
        })
        setState({
            ...state,
            credentials: updatedCreds
        })
    }

    const onInputChange = (event:any) => {
        setState({
            ...state,
            credName: event.target.value,
        })
    }

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
                        onChange={onInputChange}
                        label="Add another credential in case you lose one."
                        placeHolder="Name for credential (e.g. bluetooth key)"
                    />
                    <Button onClick={addCredential} style={{width:"100px"}}>Add</Button>
                </div>
                {state.loaded ?
                    state.credentials.map((cred: any, inc: number) => {
                        const updated = new Date(cred.updated_at)
                        return (
                            <Credential 
                                name={cred.nickname}
                                key={cred.id} 
                                lastUsed={updated.toLocaleString()} 
                                useCount={cred.sign_count} 
                                onDelete={() => deleteCredential(cred.id)}
                            />
                        )
                    })
                    : <p>{state.errors ? "Error getting credentials" : "Loading ..."}</p>
                }
            </div>
        </div>
    )
}