import {Component, OnInit} from 'angular2/core';
import {Router} from 'angular2/router';
import {SwarmNode} from './swarm-node';
import {ConsulService} from './consul.service';

@Component({
  selector: 'avast-dashboard',
  templateUrl: 'app/dashboard.component.html',
  styleUrls: ['app/dashboard.component.css'],
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
