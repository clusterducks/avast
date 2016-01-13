import {provide} from 'angular2/core';
import {bootstrap, ELEMENT_PROBE_PROVIDERS} from 'angular2/platform/browser';
import {ROUTER_PROVIDERS, LocationStrategy, HashLocationStrategy} from 'angular2/router';
import {HTTP_PROVIDERS} from 'angular2/http';
import 'rxjs/add/operator/map';

import {AppComponent} from './app/components/app.component';
import {ConsulService} from './app/components/consul/providers/consul.service';

document.addEventListener('DOMContentLoaded', function main() {
  bootstrap(AppComponent, [
    ...(process.env.ENV === 'production' ? [] : ELEMENT_PROBE_PROVIDERS),
    ...HTTP_PROVIDERS,
    ...ROUTER_PROVIDERS,
    ConsulService,
    provide(LocationStrategy, {
      useClass: HashLocationStrategy
    })
  ])
  .catch(err => console.error(err));
});
