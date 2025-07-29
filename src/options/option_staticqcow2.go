package options

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"

    "github.com/chigopher/pathlib"
    log "github.com/sirupsen/logrus"

    "carrierpigeondev/machv/src/lib"
)

func OptionCreateNewStaticQCOW2(staticDir *pathlib.Path, isoTomlPath *pathlib.Path, isoDir *pathlib.Path) (error) {
    isoPath, err := _selectISO(isoTomlPath, isoDir)
	if err != nil {
		return fmt.Errorf("selecting iso path: %w", err)
	}

    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter name:\n(.qcow2 extension is added automatically)\n\n$ ")
    diskName, err := reader.ReadString('\n')
    if err != nil {
        log.WithError(err).WithField("diskName", diskName).Fatal("A fatal error has occurred while reading diskName from stdin")
    }
    diskName = strings.TrimSpace(diskName)
    diskPath := _createStaticQCOW2(staticDir, diskName, 20480) // does not return error as errors are handled inside the function

    _createVM(diskPath, isoPath)

	return nil
}

func _createStaticQCOW2(staticDir *pathlib.Path, name string, sizeMiB int) *pathlib.Path {
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

func _selectISO(isoTomlPath *pathlib.Path, isoDir *pathlib.Path) (*pathlib.Path, error) {

    allIsos, err := lib.ParseIsoTomlToIsoEntrySlice(isoTomlPath)
    if err != nil {
		return nil, fmt.Errorf("parsing %v to entries: %w", isoTomlPath, err)
	}

	lib.DisplayOptions(allIsos, "ISOs:")
    chosenIso, err := lib.SelectOption(allIsos)
	if err != nil {
		return nil, fmt.Errorf("selecting option: %w", err)
	}

    url := chosenIso.Url
    isoName := chosenIso.FriendlyName()

    isoPath := isoDir.Join(isoName)
    doesExist, err := isoPath.Exists()
    if err != nil {
		return nil, fmt.Errorf("checking path %v: %w", isoPath, err)
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
			return nil, fmt.Errorf("downloading iso: %w", err)
        }
    } else {
        log.Info("Requested iso already exists")
    }

    return isoPath, nil
}

func _createVM(uniqueDiskPath *pathlib.Path, isoPath *pathlib.Path) {
    log.WithField("isoPath", isoPath).Info("::LOOK::")
    qemuCreateCmd := fmt.Sprintf("qemu-system-x86_64 %v -hda %v -boot d -cdrom %v", lib.QemuGlobalArgs, uniqueDiskPath, isoPath)
    log.WithField("isoPath", qemuCreateCmd).Info("::LOOK::")
    cmd := exec.Command("bash", "-c", qemuCreateCmd)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    if err := cmd.Run(); err != nil {
        log.WithError(err).Fatal("A fatal error has occurred while running qemu-system-x86_64")
    }
}
