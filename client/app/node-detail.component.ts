import {Component, OnInit} from 'angular2/core';
import {RouteParams} from 'angular2/router';
import {SwarmNode} from './swarm-node';
import {ConsulService} from './consul.service';

@Component({
  selector: 'avast-node-detail',
  templateUrl: 'app/node-detail.component.html',
  styleUrls: ['app/node-detail.component.css'],
  //inputs: ['node'],
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
