// Simple client-side hash router (avoids server 404 on refresh)
export class Router {
    constructor(routes) {
        this.routes = routes;
        this.init();
    }

    init() {
        window.addEventListener('hashchange', () => {
            this.handleRoute(this.currentPath());
        });

        this.handleRoute(this.currentPath());
    }

    currentPath() {
        const hash = window.location.hash || '#/';
        return hash.replace('#', '') || '/';
    }

    navigate(path) {
        window.location.hash = path;
    }

    handleRoute(path) {
        const route = this.routes[path] || this.routes['/'];
        if (route) {
            route();
        }
    }
}

// Navigate programmatically from anywhere
export function navigateTo(path) {
    window.location.hash = path;
}
