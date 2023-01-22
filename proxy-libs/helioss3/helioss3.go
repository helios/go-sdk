package helioss3

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	origin_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

type ListBucketsInput = origin_s3.ListBucketsInput

type ListBucketsOutput = origin_s3.ListBucketsOutput

type PutBucketReplicationInput = origin_s3.PutBucketReplicationInput

type PutBucketReplicationOutput = origin_s3.PutBucketReplicationOutput

type PutObjectTaggingInput = origin_s3.PutObjectTaggingInput

type PutObjectTaggingOutput = origin_s3.PutObjectTaggingOutput

const ServiceID = origin_s3.ServiceID

const ServiceAPIVersion = origin_s3.ServiceAPIVersion

type Client = origin_s3.Client

func New(options Options,optFns ...func(*Options) ) (*Client) {
	otelaws.AppendMiddlewares(&options.APIOptions)
	return origin_s3.New(options, optFns...)
 }

type Options = origin_s3.Options

func WithAPIOptions(optFns ...func(*middleware.Stack) error) (func(*Options) ) {
	return origin_s3.WithAPIOptions(optFns...)
 }

func WithEndpointResolver(v EndpointResolver) (func(*Options) ) {
	return origin_s3.WithEndpointResolver(v)
 }

type HTTPClient = origin_s3.HTTPClient

func NewFromConfig(cfg aws.Config,optFns ...func(*Options) ) (*Client) {
	otelaws.AppendMiddlewares(&cfg.APIOptions)
	return origin_s3.NewFromConfig(cfg,optFns...)
 }

type HTTPSignerV4 = origin_s3.HTTPSignerV4

type ComputedInputChecksumsMetadata = origin_s3.ComputedInputChecksumsMetadata

func GetComputedInputChecksumsMetadata(m middleware.Metadata) (ComputedInputChecksumsMetadata,bool) {
	return origin_s3.GetComputedInputChecksumsMetadata(m)
 }

type ChecksumValidationMetadata = origin_s3.ChecksumValidationMetadata

func GetChecksumValidationMetadata(m middleware.Metadata) (ChecksumValidationMetadata,bool) {
	return origin_s3.GetChecksumValidationMetadata(m)
 }

type ResponseError = origin_s3.ResponseError

func GetHostIDMetadata(metadata middleware.Metadata) (string,bool) {
	return origin_s3.GetHostIDMetadata(metadata)
 }

type HTTPPresignerV4 = origin_s3.HTTPPresignerV4

type PresignOptions = origin_s3.PresignOptions

func WithPresignClientFromClientOptions(optFns ...func(*Options) ) (func(*PresignOptions) ) {
	return origin_s3.WithPresignClientFromClientOptions(optFns...)
 }

func WithPresignExpires(dur time.Duration) (func(*PresignOptions) ) {
	return origin_s3.WithPresignExpires(dur)
 }

type PresignClient = origin_s3.PresignClient

func NewPresignClient(c *Client,optFns ...func(*PresignOptions) ) (*PresignClient) {
	return origin_s3.NewPresignClient(c,optFns...)
 }

type DeletePublicAccessBlockInput = origin_s3.DeletePublicAccessBlockInput

type DeletePublicAccessBlockOutput = origin_s3.DeletePublicAccessBlockOutput

type ListBucketAnalyticsConfigurationsInput = origin_s3.ListBucketAnalyticsConfigurationsInput

type ListBucketAnalyticsConfigurationsOutput = origin_s3.ListBucketAnalyticsConfigurationsOutput

type PutBucketAccelerateConfigurationInput = origin_s3.PutBucketAccelerateConfigurationInput

type PutBucketAccelerateConfigurationOutput = origin_s3.PutBucketAccelerateConfigurationOutput

type DeleteObjectTaggingInput = origin_s3.DeleteObjectTaggingInput

type DeleteObjectTaggingOutput = origin_s3.DeleteObjectTaggingOutput

