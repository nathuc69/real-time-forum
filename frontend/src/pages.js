// Page rendering functions
import { setupEventListeners } from './login.js';

export function renderHome() {
    document.body.innerHTML = `
        <div id="MenuPage">
            <h1>Welcome to the Real-Time Forum</h1>
            <button id="loginBtn">Login</button>
            <button id="registerBtn">Register</button>
        </div>
    `;
}

export function renderLogin() {
    document.body.innerHTML = `
        <div id="loginMenu" style="display: block;">
            <form id="loginForm" method="post">
                <h2>Login</h2>
                <label for="username">Username or email</label>
                <input type="text" id="username" name="username" required>
                
                <label for="password">Password:</label>
                <input type="password" id="password" name="password" autocomplete="on" required>
                
                <button id="SubmitLogin" type="submit">Login</button>
                
                <label for="register">New user?</label>
                <button id="registerBtn" type="button">Register</button>
                <button id="backBtn" type="button">Back to Home</button>
            </form>
        </div>
    `;
    setupEventListeners();
}

export function renderRegister() {
    document.body.innerHTML = `
        <div id="loginMenu" style="display: block;">
            <form id="registerForm" method="post">
                <h2>Register</h2>
                <label for="username">Username:</label>
                <input type="text" id="username" name="username" required>
                
                <label for="email">Email:</label>
                <input type="email" id="email" name="email" required>

                <label for="age">Age:</label>
                <input type="number" id="age" name="age" required>

                <label for="firstName">First Name:</label>
                <input type="text" id="firstName" name="firstName" required>

                <label for="lastName">Last Name:</label>
                <input type="text" id="lastName" name="lastName" required>

                <label for="gender">Gender:</label>
                <select id="gender" name="gender" required>
                    <option value="male">Male</option>
                    <option value="female">Female</option>
                    <option value="other">Other</option>
                </select>

                <label for="password">Password:</label>
                <input type="password" id="password" name="password" required>
                
                <button id="SubmitRegister" type="submit">Register</button>
                <button id="backBtn" type="button">Back to Home</button>
            </form>
        </div>
    `;
    setupEventListeners();
}
