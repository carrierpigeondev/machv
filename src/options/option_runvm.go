package options

import (
	log "github.com/sirupsen/logrus"
	"github.com/chigopher/pathlib"
	"fmt"
	"os"
	"os/exec"

	"carrierpigeondev/machv/src/lib"
)

func OptionLaunchVirtualMachineFromUsableQCOW2(disksDir *pathlib.Path) {
	log.Info("Loading all usable virtual machine disks...")
	diskPath := lib.SelectDisk(disksDir)
	_launchVM(diskPath)

}

func _launchVM(diskPath *pathlib.Path) {
	qemuRunCmd := fmt.Sprintf("qemu-system-x86_64 -%v -hda %v -boot a", lib.QemuGlobalArgs, diskPath)
	cmd := exec.Command("bash", "-c", qemuRunCmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.WithError(err).Fatal("A fatal error has occurred while running qemu-system-x86_64")
	}
}