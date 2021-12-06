import {Component, Input, OnInit, QueryList, ViewChild, ViewChildren} from "@angular/core";
import {FormControl} from "@angular/forms";
import {GoogleMap, MapInfoWindow, MapMarker} from "@angular/google-maps";
import {MatAutocompleteSelectedEvent} from "@angular/material/autocomplete";
import {filter, find, isEmpty, isNil} from "lodash-es";
import {map, Observable} from "rxjs";
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

  @ViewChildren(MapMarker)
  mapMarkers: QueryList<MapMarker>;

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
  autoCompleteCtrl: FormControl = new FormControl();
  filteredMarkers: Observable<PeerMapMarker[]>;

  ngOnInit(): void {
    this.peers$.subscribe(peers => {
      if (!isNil(peers)) {
        const toMark = filter(peers, (p) => !find(this.markers, marker => find(marker.peers, mp => mp.nodeId === p.nodeId)));
        this.markPeers(toMark);
      }
    });
    this.filteredMarkers = this.autoCompleteCtrl.valueChanges.pipe(
      map((value: string) => this._filter(value)),
    );
  }

  markPeers(toMark: Peer[]) {
    toMark?.forEach(peer => {
      // if there is already a marker at this position
      let existingMarker = this.markers.find(m =>
        // round position to 2 digits to mitigate ip geolocalization errors
        (m.position.lat as number).toFixed(2) == peer.lat.toFixed(2)
        && (m.position.lng as number).toFixed(2) == peer.lon.toFixed(2)
        && m.peers[0].isp === peer.isp
      );
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
    if (mapMarker.hasOwnProperty("position")) {
      this.map.panTo(mapMarker.position);
    } else if (mapMarker.hasOwnProperty("marker")) {
      // @ts-ignore
      this.map.panTo(mapMarker.marker.position);
    }
    infoWindow.open(mapMarker);
  }

  closeInfoWindows() {
    this.infoWindows?.forEach(win => win.close());
  }

  private _filter(value: string): PeerMapMarker[] {
    const query = new RegExp(value, "i");
    return this.markers.filter(option => {
        let m = filter(option.peers, p => {
          const matched = p.moniker.match(query) || p.nodeId.match(query);
          if (matched) {
            option.filteredMoniker = p.moniker;
          }
          return matched;
        });
        return !isEmpty(m);
      }
    );
  }


  private locateMarkerOnMap(value: PeerMapMarker) {
    const marker = find(this.mapMarkers.toArray(), item => {
      // thanks gmaps not to have a relation between markers and info windows
      return item?.marker?.getPosition()!.lat() === value.position.lat && item.marker.getPosition()!.lng() === value.position.lng;
    });
    google.maps.event.trigger(marker!.marker!, "click");
  }

  autocompleteOptionSelected($event: MatAutocompleteSelectedEvent) {
    this.locateMarkerOnMap($event.option.value);
  }

  getMarkerText(peer: PeerMapMarker) {
    return peer?.filteredMoniker;
  }
}
