package registry

import (
	"LamodaTest/internal/entity/goods"
	"LamodaTest/internal/entity/storages"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

func TestDatabase_AvailableGoods(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx context.Context
	}
	columns := []string{"name", "size", "uniq_code", "storage_id", "avail"}
	sqlStr := "select goods.name, goods.size, goods.uniq_code, storages.id AS storage_id, remains.count - remains.reserved AS avail FROM goods JOIN remains ON goods.id = remains.good_id JOIN storages ON remains.storage_id = storages.id WHERE remains.count > reserved AND available = 1"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[int]goods.RemainsDTO
		wantErr bool
	}{
		{
			name: "normal",
			fields: func() fields {
				db, mock, _ := sqlmock.New()
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("Test,xs,1,1,500\nTest2,l,2,4,200\nTest,xs,1,3,500\nTest2,l,2,3,200"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args: args{ctx: context.TODO()},
			want: map[int]goods.RemainsDTO{
				1: {
					Name: "Test",
					Size: "xs",
					StorageAvailable: map[int]int{
						1: 500,
						3: 500,
					},
				},
				2: {
					Name: "Test2",
					Size: "l",
					StorageAvailable: map[int]int{
						4: 200,
						3: 200,
					},
				},
			},
			wantErr: false,
		}, {
			name: "invalid req",
			fields: func() fields {
				db, mock, _ := sqlmock.New()
				mock.ExpectPrepare("select * from sys").ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("Test,xs,1,1,500\nTest2,l,2,4,200\nTest,xs,1,3,500\nTest2,l,2,3,200"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: true,
		}, {
			name: "empty table",
			fields: func() fields {
				db, mock, _ := sqlmock.New()
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString(""))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    map[int]goods.RemainsDTO{},
			wantErr: false,
		}, {
			name: "err then doing req",
			fields: func() fields {
				db, mock, _ := sqlmock.New()
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: true,
		}, {
			name: "err then doing get rows ",
			fields: func() fields {
				db, mock, _ := sqlmock.New()
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("null"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: true,
		}, {
			name: "rows.Err()",
			fields: func() fields {
				db, mock, _ := sqlmock.New()
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("null").RowError(0, errors.New("test")))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.AvailableGoods(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("AvailableGoods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AvailableGoods() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GoodAdd(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx      context.Context
		name     string
		size     string
		uniqCode int
	}
	sqlStr := "insert into goods (name, size, uniq_code) values (?, ?, ?)"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "normal",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs("test", "xs", 1).WillReturnResult(sqlmock.NewResult(1, 1))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), "test", "xs", 1},
			want:    1,
			wantErr: false,
		}, {
			name: "err sql",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs("test", "xs", 1).WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), "test", "xs", 1},
			want:    -1,
			wantErr: true,
		}, {
			name: "err last inserted id",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs("test", "xs", 1).WillReturnResult(sqlmock.NewErrorResult(errors.New("test")))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), "test", "xs", 1},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.GoodAdd(tt.args.ctx, tt.args.name, tt.args.size, tt.args.uniqCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoodAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GoodAdd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GoodDelete(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx      context.Context
		uniqCode int
	}
	sqlStr := "delete from goods where uniq_code = ?"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "normal",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1},
			want:    1,
			wantErr: false,
		}, {
			name: "err sql",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs(1).WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1},
			want:    -1,
			wantErr: true,
		}, {
			name: "err last inserted id",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(errors.New("test")))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.GoodDelete(tt.args.ctx, tt.args.uniqCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoodDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GoodDelete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_Goods(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx context.Context
	}
	sqlStr := "select * from goods;"
	columns := []string{"id", "name", "size", "uniq_code"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []goods.Good
		wantErr bool
	}{
		{
			name: "normal",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,Test,xs,1\n2,Test2,l,2"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args: args{ctx: context.TODO()},
			want: []goods.Good{{
				Id:       1,
				Name:     "Test",
				Size:     "xs",
				UniqCode: 1,
			},
				{
					Id:       2,
					Name:     "Test2",
					Size:     "l",
					UniqCode: 2,
				},
			},
			wantErr: false,
		}, {
			name: "invalid req",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare("select * from sys;").ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,Test,xs,1\n2,Test2,l,2"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: true,
		}, {
			name: "empty table",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString(""))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    []goods.Good{},
			wantErr: false,
		}, {
			name: "err then doing req",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: true,
		}, {
			name: "err then doing get rows ",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("null"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: true,
		}, {
			name: "rows.Err()",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("null").RowError(0, errors.New("test")))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.Goods(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Goods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Goods() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_ReleaseGood(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx    context.Context
		uniqId int
		count  int
	}
	columns := []string{"id", "storage_id", "reserved"}
	sqlStr := "SELECT remains.id, remains.storage_id, remains.reserved from remains JOIN storages ON storages.id = remains.storage_id where good_id = ? AND storages.available = 1"
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved - ? WHERE id = ?").WithArgs(15, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: false,
		},
		{
			name: "not enough reserved",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved - ? WHERE id = ?").WithArgs(15, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 20},
			wantErr: true,
		},
		{
			name: "err while begin transaction",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin().WillReturnError(errors.New("test"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "empty remains",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString(""))
				mock.ExpectExec("UPDATE remains SET reserved = reserved - ? WHERE id = ?").WithArgs(15, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "err while find id from uniq_code",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnError(errors.New("test"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		}, {
			name: "err [sql.ErrNoRows] while find id from uniq_code",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "err rows",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "err in rows scan",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnError(errors.New("test"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved - ? WHERE id = ?").WithArgs(15, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "err in update if reserved < count",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved - ? WHERE id = ?").WithArgs(15, 1).WillReturnError(errors.New("test"))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 20},
			wantErr: true,
		},
		{
			name: "err in update if reserved >= count",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved - ? WHERE id = ?").WithArgs(15, 1).WillReturnError(errors.New("test"))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "transaction commit error",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved - ? WHERE id = ?").WithArgs(15, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(errors.New("test"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			if err := d.ReleaseGood(tt.args.ctx, tt.args.uniqId, tt.args.count); (err != nil) != tt.wantErr {
				t.Errorf("ReleaseGood() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_ReserveGood(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx    context.Context
		uniqId int
		count  int
	}
	columns := []string{"id", "storage_id", "avail"}
	sqlStr := "SELECT remains.id, remains.storage_id, remains.count - remains.reserved AS avail from remains JOIN storages ON storages.id = remains.storage_id where good_id = ? AND storages.available = 1"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[int]int
		wantErr bool
	}{
		{
			name: "normal avail >= count",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved + ? WHERE id = ?").WithArgs(15, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			want:    map[int]int{1: 15},
			wantErr: false,
		},
		{
			name: "normal avail < count",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,10\n2,2,10"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved + ? WHERE id = ?").WithArgs(10, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE remains SET reserved = reserved + ? WHERE id = ?").WithArgs(5, 2).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			want:    map[int]int{1: 10, 2: 5},
			wantErr: false,
		},
		{
			name: "not enough reserved",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 20},
			wantErr: true,
		},
		{
			name: "err while begin transaction",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin().WillReturnError(errors.New("test"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "empty remains",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString(""))
				mock.ExpectExec("UPDATE remains SET reserved = reserved + ? WHERE id = ?").WithArgs(15, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "err while find id from uniq_code",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnError(errors.New("test"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		}, {
			name: "err [sql.ErrNoRows] while find id from uniq_code",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "err rows",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "err in rows scan",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnError(errors.New("test"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved + ? WHERE id = ?").WithArgs(15, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "err in update if reserved < count",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved + ? WHERE id = ?").WithArgs(15, 1).WillReturnError(errors.New("test"))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 20},
			wantErr: true,
		},
		{
			name: "err in update if reserved >= count",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved + ? WHERE id = ?").WithArgs(15, 1).WillReturnError(errors.New("test"))
				mock.ExpectCommit()
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
		{
			name: "transaction commit error",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id from goods where uniq_code = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
				mock.ExpectQuery(sqlStr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,15"))
				mock.ExpectExec("UPDATE remains SET reserved = reserved + ? WHERE id = ?").WithArgs(15, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(errors.New("test"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, 15},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.ReserveGood(tt.args.ctx, tt.args.uniqId, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReserveGood() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReserveGood() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_Storages(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx context.Context
		all bool
	}
	sqlStr := "select * from storages"
	columns := []string{"id", "name", "available"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []storages.Storage
		wantErr bool
	}{
		{
			name: "normal all",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,test1,1\n2,test2,1\n3,test2,0"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args: args{context.TODO(), true},
			want: []storages.Storage{
				{1, "test1", "1", true},
				{2, "test2", "1", true},
				{3, "test2", "0", false},
			},
			wantErr: false,
		}, {
			name: "normal available",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(fmt.Sprintf("%s where available = 1", sqlStr)).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,test1,1\n2,test2,1\n3,test2,0"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args: args{context.TODO(), false},
			want: []storages.Storage{
				{1, "test1", "1", true},
				{2, "test2", "1", true},
				{3, "test2", "0", false},
			},
			wantErr: false,
		},
		{
			name: "incorrect sql",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,test1,1\n2,test2,1\n3,test2,0"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), false},
			want:    nil,
			wantErr: true,
		},
		{
			name: "err doing query",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), true},
			want:    nil,
			wantErr: true,
		},
		{
			name: "err scan from rows",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("null"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), true},
			want:    nil,
			wantErr: true,
		},
		{
			name: "arr in ParseBool",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,test1,1\n2,test2,2\n3,test2,0"))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), true},
			want:    nil,
			wantErr: true,
		}, {
			name: "arr in ParseBool",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectPrepare(sqlStr).ExpectQuery().WillReturnRows(sqlmock.NewRows(columns).CloseError(errors.New("test")))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), true},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.Storages(tt.args.ctx, tt.args.all)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_StoragesAdd(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx       context.Context
		name      string
		available bool
	}
	sqlStr := "insert into storages (name, available) values (?, ?)"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "normal",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs("test", true).WillReturnResult(sqlmock.NewResult(1, 1))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), "test", true},
			want:    1,
			wantErr: false,
		},
		{
			name: "err sql",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs("test", true).WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), "test", true},
			want:    -1,
			wantErr: true,
		}, {
			name: "err last inserted id",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs("test", true).WillReturnResult(sqlmock.NewErrorResult(errors.New("test")))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), "test", true},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.StoragesAdd(tt.args.ctx, tt.args.name, tt.args.available)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoragesAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoragesAdd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_StoragesChangeAccess(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx       context.Context
		id        int
		available bool
	}
	sqlStr := "update storages set available = ? where id = ?"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "normal",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs(false, 1).WillReturnResult(sqlmock.NewResult(0, 1))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, false},
			want:    1,
			wantErr: false,
		}, {
			name: "err sql",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs(false, 1).WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, false},
			want:    -1,
			wantErr: true,
		}, {
			name: "err RowsAffected",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs(false, 1).WillReturnResult(sqlmock.NewErrorResult(errors.New("test")))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1, false},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.StoragesChangeAccess(tt.args.ctx, tt.args.id, tt.args.available)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoragesChangeAccess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoragesChangeAccess() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_StoragesDelete(t *testing.T) {
	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	type args struct {
		ctx context.Context
		id  int
	}
	sqlStr := "delete from storages where id = ?"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "normal",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1},
			want:    1,
			wantErr: false,
		}, {
			name: "err sql",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs(1).WillReturnError(sql.ErrNoRows)
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1},
			want:    -1,
			wantErr: true,
		}, {
			name: "err RowsAffected",
			fields: func() fields {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				mock.ExpectExec(sqlStr).WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(errors.New("test")))
				tmp := fields{
					conn: db,
					mock: mock,
				}
				return tmp
			}(),
			args:    args{context.TODO(), 1},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.StoragesDelete(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoragesDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoragesDelete() got = %v, want %v", got, tt.want)
			}
		})
	}
}
