import {AfterViewInit, Component, ElementRef, Input, ViewChild} from '@angular/core';
import {ModalDismissReasons, NgbModal} from '@ng-bootstrap/ng-bootstrap';
import {GraphElement} from '../../../externalized/definitions/definitions';
import {getGraphElementByIdentifier} from '../../../externalized/functionalities/quality-of-life-functions';
import {ActivatedRoute, Router} from '@angular/router';
import {Location} from '@angular/common';
import {SharedService} from '../../../shared-service';

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
    private sharedService: SharedService
  ) {
  }

  @Input() currentGraphElementsWithStacksMap: Map<string, GraphElement> = new Map();

  currentGraphElement: GraphElement | undefined;
  selectedIdentifier: string | undefined;

  @ViewChild('content', {static: false}) theModal: ElementRef | undefined;

  closeResult: string | undefined;

  idParam = undefined;

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
    this.route.paramMap.subscribe(params => {
      this.idParam = params.get('id');
      this.openModal(this.idParam);
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

  openModal(identifier: string) { // TODO When following links, url is not updated
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

}
