package check

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/willie68/gowillie68/pkg/fileutils"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/osml"
)

func (c *Checker) Touch(sdCardFolder string, files []string) (*model.GeneralResult, error) {
	fs, err := os.Stat(sdCardFolder)
	if err != nil {
		return nil, err
	}

	if fs.IsDir() && (len(files) == 0) {
		files, err = osml.GetDataFiles(sdCardFolder)
		if err != nil {
			return nil, err
		}
	}

	if !fs.IsDir() {
		files = append(files, filepath.Base(sdCardFolder))
		sdCardFolder = filepath.Dir(sdCardFolder)
	}

	gr := model.NewGeneralResult()
	for _, file := range files {
		off := filepath.Join(sdCardFolder, strings.TrimSpace(file))
		if !fileutils.FileExists(off) {
			return nil, fmt.Errorf("file %s not exists", off)
		}
		ts, err := c.getFirstTimestamp(off)
		if err != nil {
			gr.Result = false
			gr.Messages = append(gr.Messages, fmt.Sprintf("error getting timestamp from file %s: %s\r\n", off, err.Error()))
			continue
		}
		err = os.Chtimes(off, time.Unix(0, 0), ts)
		if err != nil {
			gr.Result = false
			gr.Messages = append(gr.Messages, fmt.Sprintf("error touching file %s: %s\r\n", off, err.Error()))
			continue
		}
		gr.Messages = append(gr.Messages, fmt.Sprintf("touched file %s to %s\r\n", off, ts.String()))
	}
	gr.Result = true
	return gr, nil
}

func (c *Checker) getFirstTimestamp(file string) (time.Time, error) {
	f, err := os.Open(file)
	if err != nil {
		return time.Unix(0, 0), err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// Loop through the file and read each line
	for scanner.Scan() {
		line := scanner.Text() // Get the line as a string
		ll, ok, _ := model.ParseLogLine(line)
		if ok {
			ts, ok := c.getRMCTime(ll, time.Time{})
			if ok {
				return ts, nil
			}
		}
	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		c.log.Fatalf("error reading file: %v", err)
		return time.Unix(0, 0), err
	}

	return time.Unix(0, 0), errors.New("no valid timestamp found in file")
}
