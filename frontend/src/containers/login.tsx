import * as React from 'react';
import Cookies from 'universal-cookie';
import { axiosInstance } from '../utils/axios';
import { Input } from '../components/input/input';
import { bufferDecode, bufferEncode } from '../utils/webauthn';

export const Login = () => {
    const cookies = new Cookies();
    const [inputs, setInputs] = React.useState<any>({
        email: "2@2.com",
    });
    const onInputUpdate = (event: any, name: string) => {
        setInputs({
            ...inputs,
            [name]: event.target.value
        });
    }

    const Login = async () => {
        //Contact api and login
        try {
            const start = await axiosInstance.post("/auth/login/start", {
                username: inputs.email,
            })

            //After receiving login challenge, decode/mutate response
            let stringSession = JSON.stringify(start.data.session_data);
            let base64Session = Buffer.from(stringSession).toString("base64");
            cookies.set("login-token", base64Session)
            let credentialRequestOptions = start.data.options;
            let { challenge, allowCredentials } = credentialRequestOptions.publicKey;
            credentialRequestOptions.publicKey.challenge = bufferDecode(challenge);
            allowCredentials.map((listItem:any) => {
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
            const finish = await axiosInstance.post(`/auth/login/finish/${inputs.email}/${sessionData}`, {
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
                cookies.set("token", finish.data.token);
                const todos = await axiosInstance.get('/todos', {
                    headers: {
                        'Authorization': `Basic ${finish.data.token}` 
                    }
                })
                console.log(todos)
            }
        } catch (error) {
            console.log(error);
        }
    }

    return (
        <div>
            <h1>Login</h1>

            <p>Email</p>
            <Input
                name="email"
                onChange={(event: any) => onInputUpdate(event, "email")}
                value={inputs.email}
            />
            <button onClick={Login}>Login</button>
        </div>
    );
}