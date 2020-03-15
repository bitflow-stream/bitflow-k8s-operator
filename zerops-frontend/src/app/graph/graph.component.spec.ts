import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {getGraphVisualization, GraphComponent} from './graph.component';
import {ConfigModalComponent} from "./config-modal/config-modal.component";
import {AppRoutingModule} from "../app-routing.module";
import {FormsModule} from "@angular/forms";
import {SharedService} from "../../shared-service";
import {DebugElement} from "@angular/core";
import {
  setAllCurrentGraphElementsWithStacks,
  setCurrentGraphElements
} from "../../externalized/functionalities/quality-of-life-functions";
import {Pod, Step} from "../../externalized/definitions/definitions";

describe('GraphComponent', () => {
  let fixture: ComponentFixture<any>;
  let component: GraphComponent;
  let de: DebugElement;


  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        AppRoutingModule,
        FormsModule
      ],
      providers: [
        SharedService
      ],
      declarations: [
        GraphComponent,
        ConfigModalComponent
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GraphComponent);
    component = fixture.componentInstance;
    de = fixture.debugElement;
  });

  it('should create graph component', () => {
    expect(component).toBeTruthy();
  });

  it('should create graph visualization correctly', () => {

    let step1: Step = {
      name: 'stepName-1',
      ingests: [],
      outputs: [],
      validationError: '',
      template: 'template-1',
      podType: 'pod',
      pods: [],
      raw: 'stepRaw-1'
    };

    let step2: Step = {
      name: 'stepName-2',
      ingests: [],
      outputs: [],
      validationError: '',
      template: 'template-2',
      podType: 'pod',
      pods: [],
      raw: 'stepRaw-2'
    };

    let step3: Step = {
      name: 'stepName-3',
      ingests: [],
      outputs: [],
      validationError: '',
      template: 'template-3',
      podType: 'pod',
      pods: [],
      raw: 'stepRaw-3'
    };

    let pod1: Pod = {
      name: 'podName-1',
      phase: 'podPhase-1',
      hasCreatorStep: true,
      creatorStep: step1,
      creatorDataSources: [],
      createdDataSources: [],
      raw: 'podRaw-1'
    };

    let pod2: Pod = {
      name: 'podName-2',
      phase: 'podPhase-2',
      hasCreatorStep: true,
      creatorStep: step1,
      creatorDataSources: [],
      createdDataSources: [],
      raw: 'podRaw-2'
    };

    let pod3: Pod = {
      name: 'podName-3',
      phase: 'podPhase-3',
      hasCreatorStep: true,
      creatorStep: step1,
      creatorDataSources: [],
      createdDataSources: [],
      raw: 'podRaw-3'
    };

    let pod4: Pod = {
      name: 'podName-4',
      phase: 'podPhase-4',
      hasCreatorStep: true,
      creatorStep: step2,
      creatorDataSources: [],
      createdDataSources: [],
      raw: 'podRaw-4'
    };

    let pod5: Pod = {
      name: 'podName-5',
      phase: 'podPhase-5',
      hasCreatorStep: true,
      creatorStep: step3,
      creatorDataSources: [],
      createdDataSources: [],
      raw: 'podRaw-5'
    };

    step1.pods.push(pod1);
    step1.pods.push(pod1);
    step1.pods.push(pod1);
    step2.pods.push(pod4);
    step3.pods.push(pod5);

    setCurrentGraphElements([], [step1, step2, step3], [pod1, pod2, pod3, pod4, pod5]);
    setAllCurrentGraphElementsWithStacks();
    let graphVisualization = getGraphVisualization();

    expect(graphVisualization).toBeTruthy();
    expect(graphVisualization.graphColumns[0]).toBeTruthy();
    expect(graphVisualization.graphColumns[0].graphElements.length).toEqual(8);
    expect(graphVisualization.graphColumns[0].graphElements[0].type).toEqual('step');
    expect(graphVisualization.graphColumns[0].graphElements[1].type).toEqual('pod');
    expect(graphVisualization.graphColumns[0].graphElements[2].type).toEqual('pod');
    expect(graphVisualization.graphColumns[0].graphElements[3].type).toEqual('pod');
    expect(graphVisualization.graphColumns[0].graphElements[4].type).toEqual('step');
    expect(graphVisualization.graphColumns[0].graphElements[5].type).toEqual('pod');
    expect(graphVisualization.graphColumns[0].graphElements[6].type).toEqual('step');
    expect(graphVisualization.graphColumns[0].graphElements[7].type).toEqual('pod');
  });

});
