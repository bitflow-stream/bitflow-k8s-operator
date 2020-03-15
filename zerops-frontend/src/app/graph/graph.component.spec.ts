import {async, TestBed} from '@angular/core/testing';
import {GraphComponent} from './graph.component';
import {ConfigModalComponent} from "./config-modal/config-modal.component";
import {AppRoutingModule} from "../app-routing.module";
import {FormsModule} from "@angular/forms";
import {SharedService} from "../../shared-service";

describe('GraphComponent', () => {
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

  it('should create the app', () => {
    const fixture = TestBed.createComponent(GraphComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  });

});
