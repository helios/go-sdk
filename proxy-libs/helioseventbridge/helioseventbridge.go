package helioseventbridge

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	origin_eventbridge "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/smithy-go/middleware"
	"github.com/helios/opentelemetry-go-contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/otel/attribute"
)

type DescribeConnectionInput = origin_eventbridge.DescribeConnectionInput

type DescribeConnectionOutput = origin_eventbridge.DescribeConnectionOutput

type PutRuleInput = origin_eventbridge.PutRuleInput

type PutRuleOutput = origin_eventbridge.PutRuleOutput

type CancelReplayInput = origin_eventbridge.CancelReplayInput

type CancelReplayOutput = origin_eventbridge.CancelReplayOutput

type DescribePartnerEventSourceInput = origin_eventbridge.DescribePartnerEventSourceInput

type DescribePartnerEventSourceOutput = origin_eventbridge.DescribePartnerEventSourceOutput

type ListEventSourcesInput = origin_eventbridge.ListEventSourcesInput

type ListEventSourcesOutput = origin_eventbridge.ListEventSourcesOutput

type TestEventPatternInput = origin_eventbridge.TestEventPatternInput

type TestEventPatternOutput = origin_eventbridge.TestEventPatternOutput

type DeleteEventBusInput = origin_eventbridge.DeleteEventBusInput

type DeleteEventBusOutput = origin_eventbridge.DeleteEventBusOutput

type DeleteArchiveInput = origin_eventbridge.DeleteArchiveInput

type DeleteArchiveOutput = origin_eventbridge.DeleteArchiveOutput

type ListPartnerEventSourceAccountsInput = origin_eventbridge.ListPartnerEventSourceAccountsInput

type ListPartnerEventSourceAccountsOutput = origin_eventbridge.ListPartnerEventSourceAccountsOutput

type ListTargetsByRuleInput = origin_eventbridge.ListTargetsByRuleInput

type ListTargetsByRuleOutput = origin_eventbridge.ListTargetsByRuleOutput

type PutPartnerEventsInput = origin_eventbridge.PutPartnerEventsInput

type PutPartnerEventsOutput = origin_eventbridge.PutPartnerEventsOutput

type ListEndpointsInput = origin_eventbridge.ListEndpointsInput

type ListEndpointsOutput = origin_eventbridge.ListEndpointsOutput

type PutTargetsInput = origin_eventbridge.PutTargetsInput

type PutTargetsOutput = origin_eventbridge.PutTargetsOutput

type DeletePartnerEventSourceInput = origin_eventbridge.DeletePartnerEventSourceInput

type DeletePartnerEventSourceOutput = origin_eventbridge.DeletePartnerEventSourceOutput

type ListConnectionsInput = origin_eventbridge.ListConnectionsInput

type ListConnectionsOutput = origin_eventbridge.ListConnectionsOutput

type ListRuleNamesByTargetInput = origin_eventbridge.ListRuleNamesByTargetInput

type ListRuleNamesByTargetOutput = origin_eventbridge.ListRuleNamesByTargetOutput

type EndpointResolverOptions = origin_eventbridge.EndpointResolverOptions

type EndpointResolver = origin_eventbridge.EndpointResolver

func NewDefaultEndpointResolver() interface{} {
	return origin_eventbridge.NewDefaultEndpointResolver()
}

type EndpointResolverFunc = origin_eventbridge.EndpointResolverFunc

func EndpointResolverFromURL(url string, optFns ...func(*aws.Endpoint)) EndpointResolver {
	return origin_eventbridge.EndpointResolverFromURL(url, optFns...)
}

type ResolveEndpoint = origin_eventbridge.ResolveEndpoint

type CreateArchiveInput = origin_eventbridge.CreateArchiveInput

type CreateArchiveOutput = origin_eventbridge.CreateArchiveOutput

type ListApiDestinationsInput = origin_eventbridge.ListApiDestinationsInput

