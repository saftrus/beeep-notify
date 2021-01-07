// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package beeep

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"syscall"
	"time"

	toast "github.com/go-toast/toast"
	"github.com/tadvi/systray"
	"golang.org/x/sys/windows/registry"
)

var isWindows10 bool
var applicationID string
var drivePath string

func init() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer k.Close()

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return
	}

	isWindows10 = maj == 10

	if isWindows10 {
		applicationID = appID()
		drivePath = appPath()
	}
}

func conv(actionTemp [][]string) []toast.Action {
	var action []toast.Action
	for _, h := range actionTemp {
		var a toast.Action
		a.Type = h[0]
		a.Label = h[1]
		a.Arguments = h[2]
		action = append(action, a)
	}
	return action
}

// Notify sends desktop notification.
func Notify(title, message, appIcon string, actionTemp [][]string) error {
	for i := range actionTemp {
		for j := range actionTemp[i] {
			if actionTemp[i][j] == "getPath" {
				actionTemp[i][j] = drivePath
			}
		}
	}
	if isWindows10 {
		return toastNotify(title, message, appIcon, actionTemp)
	}

	err := baloonNotify(title, message, appIcon, false)
	if err != nil {
		e := msgNotify(title, message)
		if e != nil {
			return errors.New("beeep: " + err.Error() + "; " + e.Error())
		}
	}

	return nil

}

func msgNotify(title, message string) error {
	msg, err := exec.LookPath("msg")
	if err != nil {
		return err
	}
	cmd := exec.Command(msg, "*", "/TIME:3", title+"\n\n"+message)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func baloonNotify(title, message, appIcon string, bigIcon bool) error {
	tray, err := systray.New()
	if err != nil {
		return err
	}

	err = tray.ShowCustom(pathAbs(appIcon), title)
	if err != nil {
		return err
	}

	go func() {
		tray.Run()
		time.Sleep(3 * time.Second)
		tray.Stop()
	}()

	return tray.ShowMessage(title, message, bigIcon)
}

func toastNotify(title, message, appIcon string, action [][]string) error {
	notification := toastNotification(title, message, pathAbs(appIcon), action)
	return notification.Push()
}

func toastNotification(title, message, appIcon string, actionTemp [][]string) toast.Notification {
	action := conv(actionTemp)
	return toast.Notification{
		AppID:   applicationID,
		Title:   title,
		Message: message,
		Icon:    appIcon,
		Actions: action,
	}
}

func appID() string {
	// defID := "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe"
	defID := "{7C5A40EF-A0FB-4BFC-874A-C0F2E0B9FA8E}\\CyDrive\\SystemTrayApp.exe"
	cmd := exec.Command("powershell", "Get-StartApps")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return defID
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "powershell.exe") {
			sp := strings.Split(line, " ")
			if len(sp) > 0 {
				return sp[len(sp)-1]
			}
		}
	}

	return defID
}

func appPath() string {

	cmd := exec.Command("powershell", "[System.IO.DriveInfo]::getdrives()[-1].Name")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return string(out)
}
