package track

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/willie68/gowillie68/pkg/extstrgutils"
	"github.com/willie68/osmltools/internal/model"
)

func (m *manager) routeProps2track(track *model.Track, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("error opening route properties file %s: %v", f.Name, err)
	}
	defer rc.Close()
	props, err := m.readProps(rc, err)
	if err != nil {
		return err
	}

	m.props2Track(track, props)
	return nil
}

func (m *manager) readProps(rc io.ReadCloser, err error) (map[string]string, error) {
	props := make(map[string]string)
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value, err = strconv.Unquote("\"" + strings.TrimSpace(line[equal+1:]) + "\"")
					if err != nil {
						return nil, err
					}
				}
				props[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return props, nil
}

func (m *manager) props2Track(track *model.Track, props map[string]string) {
	track.Description = props["comment"]
	track.Name = props["name"]
	track.Files = make([]model.SourceData, 0)
	files := extstrgutils.SplitMultiValueParam(props["dataFiles"])
	md5s := extstrgutils.SplitMultiValueParam(props["dataMD5"])
	for x, f := range files {
		sd := model.SourceData{
			FileName: f,
			Hash:     fmt.Sprintf("md5:%s", md5s[x]),
		}
		track.Files = append(track.Files, sd)
	}
	track.MapFile = props["mapFile"]
}
