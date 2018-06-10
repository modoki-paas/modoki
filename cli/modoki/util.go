package main

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"

	modoki "github.com/cs3238-tsuzu/modoki/client"
	"golang.org/x/net/websocket"

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

// Get stdout and stderr logs from a container.
func modokiLogsContainer(c *modoki.Client, ctx context.Context, path string, id string, follow *bool, since *time.Time, stderr *bool, stdout *bool, tail *string, timestamps *bool, until *time.Time) (*websocket.Conn, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "ws"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("id", id)
	if follow != nil {
		tmp26 := strconv.FormatBool(*follow)
		values.Set("follow", tmp26)
	}
	if since != nil {
		tmp27 := since.Format(time.RFC3339)
		values.Set("since", tmp27)
	}
	if stderr != nil {
		tmp28 := strconv.FormatBool(*stderr)
		values.Set("stderr", tmp28)
	}
	if stdout != nil {
		tmp29 := strconv.FormatBool(*stdout)
		values.Set("stdout", tmp29)
	}
	if tail != nil {
		values.Set("tail", *tail)
	}
	if timestamps != nil {
		tmp30 := strconv.FormatBool(*timestamps)
		values.Set("timestamps", tmp30)
	}
	if until != nil {
		tmp31 := until.Format(time.RFC3339)
		values.Set("until", tmp31)
	}
	u.RawQuery = values.Encode()

	var header http.Header
	if c.JWTSigner != nil {
		req, _ := http.NewRequest("ws", "", nil)
		req.URL = &u
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}

		u = *req.URL
		header = req.Header
	}
	url_ := u.String()
	cfg, err := websocket.NewConfig(url_, url_)
	if err != nil {
		return nil, err
	}

	cfg.Header = header

	return websocket.DialConfig(cfg)
}
