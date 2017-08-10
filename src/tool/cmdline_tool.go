package main

import (
    "bufio"
    "os"
    "fmt"
    "time"
    "strings"
    "os/exec"
    "github.com/lukesampson/figlet/figletlib"
    "flag"
)

var (
    SHELL_NAME = "ShootMan"
    INNER_COMMANDS = map[string]int {"session": 1, "help": 2, "clear": 3, "leave": 4}

    session = map[string]string{}
    current_session string

    fontPath string
)

func init () {
    flag.StringVar(&fontPath, "font-path", "", "set fonts path")
    flag.Parse()
}

func main() {
    f, _ := figletlib.GetFontByName(fontPath, "larry3d")
    figletlib.PrintMsg(SHELL_NAME, f, 80, f.Settings(), "left")

    fmt.Printf("%70s", "v1.0\n\n")

    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("[%s(session: %s): %s]--> ", SHELL_NAME, current_session ,string([]byte(time.Now().String())[:19]))
        str, _ := reader.ReadString('\n')
        if cm := newCommand(str); cm != nil {
            cm.run()
        }
    }
}

func newCommand (command string) *ShellCommand {
    commandSegments := strings.Fields(command)
    if len(commandSegments) == 0 {
        return nil
    } else {
        verb := commandSegments[0]
        return &ShellCommand{verb: verb, arguments: commandSegments[1:], typ: judgeType(verb)}
    }
}

func judgeType(verb string) int {
    if value, ok := INNER_COMMANDS[verb]; ok {
        return value
    }

    return 0
}


type ShellCommand struct {
    verb      string
    arguments []string
    typ       int
}

func (s *ShellCommand) run () {
    switch s.typ {
    case 1:
        s.runSessionCommand()
    case 2:
        s.runHelpCommand()
    case 3:
        s.runClearCommand()
    case 4:
        s.runLeaveCommand()
    case 0: fallthrough
    default:
        s.runRemoteShellCommand()
    }
}

func (s *ShellCommand) runSessionCommand() {
    sessionCommandLen := len(s.arguments)
    if 1 == sessionCommandLen {
        if value, ok := session[s.arguments[0]]; ok {
            current_session = value
            fmt.Println("[OK] load session: " + current_session)
        } else {
            fmt.Println("[ERROR] load session error, please check")
        }
    } else if 2 == sessionCommandLen {
        session[s.arguments[0]] = s.arguments[1]
        fmt.Println("[OK] set session: [key] -- " + s.arguments[0] + " [value] -- "+ s.arguments[1] +" done!")
        fmt.Println("[INFO] please remember to load session")
    } else {
        fmt.Println("[ERROR] wrong session command, use `help`")
    }
}

func (s *ShellCommand) runClearCommand() {
    current_session = ""
    for k := range session {
        delete(session, k)
    }
    fmt.Println("[OK] clear sessions done!")
}

func (s *ShellCommand) runHelpCommand(){
    fmt.Println(`
 __  __          ___
/\ \/\ \        /\_ \
\ \ \_\ \     __\//\ \    _____
 \ \  _  \  / __ \\ \ \  /\  __ \
  \ \ \ \ \/\  __/ \_\ \_\ \ \L\ \
   \ \_\ \_\ \____\/\____\\ \ ,__/
    \/_/\/_/\/____/\/____/ \ \ \/
                            \ \_\
                             \/_/
session [key] [user@ip]  "set session"
session [key]            "load session"
clear                    "clear sessions"
help                     "show help info"

"when set session you can use remote linux shell"
`)
}

func (s *ShellCommand) runLeaveCommand() {
    fmt.Println("[INFO] leaving " + SHELL_NAME + "...")
    os.Exit(1)
}

func (s *ShellCommand) runRemoteShellCommand() {
    if current_session == "" {
        fmt.Println("[ERROR] please set/load session first ! use `help`")
        return
    }

    args := append([]string{s.verb}, s.arguments...)

    remoteCmd := exec.Command("ssh ", args...)
    remoteCmd.Stdout = os.Stdout
    remoteCmd.Stderr = os.Stderr

    if err:= remoteCmd.Run(); err != nil {
        fmt.Println("[ERROR] run remote command error !")
    }
}