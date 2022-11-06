var packets = new Map();
export function updateReceiveTable(packet) {
  packets.set(packet.id, packet);
  let tableBody = document.getElementById("tableBody");
  tableBody.textContent = "";
  for (let [id, packet2] of packets) {
    let packetRow = createPacketRow(packet2);
    let measurementsRow = createMeasurementsRow(packet2.measurements);
  }
}
function createPacketRow(packet) {
  let row = new HTMLTableRowElement();
  let id = new HTMLTableCellElement().innerHTML = packet.id.toString();
  let name = new HTMLTableCellElement().innerHTML = packet.name.toString();
  let hexValue = new HTMLTableCellElement().innerHTML = packet.hexValue.toString(16);
  let count = new HTMLTableCellElement().innerHTML = packet.cycleTime.toString();
  row.append(row);
  row.append(id);
  row.append(name);
  row.append(hexValue);
  row.append(count);
  return row;
}
function createMeasurementsRow(measurements) {
  let row = new HTMLTableRowElement();
  for (let measurement of measurements) {
    let dataString = `${measurement.name}: ${measurement.value} ${measurement.units}`;
    row.innerHTML += dataString + "\n";
  }
  return row;
}
