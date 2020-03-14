package file

import (
	"os"
	"time"

	"github.com/pkg/errors"
)

func SafeStat(filename string) (os.FileInfo, bool, error) {
	s, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return s, false, nil
	} else if err != nil {
		return s, true, err
	}
	return s, true, nil
}

func Exists(filename string) (bool, error) {
	_, exists, err := SafeStat(filename)
	if err != nil {
		return false, errors.Wrapf(err, "error stating file %v", filename)
	}
	return exists, nil
}

func ModTime(filename string) (*time.Time, error) {
	s, exists, err := SafeStat(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "error stating %v", filename)
	} else if !exists {
		return nil, nil
	}

	t := s.ModTime()
	return &t, nil
}
