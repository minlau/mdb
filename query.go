package main

import (
	"encoding/json"
)

type GroupQueryResult struct {
	GroupId int          `json:"groupId"`
	Data    *QueryResult `json:"data"`
	Error   *QueryError  `json:"error"`
}

type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

type QueryError struct {
	Message string `json:"message"`
	Err     error  `json:"err"`
}

func NewQueryError(err error) *QueryError {
	if err == nil {
		return nil
	}
	if errJson, _ := json.Marshal(err); string(errJson) == "{}" {
		return &QueryError{Message: err.Error(), Err: nil}
	}
	return &QueryError{Message: err.Error(), Err: err}
}
