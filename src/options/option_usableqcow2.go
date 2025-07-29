package options

import (
	"errors"
	"fmt"
	"os"

	"github.com/chigopher/pathlib"
	log "github.com/sirupsen/logrus"

	"carrierpigeondev/machv/src/lib"
)

func OptionCreateNewUsableQCOW2(staticDir *pathlib.Path, disksDir *pathlib.Path) (error) {
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

    if err := _createUsableQCOW2(staticDiskPath, usableDiskPath); err != nil {
        return fmt.Errorf("creating usable qcow2: %w", err)
    }

    return nil
}

func _createUsableQCOW2(
    staticDiskPath *pathlib.Path,
    usableDiskPath *pathlib.Path,
) (error) {
    doesExist, err := usableDiskPath.Exists()
    if err != nil {
        return fmt.Errorf("checking is disk exists: %w", err)
    } else if doesExist {
        return fmt.Errorf("usable disk already exists")
    }

    if _, err := staticDiskPath.Copy(newDiskPath); err != nil {
        log.WithError(err).WithFields(log.Fields{"staticDiskPath": staticDiskPath, "newDiskPath": newDiskPath}).Fatal("A fatal error has occurred while copying staticDiskPath to newDiskPath")
    }

    bytesInNewDiskPath, err := staticDiskPath.Size()
    if err != nil {
        log.WithError(err).Error("A non-fatal error has occurred while getting bytes in newDiskPath")
        bytesInNewDiskPath = 0
    }

    log.WithFields(log.Fields{"bytesInStaticDiskPath": bytesInStaticDiskPath, "bytesCopied": bytesCopied, "bytesInNewDiskPath": bytesInNewDiskPath}).Info("Copied staticDiskPath to newDiskPath and logged number of bytes in both locations and transfer")
}
