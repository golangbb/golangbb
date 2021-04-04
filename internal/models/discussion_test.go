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

		When("inserting a Discussion that errors", func() {
			It("should rollback transaction and return error", func() {
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

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `discussions` (`created_at`,`updated_at`,`deleted_at`,`title`,`author_id`,`topic_id`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, discussion.Title, discussion.AuthorID, discussion.TopicID).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()

				err := CreateDiscussion(discussion)
				Expect(err).Should(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Post, for the Discussion, that errors ", func() {
			It("should rollback transaction and return error", func() {
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
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()

				err := CreateDiscussion(discussion)
				Expect(err).Should(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Discussion with a Post and a Author", func() {
			It("should not attempt to insert a new Author", func() {
				userID := uint(10)
				discussion := &Discussion{
					Author: User{
						Model: gorm.Model{ID: userID},
					},
					AuthorID: userID,
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

		When("inserting a Discussion with a Post and Topic", func() {
			It("should not attempt to insert a new Topic", func() {
				topicID := uint(20)
				discussion := &Discussion{
					AuthorID: 10,
					Title:    "Marvel vs DC",
					TopicID: topicID,
					Topic: Topic{Model: gorm.Model{ID: topicID}},
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

		When("inserting a Discussion with a Post with an Author", func() {
			It("should not attempt to insert a new Author", func() {
				userID := uint(10)
				discussion := &Discussion{
					AuthorID: userID,
					Title:    "Marvel vs DC",
					TopicID: 20,
					Posts: []Post{
						{
							Content: "some content",
							Author: User{
								Model: gorm.Model{ID: userID},
							},
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

		When("inserting a Discussion with a Post with a Discussion", func() {
			It("should not attempt to insert a new Discussion", func() {
				discussionID := uint(10)
				discussion := &Discussion{
					AuthorID: 10,
					Title:    "Marvel vs DC",
					TopicID: 20,
					Posts: []Post{
						{
							Content: "some content",
							Discussion: Discussion{Model: gorm.Model{ID: discussionID}},
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
	})
})
