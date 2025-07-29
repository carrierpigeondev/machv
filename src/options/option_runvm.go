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

func OptionLaunchVirtualMachineFromUsableQCOW2(disksDir *pathlib.Path, sharePath *pathlib.Path) {
    log.Info("Loading all usable virtual machine disks...")
    diskPath := lib.SelectDisk(disksDir)

    var extraArgs string

    runWithDir := lib.SelectOption([]string { "No", "Yes" }, fmt.Sprintf("Share %v with virtual machine?", sharePath))
    if runWithDir == "Yes" {
        extraArgs = strings.Join([]string{ extraArgs, fmt.Sprintf("-virtfs local,path=%v,mount_tag=host0,security_model=mapped,id=host0", sharePath) }, " ")
    }

    runWithSpice := lib.SelectOption([]string { "No", "Yes" }, "Run with spice?")
    if runWithSpice == "Yes" {
        extraArgs = strings.Join([]string{ extraArgs, "-spice port=3001,disable-ticketing=on -device virtio-serial-pci -chardev spicevmc,id=vdagent,name=vdagent -device virtserialport,chardev=vdagent,name=com.redhat.spice.0" }, " ")
    }

    log.WithFields(log.Fields{ "diskPath": diskPath, "extraArgs": extraArgs }).Info("Observe the diskPath and extraArgs")
    _launchVM(diskPath, extraArgs)
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