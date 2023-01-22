package heliosdynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	origin_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)
type BatchExecuteStatementInput = origin_dynamodb.BatchExecuteStatementInput

type BatchExecuteStatementOutput = origin_dynamodb.BatchExecuteStatementOutput

type DeleteItemInput = origin_dynamodb.DeleteItemInput

type DeleteItemOutput = origin_dynamodb.DeleteItemOutput

type DescribeGlobalTableInput = origin_dynamodb.DescribeGlobalTableInput

type DescribeGlobalTableOutput = origin_dynamodb.DescribeGlobalTableOutput

const ServiceID = origin_dynamodb.ServiceID

const ServiceAPIVersion = origin_dynamodb.ServiceAPIVersion

type Client = origin_dynamodb.Client

func New(options Options,optFns ...func(*Options) ) (*Client) {
	otelaws.AppendMiddlewares(&options.APIOptions)
	return origin_dynamodb.New(options,optFns...)
 }

type Options = origin_dynamodb.Options

func WithAPIOptions(optFns ...func(*middleware.Stack) error) (func(*Options) ) {
	return origin_dynamodb.WithAPIOptions(optFns...)
 }

func WithEndpointResolver(v EndpointResolver) (func(*Options) ) {
	return origin_dynamodb.WithEndpointResolver(v)
 }

type HTTPClient = origin_dynamodb.HTTPClient

func NewFromConfig(cfg aws.Config,optFns ...func(*Options) ) (*Client) {
	otelaws.AppendMiddlewares(&cfg.APIOptions)
	return origin_dynamodb.NewFromConfig(cfg,optFns...)
 }

type HTTPSignerV4 = origin_dynamodb.HTTPSignerV4

type EndpointDiscoveryOptions = origin_dynamodb.EndpointDiscoveryOptions

type IdempotencyTokenProvider = origin_dynamodb.IdempotencyTokenProvider

type DescribeContinuousBackupsInput = origin_dynamodb.DescribeContinuousBackupsInput

type DescribeContinuousBackupsOutput = origin_dynamodb.DescribeContinuousBackupsOutput

type ListBackupsInput = origin_dynamodb.ListBackupsInput

type ListBackupsOutput = origin_dynamodb.ListBackupsOutput

type ScanInput = origin_dynamodb.ScanInput

type ScanOutput = origin_dynamodb.ScanOutput

type ScanAPIClient = origin_dynamodb.ScanAPIClient

type ScanPaginatorOptions = origin_dynamodb.ScanPaginatorOptions

type ScanPaginator = origin_dynamodb.ScanPaginator

func NewScanPaginator(client ScanAPIClient,params *ScanInput,optFns ...func(*ScanPaginatorOptions) ) (*ScanPaginator) {
	return origin_dynamodb.NewScanPaginator(client,params,optFns...)
 }

type UpdateTableInput = origin_dynamodb.UpdateTableInput

type UpdateTableOutput = origin_dynamodb.UpdateTableOutput

type DescribeGlobalTableSettingsInput = origin_dynamodb.DescribeGlobalTableSettingsInput

type DescribeGlobalTableSettingsOutput = origin_dynamodb.DescribeGlobalTableSettingsOutput

type DescribeImportInput = origin_dynamodb.DescribeImportInput

type DescribeImportOutput = origin_dynamodb.DescribeImportOutput

type ExecuteStatementInput = origin_dynamodb.ExecuteStatementInput

type ExecuteStatementOutput = origin_dynamodb.ExecuteStatementOutput

type UpdateContinuousBackupsInput = origin_dynamodb.UpdateContinuousBackupsInput

type UpdateContinuousBackupsOutput = origin_dynamodb.UpdateContinuousBackupsOutput

type CreateGlobalTableInput = origin_dynamodb.CreateGlobalTableInput

type CreateGlobalTableOutput = origin_dynamodb.CreateGlobalTableOutput

type DescribeKinesisStreamingDestinationInput = origin_dynamodb.DescribeKinesisStreamingDestinationInput

type DescribeKinesisStreamingDestinationOutput = origin_dynamodb.DescribeKinesisStreamingDestinationOutput

type DescribeTableReplicaAutoScalingInput = origin_dynamodb.DescribeTableReplicaAutoScalingInput

type DescribeTableReplicaAutoScalingOutput = origin_dynamodb.DescribeTableReplicaAutoScalingOutput

type ListExportsInput = origin_dynamodb.ListExportsInput

type ListExportsOutput = origin_dynamodb.ListExportsOutput

type ListExportsAPIClient = origin_dynamodb.ListExportsAPIClient

type ListExportsPaginatorOptions = origin_dynamodb.ListExportsPaginatorOptions

type ListExportsPaginator = origin_dynamodb.ListExportsPaginator

func NewListExportsPaginator(client ListExportsAPIClient,params *ListExportsInput,optFns ...func(*ListExportsPaginatorOptions) ) (*ListExportsPaginator) {
	return origin_dynamodb.NewListExportsPaginator(client,params,optFns...)
 }

type BatchGetItemInput = origin_dynamodb.BatchGetItemInput

