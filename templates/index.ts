import morphdom from "morphdom";

// connect to websocket
const socket = new WebSocket("ws://localhost:8080/ws");

// initilize the connection id
let connectionId: string;

socket.onopen = () => {
  console.log("Connected");
};

// message types
enum MessageType {
  MorphData = "morph_data",
  Connected = "connected",
}

socket.onmessage = (event) => {
  const data = JSON.parse(event.data);
  const id = data.id;
  const html = data.html;
  const type = data.type as MessageType;
  switch (type) {
    case MessageType.MorphData:
      morphdom(document.getElementById(id) as Node, html);
      break;
    case MessageType.Connected:
      console.log("Connected: ", id);
      connectionId = id;
      break;
  }
};

// send mouse position to the server
window.onmousemove = (event: { clientX: any; clientY: any; }) => {
  socket.send(JSON.stringify({ id: connectionId, x: event.clientX, y: event.clientY }));
}