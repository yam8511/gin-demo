<template>
  <div class="todo">
    <h1>{{ msg }}</h1>
    <h2 v-if="isLoading">Loading...</h2>
    <ul>
      <li v-for="todo in todos">
        {{ todo.text }}
      </li>
    </ul>
  </div>
</template>

<script>
import gql from 'graphql-tag'
export default {
  name: 'todo',
  apollo: {
    todos: {
      query: gql`{
        todos {
          id
          text
          done
        }
      }`
    }
  },
  data () {
    return {
      msg: 'GraphQL Demo - Todo',
      isLoading: true,
      todos: []
    }
  },
  methods: {
    updateLoading (loading) {
      console.log('before', this.isLoading)
      this.isLoading = loading
      console.log('after', this.isLoading)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1,
h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li.done {
  /* display: inline-block; */
  color: #42b983;
  margin: 0 10px;
}
</style>
