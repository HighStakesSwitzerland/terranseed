import {HttpClientModule} from "@angular/common/http";
import {NgModule} from "@angular/core";
import {GoogleMapsModule} from "@angular/google-maps";
import {MatButtonModule} from "@angular/material/button";
import {MatCardModule} from "@angular/material/card";
import {MatDividerModule} from "@angular/material/divider";
import {MatExpansionModule} from "@angular/material/expansion";
import {MatGridListModule} from "@angular/material/grid-list";
import {MatIconModule} from "@angular/material/icon";
import {MatToolbarModule} from "@angular/material/toolbar";
import {BrowserModule} from "@angular/platform-browser";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {NgxChartsModule} from "@swimlane/ngx-charts";
import {PeerService} from "../lib/infra/peer-service";
import {AppComponent} from "./app.component";
import {GmapsComponent} from "./gmaps/gmaps.component";
import {InfoWindowComponent} from "./gmaps/info-window/info-window.component";
import {HeaderComponent} from "./header/header.component";
import {PieChartComponent} from "./pie-chart/pie-chart.component";

@NgModule({
  declarations: [
    AppComponent,
    GmapsComponent,
    HeaderComponent,
    PieChartComponent,
    InfoWindowComponent,

  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    GoogleMapsModule,
    BrowserAnimationsModule,
    MatCardModule,
    MatGridListModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    NgxChartsModule,
    MatDividerModule,
    MatExpansionModule,
  ],
  providers: [
    PeerService
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
