import * as React from 'react';
import { Link } from 'react-router-dom';
import { axiosAuth } from '../utils/axios';


export const Profile = () => {
    const [state, setState] = React.useState<any>({
        username: "",
        name: "",
        loading: true,
        errors: false,
    });
    React.useEffect(() => {
        async function getCreds() {
            try {
                const req = await axiosAuth.get("/api/profile")
                setState({
                    ...state,
                    username: req.data.username,
                    name: req.data.name,
                    loading: false,
                });
            } catch (error) {
                setState({
                    ...state,
                    errors: true,
                })
            }
        }
        getCreds()
    }, []);
    return (
        <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center" }}>
            <div style={{ flex: "0 1 800px", margin: "5px" }}>
                <h1>Personal Info</h1>
                <div style={{ padding: "2em", border: "1px solid #CCCCCC", borderRadius: "10px" }}>
                    {state.loading ?
                        <p>Loading</p>
                        :
                        <React.Fragment>
                            <p>Username: {state.username}</p>
                            <p>Name: {state.name}</p>
                            <Link to="/credentials">Credentials </Link>
                        </React.Fragment>
                    }
                </div>
                {state.errors && <p>Some kind of error occured try again later</p>}
            </div>
        </div>
    )
}