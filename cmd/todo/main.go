package main

import (
	todo "akwadit/ahmed-todo"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	fileName = ".todos.json"
)

func main() {
	add := flag.Bool("add", false, "Add a new todo")
	complete := flag.Int("complete", 0, "Mark todo as completed")
	del := flag.Int("del", 0, "delete todo")
	list := flag.Bool("list", false, "List todos")
	flag.Parse()
	todos := &todo.Todos{}

	if err := todos.Load(fileName); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:

		task, err := GetInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Add(task)
		err = todos.Store(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *list:
		todos.Print()
	case *complete > 0:
		err := todos.Completed(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		errS := todos.Store(fileName)
		if errS != nil {
			fmt.Fprintln(os.Stderr, errS.Error())
			os.Exit(1)
		}
	case *del > 0:
		err := todos.Delete(*del)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		errS := todos.Store(fileName)
		if errS != nil {
			fmt.Fprintln(os.Stderr, errS.Error())
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stdin, "invalid command")
		os.Exit(0)
	}

}
func GetInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil

	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err

	}

	if len(scanner.Text()) == 0 {
		return "", errors.New("empty task")
	}

	return scanner.Text(), nil

}
