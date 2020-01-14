import {async, TestBed} from '@angular/core/testing';
import {AppComponent} from './app.component';

describe('AppComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        AppComponent
      ],
    }).compileComponents();
  }));

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  });

  it(`should have as title 'zerops-frontend'`, () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app.title).toEqual('zerops-frontend');
  });

  it('should return list of all DataSourceGraphElements in KubernetesGraph', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;

    let kubernetesGraph = {
      dataSourceGraphElements: [
        {
          uuid: "source-abc",
          stepGraphElements: [
            {
              uuid: "step-ghi",
              outputDataSourceGraphElements: [
                {
                  uuid: "source-def",
                  stepGraphElements: [
                    {
                      uuid: "step-jkl",
                      outputDataSourceGraphElements: [],
                      sourceDataSourceGraphElements: []
                    },
                    {
                      uuid: "step-mno",
                      outputDataSourceGraphElements: [],
                      sourceDataSourceGraphElements: []
                    }
                  ],
                  creatorStepGraphElement: null
                },
                {
                  uuid: "source-ghi",
                  stepGraphElements: [],
                  creatorStepGraphElement: null
                }
              ],
              sourceDataSourceGraphElements: []
            }
          ],
          creatorStepGraphElement: null
        }
      ],
      unusedStepGraphElements: [
        {
          uuid: "step-abc",
          outputDataSourceGraphElements: [],
          sourceDataSourceGraphElements: []
        },
        {
          uuid: "step-def",
          outputDataSourceGraphElements: [],
          sourceDataSourceGraphElements: []
        }
      ]
    };

    let dataSourceGraphElements = app.getAllDataSourceGraphElements(kubernetesGraph);

    expect(dataSourceGraphElements.length).toEqual(3);
    expect(dataSourceGraphElements.map(element => element.uuid).sort()).toEqual(["source-abc", "source-def", "source-ghi"].sort())
  });

  it('should return list of all StepGraphElements in KubernetesGraph', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;

    let kubernetesGraph = {
      dataSourceGraphElements: [
        {
          uuid: "source-abc",
          stepGraphElements: [
            {
              uuid: "step-ghi",
              outputDataSourceGraphElements: [
                {
                  uuid: "source-def",
                  stepGraphElements: [
                    {
                      uuid: "step-jkl",
                      outputDataSourceGraphElements: [],
                      sourceDataSourceGraphElements: []
                    },
                    {
                      uuid: "step-mno",
                      outputDataSourceGraphElements: [],
                      sourceDataSourceGraphElements: []
                    }
                  ],
                  creatorStepGraphElement: null
                },
                {
                  uuid: "source-ghi",
                  stepGraphElements: [],
                  creatorStepGraphElement: null
                }
              ],
              sourceDataSourceGraphElements: []
            }
          ],
          creatorStepGraphElement: null
        }
      ],
      unusedStepGraphElements: [
        {
          uuid: "step-abc",
          outputDataSourceGraphElements: [],
          sourceDataSourceGraphElements: []
        },
        {
          uuid: "step-def",
          outputDataSourceGraphElements: [],
          sourceDataSourceGraphElements: []
        }
      ]
    };

    let dataSourceGraphElements = app.getAllStepGraphElements(kubernetesGraph);

    expect(dataSourceGraphElements.length).toEqual(5);
    expect(dataSourceGraphElements.map(element => element.uuid).sort()).toEqual(["step-abc", "step-def", "step-ghi", "step-jkl", "step-mno"].sort())
  });
});
