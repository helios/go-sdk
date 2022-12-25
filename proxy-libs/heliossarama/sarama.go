package heliossarama

import (
	"context"
	"hash"
	"net"

	originalSarama "github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
)

const APIKeySASLAuth = originalSarama.APIKeySASLAuth
const AclOperationAll = originalSarama.AclOperationAll
const AclOperationAlter = originalSarama.AclOperationAlter
const AclOperationAlterConfigs = originalSarama.AclOperationAlterConfigs
const AclOperationAny = originalSarama.AclOperationAny
const AclOperationClusterAction = originalSarama.AclOperationClusterAction
const AclOperationCreate = originalSarama.AclOperationCreate
const AclOperationDelete = originalSarama.AclOperationDelete
const AclOperationDescribe = originalSarama.AclOperationDescribe
const AclOperationDescribeConfigs = originalSarama.AclOperationDescribeConfigs
const AclOperationIdempotentWrite = originalSarama.AclOperationIdempotentWrite
const AclOperationRead = originalSarama.AclOperationRead
const AclOperationUnknown = originalSarama.AclOperationUnknown
const AclOperationWrite = originalSarama.AclOperationWrite
const AclPatternAny = originalSarama.AclPatternAny
const AclPatternLiteral = originalSarama.AclPatternLiteral
const AclPatternMatch = originalSarama.AclPatternMatch
const AclPatternPrefixed = originalSarama.AclPatternPrefixed
const AclPatternUnknown = originalSarama.AclPatternUnknown
const AclPermissionAllow = originalSarama.AclPermissionAllow
const AclPermissionAny = originalSarama.AclPermissionAny
const AclPermissionDeny = originalSarama.AclPermissionDeny
const AclPermissionUnknown = originalSarama.AclPermissionUnknown
const AclResourceAny = originalSarama.AclResourceAny
const AclResourceCluster = originalSarama.AclResourceCluster
const AclResourceDelegationToken = originalSarama.AclResourceDelegationToken
const AclResourceGroup = originalSarama.AclResourceGroup
const AclResourceTopic = originalSarama.AclResourceTopic
const AclResourceTransactionalID = originalSarama.AclResourceTransactionalID
const AclResourceUnknown = originalSarama.AclResourceUnknown
const BrokerLoggerResource = originalSarama.BrokerLoggerResource
const BrokerResource = originalSarama.BrokerResource
const CompressionGZIP = originalSarama.CompressionGZIP
const CompressionLZ4 = originalSarama.CompressionLZ4
const CompressionLevelDefault = originalSarama.CompressionLevelDefault
const CompressionNone = originalSarama.CompressionNone
const CompressionSnappy = originalSarama.CompressionSnappy
const CompressionZSTD = originalSarama.CompressionZSTD
const ControlRecordAbort = originalSarama.ControlRecordAbort
const ControlRecordCommit = originalSarama.ControlRecordCommit
const ControlRecordUnknown = originalSarama.ControlRecordUnknown
const CoordinatorGroup = originalSarama.CoordinatorGroup
const CoordinatorTransaction = originalSarama.CoordinatorTransaction
const ErrBrokerNotAvailable = originalSarama.ErrBrokerNotAvailable
const ErrClusterAuthorizationFailed = originalSarama.ErrClusterAuthorizationFailed
const ErrConcurrentTransactions = originalSarama.ErrConcurrentTransactions
const ErrConsumerCoordinatorNotAvailable = originalSarama.ErrConsumerCoordinatorNotAvailable
const ErrDelegationTokenAuthDisabled = originalSarama.ErrDelegationTokenAuthDisabled
const ErrDelegationTokenAuthorizationFailed = originalSarama.ErrDelegationTokenAuthorizationFailed
const ErrDelegationTokenExpired = originalSarama.ErrDelegationTokenExpired
const ErrDelegationTokenNotFound = originalSarama.ErrDelegationTokenNotFound
const ErrDelegationTokenOwnerMismatch = originalSarama.ErrDelegationTokenOwnerMismatch
const ErrDelegationTokenRequestNotAllowed = originalSarama.ErrDelegationTokenRequestNotAllowed
const ErrDuplicateSequenceNumber = originalSarama.ErrDuplicateSequenceNumber
const ErrElectionNotNeeded = originalSarama.ErrElectionNotNeeded
const ErrEligibleLeadersNotAvailable = originalSarama.ErrEligibleLeadersNotAvailable
const ErrFencedInstancedId = originalSarama.ErrFencedInstancedId
const ErrFencedLeaderEpoch = originalSarama.ErrFencedLeaderEpoch
const ErrFetchSessionIDNotFound = originalSarama.ErrFetchSessionIDNotFound
const ErrGroupAuthorizationFailed = originalSarama.ErrGroupAuthorizationFailed
const ErrGroupIDNotFound = originalSarama.ErrGroupIDNotFound
const ErrGroupMaxSizeReached = originalSarama.ErrGroupMaxSizeReached
const ErrGroupSubscribedToTopic = originalSarama.ErrGroupSubscribedToTopic
const ErrIllegalGeneration = originalSarama.ErrIllegalGeneration
const ErrIllegalSASLState = originalSarama.ErrIllegalSASLState
const ErrInconsistentGroupProtocol = originalSarama.ErrInconsistentGroupProtocol
const ErrInvalidCommitOffsetSize = originalSarama.ErrInvalidCommitOffsetSize
const ErrInvalidConfig = originalSarama.ErrInvalidConfig
const ErrInvalidFetchSessionEpoch = originalSarama.ErrInvalidFetchSessionEpoch
const ErrInvalidGroupId = originalSarama.ErrInvalidGroupId
const ErrInvalidMessage = originalSarama.ErrInvalidMessage
const ErrInvalidMessageSize = originalSarama.ErrInvalidMessageSize
const ErrInvalidPartitions = originalSarama.ErrInvalidPartitions
const ErrInvalidPrincipalType = originalSarama.ErrInvalidPrincipalType
const ErrInvalidProducerEpoch = originalSarama.ErrInvalidProducerEpoch
const ErrInvalidProducerIDMapping = originalSarama.ErrInvalidProducerIDMapping
const ErrInvalidRecord = originalSarama.ErrInvalidRecord
const ErrInvalidReplicaAssignment = originalSarama.ErrInvalidReplicaAssignment
const ErrInvalidReplicationFactor = originalSarama.ErrInvalidReplicationFactor
const ErrInvalidRequest = originalSarama.ErrInvalidRequest
const ErrInvalidRequiredAcks = originalSarama.ErrInvalidRequiredAcks
const ErrInvalidSessionTimeout = originalSarama.ErrInvalidSessionTimeout
const ErrInvalidTimestamp = originalSarama.ErrInvalidTimestamp
const ErrInvalidTopic = originalSarama.ErrInvalidTopic
const ErrInvalidTransactionTimeout = originalSarama.ErrInvalidTransactionTimeout
const ErrInvalidTxnState = originalSarama.ErrInvalidTxnState
const ErrKafkaStorageError = originalSarama.ErrKafkaStorageError
const ErrLeaderNotAvailable = originalSarama.ErrLeaderNotAvailable
const ErrListenerNotFound = originalSarama.ErrListenerNotFound
const ErrLogDirNotFound = originalSarama.ErrLogDirNotFound
const ErrMemberIdRequired = originalSarama.ErrMemberIdRequired
const ErrMessageSetSizeTooLarge = originalSarama.ErrMessageSetSizeTooLarge
const ErrMessageSizeTooLarge = originalSarama.ErrMessageSizeTooLarge
const ErrNetworkException = originalSarama.ErrNetworkException
const ErrNoError = originalSarama.ErrNoError
const ErrNoReassignmentInProgress = originalSarama.ErrNoReassignmentInProgress
const ErrNonEmptyGroup = originalSarama.ErrNonEmptyGroup
const ErrNotController = originalSarama.ErrNotController
const ErrNotCoordinatorForConsumer = originalSarama.ErrNotCoordinatorForConsumer
const ErrNotEnoughReplicas = originalSarama.ErrNotEnoughReplicas
const ErrNotEnoughReplicasAfterAppend = originalSarama.ErrNotEnoughReplicasAfterAppend
const ErrNotLeaderForPartition = originalSarama.ErrNotLeaderForPartition
const ErrOffsetMetadataTooLarge = originalSarama.ErrOffsetMetadataTooLarge
const ErrOffsetNotAvailable = originalSarama.ErrOffsetNotAvailable
const ErrOffsetOutOfRange = originalSarama.ErrOffsetOutOfRange
const ErrOffsetsLoadInProgress = originalSarama.ErrOffsetsLoadInProgress
const ErrOperationNotAttempted = originalSarama.ErrOperationNotAttempted
const ErrOutOfOrderSequenceNumber = originalSarama.ErrOutOfOrderSequenceNumber
const ErrPolicyViolation = originalSarama.ErrPolicyViolation
const ErrPreferredLeaderNotAvailable = originalSarama.ErrPreferredLeaderNotAvailable
const ErrProducerFenced = originalSarama.ErrProducerFenced
const ErrReassignmentInProgress = originalSarama.ErrReassignmentInProgress
const ErrRebalanceInProgress = originalSarama.ErrRebalanceInProgress
const ErrReplicaNotAvailable = originalSarama.ErrReplicaNotAvailable
const ErrRequestTimedOut = originalSarama.ErrRequestTimedOut
const ErrSASLAuthenticationFailed = originalSarama.ErrSASLAuthenticationFailed
const ErrSecurityDisabled = originalSarama.ErrSecurityDisabled
const ErrStaleBrokerEpoch = originalSarama.ErrStaleBrokerEpoch
const ErrStaleControllerEpochCode = originalSarama.ErrStaleControllerEpochCode
const ErrThrottlingQuotaExceeded = originalSarama.ErrThrottlingQuotaExceeded
const ErrTopicAlreadyExists = originalSarama.ErrTopicAlreadyExists
const ErrTopicAuthorizationFailed = originalSarama.ErrTopicAuthorizationFailed
const ErrTopicDeletionDisabled = originalSarama.ErrTopicDeletionDisabled
const ErrTransactionCoordinatorFenced = originalSarama.ErrTransactionCoordinatorFenced
const ErrTransactionalIDAuthorizationFailed = originalSarama.ErrTransactionalIDAuthorizationFailed
const ErrUnknown = originalSarama.ErrUnknown
const ErrUnknownLeaderEpoch = originalSarama.ErrUnknownLeaderEpoch
const ErrUnknownMemberId = originalSarama.ErrUnknownMemberId
const ErrUnknownProducerID = originalSarama.ErrUnknownProducerID
const ErrUnknownTopicOrPartition = originalSarama.ErrUnknownTopicOrPartition
const ErrUnstableOffsetCommit = originalSarama.ErrUnstableOffsetCommit
const ErrUnsupportedCompressionType = originalSarama.ErrUnsupportedCompressionType
const ErrUnsupportedForMessageFormat = originalSarama.ErrUnsupportedForMessageFormat
const ErrUnsupportedSASLMechanism = originalSarama.ErrUnsupportedSASLMechanism
const ErrUnsupportedVersion = originalSarama.ErrUnsupportedVersion
const GSS_API_FINISH = originalSarama.GSS_API_FINISH
const GSS_API_GENERIC_TAG = originalSarama.GSS_API_GENERIC_TAG
const GSS_API_INITIAL = originalSarama.GSS_API_INITIAL
const GSS_API_VERIFY = originalSarama.GSS_API_VERIFY
const GroupGenerationUndefined = originalSarama.GroupGenerationUndefined
const IncrementalAlterConfigsOperationAppend = originalSarama.IncrementalAlterConfigsOperationAppend
const IncrementalAlterConfigsOperationDelete = originalSarama.IncrementalAlterConfigsOperationDelete
const IncrementalAlterConfigsOperationSet = originalSarama.IncrementalAlterConfigsOperationSet
const IncrementalAlterConfigsOperationSubtract = originalSarama.IncrementalAlterConfigsOperationSubtract
const KRB5_KEYTAB_AUTH = originalSarama.KRB5_KEYTAB_AUTH
const KRB5_USER_AUTH = originalSarama.KRB5_USER_AUTH
const MAX_GROUP_INSTANCE_ID_LENGTH = originalSarama.MAX_GROUP_INSTANCE_ID_LENGTH
const NoResponse = originalSarama.NoResponse
const OffsetNewest = originalSarama.OffsetNewest
const OffsetOldest = originalSarama.OffsetOldest
const ProducerTxnFlagAbortableError = originalSarama.ProducerTxnFlagAbortableError
const ProducerTxnFlagAbortingTransaction = originalSarama.ProducerTxnFlagAbortingTransaction
const ProducerTxnFlagCommittingTransaction = originalSarama.ProducerTxnFlagCommittingTransaction
const ProducerTxnFlagEndTransaction = originalSarama.ProducerTxnFlagEndTransaction
const ProducerTxnFlagFatalError = originalSarama.ProducerTxnFlagFatalError
const ProducerTxnFlagInError = originalSarama.ProducerTxnFlagInError
const ProducerTxnFlagInTransaction = originalSarama.ProducerTxnFlagInTransaction
const ProducerTxnFlagInitializing = originalSarama.ProducerTxnFlagInitializing
const ProducerTxnFlagReady = originalSarama.ProducerTxnFlagReady
const ProducerTxnFlagUninitialized = originalSarama.ProducerTxnFlagUninitialized
const QuotaEntityClientID = originalSarama.QuotaEntityClientID
const QuotaEntityIP = originalSarama.QuotaEntityIP
const QuotaEntityUser = originalSarama.QuotaEntityUser
const QuotaMatchAny = originalSarama.QuotaMatchAny
const QuotaMatchDefault = originalSarama.QuotaMatchDefault
const QuotaMatchExact = originalSarama.QuotaMatchExact
const RangeBalanceStrategyName = originalSarama.RangeBalanceStrategyName
const ReadCommitted = originalSarama.ReadCommitted
const ReadUncommitted = originalSarama.ReadUncommitted
const ReceiveTime = originalSarama.ReceiveTime
const RoundRobinBalanceStrategyName = originalSarama.RoundRobinBalanceStrategyName
const SASLExtKeyAuth = originalSarama.SASLExtKeyAuth
const SASLHandshakeV0 = originalSarama.SASLHandshakeV0
const SASLHandshakeV1 = originalSarama.SASLHandshakeV1
const SASLTypeGSSAPI = originalSarama.SASLTypeGSSAPI
const SASLTypeOAuth = originalSarama.SASLTypeOAuth
const SASLTypePlaintext = originalSarama.SASLTypePlaintext
const SASLTypeSCRAMSHA256 = originalSarama.SASLTypeSCRAMSHA256
const SASLTypeSCRAMSHA512 = originalSarama.SASLTypeSCRAMSHA512
const SCRAM_MECHANISM_SHA_256 = originalSarama.SCRAM_MECHANISM_SHA_256
const SCRAM_MECHANISM_SHA_512 = originalSarama.SCRAM_MECHANISM_SHA_512
const SCRAM_MECHANISM_UNKNOWN = originalSarama.SCRAM_MECHANISM_UNKNOWN
const SourceDefault = originalSarama.SourceDefault
const SourceDynamicBroker = originalSarama.SourceDynamicBroker
const SourceDynamicDefaultBroker = originalSarama.SourceDynamicDefaultBroker
const SourceStaticBroker = originalSarama.SourceStaticBroker
const SourceTopic = originalSarama.SourceTopic
const SourceUnknown = originalSarama.SourceUnknown
const StickyBalanceStrategyName = originalSarama.StickyBalanceStrategyName
const TOK_ID_KRB_AP_REQ = originalSarama.TOK_ID_KRB_AP_REQ
const TopicResource = originalSarama.TopicResource
const UnknownResource = originalSarama.UnknownResource
const WaitForAll = originalSarama.WaitForAll
const WaitForLocal = originalSarama.WaitForLocal

