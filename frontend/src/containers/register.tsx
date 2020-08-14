import * as React from 'react';
import Cookies from 'universal-cookie';
import { axiosInstance } from '../utils/axios';
import { Input } from '../components/input/input';
import { bufferDecode, bufferEncode } from '../utils/webauthn';

export const Register = () => {
    const cookies = new Cookies();
    const [inputs, setInputs] = React.useState<any>({
        name: "name",
        email: "2@2.com",
    });
    const onInputUpdate = (event: any, name: string) => {
        setInputs({
            ...inputs,
            [name]: event.target.value
        });
    }

    const Register = async () => {
        //Contact api and register
        try {
            const req = await axiosInstance.post("/auth/register/start", {
                name: inputs.name,
                username: inputs.email,
            })

            //After receiving registration data, decode/mutate response
            console.log(req)
            let objJSONStr = JSON.stringify(req.data.session_data);
            let objJSONB64 = Buffer.from(objJSONStr).toString("base64");
            cookies.set("register-token", objJSONB64)
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
            const sessionData = cookies.get("register-token")
            console.log(credential)
            let { attestationObject, clientDataJSON, rawId } = credential.response;
            const success = await axiosInstance.post(`/auth/register/finish/${inputs.email}/${sessionData}`, {
                id: credential.id,
                rawId: bufferEncode(rawId),
                type: credential.type,
                response: {
                    attestationObject: bufferEncode(attestationObject),
                    clientDataJSON: bufferEncode(clientDataJSON),
                },
            })
            console.log("Success\n", success);
        } catch (error) {
            console.log(error);
        }
    }

    return (
        <div>
            <h1>Register</h1>

            <p>Name</p>
            <Input
                name="name"
                onChange={(event: any) => onInputUpdate(event, "name")}
                value={inputs.name}
            />
            <p>Email</p>
            <Input
                name="email"
                onChange={(event: any) => onInputUpdate(event, "email")}
                value={inputs.email}
            />
            <button onClick={Register}>Register</button>
        </div>
    )
}