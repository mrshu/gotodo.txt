package main

import  (
        "fmt"
        "../go-todotxt"
        "github.com/docopt/docopt.go"
)


func main() {

        usage := `Go Todo.txt

Usage:
    gotodo [--sort=<prio|date|len>]
    gotodo list [--sort=<prio|date|len>]
    gotodo add <task>
    gotodo (finish|done) <id>
    gotodo --num-tasks
    gotodo -h | --help
    gotodo -v | --version

Options:
    -h --help     Show this screen.
    -v --version  Show version.
    --num-tasks   Show number of tasks.`

        arguments, _ := docopt.Parse(usage, nil, true, "Go Todo.txt 0.1", false)

        //fmt.Println(arguments)

        tasks := todotxt.LoadTaskList("todo.txt")


        if (arguments["--num-tasks"].(bool)) {
                fmt.Println(tasks.Len())
        } else {

                if arguments["--sort"] == nil {
                        tasks.Sort()
                } else {
                        switch arguments["--sort"].(string) {
                        case "prio":
                                tasks.Sort()
                        case "date":
                                tasks.SortByCreateDate()
                        case "len":
                                tasks.SortByLength()
                        default:
                                tasks.Sort()
                        }
                }

                for i, task := range tasks {
                        fmt.Println(i, task.RawText())
                }
        }

        //fmt.Println(tasks)
        //fmt.Println(tasks.Count())
}
