package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

//CREATE

func TestFindUser(t *testing.T){

	///creating mock db and mock query
	db,mock,err:=sqlmock.New()
	if err!=nil{
		t.Fatalf("Error in creating mock DB :%v",err)
	}
	defer db.Close()
	mockQuery:="SELECT * FROM users WHERE email=\\$1 and user_role=\\$2"

	mock.ExpectQuery(mockQuery).WithArgs("abcd","user").WillReturnRows(sqlmock.NewRows([]string{}))
}