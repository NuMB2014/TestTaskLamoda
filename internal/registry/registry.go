package registry

import (
	"LamodaTest/internal/entity/storages"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Db interface {
	Storages(ctx context.Context) ([]storages.Storage, error)
	AvailableGoods(ctx context.Context, log logrus.FieldLogger)
	ReserveGoods(ctx context.Context, log logrus.FieldLogger)
	ReleaseGoods(ctx context.Context, log logrus.FieldLogger)
}

type Database struct {
	db *sql.DB
}

func New(connect *sql.DB) *Database {
	return &Database{db: connect}
}

func (d *Database) Storages(ctx context.Context) ([]storages.Storage, error) {
	cmd, err := d.db.Prepare("select * from storages;")
	rows, err := cmd.QueryContext(ctx)
	var result []storages.Storage
	ctx.Done()
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("can't scan from storage list: %s", err.Error())
	}
	for rows.Next() {
		values := storages.Storage{}
		err = rows.Scan(&values.ID, &values.Name, &values.Available)
		if err != nil {
			return nil, fmt.Errorf("can't scan from storage list: %s", err.Error())
		}
		result = append(result, values)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error when try get all storages: %v", err)
	}
	return result, nil
}

func (d *Database) AvailableGoods(ctx context.Context, log logrus.FieldLogger) (map[int]interface{}, error) {
	cmd, err := d.db.Prepare("select goods.name, goods.size, goods.uniq_code, storages.id AS storage_id, remains.count - remains.reserved AS avail FROM goods JOIN remains ON goods.id = remains.good_id JOIN storages ON remains.storage_id = storages.id WHERE remains.count > reserved AND available = 1")
	rows, err := cmd.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't query avail goods: %s", err.Error())
	}
	result := map[int]interface{}{}
	ctx.Done()
	defer rows.Close()
	for rows.Next() {
		var tmp struct {
			Name    string
			Size    string
			UniqId  int
			Storage int
			Avail   int
		}
		err = rows.Scan(&tmp.Name, &tmp.Size, &tmp.UniqId, &tmp.Storage, &tmp.Avail)
		if err != nil {
			return nil, fmt.Errorf("can't scan from rows: %s", err.Error())
		}

		if data, ok := result[tmp.UniqId]; !ok {
			result[tmp.UniqId] = map[string]interface{}{
				"name": tmp.Name,
				"size": tmp.Size,
				"storage_available": map[int]interface{}{
					tmp.Storage: tmp.Avail,
				},
			}
		} else {
			note := data.(map[string]interface{})
			note["storage_available"].(map[int]interface{})[tmp.Storage] = tmp.Avail
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error when try get all storages: %v", err)
	}
	return result, nil
}

func (d *Database) ReserveGoods(ctx context.Context, log logrus.FieldLogger) {

}

func (d *Database) ReleaseGoods(ctx context.Context, log logrus.FieldLogger) {

}
