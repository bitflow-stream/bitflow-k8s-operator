import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppComponent} from './app.component';
import {NgbModule} from '@ng-bootstrap/ng-bootstrap';
import {GraphComponent} from './graph/graph.component';
import {ConfigModalComponent} from './graph/config-modal/config-modal.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {AppRoutingModule} from "./app-routing.module";
import {SharedService} from "../shared-service";

@NgModule({
  declarations: [
    AppComponent,
    GraphComponent,
    ConfigModalComponent
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    NgbModule,
    FormsModule,
    ReactiveFormsModule
  ],
  providers: [
    SharedService
  ],
  bootstrap: [
    AppComponent
  ]
})
export class AppModule { }
