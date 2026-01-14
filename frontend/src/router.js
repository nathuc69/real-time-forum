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
        // D'abord, essayer de trouver une correspondance exacte
        if (this.routes[path]) {
            this.routes[path]();
            return;
        }

        // Ensuite, essayer de matcher avec des routes paramétrées
        for (const routePath in this.routes) {
            const params = this.matchRoute(routePath, path);
            if (params) {
                this.routes[routePath](params);
                return;
            }
        }

        // Si aucune route ne correspond, utiliser la route par défaut
        if (this.routes['/']) {
            this.routes['/']();
        }
    }

    matchRoute(routePath, actualPath) {
        const routeParts = routePath.split('/');
        const actualParts = actualPath.split('/');

        if (routeParts.length !== actualParts.length) {
            return null;
        }

        const params = {};
        for (let i = 0; i < routeParts.length; i++) {
            if (routeParts[i].startsWith(':')) {
                const paramName = routeParts[i].substring(1);
                params[paramName] = actualParts[i];
            } else if (routeParts[i] !== actualParts[i]) {
                return null;
            }
        }

        return Object.keys(params).length > 0 ? params : null;
    }
}

// Navigate programmatically from anywhere
export function navigateTo(path) {
    window.location.hash = path;
}
