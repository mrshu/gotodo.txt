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
        var filename string

        var cmdList = &cobra.Command{
            Use:   "list [keyword]",
            Short: "Lists tasks that contain keyword, if any",
            Long:  `List is the most basic command that is used for listing tasks.
                    You can specify a keyword as well as other options.`,
            Run: func(cmd *cobra.Command, args []string) {
                tasks := todotxt.LoadTaskList(filename)

                if numtasks {
                    fmt.Println(tasks.Len())
                } else {
                    tasks.Sort(sortby)

                    var filteredTasks todotxt.TaskList
                    for _, task := range tasks {
                        if (!task.Finished() && !finished) ||
                           (task.Finished() && finished) {
                           filteredTasks = append(filteredTasks, task)
                        }
                    }

                    for _, task := range filteredTasks {
                            task.SetIdPaddingBy(filteredTasks)
                            fmt.Println(task.PrettyPrint(prettyformat))
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

        GotodoCmd.Flags().StringVarP(&filename, "filename", "", "todo.txt",
                                     "Load tasks from this file.")

        GotodoCmd.AddCommand(cmdList)
        GotodoCmd.Execute()
}
