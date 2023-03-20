package heliospg

import (
	"context"

	"github.com/go-pg/pg/extra/pgotel/v10"
	origin_pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)
type PoolStats = origin_pg.PoolStats

var ErrNoRows = origin_pg.ErrNoRows

var ErrMultiRows = origin_pg.ErrMultiRows

type Error = origin_pg.Error

type Options = origin_pg.Options

func ParseURL(sURL string) (*Options,error) {
	return origin_pg.ParseURL(sURL)
 }

var Discard = origin_pg.Discard

type NullTime = origin_pg.NullTime

func Scan(values ...interface{}) (orm.ColumnScanner) {
	return origin_pg.Scan(values)
 }

type Safe = origin_pg.Safe

type Ident = origin_pg.Ident

func SafeQuery(query string,params ...interface{}) (*orm.SafeQueryAppender) {
	return origin_pg.SafeQuery(query,params)
 }

func In(slice interface{}) (types.ValueAppender) {
	return origin_pg.In(slice)
 }

func InMulti(values ...interface{}) (types.ValueAppender) {
	return origin_pg.InMulti(values)
 }

func Array(v interface{}) (*types.Array) {
	return origin_pg.Array(v)
 }

func Hstore(v interface{}) (*types.Hstore) {
	return origin_pg.Hstore(v)
 }

 type Logging interface {
	Printf(ctx context.Context, format string, v ...interface{})
}

func SetLogger(logger Logging){
	origin_pg.SetLogger(logger)
 }

type Query = origin_pg.Query

func Model(model ...interface{}) (*Query) {
	return origin_pg.Model(model)
 }

func ModelContext(c context.Context,model ...interface{}) (*Query) {
	return origin_pg.ModelContext(c,model)
 }

type DBI = origin_pg.DBI

type Strings = origin_pg.Strings

type Ints = origin_pg.Ints

type IntSet = origin_pg.IntSet

type Result = origin_pg.Result

func Connect(opt *Options) (*DB) {
	db := origin_pg.Connect(opt)
	db.AddQueryHook(pgotel.NewTracingHook())
	return db
 }

type DB = origin_pg.DB

type Conn = origin_pg.Conn

type BeforeScanHook = origin_pg.BeforeScanHook

type AfterScanHook = origin_pg.AfterScanHook

type AfterSelectHook = origin_pg.AfterSelectHook

type BeforeInsertHook = origin_pg.BeforeInsertHook

type AfterInsertHook = origin_pg.AfterInsertHook

type BeforeUpdateHook = origin_pg.BeforeUpdateHook

type AfterUpdateHook = origin_pg.AfterUpdateHook

type BeforeDeleteHook = origin_pg.BeforeDeleteHook

type AfterDeleteHook = origin_pg.AfterDeleteHook

type QueryEvent = origin_pg.QueryEvent

type QueryHook = origin_pg.QueryHook

type Notification = origin_pg.Notification

type Listener = origin_pg.Listener

type Stmt = origin_pg.Stmt

var ErrTxDone = origin_pg.ErrTxDone

type Tx = origin_pg.Tx

func Version() (string) {
	return origin_pg.Version()
 }
