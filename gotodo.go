package main

import  (
        "fmt"
        "../go-todotxt"
        "github.com/docopt/docopt.go"
)


func main() {

        usage := `Go Todo.txt

Usage:
    gotodo [--sort=<prio|date|len|prio-rev|date-rev|len-rev|id>]
    gotodo list [--sort=<prio|date|len|prio-rev|date-rev|len-rev|id>] [<keyword>]
    gotodo add <task>
    gotodo (finish|done) <id>
    gotodo --num-tasks
    gotodo -h | --help
    gotodo -v | --version

Options:
    -h --help     Show this screen.
    -v --version  Show version.
    --num-tasks   Show number of tasks.
    -f --finished Show the finished tasks.`

        arguments, _ := docopt.Parse(usage, nil, true, "Go Todo.txt 0.1", false)

        //fmt.Println(arguments)

        tasks := todotxt.LoadTaskList("todo.txt")


        if (arguments["--num-tasks"].(bool)) {
                fmt.Println(tasks.Len())
        } else {

                var by string
                if arguments["--sort"] == nil {
                        by = "prio"
                } else {
                        by = arguments["--sort"].(string)
                }
                tasks.Sort(by)

                for _, task := range tasks {
                        if !task.Finished() {
                                fmt.Println(task.Id(), task.RawText())
                        }
                }
        }

        //fmt.Println(tasks)
        //fmt.Println(tasks.Count())
}
