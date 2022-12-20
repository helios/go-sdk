package heliosmongo
import (	"context"

"go.mongodb.org/mongo-driver/bson/bsoncodec"
origin_mongo "go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
"go.mongodb.org/mongo-driver/x/mongo/driver"
"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)
type Database = origin_mongo.Database

type IndexOptionsBuilder = origin_mongo.IndexOptionsBuilder

func NewIndexOptionsBuilder() (*IndexOptionsBuilder) {
	return origin_mongo.NewIndexOptionsBuilder()
 }

var ErrInvalidIndexValue = origin_mongo.ErrInvalidIndexValue

var ErrNonStringIndexName = origin_mongo.ErrNonStringIndexName

var ErrMultipleIndexDrop = origin_mongo.ErrMultipleIndexDrop

type IndexView = origin_mongo.IndexView

type IndexModel = origin_mongo.IndexModel

var ErrMissingResumeToken = origin_mongo.ErrMissingResumeToken

var ErrNilCursor = origin_mongo.ErrNilCursor

type ChangeStream = origin_mongo.ChangeStream

type StreamType = origin_mongo.StreamType

const CollectionStream = origin_mongo.CollectionStream

const DatabaseStream = origin_mongo.DatabaseStream

const ClientStream = origin_mongo.ClientStream

type BulkWriteResult = origin_mongo.BulkWriteResult

type InsertOneResult = origin_mongo.InsertOneResult

type InsertManyResult = origin_mongo.InsertManyResult

type DeleteResult = origin_mongo.DeleteResult

type ListDatabasesResult = origin_mongo.ListDatabasesResult

type DatabaseSpecification = origin_mongo.DatabaseSpecification

type UpdateResult = origin_mongo.UpdateResult

type IndexSpecification = origin_mongo.IndexSpecification

type CollectionSpecification = origin_mongo.CollectionSpecification

type Dialer = origin_mongo.Dialer

type BSONAppender = origin_mongo.BSONAppender

type BSONAppenderFunc = origin_mongo.BSONAppenderFunc

type MarshalError = origin_mongo.MarshalError

type Pipeline = origin_mongo.Pipeline

type ClientEncryption = origin_mongo.ClientEncryption

func setMonitor(opts []*options.ClientOptions) []*options.ClientOptions {
	monitor := otelmongo.NewMonitor()

	for _, opt := range opts {
		if opt != nil {
			opt.Monitor = monitor
		}
	}

	return opts
}

func NewClientEncryption(keyVaultClient *Client,opts ...*options.ClientEncryptionOptions) (*ClientEncryption,error) {
	return origin_mongo.NewClientEncryption(keyVaultClient,opts ...)
 }

type Cursor = origin_mongo.Cursor

func NewCursorFromDocuments(documents []interface{},err error,registry *bsoncodec.Registry) (*Cursor,error) {
	return origin_mongo.NewCursorFromDocuments(documents,err,registry)
 }

func BatchCursorFromCursor(c *Cursor) (*driver.BatchCursor) {
	return origin_mongo.BatchCursorFromCursor(c)
 }

var ErrWrongClient = origin_mongo.ErrWrongClient

type SessionContext = origin_mongo.SessionContext

func NewSessionContext(ctx context.Context,sess Session) (SessionContext) {
	return origin_mongo.NewSessionContext(ctx,sess)
 }

func SessionFromContext(ctx context.Context) (Session) {
	return origin_mongo.SessionFromContext(ctx)
 }

type Session = origin_mongo.Session

type XSession = origin_mongo.XSession

var ErrNoDocuments = origin_mongo.ErrNoDocuments

type SingleResult = origin_mongo.SingleResult

func NewSingleResultFromDocument(document interface{},err error,registry *bsoncodec.Registry) (*SingleResult) {
	return origin_mongo.NewSingleResultFromDocument(document,err,registry)
 }

type Client = origin_mongo.Client

func Connect(ctx context.Context,opts ...*options.ClientOptions) (*Client,error) {
	return origin_mongo.Connect(ctx,setMonitor(opts)...)
 }

func NewClient(opts ...*options.ClientOptions) (*Client,error) {
	return origin_mongo.NewClient(opts ...)
 }

func WithSession(ctx context.Context,sess Session,fn func(SessionContext) error) (error) {
	return origin_mongo.WithSession(ctx,sess,fn)
 }

type WriteModel = origin_mongo.WriteModel

type InsertOneModel = origin_mongo.InsertOneModel

func NewInsertOneModel() (*InsertOneModel) {
	return origin_mongo.NewInsertOneModel()
 }

type DeleteOneModel = origin_mongo.DeleteOneModel

func NewDeleteOneModel() (*DeleteOneModel) {
	return origin_mongo.NewDeleteOneModel()
 }

type DeleteManyModel = origin_mongo.DeleteManyModel

func NewDeleteManyModel() (*DeleteManyModel) {
	return origin_mongo.NewDeleteManyModel()
 }

type ReplaceOneModel = origin_mongo.ReplaceOneModel

func NewReplaceOneModel() (*ReplaceOneModel) {
	return origin_mongo.NewReplaceOneModel()
 }

type UpdateOneModel = origin_mongo.UpdateOneModel

func NewUpdateOneModel() (*UpdateOneModel) {
	return origin_mongo.NewUpdateOneModel()
 }

type UpdateManyModel = origin_mongo.UpdateManyModel

func NewUpdateManyModel() (*UpdateManyModel) {
	return origin_mongo.NewUpdateManyModel()
 }

type Collection = origin_mongo.Collection

var ErrUnacknowledgedWrite = origin_mongo.ErrUnacknowledgedWrite

var ErrClientDisconnected = origin_mongo.ErrClientDisconnected

var ErrNilDocument = origin_mongo.ErrNilDocument

var ErrNilValue = origin_mongo.ErrNilValue

var ErrEmptySlice = origin_mongo.ErrEmptySlice

type ErrMapForOrderedArgument = origin_mongo.ErrMapForOrderedArgument

func IsDuplicateKeyError(err error) (bool) {
	return origin_mongo.IsDuplicateKeyError(err)
 }

func IsTimeout(err error) (bool) {
	return origin_mongo.IsTimeout(err)
 }

func IsNetworkError(err error) (bool) {
	return origin_mongo.IsNetworkError(err)
 }

type MongocryptError = origin_mongo.MongocryptError

type EncryptionKeyVaultError = origin_mongo.EncryptionKeyVaultError

type MongocryptdError = origin_mongo.MongocryptdError

type ServerError = origin_mongo.ServerError

type CommandError = origin_mongo.CommandError

type WriteError = origin_mongo.WriteError

type WriteErrors = origin_mongo.WriteErrors

type WriteConcernError = origin_mongo.WriteConcernError

type WriteException = origin_mongo.WriteException

type BulkWriteError = origin_mongo.BulkWriteError

type BulkWriteException = origin_mongo.BulkWriteException

