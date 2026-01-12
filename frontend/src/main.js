import { Router, navigateTo } from './router.js';
import { renderHome, renderLogin, renderRegister } from './pages.js';

let isLoggedIn = false;
let username = '';

// Vérifier l'authentification au chargement de la page
async function checkAuthentication() {
    try {
        const response = await fetch('http://localhost:8086/api/check-auth', {
            method: 'GET',
            credentials: 'include', // Important pour envoyer les cookies
        });

        const data = await response.json();

        if (data.authenticated) {
            console.log('✅ Utilisateur connecté:', data.username);
            isLoggedIn = true;
            username = data.username
            // L'utilisateur est connecté, rester sur la page actuelle ou rediriger vers home
            return true;
        } else {
            console.log('❌ Pas de session active');
            // Rediriger vers la page de login si pas sur une page publique
            const currentPath = window.location.hash.replace('#', '') || '/';
            if (currentPath !== '/login' && currentPath !== '/register') {
                navigateTo('/login');
            }
            return false;
        }
    } catch (error) {
        console.error('Erreur lors de la vérification de l\'authentification:', error);
        return false;
    }
}

// Initialize app after DOM is loaded
document.addEventListener('DOMContentLoaded', async () => {
    await checkAuthentication();
    
    const routes = {
        '/': () => renderHome(isLoggedIn, username),
        '/login': renderLogin,
        '/register': renderRegister,
    };
    
    const router = new Router(routes);
});

// Re-check authentication and re-render on page load
window.addEventListener('load', async () => {
    await checkAuthentication();
    // Force the router to re-render the current route
    const currentPath = window.location.hash.replace('#', '') || '/';
    window.location.hash = '#/temp';
    setTimeout(() => {
        window.location.hash = currentPath;
    }, 0);
});

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
    }else if (e.target.id === 'LogoutBtn') {
        e.preventDefault();
        const logoutPopup = document.getElementById("logoutPopup");
        const logoutOverlay = document.getElementById("logoutOverlay");
        logoutOverlay.style.display = "block";
        logoutPopup.style.display = "block";
        console.log("Logout button clicked");
    }

});