import {Injectable} from 'angular2/core';
import {Http, Response} from 'angular2/http';
import {Actions} from 'angular2-redux';
import 'rxjs/add/operator/map';

import {API_VERSION} from '../constants';

export const REQUEST_DATACENTERS = 'REQUEST_DATACENTERS';
export const RECEIVE_DATACENTERS = 'RECEIVE_DATACENTERS';

@Injectable()
export class ConsulActions extends Actions {

  constructor(private _http: Http) {
    super();
  }

  fetchDatacenters() {
    return (dispatch) => {
      dispatch(this.requestDatacenters());

      this._http.get(`/api/${API_VERSION}/consul/datacenters`)
        .map((res: Response) => res.json())
        .map(res => dispatch(this.receiveDatacenters(res)))
        .subscribe();
    };
  }

  requestDatacenters() {
    return {type: REQUEST_DATACENTERS};
  }

  receiveDatacenters(datacenters: string[]) {
    return {
      type: RECEIVE_DATACENTERS,
      datacenters
    };
  }

}
