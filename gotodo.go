package main

import  (
        "fmt"
        "os"
        "../go-todotxt"
        "github.com/spf13/cobra"
        "os/user"
        "strings"
        "strconv"
        "github.com/rakyll/globalconf"
        "flag"
        "regexp"
        "os/exec"
        "io/ioutil"
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

type flagValue struct {
        str string
}

func (f *flagValue) String() string {
        return f.str
}

func (f *flagValue) Set(value string) error {
        f.str = value
        return nil
}

func newFlagValue(val string) *flagValue {
        return &flagValue{str: val}
}


func main() {

        conf, _ := globalconf.New("gotodo")

        var numtasks bool
        var sortby string
        var finished bool
        var prettyformat string
        var filename string
        var no_color bool


        var flagFilename = flag.String("file", "", "Location of the todo.txt file.")

        var flags = make(map[string]string)

        var cmdConfig = &cobra.Command{
            Use:   "config [key] [value]",
            Short: "Show and sets config values",
            Long:  `Config can be used to see and also set configuration variables.`,
            Run: func(cmd *cobra.Command, args []string) {
                    if len(args) == 0{
                            fmt.Printf("Available config variables:\n\n")
                            for k := range flags {
                                    fmt.Printf("%s\n", k)
                            }
                    }
                    if len(args) == 1 {
                            val, exists := flags[args[0]]
                            if exists {
                                    fmt.Printf("%s\n", val)
                            } else {
                                    // otherwise exit with non-zero status
                                    os.Exit(1)
                            }
                    }
                    if len(args) == 2 {
                            _, exists := flags[args[0]]
                            if exists {
                                    f := &flag.Flag{Name: args[0], Value: newFlagValue(args[1])}
                                    conf.Set("", f)
                            } else {
                                    os.Exit(1)
                            }
                    }
            },
        }

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

                    if no_color {
                        prettyformat = "%i %p %t%r"
                    }

                    var filteredTasks todotxt.TaskList
                    for _, task := range tasks {
                        if (!task.Finished() && !finished) ||
                           (task.Finished() && finished) {
                           if (len(args) > 0) {
                                   match, _ := regexp.MatchString(args[0], task.Text())
                                   if match {
                                           filteredTasks = append(filteredTasks, task)
                                   }
                           } else {
                                   filteredTasks = append(filteredTasks, task)
                           }
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
                                   "Sort tasks by parameter (prio|date|len|prio-rev|date-rev|len-rev|id|rand)")
        cmdList.Flags().StringVarP(&prettyformat, "pretty", "", "%i%c %p %t%r",
                                   "Pretty print tasks")
        cmdList.Flags().BoolVarP(&no_color, "no-color", "c", false,
                                 "Do not use colored output when pretty-printing")

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

        var cmdArchive = &cobra.Command{
            Use:   "archive [taskid]",
            Short: "Archives task.",
            Long:  `Archives task.`,
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

                fmt.Printf("Archiving task %v\n", taskid)

                tasks.Save(filename)
            },
        }


        var setprio string
        var settodo string
        var cmdSet = &cobra.Command{
            Use:   "set [taskid]",
            Short: "Sets given task's parameters.",
            Long:  `Sets given task's parameters.`,
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

                if len(setprio) > 0 {
                        tasks[taskid].SetPriority(setprio[0])
                        tasks[taskid].RebuildRawTodo()
                }

                if len(settodo) > 0 {
                        tasks[taskid].SetTodo(settodo)
                        tasks[taskid].RebuildRawTodo()
                }

                tasks.Save(filename)
            },
        }

        var cmdEdit = &cobra.Command{
            Use:   "edit [taskid]",
            Short: "Edit given task",
            Long:  `Edit given task with your prefered editor.`,
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

                text := tasks[taskid].RawText()

                file, err := ioutil.TempFile(os.TempDir(), "gotodo")
                defer os.Remove(file.Name())

                if err != nil {
                        panic(err)
                }

                e := ioutil.WriteFile(file.Name(), []byte(text), 0644)
                if e != nil {
                        panic(e)
                }

                editor := os.Getenv("EDITOR")
                if len(editor) == 0 {
                        editor = "nano" //FIXME: saner default?
                }

                c := exec.Command(editor, file.Name())

                // nasty hack, see http://stackoverflow.com/a/12089980
                c.Stdin = os.Stdin
                c.Stdout = os.Stdout
                c.Stderr = os.Stderr

                er := c.Run()

                if er != nil {
                        fmt.Println(er.Error())
                        panic(er)
                }

                dat, err := ioutil.ReadFile(file.Name())
                if err != nil {
                        panic(err)
                }

                lines := strings.Split(string(dat), "\n")
                if len (lines) > 0 {
                        if tasks[taskid].RawText() != lines[0] {
                                tasks[taskid] = todotxt.ParseTask(lines[0], taskid)
                                fmt.Printf("Task %d updated to:\n%s\n", taskid, lines[0])
                        }
                }

                tasks.Save(filename)
            },
        }


        cmdSet.PersistentFlags().StringVarP(&setprio, "priority", "p", "",
                                     "Sets task's priority.")
        cmdSet.PersistentFlags().StringVarP(&settodo, "todo", "t", "",
                                     "Sets task's todo.")

        var GotodoCmd = &cobra.Command{
            Use:   "gotodo",
            Short: "Gotodo is a go implementation of todo.txt.",
            Long: `A small, fast and fun implementation of todo.txt`,
            Run: func(cmd *cobra.Command, args []string) {
                cmdList.Run(cmd, nil)
            },
        }

        GotodoCmd.PersistentFlags().StringVarP(&filename, "filename", "", "",
                                     "Load tasks from this file.")

        conf.ParseAll()

        // values here
        flags["file"] = *flagFilename

        // sadly, this is the best we can do right now
        if filename == "" {
                if *flagFilename == "" {
                        filename = "todo.txt"
                } else {
                        filename = *flagFilename
                }
        }

        GotodoCmd.AddCommand(cmdList)
        GotodoCmd.AddCommand(cmdAdd)
        GotodoCmd.AddCommand(cmdDone)
        GotodoCmd.AddCommand(cmdArchive)
        GotodoCmd.AddCommand(cmdSet)
        GotodoCmd.AddCommand(cmdEdit)
        GotodoCmd.AddCommand(cmdConfig)
        GotodoCmd.Execute()
}
