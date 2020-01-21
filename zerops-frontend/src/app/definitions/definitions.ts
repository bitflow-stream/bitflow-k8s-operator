const dataSourceMap: Map<string, DataSource> = new Map();
const stepMap: Map<string, Step> = new Map();
const podMap: Map<string, Pod> = new Map();

export {dataSourceMap, stepMap, podMap};

export interface Step {
  name: string;
  podNames: string[];
}

export interface DataSource {
  name: string;
  creatorPodName?: string;
}

export interface Pod {
  name: string;
  creatorStepName: string;
  creatorDataSourceName: string;
}

export interface D3Node {
  id: string;
  text: string;
  x: number;
  y: number;
  width: number;
  height: number;
  type: string;
}

export interface D3Edge {
  start: string;
  stop: string;
}