type AbortedTransaction = originalSarama.AbortedTransaction
type AccessToken = originalSarama.AccessToken
type AccessTokenProvider = originalSarama.AccessTokenProvider
type Acl = originalSarama.Acl
type AclCreation = originalSarama.AclCreation
type AclCreationResponse = originalSarama.AclCreationResponse
type AclFilter = originalSarama.AclFilter
type AclOperation = originalSarama.AclOperation
type AclPermissionType = originalSarama.AclPermissionType
type AclResourcePatternType = originalSarama.AclResourcePatternType
type AclResourceType = originalSarama.AclResourceType
type AddOffsetsToTxnRequest = originalSarama.AddOffsetsToTxnRequest
type AddOffsetsToTxnResponse = originalSarama.AddOffsetsToTxnResponse
type AddPartitionsToTxnRequest = originalSarama.AddPartitionsToTxnRequest
type AddPartitionsToTxnResponse = originalSarama.AddPartitionsToTxnResponse
type AlterClientQuotasEntry = originalSarama.AlterClientQuotasEntry
type AlterClientQuotasEntryResponse = originalSarama.AlterClientQuotasEntryResponse
type AlterClientQuotasRequest = originalSarama.AlterClientQuotasRequest
type AlterClientQuotasResponse = originalSarama.AlterClientQuotasResponse
type AlterConfigsRequest = originalSarama.AlterConfigsRequest
type AlterConfigsResource = originalSarama.AlterConfigsResource
type AlterConfigsResourceResponse = originalSarama.AlterConfigsResourceResponse
type AlterConfigsResponse = originalSarama.AlterConfigsResponse
type AlterPartitionReassignmentsRequest = originalSarama.AlterPartitionReassignmentsRequest
type AlterPartitionReassignmentsResponse = originalSarama.AlterPartitionReassignmentsResponse
type AlterUserScramCredentialsDelete = originalSarama.AlterUserScramCredentialsDelete
type AlterUserScramCredentialsRequest = originalSarama.AlterUserScramCredentialsRequest
type AlterUserScramCredentialsResponse = originalSarama.AlterUserScramCredentialsResponse
type AlterUserScramCredentialsResult = originalSarama.AlterUserScramCredentialsResult
type AlterUserScramCredentialsUpsert = originalSarama.AlterUserScramCredentialsUpsert
type ApiVersionsRequest = originalSarama.ApiVersionsRequest
type ApiVersionsResponse = originalSarama.ApiVersionsResponse
type ApiVersionsResponseKey = originalSarama.ApiVersionsResponseKey
type AsyncProducer = originalSarama.AsyncProducer
type BalanceStrategy = originalSarama.BalanceStrategy
type BalanceStrategyPlan = originalSarama.BalanceStrategyPlan
type Broker = originalSarama.Broker
type ByteEncoder = originalSarama.ByteEncoder
type Client = originalSarama.Client
type ClientQuotasOp = originalSarama.ClientQuotasOp
type ClusterAdmin = originalSarama.ClusterAdmin
type CompressionCodec = originalSarama.CompressionCodec
type Config = originalSarama.Config
type ConfigEntry = originalSarama.ConfigEntry
type ConfigResource = originalSarama.ConfigResource
type ConfigResourceType = originalSarama.ConfigResourceType
type ConfigSource = originalSarama.ConfigSource
type ConfigSynonym = originalSarama.ConfigSynonym
type ConfigurationError = originalSarama.ConfigurationError
type Consumer = originalSarama.Consumer
type ConsumerError = originalSarama.ConsumerError
type ConsumerErrors = originalSarama.ConsumerErrors
type ConsumerGroup = originalSarama.ConsumerGroup
type ConsumerGroupClaim = originalSarama.ConsumerGroupClaim
type ConsumerGroupHandler = originalSarama.ConsumerGroupHandler
type ConsumerGroupMemberAssignment = originalSarama.ConsumerGroupMemberAssignment
type ConsumerGroupMemberMetadata = originalSarama.ConsumerGroupMemberMetadata
type ConsumerGroupSession = originalSarama.ConsumerGroupSession
type ConsumerInterceptor = originalSarama.ConsumerInterceptor
type ConsumerMessage = originalSarama.ConsumerMessage
type ConsumerMetadataRequest = originalSarama.ConsumerMetadataRequest
type ConsumerMetadataResponse = originalSarama.ConsumerMetadataResponse
type ControlRecord = originalSarama.ControlRecord
type ControlRecordType = originalSarama.ControlRecordType
type CoordinatorType = originalSarama.CoordinatorType
type CreateAclsRequest = originalSarama.CreateAclsRequest
type CreateAclsResponse = originalSarama.CreateAclsResponse
type CreatePartitionsRequest = originalSarama.CreatePartitionsRequest
type CreatePartitionsResponse = originalSarama.CreatePartitionsResponse
type CreateTopicsRequest = originalSarama.CreateTopicsRequest
type CreateTopicsResponse = originalSarama.CreateTopicsResponse
type DeleteAclsRequest = originalSarama.DeleteAclsRequest
type DeleteAclsResponse = originalSarama.DeleteAclsResponse
type DeleteGroupsRequest = originalSarama.DeleteGroupsRequest
type DeleteGroupsResponse = originalSarama.DeleteGroupsResponse
type DeleteOffsetsRequest = originalSarama.DeleteOffsetsRequest
type DeleteOffsetsResponse = originalSarama.DeleteOffsetsResponse
type DeleteRecordsRequest = originalSarama.DeleteRecordsRequest
type DeleteRecordsRequestTopic = originalSarama.DeleteRecordsRequestTopic
type DeleteRecordsResponse = originalSarama.DeleteRecordsResponse
type DeleteRecordsResponsePartition = originalSarama.DeleteRecordsResponsePartition
type DeleteRecordsResponseTopic = originalSarama.DeleteRecordsResponseTopic
type DeleteTopicsRequest = originalSarama.DeleteTopicsRequest
type DeleteTopicsResponse = originalSarama.DeleteTopicsResponse
type DescribeAclsRequest = originalSarama.DescribeAclsRequest
type DescribeAclsResponse = originalSarama.DescribeAclsResponse
type DescribeClientQuotasEntry = originalSarama.DescribeClientQuotasEntry
type DescribeClientQuotasRequest = originalSarama.DescribeClientQuotasRequest
type DescribeClientQuotasResponse = originalSarama.DescribeClientQuotasResponse
type DescribeConfigsRequest = originalSarama.DescribeConfigsRequest
type DescribeConfigsResponse = originalSarama.DescribeConfigsResponse
type DescribeGroupsRequest = originalSarama.DescribeGroupsRequest
type DescribeGroupsResponse = originalSarama.DescribeGroupsResponse
type DescribeLogDirsRequest = originalSarama.DescribeLogDirsRequest
type DescribeLogDirsRequestTopic = originalSarama.DescribeLogDirsRequestTopic
type DescribeLogDirsResponse = originalSarama.DescribeLogDirsResponse
type DescribeLogDirsResponseDirMetadata = originalSarama.DescribeLogDirsResponseDirMetadata
type DescribeLogDirsResponsePartition = originalSarama.DescribeLogDirsResponsePartition
type DescribeLogDirsResponseTopic = originalSarama.DescribeLogDirsResponseTopic
type DescribeUserScramCredentialsRequest = originalSarama.DescribeUserScramCredentialsRequest
type DescribeUserScramCredentialsRequestUser = originalSarama.DescribeUserScramCredentialsRequestUser
type DescribeUserScramCredentialsResponse = originalSarama.DescribeUserScramCredentialsResponse
type DescribeUserScramCredentialsResult = originalSarama.DescribeUserScramCredentialsResult
type DynamicConsistencyPartitioner = originalSarama.DynamicConsistencyPartitioner
type Encoder = originalSarama.Encoder
type EndTxnRequest = originalSarama.EndTxnRequest
type EndTxnResponse = originalSarama.EndTxnResponse
type FetchRequest = originalSarama.FetchRequest
type FetchResponse = originalSarama.FetchResponse
type FetchResponseBlock = originalSarama.FetchResponseBlock
type FilterResponse = originalSarama.FilterResponse
type FindCoordinatorRequest = originalSarama.FindCoordinatorRequest
type FindCoordinatorResponse = originalSarama.FindCoordinatorResponse
type GSSAPIConfig = originalSarama.GSSAPIConfig
type GSSAPIKerberosAuth = originalSarama.GSSAPIKerberosAuth
type GSSApiHandlerFunc = originalSarama.GSSApiHandlerFunc
type GroupDescription = originalSarama.GroupDescription
type GroupMember = originalSarama.GroupMember
type GroupMemberDescription = originalSarama.GroupMemberDescription
type GroupProtocol = originalSarama.GroupProtocol
type HashPartitionerOption = originalSarama.HashPartitionerOption
type HeartbeatRequest = originalSarama.HeartbeatRequest
type HeartbeatResponse = originalSarama.HeartbeatResponse
type IncrementalAlterConfigsEntry = originalSarama.IncrementalAlterConfigsEntry
type IncrementalAlterConfigsOperation = originalSarama.IncrementalAlterConfigsOperation
type IncrementalAlterConfigsRequest = originalSarama.IncrementalAlterConfigsRequest
type IncrementalAlterConfigsResource = originalSarama.IncrementalAlterConfigsResource
type IncrementalAlterConfigsResponse = originalSarama.IncrementalAlterConfigsResponse
type InitProducerIDRequest = originalSarama.InitProducerIDRequest
type InitProducerIDResponse = originalSarama.InitProducerIDResponse
type IsolationLevel = originalSarama.IsolationLevel
type JoinGroupRequest = originalSarama.JoinGroupRequest
type JoinGroupResponse = originalSarama.JoinGroupResponse
type KError = originalSarama.KError
type KafkaGSSAPIHandler = originalSarama.KafkaGSSAPIHandler
type KafkaVersion = originalSarama.KafkaVersion
type KerberosClient = originalSarama.KerberosClient
type KerberosGoKrb5Client = originalSarama.KerberosGoKrb5Client
type LeaveGroupRequest = originalSarama.LeaveGroupRequest
type LeaveGroupResponse = originalSarama.LeaveGroupResponse
type ListGroupsRequest = originalSarama.ListGroupsRequest
type ListGroupsResponse = originalSarama.ListGroupsResponse
type ListPartitionReassignmentsRequest = originalSarama.ListPartitionReassignmentsRequest
type ListPartitionReassignmentsResponse = originalSarama.ListPartitionReassignmentsResponse
type MatchingAcl = originalSarama.MatchingAcl
type MemberIdentity = originalSarama.MemberIdentity
type MemberResponse = originalSarama.MemberResponse
type Message = originalSarama.Message
type MessageBlock = originalSarama.MessageBlock
type MessageSet = originalSarama.MessageSet
type MetadataRequest = originalSarama.MetadataRequest
type MetadataResponse = originalSarama.MetadataResponse
type MockAlterConfigsResponse = originalSarama.MockAlterConfigsResponse
type MockAlterConfigsResponseWithErrorCode = originalSarama.MockAlterConfigsResponseWithErrorCode
type MockAlterPartitionReassignmentsResponse = originalSarama.MockAlterPartitionReassignmentsResponse
type MockApiVersionsResponse = originalSarama.MockApiVersionsResponse
type MockBroker = originalSarama.MockBroker
type MockConsumerMetadataResponse = originalSarama.MockConsumerMetadataResponse
type MockCreateAclsResponse = originalSarama.MockCreateAclsResponse
type MockCreateAclsResponseError = originalSarama.MockCreateAclsResponseError
type MockCreatePartitionsResponse = originalSarama.MockCreatePartitionsResponse
type MockCreateTopicsResponse = originalSarama.MockCreateTopicsResponse
type MockDeleteAclsResponse = originalSarama.MockDeleteAclsResponse
type MockDeleteGroupsResponse = originalSarama.MockDeleteGroupsResponse
type MockDeleteOffsetResponse = originalSarama.MockDeleteOffsetResponse
type MockDeleteRecordsResponse = originalSarama.MockDeleteRecordsResponse
type MockDeleteTopicsResponse = originalSarama.MockDeleteTopicsResponse
type MockDescribeConfigsResponse = originalSarama.MockDescribeConfigsResponse
type MockDescribeConfigsResponseWithErrorCode = originalSarama.MockDescribeConfigsResponseWithErrorCode
type MockDescribeGroupsResponse = originalSarama.MockDescribeGroupsResponse
type MockDescribeLogDirsResponse = originalSarama.MockDescribeLogDirsResponse
type MockFetchResponse = originalSarama.MockFetchResponse
type MockFindCoordinatorResponse = originalSarama.MockFindCoordinatorResponse
type MockHeartbeatResponse = originalSarama.MockHeartbeatResponse
type MockIncrementalAlterConfigsResponse = originalSarama.MockIncrementalAlterConfigsResponse
type MockIncrementalAlterConfigsResponseWithErrorCode = originalSarama.MockIncrementalAlterConfigsResponseWithErrorCode
type MockJoinGroupResponse = originalSarama.MockJoinGroupResponse
type MockKerberosClient = originalSarama.MockKerberosClient
type MockLeaveGroupResponse = originalSarama.MockLeaveGroupResponse
type MockListAclsResponse = originalSarama.MockListAclsResponse
type MockListGroupsResponse = originalSarama.MockListGroupsResponse
type MockListPartitionReassignmentsResponse = originalSarama.MockListPartitionReassignmentsResponse
type MockMetadataResponse = originalSarama.MockMetadataResponse
type MockOffsetCommitResponse = originalSarama.MockOffsetCommitResponse
type MockOffsetFetchResponse = originalSarama.MockOffsetFetchResponse
type MockOffsetResponse = originalSarama.MockOffsetResponse
type MockProduceResponse = originalSarama.MockProduceResponse
type MockResponse = originalSarama.MockResponse
type MockSaslAuthenticateResponse = originalSarama.MockSaslAuthenticateResponse
type MockSaslHandshakeResponse = originalSarama.MockSaslHandshakeResponse
type MockSequence = originalSarama.MockSequence
type MockSyncGroupResponse = originalSarama.MockSyncGroupResponse
type MockWrapper = originalSarama.MockWrapper
type OffsetCommitRequest = originalSarama.OffsetCommitRequest
type OffsetCommitResponse = originalSarama.OffsetCommitResponse
type OffsetFetchRequest = originalSarama.OffsetFetchRequest
type OffsetFetchResponse = originalSarama.OffsetFetchResponse
type OffsetFetchResponseBlock = originalSarama.OffsetFetchResponseBlock
type OffsetManager = originalSarama.OffsetManager
type OffsetRequest = originalSarama.OffsetRequest
type OffsetResponse = originalSarama.OffsetResponse
type OffsetResponseBlock = originalSarama.OffsetResponseBlock
type OwnedPartition = originalSarama.OwnedPartition
type PacketDecodingError = originalSarama.PacketDecodingError
type PacketEncodingError = originalSarama.PacketEncodingError
type PartitionConsumer = originalSarama.PartitionConsumer
type PartitionError = originalSarama.PartitionError
type PartitionMetadata = originalSarama.PartitionMetadata
type PartitionOffsetManager = originalSarama.PartitionOffsetManager
type PartitionOffsetMetadata = originalSarama.PartitionOffsetMetadata
type PartitionReplicaReassignmentsStatus = originalSarama.PartitionReplicaReassignmentsStatus
type Partitioner = originalSarama.Partitioner
type PartitionerConstructor = originalSarama.PartitionerConstructor
type ProduceCallback = originalSarama.ProduceCallback
type ProduceRequest = originalSarama.ProduceRequest
type ProduceResponse = originalSarama.ProduceResponse
type ProduceResponseBlock = originalSarama.ProduceResponseBlock
type ProducerError = originalSarama.ProducerError
type ProducerErrors = originalSarama.ProducerErrors
type ProducerInterceptor = originalSarama.ProducerInterceptor
type ProducerMessage = originalSarama.ProducerMessage
type ProducerTxnStatusFlag = originalSarama.ProducerTxnStatusFlag
type QuotaEntityComponent = originalSarama.QuotaEntityComponent
type QuotaEntityType = originalSarama.QuotaEntityType
type QuotaFilterComponent = originalSarama.QuotaFilterComponent
type QuotaMatchType = originalSarama.QuotaMatchType
type Record = originalSarama.Record
type RecordBatch = originalSarama.RecordBatch
type RecordHeader = originalSarama.RecordHeader
type Records = originalSarama.Records
type RequestNotifierFunc = originalSarama.RequestNotifierFunc
type RequestResponse = originalSarama.RequestResponse
type RequiredAcks = originalSarama.RequiredAcks
type Resource = originalSarama.Resource
type ResourceAcls = originalSarama.ResourceAcls
type ResourceResponse = originalSarama.ResourceResponse
type SASLMechanism = originalSarama.SASLMechanism
type SCRAMClient = originalSarama.SCRAMClient
type SaslAuthenticateRequest = originalSarama.SaslAuthenticateRequest
type SaslAuthenticateResponse = originalSarama.SaslAuthenticateResponse
type SaslHandshakeRequest = originalSarama.SaslHandshakeRequest
type SaslHandshakeResponse = originalSarama.SaslHandshakeResponse
type ScramMechanismType = originalSarama.ScramMechanismType
type StdLogger = originalSarama.StdLogger
type StickyAssignorUserData = originalSarama.StickyAssignorUserData
type StickyAssignorUserDataV0 = originalSarama.StickyAssignorUserDataV0
type StickyAssignorUserDataV1 = originalSarama.StickyAssignorUserDataV1
type StringEncoder = originalSarama.StringEncoder
type SyncGroupRequest = originalSarama.SyncGroupRequest
type SyncGroupRequestAssignment = originalSarama.SyncGroupRequestAssignment
type SyncGroupResponse = originalSarama.SyncGroupResponse
type SyncProducer = originalSarama.SyncProducer
type TestReporter = originalSarama.TestReporter
type Timestamp = originalSarama.Timestamp
type TopicDetail = originalSarama.TopicDetail
type TopicError = originalSarama.TopicError
type TopicMetadata = originalSarama.TopicMetadata
type TopicPartition = originalSarama.TopicPartition
type TopicPartitionError = originalSarama.TopicPartitionError
type TxnOffsetCommitRequest = originalSarama.TxnOffsetCommitRequest
type TxnOffsetCommitResponse = originalSarama.TxnOffsetCommitResponse
type UserScramCredentialsResponseInfo = originalSarama.UserScramCredentialsResponseInfo
type ZstdDecoderParams = originalSarama.ZstdDecoderParams
type ZstdEncoderParams = originalSarama.ZstdEncoderParams

