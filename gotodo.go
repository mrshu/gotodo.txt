package main

import  (
        "fmt"
        "../go-todotxt"
        "github.com/spf13/cobra"
)


func main() {

        var numtasks bool
        var sortby string
        var finished bool
        var prettyformat string

        var cmdList = &cobra.Command{
            Use:   "list [keyword]",
            Short: "Lists tasks that contain keyword, if any",
            Long:  `List is the most basic command that is used for listing tasks.
                    You can specify a keyword as well as other options.`,
            Run: func(cmd *cobra.Command, args []string) {
                tasks := todotxt.LoadTaskList("todo.txt")

                if numtasks {
                    fmt.Println(tasks.Len())
                } else {
                    tasks.Sort(sortby)

                    for _, task := range tasks {
                        if (!task.Finished() && !finished) ||
                           (task.Finished() && finished) {
                           fmt.Println(task.PrettyPrint(prettyformat))
                        }
                    }
                }
            },
        }
        cmdList.Flags().BoolVarP(&numtasks, "num-tasks", "n", false,
                                 "Show the number of tasks")
        cmdList.Flags().BoolVarP(&finished, "finished", "f", false,
                                 "Show finished tasks")
        cmdList.Flags().StringVarP(&sortby, "sort", "s", "prio",
                                   "Sort tasks by parameter (prio|date|len|prio-rev|date-rev|len-rev|id)")
        cmdList.Flags().StringVarP(&prettyformat, "pretty", "", "%i %p %t",
                                   "Pretty print tasks")

        var GotodoCmd = &cobra.Command{
            Use:   "gotodo",
            Short: "Gotodo is a go implementation of todo.txt.",
            Long: `A small, fast and fun implementation of todo.txt`,
            Run: func(cmd *cobra.Command, args []string) {
                cmdList.Run(cmd, nil)
            },
        }


//              usage := `Go Todo.txt

//      Usage:
//          gotodo [--sort=prio|date|len|prio-rev|date-rev|len-rev|id]
//          gotodo list [--sort=prio|date|len|prio-rev|date-rev|len-rev|id] [<keyword>]
//          gotodo add <task>
//          gotodo (finish|done) <id>
//          gotodo --num-tasks
//          gotodo -h | --help
//          gotodo -v | --version

//      Options:
//          -h --help      Show this screen.
//          -v --version   Show version.
//          --num-tasks    Show number of tasks.
//          --pretty=STR   Pretty print the task list.
//          -f --finished  Show the finished tasks.`

        //arguments, _ := docopt.Parse(usage, nil, true, "Go Todo.txt 0.1", false)

        //fmt.Println(arguments)


        GotodoCmd.AddCommand(cmdList)
        GotodoCmd.Execute()

      //if (arguments["--num-tasks"].(bool)) {
      //        fmt.Println(tasks.Len())
      //} else {

      //        var by string
      //        if arguments["--sort"] == nil {
      //                by = "prio"
      //        } else {
      //                by = arguments["--sort"].(string)
      //        }
      //        tasks.Sort(by)

      //        for _, task := range tasks {
      //                if !task.Finished() {
      //                        fmt.Println(task.Id(), task.RawText())
      //                }
      //        }
      //}

        //fmt.Println(tasks)
        //fmt.Println(tasks.Count())
}
