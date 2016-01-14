import * as ConsulActions from '../actions/consul.actions';

export default (state = [], action: any = {}) => {
  switch (action.type) {
    case ConsulActions.REQUEST_DATACENTERS:
      return Object.assign({}, state, {
        isFetchingDatacenters: true
      });
    case ConsulActions.RECEIVE_DATACENTERS:
      return Object.assign({}, state, {
        isFetchingDatacenters: false,
        list: action.datacenters
      });
    default:
      return state;
  }
};
