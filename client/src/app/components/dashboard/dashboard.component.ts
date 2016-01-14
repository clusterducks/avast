import {Component, OnInit} from 'angular2/core';
import {Router} from 'angular2/router';
import {AppStore} from 'angular2-redux';

import {ConsulActions} from '../../actions/consul.actions';
import {ConsulService} from '../consul/providers/consul.service';
import {SwarmNode} from '../nodes/interfaces/swarm-node';

@Component({
  selector: 'avast-dashboard',
  template: require('./dashboard.component.html'),
  styles: [
    require('./dashboard.component.css')
  ]
})

export class DashboardComponent implements OnInit {
  public datacenters: string[] = [];
  public nodes: SwarmNode[];
  private isFetchingDatacenters = false;
  private isFetchingNodes = false;

  constructor (private _router: Router,
               private _appStore: AppStore,
               private _consulActions: ConsulActions,
               private _consulService: ConsulService) {
  }

  ngOnInit() {
    this._appStore.subscribe((state) => {
      this.datacenters = state.datacenters;
      this.isFetchingDatacenters = state.isFetchingDatacenters;
    });

    this._appStore.dispatch(this._consulActions.fetchDatacenters());
  }

  setDatacenter(dc: string) {
    this._consulService.getNodes(dc)
      .subscribe(res => this.nodes = res);
  }

  gotoNode(name: string) {
    this._router.navigate(['NodeDetail', { name: name }]);
  }
}
