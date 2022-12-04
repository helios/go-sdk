package heliosmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	originalMongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

const ClientStream = originalMongo.ClientStream
const CollectionStream = originalMongo.CollectionStream
const DatabaseStream = originalMongo.DatabaseStream

type BSONAppender = originalMongo.BSONAppender
type BSONAppenderFunc = originalMongo.BSONAppenderFunc
type BulkWriteError = originalMongo.BulkWriteError
type BulkWriteException = originalMongo.BulkWriteException
type BulkWriteResult = originalMongo.BulkWriteResult
type ChangeStream = originalMongo.ChangeStream
type Client = originalMongo.Client
type ClientEncryption = originalMongo.ClientEncryption
type Collection = originalMongo.Collection
type CollectionSpecification = originalMongo.CollectionSpecification
type CommandError = originalMongo.CommandError
type Cursor = originalMongo.Cursor
type Database = originalMongo.Database
type DatabaseSpecification = originalMongo.DatabaseSpecification
type DeleteManyModel = originalMongo.DeleteManyModel
type DeleteOneModel = originalMongo.DeleteOneModel
type DeleteResult = originalMongo.DeleteResult
type Dialer = originalMongo.Dialer
type EncryptionKeyVaultError = originalMongo.EncryptionKeyVaultError
type ErrMapForOrderedArgument = originalMongo.ErrMapForOrderedArgument
type IndexModel = originalMongo.IndexModel
type IndexOptionsBuilder = originalMongo.IndexOptionsBuilder
type IndexSpecification = originalMongo.IndexSpecification
type IndexView = originalMongo.IndexView
type InsertManyResult = originalMongo.InsertManyResult
type InsertOneModel = originalMongo.InsertOneModel
type InsertOneResult = originalMongo.InsertOneResult
type ListDatabasesResult = originalMongo.ListDatabasesResult
type MarshalError = originalMongo.MarshalError
type MongocryptError = originalMongo.MongocryptError
type MongocryptdError = originalMongo.MongocryptdError
type Pipeline = originalMongo.Pipeline
type ReplaceOneModel = originalMongo.ReplaceOneModel
type RewrapManyDataKeyResult = originalMongo.RewrapManyDataKeyResult
type ServerError = originalMongo.ServerError
type Session = originalMongo.Session
type SessionContext = originalMongo.SessionContext
type SingleResult = originalMongo.SingleResult
type StreamType = originalMongo.StreamType
type UpdateManyModel = originalMongo.UpdateManyModel
type UpdateOneModel = originalMongo.UpdateOneModel
type UpdateResult = originalMongo.UpdateResult
type WriteConcernError = originalMongo.WriteConcernError
type WriteError = originalMongo.WriteError
type WriteErrors = originalMongo.WriteErrors
type WriteException = originalMongo.WriteException
type WriteModel = originalMongo.WriteModel
type XSession = originalMongo.XSession

var ErrClientDisconnected = originalMongo.ErrClientDisconnected
var ErrEmptySlice = originalMongo.ErrEmptySlice
var ErrInvalidIndexValue = originalMongo.ErrInvalidIndexValue
var ErrMissingResumeToken = originalMongo.ErrMissingResumeToken
var ErrMultipleIndexDrop = originalMongo.ErrMultipleIndexDrop
var ErrNilCursor = originalMongo.ErrNilCursor
var ErrNilDocument = originalMongo.ErrNilDocument
var ErrNilValue = originalMongo.ErrNilValue
var ErrNoDocuments = originalMongo.ErrNoDocuments
var ErrNonStringIndexName = originalMongo.ErrNonStringIndexName
var ErrUnacknowledgedWrite = originalMongo.ErrUnacknowledgedWrite
var ErrWrongClient = originalMongo.ErrWrongClient

func setMonitor(opts []*options.ClientOptions) []*options.ClientOptions {
	monitor := otelmongo.NewMonitor()

	for _, opt := range opts {
		if opt != nil {
			opt.Monitor = monitor
		}
	}

	return opts
}

func BatchCursorFromCursor(c *Cursor) *driver.BatchCursor {
	return originalMongo.BatchCursorFromCursor(c)
}

func Connect(ctx context.Context, opts ...*options.ClientOptions) (*Client, error) {
	return originalMongo.Connect(ctx, setMonitor(opts)...)
}

func IsDuplicateKeyError(err error) bool {
	return originalMongo.IsDuplicateKeyError(err)
}

func IsNetworkError(err error) bool {
	return originalMongo.IsNetworkError(err)
}

func IsTimeout(err error) bool {
	return originalMongo.IsTimeout(err)
}

func NewClient(opts ...*options.ClientOptions) (*Client, error) {
	return originalMongo.NewClient(setMonitor(opts)...)
}

func NewClientEncryption(keyVaultClient *Client, opts ...*options.ClientEncryptionOptions) (*ClientEncryption, error) {
	return originalMongo.NewClientEncryption(keyVaultClient, opts...)
}

func NewCursorFromDocuments(documents []interface{}, err error, registry *bsoncodec.Registry) (*Cursor, error) {
	return originalMongo.NewCursorFromDocuments(documents, err, registry)
}

func NewDeleteManyModel() *DeleteManyModel {
	return originalMongo.NewDeleteManyModel()
}

func NewDeleteOneModel() *DeleteOneModel {
	return originalMongo.NewDeleteOneModel()
}

func NewIndexOptionsBuilder() *IndexOptionsBuilder {
	return originalMongo.NewIndexOptionsBuilder()
}

func NewInsertOneModel() *InsertOneModel {
	return originalMongo.NewInsertOneModel()
}

func NewReplaceOneModel() *ReplaceOneModel {
	return originalMongo.NewReplaceOneModel()
}

func NewSessionContext(ctx context.Context, sess Session) SessionContext {
	return originalMongo.NewSessionContext(ctx, sess)
}

func NewSingleResultFromDocument(document interface{}, err error, registry *bsoncodec.Registry) *SingleResult {
	return originalMongo.NewSingleResultFromDocument(document, err, registry)
}

func NewUpdateManyModel() *UpdateManyModel {
	return originalMongo.NewUpdateManyModel()
}

func NewUpdateOneModel() *UpdateOneModel {
	return originalMongo.NewUpdateOneModel()
}

func SessionFromContext(ctx context.Context) Session {
	return originalMongo.SessionFromContext(ctx)
}

func WithSession(ctx context.Context, sess Session, fn func(SessionContext) error) error {
	return originalMongo.WithSession(ctx, sess, fn)
}
