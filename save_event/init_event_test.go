package saveevent

import (
	"sort"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInitEventsEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "key", "name", "description"})
	mock.ExpectQuery("^SELECT (.+) FROM events_system where `key`=?").WithArgs().WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO events_system").WithArgs("blzone", "Блокировка опасной зоны").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	type_events := map[string]string{
		"blzone.state":    "Статус блокировки",
		"blzone.mode":     "Режим работы",
		"blzone.error":    "Неисправность",
		"blzone.evreg_99": "Обнаружение тага в опасной зоне поля 8 кГц",
	}

	var keys []string
	for k := range type_events {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var id int64
	for _, key := range keys {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO events_typeevent").WithArgs(key, type_events[key], "", 0, "", 1).WillReturnResult(sqlmock.NewResult(id, 1))
		mock.ExpectCommit()
		id++
	}

	err = InitEvents(db, "blzone", "Блокировка опасной зоны", type_events)

	if err != nil {
		t.Error("")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestInitEventsUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "key", "name", "description"}).
		AddRow(14, "blzone", "asdfaf", "")
	mock.ExpectQuery("^SELECT (.+) FROM events_system where `key`=?").WithArgs().WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE events_system").WithArgs("Блокировка опасной зоны", 14).WillReturnResult(sqlmock.NewResult(int64(14), 1))
	mock.ExpectCommit()

	type_events := map[string]string{
		"blzone.state":    "Статус блокировки",
		"blzone.mode":     "Режим работы",
		"blzone.error":    "Неисправность",
		"blzone.evreg_99": "Обнаружение тага в опасной зоне поля 8 кГц",
	}

	var keys []string
	for k := range type_events {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var id int64 = 1
	for _, key := range keys {
		rows := sqlmock.NewRows([]string{"id", "key", "name", "description", "type"}).
			AddRow(id, key, "asdfasdf", "", 0)
		mock.ExpectQuery("^SELECT (.+) FROM events_typeevent").WithArgs(key, 14).WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE events_typeevent").WithArgs(type_events[key], id).WillReturnResult(sqlmock.NewResult(int64(14), 1))
		mock.ExpectCommit()
		id++
	}

	err = InitEvents(db, "blzone", "Блокировка опасной зоны", type_events)
	if err != nil {
		t.Error("")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
