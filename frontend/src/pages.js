// Page rendering functions
import { setupEventListeners } from './login.js';
import { handlePosts } from './posts.js';
import { createReactionButtons } from './reaction.js';


export function renderHome(loggedIn, username) {
    const filterHTML = `
        <div class="filters" style="margin: 20px 0; display: flex; justify-content: center; gap: 10px;">
            <select id="sortBy" style="padding: 8px; border-radius: 4px; border: 1px solid #ddd;">
                <option value="recent">Newest First</option>
                <option value="oldest">Oldest First</option>
                <option value="likes">Most Liked</option>
            </select>
            <select id="categoryFilter" style="padding: 8px; border-radius: 4px; border: 1px solid #ddd;">
                <option value="">All Categories</option>
            </select>
            <button id="filterBtn" style="padding: 8px 16px; background-color: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer;">Filter</button>
        </div>
    `;

    if (!loggedIn) {
        document.body.innerHTML = `
        <div id="MenuPage">
            <h1>Welcome to the Real-Time Forum</h1>
            <button id="loginBtn">Login</button>
            <button id="registerBtn">Register</button>
            ${filterHTML}
            <div id="postsContainer"></div>
        </div>
        
    `;
        loadCategories();
        handlePosts();

    } else {
        document.body.innerHTML = `
        <div id="logoutOverlay"></div>
        <div id="MenuPage">
            <h1>Welcome to the Real-Time Forum, ${username}</h1>
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
                <div style="display: flex; gap: 10px;">
                    <button id="createPostBtn" style="padding: 10px 20px; background-color: #2196F3; color: white; border: none; border-radius: 4px; cursor: pointer;">+ Create Post</button>
                    <button id="chatBtn" style="padding: 10px 20px; background-color: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer;">💬 Chat</button>
                </div>
                <button id="LogoutBtn">Logout</button>
            </div>
            
            <div id="createPostForm" style="display: none; background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin-bottom: 20px;">
                <h3 style="margin-top: 0;">Create a New Post</h3>
                <form id="newPostForm">
                    <div style="margin-bottom: 15px;">
                        <label style="display: block; margin-bottom: 5px;">Title:</label>
                        <input type="text" id="postTitle" required style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">
                    </div>
                    <div style="margin-bottom: 15px;">
                        <label style="display: block; margin-bottom: 5px;">Content:</label>
                        <textarea id="postContent" required rows="4" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;"></textarea>
                    </div>
                    <div style="margin-bottom: 15px;">
                        <label style="display: block; margin-bottom: 5px;">Categories (comma separated):</label>
                        <input type="text" id="postCategories" placeholder="tech, news, general" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">
                    </div>
                    <div style="text-align: right;">
                        <button type="button" id="cancelPostBtn" style="padding: 8px 16px; margin-right: 10px; background: #ddd; border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
                        <button type="submit" style="padding: 8px 16px; background: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer;">Publish</button>
                    </div>
                </form>
            </div>

            ${filterHTML}
            <div id="postsContainer"></div>
        </div>
        <div id="logoutPopup">
            <h3>Are you sure you want to logout?</h3>
            <button id="SubmitLogoutBtn">Yes, Logout</button>
            <button id="cancelLogoutBtn">Cancel</button>
        </div>
    `;
        setupEventListeners();

        // Setup Create Post listeners
        const createBtn = document.getElementById('createPostBtn');
        const form = document.getElementById('createPostForm');
        const cancelBtn = document.getElementById('cancelPostBtn');
        const postForm = document.getElementById('newPostForm');

        createBtn.addEventListener('click', () => {
            form.style.display = 'block';
            createBtn.style.display = 'none';
        });

        cancelBtn.addEventListener('click', () => {
            form.style.display = 'none';
            createBtn.style.display = 'block';
            postForm.reset();
        });

        postForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const title = document.getElementById('postTitle').value;
            const content = document.getElementById('postContent').value;
            const categoriesStr = document.getElementById('postCategories').value;
            const categories = categoriesStr.split(',').map(c => c.trim()).filter(c => c);

            fetch('http://localhost:8086/api/posts/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ title, content, categories })
            })
                .then(res => {
                    if (!res.ok) throw new Error('Failed to create post');
                    return res.json();
                })
                .then(() => {
                    form.style.display = 'none';
                    createBtn.style.display = 'block';
                    postForm.reset();
                    handlePosts(); // Refresh posts
                })
                .catch(err => {
                    console.error(err);
                    alert('Error creating post: ' + err.message);
                });
        });

        loadCategories();
        handlePosts();
    }

    document.getElementById('filterBtn')?.addEventListener('click', () => {
        const sortBy = document.getElementById('sortBy').value;
        const category = document.getElementById('categoryFilter').value;
        handlePosts(sortBy, category);
    });
}

