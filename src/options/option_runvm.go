package options

import (
    "fmt"
    "os"
    "os/exec"
    "strings"

    "github.com/chigopher/pathlib"
    log "github.com/sirupsen/logrus"

    "carrierpigeondev/machv/src/lib"
)

func OptionLaunchVirtualMachineFromUsableQCOW2(disksDir *pathlib.Path, sharePath *pathlib.Path) (error) {
    log.Info("Loading all usable virtual machine disks...")

	selectedDiskPath, err := _selectUsableDiskPath(disksDir)
	if err != nil {
		return fmt.Errorf("selecting usable disk path: %w", err)
	}

    var extraArgs string

	extraArgs, err = _addExtraArgs(
		fmt.Sprintf("Share %v with virtual machine?", sharePath),
		extraArgs,
		fmt.Sprintf("-virtfs local,path=%v,mount_tag=host0,security_model=mapped,id=host0", sharePath))
	if err != nil {
		return fmt.Errorf("adding args share path: %w", err)
	}

	extraArgs, err = _addExtraArgs(
		"Run with spice?",
		extraArgs,
		"-spice port=3001,disable-ticketing=on -device virtio-serial-pci -chardev spicevmc,id=vdagent,name=vdagent -device virtserialport,chardev=vdagent,name=com.redhat.spice.0",
	)
	if err != nil {
		return fmt.Errorf("adding args spice: %w", err)
	}

    log.WithFields(
		log.Fields{ 
			"diskPath": selectedDiskPath,
			"extraArgs": extraArgs }).
		Info("Observe the diskPath and extraArgs")
	
    _launchVM(selectedDiskPath, extraArgs)

	return nil
}

func _selectUsableDiskPath(disksDir *pathlib.Path) (*pathlib.Path, error) {
	disks, err := lib.ReadFilesInDirectory(disksDir)
	if err != nil {
		return nil, fmt.Errorf("reading files in disksDir: %w", err)
	}

	lib.DisplayOptions(disks, "Usable disks:")

    selectedDiskPath, err := lib.SelectOption(disks)
	if err != nil {
		return nil, fmt.Errorf("selecting disk option: %w", err)
	}

	return selectedDiskPath, nil
}

func _addExtraArgs(promptText string, args string, extraArgs string) (string, error) {
	lib.DisplayBool(promptText)
	runWithSpice, err := lib.SelectBool()
	if err != nil {
		return "", fmt.Errorf("getting bool: %w", err)
	}
	if runWithSpice {
		return strings.Join([]string{ args, extraArgs }, " "), nil
	}
	return args, nil
}

func _launchVM(diskPath *pathlib.Path, extraArgs string) {
    qemuRunCmd := fmt.Sprintf("qemu-system-x86_64 -%v -hda %v -boot a %v", lib.QemuGlobalArgs, diskPath, extraArgs)
    log.WithField("qemuRunCmd", qemuRunCmd).Info("Running qemu with the following command")
    cmd := exec.Command("bash", "-c", qemuRunCmd)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    if err := cmd.Run(); err != nil {
        log.WithError(err).Fatal("A fatal error has occurred while running qemu-system-x86_64")
    }
}