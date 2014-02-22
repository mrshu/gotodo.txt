package main

import  (
        "fmt"
        "../go-todotxt"
        "github.com/spf13/cobra"
        "os/user"
        "strings"
        "strconv"
)

func extendedLoader(filename string) (todotxt.TaskList, error) {
        usr, err := user.Current()
        if err != nil {
                return nil, err
        }

        filename = strings.Replace(filename, "~", usr.HomeDir, -1)
        tasks := todotxt.LoadTaskList(filename)

        return tasks, nil
}

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
                tasks, err := extendedLoader(filename)
                if err != nil {
                        fmt.Println(err)
                        return
                }

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
                            task.SetIdPaddingBy(tasks)
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

        var cmdAdd = &cobra.Command{
            Use:   "add [task]",
            Short: "Adds a task to the todo list.",
            Long:  `Adds a task to the todo list.`,
            Run: func(cmd *cobra.Command, args []string) {
                tasks, err := extendedLoader(filename)
                if err != nil {
                        fmt.Println(err)
                        return
                }

                task := strings.Join(args, " ")
                tasks.Add(task)

                tasks.Save(filename)
            },
        }

        var nofinishdate bool
        var cmdDone = &cobra.Command{
            Use:   "done [taskid]",
            Short: "Marks task as done.",
            Long:  `Marks task as done.`,
            Run: func(cmd *cobra.Command, args []string) {
                tasks, err := extendedLoader(filename)
                if err != nil {
                        fmt.Println(err)
                        return
                }

                if len(args) < 1 {
                        fmt.Println("So what needs to be done?")
                        return
                }

                taskid, err := strconv.Atoi(args[0])
                if err != nil {
                        fmt.Printf("Do you really consider that a number? %v\n", err)
                        return
                }

                err = tasks.Done(taskid, !nofinishdate)
                if err != nil {
                        fmt.Printf("There was an error %v\n", err)
                }

                tasks.Save(filename)
            },
        }
        cmdDone.Flags().BoolVarP(&nofinishdate, "no-finish-date", "D", false,
                                        "Do not mark finished tasks with date.")


        var editprio string
        var cmdEdit = &cobra.Command{
            Use:   "edit [taskid]",
            Short: "Edits given task.",
            Long:  `Edits given task.`,
            Run: func(cmd *cobra.Command, args []string) {
                tasks, err := extendedLoader(filename)
                if err != nil {
                        fmt.Println(err)
                        return
                }

                if len(args) < 1 {
                        fmt.Println("So what do you want to edit?")
                        return
                }

                taskid, err := strconv.Atoi(args[0])
                if err != nil {
                        fmt.Printf("Do you really consider that a number? %v\n", err)
                        return
                }

                if len(editprio) > 0 {
                        tasks[taskid].SetPriority(editprio[0])
                        tasks[taskid].RebuildRawTodo()
                }

                tasks.Save(filename)
            },
        }

        cmdEdit.PersistentFlags().StringVarP(&editprio, "priority", "p", "",
                                     "Sets task's priority.")


        var GotodoCmd = &cobra.Command{
            Use:   "gotodo",
            Short: "Gotodo is a go implementation of todo.txt.",
            Long: `A small, fast and fun implementation of todo.txt`,
            Run: func(cmd *cobra.Command, args []string) {
                cmdList.Run(cmd, nil)
            },
        }

        GotodoCmd.PersistentFlags().StringVarP(&filename, "filename", "", "todo.txt",
                                     "Load tasks from this file.")

        GotodoCmd.AddCommand(cmdList)
        GotodoCmd.AddCommand(cmdAdd)
        GotodoCmd.AddCommand(cmdDone)
        GotodoCmd.Execute()
}
