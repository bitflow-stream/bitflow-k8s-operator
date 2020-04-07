// Flags
export const offlineMode = false;
export const useLocalDataSources = offlineMode || false;
export const useLocalSteps = offlineMode || false;
export const useLocalPods = offlineMode || false;
export const useLocalMatchingSources = offlineMode || false;

// Proxy config
export const dataSourcesLink = 'http://localhost:8080/datasources';
export const podsLink = 'http://localhost:8080/pods';
export const stepsLink = 'http://localhost:8080/steps';
export const matchingSourcesLink = 'http://localhost:8080/matchingSources';

// Graph config
export const svgNodeWidth: number = 200;
export const svgNodeHeight: number = 100;
export const svgHorizontalGap: number = 100;
export const svgVerticalGap: number = 50;
export const svgPodNodeMargin: number = 20;
export const svgNodeMargin: number = 30;
export const maxNumberOfSeparateGraphElements: number = 3;
export const maxNodeTextLength: number = 23;
