let webSocket = (() => {
  let socket = null;

  function connect(url) {
    if (!socket || socket.readyState === WebSocket.CLOSED) {
      socket = new WebSocket(url);

      socket.onopen = () => {
        console.log("WebSocket connected.");
      };

      socket.onclose = () => console.log("WebSocket disconnected.");
      socket.onerror = (error) => console.error("WebSocket error:", error);
    }
  }

  function sendMessage(message) {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(message));
    } else {
      console.error("WebSocket is not connected.");
    }
  }

  function onMessage(callback) {
    if (socket) {
      socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        callback(data);
      };
    }
  }

  function getSocket() {
    return socket;
  }

  return {
    connect,
    sendMessage,
    onMessage,
    getSocket,
  };
})();
