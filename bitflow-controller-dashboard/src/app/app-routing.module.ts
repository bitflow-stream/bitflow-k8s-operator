import {NgModule} from '@angular/core';
import {RouterModule, Routes} from "@angular/router";
import {ConfigModalComponent} from "./graph/config-modal/config-modal.component";

const routes: Routes = [
  {path: "id/:id", component: ConfigModalComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {
}
