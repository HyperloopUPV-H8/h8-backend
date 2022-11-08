import { Packet } from "./modals";
import { updateReceiveTable } from "./receiveTableBuilder";

let dataSocket = new WebSocket("ws://127.0.0.1:4000/data");

dataSocket.onopen = (ev) => {
  alert("Established WS connection");
};

dataSocket.onmessage = (ev) => {
  let packetsObject = JSON.parse(ev.data);
  let packetMap = new Map<number, Packet>();
  for (let [key, value] of Object.entries(packetsObject)) {
    packetMap.set(Number.parseInt(key), value as Packet);
  }
  updateReceiveTable(packetMap);
};
