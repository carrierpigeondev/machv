package lib

import (
    "github.com/chigopher/pathlib"
    "strings"
    "fmt"
    log "github.com/sirupsen/logrus"
    "os"
    "os/exec"
)

func CreateDiskPath(fileName string, directory *pathlib.Path) (*pathlib.Path) {
    if !strings.HasSuffix(fileName, ".qcow2") {
        fileName = strings.Join([]string{ fileName, ".qcow2" }, "")
    }

    return directory.Join(fileName)
}

func CreateStaticQCOW2(
    diskPath *pathlib.Path,
    sizeMiB int,
) (*pathlib.Path, error) {
    qemuImgCmdString := fmt.Sprintf(
        "qemu-img create -f qcow2 %v %vM",
        diskPath,
        sizeMiB,
    )
    log.WithField("command", qemuImgCmdString).
        Info("Assembled command to create disk with qemu-img")

    cmd := exec.Command("bash", "-c", qemuImgCmdString)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("running qemu-img command: %w", err)
    }

    return diskPath, nil
}

func CreateUsableQCOW2(
    staticDiskPath *pathlib.Path,
    usableDiskPath *pathlib.Path,
) (error) {
    if _, err := staticDiskPath.Copy(usableDiskPath); err != nil {
        return fmt.Errorf("copying static to usable disk path: %w", err)
    }
    return nil
}