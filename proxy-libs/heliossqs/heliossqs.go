package heliossqs

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	origin_sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/smithy-go/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)
type DeleteMessageBatchInput = origin_sqs.DeleteMessageBatchInput

type DeleteMessageBatchOutput = origin_sqs.DeleteMessageBatchOutput

type RemovePermissionInput = origin_sqs.RemovePermissionInput

type RemovePermissionOutput = origin_sqs.RemovePermissionOutput

type DeleteQueueInput = origin_sqs.DeleteQueueInput

type DeleteQueueOutput = origin_sqs.DeleteQueueOutput

type SendMessageBatchInput = origin_sqs.SendMessageBatchInput

type SendMessageBatchOutput = origin_sqs.SendMessageBatchOutput

const ServiceID = origin_sqs.ServiceID

const ServiceAPIVersion = origin_sqs.ServiceAPIVersion

type Client = origin_sqs.Client

func New(options Options,optFns ...func(*Options) ) (*Client) {
	otelaws.AppendMiddlewares(&options.APIOptions)
	return origin_sqs.New(options,optFns...)
 }

type Options = origin_sqs.Options

func WithAPIOptions(optFns ...func(*middleware.Stack) error) (func(*Options) ) {
	return origin_sqs.WithAPIOptions(optFns...)
 }

func WithEndpointResolver(v EndpointResolver) (func(*Options) ) {
	return origin_sqs.WithEndpointResolver(v)
 }

type HTTPClient = origin_sqs.HTTPClient

func NewFromConfig(cfg aws.Config,optFns ...func(*Options) ) (*Client) {
	otelaws.AppendMiddlewares(&cfg.APIOptions)
	return origin_sqs.NewFromConfig(cfg,optFns...)
 }

type HTTPSignerV4 = origin_sqs.HTTPSignerV4

type SendMessageInput = origin_sqs.SendMessageInput

type SendMessageOutput = origin_sqs.SendMessageOutput

type SetQueueAttributesInput = origin_sqs.SetQueueAttributesInput

type SetQueueAttributesOutput = origin_sqs.SetQueueAttributesOutput

type EndpointResolverOptions = origin_sqs.EndpointResolverOptions

type EndpointResolver = origin_sqs.EndpointResolver

func NewDefaultEndpointResolver() (interface {}) {
	return origin_sqs.NewDefaultEndpointResolver()
 }

type EndpointResolverFunc = origin_sqs.EndpointResolverFunc

func EndpointResolverFromURL(url string,optFns ...func(*aws.Endpoint) ) (EndpointResolver) {
	return origin_sqs.EndpointResolverFromURL(url,optFns...)
 }

type ResolveEndpoint = origin_sqs.ResolveEndpoint

type UntagQueueInput = origin_sqs.UntagQueueInput

type UntagQueueOutput = origin_sqs.UntagQueueOutput

type ChangeMessageVisibilityInput = origin_sqs.ChangeMessageVisibilityInput

type ChangeMessageVisibilityOutput = origin_sqs.ChangeMessageVisibilityOutput

type GetQueueAttributesInput = origin_sqs.GetQueueAttributesInput

type GetQueueAttributesOutput = origin_sqs.GetQueueAttributesOutput

type ListQueueTagsInput = origin_sqs.ListQueueTagsInput

type ListQueueTagsOutput = origin_sqs.ListQueueTagsOutput

type TagQueueInput = origin_sqs.TagQueueInput

type TagQueueOutput = origin_sqs.TagQueueOutput

type CreateQueueInput = origin_sqs.CreateQueueInput

type CreateQueueOutput = origin_sqs.CreateQueueOutput

type GetQueueUrlInput = origin_sqs.GetQueueUrlInput

type GetQueueUrlOutput = origin_sqs.GetQueueUrlOutput

type ListQueuesInput = origin_sqs.ListQueuesInput

type ListQueuesOutput = origin_sqs.ListQueuesOutput

type ListQueuesAPIClient = origin_sqs.ListQueuesAPIClient

type ListQueuesPaginatorOptions = origin_sqs.ListQueuesPaginatorOptions

type ListQueuesPaginator = origin_sqs.ListQueuesPaginator

func NewListQueuesPaginator(client ListQueuesAPIClient,params *ListQueuesInput,optFns ...func(*ListQueuesPaginatorOptions) ) (*ListQueuesPaginator) {
	return origin_sqs.NewListQueuesPaginator(client,params,optFns...)
 }

type PurgeQueueInput = origin_sqs.PurgeQueueInput

type PurgeQueueOutput = origin_sqs.PurgeQueueOutput

type ReceiveMessageInput = origin_sqs.ReceiveMessageInput

type ReceiveMessageOutput = origin_sqs.ReceiveMessageOutput

type AddPermissionInput = origin_sqs.AddPermissionInput

type AddPermissionOutput = origin_sqs.AddPermissionOutput

type ListDeadLetterSourceQueuesInput = origin_sqs.ListDeadLetterSourceQueuesInput

type ListDeadLetterSourceQueuesOutput = origin_sqs.ListDeadLetterSourceQueuesOutput

type ListDeadLetterSourceQueuesAPIClient = origin_sqs.ListDeadLetterSourceQueuesAPIClient

type ListDeadLetterSourceQueuesPaginatorOptions = origin_sqs.ListDeadLetterSourceQueuesPaginatorOptions

type ListDeadLetterSourceQueuesPaginator = origin_sqs.ListDeadLetterSourceQueuesPaginator

func NewListDeadLetterSourceQueuesPaginator(client ListDeadLetterSourceQueuesAPIClient,params *ListDeadLetterSourceQueuesInput,optFns ...func(*ListDeadLetterSourceQueuesPaginatorOptions) ) (*ListDeadLetterSourceQueuesPaginator) {
	return origin_sqs.NewListDeadLetterSourceQueuesPaginator(client,params,optFns...)
 }

type ChangeMessageVisibilityBatchInput = origin_sqs.ChangeMessageVisibilityBatchInput

type ChangeMessageVisibilityBatchOutput = origin_sqs.ChangeMessageVisibilityBatchOutput

type DeleteMessageInput = origin_sqs.DeleteMessageInput

type DeleteMessageOutput = origin_sqs.DeleteMessageOutput

