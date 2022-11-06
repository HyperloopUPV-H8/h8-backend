import { Packet, Packets } from "./modals";
import { updateReceiveTable } from "./receiveTableBuilder";

let dataSocket = new WebSocket("ws://127.0.0.1:4000/backend/data");

dataSocket.onopen = (ev) => {
  alert("Established WS connection");
};

dataSocket.onmessage = (ev) => {
  let packet: Packets = JSON.parse(ev.data) as Packet;
  updateReceiveTable(packet);
};
