package luarocks

import (
	"io"
	"strings"
)

type (
	Writer struct {
		io.Writer
	}
)
func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
}

func (w *Writer) WriteRepositoryPackages(list RocksList) (err error) {
	if _, err = w.Write([]byte("commands = {}\n")); err != nil {
		return
	}

	if _, err = w.Write([]byte("modules = {}\n")); err != nil {
		return
	}

	return w.WriteRepository(func(b *Writer) (rErr error) {
		for _, rock := range list {
			rErr = b.WritePackage(rock.Name, func(b *Writer) (pErr error) {
				var version *Version
				for _, version = range rock.Versions {
					pErr = b.WriteVersion(version.Name, func(b *Writer) (vErr error) {
						var arch string
						for _, arch = range version.Arch {
							if vErr = b.WriteVersionArch(arch); vErr != nil {
								return
							}
						}

						return
					})

					if pErr != nil {
						return
					}
				}
				return
			})

			if rErr != nil {
				return
			}
		}

		return
	})
}

func (w *Writer) WriteRepository(cb func (b *Writer) error) (err error) {
	_, err = w.Write([]byte("repository = {"))
	if err != nil {
		return
	}

	_, err = w.Write([]byte("\n"))
	if err != nil {
		return
	}

	defer w.Write([]byte("}"))
	return cb(w)
}

func (w *Writer) WritePackage(name string, cb func (b *Writer) error) (err error) {
	_, err = w.Write([]byte("\t"))
	if err != nil {
		return err
	}

	if strings.ContainsRune(name, '-') || strings.ContainsRune(name, '.'){
		if _, err = w.WriteMultilineString(name); err != nil {
			return
		}
	} else if _, err = w.Write([]byte(name)); err != nil {
		return
	}

	if _, err = w.Write([]byte(" = {\n")); err != nil {
		return
	}

	defer w.Write([]byte("\t},\n"))
	return cb(w)
}

func (w *Writer) WriteVersion(version string, cb func (b *Writer) error) (err error){
	if _, err = w.Write([]byte("\t\t")); err != nil {
		return
	}

	if _, err = w.WriteMultilineString(version); err != nil {
		return err
	}

	if _, err = w.Write([]byte(" = {\n")); err != nil {
		return err
	}

	defer w.Write([]byte("\t\t},\n"))

	return cb(w)
}

func (w *Writer) WriteVersionArch(arch string) (err error){
	if _, err = w.Write([]byte("\t\t\t{arch = \"")); err != nil {
		return
	}

	if _, err = w.Write([]byte(arch)); err != nil {
		return
	}

	if _, err = w.Write([]byte("\"},\n")); err != nil {
		return
	}

	return
}

func (w *Writer) WriteMultilineString(s string) (n int, err error) {
	var (
		part string
		i int
	)
	for _, part = range []string{"['", s, "']"} {
		if i, err = w.Write([]byte(part)); err != nil {
			return
		}
		n += i
	}

	return
}

