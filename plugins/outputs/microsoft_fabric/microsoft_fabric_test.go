package microsoft_fabric

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMicrosoftFabric_Init_EmptyConnectionString(t *testing.T) {
	mf := &MicrosoftFabric{
		ConnectionString: "",
		Log:              testutil.Logger{},
	}
	err := mf.Init()
	require.Error(t, err)
	assert.Equal(t, "endpoint must not be empty", err.Error())
}
