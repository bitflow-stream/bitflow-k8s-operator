import {AfterViewInit, Component, ElementRef, Input, ViewChild} from '@angular/core';
import {ModalDismissReasons, NgbModal} from '@ng-bootstrap/ng-bootstrap';
import {GraphElement} from '../../../externalized/definitions/definitions';
import {
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
    private fb: FormBuilder
  ) {
  }

  @Input() currentGraphElementsWithStacksMap: Map<string, GraphElement> = new Map();

  currentGraphElement: GraphElement | undefined;
  selectedIdentifier: string | undefined;

  @ViewChild('content', {static: false}) theModal: ElementRef | undefined;

  closeResult: string | undefined;

  private static getDismissReason(reason: any): string {
    if (reason === ModalDismissReasons.ESC) {
      return 'by pressing ESC';
    } else if (reason === ModalDismissReasons.BACKDROP_CLICK) {
      return 'by clicking on a backdrop';
    } else {
      return `with: ${reason}`;
    }
  }

  selectedElement = () => getGraphElementByIdentifier(this.selectedIdentifier);

  ngAfterViewInit(): void {
    // TODO URL params don't open modals since commit 24c2f58b8aaa8b591b10101fcdc4c0938bee8279, probably has to do with async functions
    this.route.paramMap.subscribe(params => {
      let idParam = params.get('id');
      if (idParam == null) {
        return;
      }
      this.openModal(idParam);
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

  openModal(identifier: string) {
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

    this.selectedIdentifier = identifier;

    this.modalService.open(this.theModal, {ariaLabelledBy: 'modal-basic-title', size: 'lg'}).result.then((result) => {
      this.router.navigate(['']).then(() => {
      });
      this.closeResult = `Closed with: ${result}`;
    }, (reason) => {
      this.router.navigate(['']).then(() => {
      });
      this.closeResult = `Dismissed ${ConfigModalComponent.getDismissReason(reason)}`;
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
    }
    if (graphElement.type === 'data-source') {
      let dataSource = graphElement.dataSource;

      let specUrl: string = this.dataSourceFormData.value['specUrl'];
      if (specUrl != undefined) {
        dataSource.specUrl = specUrl;
      }

      console.log(getRawDataFromDataSource(dataSource));
    }
    if (graphElement.type === 'pod') {
      let pod = graphElement.pod;

      let raw: string = this.podFormData.value['raw'];
      if (raw != undefined) {
        pod.raw = raw;
      }

      console.log(getRawDataFromPod(pod));
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

}
