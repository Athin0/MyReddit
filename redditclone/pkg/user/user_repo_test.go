package user

import (
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
)

// cd

func TestAuthorize(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID int64 = 0
	login, password := "Athin", "asdfghjk"

	// good query
	rows := sqlmock.NewRows([]string{"id", "login", "password"})
	expect := []*User{
		{elemID, login, password},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Login, item.Password)
	}

	mock.
		ExpectQuery("SELECT id, login, password FROM  users WHERE").
		WithArgs(login).
		WillReturnRows(rows)

	repo := &UserMysqlRepository{
		DB: db,
	}

	item, err := repo.Authorize(login, password)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// BabPass error
	badPass := "lol"
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Login, item.Password)
	}
	mock.
		ExpectQuery("SELECT id, login, password FROM  users WHERE").
		WithArgs(login).
		WillReturnRows(rows)

	_, err = repo.Authorize(login, badPass)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err != ErrBadPass {
		t.Errorf("expected error: %v, got: %v", ErrBadPass, err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery("SELECT id, login, password FROM  users WHERE").
		WithArgs(login).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Authorize(login, password)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "title")

	mock.
		ExpectQuery("SELECT id, login, password FROM  users WHERE").
		WithArgs(login).
		WillReturnRows(rows)

	_, err = repo.Authorize(login, password)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &UserMysqlRepository{
		DB: db,
	}

	login := "Athin"
	pass := "asdfghjk"

	//ok query
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(login, pass).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.AddUserInRepo(login, pass)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if id.ID != 1 {
		t.Errorf("bad id: want %v, have %v", id, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(login, pass).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.AddUserInRepo(login, pass)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// result error
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(login, pass).
		WillReturnError(fmt.Errorf("bad_result"))

	_, err = repo.AddUserInRepo(login, pass)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
