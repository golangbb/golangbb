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

var _ = Describe("Discussion", func() {
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

	Context("CreateDiscussion", func() {
		When("inserting a Discussion with a Post", func() {
			It("should insert a new Discussion record and Post Record with the Author for the Discussion being the same as the Author for the Post", func() {
				discussion := &Discussion{
					AuthorID: 10,
					Title:    "Marvel vs DC",
					TopicID: 20,
					Posts: []Post{
						{
							Content: "some content",
						},
					},
				}

				newDiscussionID := int64(1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `discussions` (`created_at`,`updated_at`,`deleted_at`,`title`,`author_id`,`topic_id`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, discussion.Title, discussion.AuthorID, discussion.TopicID).
					WillReturnResult(sqlmock.NewResult(newDiscussionID, 1))
				mock.ExpectExec("SAVEPOINT .").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `posts` (`created_at`,`updated_at`,`deleted_at`,`content`,`author_id`,`discussion_id`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, discussion.Posts[0].Content, discussion.AuthorID, newDiscussionID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := CreateDiscussion(discussion)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Discussion with more than one Post", func() {
			It("should not attempt to insert a new Discussion record and return an error", func() {
				discussion := &Discussion{
					AuthorID: 10,
					Title:    "Marvel vs DC",
					TopicID: 20,
					Posts: []Post{
						{
							Content: "some content",
						},
						{
							Content: "more content",
						},
					},
				}

				err := CreateDiscussion(discussion)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrDiscussionWithoutSinglePost))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Discussion without a Post", func() {
			It("should not attempt to insert a new Discussion record and return an error", func() {
				discussion := &Discussion{
					AuthorID: 10,
					Title:    "Marvel vs DC",
					TopicID: 20,
				}

				err := CreateDiscussion(discussion)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrDiscussionWithoutSinglePost))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Discussion without a Title", func() {
			It("should not attempt to insert a new Discussion record and return an error", func() {
				discussion := &Discussion{
					AuthorID: 10,
					TopicID: 20,
				}

				err := CreateDiscussion(discussion)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyTitle))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Discussion without an AuthorID", func() {
			It("should not attempt to insert a new Discussion record and return an error", func() {
				discussion := &Discussion{
					Title:    "Marvel vs DC",
					TopicID: 20,
				}

				err := CreateDiscussion(discussion)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyUserID))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Discussion without a TopicID", func() {
			It("should not attempt to insert a new Discussion record and return an error", func() {
				discussion := &Discussion{
					AuthorID: 10,
					Title:    "Marvel vs DC",
				}

				err := CreateDiscussion(discussion)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyTopicID))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Discussion with a Post", func() {
			It("should insert a new Discussion record and Post Record with the Author for the Discussion being the same as the Author for the Post", func() {
				discussion := &Discussion{
					AuthorID: 10,
					Title:    "Marvel vs DC",
					TopicID: 20,
					Posts: []Post{
						{
							Content: "some content",
						},
					},
				}

				newDiscussionID := int64(1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `discussions` (`created_at`,`updated_at`,`deleted_at`,`title`,`author_id`,`topic_id`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, discussion.Title, discussion.AuthorID, discussion.TopicID).
					WillReturnResult(sqlmock.NewResult(newDiscussionID, 1))
				mock.ExpectExec("SAVEPOINT .").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `posts` (`created_at`,`updated_at`,`deleted_at`,`content`,`author_id`,`discussion_id`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, discussion.Posts[0].Content, discussion.AuthorID, newDiscussionID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := CreateDiscussion(discussion)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Topic that errors", func() {
			It("should rollback transaction and return error", func() {
				topic := &Topic{
					AuthorID: 10,
					Title:    "Marvel",
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
					Title:    "Marvel",
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
					Title:    "Marvel",
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
