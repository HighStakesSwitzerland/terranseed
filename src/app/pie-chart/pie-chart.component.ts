import {Component, Input, OnDestroy, OnInit} from "@angular/core";
import {remove, sortBy} from "lodash-es";
import {Observable, Subject, takeUntil} from "rxjs";
import {Peer} from "../../lib/domain/peer";

@Component({
  selector: "app-pie-chart",
  templateUrl: "./pie-chart.component.html",
  styleUrls: ["./pie-chart.component.css"]
})
export class PieChartComponent implements OnInit, OnDestroy {
  @Input()
  peers$: Observable<Peer[]>;

  pieData: PieData[] = [];

  private _destroy$ = new Subject();
  private _providers = ["Amazon", "Google", "Digital", "Hetzner", "Microsoft"];

  ngOnInit(): void {
    this.peers$.pipe(
      takeUntil(this._destroy$)
    ).subscribe(peers => {
      this.analyzePeers(peers.slice());
    });
  }

  private analyzePeers(peers: Peer[]) {
    const newPieData: PieData[] = [];
    this._providers.forEach(provider => {
      let groupedPeers = remove(peers, (peer => peer.isp.startsWith(provider)));
      if (groupedPeers.length > 0) {
        newPieData.push({
          name: groupedPeers[0].isp,
          value: groupedPeers.length
        });
      }
    });
    if (peers.length > 0) {
      newPieData.push({
        name: "Others",
        value: peers.length
      });
    }
    this.pieData = sortBy(newPieData, "value").reverse();
  }

  ngOnDestroy(): void {
    this._destroy$.complete();
  }

}

interface PieData {
  name: string;
  value: number;
}
