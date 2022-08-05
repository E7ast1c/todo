package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	Title       string `json:"title" validate:"required,min=3,max=32"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed,omitempty"`
	Favorites   bool   `json:"favorites,omitempty"`
}

var toDoList ToDoList = map[uint16]Task{}

func GetAll() *ToDoList {
	return &toDoList
}

func Add(c *fiber.Ctx) error {
	t := new(Task)
	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if errors := ValidateStruct(t); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	for _, et := range toDoList {
		if t.Title == et.Title {
			return c.Status(fiber.StatusBadRequest).JSON(fmt.Errorf("duplicate title: %+v", et))
		}
	}

	toDoList[uint16(len(toDoList)+1)] = *t
	return c.Status(fiber.StatusOK).JSON(&toDoList)
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
