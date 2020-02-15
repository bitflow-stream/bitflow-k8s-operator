import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PodModalComponent } from './pod-modal.component';

describe('PodModalComponent', () => {
  let component: PodModalComponent;
  let fixture: ComponentFixture<PodModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PodModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PodModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
