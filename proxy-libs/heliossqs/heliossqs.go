package heliossqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	origin_sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/helios/opentelemetry-go-contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

type sqsMessageCarrier struct {
	messageAttrs map[string]types.MessageAttributeValue
}

func (c sqsMessageCarrier) Get(key string) string {
	if c.messageAttrs == nil {
		return ""
	}

	for attrKey, val := range c.messageAttrs {
		if attrKey == key {
			return *val.StringValue
		}
	}

	return ""
}

func (c sqsMessageCarrier) Set(key, val string) {
	dataType := "String"
	c.messageAttrs[key] = types.MessageAttributeValue{DataType: &dataType, StringValue: &val}
}

func (c sqsMessageCarrier) Keys() []string {
	result := []string{}
	for key := range c.messageAttrs {
		result = append(result, key)
	}

	return result
}

const maxSqsAttributesCount = 10

func attributeSetter(ctx context.Context, ii middleware.InitializeInput) []attribute.KeyValue {
	result := []attribute.KeyValue{}
	switch castParams := ii.Parameters.(type) {
	case *origin_sqs.SendMessageInput:
		{
			attrs := castParams.MessageAttributes
			if attrs == nil {
				attrs = map[string]types.MessageAttributeValue{}
				castParams.MessageAttributes = attrs
			}

			if len(attrs) < maxSqsAttributesCount {
				carrier := sqsMessageCarrier{attrs}
				otel.GetTextMapPropagator().Inject(ctx, carrier)
			}

			result = append(result, attribute.KeyValue{Key: "messaging.payload", Value: attribute.StringValue(*castParams.MessageBody)})
		}
	case *origin_sqs.SendMessageBatchInput:
		{
			entries := castParams.Entries
			for index := range entries {
				entry := &entries[index]
				attrs := entry.MessageAttributes
				if attrs == nil {
					attrs = map[string]types.MessageAttributeValue{}
					entry.MessageAttributes = attrs
				}

				if len(attrs) < maxSqsAttributesCount {
					carrier := sqsMessageCarrier{attrs}
					otel.GetTextMapPropagator().Inject(ctx, carrier)
				}

				carrier := sqsMessageCarrier{attrs}
				otel.GetTextMapPropagator().Inject(ctx, carrier)
			}
		}
	case *origin_sqs.ReceiveMessageInput:
		{
			attrNames := castParams.MessageAttributeNames
			if attrNames == nil {
				attrNames = []string{}
			}

			attrNames = append(attrNames, "traceparent", "baggage")
			castParams.MessageAttributeNames = attrNames
		}
	}
	return result
}

func New(options Options, optFns ...func(*Options)) *Client {
	attributeSetterOpt := otelaws.WithAttributeSetter(attributeSetter)
	otelaws.AppendMiddlewares(&options.APIOptions, attributeSetterOpt)
	return origin_sqs.New(options, optFns...)
}

type Options = origin_sqs.Options

func WithAPIOptions(optFns ...func(*middleware.Stack) error) func(*Options) {
	return origin_sqs.WithAPIOptions(optFns...)
}

func WithEndpointResolver(v EndpointResolver) func(*Options) {
	return origin_sqs.WithEndpointResolver(v)
}

type HTTPClient = origin_sqs.HTTPClient

func NewFromConfig(cfg aws.Config, optFns ...func(*Options)) *Client {
	attributeSetterOpt := otelaws.WithAttributeSetter(attributeSetter)
	otelaws.AppendMiddlewares(&cfg.APIOptions, attributeSetterOpt)
	return origin_sqs.NewFromConfig(cfg, optFns...)
}

type HTTPSignerV4 = origin_sqs.HTTPSignerV4

type SendMessageInput = origin_sqs.SendMessageInput

type SendMessageOutput = origin_sqs.SendMessageOutput

type SetQueueAttributesInput = origin_sqs.SetQueueAttributesInput

type SetQueueAttributesOutput = origin_sqs.SetQueueAttributesOutput

type EndpointResolverOptions = origin_sqs.EndpointResolverOptions

type EndpointResolver = origin_sqs.EndpointResolver

func NewDefaultEndpointResolver() interface{} {
	return origin_sqs.NewDefaultEndpointResolver()
}

type EndpointResolverFunc = origin_sqs.EndpointResolverFunc

func EndpointResolverFromURL(url string, optFns ...func(*aws.Endpoint)) EndpointResolver {
	return origin_sqs.EndpointResolverFromURL(url, optFns...)
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

func NewListQueuesPaginator(client ListQueuesAPIClient, params *ListQueuesInput, optFns ...func(*ListQueuesPaginatorOptions)) *ListQueuesPaginator {
	return origin_sqs.NewListQueuesPaginator(client, params, optFns...)
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

func NewListDeadLetterSourceQueuesPaginator(client ListDeadLetterSourceQueuesAPIClient, params *ListDeadLetterSourceQueuesInput, optFns ...func(*ListDeadLetterSourceQueuesPaginatorOptions)) *ListDeadLetterSourceQueuesPaginator {
	return origin_sqs.NewListDeadLetterSourceQueuesPaginator(client, params, optFns...)
}

type ChangeMessageVisibilityBatchInput = origin_sqs.ChangeMessageVisibilityBatchInput

type ChangeMessageVisibilityBatchOutput = origin_sqs.ChangeMessageVisibilityBatchOutput

type DeleteMessageInput = origin_sqs.DeleteMessageInput

type DeleteMessageOutput = origin_sqs.DeleteMessageOutput