var BalanceStrategyRange = originalSarama.BalanceStrategyRange
var BalanceStrategyRoundRobin = originalSarama.BalanceStrategyRoundRobin
var BalanceStrategySticky = originalSarama.BalanceStrategySticky
var DebugLogger = originalSarama.DebugLogger
var DefaultVersion = originalSarama.DefaultVersion
var ErrAddPartitionsToTxn = originalSarama.ErrAddPartitionsToTxn
var ErrAlreadyConnected = originalSarama.ErrAlreadyConnected
var ErrBrokerNotFound = originalSarama.ErrBrokerNotFound
var ErrCannotTransitionNilError = originalSarama.ErrCannotTransitionNilError
var ErrClosedClient = originalSarama.ErrClosedClient
var ErrClosedConsumerGroup = originalSarama.ErrClosedConsumerGroup
var ErrConsumerOffsetNotAdvanced = originalSarama.ErrConsumerOffsetNotAdvanced
var ErrControllerNotAvailable = originalSarama.ErrControllerNotAvailable
var ErrCreateACLs = originalSarama.ErrCreateACLs
var ErrDeleteRecords = originalSarama.ErrDeleteRecords
var ErrIncompleteResponse = originalSarama.ErrIncompleteResponse
var ErrInsufficientData = originalSarama.ErrInsufficientData
var ErrInvalidPartition = originalSarama.ErrInvalidPartition
var ErrMessageTooLarge = originalSarama.ErrMessageTooLarge
var ErrNoTopicsToUpdateMetadata = originalSarama.ErrNoTopicsToUpdateMetadata
var ErrNonTransactedProducer = originalSarama.ErrNonTransactedProducer
var ErrNotConnected = originalSarama.ErrNotConnected
var ErrOutOfBrokers = originalSarama.ErrOutOfBrokers
var ErrReassignPartitions = originalSarama.ErrReassignPartitions
var ErrShuttingDown = originalSarama.ErrShuttingDown
var ErrTransactionNotReady = originalSarama.ErrTransactionNotReady
var ErrTransitionNotAllowed = originalSarama.ErrTransitionNotAllowed
var ErrTxnOffsetCommit = originalSarama.ErrTxnOffsetCommit
var ErrTxnUnableToParseResponse = originalSarama.ErrTxnUnableToParseResponse
var ErrUnknownScramMechanism = originalSarama.ErrUnknownScramMechanism
var GROUP_INSTANCE_ID_REGEXP = originalSarama.GROUP_INSTANCE_ID_REGEXP
var Logger = originalSarama.Logger
var MaxRequestSize = originalSarama.MaxRequestSize
var MaxResponseSize = originalSarama.MaxResponseSize
var MaxVersion = originalSarama.MaxVersion
var MinVersion = originalSarama.MinVersion
var MultiErrorFormat = originalSarama.MultiErrorFormat
var NoNode = originalSarama.NoNode
var PanicHandler = originalSarama.PanicHandler
var SupportedVersions = originalSarama.SupportedVersions
var V0_10_0_0 = originalSarama.V0_10_0_0
var V0_10_0_1 = originalSarama.V0_10_0_1
var V0_10_1_0 = originalSarama.V0_10_1_0
var V0_10_1_1 = originalSarama.V0_10_1_1
var V0_10_2_0 = originalSarama.V0_10_2_0
var V0_10_2_1 = originalSarama.V0_10_2_1
var V0_10_2_2 = originalSarama.V0_10_2_2
var V0_11_0_0 = originalSarama.V0_11_0_0
var V0_11_0_1 = originalSarama.V0_11_0_1
var V0_11_0_2 = originalSarama.V0_11_0_2
var V0_8_2_0 = originalSarama.V0_8_2_0
var V0_8_2_1 = originalSarama.V0_8_2_1
var V0_8_2_2 = originalSarama.V0_8_2_2
var V0_9_0_0 = originalSarama.V0_9_0_0
var V0_9_0_1 = originalSarama.V0_9_0_1
var V1_0_0_0 = originalSarama.V1_0_0_0
var V1_0_1_0 = originalSarama.V1_0_1_0
var V1_0_2_0 = originalSarama.V1_0_2_0
var V1_1_0_0 = originalSarama.V1_1_0_0
var V1_1_1_0 = originalSarama.V1_1_1_0
var V2_0_0_0 = originalSarama.V2_0_0_0
var V2_0_1_0 = originalSarama.V2_0_1_0
var V2_1_0_0 = originalSarama.V2_1_0_0
var V2_1_1_0 = originalSarama.V2_1_1_0
var V2_2_0_0 = originalSarama.V2_2_0_0
var V2_2_1_0 = originalSarama.V2_2_1_0
var V2_2_2_0 = originalSarama.V2_2_2_0
var V2_3_0_0 = originalSarama.V2_3_0_0
var V2_3_1_0 = originalSarama.V2_3_1_0
var V2_4_0_0 = originalSarama.V2_4_0_0
var V2_4_1_0 = originalSarama.V2_4_1_0
var V2_5_0_0 = originalSarama.V2_5_0_0
var V2_5_1_0 = originalSarama.V2_5_1_0
var V2_6_0_0 = originalSarama.V2_6_0_0
var V2_6_1_0 = originalSarama.V2_6_1_0
var V2_6_2_0 = originalSarama.V2_6_2_0
var V2_6_3_0 = originalSarama.V2_6_3_0
var V2_7_0_0 = originalSarama.V2_7_0_0
var V2_7_1_0 = originalSarama.V2_7_1_0
var V2_7_2_0 = originalSarama.V2_7_2_0
var V2_8_0_0 = originalSarama.V2_8_0_0
var V2_8_1_0 = originalSarama.V2_8_1_0
var V2_8_2_0 = originalSarama.V2_8_2_0
var V3_0_0_0 = originalSarama.V3_0_0_0
var V3_0_1_0 = originalSarama.V3_0_1_0
var V3_0_2_0 = originalSarama.V3_0_2_0
var V3_1_0_0 = originalSarama.V3_1_0_0
var V3_1_1_0 = originalSarama.V3_1_1_0
var V3_1_2_0 = originalSarama.V3_1_2_0
var V3_2_0_0 = originalSarama.V3_2_0_0
var V3_2_1_0 = originalSarama.V3_2_1_0
var V3_2_2_0 = originalSarama.V3_2_2_0
var V3_2_3_0 = originalSarama.V3_2_3_0

