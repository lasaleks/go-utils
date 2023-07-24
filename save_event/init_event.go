package saveevent

import (
	"database/sql"
	"sort"
)

func insert(db *sql.DB, sql_insert string, args ...interface{}) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()
	var result sql.Result
	// var err error
	if result, err = tx.Exec(sql_insert, args...); err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func select_one_row(db *sql.DB, sql string, args []interface{}, result ...interface{}) error {
	row := db.QueryRow(sql, args...)
	err := row.Scan(result...)
	if err != nil {
		return err
	}
	return nil
}

func update(db *sql.DB, sql string, args ...interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _, err := tx.Exec(sql, args...); err != nil {
		return err
	}
	return nil
}

func InitEvents(db *sql.DB, system_key string, system_name string, type_events map[string]string) error {
	event_system := struct {
		id          int64
		key         string
		name        string
		description string
	}{}
	err := select_one_row(db, "SELECT id, `key`, name, description FROM events_system where `key`=?", []interface{}{"blzone"}, &event_system.id, &event_system.key, &event_system.name, &event_system.description)
	if err != nil {
		event_system.key = "blzone"
		event_system.name = "Блокировка опасной зоны"
		id, err := insert(db, "INSERT INTO events_system(`key`, name) VALUES (?, ?)", event_system.key, event_system.name)
		if err != nil {
			return err
		}
		event_system.id = id
	} else {
		err = update(db, "UPDATE events_system SET name=? WHERE id=?", system_name, event_system.id)
		if err != nil {
			return err
		}
	}
	type TypeEvent struct {
		id          int64
		key         string
		name        string
		description string
		typeev      int64
	}

	var keys []string
	for k := range type_events {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		type_event := TypeEvent{}
		err := select_one_row(
			db,
			"SELECT id, `key`, name, description, `type` FROM events_typeevent WHERE `key`=? and system_id=?",
			[]interface{}{key, event_system.id},
			&type_event.id, &type_event.key, &type_event.name, &type_event.description, &type_event.typeev,
		)
		if err != nil {
			type_event.key = key
			type_event.name = type_events[key]
			_, err := insert(db, "INSERT INTO events_typeevent(`key`, name, description, type, category, system_id) VALUES (?, ?, ?, ?, ?, ?)", type_event.key, type_event.name, type_event.description, type_event.typeev, "", event_system.id)
			if err != nil {
				return err
			}
		} else {
			err = update(db, "UPDATE events_typeevent SET name=? WHERE id=?", type_events[key], type_event.id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