function loadCategories() {
    fetch('http://localhost:8086/api/categories')
        .then(res => res.json())
        .then(categories => {
            const select = document.getElementById('categoryFilter');
            if (select) {
                // Clear existing options except "All Categories"
                select.innerHTML = '<option value="">All Categories</option>';
                categories.forEach(cat => {
                    const option = document.createElement('option');
                    option.value = cat;
                    option.textContent = cat;
                    select.appendChild(option);
                });
            }
        })
        .catch(err => console.error('Error loading categories:', err));
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

export function PopupLogout() {
    `
    <div id="logoutPopup">
    <h3>Are you sure you want to logout?</h3>
    <button id="SubmitLogoutBtn">Yes, Logout</button>
    <button id="cancelLogoutBtn">Cancel</button>
    </div>`;
    setupEventListeners();
}

export function renderPostDetails(params, isLoggedIn = false, username = '') {
    const postId = params?.id;

    document.body.innerHTML = `
        <div id="logoutOverlay"></div>
        <div id="postDetailsPage">
            <button id="backBtn" style="margin: 20px; padding: 10px 20px; cursor: pointer;">← Back to Home</button>
            <button id="LogoutBtn" style="margin: 20px; padding: 10px 20px; cursor: pointer; float: right; display: ${isLoggedIn ? 'block' : 'none'};">Logout</button>
            <div id="postDetailsContainer" style="max-width: 900px; margin: 0 auto; padding: 20px;">
                <div style="text-align: center; padding: 40px;">
                    <p>Loading post...</p>
                </div>
            </div>
        </div>
        <div id="logoutPopup">
            <h3>Are you sure you want to logout?</h3>
            <button id="SubmitLogoutBtn">Yes, Logout</button>
            <button id="cancelLogoutBtn">Cancel</button>
        </div>
    `;

    // Récupérer les détails du post
    if (postId) {
        fetch(`http://localhost:8086/api/posts/${postId}`)
            .then((response) => {
                if (!response.ok) {
                    throw new Error('Post not found');
                }
                return response.json();
            })
            .then((post) => {
                const container = document.getElementById('postDetailsContainer');

                // Créer le contenu principal du post
                const postDetailDiv = document.createElement('div');
                postDetailDiv.className = 'post-detail';
                postDetailDiv.innerHTML = `
                    <div class="post-header">
                        <div class="post-avatar">${(post.username || '?')[0].toUpperCase()}</div>
                        <div class="post-author-info">
                            <span class="post-author">${post.username || 'Anonymous'}</span>
                            <span class="post-date">${new Date(post.createdAt).toLocaleString()}</span>
                        </div>
                    </div>
                    <h1 class="post-detail-title">${post.title}</h1>
                    ${post.categories && post.categories.length > 0 ?
                        `<div class="post-categories" style="margin-bottom: 20px;">
                            ${post.categories.map(cat => `<span class="category-tag">${cat}</span>`).join('')}
                        </div>`
                        : ''}
                    <div class="post-detail-content">${post.content}</div>
                `;

                // Créer le footer avec les réactions
                const postFooter = document.createElement('div');
                postFooter.className = 'post-footer';

                // Compteur de commentaires
                const commentStat = document.createElement('span');
                commentStat.className = 'post-stat';
                commentStat.innerHTML = `💬 ${post.comments || 0} comments`;
                postFooter.appendChild(commentStat);

                // Boutons de réaction
                const reactionButtons = createReactionButtons(
                    post.id,
                    post.likes || 0,
                    post.dislikes || 0,
                    isLoggedIn,
                    post.islikeordislike || null
                );
                postFooter.appendChild(reactionButtons);

                postDetailDiv.appendChild(postFooter);
                container.innerHTML = '';
                container.appendChild(postDetailDiv);

                if (post.commentsList && post.commentsList.length > 0) {
                    const commentsContainer = document.createElement('div');
                    commentsContainer.id = 'commentsContainer';
                    commentsContainer.style.marginTop = '40px';
                    commentsContainer.innerHTML = '<h2 style="font-size: 1.5em; margin-bottom: 20px; color: #333;">Comments</h2>';
                    document.getElementById('postDetailsContainer').appendChild(commentsContainer);
                    post.commentsList.forEach(comment => {
                        const commentElement = document.createElement('div');
                        commentElement.className = 'comment';
                        commentElement.style.cssText = `
                        background-color: #f9f9f9;
                        border-left: 4px solid #4CAF50;
                        padding: 15px;
                        margin-bottom: 15px;
                        border-radius: 4px;
                        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                        transition: box-shadow 0.3s ease;
                    `;
                        commentElement.onmouseover = function () { this.style.boxShadow = '0 4px 8px rgba(0,0,0,0.15)'; };
                        commentElement.onmouseout = function () { this.style.boxShadow = '0 2px 4px rgba(0,0,0,0.1)'; };
                        if (comment.username == username) {
                            commentElement.innerHTML = `
                        <div class="comment-header" style="margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center;">
                            <strong style="color: #11910bff; font-size: 1.05em;">${'Vous'}</strong> 
                            <span style="color: #999; font-size: 0.85em;">${new Date(comment.createdAt).toLocaleString()}</span>
                        </div>
                        <div class="comment-content" style="color: #555; line-height: 1.6;">${comment.content}</div>
                    `;
                        } else {
                            commentElement.innerHTML = `
                        <div class="comment-header" style="margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center;">
                            <strong style="color: #333; font-size: 1.05em;">${comment.username || 'Anonymous'}</strong> 
                            <span style="color: #999; font-size: 0.85em;">${new Date(comment.createdAt).toLocaleString()}</span>
                        </div>
                        <div class="comment-content" style="color: #555; line-height: 1.6;">${comment.content}</div>
                    `;
                        }
                        commentsContainer.appendChild(commentElement);
                    });
                }
                //#region Formulaire d'ajout de commentaire
                // N'afficher le formulaire que si l'utilisateur est connecté
                if (isLoggedIn) {
                    const commentForm = document.createElement('div');
                    commentForm.id = 'commentForm';
                    commentForm.style.marginTop = '40px';
                    commentForm.innerHTML = `
                    <h2 style="font-size: 1.5em; margin-bottom: 20px; color: #333;">Add a Comment</h2>
                    <textarea id="commentContent" rows="4" style="width: 100%; padding: 10px; border: 1px solid #ccc; border-radius: 4px; resize: vertical;" placeholder="Write your comment here..."></textarea>
                    <button id="submitCommentBtn" style="margin-top: 10px; padding: 10px 20px; background-color: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer;">Submit Comment</button>
                `;
                    document.getElementById('postDetailsContainer').appendChild(commentForm);

                    document.getElementById('submitCommentBtn').addEventListener('click', (e) => {
                        e.preventDefault();
                        const commentContent = document.getElementById('commentContent').value;
                        if (commentContent.trim() === '') {
                            alert('Comment cannot be empty');
                            return;
                        }
                        fetch(`http://localhost:8086/api/posts/${postId}/comments`, {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                            },
                            credentials: 'include',
                            body: JSON.stringify({ content: commentContent }),
                        })
                            .then((response) => {
                                if (!response.ok) {
                                    if (response.status === 401) {
                                        throw new Error('You must be logged in to post a comment');
                                    }
                                    throw new Error('Failed to submit comment');
                                }
                                return response.json();
                            })
                            .then((newComment) => {
                                // Ajouter le nouveau commentaire à la liste des commentaires affichés
                                const commentsContainer = document.getElementById('commentsContainer');
                                if (commentsContainer) {
                                    const commentElement = document.createElement('div');
                                    commentElement.className = 'comment';
                                    commentElement.style.cssText = `
                                background-color: #f9f9f9;
                                border-left: 4px solid #4CAF50;
                                padding: 15px;
                                margin-bottom: 15px;
                                border-radius: 4px;
                                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                                transition: box-shadow 0.3s ease;
                            `;
                                    commentElement.onmouseover = function () { this.style.boxShadow = '0 4px 8px rgba(0,0,0,0.15)'; };
                                    commentElement.onmouseout = function () { this.style.boxShadow = '0 2px 4px rgba(0,0,0,0.1)'; };
                                    commentElement.innerHTML = `
                                <div class="comment-header" style="margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center;">
                                    <strong style="color: #333; font-size: 1.05em;">${newComment.username || 'Anonymous'}</strong> 
                                    <span style="color: #999; font-size: 0.85em;">${new Date(newComment.createdAt).toLocaleString()}</span>
                                </div>
                                <div class="comment-content" style="color: #555; line-height: 1.6;">${newComment.content}</div>
                            `;
                                    commentsContainer.appendChild(commentElement);
                                }
                                document.getElementById('commentContent').value = '';
                            })
                            .catch((error) => {
                                console.error('Error submitting comment:', error);
                                alert('Error submitting comment: ' + error.message);
                            });
                    });
                } else {
                    // Afficher un message pour les utilisateurs non connectés
                    const loginMessage = document.createElement('div');
                    loginMessage.style.cssText = 'margin-top: 40px; padding: 20px; background-color: #f0f0f0; border-radius: 4px; text-align: center;';
                    loginMessage.innerHTML = '<p>Please <a href="#/login" style="color: #4CAF50; text-decoration: underline;">login</a> to post a comment.</p>';
                    document.getElementById('postDetailsContainer').appendChild(loginMessage);
                }
                //#endregion
            })
            .catch((error) => {
                console.error('Error fetching post details:', error);
                const container = document.getElementById('postDetailsContainer');
                container.innerHTML = `
                    <div style="text-align: center; padding: 40px; color: #ff6b6b;">
                        <h2>Error loading post</h2>
                        <p>${error.message}</p>
                    </div>
                `;
            });
    }

    // Appeler setupEventListeners pour gérer les événements de logout
    setupEventListeners();
}

