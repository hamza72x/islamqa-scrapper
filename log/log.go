package log

import (
	"fmt"
	"islam-qa-scrapper/helper"
	"os"
)

const (
	ERROR_EMOJI    = "⛔️"
	FATAL_EMOJI    = "🔥"
	WARN_EMOJI     = "⚠️ "
	INFO_EMOJI     = "📌"
	OK_EMOJI       = "✅"
	DOWNLOAD_EMOJI = "⬇️"
	LOG_FILE       = "log.log"
)

var (
	logFile *os.File
)

func Initialize() {

	helper.RemoveFileIfExists(LOG_FILE)

	var err error
	logFile, err = os.OpenFile(LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
}

func lines(str string) {
	fmt.Println("-------------------------------------------------------" + str)
}

func Info(a ...interface{}) {
	print(INFO_EMOJI, a...)
}

func Err(a ...interface{}) {
	print(ERROR_EMOJI, a...)
	logFile.WriteString(fmt.Sprintln(ERROR_EMOJI, a))
}

func Warn(a ...interface{}) {
	print(WARN_EMOJI, a...)
	logFile.WriteString(fmt.Sprintln(WARN_EMOJI, a))
}

func Ok(a ...interface{}) {
	print(OK_EMOJI, a...)
}

func Fatal(a ...interface{}) {
	print(FATAL_EMOJI, a...)
	logFile.WriteString(fmt.Sprintln(FATAL_EMOJI, a))
	panic("")
}

func print(emoji string, a ...interface{}) {
	fmt.Print("\n")
	lines("")
	fmt.Fprint(os.Stdout, emoji)
	fmt.Fprint(os.Stdout, " ")
	fmt.Fprintln(os.Stdout, a...)
	lines("")
}
