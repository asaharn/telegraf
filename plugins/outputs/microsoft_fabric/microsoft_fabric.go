//go:generate ../../../tools/readme_config_includer/generator
package microsoft_fabric

import (
	_ "embed"
	"errors"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/outputs"
	ADX "github.com/influxdata/telegraf/plugins/outputs/azure_data_explorer"
	EH "github.com/influxdata/telegraf/plugins/outputs/event_hubs"
)

//go:embed sample.conf
var sampleConfig string

type MicrosoftFabric struct {
	ConnectionString  string                 `toml:"connection_string"`
	Log               telegraf.Logger        `toml:"-"`
	ADXConf           *ADX.AzureDataExplorer `toml:"adx_conf"`
	EHConf            *EH.EventHubs          `toml:"eh_conf"`
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
	return sampleConfig
}

// Write implements telegraf.Output.
func (m *MicrosoftFabric) Write(metrics []telegraf.Metric) error {
	return m.FabricSinkService.Write(metrics)
}

func (m *MicrosoftFabric) Init() error {
	ConnectionString := m.ConnectionString

	if ConnectionString == "" {
		return errors.New("endpoint must not be empty. For Kusto refer : https://learn.microsoft.com/kusto/api/connection-strings/kusto?view=microsoft-fabric for EventHouse refer : https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/add-manage-eventstream-sources?pivots=enhanced-capabilities")
	}

	if strings.HasPrefix(ConnectionString, "Endpoint=sb") {
		m.Log.Info("Detected EventHouse endpoint, using EventHouse output plugin")
		m.EHConf.ConnectionString = ConnectionString
		m.EHConf.Log = m.Log
		m.EHConf.Init()
		m.FabricSinkService = m.EHConf
	} else if isKustoEndpoint(strings.ToLower(ConnectionString)) {
		m.Log.Info("Detected Kusto endpoint, using Kusto output plugin")
		//Setting up the AzureDataExplorer plugin initial properties
		m.ADXConf.Endpoint = ConnectionString
		m.ADXConf.Log = m.Log
		m.ADXConf.Init()
		m.FabricSinkService = m.ADXConf
	} else {
		return errors.New("invalid connection string. For Kusto refer : https://learn.microsoft.com/kusto/api/connection-strings/kusto?view=microsoft-fabric for EventHouse refer : https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/add-manage-eventstream-sources?pivots=enhanced-capabilities")
	}
	return nil
}

func isKustoEndpoint(endpoint string) bool {
	prefixes := []string{
		"https://",
		"data source=",
		"addr=",
		"address=",
		"network address=",
		"server=",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(endpoint, prefix) {
			return true
		}
	}
	return false
}

func init() {

	outputs.Add("microsoft_fabric", func() telegraf.Output {
		return &MicrosoftFabric{
			ADXConf: &ADX.AzureDataExplorer{
				Timeout:      config.Duration(20 * time.Second),
				CreateTables: true,
			},
			EHConf: &EH.EventHubs{
				Timeout: config.Duration(30 * time.Second),
			},
		}
	})
}
