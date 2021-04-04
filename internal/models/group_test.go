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

var _ = Describe("Group", func() {
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

	Context("CreateGroup", func() {
		When("inserting a Group with a valid AuthorID and Name", func() {
			It("should insert a new Group record", func() {
				group := &Group{
					AuthorID: 10,
					Name:     "The Avengers",
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `groups` (`created_at`,`updated_at`,`deleted_at`,`name`,`author_id`) VALUES (?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, group.Name, group.AuthorID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := CreateGroup(group)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Group without a AuthorID", func() {
			It("should not attempt to insert a new Group record", func() {
				group := &Group{
					Name: "The Avengers",
				}

				err := CreateGroup(group)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyUserID))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Group without a Name", func() {
			It("should not attempt to insert a new Group record", func() {
				group := &Group{
					AuthorID: 10,
				}

				err := CreateGroup(group)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyName))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Group that errors", func() {
			It("should rollback transaction and return error", func() {
				group := &Group{
					AuthorID: 10,
					Name:     "The Avengers",
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `groups` (`created_at`,`updated_at`,`deleted_at`,`name`,`author_id`) VALUES (?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, group.Name, group.AuthorID).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()

				err := CreateGroup(group)
				Expect(err).Should(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Group with an Author", func() {
			It("should not attempt to insert a new User/Author record", func() {
				userID := uint(10)
				group := &Group{
					AuthorID: userID,
					Name:     "The Avengers",
					Author: User{
						Model: gorm.Model{ID: userID},
					},
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `groups` (`created_at`,`updated_at`,`deleted_at`,`name`,`author_id`) VALUES (?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, group.Name, group.AuthorID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := CreateGroup(group)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Group with User(s)", func() {
			It("should not attempt to insert a new User record", func() {
				userID := uint(10)
				group := &Group{
					AuthorID: 10,
					Name:     "The Avengers",
					Users: []User{
						{
							Model: gorm.Model{ID: userID},
						},
					},
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `groups` (`created_at`,`updated_at`,`deleted_at`,`name`,`author_id`) VALUES (?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, group.Name, group.AuthorID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := CreateGroup(group)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
