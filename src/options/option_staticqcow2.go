package options

import (
	"bufio"
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/chigopher/pathlib"
	"strings"
	"fmt"
	"os/exec"

	"carrierpigeondev/machv/src/lib"
)

func OptionCreateNewStaticQCOW2(staticDir *pathlib.Path, isoDir *pathlib.Path) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter name:\n(.qcow2 extension is added automatically)\n\n$ ")
	diskName, err := reader.ReadString('\n')
	if err != nil {
		log.WithError(err).WithField("diskName", diskName).Fatal("A fatal error has occurred while reading diskName from stdin")
	}
	diskName = strings.TrimSpace(diskName)

	diskPath := _createStaticQCOW2(staticDir, diskName, 20480) // does not return error as errors are handled inside the function
	isoPath := _selectISO(isoDir) // does not return error as errors are handled inside the function
	_createVM(diskPath, isoPath)
}

func _createStaticQCOW2(staticDir *pathlib.Path, name string, sizeMiB int) (*pathlib.Path) {
	// use a uuid as the base of a qcow2 file to ensure it's unique and create a path
	uniqueFileName := fmt.Sprintf("%v.qcow2", name)
	uniqueDiskPath := staticDir.Join(uniqueFileName)
	log.WithField("uniqueDiskPath", uniqueDiskPath).Info("Creating disk at the unique disk path and with specified size")

	// create and log the created command
	qemuImgCmdString := fmt.Sprintf("qemu-img create -f qcow2 %v %vM", uniqueDiskPath, sizeMiB)
	log.WithField("command", qemuImgCmdString).Info("Assembled command to create disk with qemu-img")

	// run the command
	cmd := exec.Command("bash", "-c", qemuImgCmdString)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.WithError(err).Fatal("A fatal error has occurred while running qemu-img")
	}

	return uniqueDiskPath
}

func _selectISO(isoDir *pathlib.Path) (*pathlib.Path) {
	// TEMP: no interaction, hardcoding it for now
	// eventually I want to make a dict that is hosted with isos and associated urls
	isoName := "debian-12.11.0-amd64-netinst.iso"
	url := "https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-12.11.0-amd64-netinst.iso"

	isoPath := isoDir.Join(isoName)
	doesExist, err := isoPath.Exists()
	if err != nil {
		log.WithError(err).WithField("path", isoPath).Fatal("A fatal error has occurred while checking path")
	} else if !doesExist {
		// download the iso since it does not exist
		curlCmdString := fmt.Sprintf("curl -LO %v", url)

		log.WithField("command", curlCmdString).Info("Downloading requested iso from url to isoDir using constructed command")

		cmd := exec.Command("bash", "-c", curlCmdString)
		cmd.Dir = isoDir.String()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			log.WithError(err).Fatal("A fatal error has occurred while downloading iso")
		}
	} else {
		log.Info("Requested iso already exists")
	}

	return isoPath
}

func _createVM(uniqueDiskPath *pathlib.Path, isoPath *pathlib.Path) {
	qemuCreateCmd := fmt.Sprintf("qemu-system-x86_64 %v -hda %v -boot d -cdrom %v", lib.QemuGlobalArgs, uniqueDiskPath, isoPath)
	cmd := exec.Command("bash", "-c", qemuCreateCmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.WithError(err).Fatal("A fatal error has occurred while running qemu-system-x86_64")
	}
}