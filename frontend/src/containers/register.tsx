import * as React from 'react';
import { axiosInstance } from '../utils/axios';
import { Input } from '../components/input/input';
import { bufferDecode, bufferEncode } from '../utils/webauthn';


export const Register = () => {
    const [inputs, setInputs] = React.useState<any>({
        name: "name",
        email: "2@2.com",
    });
    const onInputUpdate = (event: any, name: string) => {
        setInputs({
            ...inputs,
            [name]: event.target.value
        })
    }

    const Register = async () => {
        //Contact api and 
        try {
            const req = await axiosInstance.post("/auth/register/start", {
                name: inputs.name,
                username: inputs.email,
            })
            console.log(req.data)

            //After receiving registration data, decode/mutate response
            let credentialCreationOptions = req.data;
            let { user, challenge, excludeCredentials } = credentialCreationOptions.publicKey;
            credentialCreationOptions.publicKey.challenge = bufferDecode(challenge);
            credentialCreationOptions.publicKey.user.id = bufferDecode(user.id);
            if (excludeCredentials) {
                excludeCredentials.map((cred: any) => {
                    cred.id = bufferDecode(cred.id)
                })
            }

            //Call browser to read security key
            const credential: any = await navigator.credentials.create({
                publicKey: credentialCreationOptions.publicKey
            })
            let attestationObject = credential.response.attestationObject;
            let clientDataJSON = credential.response.clientDataJSON;
            let rawId = credential.rawId;
            const debug = {
                id: credential.id,
                rawId: bufferEncode(rawId),
                type: credential.type,
                response: {
                    attestationObject: bufferEncode(attestationObject),
                    clientDataJSON: bufferEncode(clientDataJSON),
                },
            }
            console.log(debug)
            const success = await axiosInstance.post("/auth/register/finish/" + inputs.email, {
                id: credential.id,
                rawId: bufferEncode(rawId),
                type: credential.type,
                response: {
                    attestationObject: bufferEncode(attestationObject),
                    clientDataJSON: bufferEncode(clientDataJSON),
                },
            })
            console.log("succesfully registered", success)
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