type GetBucketMetricsConfigurationInput = origin_s3.GetBucketMetricsConfigurationInput

type GetBucketMetricsConfigurationOutput = origin_s3.GetBucketMetricsConfigurationOutput

type ListPartsInput = origin_s3.ListPartsInput

type ListPartsOutput = origin_s3.ListPartsOutput

type ListPartsAPIClient = origin_s3.ListPartsAPIClient

type ListPartsPaginatorOptions = origin_s3.ListPartsPaginatorOptions

type ListPartsPaginator = origin_s3.ListPartsPaginator

func NewListPartsPaginator(client ListPartsAPIClient,params *ListPartsInput,optFns ...func(*ListPartsPaginatorOptions) ) (*ListPartsPaginator) {
	return origin_s3.NewListPartsPaginator(client,params,optFns...)
 }

type PutBucketOwnershipControlsInput = origin_s3.PutBucketOwnershipControlsInput

type PutBucketOwnershipControlsOutput = origin_s3.PutBucketOwnershipControlsOutput

type PutObjectLockConfigurationInput = origin_s3.PutObjectLockConfigurationInput

type PutObjectLockConfigurationOutput = origin_s3.PutObjectLockConfigurationOutput

type PutPublicAccessBlockInput = origin_s3.PutPublicAccessBlockInput

type PutPublicAccessBlockOutput = origin_s3.PutPublicAccessBlockOutput

type RestoreObjectInput = origin_s3.RestoreObjectInput

type RestoreObjectOutput = origin_s3.RestoreObjectOutput

type UploadPartCopyInput = origin_s3.UploadPartCopyInput

type UploadPartCopyOutput = origin_s3.UploadPartCopyOutput

type GetBucketAclInput = origin_s3.GetBucketAclInput

type GetBucketAclOutput = origin_s3.GetBucketAclOutput

type GetBucketTaggingInput = origin_s3.GetBucketTaggingInput

type GetBucketTaggingOutput = origin_s3.GetBucketTaggingOutput

type PutBucketLifecycleConfigurationInput = origin_s3.PutBucketLifecycleConfigurationInput

type PutBucketLifecycleConfigurationOutput = origin_s3.PutBucketLifecycleConfigurationOutput

type GetBucketLocationInput = origin_s3.GetBucketLocationInput

type GetBucketLocationOutput = origin_s3.GetBucketLocationOutput

type GetObjectTorrentInput = origin_s3.GetObjectTorrentInput

type GetObjectTorrentOutput = origin_s3.GetObjectTorrentOutput

type PutBucketAclInput = origin_s3.PutBucketAclInput

type PutBucketAclOutput = origin_s3.PutBucketAclOutput

type PutBucketVersioningInput = origin_s3.PutBucketVersioningInput

type PutBucketVersioningOutput = origin_s3.PutBucketVersioningOutput

type SelectObjectContentEventStreamReader = origin_s3.SelectObjectContentEventStreamReader

type UnknownEventMessageError = origin_s3.UnknownEventMessageError

type DeleteBucketInventoryConfigurationInput = origin_s3.DeleteBucketInventoryConfigurationInput

type DeleteBucketInventoryConfigurationOutput = origin_s3.DeleteBucketInventoryConfigurationOutput

type DeleteBucketPolicyInput = origin_s3.DeleteBucketPolicyInput

type DeleteBucketPolicyOutput = origin_s3.DeleteBucketPolicyOutput

type GetBucketPolicyStatusInput = origin_s3.GetBucketPolicyStatusInput

type GetBucketPolicyStatusOutput = origin_s3.GetBucketPolicyStatusOutput

type GetBucketWebsiteInput = origin_s3.GetBucketWebsiteInput

type GetBucketWebsiteOutput = origin_s3.GetBucketWebsiteOutput

type GetObjectAclInput = origin_s3.GetObjectAclInput

