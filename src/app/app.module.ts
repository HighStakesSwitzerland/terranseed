import { NgModule } from '@angular/core';
import {GoogleMapsModule} from "@angular/google-maps";
import { BrowserModule } from '@angular/platform-browser';

import { AppComponent } from './app.component';
import { GmapsComponent } from './gmaps/gmaps.component';

@NgModule({
  declarations: [
    AppComponent,
    GmapsComponent
  ],
  imports: [
    BrowserModule,
    GoogleMapsModule
  ],
  providers: [
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
