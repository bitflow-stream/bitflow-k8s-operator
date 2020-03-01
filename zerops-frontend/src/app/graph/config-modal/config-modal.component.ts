import {AfterViewInit, Component, ElementRef, EventEmitter, Input, Output, ViewChild} from '@angular/core';
import {ModalDismissReasons, NgbModal} from "@ng-bootstrap/ng-bootstrap";
import {GraphElement} from "../../../externalized/definitions/definitions";
import {getGraphElementByIdentifier} from "../../../externalized/functionalities/quality-of-life-functions";
import {ActivatedRoute, Router} from "@angular/router";

@Component({
  selector: 'app-pod-modal',
  templateUrl: './config-modal.component.html',
  styleUrls: ['./config-modal.component.css']
})
export class ConfigModalComponent implements AfterViewInit {
  @Input() currentGraphElementsWithStacksMap: Map<string, GraphElement> = new Map();
  @Output() updateGraphEvent = new EventEmitter<GraphElement>();

  currentGraphElement: GraphElement | undefined;
  selectedIdentifier: string | undefined;
  selectedElement = () => getGraphElementByIdentifier(this.selectedIdentifier);

  @ViewChild('content', {static: false}) theModal: ElementRef | undefined;

  closeResult: string | undefined;

  idParam = undefined;

  constructor(
    private modalService: NgbModal,
    private readonly route: ActivatedRoute,
    private readonly router: Router,
  ) {
  }

  ngAfterViewInit(): void {
    this.route.paramMap.subscribe(params => {
      this.idParam = params.get("id");
      this.openModal(this.idParam);
    });
  }

  goto(identifier: string): void {
    this.router.navigate(["id", identifier]);
  }

  openModal(identifier: string) {
    this.currentGraphElement = getGraphElementByIdentifier(identifier);
    if (this.currentGraphElement == undefined) {
      return;
    }
    this.selectedIdentifier = undefined;

    this.modalService.open(this.theModal, {ariaLabelledBy: 'modal-basic-title', size: 'lg'}).result.then((result) => {
      this.closeResult = `Closed with: ${result}`;
    }, (reason) => {
      this.closeResult = `Dismissed ${ConfigModalComponent.getDismissReason(reason)}`;
    });
  }

  private static getDismissReason(reason: any): string {
    if (reason === ModalDismissReasons.ESC) {
      return 'by pressing ESC';
    } else if (reason === ModalDismissReasons.BACKDROP_CLICK) {
      return 'by clicking on a backdrop';
    } else {
      return `with: ${reason}`;
    }
  }

  filterGraph(graphElement: GraphElement) {
    this.updateGraphEvent.next(graphElement)
  }

}
