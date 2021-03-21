package models

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golangbb/golangbb/v2/internal/database"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"regexp"
)

var _ = Describe("Email", func() {
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

	Context("CreateEmail", func() {
		When("inserting an Email with a valid UserID", func() {
			It("should insert a new Email record", func() {
				email := &Email{
					UserID: 10,
					Email: "ironman@mcu.com",
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `emails` (`email`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(email.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, email.UserID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := CreateEmail(email)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting an Email without a UserID", func() {
			It("should not attempt to insert a new Email record", func() {
				email := &Email{
					Email: "ironman@mcu.com",
				}

				err := CreateEmail(email)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyUserID))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting an Email that errors", func() {
			It("should rollback transaction and return error", func() {
				email := &Email{
					UserID: 10,
					Email: "ironman@mcu.com",
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `emails` (`email`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(email.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, email.UserID).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()

				err := CreateEmail(email)
				Expect(err).Should(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting an Email with a User", func() {
			It("should not attempt to insert a new user record", func() {
				userID := uint(10)
				email := &Email{
					Email: "ironman@mcu.com",
					UserID: userID,
					User: User{
						Model:       gorm.Model{
							ID:        userID,
						},
					},
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `emails` (`email`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(email.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, email.UserID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := CreateEmail(email)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
