export type Packet = {
  Id: number;
  Name: string;
  HexValue: number;
  Count: number;
  CycleTime: number;
  Measurements: Measurement[];
};

export interface Packets {
    [key: number]: Packet
};

export type Measurement = {
  Name: string;
  Value: string;
  Units: string;
};
