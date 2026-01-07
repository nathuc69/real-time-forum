import { Router, navigateTo } from './router.js';
import { renderHome, renderLogin, renderRegister } from './pages.js';

// Define routes
const routes = {
    '/': renderHome,
    '/login': renderLogin,
    '/register': renderRegister,
};

// Initialize router
const router = new Router(routes);

// Handle clicks on navigation buttons (event delegation)
document.addEventListener('click', (e) => {
    if (e.target.id === 'loginBtn') {
        e.preventDefault();
        navigateTo('/login');
    } else if (e.target.id === 'registerBtn') {
        e.preventDefault();
        navigateTo('/register');
    } else if (e.target.id === 'backBtn') {
        e.preventDefault();
        navigateTo('/');
    }
});