func NewAsyncProducer(addrs []string, conf *Config) (AsyncProducer, error) {
	asyncProducer, error := originalSarama.NewAsyncProducer(addrs, conf)

	if asyncProducer != nil {
		asyncProducer = asyncProducerWrapper{asyncProducer}
		asyncProducer = otelsarama.WrapAsyncProducer(conf, asyncProducer)
	}

	return asyncProducer, error
}

func NewAsyncProducerFromClient(client Client) (AsyncProducer, error) {
	asyncProducer, error := originalSarama.NewAsyncProducerFromClient(client)

	if asyncProducer != nil {
		asyncProducer = asyncProducerWrapper{asyncProducer}
		asyncProducer = otelsarama.WrapAsyncProducer(client.Config(), asyncProducer)
	}

	return asyncProducer, error
}

func NewBroker(addr string) *Broker {
	return originalSarama.NewBroker(addr)
}

func NewClient(addrs []string, conf *Config) (Client, error) {
	return originalSarama.NewClient(addrs, conf)
}

func NewClusterAdmin(addrs []string, conf *Config) (ClusterAdmin, error) {
	return originalSarama.NewClusterAdmin(addrs, conf)
}

func NewClusterAdminFromClient(client Client) (ClusterAdmin, error) {
	return originalSarama.NewClusterAdminFromClient(client)
}

