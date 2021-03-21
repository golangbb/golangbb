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

var _ = Describe("User", func() {
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
		When("inserting a User with only UserName, DisplayName and Password", func() {
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

		When("inserting a User with Groups", func() {
			It("should not attempt to insert Groups", func() {
				user := &User{
					UserName:    "MotherOfDragons",
					DisplayName: "Mother Of Dragons",
					Password:    "password",
					Groups: []Group{
						{
							Name: "some group",
						},
					},
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

		When("inserting a User without UserName", func() {
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

		When("inserting a User with only UserName and Password", func() {
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

		When("inserting a User without Password", func() {
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

		When("inserting a User and the insert fails", func() {
			It("should rollback the whole transaction", func() {
				user := &User{
					UserName:    "MotherOfDragons",
					DisplayName: "Mother Of Dragons",
					Password:    "password",
				}

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`user_name`,`display_name`,`password`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.UserName, user.DisplayName, user.Password).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()

				err := CreateUser(user)
				Expect(err).Should(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a User with an Email", func() {
			It("should insert a single Email record related to the inserted User record", func() {
				user := &User{
					UserName:    "MotherOfDragons",
					DisplayName: "Mother Of Dragons",
					Password:    "password",
					Emails: []*Email{
						{
							Email: "dany@got.com",
						},
					},
				}

				newUserID := int64(1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`user_name`,`display_name`,`password`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.UserName, user.DisplayName, user.Password).
					WillReturnResult(sqlmock.NewResult(newUserID, 1))
				mock.ExpectExec("SAVEPOINT .").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `emails` (`email`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?)")).
					WithArgs(user.Emails[0].Email, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newUserID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := CreateUser(user)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a User with multiple Emails", func() {
			It("should insert a multiple Email records related to the inserted User record", func() {
				user := &User{
					UserName:    "MotherOfDragons",
					DisplayName: "Mother Of Dragons",
					Password:    "password",
					Emails: []*Email{
						{
							Email: "dany@got.com",
						},
						{
							Email: "sarahconner@terminator.com",
						},
					},
				}

				newUserID := int64(1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`user_name`,`display_name`,`password`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.UserName, user.DisplayName, user.Password).
					WillReturnResult(sqlmock.NewResult(newUserID, 1))
				mock.ExpectExec("SAVEPOINT .").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `emails` (`email`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?),(?,?,?,?,?)")).
					WithArgs(user.Emails[0].Email, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newUserID, user.Emails[1].Email, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newUserID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := CreateUser(user)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("inserting a User with an Email where the Email insert fails", func() {
			It("should rollback the whole transaction", func() {
				user := &User{
					UserName:    "MotherOfDragons",
					DisplayName: "Mother Of Dragons",
					Password:    "password",
					Emails: []*Email{
						{
							Email: "dany@got.com",
						},
					},
				}

				newUserID := int64(1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`user_name`,`display_name`,`password`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.UserName, user.DisplayName, user.Password).
					WillReturnResult(sqlmock.NewResult(newUserID, 1))
				mock.ExpectExec("SAVEPOINT .").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `emails` (`email`,`created_at`,`updated_at`,`deleted_at`,`user_id`) VALUES (?,?,?,?,?)")).
					WithArgs(user.Emails[0].Email, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newUserID).
					WillReturnError(errors.New("some db insert error"))
				mock.ExpectExec("ROLLBACK TO SAVEPOINT .").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()

				err := CreateUser(user)
				Expect(err).Should(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
