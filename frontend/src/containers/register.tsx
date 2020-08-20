import * as React from 'react';
import Cookies from 'universal-cookie';
import { axiosInstance } from '../utils/axios';
import { Input } from '../components/input/input';
import { Button } from '../components/button/button';
import { bufferDecode, bufferEncode } from '../utils/webauthn';
import { useForm } from '../utils/useForm/useForm';
import { validate } from '../utils/useForm/registerValidations';
import { useHistory } from 'react-router-dom';

export const Register = () => {
    const cookies: Cookies = new Cookies();
    const history = useHistory();
    const [err, setErr] = React.useState<string>("")
    const [state] = React.useState<any>({
        name: "",
        username: "",
    });

    const register = async () => {
        //Contact api and register
        try {
            const req = await axiosInstance.post("/auth/register/start", {
                name: values.name,
                username: values.username,
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
            const success = await axiosInstance.post(`/auth/register/finish/${values.username}/${sessionData}`, {
                id: credential.id,
                rawId: bufferEncode(rawId),
                type: credential.type,
                response: {
                    attestationObject: bufferEncode(attestationObject),
                    clientDataJSON: bufferEncode(clientDataJSON),
                },
            })
            history.push("/login")
        } catch (error) {
            if (error.response) { //Axios error
                setErr(error.response.data.message);
            } else {
                setErr("Something went wrong try again soon");
            }
        }
    }

    const { values, errors, handleChange, handleSubmit } = useForm(
        state,
        register,
        validate
    );

    return (
        <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center" }}>
            <div style={{ flex: "0 1 350px", margin: "5px" }}>
                <h1>Register</h1>
                <Input
                    label="Name"
                    name="name"
                    placeHolder="Test User"
                    onChange={handleChange}
                    value={values.name}
                    validationText={errors.name}
                />
                <Input
                    label="Email or username"
                    name="username"
                    validationText={errors.username}
                    placeHolder="email@example.com"
                    onChange={handleChange}
                    value={values.username}
                />
                <div style={{height: "20px"}}/>
                <Button onClick={handleSubmit}>
                    Register account
                </Button>
                <p className="error">{err}</p>
            </div>
        </div>
    );
}