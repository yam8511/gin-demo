// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'

Vue.config.productionTip = false

// Apollo
import { ApolloClient, createBatchingNetworkInterface } from 'apollo-client'
import VueApollo from 'vue-apollo'

// Create the apollo client
const apolloClient = new ApolloClient({
  networkInterface: createBatchingNetworkInterface({
    uri: 'http://localhost:8000/apollo-graphql'
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
