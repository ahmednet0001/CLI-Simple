package todo

import (
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

const (
	ColorDefault = "\x1b[39m"

	ColorRed   = "\x1b[91m"
	ColorGreen = "\x1b[32m"
	ColorBlue  = "\x1b[94m"
	ColorGray  = "\x1b[90m"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*t = append(*t, todo)
}
func (t *Todos) Completed(index int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}

	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	return nil
}
func (t *Todos) Delete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("inavlid index")
	}

	*t = append(ls[:index-1], ls[index:]...)

	return nil
}

func (t *Todos) Load(fileName string) error {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)

	if err != nil {
		return err
	}
	return nil
}

func (t *Todos) Store(fileName string) error {

	data, err := json.Marshal(*t)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, data, 0644)
}
func (t *Todos) Print() {
	// table := simpletable.New()

	// table.Header = &simpletable.Header{
	// 	Cells: []*simpletable.Cell{
	// 		{Align: simpletable.AlignCenter, Text: "#"},
	// 		{Align: simpletable.AlignCenter, Text: "Task"},
	// 		{Align: simpletable.AlignCenter, Text: "Done?"},
	// 		{Align: simpletable.AlignCenter, Text: "CreatedAt"},
	// 		{Align: simpletable.AlignCenter, Text: "CompletedAt?"},
	// 	},
	// }

	// table.Body = &simpletable.Body{Cells: cells}
	// table.Footer = &simpletable.Footer{
	// 	Cells: []*simpletable.Cell{
	// 		{Align: simpletable.AlignCenter, Span: 3, Text: "And all columns is well sized"},
	// 	},
	// }
	// table.Print()

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done"},
			{Align: simpletable.AlignCenter, Text: "Created"},
			{Align: simpletable.AlignCenter, Text: "Completed"},
		},
	}

	var cells [][]*simpletable.Cell
	for idx, item := range *t {
		task := blue(item.Task)
		done := red("No")
		created := blue(item.CreatedAt.Format(time.RFC822))
		completed := blue(item.CompletedAt.Format(time.RFC822))
		if item.Done {
			task = green(fmt.Sprintf("\u2705 %s", item.Task))
			done = green("Yes")
			created = green(item.CreatedAt.Format(time.RFC822))
			completed = green(item.CompletedAt.Format(time.RFC822))
		}
		idx++
		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: task},
			{Text: done},
			{Text: created},
			{Text: completed},
		})

	}
	table.Body = &simpletable.Body{
		Cells: cells,
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: gray("Count of pending tasks is : ") + red(fmt.Sprintf("%d", t.PendingTasks()))},
		},
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func (t *Todos) PendingTasks() (total int) {
	for _, item := range *t {
		if !item.Done {
			total++
		}
	}
	return total
}
func red(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}

func green(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorDefault)
}

func blue(s string) string {
	return fmt.Sprintf("%s%s%s", ColorBlue, s, ColorDefault)
}

func gray(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGray, s, ColorDefault)
}
