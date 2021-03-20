package database

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golangbb/golangbb/v2/internal/models"
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
	})
	AfterEach(func() {
		db.Close()
	})

	Context("Connect", func() {
		When("connecting to the database", func() {
			It("should return a gorm.DB instance", func() {
				gormDB, err := Connect(sqlite.Dialector{
					DriverName: "sqlite",
					Conn:       db,
				}, gorm.Config{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(gormDB).ShouldNot(BeNil())
				Expect(gormDB.DB()).Should(BeIdenticalTo(db))
			})
		})

		When("expecting an error while connecting to the database", func() {
			It("should return an error and the gorm.DB instance should be nil", func() {
				gormDB, err := Connect(sqlite.Dialector{
					DriverName: "DOES NOT EXIST",
					Conn:       nil,
				}, gorm.Config{})
				Expect(err).Should(HaveOccurred())
				Expect(gormDB).Should(BeNil())
			})
		})
	})
	Context("Initialise", func() {
		When("initialising without a connection", func() {
			It("should return an error", func() {
				DBConnection = nil
				err := Initialise()
				Expect(err).Should(HaveOccurred())
			})
		})

		When("initialising after connecting and gorm.AutoMigrate returns an error", func() {
			It("should return an error", func() {
				_, err := Connect(sqlite.Dialector{
					DriverName: "sqlite",
					Conn:       db,
				}, gorm.Config{})
				Expect(err).ShouldNot(HaveOccurred())

				err = Initialise(struct{}{})
				Expect(err).Should(HaveOccurred())
			})
		})

		When("initialising after connecting", func() {
			It("should return nil", func() {
				_, err := Connect(sqlite.Dialector{
					DriverName: "sqlite",
					Conn:       db,
				}, gorm.Config{})
				Expect(err).ShouldNot(HaveOccurred())

				err = Initialise()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(err).Should(BeNil())
			})
		})

		When("initialising with models after connecting", func() {
			sqlStatements := []string{
				"CREATE TABLE `users` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`user_name` text,`display_name` text,`password` text NOT NULL,PRIMARY KEY (`id`))",
				"CREATE UNIQUE INDEX `idx_users_user_name` ON `users`(`user_name`)",
				"CREATE INDEX `idx_users_deleted_at` ON `users`(`deleted_at`)",
				"CREATE TABLE `emails` (`email` text,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`user_id` integer,PRIMARY KEY (`email`),CONSTRAINT `fk_users_emails` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`))",
				"CREATE INDEX `idx_emails_deleted_at` ON `emails`(`deleted_at`)",
				"CREATE TABLE `groups` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text NOT NULL,`author_id` integer,PRIMARY KEY (`id`),CONSTRAINT `fk_groups_author` FOREIGN KEY (`author_id`) REFERENCES `users`(`id`))",
				"CREATE INDEX `idx_groups_deleted_at` ON `groups`(`deleted_at`)",
				"CREATE TABLE `users_groups` (`user_id` integer,`group_id` integer,PRIMARY KEY (`user_id`,`group_id`),CONSTRAINT `fk_users_groups_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),CONSTRAINT `fk_users_groups_group` FOREIGN KEY (`group_id`) REFERENCES `groups`(`id`))",
				"CREATE TABLE `topics` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`title` text,`parent_id` integer,PRIMARY KEY (`id`))",
				"CREATE UNIQUE INDEX `idx_topics_title` ON `topics`(`title`)",
				"CREATE INDEX `idx_topics_deleted_at` ON `topics`(`deleted_at`)",
				"CREATE TABLE `discussions` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`author_id` integer,`topic_id` integer,PRIMARY KEY (`id`),CONSTRAINT `fk_discussions_author` FOREIGN KEY (`author_id`) REFERENCES `users`(`id`),CONSTRAINT `fk_discussions_topic` FOREIGN KEY (`topic_id`) REFERENCES `topics`(`id`))",
				"CREATE INDEX `idx_discussions_deleted_at` ON `discussions`(`deleted_at`)",
				"CREATE TABLE `posts` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`content` text,`author_id` integer,PRIMARY KEY (`id`),CONSTRAINT `fk_posts_author` FOREIGN KEY (`author_id`) REFERENCES `users`(`id`))",
				"CREATE INDEX `idx_posts_deleted_at` ON `posts`(`deleted_at`)",
			}
			It("should run expected migrations on database", func() {
				_, err := Connect(sqlite.Dialector{
					DriverName: "sqlite",
					Conn:       db,
				}, gorm.Config{})
				Expect(err).ShouldNot(HaveOccurred())

				mock.MatchExpectationsInOrder(false)
				for _, sql := range sqlStatements {
					migrationSql := regexp.QuoteMeta(sql)
					mock.ExpectExec(migrationSql).WillReturnResult(sqlmock.NewResult(0, 0))
				}

				err = Initialise(models.Models()...)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
