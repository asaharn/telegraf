//go:generate ../../../tools/readme_config_includer/generator
package azure_data_explorer

import (
	_ "embed"
	"errors"
	"strings"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	ADX "github.com/influxdata/telegraf/plugins/outputs/azure_data_explorer"
	ES "github.com/influxdata/telegraf/plugins/outputs/azure_event_stream"
)

type MicrosoftFabric struct {
	ConnectionString  string                 `toml:"connection_string"`
	Log               telegraf.Logger        `toml:"-"`
	ADXConf           *ADX.AzureDataExplorer `toml:"adx_conf"`
	ESConf            *ES.EventStream        `toml:"es_conf"` //TODO: implement Eventstream plugin
	FabricSinkService telegraf.Output
}

// Close implements telegraf.Output.
func (m *MicrosoftFabric) Close() error {
	return m.FabricSinkService.Close()
}

// Connect implements telegraf.Output.
func (m *MicrosoftFabric) Connect() error {
	return m.FabricSinkService.Connect()
}

// SampleConfig implements telegraf.Output.
func (m *MicrosoftFabric) SampleConfig() string {
	return m.FabricSinkService.SampleConfig()
}

// Write implements telegraf.Output.
func (m *MicrosoftFabric) Write(metrics []telegraf.Metric) error {
	return m.FabricSinkService.Write(metrics)
}

func (m *MicrosoftFabric) Init() error {
	ConnectionString := m.ConnectionString

	if ConnectionString == "" {
		return errors.New("endpoint must not be empty")
	}

	if strings.HasPrefix(ConnectionString, "Endpoint=sb") {
		m.Log.Info("Detected EventStream endpoint, using EventStream output plugin")
		m.ESConf.Init()
		m.FabricSinkService = m.ESConf
	} else {
		m.Log.Info("Detected Kusto endpoint, using Kusto output plugin")
		m.ADXConf.Endpoint = ConnectionString
		m.ADXConf.Init()
		m.FabricSinkService = m.ADXConf
	}

	return nil

}

func init() {

	outputs.Add("microsoft_fabric", func() telegraf.Output {
		return &MicrosoftFabric{}
	})
}
