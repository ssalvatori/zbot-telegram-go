package db

import "database/sql"

type ZbotDatabase interface {
	Init() error
	Close()
	Statistics() (string, error)
	Top() ([]DefinitionItem, error)
	Rand() (DefinitionItem, error)
	Last() (DefinitionItem, error)
	Get(string) (DefinitionItem, error)
	Set(DefinitionItem) error
	_set(string, DefinitionItem) (sql.Result, error)
	Find(string) ([]DefinitionItem, error)
	Search(string) ([]DefinitionItem, error)
	UserLevel(string) (string, error)
	UserIgnoreInsert(string) error
	UserCheckIgnore(string) (bool, error)
	UserCleanIgnore() error
	UserIgnoreList() ([]UserIgnore, error)
}

type DefinitionItem struct {
	Term    string
	Meaning string
	Author  string
	Date    string
	Id      int
}

type UserIgnore struct {
	Username string
	Since    string
	Until    string
}
