/*
Copyright (c) 2023 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
		func(_ string, fi os.FileInfo) error {
			if fi.Name() == root {
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
