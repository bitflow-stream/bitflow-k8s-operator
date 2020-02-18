const dataSourceMap: Map<string, DataSource> = new Map();
const stepMap: Map<string, Step> = new Map();
const podMap: Map<string, Pod> = new Map();
const currentDataSourcesMap: Map<string, DataSource> = new Map();
const currentStepsMap: Map<string, Step> = new Map();
const currentPodsMap: Map<string, Pod> = new Map();

export {dataSourceMap, stepMap, podMap, currentDataSourcesMap, currentStepsMap, currentPodsMap};

// FRONTEND
export interface FrontendData {
  nodes: D3Node[];
  edges: D3Edge[];
}

export interface D3Node {
  id: string;
  text: string;
  x: number;
  y: number;
  width: number;
  height: number;
  type: 'step' | 'data-source' | 'pod' | 'data-source-stack' | 'pod-stack';
}

export interface D3Edge {
  start: string;
  stop: string;
}

// BACKEND
export interface GraphVisualization {
  graphColumns: GraphVisualizationColumn[];
}

export interface GraphVisualizationColumn {
  maxFilledRow: number;
  graphElements: GraphElement[];
}

export interface GraphElement {
  type: 'step' | 'data-source' | 'pod' | 'data-source-stack' | 'pod-stack';
  step?: Step;
  dataSource?: DataSource
  pod?: Pod;
  dataSourceStack?: DataSourceStack
  podStack?: PodStack;
}

export interface DataSourceStack {
  stackId: string;
  dataSources: DataSource[];
}

export interface PodStack {
  stackId: string;
  pods: Pod[];
}

export interface Step {
  name: string;
  podType: 'pod' | 'pod-stack'
  pods?: Pod[]
  podStack?: PodStack
}

export interface DataSource {
  name: string;
  hasCreatorPod: boolean;
  creatorPod?: Pod;
  hasOutputName: boolean;
  outputName?: string;
}

export interface Pod {
  name: string;
  hasCreatorStep: boolean;
  creatorStep?: Step;
  creatorDataSources: DataSource[];
}

export interface StepDataSourceMatches {
  [key: string]: string[]
}

