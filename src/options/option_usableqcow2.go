package options

import (
	"fmt"
	"github.com/chigopher/pathlib"
	log "github.com/sirupsen/logrus"

	"carrierpigeondev/machv/src/lib"
)

func OptionCreateNewUsableQCOW2(
    staticDir *pathlib.Path,
    disksDir  *pathlib.Path,
) error {
    log.Info("Loading all static virtual machine disks...")

    disks, err := lib.ReadFilesInDirectory(staticDir)
    if err != nil {
        return fmt.Errorf("reading disks from dir: %w", err)
    }

    lib.DisplayOptions(disks, "Static disks:")
    staticDiskPath, err := lib.SelectOption(disks)
    if err != nil {
        return fmt.Errorf("selecting disk path: %w", err)
    }

    diskName, err := lib.GetInput("Enter name for new usable machine disk: ")
    if err != nil {
        return fmt.Errorf("getting usable disk name: %w", err)
    }

    usableDiskPath := lib.CreateDiskPath(diskName, disksDir)

    doesExist, err := usableDiskPath.Exists()
    if err != nil {
        return fmt.Errorf("checking is disk exists: %w", err)
    } else if doesExist {
        return fmt.Errorf("usable disk already exists")
    }

    err = lib.CreateUsableQCOW2(staticDiskPath, usableDiskPath)
    if err != nil {
        return fmt.Errorf("creating usable qcow2: %w", err)
    }

    return nil
}
