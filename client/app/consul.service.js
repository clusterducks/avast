System.register(['angular2/core', 'angular2/http'], function(exports_1) {
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, http_1;
    var ConsulService;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
            }],
        execute: function() {
            ConsulService = (function () {
                function ConsulService(http) {
                    this.http = http;
                }
                ConsulService.prototype.getDatacenters = function () {
                    return this.http.get('http://localhost:8080/consul/datacenters')
                        .map(function (res) { return res.json(); });
                };
                ConsulService.prototype.getNodes = function (dc) {
                    if (dc === void 0) { dc = ''; }
                    var url = 'http://localhost:8080/consul/nodes' + (dc ? '/' + dc : '');
                    return this.http.get(url)
                        .map(function (res) { return res.json(); });
                };
                ConsulService.prototype.getNode = function (name) {
                    return this.http.get('http://localhost:8080/consul/node/' + name)
                        .map(function (res) { return res.json(); });
                };
                ConsulService = __decorate([
                    core_1.Injectable(), 
                    __metadata('design:paramtypes', [http_1.Http])
                ], ConsulService);
                return ConsulService;
            })();
            exports_1("ConsulService", ConsulService);
        }
    }
});
//# sourceMappingURL=consul.service.js.map