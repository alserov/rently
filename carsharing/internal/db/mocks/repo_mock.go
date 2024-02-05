// Code generated by MockGen. DO NOT EDIT.
// Source: internal/db/repository.go

// Package repomock is a generated GoMock package.
package repomock

import (
	context "context"
	reflect "reflect"
	time "time"

	db "github.com/alserov/rently/carsharing/internal/db"
	models "github.com/alserov/rently/carsharing/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CancelRentTx mocks base method.
func (m *MockRepository) CancelRentTx(ctx context.Context, rentUUID string) (models.CancelRentInfo, db.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelRentTx", ctx, rentUUID)
	ret0, _ := ret[0].(models.CancelRentInfo)
	ret1, _ := ret[1].(db.Tx)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CancelRentTx indicates an expected call of CancelRentTx.
func (mr *MockRepositoryMockRecorder) CancelRentTx(ctx, rentUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelRentTx", reflect.TypeOf((*MockRepository)(nil).CancelRentTx), ctx, rentUUID)
}

// CheckIfCarAvailableInPeriod mocks base method.
func (m *MockRepository) CheckIfCarAvailableInPeriod(ctx context.Context, carUUID string, from, to time.Time) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIfCarAvailableInPeriod", ctx, carUUID, from, to)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIfCarAvailableInPeriod indicates an expected call of CheckIfCarAvailableInPeriod.
func (mr *MockRepositoryMockRecorder) CheckIfCarAvailableInPeriod(ctx, carUUID, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfCarAvailableInPeriod", reflect.TypeOf((*MockRepository)(nil).CheckIfCarAvailableInPeriod), ctx, carUUID, from, to)
}

// CheckRent mocks base method.
func (m *MockRepository) CheckRent(ctx context.Context, rentUUID string) (models.Rent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRent", ctx, rentUUID)
	ret0, _ := ret[0].(models.Rent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRent indicates an expected call of CheckRent.
func (mr *MockRepositoryMockRecorder) CheckRent(ctx, rentUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRent", reflect.TypeOf((*MockRepository)(nil).CheckRent), ctx, rentUUID)
}

// CreateCar mocks base method.
func (m *MockRepository) CreateCar(ctx context.Context, car models.Car) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCar", ctx, car)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCar indicates an expected call of CreateCar.
func (mr *MockRepositoryMockRecorder) CreateCar(ctx, car interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCar", reflect.TypeOf((*MockRepository)(nil).CreateCar), ctx, car)
}

// CreateCharge mocks base method.
func (m *MockRepository) CreateCharge(ctx context.Context, req models.Charge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCharge", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCharge indicates an expected call of CreateCharge.
func (mr *MockRepositoryMockRecorder) CreateCharge(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCharge", reflect.TypeOf((*MockRepository)(nil).CreateCharge), ctx, req)
}

// CreateRentTx mocks base method.
func (m *MockRepository) CreateRentTx(ctx context.Context, req models.CreateRentReq) (float32, db.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRentTx", ctx, req)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(db.Tx)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRentTx indicates an expected call of CreateRentTx.
func (mr *MockRepositoryMockRecorder) CreateRentTx(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRentTx", reflect.TypeOf((*MockRepository)(nil).CreateRentTx), ctx, req)
}

// DeleteCar mocks base method.
func (m *MockRepository) DeleteCar(ctx context.Context, uuid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCar", ctx, uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCar indicates an expected call of DeleteCar.
func (mr *MockRepositoryMockRecorder) DeleteCar(ctx, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCar", reflect.TypeOf((*MockRepository)(nil).DeleteCar), ctx, uuid)
}

// GetAvailableCars mocks base method.
func (m *MockRepository) GetAvailableCars(ctx context.Context, period models.Period) ([]models.CarMainInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableCars", ctx, period)
	ret0, _ := ret[0].([]models.CarMainInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableCars indicates an expected call of GetAvailableCars.
func (mr *MockRepositoryMockRecorder) GetAvailableCars(ctx, period interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableCars", reflect.TypeOf((*MockRepository)(nil).GetAvailableCars), ctx, period)
}

// GetCarByUUID mocks base method.
func (m *MockRepository) GetCarByUUID(ctx context.Context, uuid string) (models.Car, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCarByUUID", ctx, uuid)
	ret0, _ := ret[0].(models.Car)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCarByUUID indicates an expected call of GetCarByUUID.
func (mr *MockRepositoryMockRecorder) GetCarByUUID(ctx, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCarByUUID", reflect.TypeOf((*MockRepository)(nil).GetCarByUUID), ctx, uuid)
}

// GetCarsByParams mocks base method.
func (m *MockRepository) GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.CarMainInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCarsByParams", ctx, params)
	ret0, _ := ret[0].([]models.CarMainInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCarsByParams indicates an expected call of GetCarsByParams.
func (mr *MockRepositoryMockRecorder) GetCarsByParams(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCarsByParams", reflect.TypeOf((*MockRepository)(nil).GetCarsByParams), ctx, params)
}

// GetRentsWhatStartsOnDate mocks base method.
func (m *MockRepository) GetRentsWhatStartsOnDate(ctx context.Context, date time.Time) ([]models.RentStartData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentsWhatStartsOnDate", ctx, date)
	ret0, _ := ret[0].([]models.RentStartData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentsWhatStartsOnDate indicates an expected call of GetRentsWhatStartsOnDate.
func (mr *MockRepositoryMockRecorder) GetRentsWhatStartsOnDate(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentsWhatStartsOnDate", reflect.TypeOf((*MockRepository)(nil).GetRentsWhatStartsOnDate), ctx, date)
}

// UpdateCarPrice mocks base method.
func (m *MockRepository) UpdateCarPrice(ctx context.Context, req models.UpdateCarPriceReq) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCarPrice", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCarPrice indicates an expected call of UpdateCarPrice.
func (mr *MockRepositoryMockRecorder) UpdateCarPrice(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCarPrice", reflect.TypeOf((*MockRepository)(nil).UpdateCarPrice), ctx, req)
}

// MockAdminRepository is a mock of AdminRepository interface.
type MockAdminRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAdminRepositoryMockRecorder
}

// MockAdminRepositoryMockRecorder is the mock recorder for MockAdminRepository.
type MockAdminRepositoryMockRecorder struct {
	mock *MockAdminRepository
}

// NewMockAdminRepository creates a new mock instance.
func NewMockAdminRepository(ctrl *gomock.Controller) *MockAdminRepository {
	mock := &MockAdminRepository{ctrl: ctrl}
	mock.recorder = &MockAdminRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdminRepository) EXPECT() *MockAdminRepositoryMockRecorder {
	return m.recorder
}

// CreateCar mocks base method.
func (m *MockAdminRepository) CreateCar(ctx context.Context, car models.Car) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCar", ctx, car)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCar indicates an expected call of CreateCar.
func (mr *MockAdminRepositoryMockRecorder) CreateCar(ctx, car interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCar", reflect.TypeOf((*MockAdminRepository)(nil).CreateCar), ctx, car)
}

// DeleteCar mocks base method.
func (m *MockAdminRepository) DeleteCar(ctx context.Context, uuid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCar", ctx, uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCar indicates an expected call of DeleteCar.
func (mr *MockAdminRepositoryMockRecorder) DeleteCar(ctx, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCar", reflect.TypeOf((*MockAdminRepository)(nil).DeleteCar), ctx, uuid)
}

// UpdateCarPrice mocks base method.
func (m *MockAdminRepository) UpdateCarPrice(ctx context.Context, req models.UpdateCarPriceReq) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCarPrice", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCarPrice indicates an expected call of UpdateCarPrice.
func (mr *MockAdminRepositoryMockRecorder) UpdateCarPrice(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCarPrice", reflect.TypeOf((*MockAdminRepository)(nil).UpdateCarPrice), ctx, req)
}

// MockCarRepository is a mock of CarRepository interface.
type MockCarRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCarRepositoryMockRecorder
}

// MockCarRepositoryMockRecorder is the mock recorder for MockCarRepository.
type MockCarRepositoryMockRecorder struct {
	mock *MockCarRepository
}

// NewMockCarRepository creates a new mock instance.
func NewMockCarRepository(ctrl *gomock.Controller) *MockCarRepository {
	mock := &MockCarRepository{ctrl: ctrl}
	mock.recorder = &MockCarRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCarRepository) EXPECT() *MockCarRepositoryMockRecorder {
	return m.recorder
}

// GetAvailableCars mocks base method.
func (m *MockCarRepository) GetAvailableCars(ctx context.Context, period models.Period) ([]models.CarMainInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableCars", ctx, period)
	ret0, _ := ret[0].([]models.CarMainInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableCars indicates an expected call of GetAvailableCars.
func (mr *MockCarRepositoryMockRecorder) GetAvailableCars(ctx, period interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableCars", reflect.TypeOf((*MockCarRepository)(nil).GetAvailableCars), ctx, period)
}

// GetCarByUUID mocks base method.
func (m *MockCarRepository) GetCarByUUID(ctx context.Context, uuid string) (models.Car, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCarByUUID", ctx, uuid)
	ret0, _ := ret[0].(models.Car)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCarByUUID indicates an expected call of GetCarByUUID.
func (mr *MockCarRepositoryMockRecorder) GetCarByUUID(ctx, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCarByUUID", reflect.TypeOf((*MockCarRepository)(nil).GetCarByUUID), ctx, uuid)
}

// GetCarsByParams mocks base method.
func (m *MockCarRepository) GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.CarMainInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCarsByParams", ctx, params)
	ret0, _ := ret[0].([]models.CarMainInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCarsByParams indicates an expected call of GetCarsByParams.
func (mr *MockCarRepositoryMockRecorder) GetCarsByParams(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCarsByParams", reflect.TypeOf((*MockCarRepository)(nil).GetCarsByParams), ctx, params)
}

// MockRentRepository is a mock of RentRepository interface.
type MockRentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRentRepositoryMockRecorder
}

// MockRentRepositoryMockRecorder is the mock recorder for MockRentRepository.
type MockRentRepositoryMockRecorder struct {
	mock *MockRentRepository
}

// NewMockRentRepository creates a new mock instance.
func NewMockRentRepository(ctrl *gomock.Controller) *MockRentRepository {
	mock := &MockRentRepository{ctrl: ctrl}
	mock.recorder = &MockRentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRentRepository) EXPECT() *MockRentRepositoryMockRecorder {
	return m.recorder
}

// CancelRentTx mocks base method.
func (m *MockRentRepository) CancelRentTx(ctx context.Context, rentUUID string) (models.CancelRentInfo, db.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelRentTx", ctx, rentUUID)
	ret0, _ := ret[0].(models.CancelRentInfo)
	ret1, _ := ret[1].(db.Tx)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CancelRentTx indicates an expected call of CancelRentTx.
func (mr *MockRentRepositoryMockRecorder) CancelRentTx(ctx, rentUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelRentTx", reflect.TypeOf((*MockRentRepository)(nil).CancelRentTx), ctx, rentUUID)
}

// CheckIfCarAvailableInPeriod mocks base method.
func (m *MockRentRepository) CheckIfCarAvailableInPeriod(ctx context.Context, carUUID string, from, to time.Time) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIfCarAvailableInPeriod", ctx, carUUID, from, to)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIfCarAvailableInPeriod indicates an expected call of CheckIfCarAvailableInPeriod.
func (mr *MockRentRepositoryMockRecorder) CheckIfCarAvailableInPeriod(ctx, carUUID, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfCarAvailableInPeriod", reflect.TypeOf((*MockRentRepository)(nil).CheckIfCarAvailableInPeriod), ctx, carUUID, from, to)
}

// CheckRent mocks base method.
func (m *MockRentRepository) CheckRent(ctx context.Context, rentUUID string) (models.Rent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRent", ctx, rentUUID)
	ret0, _ := ret[0].(models.Rent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRent indicates an expected call of CheckRent.
func (mr *MockRentRepositoryMockRecorder) CheckRent(ctx, rentUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRent", reflect.TypeOf((*MockRentRepository)(nil).CheckRent), ctx, rentUUID)
}

// CreateCharge mocks base method.
func (m *MockRentRepository) CreateCharge(ctx context.Context, req models.Charge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCharge", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCharge indicates an expected call of CreateCharge.
func (mr *MockRentRepositoryMockRecorder) CreateCharge(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCharge", reflect.TypeOf((*MockRentRepository)(nil).CreateCharge), ctx, req)
}

// CreateRentTx mocks base method.
func (m *MockRentRepository) CreateRentTx(ctx context.Context, req models.CreateRentReq) (float32, db.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRentTx", ctx, req)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(db.Tx)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRentTx indicates an expected call of CreateRentTx.
func (mr *MockRentRepositoryMockRecorder) CreateRentTx(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRentTx", reflect.TypeOf((*MockRentRepository)(nil).CreateRentTx), ctx, req)
}

// GetRentsWhatStartsOnDate mocks base method.
func (m *MockRentRepository) GetRentsWhatStartsOnDate(ctx context.Context, date time.Time) ([]models.RentStartData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentsWhatStartsOnDate", ctx, date)
	ret0, _ := ret[0].([]models.RentStartData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentsWhatStartsOnDate indicates an expected call of GetRentsWhatStartsOnDate.
func (mr *MockRentRepositoryMockRecorder) GetRentsWhatStartsOnDate(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentsWhatStartsOnDate", reflect.TypeOf((*MockRentRepository)(nil).GetRentsWhatStartsOnDate), ctx, date)
}

// MockTx is a mock of Tx interface.
type MockTx struct {
	ctrl     *gomock.Controller
	recorder *MockTxMockRecorder
}

// MockTxMockRecorder is the mock recorder for MockTx.
type MockTxMockRecorder struct {
	mock *MockTx
}

// NewMockTx creates a new mock instance.
func NewMockTx(ctrl *gomock.Controller) *MockTx {
	mock := &MockTx{ctrl: ctrl}
	mock.recorder = &MockTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTx) EXPECT() *MockTxMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockTx) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTxMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTx)(nil).Commit))
}

// Rollback mocks base method.
func (m *MockTx) Rollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTxMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTx)(nil).Rollback))
}
