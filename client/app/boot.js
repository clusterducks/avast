System.register(['angular2/platform/browser', 'angular2/http', 'rxjs/add/operator/map', 'angular2/router', './app.component', './consul.service'], function(exports_1) {
    var browser_1, http_1, router_1, app_component_1, consul_service_1;
    return {
        setters:[
            function (browser_1_1) {
                browser_1 = browser_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
            },
            function (_1) {},
            function (router_1_1) {
                router_1 = router_1_1;
            },
            function (app_component_1_1) {
                app_component_1 = app_component_1_1;
            },
            function (consul_service_1_1) {
                consul_service_1 = consul_service_1_1;
            }],
        execute: function() {
            browser_1.bootstrap(app_component_1.AppComponent, [
                http_1.HTTP_PROVIDERS,
                router_1.ROUTER_PROVIDERS,
                consul_service_1.ConsulService,
            ]);
        }
    }
});
//# sourceMappingURL=boot.js.map