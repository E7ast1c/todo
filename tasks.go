package main

import (
	"fmt"
)

type ToDoList map[uint16]Task

type CombinedTask struct {
	ID   ID
	Task Task
}

type ID struct {
	Id uint16 `json:"id"`
}

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed,omitempty"`
	Favorites   bool   `json:"favorites,omitempty"`
}

var toDoList ToDoList = map[uint16]Task{}

func GetAll() *ToDoList {
	return &toDoList
}

func Add(nt Task) (*ToDoList, error) {
	for _, et := range toDoList {
		if nt.Title == et.Title {
			return nil, fmt.Errorf("duplicate title %+v", et)
		}
	}

	var ind = uint16(len(toDoList))
	for true {
		if _, ok := toDoList[ind]; !ok {
			break
		}
		ind++
	}
	toDoList[ind] = nt
	return &toDoList, nil
}

func Delete(id uint16) (*ToDoList, error) {
	if _, ok := toDoList[id]; !ok {
		return nil, fmt.Errorf("delete failed: unknown id %d", id)
	}
	delete(toDoList, id)
	return &toDoList, nil
}

func Update(id uint16, nt Task) (*ToDoList, error) {
	if v, ok := toDoList[id]; !ok {
		return nil, fmt.Errorf("update failed: unknown id %d", id)
	} else if v.Title == nt.Title {
		return nil, fmt.Errorf("duplicate title %+v", nt)
	}

	toDoList[id] = nt
	return &toDoList, nil
}
