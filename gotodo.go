package main

import  (
        "fmt"
        "flag"
        "../go-todotxt"
        "github.com/docopt/docopt.go"
)

var show_task_count = flag.Bool("n", false, "Show just task count")

func main() {

        usage := `Go Todo.txt

Usage:
    gotodo
    gotodo list
    gotodo add <task>
    gotodo (finish|done) <id>

Options:
    -h --help     Show this screen.
    --version     Show version.`

        arguments, _ := docopt.Parse(usage, nil, true, "Go Todo.txt 0.1", false)

        fmt.Println(arguments)

        flag.Parse()
        tasks := todotxt.LoadTaskList("todo.txt")

        if (*show_task_count) {
                fmt.Println(tasks.Count())
        } else {
                for i, task := range tasks {
                        fmt.Println(i, task.Text())
                }
        }

       //fmt.Println(tasks)
        //fmt.Println(tasks.Count())
}
