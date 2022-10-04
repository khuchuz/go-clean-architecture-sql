package repository

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/khuchuz/go-clean-architecture-sql/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB                *gorm.DB
	mock              sqlmock.Sqlmock
	userRepositorySQL *UserRepositorySQL
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	assert.NoError(s.T(), err, "Failed to open gorm DB")
	assert.NotNil(s.T(), db, "Mock DB is null")
	assert.NotNil(s.T(), s.mock, "SQLMock is null")

	s.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(s.T(), err)
	s.userRepositorySQL = InitUserRepositorySQL(s.DB)
	//defer db.Close()
}

func (s *Suite) TestSQLCreateUser() {
	user := &models.User{
		Username: "dummy1",
		Email:    "dummy1@gmail.com",
		Password: hashThis("Dummy123"),
	}

	s.mock.ExpectBegin() // start transaction
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`username`,`email`,`password`) VALUES (?,?,?)")).
		WithArgs(user.Username, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit() // commit transaction

	err := s.userRepositorySQL.SQLCreateUser(user)
	s.NoError(err)
}

func (s *Suite) TestSQLGetUser() {
	var (
		id       = 10
		username = "dummy10"
		email    = "akun10@email.com"
		password = hashThis("Password10")
	)

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? AND password = ?")).
		WithArgs(username, password).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(id, username, email, password))

	res, err := s.userRepositorySQL.SQLGetUser(username, password)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&models.User{ID: res.ID, Username: username, Email: email, Password: password}, res))
}

func (s *Suite) TestSQLIsUserExistByUsername_True() {
	var (
		id       = 11
		username = "dummy11"
		email    = "akun11@email.com"
		password = hashThis("Password11")
	)

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ?")).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(id, username, email, password))

	res := s.userRepositorySQL.SQLIsUserExistByUsername(username)
	require.Equal(s.T(), true, res)
}

func (s *Suite) TestSQLIsUserExistByUsername_False() {
	var (
		id       = 11
		username = "dummy11"
		email    = "akun11@email.com"
		password = hashThis("Password11")
	)

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ?")).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(id, username, email, password))

	res := s.userRepositorySQL.SQLIsUserExistByUsername(username)
	require.Equal(s.T(), true, res)
}

func (s *Suite) TestSQLIsUserExistByEmail_True() {
	var (
		id       = 12
		username = "dummy12"
		email    = "akun12@email.com"
		password = hashThis("Password12")
	)

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ?")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(id, username, email, password))

	res := s.userRepositorySQL.SQLIsUserExistByEmail(email)
	require.Equal(s.T(), true, res)
}

func (s *Suite) TestSQLUpdatePassword() {
	var (
		username = "dummy12"
		password = hashThis("Password12")
	)
	s.mock.ExpectBegin() // start transaction
	s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `password`=? WHERE username = ?")).WithArgs(password, username).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit() // commit transaction
	err := s.userRepositorySQL.SQLUpdatePassword(username, password)
	s.NoError(err)
}

func hashThis(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte("hash_salt"))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
func TestSuiteRepository(t *testing.T) {
	suite.Run(t, new(Suite))
}
