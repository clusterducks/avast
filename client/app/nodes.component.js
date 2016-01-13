System.register(['angular2/core', 'angular2/router', './consul.service', './node-detail.component'], function(exports_1) {
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, router_1, consul_service_1, node_detail_component_1;
    var NodesComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (router_1_1) {
                router_1 = router_1_1;
            },
            function (consul_service_1_1) {
                consul_service_1 = consul_service_1_1;
            },
            function (node_detail_component_1_1) {
                node_detail_component_1 = node_detail_component_1_1;
            }],
        execute: function() {
            NodesComponent = (function () {
                function NodesComponent(_consulService, _router) {
                    this._consulService = _consulService;
                    this._router = _router;
                }
                NodesComponent.prototype.ngOnInit = function () {
                    this.getNodes();
                };
                NodesComponent.prototype.getNodes = function () {
                    var _this = this;
                    this._consulService.getNodes()
                        .subscribe(function (res) { return _this.nodes = res; });
                };
                NodesComponent.prototype.onSelect = function (node) {
                    this.selectedNode = node;
                };
                NodesComponent.prototype.gotoDetail = function () {
                    this._router.navigate(['NodeDetail', {
                            name: this.selectedNode.node
                        }]);
                };
                NodesComponent = __decorate([
                    core_1.Component({
                        selector: 'avast-nodes',
                        templateUrl: 'app/nodes.component.html',
                        styleUrls: ['app/nodes.component.css'],
                        directives: [node_detail_component_1.NodeDetailComponent]
                    }), 
                    __metadata('design:paramtypes', [consul_service_1.ConsulService, router_1.Router])
                ], NodesComponent);
                return NodesComponent;
            })();
            exports_1("NodesComponent", NodesComponent);
        }
    }
});
//# sourceMappingURL=nodes.component.js.map