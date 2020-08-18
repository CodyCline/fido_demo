import * as React from 'react';
import { axiosAuth } from '../utils/axios';


export const Profile = () => {
    const [state, setState] = React.useState<any>({
        username: "",
        name: "",
        loaded: false,
    });
    React.useEffect(() => {
        async function getCreds() {
            try {
                const req = await axiosAuth.get("/api/profile")
                console.log(req);
                setState({
                    ...state,
                    username: req.data.username,
                    name: req.data.name,
                    loaded: true,
                });
            } catch (error) {

            }
        }
        getCreds()
    }, []);
    return (
        <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center" }}>
            <div style={{ flex: "0 1 800px", margin: "5px" }}>
                <h1>Personal Info</h1>
                {state.loaded ?
                <div style={{ border: "2px solid #CCCCCC", borderRadius: "10px" }}>
                    <p>Username: {state.username}</p>
                    <p>Name: {state.name}</p>
                    <p>Your credentials</p>
                    <p>Your todos</p>
                </div>
                :
                <p>Loading!!!</p>
            }
            </div>
        </div>
    )
}