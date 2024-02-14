// Code generated by MockGen. DO NOT EDIT.
// Source: internal/registry/registry.go
//
// Generated by this command:
//
//	mockgen -source=internal/registry/registry.go -destination=internal/registry/mocks/mock_registry.go
//

// Package mock_registry is a generated GoMock package.
package mock_registry

import (
	goods "LamodaTest/internal/entity/goods"
	storages "LamodaTest/internal/entity/storages"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDb is a mock of Db interface.
type MockDb struct {
	ctrl     *gomock.Controller
	recorder *MockDbMockRecorder
}

// MockDbMockRecorder is the mock recorder for MockDb.
type MockDbMockRecorder struct {
	mock *MockDb
}

// NewMockDb creates a new mock instance.
func NewMockDb(ctrl *gomock.Controller) *MockDb {
	mock := &MockDb{ctrl: ctrl}
	mock.recorder = &MockDbMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDb) EXPECT() *MockDbMockRecorder {
	return m.recorder
}

// AvailableGoods mocks base method.
func (m *MockDb) AvailableGoods(ctx context.Context) (map[int]goods.RemainsDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AvailableGoods", ctx)
	ret0, _ := ret[0].(map[int]goods.RemainsDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AvailableGoods indicates an expected call of AvailableGoods.
func (mr *MockDbMockRecorder) AvailableGoods(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvailableGoods", reflect.TypeOf((*MockDb)(nil).AvailableGoods), ctx)
}

// GoodAdd mocks base method.
func (m *MockDb) GoodAdd(ctx context.Context, name, size string, uniqCode int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GoodAdd", ctx, name, size, uniqCode)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GoodAdd indicates an expected call of GoodAdd.
func (mr *MockDbMockRecorder) GoodAdd(ctx, name, size, uniqCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoodAdd", reflect.TypeOf((*MockDb)(nil).GoodAdd), ctx, name, size, uniqCode)
}

// GoodDelete mocks base method.
func (m *MockDb) GoodDelete(ctx context.Context, uniqCode int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GoodDelete", ctx, uniqCode)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GoodDelete indicates an expected call of GoodDelete.
func (mr *MockDbMockRecorder) GoodDelete(ctx, uniqCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoodDelete", reflect.TypeOf((*MockDb)(nil).GoodDelete), ctx, uniqCode)
}

// Goods mocks base method.
func (m *MockDb) Goods(ctx context.Context) ([]goods.Good, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Goods", ctx)
	ret0, _ := ret[0].([]goods.Good)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Goods indicates an expected call of Goods.
func (mr *MockDbMockRecorder) Goods(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Goods", reflect.TypeOf((*MockDb)(nil).Goods), ctx)
}

// ReleaseGood mocks base method.
func (m *MockDb) ReleaseGood(ctx context.Context, uniqId, count int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseGood", ctx, uniqId, count)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleaseGood indicates an expected call of ReleaseGood.
func (mr *MockDbMockRecorder) ReleaseGood(ctx, uniqId, count any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseGood", reflect.TypeOf((*MockDb)(nil).ReleaseGood), ctx, uniqId, count)
}

// ReserveGood mocks base method.
func (m *MockDb) ReserveGood(ctx context.Context, uniqId, count int) (map[int]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReserveGood", ctx, uniqId, count)
	ret0, _ := ret[0].(map[int]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReserveGood indicates an expected call of ReserveGood.
func (mr *MockDbMockRecorder) ReserveGood(ctx, uniqId, count any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReserveGood", reflect.TypeOf((*MockDb)(nil).ReserveGood), ctx, uniqId, count)
}

// Storages mocks base method.
func (m *MockDb) Storages(ctx context.Context, all bool) ([]storages.Storage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Storages", ctx, all)
	ret0, _ := ret[0].([]storages.Storage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Storages indicates an expected call of Storages.
func (mr *MockDbMockRecorder) Storages(ctx, all any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Storages", reflect.TypeOf((*MockDb)(nil).Storages), ctx, all)
}

// StoragesAdd mocks base method.
func (m *MockDb) StoragesAdd(ctx context.Context, name string, available bool) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoragesAdd", ctx, name, available)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoragesAdd indicates an expected call of StoragesAdd.
func (mr *MockDbMockRecorder) StoragesAdd(ctx, name, available any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoragesAdd", reflect.TypeOf((*MockDb)(nil).StoragesAdd), ctx, name, available)
}

// StoragesChangeAccess mocks base method.
func (m *MockDb) StoragesChangeAccess(ctx context.Context, id int, available bool) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoragesChangeAccess", ctx, id, available)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoragesChangeAccess indicates an expected call of StoragesChangeAccess.
func (mr *MockDbMockRecorder) StoragesChangeAccess(ctx, id, available any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoragesChangeAccess", reflect.TypeOf((*MockDb)(nil).StoragesChangeAccess), ctx, id, available)
}

// StoragesDelete mocks base method.
func (m *MockDb) StoragesDelete(ctx context.Context, id int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoragesDelete", ctx, id)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoragesDelete indicates an expected call of StoragesDelete.
func (mr *MockDbMockRecorder) StoragesDelete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoragesDelete", reflect.TypeOf((*MockDb)(nil).StoragesDelete), ctx, id)
}