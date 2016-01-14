import {Injectable} from 'angular2/core';
import {Http, Response} from 'angular2/http';

// @TODO: pull this out into a UrlBuilder service
import {API_VERSION} from '../../../constants';

@Injectable()
export class ConsulService {

  constructor(public http: Http) {
  }

  getDatacenters() {
    return this.http.get('/api/' + API_VERSION + '/consul/datacenters')
      .map((res: Response) => res.json());
  }

  getNodes(dc: string = '') {
    let url = '/api/' + API_VERSION + '/consul/nodes' + (dc ? '/' + dc : '');
    return this.http.get(url)
      .map((res: Response) => res.json());
  }

  getNode(name: string) {
    return this.http.get('/api/' + API_VERSION + '/consul/node/' + name)
      .map((res: Response) => res.json());
  }
}