func NewConfig() *Config {
	return originalSarama.NewConfig()
}

func NewConsumer(addrs []string, config *Config) (Consumer, error) {
	return originalSarama.NewConsumer(addrs, config)
}

func NewConsumerFromClient(client Client) (Consumer, error) {
	return originalSarama.NewConsumerFromClient(client)
}

func NewConsumerGroup(addrs []string, groupID string, config *Config) (ConsumerGroup, error) {
	consumerGroup, error := originalSarama.NewConsumerGroup(addrs, groupID, config)

	if consumerGroup != nil {
		consumerGroup = consumerGroupWrapper{consumerGroup}
	}

	return consumerGroup, error
}

func NewConsumerGroupFromClient(groupID string, client Client) (ConsumerGroup, error) {
	consumerGroup, error := originalSarama.NewConsumerGroupFromClient(groupID, client)

	if consumerGroup != nil {
		consumerGroup = consumerGroupWrapper{consumerGroup}
	}

	return consumerGroup, error
}

func NewCustomHashPartitioner(hasher func() hash.Hash32) PartitionerConstructor {
	return originalSarama.NewCustomHashPartitioner(hasher)
}

func NewCustomPartitioner(options ...HashPartitionerOption) PartitionerConstructor {
	return originalSarama.NewCustomPartitioner(options...)
}

