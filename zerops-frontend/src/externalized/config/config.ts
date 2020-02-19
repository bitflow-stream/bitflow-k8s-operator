const svgNodeWidth: number = 200;
const svgNodeHeight: number = 100;
const svgHorizontalGap: number = 100;
const svgVerticalGap: number = 50;
const svgPodNodeMargin: number = 20;
const svgNodeMargin: number = 30;
const maxNumberOfSeparateGraphElements: number = 1; // Currently only works with "1", since otherwise DataSource grouping to stacks doesn't work properly. (With "1", sourceGraphElement is always the pod stack)
export {svgNodeWidth, svgNodeHeight, svgHorizontalGap, svgVerticalGap, svgPodNodeMargin, svgNodeMargin, maxNumberOfSeparateGraphElements};
