export type Packet = {
  id: number;
  name: string;
  hexValue: string;
  count: number;
  cycleTime: number;
  values: Map<string, string>;
};
