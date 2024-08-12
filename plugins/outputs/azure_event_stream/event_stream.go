//go:generate ../../../tools/readme_config_includer/generator
package event_stream

import (
	_ "embed"
	"fmt"

	"github.com/influxdata/telegraf"
)

//go:embed sample.conf
var sampleConfig string

type EventStream struct {
	Endpoint      string `toml:"endpoint_url"`
	ClientName    string `toml:"client_name"`
	EventHubName  string `toml:"event_hub_name"`
	MessageFormat string `toml:"message_format"`
}

// Close implements telegraf.Output.
func (e *EventStream) Close() error {
	panic("ES CLOSE Unimplemented")
}

// Connect implements telegraf.Output.
func (e *EventStream) Connect() error {
	fmt.Println("YES CONNECTED")
	return nil
}

// SampleConfig implements telegraf.Output.
func (e *EventStream) SampleConfig() string {
	return sampleConfig
}

// Write implements telegraf.Output.
func (e *EventStream) Write(metrics []telegraf.Metric) error {
	fmt.Println("WRITING CALLED")
	return nil
}

func (e *EventStream) Init() error {
	fmt.Println("EVENT STREAM INIT")
	return nil
}
