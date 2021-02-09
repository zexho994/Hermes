function connectWebSocket() {
    if (!"WebSocket" in window) {
        alert("sorry, your browser not support websocket.")
    }

    ws = new WebSocket('ws://localhost:994');

    ws.onopen = () => {
        console.log('connect server success')
        ws.send(JSON.stringify(msgData))
    }

    ws.onmessage = (msg) => {
        let message = JSON.parse(msg.data)
        console.log('receive msg : ', msg)
    }

    ws.onerror = () => {
        console.log('connect fail , reconnect ing ...')
        connectWebSocket();
    }

    ws.onclose = () => {
        console.log('connect close')
    }
}

connectWebSocket()