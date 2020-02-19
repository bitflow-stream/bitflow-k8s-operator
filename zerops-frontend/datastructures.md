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
    maxFilledRow: number;
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
    dataSources: DataSource[];
}
```

### PodStack

```json
{
    stackId: string;
    pods: Pod[];
}
```

### Step

```json
{
    name: string;
    podType: 'pod' | 'pod-stack'
    pods?: Pod[]
	podStack?: PodStack
}
```

### DataSource

```json
{
    name: string;
    hasCreatorPod: boolean;
    creatorPod?: Pod;
    hasOutputName: boolean;
    outputName?: string;
}
```

### Pod

```json
{
    name: string;
    hasCreatorStep: boolean;
    creatorStep?: Step;
    hasOutputName: boolean;
    outputName?: string;
    creatorDataSources: DataSource[];
}
```