type GetObjectAclOutput = origin_s3.GetObjectAclOutput

type GetPublicAccessBlockInput = origin_s3.GetPublicAccessBlockInput

type GetPublicAccessBlockOutput = origin_s3.GetPublicAccessBlockOutput

type ListObjectVersionsInput = origin_s3.ListObjectVersionsInput

type ListObjectVersionsOutput = origin_s3.ListObjectVersionsOutput

type DeleteBucketLifecycleInput = origin_s3.DeleteBucketLifecycleInput

type DeleteBucketLifecycleOutput = origin_s3.DeleteBucketLifecycleOutput

type DeleteObjectsInput = origin_s3.DeleteObjectsInput

type DeleteObjectsOutput = origin_s3.DeleteObjectsOutput

type GetBucketOwnershipControlsInput = origin_s3.GetBucketOwnershipControlsInput

type GetBucketOwnershipControlsOutput = origin_s3.GetBucketOwnershipControlsOutput

type ListObjectsV2Input = origin_s3.ListObjectsV2Input

type ListObjectsV2Output = origin_s3.ListObjectsV2Output

type ListObjectsV2APIClient = origin_s3.ListObjectsV2APIClient

type ListObjectsV2PaginatorOptions = origin_s3.ListObjectsV2PaginatorOptions

type ListObjectsV2Paginator = origin_s3.ListObjectsV2Paginator

func NewListObjectsV2Paginator(client ListObjectsV2APIClient,params *ListObjectsV2Input,optFns ...func(*ListObjectsV2PaginatorOptions) ) (*ListObjectsV2Paginator) {
	return origin_s3.NewListObjectsV2Paginator(client,params,optFns...)
 }

type GetObjectRetentionInput = origin_s3.GetObjectRetentionInput

type GetObjectRetentionOutput = origin_s3.GetObjectRetentionOutput

type HeadObjectInput = origin_s3.HeadObjectInput

type HeadObjectOutput = origin_s3.HeadObjectOutput

type HeadObjectAPIClient = origin_s3.HeadObjectAPIClient

type ObjectExistsWaiterOptions = origin_s3.ObjectExistsWaiterOptions

type ObjectExistsWaiter = origin_s3.ObjectExistsWaiter

func NewObjectExistsWaiter(client HeadObjectAPIClient,optFns ...func(*ObjectExistsWaiterOptions) ) (*ObjectExistsWaiter) {
	return origin_s3.NewObjectExistsWaiter(client,optFns...)
 }

type ObjectNotExistsWaiterOptions = origin_s3.ObjectNotExistsWaiterOptions

type ObjectNotExistsWaiter = origin_s3.ObjectNotExistsWaiter

func NewObjectNotExistsWaiter(client HeadObjectAPIClient,optFns ...func(*ObjectNotExistsWaiterOptions) ) (*ObjectNotExistsWaiter) {
	return origin_s3.NewObjectNotExistsWaiter(client,optFns...)
 }

type ListMultipartUploadsInput = origin_s3.ListMultipartUploadsInput

type ListMultipartUploadsOutput = origin_s3.ListMultipartUploadsOutput

type PutBucketInventoryConfigurationInput = origin_s3.PutBucketInventoryConfigurationInput

type PutBucketInventoryConfigurationOutput = origin_s3.PutBucketInventoryConfigurationOutput

type PutBucketRequestPaymentInput = origin_s3.PutBucketRequestPaymentInput

type PutBucketRequestPaymentOutput = origin_s3.PutBucketRequestPaymentOutput

type GetBucketCorsInput = origin_s3.GetBucketCorsInput

type GetBucketCorsOutput = origin_s3.GetBucketCorsOutput

type GetBucketEncryptionInput = origin_s3.GetBucketEncryptionInput

type GetBucketEncryptionOutput = origin_s3.GetBucketEncryptionOutput

type GetBucketLoggingInput = origin_s3.GetBucketLoggingInput