func NewHashPartitioner(topic string) Partitioner {
	return originalSarama.NewHashPartitioner(topic)
}

func NewKerberosClient(config *GSSAPIConfig) (KerberosClient, error) {
	return originalSarama.NewKerberosClient(config)
}

func NewManualPartitioner(topic string) Partitioner {
	return originalSarama.NewManualPartitioner(topic)
}

func NewMockAlterConfigsResponse(t TestReporter) *MockAlterConfigsResponse {
	return originalSarama.NewMockAlterConfigsResponse(t)
}

func NewMockAlterConfigsResponseWithErrorCode(t TestReporter) *MockAlterConfigsResponseWithErrorCode {
	return originalSarama.NewMockAlterConfigsResponseWithErrorCode(t)
}

func NewMockAlterPartitionReassignmentsResponse(t TestReporter) *MockAlterPartitionReassignmentsResponse {
	return originalSarama.NewMockAlterPartitionReassignmentsResponse(t)
}

func NewMockApiVersionsResponse(t TestReporter) *MockApiVersionsResponse {
	return originalSarama.NewMockApiVersionsResponse(t)
}

func NewMockBroker(t TestReporter, brokerID int32) *MockBroker {
	return originalSarama.NewMockBroker(t, brokerID)
}

func NewMockBrokerAddr(t TestReporter, brokerID int32, addr string) *MockBroker {
	return originalSarama.NewMockBrokerAddr(t, brokerID, addr)
}

