package lib

import (
    "bufio"
    "os"

    "github.com/chigopher/pathlib"
    log "github.com/sirupsen/logrus"

    "fmt"
    "strings"
    "strconv"
    "errors"
)

func ReadFilesInDirectory(dir *pathlib.Path) ([]*pathlib.Path) {
    disks, err := dir.ReadDir()
    if err != nil {
        log.WithError(err).WithField("dir", dir).Fatal("A fatal error has occurred while reading directory for files")
    }

    return disks
}

func SelectOption[T any](options []T, promptText string) (T) {
    reader := bufio.NewReader(os.Stdin)

    fmt.Println(promptText)

    for i, option := range options {
        fmt.Printf("  %v) %v\n", i, option)
    }
    fmt.Print("\n$ ")
    
    optionIndexString, err := reader.ReadString('\n')
    if err != nil {
        log.WithError(err).WithField("optionIndexString", optionIndexString).Fatal("A fatal error has occurred while reading option")
    }
    optionIndexString = strings.TrimSpace(optionIndexString)

    optionIndex, err := strconv.Atoi(optionIndexString)
    if err != nil {
        log.WithError(err).WithField("optionIndexString", optionIndexString).Fatal("A fatal error has occurred while converting optionIndexString to integer optionIndex")
    }
    if optionIndex < 0 || optionIndex >= len(options) {
        log.WithError(errors.New("index out of range")).Fatal("A fatal error has occurred while validating optionIndex")
    }

    return options[optionIndex]
}

func SelectDisk(dir *pathlib.Path) (*pathlib.Path) {
    return SelectFile(dir, "Disks:")
}

func SelectFile(dir *pathlib.Path, promptText string) (*pathlib.Path) {
    disks := ReadFilesInDirectory(dir)
    chosenFile := SelectOption(disks, promptText)

    log.WithField("chosenFile", chosenFile).Info("Chose file")

    return chosenFile
}