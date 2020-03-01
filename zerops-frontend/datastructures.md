# Datenstrukturen

## FrontendData

```json
{
    nodes: D3Node[];
    edges: D3Edge[];
}
```

### D3Node

```json
{
  id: string;
  text: string;
  x: number;
  y: number;
  width: number;
  height: number;
  type: 'step' | 'data-source' | 'pod' | 'data-source-stack' | 'pod-stack';
}
```

### D3Edge

```json
{
  start: string;
  stop: string;
}
```

## Backend

### GraphVisualization

```json
{
    graphColumns: GraphVisualizationColumn[];
}
```

### GraphVisualizationColumn

```json
{
    graphElements: GraphElement[];
}
```

Pods in Steps werden wie folgt dargestellt: `[Step1, Pod1_1, Pod1_2, Step2, Pod2_1]`

### GraphElement

```json
{
    type: 'step' | 'data-source' | 'pod' | 'data-source-stack' | 'pod-stack';
    step?: Step;
    dataSource?: DataSource
    pod?: Pod;
    dataSourceStack?: DataSourceStack
    podStack?: PodStack;
}
```

### DataSourceStack

```json
{
    stackId: string;
    hasSourceGraphElement: boolean;
    sourceGraphElement?: GraphElement;
    outputName: string;
    dataSources: DataSource[];
}
```

### PodStack

```json
{
    stackId: string;
    hasCreatorStep: boolean;
    creatorStep?: Step;
    pods: Pod[];
}
```

### Step

```json
{
    name: string,;
    ingests: Ingest[];
    outputs: Output[];
    validationError: string;
    template: string;
    podType: 'pod' | 'pod-stack';
    pods?: Pod[];
    podStack?: PodStack;
}
```

### Ingest
```json
{
  key: string;
  value?: string;
  check?: string;
}
```

### Output
```json
{
    name: string;
    url: string;
    "labels": Label[];
}
```

### Label
```json
{
  key: string;
  value: string;
}
```

### DataSource

```json
{
    name: string;
    labels: Label[];
    specUrl: string;
    validationError: string;
    hasCreatorPod: boolean;
    creatorPod?: Pod;
    hasOutputName: boolean;
    outputName?: string;
    createdPods: Pod[];
}
```

### Pod

```json
{
    name: string;
    phase: string;
    hasCreatorStep: boolean;
    creatorStep?: Step;
    creatorDataSources: DataSource[];
	createdDataSources: DataSource[];
}
```

### StepDataSourceMatches

```json
export interface StepDataSourceMatches {
  [key: string]: string[]
}
```
