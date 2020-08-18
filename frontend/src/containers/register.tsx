import * as React from 'react';
import Cookies from 'universal-cookie';
import { axiosInstance } from '../utils/axios';
import { Input } from '../components/input/input';
import { Button } from '../components/button/button';
import { bufferDecode, bufferEncode } from '../utils/webauthn';

export const Register = () => {
    const cookies = new Cookies();
    const [inputs, setInputs] = React.useState<any>({
        name: "",
        email: "",
    });
    const onInputUpdate = (event: React.ChangeEvent<HTMLInputElement>) => {
        setInputs({
            ...inputs,
            [event.target.name]: event.target.value
        });
    }

    const Register = async () => {
        //Contact api and register
        try {
            const req = await axiosInstance.post("/auth/register/start", {
                name: inputs.name,
                username: inputs.email,
            })

            //Todo remove, will make this jwt
            let sessJSON = JSON.stringify(req.data.session_data);
            let sessB64 = Buffer.from(sessJSON).toString("base64");
            cookies.set("register-token", sessB64)

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
        <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center" }}>
            <div style={{ flex: "0 1 350px", margin: "5px" }}>
                <h1>Register</h1>
                <Input
                    label="Name"
                    name="name"
                    placeHolder="Test User"
                    onChange={onInputUpdate}
                    value={inputs.name}
                />
                <Input
                    label="Email or username"
                    name="email"
                    // validationText="Cannot be blank"
                    placeHolder="email@example.com"
                    onChange={onInputUpdate}
                    value={inputs.email}
                />
                <div style={{height: "20px"}}/>
                <Button onClick={Register}>
                    Register account
                </Button>
                {/* <p className="validation-text">{state.mainError}</p> */}
            </div>
        </div>
    );
}