package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

type processManager struct {
	conf  *config
	oscmd *exec.Cmd
}

func (pm *processManager) formatBuildTime(duration time.Duration) string {
	return fmt.Sprintf("%.2f(s)", duration.Seconds())
}

func (pm *processManager) run() {
	logger.Debug("building application...")
	// CPrint("BUILD", aurora.BlueFg, "building application...")

	// start := time.Now()

	os.Remove(pm.conf.build)
	out, err := exec.Command("go", "build", "-o", pm.conf.build).CombinedOutput()

	if err != nil {
		logger.Errorf("build failed: %s", err.Error())
		fmt.Printf("%s", out)
		return
	}

	// build success, display build time
	// logger.Infof("build took %s", pm.formatBuildTime(time.Since(start)))

	if pm.conf.Test {
		testOut, testErr := exec.Command("go", "test", "./...").CombinedOutput()

		if testErr != nil {
			CPrint("TESTS", aurora.RedFg, aurora.Red("Failed!").String())

			separator := aurora.Bold(aurora.Red("==========================="))
			fmt.Println(separator)

			scanner := bufio.NewScanner(strings.NewReader(string(testOut)))
			for scanner.Scan() {
				line := scanner.Text()
				if !strings.HasPrefix(line, "?") && !strings.HasPrefix(line, "ok") {
					fmt.Println(line)
				}
			}

			fmt.Println(separator)
		} else {
			CPrint("TESTS", aurora.GreenFg, aurora.Green("Pass!").String())
		}
	}

	pm.oscmd = exec.Command(pm.conf.build, pm.conf.Args...)
	pm.oscmd.Stdout = os.Stdout
	pm.oscmd.Stdin = os.Stdin
	pm.oscmd.Stderr = os.Stderr

	logger.Debugf("starting application with arguments: %v", pm.conf.Args)
	err = pm.oscmd.Start()
	if err != nil {
		logger.Errorf("error while starting application! %s", err.Error())
	}
}

func (pm *processManager) stop() {
	logger.Debug("stopping application")

	if pm.oscmd == nil {
		return
	}

	err := pm.oscmd.Process.Kill()
	if err != nil {
		logger.Errorf("error while stopping application! %s", err.Error())
	}
}
