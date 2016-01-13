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
    var DashboardComponent;
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
            DashboardComponent = (function () {
                function DashboardComponent(_consulService, _router) {
                    this._consulService = _consulService;
                    this._router = _router;
                }
                DashboardComponent.prototype.ngOnInit = function () {
                    this.getDatacenters();
                };
                DashboardComponent.prototype.getDatacenters = function () {
                    var _this = this;
                    this._consulService.getDatacenters()
                        .subscribe(function (res) { return _this.datacenters = res; }, function (err) { return _this.logError(err); });
                };
                DashboardComponent.prototype.setDatacenter = function (dc) {
                    var _this = this;
                    this._consulService.getNodes(dc)
                        .subscribe(function (res) { return _this.nodes = res; });
                };
                DashboardComponent.prototype.gotoNode = function (name) {
                    this._router.navigate(['NodeDetail', { name: name }]);
                };
                DashboardComponent.prototype.logError = function (err) {
                    console.log(err);
                };
                DashboardComponent = __decorate([
                    core_1.Component({
                        selector: 'avast-dashboard',
                        templateUrl: 'app/dashboard.component.html',
                        styleUrls: ['app/dashboard.component.css'],
                    }), 
                    __metadata('design:paramtypes', [consul_service_1.ConsulService, router_1.Router])
                ], DashboardComponent);
                return DashboardComponent;
            })();
            exports_1("DashboardComponent", DashboardComponent);
        }
    }
});
//# sourceMappingURL=dashboard.component.js.map