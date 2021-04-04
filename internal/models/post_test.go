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

var _ = Describe("Post", func() {
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

	Context("CreatePost", func() {
		When("inserting a Post", func() {
			It("should insert a new Post record with an Author an a Discussion", func() {
				post := &Post{
					AuthorID:     10,
					DiscussionID: 5,
					Content:      "Marvel rules, DC drools",
				}

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `posts` (`created_at`,`updated_at`,`deleted_at`,`content`,`author_id`,`discussion_id`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, post.Content, post.AuthorID, post.DiscussionID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := CreatePost(post)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Post without an AuthorID", func() {
			It("should not attempt to insert a new Post record and return an error", func() {
				post := &Post{
					DiscussionID: 5,
					Content:      "Marvel rules, DC drools",
				}

				err := CreatePost(post)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyUserID))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Post without a DiscussionID", func() {
			It("should not attempt to insert a new Post record and return an error", func() {
				post := &Post{
					AuthorID: 10,
					Content:  "Marvel rules, DC drools",
				}

				err := CreatePost(post)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyDiscussionID))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Post without Content", func() {
			It("should not attempt to insert a new Post record and return an error", func() {
				post := &Post{
					AuthorID:     10,
					DiscussionID: 5,
				}

				err := CreatePost(post)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(Equal(ErrEmptyContent))

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Post that errors", func() {
			It("should rollback transaction and return error", func() {
				post := &Post{
					AuthorID:     10,
					DiscussionID: 5,
					Content:      "Marvel rules, DC drools",
				}

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `posts` (`created_at`,`updated_at`,`deleted_at`,`content`,`author_id`,`discussion_id`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, post.Content, post.AuthorID, post.DiscussionID).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()

				err := CreatePost(post)
				Expect(err).Should(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Post with a Discussion", func() {
			It("should not attempt to insert a new Discussion record", func() {
				discussionID := uint(10)
				post := &Post{
					AuthorID:     10,
					DiscussionID: 5,
					Content:      "Marvel rules, DC drools",
					Discussion: Discussion{
						Model: gorm.Model{ID: discussionID},
					},
				}

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `posts` (`created_at`,`updated_at`,`deleted_at`,`content`,`author_id`,`discussion_id`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, post.Content, post.AuthorID, post.DiscussionID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := CreatePost(post)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("inserting a Post with an Author", func() {
			It("should not attempt to insert a new Author record", func() {
				userID := uint(10)
				post := &Post{
					AuthorID:     userID,
					DiscussionID: 5,
					Content:      "Marvel rules, DC drools",
					Author: User{
						Model: gorm.Model{ID: userID},
					},
				}

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `posts` (`created_at`,`updated_at`,`deleted_at`,`content`,`author_id`,`discussion_id`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, post.Content, post.AuthorID, post.DiscussionID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				err := CreatePost(post)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
