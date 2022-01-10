package update

import (
	"errors"
	"fmt"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
	. "internal/logging"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type SelfUpdater struct {
	updater      *updater.Updater
	wg           sync.WaitGroup
	updateStatus updater.UpdateStatus
	updateErr    error
}

func NewSelfUpdater(version string) *SelfUpdater {
	return &SelfUpdater{
		updater: &updater.Updater{
			Provider: &provider.Github{
				RepositoryURL: "github.com/MrL0co/devtools",
				ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
			},
			ExecutableName: fmt.Sprintf("devtools_%s_%s", runtime.GOOS, runtime.GOARCH),
			Version:        version,
		},
	}
}

func (selfUpdater *SelfUpdater) verifyInstallation() error {
	latestVersion, err := selfUpdater.updater.GetLatestVersion()
	if err != nil {
		return err
	}
	executable, err := selfUpdater.updater.GetExecutable()
	if err != nil {
		return err
	}
	cmd := exec.Cmd{
		Path: executable,
		Args: []string{executable, "--version"},
	}
	// Should be replaced with Output() as soon as test project is updated
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	strOutput := string(output)

	if !strings.Contains(strOutput, latestVersion) {
		return errors.New("version not found in program output")
	}
	return nil
}

func (selfUpdater *SelfUpdater) selfUpdate() error {
	var err error

	selfUpdater.updateStatus, err = selfUpdater.updater.Update()
	if err != nil {
		return err
	}

	if selfUpdater.updateStatus == updater.Updated {
		if err := selfUpdater.verifyInstallation(); err != nil {
			Log.Error("Update failed", err)
			Log.Error("Rolling back...")

			return selfUpdater.updater.Rollback()
		}

		Log.Info("Updated to latest version!")
	}
	return nil
}

func (selfUpdater *SelfUpdater) StartUpdateCheck() {
	if selfUpdater.GetVersion() == "development" {
		return
	}

	Log.Debug("Current version: " + selfUpdater.GetVersion())
	Log.Debug("Looking for updates...")

	selfUpdater.wg.Add(1)

	go func() {
		if err := selfUpdater.selfUpdate(); err != nil {
			log.Println(err)
		}
		selfUpdater.wg.Done()
	}()
}

func (selfUpdater *SelfUpdater) Wait() {
	selfUpdater.wg.Wait() // Waiting for the update process to finish before exiting
}

func (selfUpdater *SelfUpdater) GetVersion() string {
	return selfUpdater.updater.Version
}
