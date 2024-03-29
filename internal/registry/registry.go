package registry

import (
	"LamodaTest/internal/entity/goods"
	"LamodaTest/internal/entity/storages"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Db interface {
	Storages(ctx context.Context, all bool) ([]storages.Storage, error)
	StoragesAdd(ctx context.Context, name string, available bool) (int64, error)
	StoragesDelete(ctx context.Context, id int) (int64, error)
	StoragesChangeAccess(ctx context.Context, id int, available bool) (int64, error)
	Goods(ctx context.Context) ([]goods.Good, error)
	AvailableGoods(ctx context.Context) (map[int]goods.RemainsDTO, error)
	ReserveGood(ctx context.Context, uniqId int, count int) (map[int]int, error)
	ReleaseGood(ctx context.Context, uniqId int, count int) error
	GoodAdd(ctx context.Context, name string, size string, uniqCode int) (int64, error)
	GoodDelete(ctx context.Context, uniqCode int) (int64, error)
}

type Database struct {
	conn *sql.DB
}

func New(connect *sql.DB) *Database {
	return &Database{conn: connect}
}

func (d *Database) Storages(ctx context.Context, all bool) ([]storages.Storage, error) {
	query := "select * from storages"
	if !all {
		query = fmt.Sprintf("%s where available = 1", query)
	}
	cmd, err := d.conn.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("can't prepare sql: %w", err)
	}
	rows, err := cmd.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't scan from storage list: %w", err)
	}
	defer rows.Close()
	var result []storages.Storage
	for rows.Next() {
		values := storages.Storage{}
		err = rows.Scan(&values.ID, &values.Name, &values.RawAvailable)
		if err != nil {
			return nil, fmt.Errorf("can't scan from storage list: %s", err.Error())
		}
		values.Available, err = strconv.ParseBool(values.RawAvailable)
		if err != nil {
			return nil, fmt.Errorf("can't parse bool from %s str: %w", values.RawAvailable, err)
		}
		result = append(result, values)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error when try get all storages: %v", err)
	}
	return result, nil
}

func (d *Database) StoragesAdd(ctx context.Context, name string, available bool) (int64, error) {
	result, err := d.conn.ExecContext(ctx, "insert into storages (name, available) values (?, ?)",
		name, available)
	if err != nil {
		return -1, fmt.Errorf("can't add storage [%s, %t]: %w", name, available, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("can't get last added storage id from database: %w", err)
	}
	return id, nil
}

func (d *Database) StoragesDelete(ctx context.Context, id int) (int64, error) {
	result, err := d.conn.ExecContext(ctx, "delete from storages where id = ?", id)
	if err != nil {
		return -1, fmt.Errorf("can't delete storage with id %d: %w", id, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("can't get row affected after delete storage: %w", err)
	}
	return affected, nil
}

func (d *Database) StoragesChangeAccess(ctx context.Context, id int, available bool) (int64, error) {
	result, err := d.conn.ExecContext(ctx, "update storages set available = ? where id = ?", available, id)
	if err != nil {
		return -1, fmt.Errorf("can't change storage with id %d: %w", id, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("can't get row affected after change storage: %w", err)
	}
	return affected, nil
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
	if err != nil {
		return nil, fmt.Errorf("can't prepare sql: %w", err)
	}
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
		return nil, fmt.Errorf("error when try get available goods: %v", err)
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
	rows, err := d.conn.QueryContext(ctx, `SELECT 
			remains.id, 
			remains.storage_id, 
			remains.count - remains.reserved AS avail 
		from remains 
		JOIN storages ON storages.id = remains.storage_id 
		where good_id = ? AND storages.available = 1`, id)
	if err != nil {
		return nil, fmt.Errorf("can't request avail goods for release: %w", err)
	}
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

func (d *Database) ReleaseGood(ctx context.Context, uniqId int, count int) error {
	tx, err := d.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable}) //
	if err != nil {
		return fmt.Errorf("can't init transaction: %w", err)
	}
	defer tx.Rollback()
	var id int
	if err = tx.QueryRowContext(ctx, "SELECT id from goods where uniq_code = ?",
		uniqId).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("can't found good with uniq_code %d", uniqId)
		}
		return err
	}
	rows, err := d.conn.QueryContext(ctx, `SELECT 
			remains.id, 
			remains.storage_id, 
			remains.reserved
		from remains 
		JOIN storages ON storages.id = remains.storage_id 
		where good_id = ? AND storages.available = 1`, id)
	if err != nil {
		return fmt.Errorf("can't request avail goods for release: %w", err)
	}
	for rows.Next() {
		var tmp struct {
			Id        int
			storageId int
			Reserved  int
		}
		err = rows.Scan(&tmp.Id, &tmp.storageId, &tmp.Reserved)
		if err != nil {
			return fmt.Errorf("can't get release good with id %d: %w", id, err)
		}
		if tmp.Reserved < count {
			_, err = tx.ExecContext(ctx, "UPDATE remains SET reserved = reserved - ? WHERE id = ?",
				tmp.Reserved, tmp.Id)
			if err != nil {
				return fmt.Errorf("can't update remains note with id %d: %w", tmp.Id, err)
			}
			count = count - tmp.Reserved
		} else {
			_, err = tx.ExecContext(ctx, "UPDATE remains SET reserved = reserved - ? WHERE id = ?",
				count, tmp.Id)
			if err != nil {
				return fmt.Errorf("can't update remains note with id %d: %w", tmp.Id, err)
			}
			count = 0
		}
	}
	defer rows.Close()
	if count != 0 {
		return fmt.Errorf("can't release good with id %d: not enough reserved goods on available storages", uniqId)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("can't commit release transaction: %w", err)
	}
	return nil
}

func (d *Database) Goods(ctx context.Context) ([]goods.Good, error) {
	cmd, err := d.conn.Prepare("select * from goods;")
	if err != nil {
		return nil, fmt.Errorf("can't prepare sql: %w", err)
	}
	rows, err := cmd.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't scan from goods list: %w", err)
	}
	result := []goods.Good{}
	defer rows.Close()
	for rows.Next() {
		values := goods.Good{}
		err = rows.Scan(&values.Id, &values.Name, &values.Size, &values.UniqCode)
		if err != nil {
			return nil, fmt.Errorf("can't scan from goods list: %w", err)
		}
		result = append(result, values)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error when try get all goods: %w", err)
	}
	return result, nil
}

func (d *Database) GoodAdd(ctx context.Context, name string, size string, uniqCode int) (int64, error) {
	result, err := d.conn.ExecContext(ctx, "insert into goods (name, size, uniq_code) values (?, ?, ?)",
		name, size, uniqCode)
	if err != nil {
		return -1, fmt.Errorf("can't add good [%s, %s, %d]: %w", name, size, uniqCode, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("can't get last added good id from database: %w", err)
	}
	return id, nil
}

func (d *Database) GoodDelete(ctx context.Context, uniqCode int) (int64, error) {
	result, err := d.conn.ExecContext(ctx, "delete from goods where uniq_code = ?", uniqCode)
	if err != nil {
		return -1, fmt.Errorf("can't delete good with uniq_code %d: %w", uniqCode, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("can't get row affected after delete good: %w", err)
	}
	return affected, nil
}
