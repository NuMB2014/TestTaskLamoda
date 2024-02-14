package storages

import (
	"LamodaTest/internal/entity/storages"
	"LamodaTest/internal/logger"
	"LamodaTest/internal/registry"
	mock_registry "LamodaTest/internal/registry/mocks"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Add(t *testing.T) {
	type fields struct {
		registry registry.Db
		log      logrus.FieldLogger
	}
	type args struct {
		method string
		body   string
	}
	l := logger.New(false)
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRes  map[string]interface{}
		wantCode int
	}{
		{
			name: "normal",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesAdd(context.Background(), "test", true).Return(int64(1), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"name":      "test",
						"available": true,
					})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code": 200,
				"data": 1,
			},
		}, {
			name: "err from db",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesAdd(context.Background(), "test", true).Return(int64(0), errors.New("test")).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"name":      "test",
						"available": true,
					})
					return string(marshal)
				}(),
			},
			wantCode: 500,
			wantRes: map[string]interface{}{
				"code":    500,
				"message": "Not added",
			},
		}, {
			name: "invalid json",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesAdd(context.Background(), "test", true).Return(int64(0), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"available": true,
					})
					return string(marshal)
				}(),
			},
			wantCode: http.StatusBadRequest,
			wantRes: map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid JSON",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				registry: tt.fields.registry,
				log:      tt.fields.log,
			}
			router := gin.Default()
			gin.SetMode(gin.ReleaseMode)
			router.POST(AddRoute, h.Add)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.args.method, AddRoute, strings.NewReader(tt.args.body))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			bytes, _ := json.Marshal(tt.wantRes)
			assert.Equal(t, string(bytes), w.Body.String())
		})
	}
}

func TestHandler_All(t *testing.T) {
	type fields struct {
		registry registry.Db
		log      logrus.FieldLogger
	}
	l := logger.New(false)
	tests := []struct {
		name     string
		fields   fields
		wantRes  map[string]interface{}
		wantCode int
	}{
		{
			name: "normal",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().
						Storages(context.Background(), true).
						Return([]storages.Storage{
							{
								ID:           1,
								Name:         "test",
								RawAvailable: "1",
								Available:    true,
							},
						}, nil).
						Times(1).
						AnyTimes()
					return m
				}(),
				log: l,
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code": 200,
				"data": []storages.Storage{
					{ID: 1,
						Name:         "test",
						RawAvailable: "1",
						Available:    true,
					},
				},
			},
		}, {
			name: "err from db",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().Storages(context.Background(), true).Return(nil, errors.New("test")).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			wantCode: 500,
			wantRes: map[string]interface{}{
				"code":    500,
				"message": "Internal server error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				registry: tt.fields.registry,
				log:      tt.fields.log,
			}
			router := gin.Default()
			gin.SetMode(gin.ReleaseMode)
			router.GET(AllRoute, h.All)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", AllRoute, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			bytes, _ := json.Marshal(tt.wantRes)
			assert.Equal(t, string(bytes), w.Body.String())
		})
	}
}

func TestHandler_Available(t *testing.T) {
	type fields struct {
		registry registry.Db
		log      logrus.FieldLogger
	}
	l := logger.New(false)
	tests := []struct {
		name     string
		fields   fields
		wantRes  map[string]interface{}
		wantCode int
	}{
		{
			name: "normal",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().
						Storages(context.Background(), false).
						Return([]storages.Storage{
							{
								ID:           1,
								Name:         "test",
								RawAvailable: "1",
								Available:    true,
							},
						}, nil).
						Times(1).
						AnyTimes()
					return m
				}(),
				log: l,
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code": 200,
				"data": []storages.Storage{
					{ID: 1,
						Name:         "test",
						RawAvailable: "1",
						Available:    true,
					},
				},
			},
		}, {
			name: "err from db",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().Storages(context.Background(), false).Return(nil, errors.New("test")).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			wantCode: 500,
			wantRes: map[string]interface{}{
				"code":    500,
				"message": "Internal server error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				registry: tt.fields.registry,
				log:      tt.fields.log,
			}
			router := gin.Default()
			gin.SetMode(gin.ReleaseMode)
			router.GET(AvailableRoute, h.Available)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", AvailableRoute, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			bytes, _ := json.Marshal(tt.wantRes)
			assert.Equal(t, string(bytes), w.Body.String())
		})
	}
}

func TestHandler_ChangeAccess(t *testing.T) {
	type fields struct {
		registry registry.Db
		log      logrus.FieldLogger
	}
	type args struct {
		method string
		body   string
	}
	l := logger.New(false)
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRes  map[string]interface{}
		wantCode int
	}{
		{
			name: "normal",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesChangeAccess(context.Background(), 1, true).Return(int64(1), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"id":        1,
						"available": true,
					})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code":    200,
				"message": "OK",
			},
		}, {
			name: "no one changed",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesChangeAccess(context.Background(), 1, true).Return(int64(0), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"id":        1,
						"available": true,
					})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code":    200,
				"message": "no records are changed",
			},
		}, {
			name: "err from db",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesChangeAccess(context.Background(), 1, true).Return(int64(0), errors.New("test")).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"id":        1,
						"available": true,
					})
					return string(marshal)
				}(),
			},
			wantCode: 500,
			wantRes: map[string]interface{}{
				"code":    500,
				"message": "Can't change this storage",
			},
		}, {
			name: "invalid json",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesChangeAccess(context.Background(), 1, true).Return(int64(0), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"available": true,
					})
					return string(marshal)
				}(),
			},
			wantCode: http.StatusBadRequest,
			wantRes: map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid JSON",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				registry: tt.fields.registry,
				log:      tt.fields.log,
			}
			router := gin.Default()
			gin.SetMode(gin.ReleaseMode)
			router.POST(AccessStatus, h.ChangeAccess)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.args.method, AccessStatus, strings.NewReader(tt.args.body))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			bytes, _ := json.Marshal(tt.wantRes)
			assert.Equal(t, string(bytes), w.Body.String())
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type fields struct {
		registry registry.Db
		log      logrus.FieldLogger
	}
	type args struct {
		method string
		body   string
	}
	l := logger.New(false)
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRes  map[string]interface{}
		wantCode int
	}{
		{
			name: "normal",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesDelete(context.Background(), 1).Return(int64(1), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "DELETE",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"id": 1,
					})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code":    200,
				"message": "OK",
			},
		}, {
			name: "no one deleted",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesDelete(context.Background(), 1).Return(int64(0), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "DELETE",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"id": 1,
					})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code":    200,
				"message": "no records are deleted",
			},
		}, {
			name: "err from db",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesDelete(context.Background(), 1).Return(int64(0), errors.New("test")).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "DELETE",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"id": 1,
					})
					return string(marshal)
				}(),
			},
			wantCode: 500,
			wantRes: map[string]interface{}{
				"code":    500,
				"message": "Can't delete this storage",
			},
		}, {
			name: "invalid json",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().StoragesDelete(context.Background(), 1).Return(int64(0), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "DELETE",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"d": 1,
					})
					return string(marshal)
				}(),
			},
			wantCode: http.StatusBadRequest,
			wantRes: map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid JSON",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				registry: tt.fields.registry,
				log:      tt.fields.log,
			}
			router := gin.Default()
			gin.SetMode(gin.ReleaseMode)
			router.DELETE(DeleteRoute, h.Delete)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.args.method, DeleteRoute, strings.NewReader(tt.args.body))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			bytes, _ := json.Marshal(tt.wantRes)
			assert.Equal(t, string(bytes), w.Body.String())
		})
	}
}
