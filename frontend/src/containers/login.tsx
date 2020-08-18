import * as React from 'react';
import Cookies from 'universal-cookie';
import { axiosInstance } from '../utils/axios';
import { Input } from '../components/input/input';
import { Button } from '../components/button/button';
import { bufferDecode, bufferEncode } from '../utils/webauthn';
import { Redirect, useHistory } from 'react-router-dom';

export const Login = () => {
    const history = useHistory();
    const cookies = new Cookies();
    const [state, setState] = React.useState<any>({
        email: "",
        loggingIn: false,
        mainError: "",
    })
    const onInputUpdate = (event: React.ChangeEvent<HTMLInputElement>) => {
        setState({
            ...state,
            [event.target.name]: event.target.value
        });
    }

    const Login = async () => {
        //Contact api and start login process
        try {
            const start = await axiosInstance.post("/auth/login/start", {
                username: state.email,
            })

            //If username exists handle here

            //After receiving login challenge, decode/mutate response
            let stringSession = JSON.stringify(start.data.session_data);
            let base64Session = Buffer.from(stringSession).toString("base64");
            cookies.set("login-token", base64Session)
            let credentialRequestOptions = start.data.options;
            let { challenge, allowCredentials } = credentialRequestOptions.publicKey;
            credentialRequestOptions.publicKey.challenge = bufferDecode(challenge);
            allowCredentials.map((listItem: any) => {
                listItem.id = bufferDecode(listItem.id)
            })
            //Call browser to insert key
            const assertion: any = await navigator.credentials.get({
                publicKey: credentialRequestOptions.publicKey
            })
            console.log("Assertion", assertion)
            const sessionData = cookies.get("login-token")
            let { authenticatorData, clientDataJSON, signature, userHandle } = assertion.response;
            console.log({
                id: assertion.id,
                rawId: bufferEncode(assertion.rawId),
                type: assertion.type,
                response: {
                    authenticatorData: bufferEncode(authenticatorData),
                    clientDataJSON: bufferEncode(clientDataJSON),
                    signature: bufferEncode(signature),
                    userHandle: bufferEncode(userHandle),
                },
            })
            const finish = await axiosInstance.post(`/auth/login/finish/${state.email}/${sessionData}`, {
                id: assertion.id,
                rawId: bufferEncode(assertion.rawId),
                type: assertion.type,
                response: {
                    authenticatorData: bufferEncode(authenticatorData),
                    clientDataJSON: bufferEncode(clientDataJSON),
                    signature: bufferEncode(signature),
                    userHandle: bufferEncode(userHandle),
                },
            })
            if (finish.data.success) {
                //Remove any cookies
                await cookies.set("token", finish.data.token);
            }
            setTimeout(() => history.push("/profile"), 3000)
        } catch (error) {
            console.log(error);
            setState({
                ...state, 
                mainError: `Error: ${error.message}`,
            })
        }
    }

    return (
        <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center" }}>
            <div style={{ flex: "0 1 350px", margin: "5px"}}>
                <h1>Login</h1>
                <Input
                    label="Email or username"
                    name="email"
                    // validationText="Cannot be blank"
                    placeHolder="email@example.com"
                    onChange={onInputUpdate}
                    value={state.email}
                />
                <div style={{height: "20px"}}/>
                <Button onClick={Login}>
                    {state.loggingIn? "Hold on ..." : "Sign in without password"}
                </Button>
                <p className="validation-text">{state.mainError}</p>
            </div>
        </div>
    );
}