type GetBucketLoggingOutput = origin_s3.GetBucketLoggingOutput

type GetBucketNotificationConfigurationInput = origin_s3.GetBucketNotificationConfigurationInput

type GetBucketNotificationConfigurationOutput = origin_s3.GetBucketNotificationConfigurationOutput

type ListBucketMetricsConfigurationsInput = origin_s3.ListBucketMetricsConfigurationsInput

type ListBucketMetricsConfigurationsOutput = origin_s3.ListBucketMetricsConfigurationsOutput

type PutBucketNotificationConfigurationInput = origin_s3.PutBucketNotificationConfigurationInput

type PutBucketNotificationConfigurationOutput = origin_s3.PutBucketNotificationConfigurationOutput

type DeleteBucketAnalyticsConfigurationInput = origin_s3.DeleteBucketAnalyticsConfigurationInput

type DeleteBucketAnalyticsConfigurationOutput = origin_s3.DeleteBucketAnalyticsConfigurationOutput

type DeleteBucketWebsiteInput = origin_s3.DeleteBucketWebsiteInput

type DeleteBucketWebsiteOutput = origin_s3.DeleteBucketWebsiteOutput

type DeleteObjectInput = origin_s3.DeleteObjectInput

type DeleteObjectOutput = origin_s3.DeleteObjectOutput

type GetBucketAccelerateConfigurationInput = origin_s3.GetBucketAccelerateConfigurationInput

type GetBucketAccelerateConfigurationOutput = origin_s3.GetBucketAccelerateConfigurationOutput

type PutBucketAnalyticsConfigurationInput = origin_s3.PutBucketAnalyticsConfigurationInput

type PutBucketAnalyticsConfigurationOutput = origin_s3.PutBucketAnalyticsConfigurationOutput

type PutBucketIntelligentTieringConfigurationInput = origin_s3.PutBucketIntelligentTieringConfigurationInput

type PutBucketIntelligentTieringConfigurationOutput = origin_s3.PutBucketIntelligentTieringConfigurationOutput

type PutBucketTaggingInput = origin_s3.PutBucketTaggingInput

type PutBucketTaggingOutput = origin_s3.PutBucketTaggingOutput

type UploadPartInput = origin_s3.UploadPartInput

type UploadPartOutput = origin_s3.UploadPartOutput

type CopyObjectInput = origin_s3.CopyObjectInput

type CopyObjectOutput = origin_s3.CopyObjectOutput

type DeleteBucketEncryptionInput = origin_s3.DeleteBucketEncryptionInput

type DeleteBucketEncryptionOutput = origin_s3.DeleteBucketEncryptionOutput

type DeleteBucketReplicationInput = origin_s3.DeleteBucketReplicationInput

type DeleteBucketReplicationOutput = origin_s3.DeleteBucketReplicationOutput

type PutBucketLoggingInput = origin_s3.PutBucketLoggingInput

type PutBucketLoggingOutput = origin_s3.PutBucketLoggingOutput

type AbortMultipartUploadInput = origin_s3.AbortMultipartUploadInput

type AbortMultipartUploadOutput = origin_s3.AbortMultipartUploadOutput

type DeleteBucketTaggingInput = origin_s3.DeleteBucketTaggingInput

type DeleteBucketTaggingOutput = origin_s3.DeleteBucketTaggingOutput

type GetBucketInventoryConfigurationInput = origin_s3.GetBucketInventoryConfigurationInput

type GetBucketInventoryConfigurationOutput = origin_s3.GetBucketInventoryConfigurationOutput

type GetObjectLegalHoldInput = origin_s3.GetObjectLegalHoldInput

type GetObjectLegalHoldOutput = origin_s3.GetObjectLegalHoldOutput

type PutBucketPolicyInput = origin_s3.PutBucketPolicyInput

type PutBucketPolicyOutput = origin_s3.PutBucketPolicyOutput

