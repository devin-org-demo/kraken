// Copyright (c) 2016-2019 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package diskspaceutil

import (
	"os"
	"path/filepath"
	"syscall"
)

const path = "/"

// Helper method to get disk util.
func DiskSpaceUtil() (int, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return 0, err
	}

	diskAll := fs.Blocks * uint64(fs.Bsize)
	diskFree := fs.Bfree * uint64(fs.Bsize)
	diskUsed := diskAll - diskFree
	return int(diskUsed * 100 / diskAll), nil
}

func KrakenDiskUsage(paths []string, totalSize uint64) (int, error) {
	var krakenUsed int64

	for _, path := range paths {
		size, err := calculateDirSize(path)
		if err != nil {
			continue
		}
		krakenUsed += size
	}

	if totalSize == 0 {
		fs := syscall.Statfs_t{}
		err := syscall.Statfs("/", &fs)
		if err != nil {
			return 0, err
		}
		totalSize = fs.Blocks * uint64(fs.Bsize)
	}

	if totalSize == 0 {
		return 0, nil
	}

	return int(uint64(krakenUsed) * 100 / totalSize), nil
}

func calculateDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}
