package heliosgrpc

import (
	"context"
	"net"
	"time"

	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	realGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/channelz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/tap"
)

const (
	SupportPackageIsVersion3 = realGrpc.SupportPackageIsVersion3
	SupportPackageIsVersion4 = realGrpc.SupportPackageIsVersion4
	SupportPackageIsVersion5 = realGrpc.SupportPackageIsVersion5
	SupportPackageIsVersion6 = realGrpc.SupportPackageIsVersion6
	SupportPackageIsVersion7 = realGrpc.SupportPackageIsVersion7
)

const PickFirstBalancerName = realGrpc.PickFirstBalancerName
const Version = realGrpc.Version

var DefaultBackoffConfig = realGrpc.DefaultBackoffConfig
var EnableTracing = realGrpc.EnableTracing
var ErrClientConnClosing = realGrpc.ErrClientConnClosing
var ErrClientConnTimeout = realGrpc.ErrClientConnTimeout
var ErrServerStopped = realGrpc.ErrServerStopped

type CallOption = realGrpc.CallOption
type ClientConn = realGrpc.ClientConn
type ServerStream = realGrpc.ServerStream
type ServerTransportStream = realGrpc.ServerTransportStream
type DialOption = realGrpc.DialOption
type ClientStream = realGrpc.ClientStream
type StreamDesc = realGrpc.StreamDesc
type CompressorCallOption = realGrpc.CompressorCallOption
type ConnectParams = realGrpc.ConnectParams
type ContentSubtypeCallOption = realGrpc.ContentSubtypeCallOption
type CustomCodecCallOption = realGrpc.CustomCodecCallOption
type StreamClientInterceptor = realGrpc.StreamClientInterceptor
type UnaryClientInterceptor = realGrpc.UnaryClientInterceptor
type EmptyCallOption = realGrpc.EmptyCallOption
type EmptyDialOption = realGrpc.EmptyDialOption
type EmptyServerOption = realGrpc.EmptyServerOption
type FailFastCallOption = realGrpc.FailFastCallOption
type ForceCodecCallOption = realGrpc.ForceCodecCallOption
type HeaderCallOption = realGrpc.HeaderCallOption
type MaxRecvMsgSizeCallOption = realGrpc.MaxRecvMsgSizeCallOption
type MaxRetryRPCBufferSizeCallOption = realGrpc.MaxRetryRPCBufferSizeCallOption
type MaxSendMsgSizeCallOption = realGrpc.MaxSendMsgSizeCallOption
type MethodDesc = realGrpc.MethodDesc
type MethodInfo = realGrpc.MethodInfo
type PeerCallOption = realGrpc.PeerCallOption
type PerRPCCredsCallOption = realGrpc.PerRPCCredsCallOption
type PreparedMsg = realGrpc.PreparedMsg
type Server = realGrpc.Server
type ServerOption = realGrpc.ServerOption
type StreamServerInterceptor = realGrpc.StreamServerInterceptor
type UnaryServerInterceptor = realGrpc.UnaryServerInterceptor
type StreamHandler = realGrpc.StreamHandler
type ServiceDesc = realGrpc.ServiceDesc
type ServiceInfo = realGrpc.ServiceInfo
type ServiceRegistrar = realGrpc.ServiceRegistrar
type StreamServerInfo = realGrpc.StreamServerInfo
type Streamer = realGrpc.Streamer
type TrailerCallOption = realGrpc.TrailerCallOption
type UnaryHandler = realGrpc.UnaryHandler
type UnaryInvoker = realGrpc.UnaryInvoker
type UnaryServerInfo = realGrpc.UnaryServerInfo
type BackoffConfig = realGrpc.BackoffConfig
type Codec = realGrpc.Codec
type Compressor = realGrpc.Compressor
type Decompressor = realGrpc.Decompressor
type MethodConfig = realGrpc.MethodConfig
type ServiceConfig = realGrpc.ServiceConfig
type Stream = realGrpc.Stream
type ClientConnInterface = realGrpc.ClientConnInterface

