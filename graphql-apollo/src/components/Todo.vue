<template>
  <div class="todo">
    <h1>{{ msg }}</h1>
    <h2 v-if="isLoading > 0">Loading...</h2>
    <div class="current"
      v-if="currentTodo != null"
      @click="toggleTodo"
    >
      Current Todo: <br/>
      ID: {{ currentTodo.id }}<br/>
      Text: {{ currentTodo.text }}<br/>
      Done: {{ currentTodo.done }}
    </div>
    <input placeholder="Create New Todo .."
      v-model="inputText"
      @keyup.enter="createTodo"
    />
    <ul>
      <li v-for="todo in todos"
        :key="todo.id"
        :class="todo.done ? 'done' : ''"
        @click="focusTodo(todo)"
      >
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
      variables: {
        Done: true
      },
      // query
      query: gql`{
        todos {
          id
          text
          done
        }
      }`,
      // receive the specify graphql data
      update (data) {
        console.info('todos update', data)
        // this.isLoading = false
        return data.todos
      },
      // get all the result info
      result (data) {
        console.group('Apollo Result: Todo')
        console.info('all', data)
        console.info('data', data.data)
        console.info('loader', data.loading)
        console.info('network status', data.networkStatus)
        console.info('stale', data.stale)
        console.groupEnd()
      },
      // if error
      error (error) {
        console.warn('We\'ve got an error!', error)
      },
      // bind custom loading key
      loadingKey: 'isLoading',
      // watch the loading status
      watchLoading (isLoading, countModifier) {
        console.log('watch loading', isLoading, countModifier)
        console.info('loading code', this.isLoading)
      }
      // Polling Query
      // pollInterval: 300 // ms
    }
  },
  data () {
    return {
      msg: 'GraphQL Demo - Todo',
      isLoading: 0,
      todos: [],
      currentTodo: null,
      inputText: ''
    }
  },
  methods: {
    focusTodo (todo) {
      this.currentTodo = todo
    },
    toggleTodo () {
      this.todos = this.todos.map((todo) => {
        if (todo.id === this.currentTodo.id) {
          this.currentTodo = {
            ...todo,
            done: !todo.done
          }
          return this.currentTodo
        }
        return todo
      })
    },
    createTodo () {
      if (this.inputText === '') {
        return
      }
      const newText = this.inputText
      console.log('ntext', newText)
      this.inputText = ''

      this.$apollo.mutate({
        // Query
        mutation: gql`mutation createTodo($Text: String!) {
          createTodo(text: $Text) {
            text
            done
            id
          }
        }`,
        variables: {
          Text: newText
        },
        update (data) {
          console.info('mutation update:', data)
        }
      })
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

li {
  color:blue;
  cursor: pointer;
}

li.done {
  /* display: inline-block; */
  color: #42b983;
  margin: 0 10px;
}

.current {
  color:orange;
  cursor: pointer;
}
</style>
