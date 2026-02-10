class WebSocketClient {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.listeners = {};
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000;
    }

    connect() {
        return new Promise((resolve, reject) => {
            try {
                // Les cookies sont automatiquement envoyés avec la requête WebSocket
                this.ws = new WebSocket(this.url);

                this.ws.onopen = () => {
                    console.log('✅ Connected to WebSocket server');
                    this.reconnectAttempts = 0;
                    resolve();
                };

                this.ws.onmessage = (event) => {
                    try {
                        const data = JSON.parse(event.data);
                        this.emit(data.type, data.payload);
                    } catch (error) {
                        console.error('Error parsing WebSocket message:', error);
                    }
                };

                this.ws.onerror = (error) => {
                    console.error('❌ WebSocket error:', error);
                    reject(error);
                };

                this.ws.onclose = () => {
                    console.log('🔌 Disconnected from WebSocket server');
                    this.emit('disconnected', {});
                    this.attemptReconnect();
                };
            } catch (error) {
                reject(error);
            }
        });
    }

    attemptReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`🔄 Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);

            setTimeout(() => {
                this.connect().catch(err => {
                    console.error('Reconnection failed:', err);
                });
            }, this.reconnectDelay * this.reconnectAttempts);
        } else {
            console.error('❌ Max reconnection attempts reached');
            this.emit('max_reconnect_failed', {});
        }
    }

    send(type, payload) {
        if (this.ws?.readyState === WebSocket.OPEN) {
            const message = {
                type,
                payload,
                timestamp: new Date().toISOString()
            };
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn('⚠️ WebSocket is not connected');
        }
    }

    on(event, callback) {
        if (!this.listeners[event]) {
            this.listeners[event] = [];
        }
        this.listeners[event].push(callback);
    }

    off(event, callback) {
        if (this.listeners[event]) {
            this.listeners[event] = this.listeners[event].filter(cb => cb !== callback);
        }
    }

    emit(event, data) {
        if (this.listeners[event]) {
            this.listeners[event].forEach(callback => callback(data));
        }
    }

    disconnect() {
        if (this.ws) {
            this.reconnectAttempts = this.maxReconnectAttempts; // Prevent reconnection
            this.ws.close();
            this.ws = null;
        }
    }
}

// Instance globale du client WebSocket
let wsClient = null;

// Initialiser la connexion WebSocket
function initWebSocket() {
    // Vérifier que l'utilisateur est connecté via localStorage
    const user = localStorage.getItem('user');
    if (!user) {
        console.warn('No user found, cannot connect to WebSocket');
        return null;
    }

    const wsUrl = `ws://localhost:8086/ws`;
    wsClient = new WebSocketClient(wsUrl);

    // Écouter les événements
    wsClient.on('online_users', handleOnlineUsers);
    wsClient.on('user_status', handleUserStatus);
    wsClient.on('chat_message', handleChatMessage);
    wsClient.on('typing', handleTyping);

    // Connecter
    wsClient.connect()
        .then(() => {
            console.log('✅ WebSocket connected successfully');
        })
        .catch(err => {
            console.error('❌ Failed to connect to WebSocket:', err);
        });

    return wsClient;
}

// Gestionnaires d'événements
function handleOnlineUsers(users) {
    console.log('Online users:', users);
    // Mettre à jour l'UI avec la liste des utilisateurs en ligne
    updateOnlineUsersList(users);
}

function handleUserStatus(status) {
    console.log('User status changed:', status);
    // Mettre à jour le statut d'un utilisateur dans l'UI
    updateUserStatus(status.userId, status.isOnline);
}

function handleChatMessage(message) {
    console.log('New chat message:', message);
    // Ajouter le message à la conversation
    addMessageToConversation(message);
}

function handleTyping(data) {
    console.log('User is typing:', data);
    // Afficher l'indicateur de frappe
    showTypingIndicator(data.userId, data.username);
}

// Fonctions utilitaires pour l'UI
function updateOnlineUsersList(users) {
    const usersList = document.getElementById('online-users-list');
    if (!usersList) return;

    usersList.innerHTML = '';
    users.forEach(user => {
        const userElement = document.createElement('div');
        userElement.className = 'user-item';
        userElement.innerHTML = `
            <span class="status-indicator ${user.isOnline ? 'online' : 'offline'}"></span>
            <span class="username">${user.username}</span>
        `;
        userElement.onclick = () => openChat(user.userId, user.username);
        usersList.appendChild(userElement);
    });
}

function updateUserStatus(userId, isOnline) {
    const userElements = document.querySelectorAll(`[data-user-id="${userId}"]`);
    userElements.forEach(element => {
        const statusIndicator = element.querySelector('.status-indicator');
        if (statusIndicator) {
            statusIndicator.className = `status-indicator ${isOnline ? 'online' : 'offline'}`;
        }
    });
}

function addMessageToConversation(message) {
    const messagesContainer = document.getElementById('messages-container');
    if (!messagesContainer) return;

    const messageElement = document.createElement('div');
    const currentUserId = getCurrentUserId();
    const isSent = message.senderId === currentUserId;

    messageElement.className = `message ${isSent ? 'sent' : 'received'}`;
    messageElement.innerHTML = `
        <div class="message-content">
            <p>${escapeHtml(message.content)}</p>
            <span class="message-time">${formatTime(message.timestamp || new Date())}</span>
        </div>
    `;

    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

function showTypingIndicator(userId, username) {
    const typingIndicator = document.getElementById('typing-indicator');
    if (!typingIndicator) return;

    typingIndicator.textContent = `${username} is typing...`;
    typingIndicator.style.display = 'block';

    // Masquer après 3 secondes
    setTimeout(() => {
        typingIndicator.style.display = 'none';
    }, 3000);
}

// Fonctions d'envoi de messages
function sendChatMessage(receiverId, content) {
    if (!wsClient || !content.trim()) return;

    const currentUserId = getCurrentUserId();
    wsClient.send('chat_message', {
        senderId: currentUserId,
        receiverId: receiverId,
        content: content.trim()
    });
}

function sendTypingIndicator(receiverId) {
    if (!wsClient) return;

    wsClient.send('typing', {
        receiverId: receiverId
    });
}

function markMessagesAsRead(senderId) {
    if (!wsClient) return;

    wsClient.send('mark_read', {
        senderId: senderId
    });
}

// Fonctions utilitaires
function getCurrentUserId() {
    // Récupérer l'ID de l'utilisateur actuel depuis le localStorage ou une variable globale
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

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatTime(date) {
    if (typeof date === 'string') {
        date = new Date(date);
    }
    return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
}

// Exporter pour utilisation globale
window.initWebSocket = initWebSocket;
window.wsClient = wsClient;
window.sendChatMessage = sendChatMessage;
window.sendTypingIndicator = sendTypingIndicator;
window.markMessagesAsRead = markMessagesAsRead;