type PutBucketWebsiteInput = origin_s3.PutBucketWebsiteInput

type PutBucketWebsiteOutput = origin_s3.PutBucketWebsiteOutput

type PutObjectAclInput = origin_s3.PutObjectAclInput

type PutObjectAclOutput = origin_s3.PutObjectAclOutput

type CompleteMultipartUploadInput = origin_s3.CompleteMultipartUploadInput

type CompleteMultipartUploadOutput = origin_s3.CompleteMultipartUploadOutput

type CreateBucketInput = origin_s3.CreateBucketInput

type CreateBucketOutput = origin_s3.CreateBucketOutput

type GetBucketLifecycleConfigurationInput = origin_s3.GetBucketLifecycleConfigurationInput

type GetBucketLifecycleConfigurationOutput = origin_s3.GetBucketLifecycleConfigurationOutput

type DeleteBucketOwnershipControlsInput = origin_s3.DeleteBucketOwnershipControlsInput

type DeleteBucketOwnershipControlsOutput = origin_s3.DeleteBucketOwnershipControlsOutput

type GetBucketVersioningInput = origin_s3.GetBucketVersioningInput

type GetBucketVersioningOutput = origin_s3.GetBucketVersioningOutput

type GetObjectAttributesInput = origin_s3.GetObjectAttributesInput

type GetObjectAttributesOutput = origin_s3.GetObjectAttributesOutput

type GetBucketRequestPaymentInput = origin_s3.GetBucketRequestPaymentInput

type GetBucketRequestPaymentOutput = origin_s3.GetBucketRequestPaymentOutput

type GetObjectLockConfigurationInput = origin_s3.GetObjectLockConfigurationInput

type GetObjectLockConfigurationOutput = origin_s3.GetObjectLockConfigurationOutput

type HeadBucketInput = origin_s3.HeadBucketInput

type HeadBucketOutput = origin_s3.HeadBucketOutput

type HeadBucketAPIClient = origin_s3.HeadBucketAPIClient

type BucketExistsWaiterOptions = origin_s3.BucketExistsWaiterOptions

type BucketExistsWaiter = origin_s3.BucketExistsWaiter

func NewBucketExistsWaiter(client HeadBucketAPIClient,optFns ...func(*BucketExistsWaiterOptions) ) (*BucketExistsWaiter) {
	return origin_s3.NewBucketExistsWaiter(client,optFns...)
 }

type BucketNotExistsWaiterOptions = origin_s3.BucketNotExistsWaiterOptions

type BucketNotExistsWaiter = origin_s3.BucketNotExistsWaiter

func NewBucketNotExistsWaiter(client HeadBucketAPIClient,optFns ...func(*BucketNotExistsWaiterOptions) ) (*BucketNotExistsWaiter) {
	return origin_s3.NewBucketNotExistsWaiter(client,optFns...)
 }

type ListBucketIntelligentTieringConfigurationsInput = origin_s3.ListBucketIntelligentTieringConfigurationsInput

type ListBucketIntelligentTieringConfigurationsOutput = origin_s3.ListBucketIntelligentTieringConfigurationsOutput

type PutBucketCorsInput = origin_s3.PutBucketCorsInput

type PutBucketCorsOutput = origin_s3.PutBucketCorsOutput

type DeleteBucketIntelligentTieringConfigurationInput = origin_s3.DeleteBucketIntelligentTieringConfigurationInput

type DeleteBucketIntelligentTieringConfigurationOutput = origin_s3.DeleteBucketIntelligentTieringConfigurationOutput

type GetBucketAnalyticsConfigurationInput = origin_s3.GetBucketAnalyticsConfigurationInput

type GetBucketAnalyticsConfigurationOutput = origin_s3.GetBucketAnalyticsConfigurationOutput

type GetBucketIntelligentTieringConfigurationInput = origin_s3.GetBucketIntelligentTieringConfigurationInput

