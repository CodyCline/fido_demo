import * as React from 'react';
import Cookies from 'universal-cookie';
import { axiosInstance } from '../utils/axios';
import { Input } from '../components/input/input';
import { Button } from '../components/button/button';
import { bufferDecode, bufferEncode } from '../utils/webauthn';
import { useForm } from '../utils/useForm/useForm';
import { validate } from '../utils/useForm/loginValidations';

export const Login = () => {
    const cookies = new Cookies();
    const [state] = React.useState<any>({
        username: "",
        loggingIn: false,
    });
    const [err, setErr] = React.useState<string>("");

    const login = async () => {
        //Contact api and start login process
        try {
            const start = await axiosInstance.post("/auth/login/start", {
                username: values.username,
            })


            //Todo remove            
            let stringSession = JSON.stringify(start.data.session_data);
            let base64Session = Buffer.from(stringSession).toString("base64");

            cookies.set("login-token", base64Session)
            //After receiving login challenge, decode/mutate response
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
            const finish = await axiosInstance.post(`/auth/login/finish/${values.username}/${sessionData}`, {
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
                //Remove any cookies redirect, set state.
                await cookies.set("token", finish.data.token);
            }
        } catch (error) {
            if (!error.response.data.success) { //Axios error
                setErr(error.response.data.message);
            } else {
                setErr("Something went wrong try again soon");
            }
        }
    }

    const { values, errors, handleChange, handleSubmit } = useForm(
        state,
        login,
        validate
    );

    return (
        <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center" }}>
            <div style={{ flex: "0 1 350px", margin: "5px"}}>
                <h1>Login</h1>
                <Input
                    label="Email or username"
                    name="username"
                    validationText={errors.username}
                    placeHolder="email@example.com"
                    onChange={handleChange}
                    value={values.username}
                />
                <div style={{height: "20px"}}/>
                <Button onClick={handleSubmit}>Sign in without password</Button>
                <p className="validation-text">{err}</p>
            </div>
        </div>
    );
}