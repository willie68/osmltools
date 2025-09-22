package track

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/willie68/osmltools/internal/model"
)

// IsOldVersion checks if the given track file is in the old format (contains route.properties)
func (m *manager) IsOldVersion(tf string) bool {
	if _, err := os.Stat(tf); err != nil {
		m.log.Errorf("error zip file %s does not exists: %v", tf, err)
		return false
	}

	r, err := zip.OpenReader(tf)
	if err != nil {
		m.log.Errorf("error opening zip file %s: %v", tf, err)
		return false
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "route.properties" {
			return true
		}
	}
	return false
}

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
	files := SplitMultiValueParam(props["dataFiles"])
	md5s := SplitMultiValueParam(props["dataMD5"])
	for x, f := range files {
		sd := model.SourceData{
			FileName: f,
			Hash:     fmt.Sprintf("md5:%s", md5s[x]),
		}
		track.Files = append(track.Files, sd)
	}
	track.MapFile = props["mapFile"]
}
