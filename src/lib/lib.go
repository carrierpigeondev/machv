package lib

import (
    "bufio"
    "os"
    
    "github.com/chigopher/pathlib"
    
    "fmt"
    "strings"
    "strconv"
)

func GetInput(promptText string) (string, error) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print(promptText)
    userInput, err := reader.ReadString('\n')
    if err != nil {
        return "", fmt.Errorf("reading input %w", err)
    }
    userInput = strings.TrimSpace(userInput)

    return userInput, nil
}

func ReadFilesInDirectory(dir *pathlib.Path) ([]*pathlib.Path, error) {
    disks, err := dir.ReadDir()
    if err != nil {
        return nil, fmt.Errorf("reading %q: %w", dir, err)
    }
    return disks, nil
}

func DisplayOptions[T any](options []T, promptText string) {
    fmt.Println(promptText)
    for i, option := range options {
        fmt.Printf("  %v) %v\n", i, option)
    }
    fmt.Print("\n$ ")
}

func DisplayBool(promptText string) {
    DisplayOptions([]string { "No", "Yes" }, promptText)
}

func SelectBool() (bool, error) {
    choice, err := SelectOption([]string { "No", "Yes" })
    if err != nil {
        return false, fmt.Errorf("selecting bool: %w", err)
    }

    switch choice {
    case "No":
        return false, nil
    case "Yes":
        return true, nil
    }

    return false, fmt.Errorf("returned other than No, Yes: %v", choice)
}

func SelectOption[T any](options []T) (T, error) {
    reader := bufio.NewReader(os.Stdin)
    var zero T
    
    optionIndexString, err := reader.ReadString('\n')
    if err != nil {
        return zero, fmt.Errorf("reading option: %w", err)
    }
    optionIndexString = strings.TrimSpace(optionIndexString)

    optionIndex, err := strconv.Atoi(optionIndexString)
    if err != nil {
        return zero, fmt.Errorf("converting %v to int: %w", optionIndex, err)
    }
    if optionIndex < 0 || optionIndex >= len(options) {
        return zero, fmt.Errorf("index %v out of range", optionIndex)
    }

    return options[optionIndex], nil
}
