import { Packet } from "./modals";
import { updateReceiveTable } from "./receiveTableBuilder";

let dataSocket = new WebSocket("ws://127.0.0.1:4000/data");

dataSocket.onopen = (ev) => {
  alert("Established WS connection");
};

dataSocket.onmessage = (ev) => {
  let packetObject = JSON.parse(ev.data);
  updateReceiveTable(packetObject);
};
