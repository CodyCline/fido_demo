import * as React from 'react';
import { Link } from 'react-router-dom';
import { Spoiler } from '../components/spoiler/spoiler';

export const Home = () => {
    return (
        <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center" }}>
            <div style={{ flex: "0 1 800px", margin: "5px" }}>
                <h1>Welcome to WebAuthn demo</h1>
                <p>This is an example app that implements an emerging security standard which allows passwordless registration and sign in using the WebAuthn API. All that's required is a hardware security key to test it out.</p>
                <p><Link to="/register">Register</Link></p>
                <p><Link to="/login">Login</Link></p>
                <p>FAQ</p>
                <Spoiler title="I don't have a hardware key?">
                    There are many manufacturers of hardware keys.
                        Some popular choices are the <a href="https://store.google.com/product/titan_security_key" target="_blank" rel="noopener noreferrer">Google' Titan Key</a>, <a href="https://www.yubico.com/products/yubikey-5-overview/" target="_blank" rel="noopener noreferrer">Yubico' Yubikey</a>.
                        For more advanced users, there are tutorials on how to craft your own hardware key.
                </Spoiler>
                <Spoiler title="Browser Support ***">
                    The following web browsers support this authentication method.
                    <li>Chrome 67+</li>
                    <li>Firefox 60+</li>
                    <li>Safari 13+</li>
                    <li>Opera 54+</li>
                </Spoiler>
            </div>
        </div>
    )
}