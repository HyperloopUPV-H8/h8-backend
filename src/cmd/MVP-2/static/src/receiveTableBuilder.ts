import { Packet } from "./modals";

var globalPackets = new Map<number, Packet>();

export function updateReceiveTable(packet: Packet) {
  updatePackets(packet);
  let tableBody = document.getElementById(
    "tableBody"
  ) as HTMLTableSectionElement;
  tableBody.textContent = "";
  for (let [_, packet] of globalPackets) {
    addPacketToTable(tableBody, packet);
  }
}

function updatePackets(packet: Packet) {
  globalPackets.set(packet.id, packet);
}

function addPacketToTable(tableBody: HTMLTableSectionElement, packet: Packet) {
  let packetRow = createPacketRow(packet);
  tableBody.append(packetRow);
  let measurementsRows = createMeasurementRows(packet.values);
  for (let row of measurementsRows) {
    tableBody.append(row);
  }
}

function createPacketRow(packet: Packet): HTMLTableRowElement {
  let row = document.createElement("tr");
  let id_td = document.createElement("td");
  id_td.innerHTML = packet.id.toString(10);
  let name_td = document.createElement("td");
  name_td.innerHTML = "___";
  let hexValue_td = document.createElement("td");
  hexValue_td.innerHTML = packet.hexValue.toString();
  let count_td = document.createElement("td");
  count_td.innerHTML = packet.count.toString(10);
  let cycleTime_td = document.createElement("td");
  cycleTime_td.innerHTML = packet.cycleTime.toString(10);

  row.append(id_td);
  row.append(name_td);
  row.append(hexValue_td);
  row.append(count_td);
  row.append(cycleTime_td);
  return row;
}

function createMeasurementRows(
  measurements: Map<string, string>
): HTMLTableRowElement[] {
  console.log(measurements)
  let rows = [];
  for (let measurement in measurements) {
    let row = document.createElement("tr");
    let data = document.createElement("td");
    let dataString = `${measurement}: ${measurements.get(measurement)} ___`;
    data.innerHTML += dataString;
    data.setAttribute("colspan", "5");
    row.append(data);
    rows.push(row);
  }

  return rows;
}
