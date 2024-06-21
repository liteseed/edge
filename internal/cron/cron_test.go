package cron

import (
	"log/slog"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/liteseed/aogo"
	"github.com/liteseed/edge/test"
	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/transaction/data_item"
	"github.com/liteseed/sdk-go/contract"
	"github.com/stretchr/testify/assert"
)

func TestJobPostBundle(t *testing.T) {
	dataItem := test.DataItem()

	g := test.Gateway()
	defer g.Close()

	mock, db := test.Database()

	s := test.Store()
	defer s.Shutdown()

	w := test.Wallet(g.URL)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"ID", "Status"}).AddRow("dataitem", "queued"))
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := s.Set("dataitem", dataItem.Raw)
		assert.NoError(t, err)
		crn, err := New(WithDatabase(db), WithLogger(slog.Default()), WithStore(s), WithWallet(w))
		assert.NoError(t, err)

		crn.JobPostBundle()

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success Many Rows", func(t *testing.T) {
		dataItem1 := data_item.New([]byte("test"), "", "", &[]tag.Tag{})
		dataItem2 := data_item.New([]byte("test test"), "", "", &[]tag.Tag{})
		dataItem3 := data_item.New([]byte("test test test"), "", "", &[]tag.Tag{})
		dataItem4 := data_item.New([]byte("test test test test"), "", "", &[]tag.Tag{})
		dataItem5 := data_item.New([]byte("test test test test test"), "", "", &[]tag.Tag{})

		_, _ = w.SignDataItem(dataItem1)
		_, _ = w.SignDataItem(dataItem2)
		_, _ = w.SignDataItem(dataItem3)
		_, _ = w.SignDataItem(dataItem4)
		_, _ = w.SignDataItem(dataItem5)

		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"ID", "Status"}).AddRow("dataitem1", "queued").AddRow("dataitem2", "queued").AddRow("dataitem3", "queued").AddRow("dataitem4", "queued").AddRow("dataitem5", "queued"))
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := s.Set("dataitem1", dataItem1.Raw)
		assert.NoError(t, err)
		err = s.Set("dataitem2", dataItem2.Raw)
		assert.NoError(t, err)
		err = s.Set("dataitem3", dataItem3.Raw)
		assert.NoError(t, err)
		err = s.Set("dataitem4", dataItem4.Raw)
		assert.NoError(t, err)
		err = s.Set("dataitem5", dataItem5.Raw)
		assert.NoError(t, err)

		crn, err := New(WithDatabase(db), WithLogger(slog.Default()), WithStore(s), WithWallet(w))
		assert.NoError(t, err)

		crn.JobPostBundle()

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"ID", "Status"}).AddRow(dataItem.ID, "queued"))

		err := s.Set(dataItem.ID, nil)
		assert.NoError(t, err)
		crn, err := New(WithDatabase(db), WithLogger(slog.Default()), WithStore(s), WithWallet(w))
		assert.NoError(t, err)

		crn.JobPostBundle()

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success Fail Mix", func(t *testing.T) {
		dataItem1 := data_item.New([]byte("test"), "", "", &[]tag.Tag{})
		dataItem2 := data_item.New([]byte("test test"), "", "", &[]tag.Tag{})
		dataItem3 := data_item.New([]byte("test test test"), "", "", &[]tag.Tag{})
		dataItem4 := data_item.New([]byte("test test test test"), "", "", &[]tag.Tag{})
		dataItem5 := data_item.New([]byte("test test test test test"), "", "", &[]tag.Tag{})

		_, _ = w.SignDataItem(dataItem1)
		_, _ = w.SignDataItem(dataItem2)
		_, _ = w.SignDataItem(dataItem3)
		_, _ = w.SignDataItem(dataItem4)
		_, _ = w.SignDataItem(dataItem5)

		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(dataItem1.ID, "queued").AddRow(dataItem2.ID, "queued").AddRow(dataItem3.ID, "queued").AddRow(dataItem4.ID, "queued").AddRow(dataItem5.ID, "queued"))
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := s.Set(dataItem1.ID, dataItem1.Raw)
		assert.NoError(t, err)
		err = s.Set(dataItem2.ID, dataItem2.Raw)
		assert.NoError(t, err)
		err = s.Set(dataItem3.ID, nil)
		assert.NoError(t, err)
		err = s.Set(dataItem4.ID, nil)
		assert.NoError(t, err)
		err = s.Set(dataItem5.ID, dataItem5.Raw)
		assert.NoError(t, err)

		crn, err := New(WithDatabase(db), WithLogger(slog.Default()), WithStore(s), WithWallet(w))
		assert.NoError(t, err)

		crn.JobPostBundle()

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestJobBundleConfirmations(t *testing.T) {
	g := test.Gateway()
	defer g.Close()

	mock, db := test.Database()

	w := test.Wallet(g.URL)
	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "status", "bundle_id"}).AddRow("dataitem", "sent", "bundle")
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		crn, err := New(WithDatabase(db), WithLogger(slog.Default()), WithWallet(w))
		assert.NoError(t, err)

		crn.JobBundleConfirmations()
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "status", "bundle_id"}).AddRow("dataitem", "sent", "failbundle")
		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		crn, err := New(WithDatabase(db), WithLogger(slog.Default()), WithWallet(w))
		assert.NoError(t, err)

		crn.JobBundleConfirmations()
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestJobPostUpdate(t *testing.T) {
	g := test.Gateway()
	defer g.Close()

	mock, db := test.Database()

	mu := test.MU()
	defer mu.Close()

	ao, err := aogo.New(aogo.WthMU(mu.URL))
	assert.NoError(t, err)

	w := test.Wallet(g.URL)

	c := contract.Custom(ao, "process", w.Signer)

	crn, err := New(WithContracts(c), WithDatabase(db), WithLogger(slog.Default()), WithWallet(w))
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "status"}).AddRow("dataitem", "confirmed")
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		crn.JobBundleConfirmations()
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}


func TestJobRelease(t *testing.T) {
	g := test.Gateway()
	defer g.Close()

	mock, db := test.Database()

	mu := test.MU()
	defer mu.Close()

	ao, err := aogo.New(aogo.WthMU(mu.URL))
	assert.NoError(t, err)

	w := test.Wallet(g.URL)

	c := contract.Custom(ao, "process", w.Signer)

	crn, err := New(WithContracts(c), WithDatabase(db), WithLogger(slog.Default()), WithWallet(w))
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "status"}).AddRow("dataitem", "confirmed")
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		crn.JobBundleConfirmations()
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestJobDeleteBundle(t *testing.T) {
	g := test.Gateway()
	defer g.Close()

	mock, db := test.Database()

	mu := test.MU()
	defer mu.Close()

	ao, err := aogo.New(aogo.WthMU(mu.URL))
	assert.NoError(t, err)

	w := test.Wallet(g.URL)

	c := contract.Custom(ao, "process", w.Signer)

	crn, err := New(WithContracts(c), WithDatabase(db), WithLogger(slog.Default()), WithWallet(w))
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "status"}).AddRow("dataitem", "confirmed")
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "orders" WHERE "orders"."id" = $1`)).WithArgs("dataitem").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		crn.JobDeleteBundle()
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