func Invoke(ctx context.Context, method string, args, reply interface{}, cc *ClientConn, opts ...CallOption) error {
	return realGrpc.Invoke(ctx, method, args, reply, cc, opts...)
}

func NewGZIPDecompressor() Decompressor {
	return realGrpc.NewGZIPDecompressor()
}

func NewGZIPCompressor() Compressor {
	return realGrpc.NewGZIPCompressor()
}

func NewGZIPCompressorWithLevel(level int) (Compressor, error) {
	return realGrpc.NewGZIPCompressorWithLevel(level)
}

func WithDialer(f func(string, time.Duration) (net.Conn, error)) DialOption {
	return realGrpc.WithDialer(f)
}

func WithBackoffConfig(b BackoffConfig) DialOption {
	return realGrpc.WithBackoffConfig(b)
}

func WithBackoffMaxDelay(md time.Duration) DialOption {
	return realGrpc.WithBackoffMaxDelay(md)
}

func Method(ctx context.Context) (string, bool) {
	return realGrpc.Method(ctx)
}

func MethodFromServerStream(stream ServerStream) (string, bool) {
	return realGrpc.MethodFromServerStream(stream)
}

func NewContextWithServerTransportStream(ctx context.Context, stream ServerTransportStream) context.Context {
	return realGrpc.NewContextWithServerTransportStream(ctx, stream)
}

func SendHeader(ctx context.Context, md metadata.MD) error {
	return realGrpc.SendHeader(ctx, md)
}

func SetHeader(ctx context.Context, md metadata.MD) error {
	return realGrpc.SetHeader(ctx, md)
}

func SetTrailer(ctx context.Context, md metadata.MD) error {
	return realGrpc.SetTrailer(ctx, md)
}

func CallContentSubtype(contentSubtype string) CallOption {
	return realGrpc.CallContentSubtype(contentSubtype)
}

func ForceCodec(codec encoding.Codec) CallOption {
	return realGrpc.ForceCodec(codec)
}

func Header(md *metadata.MD) CallOption {
	return realGrpc.Header(md)
}

func MaxCallRecvMsgSize(bytes int) CallOption {
	return realGrpc.MaxCallRecvMsgSize(bytes)
}

func MaxCallSendMsgSize(bytes int) CallOption {
	return realGrpc.MaxCallSendMsgSize(bytes)
}

func MaxRetryRPCBufferSize(bytes int) CallOption {
	return realGrpc.MaxRetryRPCBufferSize(bytes)
}

func Peer(p *peer.Peer) CallOption {
	return realGrpc.Peer(p)
}

func PerRPCCredentials(creds credentials.PerRPCCredentials) CallOption {
	return realGrpc.PerRPCCredentials(creds)
}

func Trailer(md *metadata.MD) CallOption {
	return realGrpc.Trailer(md)
}

func UseCompressor(name string) CallOption {
	return realGrpc.UseCompressor(name)
}

func WaitForReady(waitForReady bool) CallOption {
	return realGrpc.WaitForReady(waitForReady)
}

func Dial(target string, opts ...DialOption) (*ClientConn, error) {
	newOptions := append(opts, WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()), WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	return realGrpc.Dial(target, newOptions...)
}

func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {
	return realGrpc.DialContext(ctx, target, opts...)
}