type BatchGetItemOutput = origin_dynamodb.BatchGetItemOutput

type CreateBackupInput = origin_dynamodb.CreateBackupInput

type CreateBackupOutput = origin_dynamodb.CreateBackupOutput

type ExecuteTransactionInput = origin_dynamodb.ExecuteTransactionInput

type ExecuteTransactionOutput = origin_dynamodb.ExecuteTransactionOutput

type ListGlobalTablesInput = origin_dynamodb.ListGlobalTablesInput

type ListGlobalTablesOutput = origin_dynamodb.ListGlobalTablesOutput

type UpdateTableReplicaAutoScalingInput = origin_dynamodb.UpdateTableReplicaAutoScalingInput

type UpdateTableReplicaAutoScalingOutput = origin_dynamodb.UpdateTableReplicaAutoScalingOutput

type CreateTableInput = origin_dynamodb.CreateTableInput

type CreateTableOutput = origin_dynamodb.CreateTableOutput

type DescribeBackupInput = origin_dynamodb.DescribeBackupInput

type DescribeBackupOutput = origin_dynamodb.DescribeBackupOutput

type ListContributorInsightsInput = origin_dynamodb.ListContributorInsightsInput

type ListContributorInsightsOutput = origin_dynamodb.ListContributorInsightsOutput

type ListContributorInsightsAPIClient = origin_dynamodb.ListContributorInsightsAPIClient

type ListContributorInsightsPaginatorOptions = origin_dynamodb.ListContributorInsightsPaginatorOptions

type ListContributorInsightsPaginator = origin_dynamodb.ListContributorInsightsPaginator

func NewListContributorInsightsPaginator(client ListContributorInsightsAPIClient,params *ListContributorInsightsInput,optFns ...func(*ListContributorInsightsPaginatorOptions) ) (*ListContributorInsightsPaginator) {
	return origin_dynamodb.NewListContributorInsightsPaginator(client,params,optFns...)
 }

type UpdateContributorInsightsInput = origin_dynamodb.UpdateContributorInsightsInput

type UpdateContributorInsightsOutput = origin_dynamodb.UpdateContributorInsightsOutput

type UpdateItemInput = origin_dynamodb.UpdateItemInput

type UpdateItemOutput = origin_dynamodb.UpdateItemOutput

type GetItemInput = origin_dynamodb.GetItemInput

type GetItemOutput = origin_dynamodb.GetItemOutput

type ListTagsOfResourceInput = origin_dynamodb.ListTagsOfResourceInput

type ListTagsOfResourceOutput = origin_dynamodb.ListTagsOfResourceOutput

type DescribeLimitsInput = origin_dynamodb.DescribeLimitsInput

type DescribeLimitsOutput = origin_dynamodb.DescribeLimitsOutput

type TagResourceInput = origin_dynamodb.TagResourceInput

type TagResourceOutput = origin_dynamodb.TagResourceOutput

type DeleteTableInput = origin_dynamodb.DeleteTableInput

type DeleteTableOutput = origin_dynamodb.DeleteTableOutput

type UpdateGlobalTableInput = origin_dynamodb.UpdateGlobalTableInput

type UpdateGlobalTableOutput = origin_dynamodb.UpdateGlobalTableOutput

type DeleteBackupInput = origin_dynamodb.DeleteBackupInput

type DeleteBackupOutput = origin_dynamodb.DeleteBackupOutput

type RestoreTableFromBackupInput = origin_dynamodb.RestoreTableFromBackupInput

type RestoreTableFromBackupOutput = origin_dynamodb.RestoreTableFromBackupOutput

type QueryInput = origin_dynamodb.QueryInput

type QueryOutput = origin_dynamodb.QueryOutput

type QueryAPIClient = origin_dynamodb.QueryAPIClient

type QueryPaginatorOptions = origin_dynamodb.QueryPaginatorOptions

type QueryPaginator = origin_dynamodb.QueryPaginator

func NewQueryPaginator(client QueryAPIClient,params *QueryInput,optFns ...func(*QueryPaginatorOptions) ) (*QueryPaginator) {
	return origin_dynamodb.NewQueryPaginator(client,params,optFns...)
 }

type RestoreTableToPointInTimeInput = origin_dynamodb.RestoreTableToPointInTimeInput

type RestoreTableToPointInTimeOutput = origin_dynamodb.RestoreTableToPointInTimeOutput

type DescribeContributorInsightsInput = origin_dynamodb.DescribeContributorInsightsInput

type DescribeContributorInsightsOutput = origin_dynamodb.DescribeContributorInsightsOutput

type DescribeEndpointsInput = origin_dynamodb.DescribeEndpointsInput

type DescribeEndpointsOutput = origin_dynamodb.DescribeEndpointsOutput

type DescribeTimeToLiveInput = origin_dynamodb.DescribeTimeToLiveInput

type DescribeTimeToLiveOutput = origin_dynamodb.DescribeTimeToLiveOutput

type DisableKinesisStreamingDestinationInput = origin_dynamodb.DisableKinesisStreamingDestinationInput

