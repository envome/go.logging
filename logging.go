package logging

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"bright.sh/daemon/config"
	"github.com/fatih/color"
)

var Defaults *Options

var (
	// BGreen output
	BGreen = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	// BYellow output
	BYellow = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	// BRed output
	BRed = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	// BBlue output
	BBlue = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	// BMagenta output
	BMagenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	// BCyan output
	BCyan = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	// BReset output color
	BReset = string([]byte{27, 91, 48, 109})

	cYellow  = color.New(color.FgYellow).SprintFunc()
	cMagenta = color.New(color.FgMagenta).SprintFunc()
	cGreen   = color.New(color.FgGreen).SprintFunc()
	cRed     = color.New(color.FgRed).SprintFunc()
	cBlue    = color.New(color.FgBlue).SprintFunc()
)

type Logger struct {
	log.Logger
	Options *Options
}

type Options struct {
	Color            string `toml:"color" json:"color"`
	ScopeName        string `toml:"scope_name" json:"scope_name"`
	Path             string `toml:"path" json:"path"`
	Colorful         bool   `toml:"colorful" json:"colorful"`
	OutputToTerminal bool   `toml:"output_to_terminal" json:"output_to_terminal"`
}

func setDefaults(opts *Options) {
	if opts == nil {
		opts = &Options{
			Color:            BRed,
			ScopeName:        "app",
			Colorful:         true,
			OutputToTerminal: true,
		}
	}

	Defaults = opts
}

// GetLogger returns a logger with the given name and color
func GetLogger(name string, opts *Options) *Logger {
	l := &Logger{}

	if opts == nil {
		opts = &Options{}
	}

	if Defaults == nil {
		setDefaults(nil)
	}

	l.Options = opts

	if opts.Color == "" {
		opts.Color = Defaults.Color
	}

	if opts.ScopeName == "" {
		opts.ScopeName = Defaults.ScopeName
	}

	if opts.OutputToTerminal == true || opts.Path != "" {
		if opts.Colorful {
			l.SetPrefix(fmt.Sprintf("[ %s %10s %s ] %s ", opts.Color, name, BReset, opts.ScopeName))
		} else {
			l.SetPrefix(fmt.Sprintf("[ %10s ] %s", name, opts.ScopeName))
		}

		// var mw *io.Writer
		var logFile *os.File
		if opts.Path != "" {
			var err error
			logFile, err = os.OpenFile(opts.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				panic(err)
			}
		}

		if opts.OutputToTerminal && logFile != nil {
			mw := io.MultiWriter(os.Stdout, logFile)
			l.SetOutput(mw)
		} else if opts.OutputToTerminal {
			l.SetOutput(os.Stdout)
		} else if logFile != nil {
			l.SetOutput(logFile)
		} else {
			l.SetOutput(ioutil.Discard)
		}

		l.SetFlags(log.Ldate | log.Ltime)
	} else {
		l.SetOutput(ioutil.Discard)
	}

	return l
}

func (l *Logger) Error(str string, v ...interface{}) {
	if !l.Options.Colorful {
		l.Output(2, "ERROR "+fmt.Sprintf(str, v...))
	} else {
		l.Output(2, Red("ERROR ")+fmt.Sprintf(str, v...))
	}
}

func (l *Logger) Success(str string, v ...interface{}) {
	if !l.Options.Colorful {
		l.Output(2, "SUCCESS "+fmt.Sprintf(str, v...))
	} else {
		l.Output(2, Green("SUCCESS ")+fmt.Sprintf(str, v...))
	}
}

func (l *Logger) Warn(str string, v ...interface{}) {
	if !l.Options.Colorful {
		l.Output(2, "WARN "+fmt.Sprintf(str, v...))
	} else {
		l.Output(2, Yellow("WARN ")+fmt.Sprintf(str, v...))
	}
}

func Sprintf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}

func Yellow(args interface{}) string {
	if !Defaults.Colorful {
		return fmt.Sprint(args)
	}

	return cYellow(args)
}

func Red(args interface{}) string {
	if !Defaults.Colorful {
		return fmt.Sprint(args)
	}

	return cRed(args)
}

func Green(args interface{}) string {
	if !Defaults.Colorful {
		return fmt.Sprint(args)
	}

	return cGreen(args)
}

func Blue(args interface{}) string {
	if !Defaults.Colorful {
		return fmt.Sprint(args)
	}

	return cBlue(args)
}

// Info log info message
func Info(name string, v ...interface{}) *Logger {
	logger := GetLogger(name, &Options{Color: BMagenta})
	logger.Println(v...)
	return logger
}

// Error log error message
func Error(name string, v ...interface{}) *Logger {
	logger := GetLogger(name, &Options{Color: BRed})
	logger.Println("[ERROR]", v)
	return logger
}

func Last(length int64) *[]string {
	conf := config.GetInstance().Log
	file, err := os.Open(conf.Path)
	if err != nil {
		Error("logging", BRed, err)
		return nil
	}
	defer file.Close()

	buf := make([]byte, length)
	stat, err := os.Stat(conf.Path)
	start := stat.Size() - length
	_, err = file.ReadAt(buf, start)
	if err != nil {
		Error("logging", BRed, err)
		return nil
	}

	parts := strings.Split(string(buf), "\n")
	// sort.Strings(parts)
	reverse(parts)

	return &parts
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
