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

var _ = Describe("Topic", func() {
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

	Context("CreateTopic", func() {
		When("inserting a Topic with a valid AuthorID and Title", func() {
			It("should insert a new Topic record", func() {
				topic := &Topic{
					AuthorID: 10,
					Title:     "Marvel",
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `topics` (`created_at`,`updated_at`,`deleted_at`,`title`,`parent_id`,`author_id`) VALUES (?,?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, topic.Title, nil, topic.AuthorID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := CreateTopic(topic)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Topic with a valid AuthorID, Title and ParentID", func() {
			It("should insert a new Topic record", func() {
				parentID := uint(20)
				topic := &Topic{
					AuthorID: 10,
					Title:     "Marvel",
					ParentID: &parentID,
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `topics` (`created_at`,`updated_at`,`deleted_at`,`title`,`parent_id`,`author_id`) VALUES (?,?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, topic.Title, topic.ParentID, topic.AuthorID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := CreateTopic(topic)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Topic without a AuthorID", func() {
			It("should not attempt to insert a new Topic record", func() {
				topic := &Topic{
					Title:     "Marvel",
				}

				err := CreateTopic(topic)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyUserID))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Topic without a Title", func() {
			It("should not attempt to insert a new Topic record", func() {
				topic := &Topic{
					AuthorID: 10,
				}

				err := CreateTopic(topic)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyTitle))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Topic that errors", func() {
			It("should rollback transaction and return error", func() {
				topic := &Topic{
					AuthorID: 10,
					Title:     "Marvel",
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `topics` (`created_at`,`updated_at`,`deleted_at`,`title`,`parent_id`,`author_id`) VALUES (?,?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, topic.Title, nil, topic.AuthorID).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()

				err := CreateTopic(topic)
				Expect(err).Should(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Topic with an Author", func() {
			It("should not attempt to insert a new User/Author record", func() {
				userID := uint(10)
				topic := &Topic{
					AuthorID: userID,
					Title:     "Marvel",
					Author: User{
						Model: gorm.Model{ID: userID},
					},
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `topics` (`created_at`,`updated_at`,`deleted_at`,`title`,`parent_id`,`author_id`) VALUES (?,?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, topic.Title, nil, topic.AuthorID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := CreateTopic(topic)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Topic with a Parent Topic", func() {
			It("should not attempt to insert a new Topic record", func() {
				topic := &Topic{
					AuthorID: 10,
					Title:     "Marvel",
					Parent: &Topic{
						Model: gorm.Model{ID: 1},
					},
				}

				mock.ExpectBegin()
				sql := regexp.QuoteMeta("INSERT INTO `topics` (`created_at`,`updated_at`,`deleted_at`,`title`,`parent_id`,`author_id`) VALUES (?,?,?,?,?,?)")
				mock.ExpectExec(sql).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, topic.Title, nil, topic.AuthorID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := CreateTopic(topic)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})

