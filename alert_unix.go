// +build linux freebsd netbsd openbsd

package beeep

// Alert displays a desktop notification and plays a beep.
func Alert(title, message, appIcon string) error {
	action := [][]string
	if err := Notify(title, message, appIcon, action); err != nil {
		return err
	}
	return Beep(DefaultFreq, DefaultDuration)
}
