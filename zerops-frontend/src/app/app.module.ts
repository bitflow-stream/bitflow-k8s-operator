import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { GraphComponent } from './graph/graph.component';
import { ConfigModalComponent } from './graph/config-modal/config-modal.component';
import {FormsModule} from "@angular/forms";
import {AppRoutingModule} from "./app-routing.module";
@NgModule({
  declarations: [
    AppComponent,
    GraphComponent,
    AppComponent,
    ConfigModalComponent
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    NgbModule,
    FormsModule
  ],
  providers: [],
  bootstrap: [
    AppComponent
  ]
})
export class AppModule { }
