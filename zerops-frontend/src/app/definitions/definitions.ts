const dataSourceMap: Map<string, DataSource> = new Map();
const stepMap: Map<string, Step> = new Map();
const dataSourceGraphElementMap: Map<string, DataSourceGraphElement> = new Map();
const stepGraphElementMap: Map<string, StepGraphElement> = new Map();

export {dataSourceMap, stepMap, dataSourceGraphElementMap, stepGraphElementMap};

export enum AnalysisType {
  ALL_TO_ONE = 'all-to-one',
  ONE_TO_ONE = 'one-to-one'
}

export class KubernetesGraph {
  dataSourceGraphElements: DataSourceGraphElement[];
  stepGraphElements: StepGraphElement[];
}

export declare class DataSourceGraphElement {
  uuid: string;
  stepGraphElements: string[];
  creatorStepGraphElement: string;
}

export interface SourceDataSourceGraphElement {
  uuid: string;
  alreadyCreatedOutputFromThisDataSource: boolean;
}

export declare class StepGraphElement {
  uuid: string;
  outputDataSourceGraphElements: string[];
  sourceDataSourceGraphElements: string[];
}

export declare class DataSourceStepMatch {
  dataSource: string;
  step: string;
}

export declare class DataSourceLabelKeyValuePair {
  key: string;
  value: string;
}

export declare class StepKeyValuePair {
  regex: boolean;
  key: string;
  value: string;
}

export declare class KeyValuePair {
  key: string;
  value: string;
}

export declare class DataSource {
  uuid: string;
  name: string;
  labels: DataSourceLabelKeyValuePair[];
}

export declare class Step {
  uuid: string;
  name: string;
  keyValuePairs: StepKeyValuePair[];
  type: string;
  outputLabelsArray: KeyValuePair[][];
}

export declare class D3Node {
  id: string;
  text: string;
  x: number;
  y: number;
  width: number;
  height: number;
  type: string;
}

export declare class D3Edge {
  start: string;
  stop: string;
}
