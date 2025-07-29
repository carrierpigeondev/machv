package lib

import (
    "github.com/chigopher/pathlib"
    "strings"
)

func CreateDiskPath(fileName string, directory *pathlib.Path) (*pathlib.Path) {
    if strings.HasSuffix(fileName, ".qcow2") {
        strings.Join([]string{ fileName, ".qcow2" }, "")
    }

    return directory.Join(fileName)
}