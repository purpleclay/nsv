package nsv

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/saracen/walker"
)

var (
	vprefixed = map[string]struct{}{
		"go.mod": {},
	}

	// a custom error that is never propagated, but used to circuit break the file walker
	errLanguageDetected = errors.New("language detected")
)

func detectLanguageIsPrefixed(dir string) bool {
	root := dir
	if root == "" {
		root = "."
	}

	err := walker.Walk(root,
		func(pathname string, fi os.FileInfo) error {
			if pathname == root {
				return nil
			}

			if fi.IsDir() {
				return filepath.SkipDir
			}

			if _, matched := vprefixed[fi.Name()]; matched {
				return errLanguageDetected
			}
			return nil
		})
	return errors.Is(err, errLanguageDetected)
}
