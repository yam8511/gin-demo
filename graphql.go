package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

// Todo : 工作項目
type Todo struct {
	ID   string `json:"id"`
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
			Type: graphql.String,
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
		Todo{ID: "a", Text: "A todo not to forget", Done: true},
		Todo{ID: "b", Text: "This is the most important", Done: false},
		Todo{ID: "c", Text: "Please do this or else", Done: false},
	}
}

// GraphQLHandle : GraphQL Schema
func GraphQLHandle(c *gin.Context) {
	todoInit()
	rootQuery := schemaQuerySetting()

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contentType := c.Request.Header.Get("content-type")
	queryString := ""
	if contentType == "application/json" {
		jsonData := map[string]interface{}{}
		err = json.NewDecoder(c.Request.Body).Decode(&jsonData)
		if err == nil {
			queryString, _ = jsonData["query"].(string)
		}
	} else {
		queryString = c.Query("query")
	}

	result := executeQuery(queryString, schema)

	c.JSON(http.StatusOK, result)
}

// GraphIQLHandle : GraphQL Interview
func GraphIQLHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "graphiql.html", nil)
}

// 設定 Query 設定
func schemaQuerySetting() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"todo": &graphql.Field{
				Type:        todoType,
				Description: "取單一個 todo",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
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

// todo 處理區
func todoResolveFn(params graphql.ResolveParams) (interface{}, error) {
	idQuery, isOK := params.Args["id"].(string)
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

// 執行Query查詢
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	return result
}
