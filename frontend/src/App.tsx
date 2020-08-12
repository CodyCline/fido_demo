import React from 'react';
import { Route, Switch } from 'react-router-dom';
import { Home } from './containers/home';
import { Login } from './containers/login';
import { Register } from './containers/register';
import { Profile } from './containers/profile';
import { Credentials } from './containers/credentials'; 
import './App.css';


function App() {
	return (
		<Switch>
			<Route exact path="/" component={Home} />
			<Route path="/login" component={Login} />
			<Route path="/register" component={Register} />
			{/* Requires authentication */}
			<Route path="/profile" component={Profile} />
			<Route path="/credentials" component={Credentials} />
			<Route/>
		</Switch>
	);
}

export default App;
