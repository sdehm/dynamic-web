import morphdom from "morphdom";

// connect to websocket
const socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = () => {
  console.log("Connected");
};

socket.onmessage = (event) => {
  const data = JSON.parse(event.data);
  const id = data.id;
  const html = data.html;
  morphdom(document.getElementById(id), html);
};