type ListApiDestinationsOutput = origin_eventbridge.ListApiDestinationsOutput

type PutEventsInput = origin_eventbridge.PutEventsInput

type PutEventsOutput = origin_eventbridge.PutEventsOutput

type DescribeEventSourceInput = origin_eventbridge.DescribeEventSourceInput

type DescribeEventSourceOutput = origin_eventbridge.DescribeEventSourceOutput

type UpdateEndpointInput = origin_eventbridge.UpdateEndpointInput

type UpdateEndpointOutput = origin_eventbridge.UpdateEndpointOutput

type DisableRuleInput = origin_eventbridge.DisableRuleInput

type DisableRuleOutput = origin_eventbridge.DisableRuleOutput

type RemoveTargetsInput = origin_eventbridge.RemoveTargetsInput

type RemoveTargetsOutput = origin_eventbridge.RemoveTargetsOutput

type DeleteConnectionInput = origin_eventbridge.DeleteConnectionInput

type DeleteConnectionOutput = origin_eventbridge.DeleteConnectionOutput

type DescribeReplayInput = origin_eventbridge.DescribeReplayInput

type DescribeReplayOutput = origin_eventbridge.DescribeReplayOutput

type PutPermissionInput = origin_eventbridge.PutPermissionInput

type PutPermissionOutput = origin_eventbridge.PutPermissionOutput

type DeleteApiDestinationInput = origin_eventbridge.DeleteApiDestinationInput

type DeleteApiDestinationOutput = origin_eventbridge.DeleteApiDestinationOutput

type StartReplayInput = origin_eventbridge.StartReplayInput

type StartReplayOutput = origin_eventbridge.StartReplayOutput

type CreatePartnerEventSourceInput = origin_eventbridge.CreatePartnerEventSourceInput

type CreatePartnerEventSourceOutput = origin_eventbridge.CreatePartnerEventSourceOutput

type DeactivateEventSourceInput = origin_eventbridge.DeactivateEventSourceInput

type DeactivateEventSourceOutput = origin_eventbridge.DeactivateEventSourceOutput

type DeleteEndpointInput = origin_eventbridge.DeleteEndpointInput

type DeleteEndpointOutput = origin_eventbridge.DeleteEndpointOutput

type RemovePermissionInput = origin_eventbridge.RemovePermissionInput

type RemovePermissionOutput = origin_eventbridge.RemovePermissionOutput

type UpdateApiDestinationInput = origin_eventbridge.UpdateApiDestinationInput

type UpdateApiDestinationOutput = origin_eventbridge.UpdateApiDestinationOutput

type CreateEndpointInput = origin_eventbridge.CreateEndpointInput

type CreateEndpointOutput = origin_eventbridge.CreateEndpointOutput

type CreateApiDestinationInput = origin_eventbridge.CreateApiDestinationInput

type CreateApiDestinationOutput = origin_eventbridge.CreateApiDestinationOutput

type CreateConnectionInput = origin_eventbridge.CreateConnectionInput

type CreateConnectionOutput = origin_eventbridge.CreateConnectionOutput

type DescribeEndpointInput = origin_eventbridge.DescribeEndpointInput

type DescribeEndpointOutput = origin_eventbridge.DescribeEndpointOutput

type DescribeEventBusInput = origin_eventbridge.DescribeEventBusInput

type DescribeEventBusOutput = origin_eventbridge.DescribeEventBusOutput

type ListPartnerEventSourcesInput = origin_eventbridge.ListPartnerEventSourcesInput

type ListPartnerEventSourcesOutput = origin_eventbridge.ListPartnerEventSourcesOutput

type ListRulesInput = origin_eventbridge.ListRulesInput

type ListRulesOutput = origin_eventbridge.ListRulesOutput

type UntagResourceInput = origin_eventbridge.UntagResourceInput

type UntagResourceOutput = origin_eventbridge.UntagResourceOutput

