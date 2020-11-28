package nfs

import (
	"bytes"
	"context"
	"log"
	"os"
	"reflect"
	"syscall"

	"github.com/go-git/go-billy/v5"
	"github.com/willscott/go-nfs-client/nfs/xdr"
)

func onSetAttr(ctx context.Context, w *response, userHandle Handler) error {
	w.errorFmt = wccDataErrorFormatter
	handle, err := xdr.ReadOpaque(w.req.Body)
	if err != nil {
		return &NFSStatusError{NFSStatusInval, err}
	}

	fs, path, err := userHandle.FromHandle(handle)
	if err != nil {
		return &NFSStatusError{NFSStatusStale, err}
	}
	attrs, err := ReadSetFileAttributes(w.req.Body)
	if err != nil {
		return &NFSStatusError{NFSStatusInval, err}
	}

	info, err := fs.Lstat(fs.Join(path...))
	if err != nil {
		if os.IsNotExist(err) {
			return &NFSStatusError{NFSStatusNoEnt, err}
		}
		return &NFSStatusError{NFSStatusAccess, err}
	}

	// see if there's a "guard"
	if guard, err := xdr.ReadUint32(w.req.Body); err != nil {
		return &NFSStatusError{NFSStatusInval, err}
	} else if guard != 0 {
		// read the ctime.
		t := FileTime{}
		if err := xdr.Read(w.req.Body, &t); err != nil {
			return &NFSStatusError{NFSStatusInval, err}
		}
		if info.Sys() != nil {
			extra := reflect.ValueOf(info.Sys())
			if extra.Kind() == reflect.Struct {
				ctimeField := extra.FieldByName("Ctimespec")
				if ts, ok := ctimeField.Interface().(syscall.Timespec); ok {
					if !t.EqualTimespec(ts.Unix()) {
						return &NFSStatusError{NFSStatusNotSync, nil}
					}
					goto TIME_GOOD
				} else {
					log.Printf("Ctimespec field isn't a timespec")
				}
			}
		}
		return &NFSStatusError{NFSStatusNotSupp, nil}
	TIME_GOOD:
	}

	if !billy.CapabilityCheck(fs, billy.WriteCapability) {
		return &NFSStatusError{NFSStatusROFS, os.ErrPermission}
	}

	changer := userHandle.Change(fs)
	if err := attrs.Apply(changer, fs, fs.Join(path...)); err != nil {
		// Already an nfsstatuserror
		return err
	}

	preAttr := ToFileAttribute(info).AsCache()

	writer := bytes.NewBuffer([]byte{})
	if err := xdr.Write(writer, uint32(NFSStatusOk)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	if err := WriteWcc(writer, preAttr, tryStat(fs, path)); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}

	if err := w.Write(writer.Bytes()); err != nil {
		return &NFSStatusError{NFSStatusServerFault, err}
	}
	return nil
}
