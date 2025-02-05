package adx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-kusto-go/kusto"
	kustoerrors "github.com/Azure/azure-kusto-go/kusto/data/errors"
	"github.com/Azure/azure-kusto-go/kusto/ingest"
	"github.com/Azure/azure-kusto-go/kusto/kql"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/internal/choice"
)

type AzureDataExplorer struct {
	Endpoint        string          `toml:"endpoint_url"`
	Database        string          `toml:"database"`
	Timeout         config.Duration `toml:"timeout"`
	MetricsGrouping string          `toml:"metrics_grouping_type"`
	TableName       string          `toml:"table_name"`
	CreateTables    bool            `toml:"create_tables"`
	IngestionType   string          `toml:"ingestion_type"`
	kustoClient     *kusto.Client
	metricIngestors map[string]ingest.Ingestor
	AppName         string
	logger          telegraf.Logger
}

const (
	TablePerMetric = "tablepermetric"
	SingleTable    = "singletable"
	// These control the amount of memory we use when ingesting blobs
	bufferSize = 1 << 20 // 1 MiB
	maxBuffers = 5
)

const ManagedIngestion = "managed"
const QueuedIngestion = "queued"

// Initialize the client and the ingestor
func (adx *AzureDataExplorer) Connect() error {
	conn := kusto.NewConnectionStringBuilder(adx.Endpoint).WithDefaultAzureCredential()
	// Since init is called before connect, we can set the connector details here including the type. This will be used for telemetry and tracing.
	conn.SetConnectorDetails("Telegraf", internal.ProductToken(), adx.AppName, "", false, "")
	client, err := kusto.New(conn)
	if err != nil {
		return err
	}
	adx.kustoClient = client
	adx.metricIngestors = make(map[string]ingest.Ingestor)

	return nil
}

// Clean up and close the ingestor
func (adx *AzureDataExplorer) Close() error {
	var errs []error
	for _, v := range adx.metricIngestors {
		if err := v.Close(); err != nil {
			// accumulate errors while closing ingestors
			errs = append(errs, err)
		}
	}
	if err := adx.kustoClient.Close(); err != nil {
		errs = append(errs, err)
	}

	adx.kustoClient = nil
	adx.metricIngestors = nil

	if len(errs) == 0 {
		adx.logger.Info("Closed ingestors and client")
		return nil
	}
	// Combine errors into a single object and return the combined error
	return kustoerrors.GetCombinedError(errs...)
}

func (adx *AzureDataExplorer) PushMetrics(format ingest.FileOption, tableName string, metricsArray []byte) error {
	var metricIngestor ingest.Ingestor
	var err error
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(adx.Timeout))
	defer cancel()
	metricIngestor, err = adx.GetMetricIngestor(ctx, tableName)
	if err != nil {
		return err
	}

	length := len(metricsArray)
	adx.logger.Debugf("Writing %d metrics to table %q", length, tableName)
	reader := bytes.NewReader(metricsArray)
	mapping := ingest.IngestionMappingRef(tableName+"_mapping", ingest.JSON)
	if metricIngestor != nil {
		if _, err := metricIngestor.FromReader(ctx, reader, format, mapping); err != nil {
			adx.logger.Errorf("sending ingestion request to Azure Data Explorer for table %q failed: %v", tableName, err)
		}
	}
	return nil
}

func (adx *AzureDataExplorer) GetMetricIngestor(ctx context.Context, tableName string) (ingest.Ingestor, error) {
	ingestor := adx.metricIngestors[tableName]

	if ingestor == nil {
		if err := adx.createAzureDataExplorerTable(ctx, tableName); err != nil {
			return nil, fmt.Errorf("creating table for %q failed: %w", tableName, err)
		}
		// create a new ingestor client for the table
		tempIngestor, err := createIngestorByTable(adx.kustoClient, adx.Database, tableName, adx.IngestionType)
		if err != nil {
			return nil, fmt.Errorf("creating ingestor for %q failed: %w", tableName, err)
		}
		adx.metricIngestors[tableName] = tempIngestor
		adx.logger.Debugf("Ingestor for table %s created", tableName)
		ingestor = tempIngestor
	}
	return ingestor, nil
}

