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
    "github.com/logrusorgru/aurora"
)

type colorFunc func(interface{}) aurora.Value

var (
    SHELL_NAME = "ShootMan"
    SHELL_VERSION = "v1.0"
    INNER_COMMANDS = map[string]int {"session": 1, "help": 2, "clean": 3, "leave": 4}

    session = map[string]string{}
    current_session string

    fontPath string
)

func init () {
    flag.StringVar(&fontPath, "font-path", "", "set fonts path")
    flag.Parse()
}

func main() {
    if fontPath != "" {
        f, _ := figletlib.GetFontByName(fontPath, "big")
        PrintMsgWithColor(SHELL_NAME, f, 80, f.Settings(), "left", aurora.Cyan)
        fmt.Printf("%50s", aurora.Cyan(SHELL_VERSION + "\n\n"))
    } else {
        fmt.Println(aurora.Cyan(SHELL_NAME + " " + SHELL_VERSION))
        fmt.Println()
    }

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
    if 0 == sessionCommandLen {
        fmt.Println(aurora.Blue("[INFO] session list {key, value}: "))
        for k, v := range session {
            fmt.Println(aurora.Blue("{" + k + ", " + v + "}"))
        }
    } else if 1 == sessionCommandLen {
        if value, ok := session[s.arguments[0]]; ok {
            current_session = value
            fmt.Println(aurora.Green("[OK] load session: " + current_session))
        } else {
            fmt.Println(aurora.Red("[ERROR] load session error, please check"))
        }
    } else if 2 == sessionCommandLen {
        session[s.arguments[0]] = s.arguments[1]
        fmt.Println(aurora.Green("[OK] set session: [key] -- " + s.arguments[0] + " [value] -- "+ s.arguments[1] +" done!"))
        fmt.Println(aurora.Blue("[INFO] please remember to load session"))
    } else {
        fmt.Println(aurora.Red("[ERROR] wrong session command, use `help`"))
    }
}

func (s *ShellCommand) runClearCommand() {
    current_session = ""
    for k := range session {
        delete(session, k)
    }
    fmt.Println(aurora.Green("[OK] clean sessions done!"))
}

func (s *ShellCommand) runHelpCommand(){
    fmt.Println(aurora.Brown(`
 __  __          ___
/\ \/\ \        /\_ \
\ \ \_\ \     __\//\ \    _____
 \ \  _  \  / __ \\ \ \  /\  __ \
  \ \ \ \ \/\  __/ \_\ \_\ \ \L\ \
   \ \_\ \_\ \____\/\____\\ \ ,__/
    \/_/\/_/\/____/\/____/ \ \ \/
                            \ \_\
                             \/_/
session                  "list all sessions"
session [key] [user@ip]  "set session"
session [key]            "load session"
clean                    "clean sessions"
help                     "show help info"
leave                    "leave the command shell"

"when set session you can use remote linux shell"
`))
}

func (s *ShellCommand) runLeaveCommand() {
    fmt.Println(aurora.Blue("[INFO] leaving " + SHELL_NAME + "..."))
    os.Exit(1)
}

func (s *ShellCommand) runRemoteShellCommand() {
    if current_session == "" {
        fmt.Println(aurora.Red("[ERROR] please set/load session first ! use `help`"))
        return
    }

    args := append([]string{s.verb}, s.arguments...)

    remoteCmd := exec.Command("ssh", current_session, strings.Join(args, " "))
    remoteCmd.Stdout = os.Stdout
    remoteCmd.Stderr = os.Stderr

    if err:= remoteCmd.Run(); err != nil {
        fmt.Println(aurora.Red("[ERROR] run remote command error !"))
    }
}

func PrintMsgWithColor(msg string, f *figletlib.Font, maxwidth int, s figletlib.Settings, align string, color colorFunc) {
    lines := figletlib.GetLines(msg, f, maxwidth, s)
    PrintLineWithColor(lines, s.HardBlank(), maxwidth, align, color)
}

func PrintLineWithColor(lines []figletlib.FigText, hardblank rune, maxwidth int, align string, color colorFunc) {
    padleft := func(linelen int) {
        switch align {
        case "right":
            fmt.Print(strings.Repeat(" ", maxwidth-linelen))
        case "center":
            fmt.Print(strings.Repeat(" ", (maxwidth-linelen)/2))
        }
    }

    for _, line := range lines {
        for _, subline := range line.Art() {
            padleft(len(subline))
            for _, outchar := range subline {
                if outchar == hardblank {
                    outchar = ' '
                }
                fmt.Printf("%c", color(outchar))
            }
            if len(subline) < maxwidth && align != "right" {
                fmt.Println()
            }
        }
    }
}