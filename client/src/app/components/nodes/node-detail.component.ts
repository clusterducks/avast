import {Component, OnInit} from 'angular2/core';
import {RouteParams} from 'angular2/router';

import {ConsulService} from '../consul/providers/consul.service';
import {SwarmNode} from './interfaces/swarm-node';

@Component({
  selector: 'avast-node-detail',
  template: require('./node-detail.component.html'),
  styles: [
    require('./node-detail.component.css')
  ]
})

export class NodeDetailComponent implements OnInit {
  public node: SwarmNode;

  constructor(private _consulService: ConsulService,
              private _routeParams: RouteParams) {
  }

  ngOnInit() {
    if (!this.node) {
      let name = this._routeParams.get('name');
      this._consulService.getNode(name)
        .subscribe(res => this.node = res);
    }
  }

  goBack() {
    window.history.back();
  }
}