func (adx *AzureDataExplorer) createAzureDataExplorerTable(ctx context.Context, tableName string) error {
	if !adx.CreateTables {
		adx.logger.Info("skipped table creation")
		return nil
	}

	if _, err := adx.kustoClient.Mgmt(ctx, adx.Database, createTableCommand(tableName)); err != nil {
		return err
	}

	if _, err := adx.kustoClient.Mgmt(ctx, adx.Database, createTableMappingCommand(tableName)); err != nil {
		return err
	}

	return nil
}

func (adx *AzureDataExplorer) Init() error {
	if adx.Endpoint == "" {
		return errors.New("endpoint configuration cannot be empty")
	}
	if adx.Database == "" {
		return errors.New("database configuration cannot be empty")
	}

	adx.MetricsGrouping = strings.ToLower(adx.MetricsGrouping)
	if adx.MetricsGrouping == SingleTable && adx.TableName == "" {
		return errors.New("table name cannot be empty for SingleTable metrics grouping type")
	}

	if adx.MetricsGrouping == "" {
		adx.MetricsGrouping = TablePerMetric
	}

	if !(adx.MetricsGrouping == SingleTable || adx.MetricsGrouping == TablePerMetric) {
		return errors.New("metrics grouping type is not valid")
	}

	if adx.Timeout == 0 {
		adx.Timeout = config.Duration(20 * time.Second)
	}

	if adx.IngestionType == "" {
		adx.IngestionType = QueuedIngestion
	} else if !(choice.Contains(adx.IngestionType, []string{ManagedIngestion, QueuedIngestion})) {
		return fmt.Errorf("unknown ingestion type %q", adx.IngestionType)
	}
	return nil
}

// For each table create the ingestor
func createIngestorByTable(client *kusto.Client, database, tableName, ingestionType string) (ingest.Ingestor, error) {
	switch strings.ToLower(ingestionType) {
	case ManagedIngestion:
		mi, err := ingest.NewManaged(client, database, tableName)
		return mi, err
	case QueuedIngestion:
		qi, err := ingest.New(client, database, tableName, ingest.WithStaticBuffer(bufferSize, maxBuffers))
		return qi, err
	}
	return nil, fmt.Errorf(`ingestion_type has to be one of %q or %q`, ManagedIngestion, QueuedIngestion)
}

func createTableCommand(table string) kusto.Statement {
	builder := kql.New(`.create-merge table ['`).AddTable(table).AddLiteral(`'] `)
	builder.AddLiteral(`(['fields']:dynamic, ['name']:string, ['tags']:dynamic, ['timestamp']:datetime);`)

	return builder
}

func createTableMappingCommand(table string) kusto.Statement {
	builder := kql.New(`.create-or-alter table ['`).AddTable(table).AddLiteral(`'] `)
	builder.AddLiteral(`ingestion json mapping '`).AddTable(table + "_mapping").AddLiteral(`' `)
	builder.AddLiteral(`'[{"column":"fields", `)
	builder.AddLiteral(`"Properties":{"Path":"$[\'fields\']"}},{"column":"name", `)
	builder.AddLiteral(`"Properties":{"Path":"$[\'name\']"}},{"column":"tags", `)
	builder.AddLiteral(`"Properties":{"Path":"$[\'tags\']"}},{"column":"timestamp", `)
	builder.AddLiteral(`"Properties":{"Path":"$[\'timestamp\']"}}]'`)

	return builder
}

// Setters for testing
func (adx *AzureDataExplorer) SetLogger(logger telegraf.Logger) {
	adx.logger = logger
}

func (adx *AzureDataExplorer) SetKustoClient(client *kusto.Client) {
	adx.kustoClient = client
}

func (adx *AzureDataExplorer) SetMetricsIngestors(ingestors map[string]ingest.Ingestor) {
	adx.metricIngestors = ingestors
}
