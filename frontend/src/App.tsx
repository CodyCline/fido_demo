import React from 'react';
import { Route, Switch } from 'react-router-dom';
import { Layout } from './components/layout/layout';
import { AuthRequired } from './components/auth/auth';
import { Home } from './containers/home';
import { Login } from './containers/login';
import { Register } from './containers/register';
import { Profile } from './containers/profile';
import { Credentials } from './containers/credentials';
import './App.css';
import Cookies from 'universal-cookie';

function App() {
	return (
		<Layout>
			<Switch>
				<Route exact path="/" component={Home} />
				<Route path="/login" component={Login} />
				<Route path="/register" component={Register} />
				{/* Requires authentication */}
				<AuthRequired>
					<Route path="/profile" component={Profile} />
					<Route path="/credentials" component={Credentials} />
				</AuthRequired>
				<Route />
			</Switch>
		</Layout>
	);
}

export default App;
