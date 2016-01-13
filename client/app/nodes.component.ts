import {Component, OnInit} from 'angular2/core';
import {Router} from 'angular2/router';
import {SwarmNode} from './swarm-node';
import {ConsulService} from './consul.service';
import {NodeDetailComponent} from './node-detail.component';

@Component({
  selector: 'avast-nodes',
  templateUrl: 'app/nodes.component.html',
  styleUrls: ['app/nodes.component.css'],
  directives: [NodeDetailComponent]
})

/* @TODO: make this a shared directive */

export class NodesComponent implements OnInit {
  public nodes: SwarmNode[];
  public selectedNode: SwarmNode;

  constructor(private _consulService: ConsulService,
              private _router: Router) {
  }

  ngOnInit() {
    this.getNodes();
  }

  getNodes() {
    this._consulService.getNodes()
      .subscribe(res => this.nodes = res);
  }

  onSelect(node: SwarmNode) {
    this.selectedNode = node;
  }

  gotoDetail() {
    this._router.navigate(['NodeDetail', {
      name: this.selectedNode.node
    }]);
  }
}
