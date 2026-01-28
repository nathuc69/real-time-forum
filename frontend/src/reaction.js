// Gestion des réactions (likes/dislikes) pour les posts

/**
 * Récupère la réaction actuelle de l'utilisateur pour un post
 * @param {number} postId - L'ID du post
 * @returns {Promise<Object>} - La réaction de l'utilisateur
 */
export async function getUserReaction(postId) {
    try {
        const response = await fetch(`http://localhost:8086/api/posts/${postId}/reaction`, {
            method: 'GET',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                return { hasReacted: false, value: null };
            }
            throw new Error('Failed to fetch user reaction');
        }

        return await response.json();
    } catch (error) {
        console.error('Error fetching user reaction:', error);
        return { hasReacted: false, value: null };
    }
}

/**
 * Ajoute ou modifie une réaction
 * @param {number} postId - L'ID du post
 * @param {number} value - 1 pour like, -1 pour dislike
 * @returns {Promise<Object>} - La réponse du serveur
 */
export async function addReaction(postId, value) {
    try {
        const response = await fetch(`http://localhost:8086/api/posts/${postId}/reaction`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ value })
        });

        if (!response.ok) {
            if (response.status === 401) {
                throw new Error('You must be logged in to react to posts');
            }
            throw new Error('Failed to add reaction');
        }

        return await response.json();
    } catch (error) {
        console.error('Error adding reaction:', error);
        throw error;
    }
}

/**
 * Supprime une réaction
 * @param {number} postId - L'ID du post
 * @returns {Promise<Object>} - La réponse du serveur
 */
export async function removeReaction(postId) {
    try {
        const response = await fetch(`http://localhost:8086/api/posts/${postId}/reaction`, {
            method: 'DELETE',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            }
        });

        if (!response.ok) {
            throw new Error('Failed to remove reaction');
        }

        return await response.json();
    } catch (error) {
        console.error('Error removing reaction:', error);
        throw error;
    }
}

/**
 * Crée les boutons de réaction pour un post
 * @param {number} postId - L'ID du post
 * @param {number} likes - Nombre de likes
 * @param {number} dislikes - Nombre de dislikes
 * @param {boolean} isLoggedIn - Si l'utilisateur est connecté
 * @param {number|null} userReaction - La réaction actuelle de l'utilisateur (1, -1, ou null)
 * @returns {HTMLElement} - L'élément contenant les boutons de réaction
 */
export function createReactionButtons(postId, likes = 0, dislikes = 0, isLoggedIn = false, userReaction = null) {
    const reactionContainer = document.createElement('div');
    reactionContainer.className = 'reaction-container';
    reactionContainer.setAttribute('data-post-id', postId);

    const likeBtn = document.createElement('button');
    likeBtn.className = `reaction-btn like-btn ${userReaction === 1 ? 'active' : ''}`;
    likeBtn.innerHTML = `
        <span class="reaction-icon">👍</span>
        <span class="reaction-count">${likes}</span>
    `;
    likeBtn.setAttribute('data-reaction', '1');
    likeBtn.disabled = !isLoggedIn;

    const dislikeBtn = document.createElement('button');
    dislikeBtn.className = `reaction-btn dislike-btn ${userReaction === -1 ? 'active' : ''}`;
    dislikeBtn.innerHTML = `
        <span class="reaction-icon">👎</span>
        <span class="reaction-count">${dislikes}</span>
    `;
    dislikeBtn.setAttribute('data-reaction', '-1');
    dislikeBtn.disabled = !isLoggedIn;

    // Gestion des clics
    if (isLoggedIn) {
        likeBtn.addEventListener('click', async (e) => {
            e.stopPropagation(); // Empêcher la navigation vers le post
            await handleReactionClick(postId, 1, likeBtn, dislikeBtn, userReaction);
        });

        dislikeBtn.addEventListener('click', async (e) => {
            e.stopPropagation(); // Empêcher la navigation vers le post
            await handleReactionClick(postId, -1, likeBtn, dislikeBtn, userReaction);
        });
    } else {
        // Afficher un message pour les utilisateurs non connectés
        likeBtn.title = 'Please login to react';
        dislikeBtn.title = 'Please login to react';
    }

    reactionContainer.appendChild(likeBtn);
    reactionContainer.appendChild(dislikeBtn);

    return reactionContainer;
}

/**
 * Gère le clic sur un bouton de réaction
 * @param {number} postId - L'ID du post
 * @param {number} newValue - 1 pour like, -1 pour dislike
 * @param {HTMLElement} likeBtn - Le bouton like
 * @param {HTMLElement} dislikeBtn - Le bouton dislike
 * @param {number|null} currentReaction - La réaction actuelle
 */
async function handleReactionClick(postId, newValue, likeBtn, dislikeBtn, currentReaction) {
    try {
        let result;

        // Si l'utilisateur clique sur la même réaction, on la supprime
        if (currentReaction === newValue) {
            result = await removeReaction(postId);

            // Mettre à jour l'interface
            likeBtn.classList.remove('active');
            dislikeBtn.classList.remove('active');

            // Mettre à jour les compteurs
            updateReactionCounts(likeBtn, dislikeBtn, result.totalLikes, result.totalDislikes);

            // Mettre à jour la réaction actuelle
            currentReaction = null;
        } else {
            // Sinon, on ajoute/modifie la réaction
            result = await addReaction(postId, newValue);

            // Mettre à jour l'interface
            if (newValue === 1) {
                likeBtn.classList.add('active');
                dislikeBtn.classList.remove('active');
            } else {
                dislikeBtn.classList.add('active');
                likeBtn.classList.remove('active');
            }

            // Mettre à jour les compteurs
            updateReactionCounts(likeBtn, dislikeBtn, result.totalLikes, result.totalDislikes);

            // Mettre à jour la réaction actuelle
            currentReaction = newValue;
        }

        // Animation de feedback
        const activeBtn = newValue === 1 ? likeBtn : dislikeBtn;
        activeBtn.classList.add('reaction-pulse');
        setTimeout(() => activeBtn.classList.remove('reaction-pulse'), 600);

    } catch (error) {
        console.error('Error handling reaction:', error);
        alert(error.message);
    }
}

/**
 * Met à jour les compteurs de réactions
 * @param {HTMLElement} likeBtn - Le bouton like
 * @param {HTMLElement} dislikeBtn - Le bouton dislike
 * @param {number} likes - Nouveau nombre de likes
 * @param {number} dislikes - Nouveau nombre de dislikes
 */
function updateReactionCounts(likeBtn, dislikeBtn, likes, dislikes) {
    const likeCount = likeBtn.querySelector('.reaction-count');
    const dislikeCount = dislikeBtn.querySelector('.reaction-count');

    if (likeCount) likeCount.textContent = likes;
    if (dislikeCount) dislikeCount.textContent = dislikes;
}
