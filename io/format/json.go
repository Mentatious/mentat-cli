package format

import (
	"encoding/json"
	"github.com/Mentatious/mentat-cli/io"
	"github.com/Mentatious/mentat-cli/schema"
	"io/ioutil"
	"os"
)

const (
	dumpFilePermissions = 0644 // // TODO: think of moving to "config"
)

// DumpJSON ... dump entries as JSON
func DumpJSON(entries []schema.Entry, sinkName string) error {
	log := io.GetLog()
	b, err := json.MarshalIndent(entries, "", "\t")
	b = append(b, '\n')
	if err != nil {
		log.Infof("error dumping as JSON: %s", err.Error())
		return err
	}
	if sinkName == "stdout" {
		os.Stdout.Write(b)
	} else {
		err := ioutil.WriteFile(sinkName, b, dumpFilePermissions)
		if err != nil {
			log.Infof("error writing dumped data to '%s': %s", sinkName, err.Error())
			return err
		}
	}
	return nil
}
