import {Component, OnDestroy} from "@angular/core";
import {isNil} from "lodash-es";
import {BehaviorSubject, filter, Subject, switchMap, takeUntil, timer} from "rxjs";
import {Peer} from "../lib/domain/peer";
import {PeerService} from "../lib/infra/peer-service";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css"]
})
export class AppComponent implements OnDestroy {
  peers$ = new BehaviorSubject<Peer[]>([]);

  private _destroy$ = new Subject();

  constructor(private readonly _peerService: PeerService) {
    timer(0, 5000)
      .pipe(
        switchMap(() => this._peerService.getAllPeers()),
        filter(peers => !isNil(peers)),
        takeUntil(this._destroy$))
      .subscribe(peers => this.peers$.next(peers));
  }

  ngOnDestroy(): void {
    this._destroy$.complete();
  }

}
