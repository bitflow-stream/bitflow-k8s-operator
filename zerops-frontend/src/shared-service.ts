import {Injectable} from "@angular/core";
import {Observable, Subject} from "rxjs";
import {GraphElement} from "./externalized/definitions/definitions";

@Injectable()
export class SharedService {
  // Observable string sources
  private emitChangeSource = new Subject<GraphElement>();
  // Observable string streams
  changeEmitted$: Observable<GraphElement> = this.emitChangeSource.asObservable();

  // Service message commands
  filterGraph(change: GraphElement) {
    this.emitChangeSource.next(change);
  }
}
