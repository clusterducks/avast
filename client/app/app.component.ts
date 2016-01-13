import {Component} from 'angular2/core';
import {RouteConfig, ROUTER_DIRECTIVES} from 'angular2/router';
import {DashboardComponent} from './dashboard.component';
import {ConsulService} from './consul.service';
import {NodeDetailComponent} from './node-detail.component';
import {NodesComponent} from './nodes.component';

@Component({
  selector: 'avast',
  template: `
    <h1>{{title}}</h1>
    <a [routerLink]="['Dashboard']">Dashboard</a>
    <a [routerLink]="['Nodes']">Nodes</a>
    <a [routerLink]="['Containers']">Containers</a>
    <a [routerLink]="['Images']">Images</a>
    <router-outlet></router-outlet>
  `,
  styleUrls: ['app/app.component.css'],
  directives: [ROUTER_DIRECTIVES],
  providers: [ConsulService],
})

@RouteConfig([
  {
  //  path: '/',
  //  redirectTo: ['Dashboard']},
  //}, {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardComponent,
    useAsDefault: true,
  }, {
    path: '/nodes',
    name: 'Nodes',
    component: NodesComponent,
  }, {
    path: '/containers',
    name: 'Containers',
    component: NodesComponent,
  }, {
    path: '/images',
    name: 'Images',
    component: NodesComponent,
  }, {
    path: '/node/detail/:name',
    name: 'NodeDetail',
    component: NodeDetailComponent,
  }
])

export class AppComponent {
  public title = 'Avast'
}