type GetBucketIntelligentTieringConfigurationOutput = origin_s3.GetBucketIntelligentTieringConfigurationOutput

type PutObjectLegalHoldInput = origin_s3.PutObjectLegalHoldInput

type PutObjectLegalHoldOutput = origin_s3.PutObjectLegalHoldOutput

type ListObjectsInput = origin_s3.ListObjectsInput

type ListObjectsOutput = origin_s3.ListObjectsOutput

type PutBucketEncryptionInput = origin_s3.PutBucketEncryptionInput

type PutBucketEncryptionOutput = origin_s3.PutBucketEncryptionOutput

type SelectObjectContentInput = origin_s3.SelectObjectContentInput

type SelectObjectContentOutput = origin_s3.SelectObjectContentOutput

type SelectObjectContentEventStream = origin_s3.SelectObjectContentEventStream

func NewSelectObjectContentEventStream(optFns ...func(*SelectObjectContentEventStream) ) (*SelectObjectContentEventStream) {
	return origin_s3.NewSelectObjectContentEventStream(optFns...)
 }

type WriteGetObjectResponseInput = origin_s3.WriteGetObjectResponseInput

type WriteGetObjectResponseOutput = origin_s3.WriteGetObjectResponseOutput

type DeleteBucketMetricsConfigurationInput = origin_s3.DeleteBucketMetricsConfigurationInput

type DeleteBucketMetricsConfigurationOutput = origin_s3.DeleteBucketMetricsConfigurationOutput

type GetObjectTaggingInput = origin_s3.GetObjectTaggingInput

type GetObjectTaggingOutput = origin_s3.GetObjectTaggingOutput

type PutBucketMetricsConfigurationInput = origin_s3.PutBucketMetricsConfigurationInput

type PutBucketMetricsConfigurationOutput = origin_s3.PutBucketMetricsConfigurationOutput

type DeleteBucketCorsInput = origin_s3.DeleteBucketCorsInput

type DeleteBucketCorsOutput = origin_s3.DeleteBucketCorsOutput

type GetBucketPolicyInput = origin_s3.GetBucketPolicyInput

type GetBucketPolicyOutput = origin_s3.GetBucketPolicyOutput

type PutObjectInput = origin_s3.PutObjectInput

type PutObjectOutput = origin_s3.PutObjectOutput

type GetObjectInput = origin_s3.GetObjectInput

type GetObjectOutput = origin_s3.GetObjectOutput

type ListBucketInventoryConfigurationsInput = origin_s3.ListBucketInventoryConfigurationsInput

type ListBucketInventoryConfigurationsOutput = origin_s3.ListBucketInventoryConfigurationsOutput

type PutObjectRetentionInput = origin_s3.PutObjectRetentionInput

type PutObjectRetentionOutput = origin_s3.PutObjectRetentionOutput

type EndpointResolverOptions = origin_s3.EndpointResolverOptions

type EndpointResolver = origin_s3.EndpointResolver

func NewDefaultEndpointResolver() (interface {}) {
	return origin_s3.NewDefaultEndpointResolver()
 }

type EndpointResolverFunc = origin_s3.EndpointResolverFunc

func EndpointResolverFromURL(url string,optFns ...func(*aws.Endpoint) ) (EndpointResolver) {
	return origin_s3.EndpointResolverFromURL(url,optFns...)
 }

type ResolveEndpoint = origin_s3.ResolveEndpoint

type CreateMultipartUploadInput = origin_s3.CreateMultipartUploadInput

type CreateMultipartUploadOutput = origin_s3.CreateMultipartUploadOutput

type DeleteBucketInput = origin_s3.DeleteBucketInput

type DeleteBucketOutput = origin_s3.DeleteBucketOutput

type GetBucketReplicationInput = origin_s3.GetBucketReplicationInput

type GetBucketReplicationOutput = origin_s3.GetBucketReplicationOutput