type DisableKinesisStreamingDestinationOutput = origin_dynamodb.DisableKinesisStreamingDestinationOutput

type ListTablesInput = origin_dynamodb.ListTablesInput

type ListTablesOutput = origin_dynamodb.ListTablesOutput

type ListTablesAPIClient = origin_dynamodb.ListTablesAPIClient

type ListTablesPaginatorOptions = origin_dynamodb.ListTablesPaginatorOptions

type ListTablesPaginator = origin_dynamodb.ListTablesPaginator

func NewListTablesPaginator(client ListTablesAPIClient,params *ListTablesInput,optFns ...func(*ListTablesPaginatorOptions) ) (*ListTablesPaginator) {
	return origin_dynamodb.NewListTablesPaginator(client,params,optFns...)
 }

type TransactGetItemsInput = origin_dynamodb.TransactGetItemsInput

type TransactGetItemsOutput = origin_dynamodb.TransactGetItemsOutput

type EndpointResolverOptions = origin_dynamodb.EndpointResolverOptions

type EndpointResolver = origin_dynamodb.EndpointResolver

func NewDefaultEndpointResolver() (interface {}) {
	return origin_dynamodb.NewDefaultEndpointResolver()
 }

type EndpointResolverFunc = origin_dynamodb.EndpointResolverFunc

func EndpointResolverFromURL(url string,optFns ...func(*aws.Endpoint) ) (EndpointResolver) {
	return origin_dynamodb.EndpointResolverFromURL(url,optFns...)
 }

type ResolveEndpoint = origin_dynamodb.ResolveEndpoint

type BatchWriteItemInput = origin_dynamodb.BatchWriteItemInput

type BatchWriteItemOutput = origin_dynamodb.BatchWriteItemOutput

type DescribeExportInput = origin_dynamodb.DescribeExportInput

type DescribeExportOutput = origin_dynamodb.DescribeExportOutput

type DescribeTableInput = origin_dynamodb.DescribeTableInput

type DescribeTableOutput = origin_dynamodb.DescribeTableOutput

type DescribeTableAPIClient = origin_dynamodb.DescribeTableAPIClient

type TableExistsWaiterOptions = origin_dynamodb.TableExistsWaiterOptions

type TableExistsWaiter = origin_dynamodb.TableExistsWaiter

func NewTableExistsWaiter(client DescribeTableAPIClient,optFns ...func(*TableExistsWaiterOptions) ) (*TableExistsWaiter) {
	return origin_dynamodb.NewTableExistsWaiter(client,optFns...)
 }

type TableNotExistsWaiterOptions = origin_dynamodb.TableNotExistsWaiterOptions

type TableNotExistsWaiter = origin_dynamodb.TableNotExistsWaiter

func NewTableNotExistsWaiter(client DescribeTableAPIClient,optFns ...func(*TableNotExistsWaiterOptions) ) (*TableNotExistsWaiter) {
	return origin_dynamodb.NewTableNotExistsWaiter(client,optFns...)
 }

type ExportTableToPointInTimeInput = origin_dynamodb.ExportTableToPointInTimeInput

type ExportTableToPointInTimeOutput = origin_dynamodb.ExportTableToPointInTimeOutput

type PutItemInput = origin_dynamodb.PutItemInput

type PutItemOutput = origin_dynamodb.PutItemOutput

type ListImportsInput = origin_dynamodb.ListImportsInput

type ListImportsOutput = origin_dynamodb.ListImportsOutput

type ListImportsAPIClient = origin_dynamodb.ListImportsAPIClient

type ListImportsPaginatorOptions = origin_dynamodb.ListImportsPaginatorOptions

type ListImportsPaginator = origin_dynamodb.ListImportsPaginator

func NewListImportsPaginator(client ListImportsAPIClient,params *ListImportsInput,optFns ...func(*ListImportsPaginatorOptions) ) (*ListImportsPaginator) {
	return origin_dynamodb.NewListImportsPaginator(client,params,optFns...)
 }

type UpdateTimeToLiveInput = origin_dynamodb.UpdateTimeToLiveInput

type UpdateTimeToLiveOutput = origin_dynamodb.UpdateTimeToLiveOutput

type EnableKinesisStreamingDestinationInput = origin_dynamodb.EnableKinesisStreamingDestinationInput

type EnableKinesisStreamingDestinationOutput = origin_dynamodb.EnableKinesisStreamingDestinationOutput

type ImportTableInput = origin_dynamodb.ImportTableInput

type ImportTableOutput = origin_dynamodb.ImportTableOutput

type TransactWriteItemsInput = origin_dynamodb.TransactWriteItemsInput

type TransactWriteItemsOutput = origin_dynamodb.TransactWriteItemsOutput

type UntagResourceInput = origin_dynamodb.UntagResourceInput

type UntagResourceOutput = origin_dynamodb.UntagResourceOutput

type UpdateGlobalTableSettingsInput = origin_dynamodb.UpdateGlobalTableSettingsInput

type UpdateGlobalTableSettingsOutput = origin_dynamodb.UpdateGlobalTableSettingsOutput

