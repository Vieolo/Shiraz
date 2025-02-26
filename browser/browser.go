// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package browser provides utilities for interacting with users' browsers.
package browser

import (
	"os"
	"os/exec"
	"time"
	"github.com/skratchdot/open-golang/open"
)


func checkEnvVarOverride() string {
	return os.Getenv("BROWSER")
}

// Open tries to open url in a browser and reports whether it succeeded.
func Open(url string) bool {
	var err error;
	if browserVar := checkEnvVarOverride(); browserVar != "" {
		err = open.RunWith(url, browserVar)
	} else {
		err = open.Run(url)
	}
	if err != nil {
		return false
	}
	return true
}

// appearsSuccessful reports whether the command appears to have run successfully.
// If the command runs longer than the timeout, it's deemed successful.
// If the command runs within the timeout, it's deemed successful if it exited cleanly.
func appearsSuccessful(cmd *exec.Cmd, timeout time.Duration) bool {
	errc := make(chan error, 1)
	go func() {
		errc <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		return true
	case err := <-errc:
		return err == nil
	}
}
