package options

import (
    "fmt"
    "os"
    "os/exec"

    "github.com/chigopher/pathlib"
    log "github.com/sirupsen/logrus"

    "carrierpigeondev/machv/src/lib"
)

func OptionCreateNewStaticQCOW2(
    staticDir   *pathlib.Path,
    isoTomlPath *pathlib.Path,
    isoDir      *pathlib.Path,
) (error) {
    allIsos, err := lib.ParseIsoToml(isoTomlPath)
    if err != nil {
        return fmt.Errorf("parsing %v to entries: %w", isoTomlPath, err)
    }

    lib.DisplayOptions(allIsos, "ISOs:")
    chosenIso, err := lib.SelectOption(allIsos)
    if err != nil {
        return fmt.Errorf("selecting option: %w", err)
    }

    isoPath := lib.GetIsoPath(chosenIso, isoDir)
    doesExist, err := isoPath.Exists()
    if err != nil {
        return fmt.Errorf("checking path %v: %w", isoPath, err)
    }
    if !doesExist {
        err = lib.DownloadIso(chosenIso, isoDir)
        if err != nil {
            return fmt.Errorf("downloading iso: %w", err)
        }
    }

    diskName, err := lib.GetInput("Enter name for new static disk: ")
    if err != nil {
        return fmt.Errorf("getting disk name: %w", err)
    }

    wantedDiskPath := lib.CreateDiskPath(diskName, staticDir)

    staticDiskPath, err := lib.CreateStaticQCOW2(wantedDiskPath, 20480)
    if err != nil {
        return fmt.Errorf("creating static qcwo2: %w", err)
    }

    if err := _createVM(staticDiskPath, isoPath); err != nil {
        return fmt.Errorf("creating vm: %w", err)
    }

    return nil
}

func _createVM(uniqueDiskPath *pathlib.Path, isoPath *pathlib.Path) (error) {
    log.WithField("isoPath", isoPath).Info("::LOOK::")
    qemuCreateCmd := fmt.Sprintf(
        "qemu-system-x86_64 %v -hda %v -boot d -cdrom %v",
        lib.QemuGlobalArgs, uniqueDiskPath, isoPath)
    log.WithField("isoPath", qemuCreateCmd).Info("::LOOK::")
    cmd := exec.Command("bash", "-c", qemuCreateCmd)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("running qemu-system-x86_64: %w", err)
    }
    
    return nil
}
