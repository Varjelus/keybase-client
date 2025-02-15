// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

// +build dragonfly freebsd linux netbsd openbsd solaris

package client

func (c *CmdStatus) osSpecific(status *fstatus) error {
	// TODO: on darwin, install.KeybaseServiceStatus() is implemented to get pid for service and kbfs.
	// This is currently the best way to determine if KBFS is running, so other OS's should implement
	// it.
	return nil
}

func (c *CmdStatus) serviceLogFilename() string {
	return "keybase.service.log"
}

func (c *CmdStatus) kbfsLogFilename() string {
	return "keybase.kbfs.log"
}

func (c *CmdStatus) desktopLogFilename() string {
	return "keybase.gui.log"
}
