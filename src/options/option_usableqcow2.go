package options

import (
	log "github.com/sirupsen/logrus"
	"github.com/chigopher/pathlib"
	"bufio"
	"os"
	"fmt"
	"strings"
	"errors"

	"carrierpigeondev/machv/src/lib"
)

func OptionCreateNewUsableQCOW2(staticDir *pathlib.Path, disksDir *pathlib.Path) {
	reader := bufio.NewReader(os.Stdin)

	log.Info("Loading all static virtual machine disks...")

	diskPath := lib.SelectDisk(staticDir)

	fmt.Print("Enter name for new usable virtual machine disk:\n(Make sure to include .qcow2 at the end)\n\n$ ")
	diskName, err := reader.ReadString('\n')
	if err != nil {
		log.WithError(err).WithField("diskName", diskName).Fatal("A fatal error has occurred while reading diskName from stdin")
	}
	diskName = strings.TrimSpace(diskName)
	if !strings.HasSuffix(diskName, ".qcow2") {
		log.WithError(errors.New("name does not end with .qcow2")).WithField("diskName", diskName).Fatal("A fatal error has occurred while validating disName")
	}

	_createUsableQCOW2(disksDir, diskPath, diskName)
}

func _createUsableQCOW2(disksDir *pathlib.Path, staticDiskPath *pathlib.Path, diskName string) {
	bytesInStaticDiskPath, err := staticDiskPath.Size()
	if err != nil {
		log.WithError(err).Error("A non-fatal error has occurred while getting bytes in staticDiskPath")
		bytesInStaticDiskPath = 0
	}

	newDiskPath := disksDir.Join(diskName)
	log.WithFields(log.Fields{"staticDiskPath": staticDiskPath, "newDiskPath": newDiskPath}).Info("Copying staticDiskPath to newDiskPath")

	doesExist, err := newDiskPath.Exists()
	if err != nil {
		log.WithError(err).WithField("path", newDiskPath).Fatal("A fatal error has occurred while checking path")
	} else if doesExist {
		log.WithError(errors.New("file already exists")).Fatal("A fatal error has occurred while checling if the specified newDiskPath already exists")
	}

	bytesCopied, err := staticDiskPath.Copy(newDiskPath)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"staticDiskPath": staticDiskPath, "newDiskPath": newDiskPath}).Fatal("A fatal error has occurred while copying staticDiskPath to newDiskPath")
	}

	bytesInNewDiskPath, err := staticDiskPath.Size()
	if err != nil {
		log.WithError(err).Error("A non-fatal error has occurred while getting bytes in newDiskPath")
		bytesInNewDiskPath = 0
	}

	log.WithFields(log.Fields{"bytesInStaticDiskPath": bytesInStaticDiskPath, "bytesCopied": bytesCopied, "bytesInNewDiskPath": bytesInNewDiskPath}).Info("Copied staticDiskPath to newDiskPath and logged number of bytes in both locations and transfer")
}