func NewClientStream(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (ClientStream, error) {
	return realGrpc.NewClientStream(ctx, desc, cc, method, opts...)
}

func FailOnNonTempDialError(f bool) DialOption {
	return realGrpc.FailOnNonTempDialError(f)
}

func WithAuthority(a string) DialOption {
	return realGrpc.WithAuthority(a)
}

func WithBlock() DialOption {
	return realGrpc.WithBlock()
}

func WithChainStreamInterceptor(interceptors ...StreamClientInterceptor) DialOption {
	return realGrpc.WithChainStreamInterceptor(interceptors...)
}

func WithChainUnaryInterceptor(interceptors ...UnaryClientInterceptor) DialOption {
	return realGrpc.WithChainUnaryInterceptor(interceptors...)
}

func WithChannelzParentID(id *channelz.Identifier) DialOption {
	return realGrpc.WithChannelzParentID(id)
}

func WithConnectParams(p ConnectParams) DialOption {
	return realGrpc.WithConnectParams(p)
}

func WithContextDialer(f func(context.Context, string) (net.Conn, error)) DialOption {
	return realGrpc.WithContextDialer(f)
}

func WithCredentialsBundle(b credentials.Bundle) DialOption {
	return realGrpc.WithCredentialsBundle(b)
}

func WithDefaultCallOptions(cos ...CallOption) DialOption {
	return realGrpc.WithDefaultCallOptions(cos...)
}

func WithDefaultServiceConfig(s string) DialOption {
	return realGrpc.WithDefaultServiceConfig(s)
}

func WithDisableHealthCheck() DialOption {
	return realGrpc.WithDisableHealthCheck()
}

func WithDisableRetry() DialOption {
	return realGrpc.WithDisableRetry()
}

func WithDisableServiceConfig() DialOption {
	return realGrpc.WithDisableServiceConfig()
}

func WithInitialConnWindowSize(s int32) DialOption {
	return realGrpc.WithInitialConnWindowSize(s)
}

func WithInitialWindowSize(s int32) DialOption {
	return realGrpc.WithInitialWindowSize(s)
}

func WithKeepaliveParams(kp keepalive.ClientParameters) DialOption {
	return realGrpc.WithKeepaliveParams(kp)
}

func WithMaxHeaderListSize(s uint32) DialOption {
	return realGrpc.WithMaxHeaderListSize(s)
}

func WithNoProxy() DialOption {
	return realGrpc.WithNoProxy()
}

func WithPerRPCCredentials(creds credentials.PerRPCCredentials) DialOption {
	return realGrpc.WithPerRPCCredentials(creds)
}

func WithReadBufferSize(s int) DialOption {
	return realGrpc.WithReadBufferSize(s)
}

func WithResolvers(rs ...resolver.Builder) DialOption {
	return realGrpc.WithResolvers(rs...)
}

func WithReturnConnectionError() DialOption {
	return realGrpc.WithReturnConnectionError()
}

func WithStatsHandler(h stats.Handler) DialOption {
	return realGrpc.WithStatsHandler(h)
}

func WithStreamInterceptor(f StreamClientInterceptor) DialOption {
	return realGrpc.WithStreamInterceptor(f)
}

func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return realGrpc.WithTransportCredentials(creds)
}

func WithUnaryInterceptor(f UnaryClientInterceptor) DialOption {
	return realGrpc.WithUnaryInterceptor(f)
}

func WithUserAgent(s string) DialOption {
	return realGrpc.WithUserAgent(s)
}

func WithWriteBufferSize(s int) DialOption {
	return realGrpc.WithWriteBufferSize(s)
}

func NewServer(opt ...ServerOption) *Server {
	newOptions := append(opt, ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor()), ChainStreamInterceptor(otelgrpc.StreamServerInterceptor()))
	return realGrpc.NewServer(newOptions...)
}

func ChainStreamInterceptor(interceptors ...StreamServerInterceptor) ServerOption {
	return realGrpc.ChainStreamInterceptor(interceptors...)
}

func ChainUnaryInterceptor(interceptors ...UnaryServerInterceptor) ServerOption {
	return realGrpc.ChainUnaryInterceptor(interceptors...)
}

func ConnectionTimeout(d time.Duration) ServerOption {
	return realGrpc.ConnectionTimeout(d)
}

func Creds(c credentials.TransportCredentials) ServerOption {
	return realGrpc.Creds(c)
}

func ForceServerCodec(codec encoding.Codec) ServerOption {
	return realGrpc.ForceServerCodec(codec)
}

