package postgres

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeDB struct {
	row      pgx.Row
	queued   []pgx.Rows
	queryErr error
	execTag  pgconn.CommandTag
	execErr  error
	execSet  bool
	tx       pgx.Tx
	beginErr error
}

func (f *fakeDB) Begin(context.Context) (pgx.Tx, error) {
	if f.beginErr != nil {
		return nil, f.beginErr
	}
	return f.tx, nil
}

func (f *fakeDB) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	if f.execErr != nil {
		return pgconn.CommandTag{}, f.execErr
	}
	if !f.execSet {
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}
	return f.execTag, nil
}

func (f *fakeDB) Query(context.Context, string, ...any) (pgx.Rows, error) {
	if f.queryErr != nil {
		return nil, f.queryErr
	}
	if len(f.queued) == 0 {
		return &fakeRows{}, nil
	}
	rows := f.queued[0]
	f.queued = f.queued[1:]
	return rows, nil
}

func (f *fakeDB) QueryRow(context.Context, string, ...any) pgx.Row {
	return f.row
}

type errRow struct {
	err error
}

func (r errRow) Scan(...any) error {
	return r.err
}

type valueRow struct {
	values []any
}

func (r valueRow) Scan(dest ...any) error {
	return assign(dest, r.values...)
}

type fakeRows struct {
	scans []func(dest ...any) error
	n     int
	err   error
}

func (r *fakeRows) Close() {}

func (r *fakeRows) Err() error { return r.err }

func (r *fakeRows) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }

func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription { return nil }

func (r *fakeRows) Next() bool { return r.n < len(r.scans) }

func (r *fakeRows) Scan(dest ...any) error {
	fn := r.scans[r.n]
	r.n++
	return fn(dest...)
}

func (r *fakeRows) Values() ([]any, error) { return nil, nil }

func (r *fakeRows) RawValues() [][]byte { return nil }

func (r *fakeRows) Conn() *pgx.Conn { return nil }

type fakeTx struct {
	row       pgx.Row
	execErrs  []error
	execTags  []pgconn.CommandTag
	n         int
	commitErr error
}

func (f *fakeTx) Begin(context.Context) (pgx.Tx, error) {
	panic("unexpected Begin")
}

func (f *fakeTx) Commit(context.Context) error { return f.commitErr }

func (f *fakeTx) Rollback(context.Context) error { return nil }

func (f *fakeTx) CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error) {
	panic("unexpected CopyFrom")
}

func (f *fakeTx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults {
	panic("unexpected SendBatch")
}

func (f *fakeTx) LargeObjects() pgx.LargeObjects {
	panic("unexpected LargeObjects")
}

func (f *fakeTx) Prepare(context.Context, string, string) (*pgconn.StatementDescription, error) {
	panic("unexpected Prepare")
}

func (f *fakeTx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	var err error
	tag := pgconn.NewCommandTag("UPDATE 1")
	if f.n < len(f.execErrs) {
		err = f.execErrs[f.n]
	}
	if f.n < len(f.execTags) {
		tag = f.execTags[f.n]
	}
	f.n++
	return tag, err
}

func (f *fakeTx) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return &fakeRows{}, nil
}

func (f *fakeTx) QueryRow(context.Context, string, ...any) pgx.Row {
	return f.row
}

func (f *fakeTx) Conn() *pgx.Conn { return nil }

func assign(dest []any, values ...any) error {
	for i, value := range values {
		target := reflect.ValueOf(dest[i]).Elem()
		if value == nil {
			target.Set(reflect.Zero(target.Type()))
			continue
		}
		incoming := reflect.ValueOf(value)
		if incoming.Type().AssignableTo(target.Type()) {
			target.Set(incoming)
			continue
		}
		if !incoming.Type().ConvertibleTo(target.Type()) {
			return fmt.Errorf("scan destination %d: %s to %s", i, incoming.Type(), target.Type())
		}
		target.Set(incoming.Convert(target.Type()))
	}
	return nil
}

func pgError(code string) error {
	return &pgconn.PgError{Code: code}
}

func updated(n string) pgconn.CommandTag {
	return pgconn.NewCommandTag(n)
}
