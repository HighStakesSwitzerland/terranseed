import {HttpClient} from "@angular/common/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs";
import {Peer} from "../domain/peer";

@Injectable({
  providedIn: "root"
})
export class PeerService {

  constructor(private readonly httpClient: HttpClient) {
  }

  public getAllPeers(): Observable<Peer[]> {
    return this.httpClient.get<Peer[]>("/api/peers");
  }

}
