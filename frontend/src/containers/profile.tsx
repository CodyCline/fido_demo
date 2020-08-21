import * as React from 'react';
import { Link } from 'react-router-dom';
import { axiosAuth } from '../utils/axios';
import { ReactComponent as Key } from '../assets/key-solid.svg';
import Cookies from 'universal-cookie';

export const Profile = () => {
    const cookies = new Cookies();
    const [state, setState] = React.useState<any>({
        username: "",
        name: "",
        loading: true,
        errors: false,
    });
    React.useEffect(() => {
        async function getCreds() {
            try {
                const req = await axiosAuth.get("/api/profile", {
                    headers: {
                        Authorization: `Bearer ${cookies.get("token")}`
                    }
                })
                setState({
                    ...state,
                    username: req.data.username,
                    name: req.data.name,
                    loaded: true,
                });
            } catch (error) {
                setState({
                    ...state,
                    errors: true,
                })
            }
        }
        // setTimeout(getCreds, 3000);
        getCreds()
    }, []);
    return (
        <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center" }}>
            <div style={{ flex: "0 1 800px", margin: "5px" }}>
                <h1>Personal Info</h1>
                <div style={{ padding: "2em", border: "1px solid #CCCCCC", borderRadius: "10px" }}>
                    {state.loaded ?
                        <React.Fragment>
                            <p>Username: {state.username}</p>
                            <p>Name: {state.name}</p>
                        </React.Fragment>
                        : 
                        <p>{state.errors ? "Error getting profile" : "Loading ..."}</p>
                    }
                </div>
                <div style={{height: "20px"}}/>
                <li style={{background: "#CCC", padding: "10px", display: "flex", flexDirection:"row", alignItems:"center"}}>
                    <Key style={{margin: "10px", height: "20px"}}/>
                    <Link to="/credentials">Manage Credentials</Link>
                </li>
            </div>
        </div>
    )
}