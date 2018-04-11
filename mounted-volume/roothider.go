package mountedvolume

import (
	"log"
	"syscall"
)

// HideRoot hides the root folder by performing a mount of a tmpfs on top of the /root folder.
func HideRoot() error {
	err := syscall.Mount("tmpfs", "/root", "tmpfs", syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, "size=1m")
	if err != nil {
		log.Printf("unable to hide /root: %s", err)
	}
	return err
}

// UnhideRoot unhides the root folder by performing a unmount of the tmpfs that is on top of the /root folder.
func UnhideRoot() error {
	err := syscall.Unmount("/root", 0)
	if err != nil {
		log.Printf("unable to unhide /root: %s", err)
	}
	return err
}
