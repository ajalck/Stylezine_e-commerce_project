package repository

import (
	"ajalck/e_commerce/domain"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// CREATE
// type MockDB struct {
// 	mock.Mock
// 	gorm.DB
// }

//	func (m *MockDB) Create(value interface{}) *gorm.DB {
//		args := m.Called(value)
//		return args.Get(0).(*gorm.DB)
//	}
func TestCreateAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()
	// mockDB := new(MockDB)
	// mockDB := mock.Mock{}

	if err != nil {
		t.Fatalf("Error creating mock database connection :%v", err)
	}
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm instance:%v", err)
	}

	mockUser := domain.User{
		User_ID:    "abc",
		First_Name: "Ajal",
		Last_Name:  "CK",
		Photo:      "nil",
		Email:      "abc@gmail.com",
		Gender:     "m",
		Phone:      "9939443874",
		Password:   "12345",
	}
	//CREATE GORM MOCK
	// mockDB.On("CreateAdmin", mockUser).Return(&gorm.DB{}).Times(1)

	expID := 1
	mockquery := "INSERT INTO \"users\" \\(\"user_id\",\"first_name\",\"last_name\",\"photo\",\"email\",\"gender\",\"phone\",\"password\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6,\\$7,\\$8) RETURNING \"id\";)"
	mock.ExpectBegin()
	mock.ExpectQuery(mockquery).
		WithArgs(mockUser.User_ID, mockUser.First_Name, mockUser.Last_Name, mockUser.Photo, mockUser.Email, mockUser.Gender, mockUser.Phone, mockUser.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expID))
	mock.ExpectCommit()

	repo := &AdminRepo{DB: gormDB}

	err = repo.CreateAdmin(mockUser)

	fmt.Println("------------------error is ---------", err)
	mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("mock expectations were not met:%v", err)
	}
}

//READ

func TestViewCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB :%v", err)
	}
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm instance:%v", err)
	}

	mockCategory := domain.Category{
		Category_ID:   1,
		Category_name: "test",
	}
	mock.ExpectQuery("SELECT \\* FROM \"categories\" WHERE Category_ID=\\$1 ORDER BY \"categories\"\\.\"category_id\" LIMIT 1").WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "category_name"}).AddRow(mockCategory.Category_ID, mockCategory.Category_name))

	repo := &AdminRepo{DB: gormDB}

	result, err := repo.ViewCategory(mockCategory)
	mockerr := mock.ExpectationsWereMet()
	if mockerr != nil {
		t.Errorf("TEST FAILED : exp.result:%v   act.result:%v\nexp.error:%v   act.error:%v", mockCategory, result, nil, err)
		t.Fatalf("Mock expressions were not met :%v", mockerr)
	}
}

//UPDATE

func TestBlockUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error in creating mock DB :%v", err)
	}
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm instance:%v", err)
	}

	mockQuery := "UPDATE users SET status=\\$1 WHERE user_id=\\$2;"

	mock.ExpectQuery(mockQuery).WithArgs("blocked", "abcd")

	repo := &AdminRepo{DB: gormDB}

	repo.BlockUser("abcd")
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("mock expectations were not met:%v", err)
	}
}
