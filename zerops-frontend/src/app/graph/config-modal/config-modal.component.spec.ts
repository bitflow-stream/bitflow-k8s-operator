import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {ConfigModalComponent} from './config-modal.component';
import {FormsModule} from "@angular/forms";
import {SharedService} from "../../../shared-service";
import {AppRoutingModule} from "../../app-routing.module";

describe('ConfigModalComponent', () => {
  let component: ConfigModalComponent;
  let fixture: ComponentFixture<ConfigModalComponent>;

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
        ConfigModalComponent
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfigModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
