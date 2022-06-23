package update

import (
	"bytes"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DoSelfUpdate(currentVersion string) bool {
	v := semver.MustParse(currentVersion)
	latest, err := selfupdate.UpdateSelf(v, "annikovk/IntelliJ-Log-Analyzer")
	if err != nil {
		log.Println("Binary update failed:", err)
		return false
	}
	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		log.Println("Current binary is the latest version", currentVersion)
		return true
	} else {
		log.Println("Successfully updated to version", latest.Version)
		log.Println("Release note:\n", latest.ReleaseNotes)
		return true
	}
}

func DoSelfUpdateMac() bool {
	latest, found, _ := selfupdate.DetectLatest("annikovk/IntelliJ-Log-Analyzer")
	if found {
		homeDir, _ := os.UserHomeDir()
		downloadPath := filepath.Join(homeDir, "Downloads", "IntelliJ-Log-Analyzer.zip")
		err := exec.Command("curl", "-L", latest.AssetURL, "-o", downloadPath).Run()
		if err != nil {
			log.Println("curl error:", err)
			return false
		}
		var appPath string
		cmdPath, err := os.Executable()
		appPath = strings.TrimSuffix(cmdPath, "IntelliJ Log Analyzer.app/Contents/MacOS/IntelliJ Log Analyzer")
		if err != nil {
			appPath = "/Applications/"
		}
		cmd := exec.Command("ditto", "-xk", downloadPath, appPath)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			log.Printf("ditto error: %s\n ouput: %s, error: %s", err, out.String(), stderr.String())
			log.Printf("downloadPath: %s \n appPath %s", downloadPath, appPath)
			return false
		}
		err = exec.Command("rm", downloadPath).Run()
		if err != nil {
			log.Println("removing error:", err)
			log.Printf("downloadPath: %s", downloadPath)
			return false
		}
		return true
	} else {
		return false
	}
}

func CheckForUpdate(currentVersion string) (bool, string, string) {
	latest, found, err := selfupdate.DetectLatest("annikovk/IntelliJ-Log-Analyzer")
	if err != nil {
		log.Println("Error occurred while detecting version:", err)
		return false, "", ""
	}

	v := semver.MustParse(currentVersion)
	if !found || latest.Version.LTE(v) {
		log.Println("Current version is the latest")
		return false, "", ""
	}

	return true, latest.Version.String(), latest.ReleaseNotes
}
