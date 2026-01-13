import {navigateTo } from './router.js';
export function EventLoginBtn(){
    const loginBtn = document.getElementById("loginBtn");
    const loginMenu = document.getElementById("loginMenu");
    const menu = document.getElementById("MenuPage");
    const logoutPopup = document.getElementById("logoutPopup");
    const logoutBtn = document.getElementById("LogoutBtn");

    if (!loginBtn || !loginMenu) {
        console.error("Login button or menu not found in DOM");
        return;
    }

    logoutBtn.addEventListener("click", (e) => {
        e.preventDefault();
        const logoutPopupElement = document.getElementById("logoutPopup");
        const logoutOverlay = document.getElementById("logoutOverlay");
        logoutOverlay.style.display = "block";
        logoutPopupElement.style.display = "block";
        console.log("Logout button clicked");
    });

    loginBtn.addEventListener("click", (e) => {
        e.preventDefault();
        loginMenu.style.display = "block";
        menu.style.display = "none";
        console.log("Login button clicked");
    });
}

export function setupEventListeners() {
    const submitLoginBtn = document.getElementById("SubmitLogin");
    const submitRegisterBtn = document.getElementById("SubmitRegister");
    const SubmitLogoutBtn = document.getElementById("SubmitLogoutBtn");
    const cancelLogoutBtn = document.getElementById("cancelLogoutBtn");
    
    if (submitLoginBtn) {
        submitLoginBtn.addEventListener("click", handleLogin);
    }
    if (submitRegisterBtn) {
        submitRegisterBtn.addEventListener("click", handleRegister);
    }
    if (SubmitLogoutBtn) {
        SubmitLogoutBtn.addEventListener("click", handleLogout);
    }
    if (cancelLogoutBtn) {
        cancelLogoutBtn.addEventListener("click", (e) => {
            e.preventDefault();
            const logoutPopupElement = document.getElementById("logoutPopup");
            const logoutOverlay = document.getElementById("logoutOverlay");
            logoutPopupElement.style.display = "none";
            logoutOverlay.style.display = "none";
            navigateTo('/');
        });
    }
}

export async function handleLogin(e) {
    e.preventDefault();

    const usernameVal = document.getElementById("username").value;
    if (usernameVal.includes("@")) {
        var email = usernameVal;
        var username = "";
    } else {
        var username = usernameVal;
        var email = "";
    }
    const password = document.getElementById("password").value;

    try {
        const response = await fetch('http://localhost:8086/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, password }),
            credentials: 'include',
        });

        if (response.ok) {  
            const data = await response.json();
            console.log('Login successful:', data);
            // Redirect to home and force page reload to update auth state
            window.location.hash = '#/';
            window.location.reload();
        } else {
            const errorData = await response.json();
            console.error('Login failed:', errorData);
            // Handle login failure (e.g., show error message)
        }
    } catch (error) {
        console.error('Error during login:', error);
    }
}

export async function handleRegister(e) {
    e.preventDefault();

    const username = document.getElementById("username").value;
    const email = document.getElementById("email").value;
    const age = document.getElementById("age").value;
    const firstName = document.getElementById("firstName").value;
    const lastName = document.getElementById("lastName").value
    const password = document.getElementById("password").value;
    const Gender = document.getElementById("gender").value;

    try {

        const response = await fetch('http://localhost:8086/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, age: Number(age), firstName, lastName , Gender, password  }),
            credentials: 'include',
        });

        if (response.ok) {
            const data = await response.json();
            console.log('Registration successful:', data);
            // Handle successful registration (e.g., redirect, update UI)
        } else {
            const errorData = await response.json();
            console.error('Registration failed:', errorData);
            // Handle registration failure (e.g., show error message)
        }
    } catch (error) {
        console.error('Error during registration:', error);
    }
}

export async function handleLogout(e) {
    e.preventDefault();

    try {
        const response = await fetch('http://localhost:8086/api/logout', {
            method: 'POST',
            credentials: 'include',
        });

        if (response.ok) {
            console.log('Logout successful');
            // Hide the logout popup and overlay
            const logoutPopupElement = document.getElementById("logoutPopup");
            const logoutOverlay = document.getElementById("logoutOverlay");
            logoutPopupElement.style.display = "none";
            logoutOverlay.style.display = "none";
            
            // Redirect to login page
            window.location.hash = '#/login';
            window.location.reload();
        } else {
            console.error('Logout failed');
        }
    } catch (error) {
        console.error('Error during logout:', error);
    }
}