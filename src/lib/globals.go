package lib

var (
    QemuGlobalArgs string
)

func init() {
    QemuGlobalArgs = "-enable-kvm -m 4G -smp 2 -netdev user,id=net0,net=192.168.0.0/24,dhcpstart=192.168.0.9 -device virtio-net-pci,netdev=net0 -vga qxl -device ich9-intel-hda"
}