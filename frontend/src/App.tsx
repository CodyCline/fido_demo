import React from 'react';
import { Route, Switch } from 'react-router-dom';
import { Home } from './containers/home';
import { Login } from './containers/login';
import { Register } from './containers/register';
import { Profile } from './containers/profile';
import { Credentials } from './containers/credentials'; 
import { AuthRequired } from './components/auth/auth';
import './App.css';
import Cookies from 'universal-cookie';


function App() {
	const cookies = new Cookies();
	const [state] = React.useState({
		token: cookies.get("token")
	});
	return (
		<Switch>
			<Route exact path="/" component={Home} />
			<Route path="/login" component={Login} />
			<Route path="/register" component={Register} />
			{/* Requires authentication */}
			<AuthRequired isAuthenticated={state.token}> 
				<Route path="/profile" component={Profile} />
				<Route path="/credentials" component={Credentials} />
			</AuthRequired>
			<Route/>
		</Switch>
	);
}

export default App;