func NewMockBrokerListener(t TestReporter, brokerID int32, listener net.Listener) *MockBroker {
	return originalSarama.NewMockBrokerListener(t, brokerID, listener)
}

func NewMockConsumerMetadataResponse(t TestReporter) *MockConsumerMetadataResponse {
	return originalSarama.NewMockConsumerMetadataResponse(t)
}

func NewMockCreateAclsResponse(t TestReporter) *MockCreateAclsResponse {
	return originalSarama.NewMockCreateAclsResponse(t)
}

func NewMockCreateAclsResponseWithError(t TestReporter) *MockCreateAclsResponseError {
	return originalSarama.NewMockCreateAclsResponseWithError(t)
}

func NewMockCreatePartitionsResponse(t TestReporter) *MockCreatePartitionsResponse {
	return originalSarama.NewMockCreatePartitionsResponse(t)
}

func NewMockCreateTopicsResponse(t TestReporter) *MockCreateTopicsResponse {
	return originalSarama.NewMockCreateTopicsResponse(t)
}

func NewMockDeleteAclsResponse(t TestReporter) *MockDeleteAclsResponse {
	return originalSarama.NewMockDeleteAclsResponse(t)
}

func NewMockDeleteGroupsRequest(t TestReporter) *MockDeleteGroupsResponse {
	return originalSarama.NewMockDeleteGroupsRequest(t)
}

