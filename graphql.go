package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

// Schema : 腳本
var Schema graphql.Schema

// Todo : 工作項目
type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// TodoList : 儲存工作項目清單
var TodoList []Todo

// 建立 Todo Type
var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"text": &graphql.Field{
			Type: graphql.String,
		},
		"done": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

func todoInit() {
	TodoList = []Todo{
		Todo{ID: 1, Text: "A todo not to forget", Done: true},
		Todo{ID: 2, Text: "This is the most important", Done: false},
		Todo{ID: 3, Text: "Please do this or else", Done: false},
	}
}

func schemaInit() {
	rootQuery := schemaQuerySetting()
	rootMutation := schemaMutationSetting()

	Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
}

// ApolloGraphQLHandle : GraphQL
func ApolloGraphQLHandle(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.Status(http.StatusNoContent)
		return
	}

	resultList := []*graphql.Result{}
	apolloData := []map[string]interface{}{}
	err := json.NewDecoder(c.Request.Body).Decode(&apolloData)
	if err != nil {
		log.Println("APOLLO DECODE ERROR", err)
	} else {
		for _, data := range apolloData {
			query, _ := data["query"].(string)
			resultList = append(resultList, executeQuery(query, Schema))
		}
	}
	c.JSON(http.StatusOK, resultList)
	return
}

// GraphQLHandle : GraphQL Schema
func GraphQLHandle(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.Status(http.StatusNoContent)
		return
	}

	contentType := c.Request.Header.Get("content-type")
	queryString := ""
	if contentType == "application/json" {
		jsonData := map[string]interface{}{}
		err := json.NewDecoder(c.Request.Body).Decode(&jsonData)
		if err == nil {
			log.Println("json", jsonData)
			queryString, _ = jsonData["query"].(string)
		}
	} else {
		queryString = c.Query("query")
	}

	result := executeQuery(queryString, Schema)

	c.JSON(http.StatusOK, result)
}

// GraphIQLHandle : GraphQL Interview
func GraphIQLHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "graphiql.html", nil)
}

// 設定 Query
func schemaQuerySetting() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"todo": &graphql.Field{
				Type:        todoType,
				Description: "取單一個 todo",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: todoResolveFn,
			},
			"todos": &graphql.Field{
				Type:        graphql.NewList(todoType),
				Description: "取多個 todo",
				Args: graphql.FieldConfigArgument{
					"done": &graphql.ArgumentConfig{
						Type: graphql.Boolean,
					},
				},
				Resolve: todosResolveFn,
			},
		},
	})
}

// 設定 Mutation
func schemaMutationSetting() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createTodo": &graphql.Field{
				Type:        todoType,
				Description: "新增 Todo",
				Args: graphql.FieldConfigArgument{
					"text": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "要新增的 Todo 文字",
					},
				},
				Resolve: actionTodoResolveFn,
			},
			"updateTodo": &graphql.Field{
				Type:        todoType,
				Description: "更新 Todo",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: "要更新的 Todo ID",
					},
					"done": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Boolean),
						Description: "要更新的狀態",
					},
				},
				Resolve: actionTodoResolveFn,
			},
			"deleteTodo": &graphql.Field{
				Type:        todoType,
				Description: "刪除 Todo",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: "要刪除的 Todo ID",
					},
				},
				Resolve: actionTodoResolveFn,
			},
		},
	})
}

// todo 處理區
func todoResolveFn(params graphql.ResolveParams) (interface{}, error) {
	idQuery, isOK := params.Args["id"].(int)
	if isOK {
		for _, todo := range TodoList {
			if todo.ID == idQuery {
				return todo, nil
			}
		}
	}

	return Todo{}, nil
}

// todos 處理區
func todosResolveFn(params graphql.ResolveParams) (interface{}, error) {
	doneQuery, isOK := params.Args["done"].(bool)
	if isOK {
		doneList := []Todo{}
		for _, todo := range TodoList {
			if todo.Done == doneQuery {
				doneList = append(doneList, todo)
			}
		}
		return doneList, nil
	}

	return TodoList, nil
}

// todo 動作區
func actionTodoResolveFn(params graphql.ResolveParams) (interface{}, error) {
	text, isCreateTodo := params.Args["text"].(string)
	done, isUpdateTodo := params.Args["done"].(bool)
	id, isDeleteTodo := params.Args["id"].(int)

	if isCreateTodo {
		newTodo := Todo{
			ID:   TodoList[len(TodoList)-1].ID + 1,
			Text: text,
			Done: false,
		}
		TodoList = append(TodoList, newTodo)
		return newTodo, nil
	}

	if isUpdateTodo {
		for index, todo := range TodoList {
			if todo.ID == id {
				TodoList[index].Done = done
				return TodoList[index], nil
			}
		}
		return Todo{}, errors.New("Todo not exists")
	}

	if isDeleteTodo {
		for index, todo := range TodoList {
			if id == todo.ID {
				TodoList = append(TodoList[:index], TodoList[index+1:]...)
				return todo, nil
			}
		}
		return Todo{}, errors.New("Todo not exists")
	}

	return Todo{}, errors.New("Incorrect Mutation")
}

// 執行Query查詢
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	return result
}
