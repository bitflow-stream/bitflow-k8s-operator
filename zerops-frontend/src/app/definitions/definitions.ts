const dataSourceMap: Map<string, DataSource> = new Map();
const stepMap: Map<string, Step> = new Map();
const podMap: Map<string, Pod> = new Map();

export {dataSourceMap, stepMap, podMap};

export interface KubernetesNode {
  dataSource?: DataSource,
  step?: Step
}

export var kubernetesGraph: KubernetesNode[][] = [];

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
  creatorDataSourceNames: string[];
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

export interface VisualizationData {
  nodes: D3Node[];
  edges: D3Edge[];
}

export interface StepDataSourceMatches {
  [key: string]: string[]
}
