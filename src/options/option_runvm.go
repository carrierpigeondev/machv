package options

import (
	log "github.com/sirupsen/logrus"
	"github.com/chigopher/pathlib"
	"fmt"
	"os"
	"os/exec"

	"carrierpigeondev/machv/src/lib"
)

func OptionLaunchVirtualMachineFromUsableQCOW2(disksDir *pathlib.Path, sharePath *pathlib.Path) {
	log.Info("Loading all usable virtual machine disks...")
	diskPath := lib.SelectDisk(disksDir)


	runWithDir := lib.SelectOption([]string { "No", "Yes" }, fmt.Sprintf("Share %v with virtual machine?", sharePath))
	if runWithDir == "Yes" {
		_launchVM(diskPath, fmt.Sprintf("-virtfs local,path=%v,mount_tag=host0,security_model=mapped,id=host0", sharePath))
	} else {
		_launchVM(diskPath, "")
	}
}

func _launchVM(diskPath *pathlib.Path, extraArgs string) {
	qemuRunCmd := fmt.Sprintf("qemu-system-x86_64 -%v -hda %v -boot a %v", lib.QemuGlobalArgs, diskPath, extraArgs)
	cmd := exec.Command("bash", "-c", qemuRunCmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.WithError(err).Fatal("A fatal error has occurred while running qemu-system-x86_64")
	}
}