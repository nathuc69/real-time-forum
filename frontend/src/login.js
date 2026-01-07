// function LoginDOM(){
//     const LoginMenu = document.createElement("div");
//     LoginMenu.id = "loginMenu";

//     const loginForm = document.createElement("form");
//     loginForm.id = "loginForm";

//     const usernameInput = document.createElement("input");
//     usernameInput.type = "text";
//     usernameInput.id = "usernameInput";
//     usernameInput.placeholder = "Username";

//     const passwordInput = document.createElement("input");
//     passwordInput.type = "password";
//     passwordInput.id = "passwordInput";
//     passwordInput.placeholder = "Password";

//     const loginBtn = document.createElement("button");
//     loginBtn.type = "submit";
//     loginBtn.id = "loginBtn";
//     loginBtn.innerText = "Login";

//     loginForm.appendChild(usernameInput);
//     loginForm.appendChild(passwordInput);
//     loginForm.appendChild(loginBtn);

//     LoginMenu.appendChild(loginForm);

//     document.body.appendChild(LoginMenu);   
//     console.log("Login DOM created"); 
// }
export function EventLoginBtn(){
    const loginBtn = document.getElementById("loginBtn");
    const loginMenu = document.getElementById("loginMenu");
    const menu = document.getElementById("MenuPage");

    if (!loginBtn || !loginMenu) {
        console.error("Login button or menu not found in DOM");
        return;
    }

    loginBtn.addEventListener("click", (e) => {
        e.preventDefault();
        loginMenu.style.display = "block";
        menu.style.display = "none";
        console.log("Login button clicked");
    });
}

addEventListener('submit', handleLogin);

export async function handleLogin() {
    const SubmitBtnLogin = document.getElementById("SubmitLogin");
    SubmitBtnLogin.preventDefault();

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
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, password }),
        });

        if (response.ok) {
            const data = await response.json();
            console.log('Login successful:', data);
            // Handle successful login (e.g., redirect, update UI)
        } else {
            const errorData = await response.json();
            console.error('Login failed:', errorData);
            // Handle login failure (e.g., show error message)
        }
    } catch (error) {
        console.error('Error during login:', error);
    }
}

export async function handleRegister() {
    const SubmitBtnRegister = document.getElementById("SubmitRegister");
    SubmitBtnRegister.preventDefault();

    const username = document.getElementById("username").value;
    const email = document.getElementById("email").value;
    const age = document.getElementById("age").value;
    const firstName = document.getElementById("firstName").value;
    const lastName = document.getElementById("lastName").value

    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, age, firstName, lastName }),
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
