package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golangbb/golangbb/v2/internal/database"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

var _ = Describe("Models", func() {
	When("Models is executed", func() {
		It("should return a slice containing a pointer to each model", func() {
			models := Models()
			Expect(models).Should(Equal([]interface{}{
				&Discussion{}, &Email{}, &Group{}, &Post{}, &Topic{}, &User{},
			}))
		})
		It("should return a slice with a length equal to the number of models defined", func() {
			cfg := &packages.Config{
				Mode: packages.NeedImports | packages.NeedTypes,
			}
			pkgs, err := packages.Load(cfg, "github.com/golangbb/golangbb/v2/internal/models")
			Expect(err).Should(BeNil())
			Expect(len(pkgs)).Should(Equal(1))

			pkg := pkgs[0]
			scope := pkg.Types.Scope()
			numberOfExportedModels := 0
			for _, name := range scope.Names() {
				obj := scope.Lookup(name)
				if !obj.Exported() {
					continue
				}

				if !strings.HasPrefix(obj.String(), "type") {
					continue
				}

				numberOfExportedModels += 1
			}

			numberOfModels := len(Models())
			Expect(numberOfModels).Should(Equal(numberOfExportedModels))
		})
	})

	Context("Migrations in database.Initialise", func() {
		When("initialising with models", func() {
			sqlStatements := []string{
				"CREATE TABLE `users` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`user_name` text,`display_name` text NOT NULL,`password` text NOT NULL,PRIMARY KEY (`id`))",
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
				db, mock, err := sqlmock.New()
				Expect(err).ShouldNot(HaveOccurred())
				defer db.Close()

				_, err = database.Connect(sqlite.Dialector{
					DriverName: "sqlite",
					Conn:       db,
				}, gorm.Config{})
				Expect(err).ShouldNot(HaveOccurred())

				mock.MatchExpectationsInOrder(false)
				for _, sql := range sqlStatements {
					migrationSql := regexp.QuoteMeta(sql)
					mock.ExpectExec(migrationSql).WillReturnResult(sqlmock.NewResult(0, 0))
				}

				err = database.Initialise(Models()...)
				Expect(err).ShouldNot(HaveOccurred())

				err = mock.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
