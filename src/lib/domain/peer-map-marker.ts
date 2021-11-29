import {MapMarker} from "@angular/google-maps";
import {Peer} from "./peer";

export interface PeerMapMarker extends MapMarker {
  peers: Peer[];
}
