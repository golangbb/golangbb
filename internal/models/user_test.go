package models

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golangbb/golangbb/v2/internal/database"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"regexp"
)

var _ = Describe("database", func() {
	var mock sqlmock.Sqlmock
	var db *sql.DB

	BeforeEach(func() {
		sqlDb, sqlMock, err := sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())

		mock = sqlMock
		db = sqlDb

		gormDB, err := database.Connect(sqlite.Dialector{
			DriverName: "sqlite",
			Conn:       db,
		}, gorm.Config{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(gormDB).ShouldNot(BeNil())
		Expect(gormDB.DB()).Should(BeIdenticalTo(db))
	})
	AfterEach(func() {
		db.Close()
	})

	Context("CreateUser", func() {
		When("inserting a user with only UserName, DisplayName and Password", func() {
			It("should execute insert all 3 fields and not error", func() {
				user := &User{
					UserName:    "MotherOfDragons",
					DisplayName: "Mother Of Dragons",
					Password:    "password",
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`user_name`,`display_name`,`password`) VALUES (?,?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.UserName, user.DisplayName, user.Password).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := CreateUser(user)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a user without UserName", func() {
			It("should return an error without executing any sql on database", func() {
				user := &User{
					DisplayName: "Mother Of Dragons",
					Password:    "password",
				}

				err := CreateUser(user)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyUserName))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a user with only UserName and Password", func() {
			It("should insert the UserName value as DisplayName", func() {
				user := &User{
					UserName: "MotherOfDragons",
					Password: "password",
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`user_name`,`display_name`,`password`) VALUES (?,?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.UserName, user.UserName, user.Password).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := CreateUser(user)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a user without Password", func() {
			It("should return an error without executing any sql on database", func() {
				user := &User{
					UserName:    "MotherOfDragons",
					DisplayName: "Mother Of Dragons",
				}

				err := CreateUser(user)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyPassword))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
