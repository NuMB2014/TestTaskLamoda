package registry

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func TestDatabase_ReserveGoods(t *testing.T) {
	type fields struct {
		conn *sql.DB
	}
	type args struct {
		ctx    context.Context
		uniqId int
		count  int
	}
	db, _ := sql.Open("mysql", "root:1@/Lamoda")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[int]int
		wantErr bool
	}{
		{
			"test1",
			fields{conn: db},
			args{
				ctx:    context.TODO(),
				uniqId: 1,
				count:  8,
			},
			nil,
			false,
		},
		{
			"test2",
			fields{conn: db},
			args{
				ctx:    context.TODO(),
				uniqId: 1,
				count:  8,
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.ReserveGoods(tt.args.ctx, tt.args.uniqId, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReserveGoods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReserveGoods() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_ReserveGoods2(t *testing.T) {
	type fields struct {
		conn *sql.DB
	}
	type args struct {
		ctx    context.Context
		uniqId int
		count  int
	}
	db, _ := sql.Open("mysql", "root:1@/Lamoda")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[int]int
		wantErr bool
	}{
		{
			"test1",
			fields{conn: db},
			args{
				ctx:    context.TODO(),
				uniqId: 1,
				count:  5,
			},
			nil,
			false,
		},
		{
			"test2",
			fields{conn: db},
			args{
				ctx:    context.TODO(),
				uniqId: 1,
				count:  5,
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				conn: tt.fields.conn,
			}
			got, err := d.ReserveGoods(tt.args.ctx, tt.args.uniqId, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReserveGoods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReserveGoods() got = %v, want %v", got, tt.want)
			}
		})
	}
}
