import {Injectable} from 'angular2/core';
import {Http, Response} from 'angular2/http';

@Injectable()
export class ConsulService {

  constructor(public http: Http) {
  }

  getDatacenters() {
    return this.http.get('http://localhost:8080/consul/datacenters')
      .map((res: Response) => res.json());
  }

  getNodes(dc: string='') {
    let url = 'http://localhost:8080/consul/nodes' + (dc ? '/' + dc : '');
    return this.http.get(url)
      .map((res: Response) => res.json());
  }

  getNode(name: string) {
    return this.http.get('http://localhost:8080/consul/node/' + name)
      .map((res: Response) => res.json());
  }
}
