import {Component, OnInit} from 'angular2/core';
import {Router} from 'angular2/router';

import {ConsulService} from '../consul/providers/consul.service';
import {SwarmNode} from '../nodes/interfaces/swarm-node';

@Component({
  selector: 'avast-dashboard',
  template: require('./dashboard.component.html'),
  styles: [
    require('./dashboard.component.css')
  ],
})

export class DashboardComponent implements OnInit {
  public datacenters: string[];
  public nodes: SwarmNode[];

  constructor (private _consulService: ConsulService,
               private _router: Router) {
  }

  ngOnInit() {
    this.getDatacenters();
  }

  getDatacenters() {
    this._consulService.getDatacenters()
      .subscribe(
        res => this.datacenters = res,
        err => this.logError(err)
      );
  }

  setDatacenter(dc: string) {
    this._consulService.getNodes(dc)
      .subscribe(res => this.nodes = res);
  }

  gotoNode(name: string) {
    this._router.navigate(['NodeDetail', { name: name }]);
  }

  logError(err) {
    console.log(err);
  }
}
