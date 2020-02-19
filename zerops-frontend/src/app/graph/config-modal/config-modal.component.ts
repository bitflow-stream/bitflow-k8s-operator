import {Component, ElementRef, EventEmitter, Input, Output, ViewChild} from '@angular/core';
import {ModalDismissReasons, NgbModal} from "@ng-bootstrap/ng-bootstrap";
import {GraphElement} from "../../../externalized/definitions/definitions";
import {getGraphElementByIdentifier} from "../../../externalized/functionalities/quality-of-life-functions";

// function getKubernetesNode(uuid: string | undefined): KubernetesNode | undefined {
//   if (uuid == undefined) {
//     return undefined;
//   }
//
//   let kubernetesNode: KubernetesNode = {};
//
//   kubernetesGraph.forEach(column => {
//     column.forEach(row => {
//       if (row.dataSources != undefined && row.dataSources.length > 0) {
//         if (row.dataSources[0].dataSourceStackId === uuid) {
//           kubernetesNode.dataSources = [];
//           row.dataSources.forEach(dataSource => {
//             kubernetesNode.dataSources?.push(dataSource);
//           });
//         }
//       } else if (row.step != undefined && row.step.podNames.length > 0) {
//         if (row.step.podStackId === uuid) {
//           kubernetesNode.step = row.step;
//         }
//       }
//     });
//   });
//
//   return kubernetesNode;
// }

// function getElementNames(kubernetesNode: KubernetesNode | undefined) {
//   if (kubernetesNode == undefined) {
//     return [];
//   }
//
//   if (kubernetesNode.step != undefined) {
//     return kubernetesNode.step.podNames;
//   }
//   if (kubernetesNode.dataSources != undefined) {
//     return kubernetesNode.dataSources.map(dataSource => dataSource.name);
//   }
//
//   return [];
// }

@Component({
  selector: 'app-pod-modal',
  templateUrl: './config-modal.component.html',
  styleUrls: ['./config-modal.component.css']
})
export class ConfigModalComponent {
  @Input() currentGraphElementsWithStacksMap: Map<string, GraphElement> = new Map();
  @Output() updateGraphEvent = new EventEmitter<GraphElement>();

  currentGraphElement: GraphElement | undefined;
  selectedIdentifier: string | undefined;
  selectedElement = () => getGraphElementByIdentifier(this.selectedIdentifier);

  @ViewChild('content', {static: false}) theModal: ElementRef | undefined;

  closeResult: string | undefined;

  constructor(private modalService: NgbModal) {
  }

  openModal(identifier: string) {
    this.currentGraphElement = getGraphElementByIdentifier(identifier);
    this.selectedIdentifier = undefined;

    this.modalService.open(this.theModal, {ariaLabelledBy: 'modal-basic-title', size: 'lg'}).result.then((result) => {
      this.closeResult = `Closed with: ${result}`;
    }, (reason) => {
      this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
    });
  }

  private getDismissReason(reason: any): string {
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
