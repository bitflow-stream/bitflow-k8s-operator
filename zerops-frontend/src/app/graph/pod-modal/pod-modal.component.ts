import {Component, ElementRef, ViewChild} from '@angular/core';
import {ModalDismissReasons, NgbModal} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'app-pod-modal',
  templateUrl: './pod-modal.component.html',
  styleUrls: ['./pod-modal.component.css']
})
export class PodModalComponent {
  @ViewChild('content', {static: false}) theModal: ElementRef | undefined;

  closeResult: string = '';

  constructor(private modalService: NgbModal) {
  }

  openModal(uuid: string) {
    console.log(uuid);
    this.modalService.open(this.theModal, {ariaLabelledBy: 'modal-basic-title'}).result.then((result) => {
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

}
