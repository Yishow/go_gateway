package diagnostics

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func (o *ownedFiles) append(data []byte) error {
	if len(data) > MaxEventBytes+1 {
		return ErrFileUnsafe
	}
	if o.needsRollback {
		if err := o.files[0].Truncate(o.rollback); err != nil {
			return err
		}
		o.size = o.rollback
		o.needsRollback = false
	}
	for _, file := range o.files {
		if err := o.repairHeader(file); err != nil {
			return err
		}
	}
	info, err := o.files[0].Stat()
	if err != nil {
		return err
	}
	o.size = info.Size()
	if o.size+int64(len(data)) > MaxFileBytes {
		if err := copyOwned(o.files[2], o.files[1]); err != nil {
			return err
		}
		if err := copyOwned(o.files[1], o.files[0]); err != nil {
			return err
		}
		if err := o.files[0].Truncate(0); err != nil {
			return err
		}
		n, err := o.files[0].WriteAt(o.header, 0)
		if err != nil {
			return err
		}
		if n != len(o.header) {
			return io.ErrShortWrite
		}
		o.size = int64(n)
	}
	offset := o.size
	n, err := o.files[0].WriteAt(data, offset)
	if err != nil || n != len(data) {
		o.rollback = offset
		o.needsRollback = true
		if err == nil {
			err = io.ErrShortWrite
		}
		rollbackErr := o.files[0].Truncate(offset)
		if rollbackErr == nil {
			o.needsRollback = false
		}
		return errors.Join(err, rollbackErr)
	}
	o.size += int64(n)
	return o.files[0].Sync()
}
func copyOwned(dst, src *os.File) error {
	info, err := src.Stat()
	if err != nil || info.Size() > MaxFileBytes {
		return ErrFileUnsafe
	}
	// Rotation never renames or deletes a pathname; only validated held inodes change.
	if err := dst.Truncate(0); err != nil {
		return err
	}
	if _, err := dst.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if _, err := io.Copy(dst, io.NewSectionReader(src, 0, info.Size())); err != nil {
		return err
	}
	return dst.Sync()
}

// A failed write can leave a partial header. Repair only an already-owned inode.
func (o *ownedFiles) repairHeader(file *os.File) error {
	header := make([]byte, len(o.header))
	if _, err := file.ReadAt(header, 0); err == nil && bytes.Equal(header, o.header) {
		return nil
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	if n, err := file.WriteAt(o.header, 0); err != nil {
		return err
	} else if n != len(o.header) {
		return io.ErrShortWrite
	}
	return file.Sync()
}
func (o *ownedFiles) verifiedPath() string {
	root, err := openSafeDirectory(filepath.Dir(o.path), false)
	if err != nil {
		return ""
	}
	defer root.Close()
	for i, name := range o.names {
		actual, err := root.Lstat(name)
		if err != nil || reparse(actual) || !actual.Mode().IsRegular() {
			return ""
		}
		held, err := o.files[i].Stat()
		if err != nil || !os.SameFile(actual, held) || !singleLink(o.files[i]) {
			return ""
		}
		header := make([]byte, len(o.header))
		if _, err := o.files[i].ReadAt(header, 0); err != nil || !bytes.Equal(header, o.header) {
			return ""
		}
	}
	return o.path
}
