import { Measurement, Packet, Packets } from "./modals";

var packets = new Map<number, Packet>();

export function updateReceiveTable(packets: Packets) {
  let tableBody = document.getElementById(
    "tableBody"
  ) as HTMLTableSectionElement;
  tableBody.textContent = "";
  for (let packet in packets) {
    console.log(packets[packet].HexValue)
    let packetRow = createPacketRow(packets[packet]);
    tableBody.append(packetRow);
    let measurementsRows = createMeasurementRows(packets[packet].Measurements);
    for (let row of measurementsRows) {
      tableBody.append(row);
    }
  }
}

function createPacketRow(packet: Packet): HTMLTableRowElement {
  let row = document.createElement("tr");
  let id_td = document.createElement("td");
  id_td.innerHTML = packet.Id.toString(10);
  let name_td = document.createElement("td");
  name_td.innerHTML = packet.Name.toString();
  let hexValue_td = document.createElement("td");
  hexValue_td.innerHTML = packet.HexValue.toString(16);
  let count_td = document.createElement("td");
  count_td.innerHTML = packet.Count.toString();
  let cycleTime_td = document.createElement("td");
  cycleTime_td.innerHTML = packet.CycleTime.toString();

  row.append(id_td);
  row.append(name_td);
  row.append(hexValue_td);
  row.append(count_td);
  row.append(cycleTime_td);
  return row;
}

function createMeasurementRows(
  measurements: Measurement[]
): HTMLTableRowElement[] {
  let rows = [];
  for (let measurement of measurements) {
    let row = document.createElement("tr");
    let data = document.createElement("td");
    let dataString = `${measurement.Name}: ${measurement.Value} ${measurement.Units}`;
    data.innerHTML += dataString;
    data.setAttribute("colspan", "5");
    row.append(data);
    rows.push(row);
  }

  return rows;
}
