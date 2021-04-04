package database

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
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

				mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE `` ()")).
					WillReturnError(errors.New("some db error"))

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
	})
})
