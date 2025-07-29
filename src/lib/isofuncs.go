package lib

import (
    "github.com/chigopher/pathlib"
    "fmt"
    log "github.com/sirupsen/logrus"
    "os"
    "os/exec"
)

func GetIsoPath(iso IsoEntry, isoDir *pathlib.Path) (*pathlib.Path) {
    return isoDir.Join(iso.FriendlyName())
}

func DownloadIso(
    iso IsoEntry,
    isoDir *pathlib.Path,
) (error) {
    curlCmdString := fmt.Sprintf("curl -LO %v", iso.Url)

    log.WithField("command", curlCmdString).
        Info("Downloading requested iso from url to isoDir using constructed command")

    cmd := exec.Command("bash", "-c", curlCmdString)
    cmd.Dir = isoDir.String()
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("downloading iso: %w", err)
    }

    return nil
}