type ActivateEventSourceInput = origin_eventbridge.ActivateEventSourceInput

type ActivateEventSourceOutput = origin_eventbridge.ActivateEventSourceOutput

type UpdateArchiveInput = origin_eventbridge.UpdateArchiveInput

type UpdateArchiveOutput = origin_eventbridge.UpdateArchiveOutput

type ListEventBusesInput = origin_eventbridge.ListEventBusesInput

type ListEventBusesOutput = origin_eventbridge.ListEventBusesOutput

type ListArchivesInput = origin_eventbridge.ListArchivesInput

type ListArchivesOutput = origin_eventbridge.ListArchivesOutput

const ServiceID = origin_eventbridge.ServiceID

const ServiceAPIVersion = origin_eventbridge.ServiceAPIVersion

type Client = origin_eventbridge.Client

func attributeSetter(ctx context.Context, ii middleware.InitializeInput) []attribute.KeyValue {
	result := []attribute.KeyValue{}
	return result
}

func New(options Options, optFns ...func(*Options)) *Client {
	attributeSetterOpt := otelaws.WithAttributeSetter(attributeSetter)
	otelaws.AppendMiddlewares(&options.APIOptions, attributeSetterOpt)
	return origin_eventbridge.New(options, optFns...)
}

type Options = origin_eventbridge.Options

func WithAPIOptions(optFns ...func(*middleware.Stack) error) func(*Options) {
	return origin_eventbridge.WithAPIOptions(optFns...)
}

func WithEndpointResolver(v EndpointResolver) func(*Options) {
	return origin_eventbridge.WithEndpointResolver(v)
}

type HTTPClient = origin_eventbridge.HTTPClient

func NewFromConfig(cfg aws.Config, optFns ...func(*Options)) *Client {
	attributeSetterOpt := otelaws.WithAttributeSetter(attributeSetter)
	otelaws.AppendMiddlewares(&cfg.APIOptions, attributeSetterOpt)
	return origin_eventbridge.NewFromConfig(cfg, optFns...)
}

type HTTPSignerV4 = origin_eventbridge.HTTPSignerV4

type DeleteRuleInput = origin_eventbridge.DeleteRuleInput

type DeleteRuleOutput = origin_eventbridge.DeleteRuleOutput

type DescribeApiDestinationInput = origin_eventbridge.DescribeApiDestinationInput

type DescribeApiDestinationOutput = origin_eventbridge.DescribeApiDestinationOutput

type DescribeArchiveInput = origin_eventbridge.DescribeArchiveInput

type DescribeArchiveOutput = origin_eventbridge.DescribeArchiveOutput

type EnableRuleInput = origin_eventbridge.EnableRuleInput

type EnableRuleOutput = origin_eventbridge.EnableRuleOutput

type ListTagsForResourceInput = origin_eventbridge.ListTagsForResourceInput

type ListTagsForResourceOutput = origin_eventbridge.ListTagsForResourceOutput

type CreateEventBusInput = origin_eventbridge.CreateEventBusInput

type CreateEventBusOutput = origin_eventbridge.CreateEventBusOutput

type DescribeRuleInput = origin_eventbridge.DescribeRuleInput

type DescribeRuleOutput = origin_eventbridge.DescribeRuleOutput

type ListReplaysInput = origin_eventbridge.ListReplaysInput

type ListReplaysOutput = origin_eventbridge.ListReplaysOutput

type TagResourceInput = origin_eventbridge.TagResourceInput

type TagResourceOutput = origin_eventbridge.TagResourceOutput

type UpdateConnectionInput = origin_eventbridge.UpdateConnectionInput

type UpdateConnectionOutput = origin_eventbridge.UpdateConnectionOutput

type DeauthorizeConnectionInput = origin_eventbridge.DeauthorizeConnectionInput

type DeauthorizeConnectionOutput = origin_eventbridge.DeauthorizeConnectionOutput
