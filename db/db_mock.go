package db

import (
	"database/sql"

	"errors"

	"gopkg.in/tucnak/telebot.v2"
)

type MockZbotDatabase struct {
	Level             string
	File              string
	Term              string
	Meaning           string
	Author            string
	Date              string
	Find_terms        []string
	Search_terms      []string
	Not_found         bool
	Rand_def          DefinitionItem
	User              telebot.User
	Ignore_list       []string
	User_ignored      []UserIgnore
	Ignore_User       bool
	Error             bool
	ErrorAppend       bool
	IgnoreListCleaned bool
}

func (d *MockZbotDatabase) GetConnectionString() string {
	return "mock"
}

func (d *MockZbotDatabase) Init() error {
	return nil
}

func (d *MockZbotDatabase) Close() {
}

func (d *MockZbotDatabase) UserIgnoreList() ([]UserIgnore, error) {
	if d.Error {
		return nil, errors.New("mock")
	}
	return d.User_ignored, nil
}

func (d *MockZbotDatabase) Statistics() (string, error) {
	if d.Error {
		return "", errors.New("mock")
	}
	return d.Level, nil
}

func (d *MockZbotDatabase) Top() ([]DefinitionItem, error) {
	var items []DefinitionItem

	if d.Error {
		return nil, errors.New("mock")
	}

	for _, value := range d.Find_terms {
		items = append(items, DefinitionItem{Term: value})
	}

	return items, nil
}

func (d *MockZbotDatabase) Rand() (DefinitionItem, error) {
	if d.Error {
		return DefinitionItem{}, errors.New("mock")
	}
	return d.Rand_def, nil
}

func (d *MockZbotDatabase) Last() (DefinitionItem, error) {

	if d.Error {
		return DefinitionItem{}, errors.New("mock")
	}

	return DefinitionItem{Term: d.Term, Meaning: d.Meaning}, nil
}

func (d *MockZbotDatabase) Set(def DefinitionItem) (string, error) {
	if d.Error {
		return "", errors.New("mock")
	}
	return def.Term, nil
}

func (d *MockZbotDatabase) Find(criteria string) ([]DefinitionItem, error) {
	if d.Not_found {
		return []DefinitionItem{}, nil
	}
	if d.Error {
		return nil, errors.New("mock")
	}
	return []DefinitionItem{{Term: d.Term}}, nil
}

func (d *MockZbotDatabase) Get(term string) (DefinitionItem, error) {
	if d.Not_found {
		return DefinitionItem{}, nil
	}
	if d.Error {
		return DefinitionItem{}, errors.New("mock")
	}
	return DefinitionItem{Term: d.Term, Meaning: d.Meaning, Author: d.Author, Date: d.Date}, nil
}

func (d *MockZbotDatabase) _set(term string, def DefinitionItem) (sql.Result, error) {
	var result sql.Result

	if d.Error {
		return nil, errors.New("mock")
	}

	return result, nil
}

func (d *MockZbotDatabase) Search(str string) ([]DefinitionItem, error) {
	var items []DefinitionItem

	if d.Error {
		return []DefinitionItem{}, errors.New("mock")
	}

	for _, value := range d.Search_terms {
		items = append(items, DefinitionItem{Term: value})
	}

	return items, nil
}

//UserLevel Mock
func (d *MockZbotDatabase) UserLevel(str string) (string, error) {
	if d.Error {
		return "", errors.New("mock")
	}
	return d.Level, nil
}

//UserCheckIgnore Mock, it will return false if error is set otherwise it will return Ignore_User value
func (d *MockZbotDatabase) UserCheckIgnore(str string) bool {

	if d.Error {
		return false
	}

	return d.Ignore_User
}

func (d *MockZbotDatabase) UserIgnoreInsert(username string) error {
	return nil
}

//UserCleanupIgnorelist Cleanup ignore list
func (d *MockZbotDatabase) UserCleanupIgnorelist() error {
	d.IgnoreListCleaned = true
	return nil
}

//Lock defintion mock
func (d *MockZbotDatabase) Lock(item DefinitionItem) error {
	if d.Error {
		return errors.New("mock")
	}
	return nil
}

func (d *MockZbotDatabase) Append(item DefinitionItem) error {
	if d.ErrorAppend {
		return errors.New("mock")
	}
	return nil
}

func (d *MockZbotDatabase) Forget(item DefinitionItem) error {
	if d.Error {
		return errors.New("mock")
	}
	return nil
}