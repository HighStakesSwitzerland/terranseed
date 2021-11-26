import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-gmaps',
  templateUrl: './gmaps.component.html',
  styleUrls: ['./gmaps.component.css']
})
export class GmapsComponent implements OnInit {

  //https://timdeschryver.dev/blog/google-maps-as-an-angular-component#mapinfowindow
  center: google.maps.LatLngLiteral = {
    lat: 46,
    lng: 6
  }
  options: google.maps.MapOptions = {
    zoomControl: true,
    scrollwheel: true,
    zoom: 12,
    clickableIcons: true,
    streetViewControl: false,
  }

  constructor() { }

  ngOnInit(): void {
  }

}
