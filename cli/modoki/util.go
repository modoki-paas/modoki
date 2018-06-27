package main

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"

	goaclient "github.com/goadesign/goa/client"
)

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

// newJWTSigner returns the request signer used for authenticating
// against the jwt security scheme.
func newJWTSigner(key string) goaclient.Signer {
	return &goaclient.APIKeySigner{
		SignQuery: false,
		KeyName:   "Authorization",
		KeyValue:  key,
		Format:    "Bearer %s",
	}

}

func stringPtr(s string) *string {
	return &s
}

func createTarArchive(src string) (string, error) {
	if _, err := os.Stat(src); err != nil {
		return "", err
	}

	fp, err := ioutil.TempFile("/tmp", "modoki_tar_")

	if err != nil {
		return "", err
	}
	defer fp.Close()

	tw := tar.NewWriter(fp)

	src, err = filepath.Abs(src)

	if err != nil {
		return "", err
	}

	src = filepath.Clean(src)

	err = filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		abs, err := filepath.Abs(file)

		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.TrimPrefix(filepath.Clean(abs), filepath.Dir(src)), string(filepath.Separator))

		if header.Name == "" {
			return nil
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(file)
		defer f.Close()
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return fp.Name(), nil
}

func extractTarArchive(reader io.Reader, target string, stat types.ContainerPathStat, verbose bool) error {
	target = filepath.Clean(target)

	tarReader := tar.NewReader(reader)
	if !stat.Mode.IsDir() {
		if s, err := os.Stat(target); err == nil && s.IsDir() {
			target = filepath.Join(target, stat.Name)
		}

		header, err := tarReader.Next()

		if verbose {
			log.Println("Extracting: ", header.Name)
		}

		info := header.FileInfo()

		fp, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
		_, err = io.Copy(fp, tarReader)

		return err
	}

	if st, err := os.Stat(target); err != nil {
		if err := os.Mkdir(target, 0774); err != nil {
			return err
		}
	} else {
		if !st.IsDir() {
			return errors.New("The path is not a directory")
		}
		target = filepath.Join(target, stat.Name)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			fmt.Println("break")

			break
		} else if err != nil {
			return err
		}
		path := filepath.Join(target, strings.TrimPrefix(header.Name, stat.Name))

		if verbose {
			log.Println("Extracting: ", header.Name)
		}
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		fp, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(fp, tarReader)
		if err != nil {
			fp.Close()
			return err
		}
		fp.Close()

	}

	return nil
}
