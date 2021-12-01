import {Component, OnDestroy, OnInit} from "@angular/core";
import {filter, groupBy, isNil, reduce, sortBy} from "lodash-es";
import {BehaviorSubject, Observable, Subject, switchMap, takeUntil, timer} from "rxjs";
import {Peer} from "../lib/domain/peer";
import {PeerService} from "../lib/infra/peer-service";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css"]
})
export class AppComponent implements OnInit, OnDestroy {
  totalPeers: number = 0;
  biggestProviderNames: string;
  peers$ = new BehaviorSubject<Peer[]>([]);

  private _destroy$ = new Subject();

  constructor(private readonly _peerService: PeerService) {
    timer(0, 5000)
      .pipe(
        switchMap(() => this._peerService.getAllPeers()),
        takeUntil(this._destroy$))
      .subscribe(peers => this.peers$.next(peers));

    this.peers$.subscribe((peers: Peer[]) => {
      try {
        this.totalPeers = peers?.length;
        const peerMap = groupBy(peers, p => p.isp);
        let peerArray: Array<any> = Object.keys(peerMap).map((key) => {
          return {
            name: key,
            value: peerMap[key]
          };
        });
        let sortedPeers = sortBy(peerArray, peer => peer.value.length).reverse();
        this.biggestProviderNames = reduce(
          filter(sortedPeers, peer => peer.value.length === sortedPeers[0].value.length),
          (names, peer) => {
            return (names.hasOwnProperty("name") ? names.name : names) + ", " + peer.name;
          }
        );
      } catch (err) {
        console.error(err);
      }
    });
  }


  ngOnInit(): void {

  }


  ngOnDestroy(): void {
    this._destroy$.complete();
  }

}
