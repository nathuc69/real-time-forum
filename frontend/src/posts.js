import { navigateTo } from './router.js';
import { createReactionButtons } from './reaction.js';


export function EventPosts() {
    document.querySelectorAll('.post-card').forEach(card => {
        card.addEventListener('click', (e) => {
            // Prevent navigation if clicking on reaction buttons
            if (e.target.closest('.reaction-btn')) return;

            const postId = card.getAttribute('data-post-id');
            console.log(`Post ID clicked: ${postId}`);
            navigateTo(`/posts/${postId}`);
        });
    });
}


export function handlePosts(sortBy = "recent", category = "") {
    // Vérifier si l'utilisateur est connecté
    const isLoggedIn = document.cookie.includes('session_token');

    let url = "http://localhost:8086/api/posts";
    const params = new URLSearchParams();
    if (sortBy) params.append("sort", sortBy);
    if (category) params.append("category", category);

    if (params.toString()) url += `?${params.toString()}`;

    fetch(url)
        .then((response) => response.json())
        .then((data) => {
            const postsContainer = document.getElementById("postsContainer");
            postsContainer.innerHTML = ""; // Clear existing posts

            data.forEach((post) => {
                const postElement = document.createElement("div");
                postElement.className = "post-card";
                postElement.setAttribute("data-post-id", post.id);

                // Header avec avatar et info auteur
                const postHeader = document.createElement("div");
                postHeader.className = "post-header";

                const avatarElement = document.createElement("div");
                avatarElement.className = "post-avatar";
                avatarElement.textContent = (post.username || post.author || "?")[0].toUpperCase();

                const authorInfo = document.createElement("div");
                authorInfo.className = "post-author-info";

                const authorElement = document.createElement("span");
                authorElement.className = "post-author";
                authorElement.textContent = post.username || post.author || "Anonymous";

                const dateElement = document.createElement("span");
                dateElement.className = "post-date";
                const postDate = new Date(post.createdAt || post.created_at);
                dateElement.textContent = formatDate(postDate);

                authorInfo.appendChild(authorElement);
                authorInfo.appendChild(dateElement);
                postHeader.appendChild(avatarElement);
                postHeader.appendChild(authorInfo);

                // Categories
                if (post.categories && post.categories.length > 0) {
                    const categoriesElement = document.createElement("div");
                    categoriesElement.className = "post-categories";
                    categoriesElement.style.marginTop = "5px";
                    categoriesElement.innerHTML = post.categories.map(cat =>
                        `<span class="category-tag" style="display: inline-block; background: #e0e0e0; padding: 2px 6px; border-radius: 4px; font-size: 0.8em; margin-right: 5px; color: #555;">${cat}</span>`
                    ).join('');
                    postHeader.appendChild(categoriesElement);
                }

                // Titre du post
                const titleElement = document.createElement("h3");
                titleElement.className = "post-title";
                titleElement.textContent = post.title;

                // Contenu du post (tronqué)
                const contentElement = document.createElement("p");
                contentElement.className = "post-content";
                const maxLength = 150;
                contentElement.textContent = post.content.length > maxLength
                    ? post.content.substring(0, maxLength) + "..."
                    : post.content;

                // Footer avec indicateurs
                const postFooter = document.createElement("div");
                postFooter.className = "post-footer";

                // Compteur de commentaires
                const commentStat = document.createElement("span");
                commentStat.className = "post-stat";
                commentStat.innerHTML = `💬 ${post.comments || 0} comments`;
                postFooter.appendChild(commentStat);

                // Boutons de réaction
                const reactionButtons = createReactionButtons(
                    post.id,
                    post.likes || 0,
                    post.dislikes || 0,
                    isLoggedIn,
                    post.userReaction || null
                );
                postFooter.appendChild(reactionButtons);

                postElement.appendChild(postHeader);
                postElement.appendChild(titleElement);
                postElement.appendChild(contentElement);
                postElement.appendChild(postFooter);
                postsContainer.appendChild(postElement);
            });

            // Activer les événements de clic sur les posts
            EventPosts();
        })
        .catch((error) => {
            console.error("Error fetching posts:", error);
        });
}


function formatDate(date) {
    const now = new Date();
    const diff = now - date;
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return "Just now";
    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    if (days < 7) return `${days}d ago`;
    return date.toLocaleDateString();
}