func NewMockDeleteOffsetRequest(t TestReporter) *MockDeleteOffsetResponse {
	return originalSarama.NewMockDeleteOffsetRequest(t)
}

func NewMockDeleteRecordsResponse(t TestReporter) *MockDeleteRecordsResponse {
	return originalSarama.NewMockDeleteRecordsResponse(t)
}

func NewMockDeleteTopicsResponse(t TestReporter) *MockDeleteTopicsResponse {
	return originalSarama.NewMockDeleteTopicsResponse(t)
}

func NewMockDescribeConfigsResponse(t TestReporter) *MockDescribeConfigsResponse {
	return originalSarama.NewMockDescribeConfigsResponse(t)
}

func NewMockDescribeConfigsResponseWithErrorCode(t TestReporter) *MockDescribeConfigsResponseWithErrorCode {
	return originalSarama.NewMockDescribeConfigsResponseWithErrorCode(t)
}

func NewMockDescribeGroupsResponse(t TestReporter) *MockDescribeGroupsResponse {
	return originalSarama.NewMockDescribeGroupsResponse(t)
}

func NewMockDescribeLogDirsResponse(t TestReporter) *MockDescribeLogDirsResponse {
	return originalSarama.NewMockDescribeLogDirsResponse(t)
}

func NewMockFetchResponse(t TestReporter, batchSize int) *MockFetchResponse {
	return originalSarama.NewMockFetchResponse(t, batchSize)
}

func NewMockFindCoordinatorResponse(t TestReporter) *MockFindCoordinatorResponse {
	return originalSarama.NewMockFindCoordinatorResponse(t)
}

func NewMockHeartbeatResponse(t TestReporter) *MockHeartbeatResponse {
	return originalSarama.NewMockHeartbeatResponse(t)
}

func NewMockIncrementalAlterConfigsResponse(t TestReporter) *MockIncrementalAlterConfigsResponse {
	return originalSarama.NewMockIncrementalAlterConfigsResponse(t)
}

func NewMockIncrementalAlterConfigsResponseWithErrorCode(t TestReporter) *MockIncrementalAlterConfigsResponseWithErrorCode {
	return originalSarama.NewMockIncrementalAlterConfigsResponseWithErrorCode(t)
}

func NewMockJoinGroupResponse(t TestReporter) *MockJoinGroupResponse {
	return originalSarama.NewMockJoinGroupResponse(t)
}

func NewMockLeaveGroupResponse(t TestReporter) *MockLeaveGroupResponse {
	return originalSarama.NewMockLeaveGroupResponse(t)
}

func NewMockListAclsResponse(t TestReporter) *MockListAclsResponse {
	return originalSarama.NewMockListAclsResponse(t)
}

func NewMockListGroupsResponse(t TestReporter) *MockListGroupsResponse {
	return originalSarama.NewMockListGroupsResponse(t)
}

func NewMockListPartitionReassignmentsResponse(t TestReporter) *MockListPartitionReassignmentsResponse {
	return originalSarama.NewMockListPartitionReassignmentsResponse(t)
}

func NewMockMetadataResponse(t TestReporter) *MockMetadataResponse {
	return originalSarama.NewMockMetadataResponse(t)
}

func NewMockOffsetCommitResponse(t TestReporter) *MockOffsetCommitResponse {
	return originalSarama.NewMockOffsetCommitResponse(t)
}

func NewMockOffsetFetchResponse(t TestReporter) *MockOffsetFetchResponse {
	return originalSarama.NewMockOffsetFetchResponse(t)
}

func NewMockOffsetResponse(t TestReporter) *MockOffsetResponse {
	return originalSarama.NewMockOffsetResponse(t)
}

func NewMockProduceResponse(t TestReporter) *MockProduceResponse {
	return originalSarama.NewMockProduceResponse(t)
}

func NewMockSaslAuthenticateResponse(t TestReporter) *MockSaslAuthenticateResponse {
	return originalSarama.NewMockSaslAuthenticateResponse(t)
}

func NewMockSaslHandshakeResponse(t TestReporter) *MockSaslHandshakeResponse {
	return originalSarama.NewMockSaslHandshakeResponse(t)
}

func NewMockSequence(responses ...interface{}) *MockSequence {
	return originalSarama.NewMockSequence(responses)
}

func NewMockSyncGroupResponse(t TestReporter) *MockSyncGroupResponse {
	return originalSarama.NewMockSyncGroupResponse(t)
}

// Cannot create a proxy for "NewMockWrapper" because the parameter's type "encoderWithHeader" is not exported.
// func NewMockWrapper(res encoderWithHeader) *MockWrapper {
// 	return originalSarama.NewMockWrapper(res)
// }

func NewOffsetManagerFromClient(group string, client Client) (OffsetManager, error) {
	return originalSarama.NewOffsetManagerFromClient(group, client)
}

func NewRandomPartitioner(topic string) Partitioner {
	return originalSarama.NewRandomPartitioner(topic)
}

func NewReferenceHashPartitioner(topic string) Partitioner {
	return originalSarama.NewReferenceHashPartitioner(topic)
}

func NewRoundRobinPartitioner(topic string) Partitioner {
	return originalSarama.NewRoundRobinPartitioner(topic)
}

func NewSyncProducer(addrs []string, config *Config) (SyncProducer, error) {
	syncProducer, error := originalSarama.NewSyncProducer(addrs, config)

	if syncProducer != nil {
		syncProducer = otelsarama.WrapSyncProducer(config, syncProducer)
	}

	return syncProducer, error
}

func NewSyncProducerFromClient(client Client) (SyncProducer, error) {
	syncProducer, error := originalSarama.NewSyncProducerFromClient(client)

	if syncProducer != nil {
		syncProducer = otelsarama.WrapSyncProducer(client.Config(), syncProducer)
	}

	return syncProducer, error
}

func ParseKafkaVersion(s string) (KafkaVersion, error) {
	return originalSarama.ParseKafkaVersion(s)
}

func WithAbsFirst() HashPartitionerOption {
	return originalSarama.WithAbsFirst()
}

func WithCustomFallbackPartitioner(randomHP Partitioner) HashPartitionerOption {
	return originalSarama.WithCustomFallbackPartitioner(randomHP)
}

func WithCustomHashFunction(hasher func() hash.Hash32) HashPartitionerOption {
	return originalSarama.WithCustomHashFunction(hasher)
}

// The original return type of "Wrap" is "sentinelError", but the proxy cannot access it, so it returns an interface.
func Wrap(sentinel error, wrapped ...error) interface{} {
	return originalSarama.Wrap(sentinel, wrapped...)
}

func InjectContextToMessage(ctx context.Context, message *originalSarama.ProducerMessage) {
	carrier := otelsarama.NewProducerMessageCarrier(message)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}
