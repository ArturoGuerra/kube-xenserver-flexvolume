package flexvolume
import (
    "syscall"
)

func (n *nodeClient) UnmountDevice(mountdir string) {
    if err := syscall.Unmount(mountdir, syscall.MNT_DETACH); err != nil {
        n.Reply(Failure(err.Error()))
    }

    n.Reply(Success("Unmounted device"))
}
