//go:build linux && amd64

package intel_pmu

import (
	"os"

	"github.com/intel/iaevents"
	"github.com/stretchr/testify/mock"
)

// mockValuesReader is an autogenerated mock type for the valuesReader type
type mockValuesReader struct {
	mock.Mock
}

// readValue provides a mock function with given fields: event
func (_m *mockValuesReader) readValue(event *iaevents.ActiveEvent) (iaevents.CounterValue, error) {
	ret := _m.Called(event)

	var r0 iaevents.CounterValue
	if rf, ok := ret.Get(0).(func(*iaevents.ActiveEvent) iaevents.CounterValue); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Get(0).(iaevents.CounterValue)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*iaevents.ActiveEvent) error); ok {
		r1 = rf(event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockEntitiesValuesReader is an autogenerated mock type for the entitiesValuesReader type
type mockEntitiesValuesReader struct {
	mock.Mock
}

// readEntities provides a mock function with given fields: _a0, _a1
func (_m *mockEntitiesValuesReader) readEntities(_a0 []*coreEventEntity, _a1 []*uncoreEventEntity) ([]coreMetric, []uncoreMetric, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []coreMetric
	if rf, ok := ret.Get(0).(func([]*coreEventEntity, []*uncoreEventEntity) []coreMetric); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]coreMetric)
		}
	}

	var r1 []uncoreMetric
	if rf, ok := ret.Get(1).(func([]*coreEventEntity, []*uncoreEventEntity) []uncoreMetric); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]uncoreMetric)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func([]*coreEventEntity, []*uncoreEventEntity) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// mockEntitiesActivator is an autogenerated mock type for the entitiesActivator type
type mockEntitiesActivator struct {
	mock.Mock
}

