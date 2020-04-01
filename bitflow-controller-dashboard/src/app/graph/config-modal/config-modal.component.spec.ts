import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {ConfigModalComponent} from './config-modal.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {SharedService} from "../../../shared-service";
import {AppRoutingModule} from "../../app-routing.module";

describe('ConfigModalComponent', () => {
  let component: ConfigModalComponent;
  let fixture: ComponentFixture<ConfigModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        AppRoutingModule,
        FormsModule,
        ReactiveFormsModule
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

  it('should create config modal component', () => {
    expect(component).toBeTruthy();
  });
});
