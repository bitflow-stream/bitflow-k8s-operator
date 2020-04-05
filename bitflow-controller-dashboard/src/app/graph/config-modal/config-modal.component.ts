import {AfterViewInit, Component, ElementRef, Input, NgZone, ViewChild} from '@angular/core';
import {NgbModal} from '@ng-bootstrap/ng-bootstrap';
import {GraphElement} from '../../../externalized/definitions/definitions';
import {
  getAllCurrentGraphElementsWithStacks,
  getGraphElementByIdentifier,
  getRawDataFromDataSource,
  getRawDataFromPod,
  getRawDataFromStep
} from '../../../externalized/functionalities/quality-of-life-functions';
import {ActivatedRoute, Router} from '@angular/router';
import {Location} from '@angular/common';
import {SharedService} from '../../../shared-service';
import {FormBuilder} from '@angular/forms';

@Component({
  selector: 'app-config-modal',
  templateUrl: './config-modal.component.html',
  styleUrls: ['./config-modal.component.css']
})
export class ConfigModalComponent implements AfterViewInit {

  constructor(
    private modalService: NgbModal,
    private readonly route: ActivatedRoute,
    private readonly router: Router,
    private location: Location,
    private sharedService: SharedService,
    private fb: FormBuilder,
    private ngZone: NgZone
  ) {
  }

  @Input() currentGraphElementsWithStacksMap: Map<string, GraphElement> = new Map();

  currentGraphElement: GraphElement | undefined;
  selectedIdentifier: string | undefined;

  @ViewChild('content', {static: false}) theModal: ElementRef | undefined;

  selectedElement = () => getGraphElementByIdentifier(this.selectedIdentifier);

  async ngAfterViewInit() {
    this.route.paramMap.subscribe(params => {
      this.ngZone.run(async () => {
        let idParam = params.get('id');
        if (idParam == null) {
          return;
        }

        await new Promise(resolve => setTimeout(() => {
          this.openModal(idParam);
          resolve();
        }, getAllCurrentGraphElementsWithStacks().length === 0 ? 500 : 0));
      });
    });
  }

  updateUrlBySelectElement(element: any) {
    this.updateUrl('/id/' + element.value);
  }

  updateUrl(url: string) {
    this.location.replaceState(url);
  }

  goto(identifier: string): void {
    this.router.navigate(['id', identifier]).then(() => {
    });
  }

  async openModal(identifier: string) {
    this.selectedIdentifier = identifier;

    this.modalService.dismissAll();
    this.currentGraphElement = getGraphElementByIdentifier(identifier);

    if (this.currentGraphElement === undefined || this.currentGraphElement === null) {
      return;
    }

    if (this.currentGraphElement.type === 'data-source-stack' && this.currentGraphElement.dataSourceStack.dataSources.length !== 0) {
      this.goto(this.currentGraphElement.dataSourceStack.dataSources[0].name);
      return;
    }
    if (this.currentGraphElement.type === 'pod-stack' && this.currentGraphElement.podStack.pods.length !== 0) {
      this.goto(this.currentGraphElement.podStack.pods[0].name);
      return;
    }

    await this.modalService.open(this.theModal, {
      ariaLabelledBy: 'modal-basic-title',
      size: 'lg'
    });
  }

  filterGraph(graphElement: GraphElement) {
    this.sharedService.filterGraph(graphElement);
  }

  save(graphElement: GraphElement) {
    if (graphElement == undefined) {
      return;
    }

    if (graphElement.type === 'step') {
      let step = graphElement.step;

      let template: string = this.stepFormData.value['template'];
      if (template != undefined) {
        step.template = template;
      }

      console.log(getRawDataFromStep(step));
      // TODO save in kubernetes
    }
    if (graphElement.type === 'data-source') {
      let dataSource = graphElement.dataSource;

      let specUrl: string = this.dataSourceFormData.value['specUrl'];
      if (specUrl != undefined) {
        dataSource.specUrl = specUrl;
      }

      console.log(getRawDataFromDataSource(dataSource));
      // TODO save in kubernetes
    }
    if (graphElement.type === 'pod') {
      let pod = graphElement.pod;

      let raw: string = this.podFormData.value['raw'];
      if (raw != undefined) {
        pod.raw = raw;
      }

      console.log(getRawDataFromPod(pod));
      // TODO save in kubernetes
    }
  }

  podFormData = this.fb.group({
    raw: []
  });

  dataSourceFormData = this.fb.group({
    specUrl: []
    // TODO labels
    // TODO removing / adding labels
  });

  stepFormData = this.fb.group({
    template: []
    // TODO ingests
    // TODO removing / adding ingests
    // TODO outputs
    // TODO removing / adding outputs
  });

  handleSubmit() {
    this.modalService.dismissAll();
    this.save(this.selectedElement())
  }

  removeLabelFromDataSource(graphElement: GraphElement, index: number) {
    graphElement.dataSource.labels.splice(index, 1);
  }

  addLabelToDataSource(graphElement: GraphElement) {
    graphElement.dataSource.labels.push({key: '', value: ''})
  }

}
