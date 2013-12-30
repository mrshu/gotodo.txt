package main

import  (
        "fmt"
        "flag"
        "../go-todotxt"
)

var show_task_count = flag.Bool("n", false, "Show just task count")

func main() {
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
