package index

import (
	"errors"
	"os"
	"syscall"
	"time"
)

type Entry struct {
	// file metadata changed
	ctime     time.Time
	ctimeNSec int64
	// file data changed
	mtime     time.Time
	mtimeNSec int64
	// ID of the hardware device the file is stored on
	dev   uint64
	inode uint64
	// file modes
	mode os.FileMode
	// user id
	uid uint32
	// group id
	gid      uint32
	fileSize uint64
}

func NewEntry(fInfo os.FileInfo) (*Entry, error) {
	// Get the underlying data source and type assert to syscall.Stat_t
	stat, ok := fInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, errors.New("not a syscall.Stat_t type")
	}

	// Get the inode number
	inode := stat.Ino

	mtime := time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
	mtimeNsec := stat.Mtim.Nsec

	// Get change time (ctime) and nanoseconds
	ctime := time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	ctimeNsec := stat.Ctim.Nsec

	return &Entry{
		ctime:     ctime,
		ctimeNSec: ctimeNsec,
		mtime:     mtime,
		mtimeNSec: mtimeNsec,
		dev:       stat.Dev,
		inode:     inode,
		mode:      fInfo.Mode(),
		uid:       stat.Uid,
		gid:       stat.Gid,
		fileSize:  uint64(stat.Size),
	}, nil
}
