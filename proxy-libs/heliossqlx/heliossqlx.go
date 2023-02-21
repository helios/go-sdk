package heliossqlx

import (
	"context"
	"database/sql"

	origin_sqlx "github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

type rowsi interface {
	Close() error
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(...interface{}) error
}

const UNKNOWN = origin_sqlx.UNKNOWN

const QUESTION = origin_sqlx.QUESTION

const DOLLAR = origin_sqlx.DOLLAR

const NAMED = origin_sqlx.NAMED

const AT = origin_sqlx.AT

func BindType(driverName string) int {
	return origin_sqlx.BindType(driverName)
}

func BindDriver(driverName string, bindType int) {
	origin_sqlx.BindDriver(driverName, bindType)
}

func Rebind(bindType int, query string) string {
	return origin_sqlx.Rebind(bindType, query)
}

func In(query string, args ...interface{}) (string, []interface{}, error) {
	return origin_sqlx.In(query, args)
}

func NamedQueryContext(ctx context.Context, e ExtContext, query string, arg interface{}) (*Rows, error) {
	return origin_sqlx.NamedQueryContext(ctx, e, query, arg)
}

func NamedExecContext(ctx context.Context, e ExtContext, query string, arg interface{}) (sql.Result, error) {
	return origin_sqlx.NamedExecContext(ctx, e, query, arg)
}

func ConnectContext(ctx context.Context, driverName, dataSourceName string, opts ...otelsql.Option) (*DB, error) {
	return otelsqlx.ConnectContext(ctx, driverName, dataSourceName, opts...)
}

type QueryerContext = origin_sqlx.QueryerContext

type PreparerContext = origin_sqlx.PreparerContext

type ExecerContext = origin_sqlx.ExecerContext

type ExtContext = origin_sqlx.ExtContext

func SelectContext(ctx context.Context, q QueryerContext, dest interface{}, query string, args ...interface{}) error {
	return origin_sqlx.SelectContext(ctx, q, dest, query, args)
}

func PreparexContext(ctx context.Context, p PreparerContext, query string) (*Stmt, error) {
	return origin_sqlx.PreparexContext(ctx, p, query)
}

func GetContext(ctx context.Context, q QueryerContext, dest interface{}, query string, args ...interface{}) error {
	return origin_sqlx.GetContext(ctx, q, dest, query, args)
}

func LoadFileContext(ctx context.Context, e ExecerContext, path string) (*sql.Result, error) {
	return origin_sqlx.LoadFileContext(ctx, e, path)
}

func MustExecContext(ctx context.Context, e ExecerContext, query string, args ...interface{}) sql.Result {
	return origin_sqlx.MustExecContext(ctx, e, query, args)
}

type NamedStmt = origin_sqlx.NamedStmt

func BindNamed(bindType int, query string, arg interface{}) (string, []interface{}, error) {
	return origin_sqlx.BindNamed(bindType, query, arg)
}

func Named(query string, arg interface{}) (string, []interface{}, error) {
	return origin_sqlx.Named(query, arg)
}

func NamedQuery(e Ext, query string, arg interface{}) (*Rows, error) {
	return origin_sqlx.NamedQuery(e, query, arg)
}

func NamedExec(e Ext, query string, arg interface{}) (sql.Result, error) {
	return origin_sqlx.NamedExec(e, query, arg)
}

var NameMapper = origin_sqlx.NameMapper

type ColScanner = origin_sqlx.ColScanner

type Queryer = origin_sqlx.Queryer

type Execer = origin_sqlx.Execer

type Ext = origin_sqlx.Ext

type Preparer = origin_sqlx.Preparer

type Row = origin_sqlx.Row

type DB = origin_sqlx.DB

func NewDb(db *sql.DB, driverName string) *DB {
	return origin_sqlx.NewDb(db, driverName)
}

func Open(driverName, dataSourceName string, opts ...otelsql.Option) (*DB, error) {
	return otelsqlx.Open(driverName, dataSourceName, opts...)
}

func MustOpen(driverName, dataSourceName string, opts ...otelsql.Option) *DB {
	return otelsqlx.MustOpen(driverName, dataSourceName, opts...)
}

type Conn = origin_sqlx.Conn

type Tx = origin_sqlx.Tx

type Stmt = origin_sqlx.Stmt

type Rows = origin_sqlx.Rows

func Connect(driverName, dataSourceName string, opts ...otelsql.Option) (*DB, error) {
	return otelsqlx.Connect(driverName, dataSourceName, opts ...)
}

func MustConnect(driverName, dataSourceName string, opts ...otelsql.Option) *DB {
	return otelsqlx.MustConnect(driverName, dataSourceName, opts...)
}

func Preparex(p Preparer, query string) (*Stmt, error) {
	return origin_sqlx.Preparex(p, query)
}

func Select(q Queryer, dest interface{}, query string, args ...interface{}) error {
	return origin_sqlx.Select(q, dest, query, args)
}

func Get(q Queryer, dest interface{}, query string, args ...interface{}) error {
	return origin_sqlx.Get(q, dest, query, args)
}

func LoadFile(e Execer, path string) (*sql.Result, error) {
	return origin_sqlx.LoadFile(e, path)
}

func MustExec(e Execer, query string, args ...interface{}) sql.Result {
	return origin_sqlx.MustExec(e, query, args)
}

func SliceScan(r ColScanner) ([]interface{}, error) {
	return origin_sqlx.SliceScan(r)
}

func MapScan(r ColScanner, dest map[string]interface{}) error {
	return origin_sqlx.MapScan(r, dest)
}

func StructScan(rows rowsi, dest interface{}) error {
	return origin_sqlx.StructScan(rows, dest)
}
