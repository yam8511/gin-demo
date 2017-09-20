package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

// GraphQLType : GraphQL 基本型態
type GraphQLType struct {
	Query         string `json:"query"`
	OperationName string `json:"operationName"`
	Variables     string `json:"variables"`
}

// ApolloGraphQLType : GraphQL 基本型態
type ApolloGraphQLType struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// ApolloGraphQLHandle : GraphQL
func ApolloGraphQLHandle(schema *graphql.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			return
		}

		resultList := []*graphql.Result{}
		apolloData := []ApolloGraphQLType{}
		err := json.NewDecoder(c.Request.Body).Decode(&apolloData)
		if err != nil {
			log.Println("APOLLO DECODE ERROR", err)
		} else {
			for _, data := range apolloData {
				result := executeQuery(data.OperationName, data.Query, data.Variables, *schema)
				resultList = append(resultList, result)
			}
		}
		c.JSON(http.StatusOK, resultList)
		return
	}
}

// GraphQLHandle : GraphQL Schema
func GraphQLHandle(schema *graphql.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			return
		}

		contentType := c.Request.Header.Get("content-type")
		var result interface{}
		if contentType == "application/json" {
			graphqlData := GraphQLType{}
			err := json.NewDecoder(c.Request.Body).Decode(&graphqlData)
			if err == nil {
				queryString := graphqlData.Query
				operationName := graphqlData.OperationName
				variablesString := graphqlData.Variables
				variables := map[string]interface{}{}
				json.Unmarshal([]byte(variablesString), &variables)
				result = executeQuery(operationName, queryString, variables, *schema)
			}
		} else {
			queryString := c.Query("query")
			result = executeQuery("", queryString, map[string]interface{}{}, *schema)
		}

		c.JSON(http.StatusOK, result)
	}
}

// 執行Query查詢
func executeQuery(operationName string, query string, variables map[string]interface{}, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		OperationName:  operationName,
		RequestString:  query,
		VariableValues: variables,
	})

	return result
}