// activateEntities provides a mock function with given fields: coreEntities, uncoreEntities
func (_m *mockEntitiesActivator) activateEntities(coreEntities []*coreEventEntity, uncoreEntities []*uncoreEventEntity) error {
	ret := _m.Called(coreEntities, uncoreEntities)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*coreEventEntity, []*uncoreEventEntity) error); ok {
		r0 = rf(coreEntities, uncoreEntities)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockEntitiesParser is an autogenerated mock type for the entitiesParser type
type mockEntitiesParser struct {
	mock.Mock
}

// parseEntities provides a mock function with given fields: coreEntities, uncoreEntities
func (_m *mockEntitiesParser) parseEntities(coreEntities []*coreEventEntity, uncoreEntities []*uncoreEventEntity) error {
	ret := _m.Called(coreEntities, uncoreEntities)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*coreEventEntity, []*uncoreEventEntity) error); ok {
		r0 = rf(coreEntities, uncoreEntities)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockEntitiesResolver is an autogenerated mock type for the entitiesResolver type
type mockEntitiesResolver struct {
	mock.Mock
}

// resolveEntities provides a mock function with given fields: coreEntities, uncoreEntities
func (_m *mockEntitiesResolver) resolveEntities(coreEntities []*coreEventEntity, uncoreEntities []*uncoreEventEntity) error {
	ret := _m.Called(coreEntities, uncoreEntities)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*coreEventEntity, []*uncoreEventEntity) error); ok {
		r0 = rf(coreEntities, uncoreEntities)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockEventsActivator is an autogenerated mock type for the eventsActivator type
type mockEventsActivator struct {
	mock.Mock
}

// activateEvent provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockEventsActivator) activateEvent(_a0 iaevents.Activator, _a1 iaevents.PlacementProvider, _a2 iaevents.Options) (*iaevents.ActiveEvent, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *iaevents.ActiveEvent
	if rf, ok := ret.Get(0).(func(iaevents.Activator, iaevents.PlacementProvider, iaevents.Options) *iaevents.ActiveEvent); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iaevents.ActiveEvent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(iaevents.Activator, iaevents.PlacementProvider, iaevents.Options) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// activateGroup provides a mock function with given fields: _a0, _a1
func (_m *mockEventsActivator) activateGroup(_a0 iaevents.PlacementProvider, _a1 []iaevents.CustomizableEvent) (*iaevents.ActiveEventGroup, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iaevents.ActiveEventGroup
	if rf, ok := ret.Get(0).(func(iaevents.PlacementProvider, []iaevents.CustomizableEvent) *iaevents.ActiveEventGroup); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iaevents.ActiveEventGroup)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(iaevents.PlacementProvider, []iaevents.CustomizableEvent) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// activateMulti provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockEventsActivator) activateMulti(
	_a0 iaevents.MultiActivator,
	_a1 []iaevents.PlacementProvider,
	_a2 iaevents.Options,
) (*iaevents.ActiveMultiEvent, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *iaevents.ActiveMultiEvent
	if rf, ok := ret.Get(0).(func(iaevents.MultiActivator, []iaevents.PlacementProvider, iaevents.Options) *iaevents.ActiveMultiEvent); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iaevents.ActiveMultiEvent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(iaevents.MultiActivator, []iaevents.PlacementProvider, iaevents.Options) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockFileInfoProvider is an autogenerated mock type for the fileInfoProvider type
type mockFileInfoProvider struct {
	mock.Mock
}

// fileLimit provides a mock function with given fields:
func (_m *mockFileInfoProvider) fileLimit() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// readFile provides a mock function with given fields: _a0
func (_m *mockFileInfoProvider) readFile(_a0 string) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// lstat provides a mock function with given fields: _a0
func (_m *mockFileInfoProvider) lstat(_a0 string) (os.FileInfo, error) {
	ret := _m.Called(_a0)

	var r0 os.FileInfo
	if rf, ok := ret.Get(0).(func(string) os.FileInfo); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(os.FileInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockPlacementMaker is an autogenerated mock type for the placementMaker type
type mockPlacementMaker struct {
	mock.Mock
}

// makeCorePlacements provides a mock function with given fields: cores, perfEvent
func (_m *mockPlacementMaker) makeCorePlacements(cores []int, factory iaevents.PlacementFactory) ([]iaevents.PlacementProvider, error) {
	ret := _m.Called(cores, factory)

	var r0 []iaevents.PlacementProvider
	if rf, ok := ret.Get(0).(func([]int, iaevents.PlacementFactory) []iaevents.PlacementProvider); ok {
		r0 = rf(cores, factory)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]iaevents.PlacementProvider)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]int, iaevents.PlacementFactory) error); ok {
		r1 = rf(cores, factory)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// makeUncorePlacements provides a mock function with given fields: factory, socket
func (_m *mockPlacementMaker) makeUncorePlacements(socket int, factory iaevents.PlacementFactory) ([]iaevents.PlacementProvider, error) {
	ret := _m.Called(factory, socket)

	var r0 []iaevents.PlacementProvider
	if rf, ok := ret.Get(0).(func(iaevents.PlacementFactory, int) []iaevents.PlacementProvider); ok {
		r0 = rf(factory, socket)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]iaevents.PlacementProvider)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(iaevents.PlacementFactory, int) error); ok {
		r1 = rf(factory, socket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockSysInfoProvider is an autogenerated mock type for the sysInfoProvider type
type mockSysInfoProvider struct {
	mock.Mock
}

// allCPUs provides a mock function with given fields:
func (_m *mockSysInfoProvider) allCPUs() ([]int, error) {
	ret := _m.Called()

	var r0 []int
	if rf, ok := ret.Get(0).(func() []int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// allSockets provides a mock function with given fields:
func (_m *mockSysInfoProvider) allSockets() ([]int, error) {
	ret := _m.Called()

	var r0 []int
	if rf, ok := ret.Get(0).(func() []int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTransformer is an autogenerated mock type for the Transformer type
type MockTransformer struct {
	mock.Mock
}

// Transform provides a mock function with given fields: reader, matcher
func (_m *MockTransformer) Transform(reader iaevents.Reader, matcher iaevents.Matcher) ([]*iaevents.PerfEvent, error) {
	ret := _m.Called(reader, matcher)

	var r0 []*iaevents.PerfEvent
	if rf, ok := ret.Get(0).(func(iaevents.Reader, iaevents.Matcher) []*iaevents.PerfEvent); ok {
		r0 = rf(reader, matcher)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*iaevents.PerfEvent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(iaevents.Reader, iaevents.Matcher) error); ok {
		r1 = rf(reader, matcher)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
