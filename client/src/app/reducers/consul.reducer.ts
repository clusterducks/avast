import * as ConsulActions from '../actions/consul.actions';

export default (state = [], action: any = {}) => {
  switch (action.type) {
    case ConsulActions.REQUEST_DATACENTERS:
      return Object.assign({}, state, {
        isFetchingDatacenters: true
      });

    case ConsulActions.RECEIVE_DATACENTERS:
      return Object.assign({}, state, {
        datacenters: action.datacenters,
        isFetchingDatacenters: false
      });

    default:
      return state;
  }
};
