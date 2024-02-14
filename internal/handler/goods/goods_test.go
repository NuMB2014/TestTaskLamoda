package goods

import (
	"LamodaTest/internal/entity/goods"
	"LamodaTest/internal/logger"
	"LamodaTest/internal/registry"
	mock_registry "LamodaTest/internal/registry/mocks"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
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
					m.EXPECT().GoodAdd(context.Background(), "test", "l", 1).Return(int64(1), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"name":      "test",
						"size":      "l",
						"uniq_code": 1,
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
					m.EXPECT().GoodAdd(context.Background(), "test", "l", 1).Return(int64(0), errors.New("test")).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"name":      "test",
						"size":      "l",
						"uniq_code": 1,
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
					m.EXPECT().GoodAdd(context.Background(), "test", "l", 1).Return(int64(0), nil).Times(1).AnyTimes()
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
						Goods(context.Background()).
						Return([]goods.Good{
							{
								Id:       1,
								Name:     "test",
								Size:     "l",
								UniqCode: 1,
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
				"data": []goods.Good{
					{
						Id:       1,
						Name:     "test",
						Size:     "l",
						UniqCode: 1,
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
					m.EXPECT().Goods(context.Background()).Return(nil, errors.New("test")).Times(1).AnyTimes()
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
					m.EXPECT().GoodDelete(context.Background(), 1).Return(int64(1), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "DELETE",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"uniq_code": 1,
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
					m.EXPECT().GoodDelete(context.Background(), 1).Return(int64(0), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "DELETE",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"uniq_code": 1,
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
					m.EXPECT().GoodDelete(context.Background(), 1).Return(int64(0), errors.New("test")).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "DELETE",
				body: func() string {
					marshal, _ := json.Marshal(map[string]interface{}{
						"uniq_code": 1,
					})
					return string(marshal)
				}(),
			},
			wantCode: 500,
			wantRes: map[string]interface{}{
				"code":    500,
				"message": "Can't delete this good",
			},
		}, {
			name: "invalid json",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().GoodDelete(context.Background(), 1).Return(int64(0), nil).Times(1).AnyTimes()
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

func TestHandler_Release(t *testing.T) {
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
					m.EXPECT().ReleaseGood(context.Background(), 1, 5).Return(nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal([]map[string]interface{}{{
						"uniq_code": 1,
						"count":     5,
					}})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code": 200,
				"data": []goods.ReleasedDTO{{
					1,
					"OK",
				}},
			},
		}, {
			name: "one normal, but one is corrupted",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().ReleaseGood(context.Background(), 1, 5).Return(nil).AnyTimes()
					m.EXPECT().ReleaseGood(context.Background(), 2, 1).Return(errors.New("test")).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal([]map[string]interface{}{
						{
							"uniq_code": 1,
							"count":     5,
						}, {
							"uniq_code": 2,
							"count":     1,
						},
					})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code": 200,
				"data": []goods.ReleasedDTO{
					{
						1,
						"OK",
					},
					{
						2,
						"can't release this good",
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
					m.EXPECT().ReleaseGood(context.Background(), 1, 5).Return(errors.New("test")).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal([]map[string]interface{}{
						{
							"uniq_code": 1,
							"count":     5,
						},
					})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code": 200,
				"data": []goods.ReleasedDTO{{
					1,
					"can't release this good",
				}},
			},
		}, {
			name: "invalid json",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().GoodDelete(context.Background(), 1).Return(int64(0), nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
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
			router.POST(ReleaseRoute, h.Release)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.args.method, ReleaseRoute, strings.NewReader(tt.args.body))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			bytes, _ := json.Marshal(tt.wantRes)
			assert.Equal(t, string(bytes), w.Body.String())
		})
	}
}

func TestHandler_Remains(t *testing.T) {
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
						AvailableGoods(context.Background()).
						Return(map[int]goods.RemainsDTO{
							1: {
								Name: "test",
								Size: "l",
								StorageAvailable: map[int]int{
									1: 1,
								},
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
				"data": map[int]goods.RemainsDTO{
					1: {
						Name: "test",
						Size: "l",
						StorageAvailable: map[int]int{
							1: 1,
						},
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
					m.EXPECT().AvailableGoods(context.Background()).Return(nil, errors.New("test")).Times(1).AnyTimes()
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
			router.GET(RemainsRoute, h.Remains)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", RemainsRoute, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			bytes, _ := json.Marshal(tt.wantRes)
			assert.Equal(t, string(bytes), w.Body.String())
		})
	}
}

func TestHandler_Reserve(t *testing.T) {
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
					m.EXPECT().ReserveGood(context.Background(), 1, 5).Return(map[int]int{1: 5}, nil).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal([]map[string]interface{}{{
						"uniq_code": 1,
						"count":     5,
					}})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code": 200,
				"data": []goods.ReservedDTO{{
					UniqCode: 1,
					Storages: []map[string]int{
						{
							"reserved": 5,
							"storage":  1,
						},
					},
				}},
			},
		}, {
			name: "one normal, but one is corrupted",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					m.EXPECT().ReserveGood(context.Background(), 1, 5).Return(map[int]int{1: 5}, nil).AnyTimes()
					m.EXPECT().ReserveGood(context.Background(), 2, 1).Return(nil, errors.New("test")).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal([]map[string]interface{}{
						{
							"uniq_code": 1,
							"count":     5,
						}, {
							"uniq_code": 2,
							"count":     1,
						},
					})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code": 200,
				"data": []goods.ReservedDTO{
					{
						UniqCode: 1,
						Storages: []map[string]int{
							{
								"reserved": 5,
								"storage":  1,
							},
						},
					},
					{
						UniqCode:       2,
						Storages:       []map[string]int{},
						AdditionalInfo: "Can't reserve this good",
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
					m.EXPECT().ReserveGood(context.Background(), 1, 5).Return(nil, errors.New("test")).Times(1).AnyTimes()
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
				body: func() string {
					marshal, _ := json.Marshal([]map[string]interface{}{
						{
							"uniq_code": 1,
							"count":     5,
						},
					})
					return string(marshal)
				}(),
			},
			wantCode: 200,
			wantRes: map[string]interface{}{
				"code": 200,
				"data": []goods.ReservedDTO{
					{
						UniqCode:       1,
						Storages:       []map[string]int{},
						AdditionalInfo: "Can't reserve this good",
					},
				},
			},
		}, {
			name: "invalid json",
			fields: fields{
				registry: func() *mock_registry.MockDb {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					m := mock_registry.NewMockDb(ctrl)
					return m
				}(),
				log: l,
			},
			args: args{
				method: "POST",
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
			router.POST(ReserveRoute, h.Reserve)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.args.method, ReserveRoute, strings.NewReader(tt.args.body))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			bytes, _ := json.Marshal(tt.wantRes)
			assert.Equal(t, string(bytes), w.Body.String())
		})
	}
}
