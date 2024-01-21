import React, { useState, useEffect } from "react";

function App() {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/echo");

    socket.onopen = () => console.log("Connection established");
    socket.onmessage = (event) => {
      setMessages((messages) => [...messages, event.data]);
    };

    return () => {
      socket.close();
    };
  }, []);

  const handleChange = (event) => {
    setInput(event.target.value);
  };

  const sendMessage = () => {
    socket.send(input);
    setInput("");
  };

  return (
    <div className="content">
      <input type="text" value={input} onChange={handleChange} />
      <button onClick={sendMessage}>send message</button>
      <pre>{messages.join("\n")}</pre>
    </div>
  );
}

export default App;