func HeaderTableSize(s uint32) ServerOption {
	return realGrpc.HeaderTableSize(s)
}

func InTapHandle(h tap.ServerInHandle) ServerOption {
	return realGrpc.InTapHandle(h)
}

func InitialConnWindowSize(s int32) ServerOption {
	return realGrpc.InitialConnWindowSize(s)
}

func InitialWindowSize(s int32) ServerOption {
	return realGrpc.InitialWindowSize(s)
}

func KeepaliveEnforcementPolicy(kep keepalive.EnforcementPolicy) ServerOption {
	return realGrpc.KeepaliveEnforcementPolicy(kep)
}

func KeepaliveParams(kp keepalive.ServerParameters) ServerOption {
	return realGrpc.KeepaliveParams(kp)
}

func MaxConcurrentStreams(n uint32) ServerOption {
	return realGrpc.MaxConcurrentStreams(n)
}

func MaxHeaderListSize(s uint32) ServerOption {
	return realGrpc.MaxHeaderListSize(s)
}

func MaxRecvMsgSize(m int) ServerOption {
	return realGrpc.MaxRecvMsgSize(m)
}

func MaxSendMsgSize(m int) ServerOption {
	return realGrpc.MaxSendMsgSize(m)
}

func NumStreamWorkers(numServerWorkers uint32) ServerOption {
	return realGrpc.NumStreamWorkers(numServerWorkers)
}

func ReadBufferSize(s int) ServerOption {
	return realGrpc.ReadBufferSize(s)
}

func StatsHandler(h stats.Handler) ServerOption {
	return realGrpc.StatsHandler(h)
}

func StreamInterceptor(i StreamServerInterceptor) ServerOption {
	return realGrpc.StreamInterceptor(i)
}

func UnaryInterceptor(i UnaryServerInterceptor) ServerOption {
	return realGrpc.UnaryInterceptor(i)
}

func UnknownServiceHandler(streamHandler StreamHandler) ServerOption {
	return realGrpc.UnknownServiceHandler(streamHandler)
}

func WriteBufferSize(s int) ServerOption {
	return realGrpc.WriteBufferSize(s)
}

func ServerTransportStreamFromContext(ctx context.Context) ServerTransportStream {
	return realGrpc.ServerTransportStreamFromContext(ctx)
}

func Code(err error) codes.Code {
	return realGrpc.Code(err)
}

func ErrorDesc(err error) string {
	return realGrpc.ErrorDesc(err)
}

func Errorf(c codes.Code, format string, a ...interface{}) error {
	return realGrpc.Errorf(c, format, a...)
}

func CallCustomCodec(codec Codec) CallOption {
	return realGrpc.CallCustomCodec(codec)
}

func FailFast(failFast bool) CallOption {
	return realGrpc.FailFast(failFast)
}

func WithCodec(c Codec) DialOption {
	return realGrpc.WithCodec(c)
}

func WithCompressor(cp Compressor) DialOption {
	return realGrpc.WithCompressor(cp)
}

func WithDecompressor(dc Decompressor) DialOption {
	return realGrpc.WithDecompressor(dc)
}

func WithInsecure() DialOption {
	return realGrpc.WithInsecure()
}

func WithMaxMsgSize(s int) DialOption {
	return realGrpc.WithMaxMsgSize(s)
}

func WithServiceConfig(c <-chan ServiceConfig) DialOption {
	return realGrpc.WithServiceConfig(c)
}

func WithTimeout(d time.Duration) DialOption {
	return realGrpc.WithTimeout(d)
}

func CustomCodec(codec Codec) ServerOption {
	return realGrpc.CustomCodec(codec)
}

func MaxMsgSize(m int) ServerOption {
	return realGrpc.MaxMsgSize(m)
}

func RPCCompressor(cp Compressor) ServerOption {
	return realGrpc.RPCCompressor(cp)
}

func RPCDecompressor(dc Decompressor) ServerOption {
	return realGrpc.RPCDecompressor(dc)
}
