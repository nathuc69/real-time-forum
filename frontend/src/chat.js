class WebSocketClient {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.listeners = {};
    }

    connect() {
        return new Promise((resolve, reject) => {
            try {
                this.ws = new WebSocket(this.url);

                this.ws.onopen = () => {
                    console.log('Connected to server');
                    resolve();
                };

                this.ws.onmessage = (event) => {
                    const data = JSON.parse(event.data);
                    this.emit(data.type, data.payload);
                };

                this.ws.onerror = (error) => {
                    console.error('WebSocket error:', error);
                    reject(error);
                };

                this.ws.onclose = () => {
                    console.log('Disconnected from server');
                };
            } catch (error) {
                reject(error);
            }
        });
    }

    send(type, payload) {
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({ type, payload }));
        }
    }

    on(event, callback) {
        if (!this.listeners[event]) {
            this.listeners[event] = [];
        }
        this.listeners[event].push(callback);
    }

    emit(event, data) {
        if (this.listeners[event]) {
            this.listeners[event].forEach(callback => callback(data));
        }
    }

    disconnect() {
        this.ws?.close();
    }
}

// Usage
const client = new WebSocketClient('ws://localhost:8080/ws');
client.connect();
client.on('message', (data) => console.log(data));
client.send('message', { text: 'Hello' });