export function renderChat(isLoggedIn, username) {
    if (!isLoggedIn) {
        renderLogin();
        return;
    }

    document.body.innerHTML = `
        <div id="chatPage" style="display: flex; height: 100vh; overflow: hidden;">
            <!-- Sidebar avec la liste des utilisateurs -->
            <div id="chatSidebar" style="width: 300px; background: #f5f5f5; border-right: 1px solid #ddd; display: flex; flex-direction: column;">
                <div style="padding: 20px; background: #4CAF50; color: white;">
                    <h2 style="margin: 0; font-size: 1.5em;">💬 Chat</h2>
                    <p style="margin: 5px 0 0 0; font-size: 0.9em;">Welcome, ${username}</p>
                </div>
                
                <div style="padding: 10px;">
                    <button id="backToHomeBtn" style="width: 100%; padding: 10px; background: #2196F3; color: white; border: none; border-radius: 4px; cursor: pointer; margin-bottom: 10px;">
                        ← Back to Home
                    </button>
                </div>

                <div style="padding: 10px; border-bottom: 1px solid #ddd;">
                    <input type="text" id="userSearch" placeholder="Search users..." style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box;">
                </div>

                <div id="chatUsersList" style="flex: 1; overflow-y: auto; padding: 10px;">
                    <div style="text-align: center; padding: 20px; color: #999;">
                        Loading users...
                    </div>
                </div>
            </div>

            <!-- Zone de conversation -->
            <div id="chatMain" style="flex: 1; display: flex; flex-direction: column; background: #fff;">
                <div id="chatHeader" style="padding: 20px; background: #f9f9f9; border-bottom: 1px solid #ddd; display: none;">
                    <h3 id="chatUsername" style="margin: 0; color: #333;"></h3>
                    <p id="chatUserStatus" style="margin: 5px 0 0 0; font-size: 0.9em; color: #999;"></p>
                </div>

                <div id="messagesContainer" style="flex: 1; overflow-y: auto; padding: 20px; background: #fafafa;">
                    <div style="text-align: center; padding: 40px; color: #999;">
                        <h3>Select a user to start chatting</h3>
                        <p>Choose someone from the list on the left</p>
                    </div>
                </div>

                <div id="typingIndicator" style="padding: 10px 20px; color: #999; font-style: italic; display: none;">
                    Someone is typing...
                </div>

                <div id="messageInputContainer" style="padding: 20px; background: #f9f9f9; border-top: 1px solid #ddd; display: none;">
                    <div style="display: flex; gap: 10px;">
                        <textarea id="messageInput" placeholder="Type your message..." style="flex: 1; padding: 10px; border: 1px solid #ddd; border-radius: 4px; resize: none; font-family: inherit;" rows="2"></textarea>
                        <button id="sendMessageBtn" style="padding: 10px 20px; background: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                            Send
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `;

    // Initialiser le WebSocket seulement s'il n'existe pas déjà
    let ws = window.wsClient;
    if (!ws || ws.ws?.readyState !== WebSocket.OPEN) {
        console.log('🔌 Initializing WebSocket connection...');
        ws = window.initWebSocket();
    } else {
        console.log('✅ WebSocket already connected, reusing existing connection');
    }

    // Variables pour suivre la conversation actuelle
    let currentChatUserId = null;
    let currentChatUsername = null;
    let messageOffset = 0;
    const messageLimit = 10;

    // Charger la liste des utilisateurs
    loadChatUsers();

    // Gestionnaire de recherche d'utilisateurs
    document.getElementById('userSearch').addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        const userItems = document.querySelectorAll('.chat-user-item');
        userItems.forEach(item => {
            const username = item.querySelector('.chat-user-name').textContent.toLowerCase();
            item.style.display = username.includes(searchTerm) ? 'flex' : 'none';
        });
    });

    // Gestionnaire du bouton retour
    document.getElementById('backToHomeBtn').addEventListener('click', () => {
        if (window.wsClient) {
            window.wsClient.disconnect();
        }
        window.location.hash = '#/';
    });

    // Gestionnaire d'envoi de message
    document.getElementById('sendMessageBtn').addEventListener('click', sendMessage);
    document.getElementById('messageInput').addEventListener('keypress', (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        } else if (e.key !== 'Enter') {
            // Envoyer l'indicateur de frappe
            if (currentChatUserId) {
                window.sendTypingIndicator(currentChatUserId);
            }
        }
    });

    // Fonction pour charger la liste des utilisateurs
    function loadChatUsers() {
        fetch('http://localhost:8086/api/users', {
            credentials: 'include'
        })
            .then(res => res.json())
            .then(users => {
                const usersList = document.getElementById('chatUsersList');
                usersList.innerHTML = '';

                if (users.length === 0) {
                    usersList.innerHTML = '<div style="text-align: center; padding: 20px; color: #999;">No users available</div>';
                    return;
                }

                users.forEach(user => {
                    const userItem = createUserItem(user);
                    usersList.appendChild(userItem);
                });
            })
            .catch(err => {
                console.error('Error loading users:', err);
                document.getElementById('chatUsersList').innerHTML = '<div style="text-align: center; padding: 20px; color: #f44336;">Error loading users</div>';
            });
    }

    // Fonction pour créer un élément utilisateur
    function createUserItem(user) {
        const userItem = document.createElement('div');
        userItem.className = 'chat-user-item';
        userItem.setAttribute('data-user-id', user.id);
        userItem.style.cssText = `
            display: flex;
            align-items: center;
            padding: 12px;
            cursor: pointer;
            border-radius: 4px;
            margin-bottom: 5px;
            transition: background 0.2s;
        `;

        userItem.innerHTML = `
            <div style="position: relative; margin-right: 12px;">
                <div style="width: 40px; height: 40px; border-radius: 50%; background: #4CAF50; color: white; display: flex; align-items: center; justify-content: center; font-weight: bold;">
                    ${user.username[0].toUpperCase()}
                </div>
                <div class="status-indicator ${user.isOnline ? 'online' : 'offline'}" style="position: absolute; bottom: 0; right: 0; width: 12px; height: 12px; border-radius: 50%; border: 2px solid #f5f5f5; background: ${user.isOnline ? '#4CAF50' : '#999'};"></div>
            </div>
            <div style="flex: 1;">
                <div class="chat-user-name" style="font-weight: 500; color: #333;">${user.username}</div>
                <div class="chat-user-status" style="font-size: 0.85em; color: #999;">${user.isOnline ? 'Online' : 'Offline'}</div>
            </div>
        `;

        userItem.addEventListener('mouseenter', () => {
            userItem.style.background = '#e8e8e8';
        });

        userItem.addEventListener('mouseleave', () => {
            userItem.style.background = 'transparent';
        });

        userItem.addEventListener('click', () => {
            openChat(user.id, user.username);
        });

        return userItem;
    }

    // Fonction pour ouvrir une conversation
    function openChat(userId, username) {
        currentChatUserId = userId;
        currentChatUsername = username;
        messageOffset = 0;

        // Mettre à jour l'interface
        document.getElementById('chatHeader').style.display = 'block';
        document.getElementById('messageInputContainer').style.display = 'block';
        document.getElementById('chatUsername').textContent = username;

        // Charger la conversation
        loadConversation(userId);

        // Marquer les messages comme lus
        window.markMessagesAsRead(userId);

        // Mettre en surbrillance l'utilisateur sélectionné
        document.querySelectorAll('.chat-user-item').forEach(item => {
            item.style.background = 'transparent';
        });
        document.querySelector(`[data-user-id="${userId}"]`).style.background = '#e0f7fa';
    }

    // Fonction pour charger une conversation
    function loadConversation(userId) {
        fetch(`http://localhost:8086/api/chat/conversation?userId=${userId}&limit=${messageLimit}&offset=${messageOffset}`, {
            credentials: 'include'
        })
            .then(res => res.json())
            .then(messages => {
                const container = document.getElementById('messagesContainer');
                container.innerHTML = '';

                if (messages.length === 0) {
                    container.innerHTML = '<div style="text-align: center; padding: 40px; color: #999;">No messages yet. Start the conversation!</div>';
                    return;
                }

                messages.forEach(message => {
                    addMessageToUI(message);
                });

                // Scroll vers le bas
                container.scrollTop = container.scrollHeight;
            })
            .catch(err => {
                console.error('Error loading conversation:', err);
            });
    }

    // Fonction pour ajouter un message à l'UI
    function addMessageToUI(message) {
        const container = document.getElementById('messagesContainer');
        const currentUserId = getCurrentUserId();
        const isSent = message.senderId === currentUserId;

        const messageDiv = document.createElement('div');
        messageDiv.className = isSent ? 'message sent' : 'message';
        messageDiv.style.cssText = `
            display: flex;
            justify-content: ${isSent ? 'flex-end' : 'flex-start'};
            margin-bottom: 15px;
        `;

        messageDiv.innerHTML = `
            <div style="max-width: 60%; padding: 10px 15px; border-radius: 18px; background: ${isSent ? '#4CAF50' : '#e0e0e0'}; color: ${isSent ? 'white' : '#333'};">
                <p style="margin: 0; word-wrap: break-word;">${escapeHtml(message.content)}</p>
                <span style="font-size: 0.75em; opacity: 0.8; margin-top: 5px; display: block;">
                    ${formatTime(message.createdAt)}
                </span>
            </div>
        `;

        container.appendChild(messageDiv);
    }

    // Fonction pour envoyer un message
    function sendMessage() {
        const input = document.getElementById('messageInput');
        const content = input.value.trim();

        if (!content || !currentChatUserId) return;

        window.sendChatMessage(currentChatUserId, content);
        input.value = '';
    }

    // Fonction pour obtenir l'ID de l'utilisateur actuel
    function getCurrentUserId() {
        const userStr = localStorage.getItem('user');
        if (userStr) {
            try {
                const user = JSON.parse(userStr);
                return user.id;
            } catch (e) {
                console.error('Error parsing user data:', e);
            }
        }
        return null;
    }

    // Fonctions utilitaires
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    function formatTime(dateStr) {
        const date = new Date(dateStr);
        return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
    }

    // Écouter les nouveaux messages via WebSocket
    if (ws) {
        // Retirer les ANCIENS listeners s'ils existent (évite les doublons à chaque renderChat)
        if (window._chatHandlers) {
            ws.off('chat_message', window._chatHandlers.onMessage);
            ws.off('user_status', window._chatHandlers.onStatus);
        }

        // Définir les handlers
        function handleChatMessageEvent(message) {
            // Si le message concerne la conversation actuelle, l'ajouter
            if (currentChatUserId &&
                (message.senderId === currentChatUserId || message.receiverId === currentChatUserId)) {
                addMessageToUI(message);

                // Scroll vers le bas
                const container = document.getElementById('messagesContainer');
                if (container) {
                    container.scrollTop = container.scrollHeight;
                }

                // Marquer comme lu si c'est un message reçu
                if (message.senderId === currentChatUserId) {
                    window.markMessagesAsRead(currentChatUserId);
                }
            }
        }

        function handleUserStatusEvent(status) {
            // Mettre à jour le statut dans la liste sans re-render
            const userItem = document.querySelector(`[data-user-id="${status.userId}"]`);
            if (userItem) {
                const statusIndicator = userItem.querySelector('.status-indicator');
                const statusText = userItem.querySelector('.chat-user-status');

                if (statusIndicator) {
                    statusIndicator.className = `status-indicator ${status.isOnline ? 'online' : 'offline'}`;
                    statusIndicator.style.background = status.isOnline ? '#4CAF50' : '#999';
                }

                if (statusText) {
                    statusText.textContent = status.isOnline ? 'Online' : 'Offline';
                }
            }

            // Mettre à jour le header si c'est la conversation actuelle
            if (currentChatUserId === status.userId) {
                const statusEl = document.getElementById('chatUserStatus');
                if (statusEl) {
                    statusEl.textContent = status.isOnline ? '🟢 Online' : '⚫ Offline';
                }
            }
        }

        // Sauvegarder les références pour pouvoir les retirer au prochain appel de renderChat
        window._chatHandlers = {
            onMessage: handleChatMessageEvent,
            onStatus: handleUserStatusEvent
        };

        // Ajouter les nouveaux listeners
        ws.on('chat_message', handleChatMessageEvent);
        ws.on('user_status', handleUserStatusEvent);
    }
}