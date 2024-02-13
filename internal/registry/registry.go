package registry

import (
	"LamodaTest/internal/entity/goods"
	"LamodaTest/internal/entity/storages"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Db interface {
	Storages(ctx context.Context) ([]storages.Storage, error)
	AvailableGoods(ctx context.Context) (map[int]goods.RemainsDTO, error)
	ReserveGood(ctx context.Context, uniqId int, count int) (map[int]int, error)
	ReleaseGood(ctx context.Context)
}

type Database struct {
	conn *sql.DB
}

func New(connect *sql.DB) *Database {
	return &Database{conn: connect}
}

func (d *Database) Storages(ctx context.Context) ([]storages.Storage, error) {
	cmd, err := d.conn.Prepare("select * from storages;")
	rows, err := cmd.QueryContext(ctx)
	var result []storages.Storage
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

func (d *Database) AvailableGoods(ctx context.Context) (map[int]goods.RemainsDTO, error) {
	cmd, err := d.conn.Prepare(`select 
			goods.name, 
			goods.size, 
			goods.uniq_code, 
			storages.id AS storage_id, 
			remains.count - remains.reserved AS avail 
		FROM 
		    goods 
		JOIN remains ON goods.id = remains.good_id 
		JOIN storages ON remains.storage_id = storages.id 
		WHERE remains.count > reserved AND available = 1`)
	rows, err := cmd.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't query avail goods: %s", err.Error())
	}
	result := map[int]goods.RemainsDTO{}
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

		if note, ok := result[tmp.UniqId]; !ok {
			result[tmp.UniqId] = goods.RemainsDTO{
				Name: tmp.Name,
				Size: tmp.Size,
				StorageAvailable: map[int]int{
					tmp.Storage: tmp.Avail,
				},
			}
		} else {
			note.StorageAvailable[tmp.Storage] = tmp.Avail
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error when try get all storages: %v", err)
	}
	return result, nil
}

func (d *Database) ReserveGood(ctx context.Context, uniqId int, count int) (map[int]int, error) {
	tx, err := d.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable}) //
	if err != nil {
		return nil, fmt.Errorf("can't init transaction: %w", err)
	}
	defer tx.Rollback()
	var id int
	if err = tx.QueryRowContext(ctx, "SELECT id from goods where uniq_code = ?",
		uniqId).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("not found goods")
		}
		return nil, err
	}
	reserved := map[int]int{}
	rows, err := d.conn.QueryContext(ctx, "SELECT remains.id, remains.storage_id, remains.count - remains.reserved AS avail from remains where good_id = ?", id)
	for rows.Next() {
		var tmp struct {
			Id        int
			storageId int
			Avail     int
		}
		err = rows.Scan(&tmp.Id, &tmp.storageId, &tmp.Avail)
		if err != nil {
			return nil, fmt.Errorf("can't get remains by %d good: %w", id, err)
		}
		if tmp.Avail > 0 {
			if tmp.Avail < count {
				_, err = tx.ExecContext(ctx, "UPDATE remains SET reserved = reserved + ? WHERE id = ?",
					tmp.Avail, tmp.Id)
				if err != nil {
					return nil, fmt.Errorf("can't reserve good by %d id: %w", tmp.Id, err)
				}
				count = count - tmp.Avail
				reserved[tmp.storageId] = tmp.Avail
			} else {
				_, err = tx.ExecContext(ctx, "UPDATE remains SET reserved = reserved + ? WHERE id = ?",
					count, tmp.Id)
				if err != nil {
					return nil, fmt.Errorf("can't reserve good by %d id: %w", tmp.Id, err)
				}
				reserved[tmp.storageId] = count
				count = 0
			}
		}
	}
	defer rows.Close()
	if len(reserved) == 0 || count != 0 {
		return nil, fmt.Errorf("can't reserve %d good: not enough goods on available storages", uniqId)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("can't commit reserve transaction: %w", err)
	}
	return reserved, nil
}

func (d *Database) ReleaseGood(ctx context.Context) {

}
