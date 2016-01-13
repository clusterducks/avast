System.register(['angular2/core', 'angular2/router', './consul.service'], function(exports_1) {
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, router_1, consul_service_1;
    var NodeDetailComponent;
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
            }],
        execute: function() {
            NodeDetailComponent = (function () {
                function NodeDetailComponent(_consulService, _routeParams) {
                    this._consulService = _consulService;
                    this._routeParams = _routeParams;
                }
                NodeDetailComponent.prototype.ngOnInit = function () {
                    var _this = this;
                    if (!this.node) {
                        var name_1 = this._routeParams.get('name');
                        this._consulService.getNode(name_1)
                            .subscribe(function (res) { return _this.node = res; });
                    }
                };
                NodeDetailComponent.prototype.goBack = function () {
                    window.history.back();
                };
                NodeDetailComponent = __decorate([
                    core_1.Component({
                        selector: 'avast-node-detail',
                        templateUrl: 'app/node-detail.component.html',
                        styleUrls: ['app/node-detail.component.css'],
                    }), 
                    __metadata('design:paramtypes', [consul_service_1.ConsulService, router_1.RouteParams])
                ], NodeDetailComponent);
                return NodeDetailComponent;
            })();
            exports_1("NodeDetailComponent", NodeDetailComponent);
        }
    }
});
//# sourceMappingURL=node-detail.component.js.map