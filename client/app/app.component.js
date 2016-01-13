System.register(['angular2/core', 'angular2/router', './dashboard.component', './consul.service', './node-detail.component', './nodes.component'], function(exports_1) {
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, router_1, dashboard_component_1, consul_service_1, node_detail_component_1, nodes_component_1;
    var AppComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (router_1_1) {
                router_1 = router_1_1;
            },
            function (dashboard_component_1_1) {
                dashboard_component_1 = dashboard_component_1_1;
            },
            function (consul_service_1_1) {
                consul_service_1 = consul_service_1_1;
            },
            function (node_detail_component_1_1) {
                node_detail_component_1 = node_detail_component_1_1;
            },
            function (nodes_component_1_1) {
                nodes_component_1 = nodes_component_1_1;
            }],
        execute: function() {
            AppComponent = (function () {
                function AppComponent() {
                    this.title = 'Avast';
                }
                AppComponent = __decorate([
                    core_1.Component({
                        selector: 'avast',
                        template: "\n    <h1>{{title}}</h1>\n    <a [routerLink]=\"['Dashboard']\">Dashboard</a>\n    <a [routerLink]=\"['Nodes']\">Nodes</a>\n    <a [routerLink]=\"['Containers']\">Containers</a>\n    <a [routerLink]=\"['Images']\">Images</a>\n    <router-outlet></router-outlet>\n  ",
                        styleUrls: ['app/app.component.css'],
                        directives: [router_1.ROUTER_DIRECTIVES],
                        providers: [consul_service_1.ConsulService],
                    }),
                    router_1.RouteConfig([
                        {
                            //  path: '/',
                            //  redirectTo: ['Dashboard']},
                            //}, {
                            path: '/dashboard',
                            name: 'Dashboard',
                            component: dashboard_component_1.DashboardComponent,
                            useAsDefault: true,
                        }, {
                            path: '/nodes',
                            name: 'Nodes',
                            component: nodes_component_1.NodesComponent,
                        }, {
                            path: '/containers',
                            name: 'Containers',
                            component: nodes_component_1.NodesComponent,
                        }, {
                            path: '/images',
                            name: 'Images',
                            component: nodes_component_1.NodesComponent,
                        }, {
                            path: '/node/detail/:name',
                            name: 'NodeDetail',
                            component: node_detail_component_1.NodeDetailComponent,
                        }
                    ]), 
                    __metadata('design:paramtypes', [])
                ], AppComponent);
                return AppComponent;
            })();
            exports_1("AppComponent", AppComponent);
        }
    }
});
//# sourceMappingURL=app.component.js.map