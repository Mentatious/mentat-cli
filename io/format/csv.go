package format

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/Mentatious/mentat-cli/io"
	"github.com/Mentatious/mentat-cli/schema"
)

var (
	// TODO: config with export formats options (e.g. exported fields, whatever) (not only for CSV)
	header  = []string{"Content", "Tags", "Type", "From", "Description", "WasAdded"} // think of moving this to "config"
	tagsSep = ":"
)

// DumpCSV ... dump entries as CSV
func DumpCSV(entries []schema.Entry, sinkName string) error {
	log := io.GetLog()

	var writer *csv.Writer
	if sinkName == "stdout" {
		writer = csv.NewWriter(os.Stdout)
	} else {
		sink, err := os.Create(sinkName)
		defer sink.Close()
		if err != nil {
			log.Infof("error dumping data to '%s': %s", sinkName, err)
		}
		writer = csv.NewWriter(sink)
	}
	defer writer.Flush()
	err := writer.Write(header) // write header
	if err != nil {
		log.Infof("error dumping data to '%s': %s", sinkName, err)
	}
	for _, entry := range entries {
		err = writer.Write([]string{
			entry.Content,
			strings.Join(entry.Tags, tagsSep),
			entry.Type,
			entry.Metadata.From,
			entry.Metadata.Description,
			entry.Metadata.TimeAddedOrigin,
		})
		if err != nil {
			log.Infof("error dumping data to '%s': %s", sinkName, err)
		}

	}
	return nil
}
