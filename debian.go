// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// build goes machine(s)
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"time"

	"github.com/platinasystems/ar"
)

func sliceToDeb(aw *ar.Writer, name string, b []byte) (err error) {
	header := &ar.Header{
		Name:    name,
		ModTime: time.Now(),
		Mode:    0600,
		Size:    int64(len(b)),
	}
	if err := aw.WriteHeader(header); err != nil {
		return fmt.Errorf("Error writing ar header for %s: %w",
			name, err)
	}
	n, err := aw.Write(b)
	if err != nil {
		return fmt.Errorf("Error writing ar data for %s: %w",
			name, err)
	}
	if n != len(b) {
		return fmt.Errorf("Error writing ar data for %s: wrote %d expecting %d",
			name, n, len(b))
	}
	return nil
}

type nb struct {
	Name string
	Body []byte
}

func newTarMember(aw *ar.Writer, name string, data []nb) (err error) {
	var databuf bytes.Buffer
	twz := gzip.NewWriter(&databuf)
	tw := tar.NewWriter(twz)

	for _, entry := range data {
		hdr := &tar.Header{
			Name: entry.Name,
			Mode: 0600,
			Size: int64(len(entry.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("Error writing header for %s in %s: %w",
				entry.Name, name, err)
		}
		cnt, err := tw.Write(entry.Body)
		if err != nil {
			return fmt.Errorf("Error writing data for %s in %s: %w",
				entry.Name, name, err)
		}
		if cnt != len(entry.Body) {
			return fmt.Errorf("Error writing data for %s in %s: wrote %d expecting %d",
				entry.Name, name, cnt, len(entry.Body))
		}
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("Error closing tar %s: %w", name, err)
	}
	if err := twz.Close(); err != nil {
		return fmt.Errorf("Error closing %s: %w", name, err)
	}
	if err := sliceToDeb(aw, name, databuf.Bytes()); err != nil {
		return err
	}
	return nil
}

func NewDebianArchive(name string, data []nb) (err error) {
	deb, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("Error creating Debian archive %s: %w",
			name, err)
	}
	aw := ar.NewWriter(deb)
	aw.WriteGlobalHeader()
	sliceToDeb(aw, "debian-binary", []byte("2.0\n"))

	controlFile := fmt.Sprintf(`Package: %s
Version: %s
Maintainer: %s
Description: Package short description will go here
 Long Description will go here.
 Long Description will continue here.
Architecture: amd64
`,
		"test-package",
		"1.0.0",
		"Platina Systems <platinasystems@platinasystems.com>",
	)

	var ctlbuf bytes.Buffer
	twz := gzip.NewWriter(&ctlbuf)
	tw := tar.NewWriter(twz)
	hdr := &tar.Header{
		Name: "control",
		Mode: 0600,
		Size: int64(len(controlFile)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("Error writing header for control in control.tar.gz: %w",
			err)
	}
	cnt, err := tw.Write([]byte(controlFile))
	if err != nil {
		return fmt.Errorf("Error writing data for control in control.tar: %w",
			err)
	}
	if cnt != len(controlFile) {
		return fmt.Errorf("Error writing data for control in control.tar: wrote %d expecting %d",
			cnt, len(controlFile))
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("Error closing control.tar: %w", err)
	}
	if err := twz.Close(); err != nil {
		return fmt.Errorf("Error closing control.tar.gz: %w", err)
	}
	if err := sliceToDeb(aw, "control.tar.gz", ctlbuf.Bytes()); err != nil {
		return err
	}

	newTarMember(aw, "data.tar.gz", data)

	if err := deb.Close(); err != nil {
		return fmt.Errorf("Error closing archive %s: %w", name, err)
	}
	return nil
}
