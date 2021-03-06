// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'

Vue.config.productionTip = false

// Apollo
import { ApolloClient, createBatchingNetworkInterface } from 'apollo-client'
import VueApollo from 'vue-apollo'

// Create the apollo client
var port = location.port
if (port !== '') {
  port = ':' + port
}
var url = location.protocol + '//' + location.host + port + '/apollo-graphql'

const apolloClient = new ApolloClient({
  networkInterface: createBatchingNetworkInterface({
    uri: url
  }),
  connectToDevTools: true
})

// Install the vue plugin
Vue.use(VueApollo)

// Create the apollo provider
const apolloProvider = new VueApollo({
  defaultClient: apolloClient
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  apolloProvider,
  template: '<App/>',
  components: {
    App
  }
})
