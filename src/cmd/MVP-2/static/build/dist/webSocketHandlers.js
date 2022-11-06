import {updateReceiveTable} from "./receiveTableBuilder.js";
let dataSocket = new WebSocket("ws://127.0.0.1:5000");
dataSocket.onopen = (ev) => {
  alert("Established WS connection");
};
dataSocket.onmessage = (ev) => {
  let packet = JSON.parse(ev.data);
  updateReceiveTable(packet);
};
