import {Component, Input, OnInit, QueryList, ViewChild, ViewChildren} from "@angular/core";
import {GoogleMap, MapInfoWindow, MapMarker} from "@angular/google-maps";
import {filter, find, isNil} from "lodash-es";
import {Observable} from "rxjs";
import {Peer} from "../../lib/domain/peer";
import {PeerMapMarker} from "../../lib/domain/peer-map-marker";
import Icon = google.maps.Icon;
import InfoWindowOptions = google.maps.InfoWindowOptions;
import MarkerLabel = google.maps.MarkerLabel;
import Size = google.maps.Size;


@Component({
  selector: "app-gmaps",
  templateUrl: "./gmaps.component.html",
  styleUrls: ["./gmaps.component.css"]
})
export class GmapsComponent implements OnInit {

  @Input()
  peers$: Observable<Peer[]>;

  @ViewChild(GoogleMap, {static: false})
  map: GoogleMap;

  @ViewChildren(MapInfoWindow)
  infoWindows: QueryList<MapInfoWindow>;

  center: google.maps.LatLngLiteral = {
    lat: 46,
    lng: 6
  };
  options: google.maps.MapOptions = {
    zoomControl: true,
    scrollwheel: true,
    zoom: 2,
    clickableIcons: true,
    streetViewControl: false,
  };
  mapInfoWindowOptions: InfoWindowOptions = {};
  markers: PeerMapMarker[] = [];

  ngOnInit(): void {
    this.peers$.subscribe(peers => {
      if (!isNil(peers)) {
        peers.forEach(p => {
          // round coords to 1 digit
          p.lat = parseFloat(p.lat.toFixed(1));
          p.lon = parseFloat(p.lon.toFixed(1));
        });
        const toMark = filter(peers, (p) => !find(this.markers, marker => find(marker.peers, mp => mp.nodeId === p.nodeId)));
        this.markPeers(toMark);
      }
    });
  }

  markPeers(toMark: Peer[]) {
    toMark?.forEach(peer => {
      // if there is already a marker at this position
      let existingMarker = this.markers.find(m => m.position.lat == peer.lat && m.position.lng == peer.lon);
      if (!isNil(existingMarker)) {
        this.updateMarker(existingMarker, peer);
      } else {
        this.addNewMarker(peer);
      }
    });
  }

  private updateMarker(existingMarker: PeerMapMarker, peer: Peer) {
    const markerLabel = existingMarker.label as MarkerLabel;
    const exitingIconUrl = existingMarker.icon as Icon;
    let numberOfHosts = parseInt(markerLabel.text, 10);

    if (isNaN(numberOfHosts)) {
      numberOfHosts = 1;
    } else {
      numberOfHosts++;
    }
    markerLabel.text = String(numberOfHosts);
    if (numberOfHosts >= 10) {
      exitingIconUrl.url = "http://maps.google.com/mapfiles/ms/micons/red.png";
    } else if (numberOfHosts >= 7) {
      exitingIconUrl.url = "http://maps.google.com/mapfiles/ms/micons/pink.png";
    } else if (numberOfHosts >= 5) {
      exitingIconUrl.url = "http://maps.google.com/mapfiles/ms/micons/orange.png";
    } else if (numberOfHosts >= 2) {
      exitingIconUrl.url = "http://maps.google.com/mapfiles/ms/micons/lightblue.png";
    } else {
      exitingIconUrl.url = "http://maps.google.com/mapfiles/ms/micons/green.png";
    }

    existingMarker.peers.push(peer);
  }

  addNewMarker(peer: Peer) {
    const marker = {
      position: {
        lat: peer.lat,
        lng: peer.lon,
      },
      options: {
        animation: google.maps.Animation.DROP,
        clickable: true,
      },
      label: {
        className: "map-maker",
        text: "1"
      },
      icon: {
        url: "http://maps.google.com/mapfiles/ms/micons/green.png",
        size: new Size(40, 40, "px", "px"),
        origin: new google.maps.Point(-5, -10)
      }
    } as MapMarker;

    const peerMapMarker = {
      peers: [peer],
      ...marker
    } as PeerMapMarker;

    this.markers.push(peerMapMarker);
  }

  openInfoWindow(mapMarker: MapMarker, infoWindow: MapInfoWindow) {
    this.closeInfoWindows();
    infoWindow.open(mapMarker);
  }

  closeInfoWindows() {
    this.infoWindows?.forEach(win => win.close());
  }

}
