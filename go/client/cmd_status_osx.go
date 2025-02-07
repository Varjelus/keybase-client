// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

// +build darwin

package client

import "github.com/keybase/client/go/install"

func (c *CmdStatus) osSpecific(status *fstatus) error {
	serviceStatus := install.KeybaseServiceStatus(c.G(), "service")
	kbfsStatus := install.KeybaseServiceStatus(c.G(), "kbfs")

	if len(serviceStatus.Pid) > 0 {
		status.Service.Running = true
		status.Service.Pid = serviceStatus.Pid
	}

	if len(kbfsStatus.Pid) > 0 {
		status.KBFS.Running = true
		status.KBFS.Pid = kbfsStatus.Pid
	}

	return nil
}

func (c *CmdStatus) serviceLogFilename() string {
	return "keybase.service.log"
}

func (c *CmdStatus) kbfsLogFilename() string {
	return "keybase.kbfs.log"
}

func (c *CmdStatus) desktopLogFilename() string {
	return "Keybase.app.log"
}
