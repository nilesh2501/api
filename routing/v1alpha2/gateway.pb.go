// Code generated by protoc-gen-go. DO NOT EDIT.
// source: routing/v1alpha2/gateway.proto

/*
Package istio_routing_v1alpha2 is a generated protocol buffer package.

It is generated from these files:
	routing/v1alpha2/gateway.proto
	routing/v1alpha2/route_rule.proto

It has these top-level messages:
	Gateway
	Server
	RouteRule
	Destination
	HTTPRoute
	TCPRoute
	HTTPMatchRequest
	DestinationWeight
	L4MatchAttributes
	HTTPRedirect
	HTTPRewrite
	StringMatch
	HTTPRetry
	CorsPolicy
	HTTPFaultInjection
	PortSelector
*/
package istio_routing_v1alpha2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// TLS modes enforced by the gateway
type Server_TLSOptions_TLSmode int32

const (
	// If set to "passthrough", the gateway will forward the connection to the
	// upstream server as is.
	Server_TLSOptions_PASSTHROUGH Server_TLSOptions_TLSmode = 0
	// If set to "simple", the gateway will secure connections with
	// standard TLS semantics (server certs only).
	Server_TLSOptions_SIMPLE Server_TLSOptions_TLSmode = 1
	// If set to "mutual", the gateway will use standard mTLS authentication.
	Server_TLSOptions_MUTUAL Server_TLSOptions_TLSmode = 2
)

var Server_TLSOptions_TLSmode_name = map[int32]string{
	0: "PASSTHROUGH",
	1: "SIMPLE",
	2: "MUTUAL",
}
var Server_TLSOptions_TLSmode_value = map[string]int32{
	"PASSTHROUGH": 0,
	"SIMPLE":      1,
	"MUTUAL":      2,
}

func (x Server_TLSOptions_TLSmode) String() string {
	return proto.EnumName(Server_TLSOptions_TLSmode_name, int32(x))
}
func (Server_TLSOptions_TLSmode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{1, 1, 0}
}

// Gateway describes a load balancer operating at the edge of the mesh
// receiving incoming or outgoing HTTP/TCP connections. The specification
// describes a set of ports that should be exposed, the type of protocol to
// use, SNI configuration for the load balancer, etc.
//
// For example, the following gateway spec sets up a proxy to act as a load
// balancer exposing port 80 and 9080 (http), 443 (https), and port 2379
// (TCP) for ingress.  While Istio will configure the proxy to listen on
// these ports, it is the responsibility of the user to ensure that
// external traffic to these ports are allowed into the mesh.
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: Gateway
//     metadata:
//       name: my-gateway
//     spec:
//       servers:
//       - port:
//           number: 80
//           name: http
//         domains:
//         - uk.bookinfo.com
//         - eu.bookinfo.com
//         tls:
//           httpsRedirect: true # sends 302 redirect for http requests
//       - port:
//           number: 443
//           name: https
//         domains:
//         - uk.bookinfo.com
//         - eu.bookinfo.com
//         tls:
//           mode: simple #enables HTTPS on this port
//           serverCert: server.crt
//           clientCABundle: client.ca-bundle
//       - port:
//           number: 9080
//           name: http-wildcard
//         # no domains implies wildcard match
//       - port:
//           number: 2379 #to expose internal service via external port 2379
//           name: Mongo
//           protocol: MONGO
//
// The gateway specification above describes the L4-L6 properties of a load
// balancer. Routing rules can then be bound to a gateway to control
// the forwarding of traffic arriving at a particular domain or gateway port.
//
// The following sample route rule splits traffic for
// https://uk.bookinfo.com/reviews, https://eu.bookinfo.com/reviews,
// http://uk.bookinfo.com:9080/reviews, http://eu.bookinfo.com:9080/reviews
// into two versions (prod and qa) of an internal reviews service on port
// 9080. In addition, requests containing the cookie user: dev-123 will be
// sent to special port 7777 in the qa version. The same rule is also
// applicable inside the mesh for requests to the reviews.prod
// service. This rule is applicable across ports 443, 9080. Note that
// http://uk.bookinfo.com gets redirected to https://uk.bookinfo.com
// (i.e. 80 redirects to 443).
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: RouteRule
//     metadata:
//       name: bookinfo-rule
//     spec:
//       hosts:
//       - reviews.prod
//       - uk.bookinfo.com
//       - eu.bookinfo.com
//       gateways:
//       - my-gateway
//       http:
//       - match:
//         - headers:
//             cookie:
//               user: dev-123
//         route:
//         - destination:
//             port:
//               number: 7777
//             name: reviews.qa
//       - match:
//           uri:
//             prefix: /reviews/
//         route:
//         - destination:
//             port:
//               number: 9080 # can be omitted if its the only port for reviews
//             name: reviews.prod
//           weight: 80
//         - destination:
//             name: reviews.qa
//           weight: 20
//
// The following routing rule forwards traffic arriving at (external)
// port 2379 from 172.17.16.* subnet to internal Mongo server on port 5555.
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: RouteRule
//     metadata:
//       name: bookinfo-Mongo
//     spec:
//       hosts:
//       - Mongosvr #name of Mongo service
//       gateways:
//       - my-gateway
//       tcp:
//       - match:
//         - port:
//             number: 2379
//           sourceSubnet: "172.17.16.0/24"
//         route:
//         - destination:
//             name: mongo.prod
//
// By default, if there is no wildcard, HTTP requests for unknown domains
// or requests that have no matching route rule will respond with a 404. If
// a specific default behavior is desired at the ingress, add a route rule
// with a wildcard host and the desired backend. For example, the following
// wildcard routing rule routes all traffic to homepage.prod by default.
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: RouteRule
//     metadata:
//       name: default-ingress
//     spec:
//       hosts:
//       - *
//       gateways:
//       - my-gateway
//       http:
//       - route:
//         - destination:
//             name: homepage.prod
//
type Gateway struct {
	// REQUIRED: A list of server specifications.
	Servers []*Server `protobuf:"bytes,1,rep,name=servers" json:"servers,omitempty"`
}

func (m *Gateway) Reset()                    { *m = Gateway{} }
func (m *Gateway) String() string            { return proto.CompactTextString(m) }
func (*Gateway) ProtoMessage()               {}
func (*Gateway) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Gateway) GetServers() []*Server {
	if m != nil {
		return m.Servers
	}
	return nil
}

// Server describes the properties of the proxy on a given load balancer port.
// For example,
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: Gateway
//     metadata:
//       name: my-ingress
//     spec:
//       servers:
//       - port:
//           number: 80
//           protocol: HTTP2
//
// Another example
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: Gateway
//     metadata:
//       name: my-tcp-ingress
//     spec:
//       servers:
//       - port:
//           number: 27018
//           protocol: MONGO
//
// The following is an example of TLS configuration for port 443
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: Gateway
//     metadata:
//       name: my-ingress
//     spec:
//       servers:
//       - port:
//           number: 443
//           protocol: HTTP
//         tls:
//           mode: simple
//           serverCertificatet: server.crt
//
type Server struct {
	// REQUIRED: The Port on which the proxy should listen for incoming
	// connections
	Port *Server_Port `protobuf:"bytes,1,opt,name=port" json:"port,omitempty"`
	// A list of domains exposed by this gateway. While
	// typically applicable to HTTP services, it can also be used for TCP
	// services using TLS with SNI. Standard DNS wildcard prefix syntax
	// is permitted.
	//
	// RouteRules that are bound to a gateway must having a matching domain
	// in their default destination. Specifically one of the route rule
	// destination domains is a strict suffix of a gateway domain or
	// a gateway domain is a suffix of one of the route rule domains.
	Domains []string `protobuf:"bytes,2,rep,name=domains" json:"domains,omitempty"`
	// Set of TLS related options that govern the server's behavior. Use
	// these options to control if all http requests should be redirected to
	// https, and the TLS modes to use.
	Tls *Server_TLSOptions `protobuf:"bytes,3,opt,name=tls" json:"tls,omitempty"`
}

func (m *Server) Reset()                    { *m = Server{} }
func (m *Server) String() string            { return proto.CompactTextString(m) }
func (*Server) ProtoMessage()               {}
func (*Server) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Server) GetPort() *Server_Port {
	if m != nil {
		return m.Port
	}
	return nil
}

func (m *Server) GetDomains() []string {
	if m != nil {
		return m.Domains
	}
	return nil
}

func (m *Server) GetTls() *Server_TLSOptions {
	if m != nil {
		return m.Tls
	}
	return nil
}

// Port describes the properties of a specific port of a service.
type Server_Port struct {
	// REQUIRED: A valid non-negative integer port number.
	Number uint32 `protobuf:"varint,1,opt,name=number" json:"number,omitempty"`
	// The protocol exposed on the port.
	// MUST BE one of HTTP|HTTPS|GRPC|HTTP2|MONGO|TCP.
	Protocol string `protobuf:"bytes,2,opt,name=protocol" json:"protocol,omitempty"`
	// Label assigned to the port.
	Name string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *Server_Port) Reset()                    { *m = Server_Port{} }
func (m *Server_Port) String() string            { return proto.CompactTextString(m) }
func (*Server_Port) ProtoMessage()               {}
func (*Server_Port) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *Server_Port) GetNumber() uint32 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *Server_Port) GetProtocol() string {
	if m != nil {
		return m.Protocol
	}
	return ""
}

func (m *Server_Port) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Server_TLSOptions struct {
	// If set to true, the load balancer will send a 302 redirect for all
	// http connections, asking the clients to use HTTPS.
	HttpsRedirect bool `protobuf:"varint,1,opt,name=https_redirect,json=httpsRedirect" json:"https_redirect,omitempty"`
	// Optional: Indicates whether connections to this port should be
	// secured using TLS.  The value of this field determines how TLS is
	// enforced.
	Mode Server_TLSOptions_TLSmode `protobuf:"varint,2,opt,name=mode,enum=istio.routing.v1alpha2.Server_TLSOptions_TLSmode" json:"mode,omitempty"`
	// REQUIRED if mode == SIMPLE/MUTUAL. The name of the file holding the
	// server-side TLS certificate to use.  It is the responsibility of the
	// underlying platform to mount the certificate as a file under
	// /etc/istio/ingress-certs with the same name as the specified in this
	// field.
	ServerCertificate string `protobuf:"bytes,3,opt,name=server_certificate,json=serverCertificate" json:"server_certificate,omitempty"`
	// REQUIRED if mode == MUTUAL. To use mutual TLS for external clients,
	// specify the name of the file holding the CA certificate to validate
	// the client's certificate. It is the responsibility of the underlying
	// platform to mount the certificate as a file under
	// /etc/istio/ingress-certs with the same name as specified in this
	// field.
	ClientCaBundle string `protobuf:"bytes,4,opt,name=client_ca_bundle,json=clientCaBundle" json:"client_ca_bundle,omitempty"`
}

func (m *Server_TLSOptions) Reset()                    { *m = Server_TLSOptions{} }
func (m *Server_TLSOptions) String() string            { return proto.CompactTextString(m) }
func (*Server_TLSOptions) ProtoMessage()               {}
func (*Server_TLSOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 1} }

func (m *Server_TLSOptions) GetHttpsRedirect() bool {
	if m != nil {
		return m.HttpsRedirect
	}
	return false
}

func (m *Server_TLSOptions) GetMode() Server_TLSOptions_TLSmode {
	if m != nil {
		return m.Mode
	}
	return Server_TLSOptions_PASSTHROUGH
}

func (m *Server_TLSOptions) GetServerCertificate() string {
	if m != nil {
		return m.ServerCertificate
	}
	return ""
}

func (m *Server_TLSOptions) GetClientCaBundle() string {
	if m != nil {
		return m.ClientCaBundle
	}
	return ""
}

func init() {
	proto.RegisterType((*Gateway)(nil), "istio.routing.v1alpha2.Gateway")
	proto.RegisterType((*Server)(nil), "istio.routing.v1alpha2.Server")
	proto.RegisterType((*Server_Port)(nil), "istio.routing.v1alpha2.Server.Port")
	proto.RegisterType((*Server_TLSOptions)(nil), "istio.routing.v1alpha2.Server.TLSOptions")
	proto.RegisterEnum("istio.routing.v1alpha2.Server_TLSOptions_TLSmode", Server_TLSOptions_TLSmode_name, Server_TLSOptions_TLSmode_value)
}

func init() { proto.RegisterFile("routing/v1alpha2/gateway.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 381 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x50, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xc5, 0x1f, 0xb2, 0x9b, 0x89, 0x1a, 0xcc, 0x1c, 0x2a, 0x2b, 0x87, 0xca, 0x0a, 0x42, 0x32,
	0x07, 0x5c, 0xd5, 0x1c, 0x40, 0xe2, 0x54, 0xa2, 0xaa, 0x41, 0x4a, 0x48, 0xb4, 0x4e, 0xce, 0xd6,
	0xc6, 0x5e, 0x92, 0x95, 0x6c, 0xaf, 0xb5, 0xde, 0x04, 0xf1, 0x0b, 0xf8, 0xbf, 0xfc, 0x02, 0xe4,
	0xb5, 0x4d, 0x2e, 0x08, 0x7a, 0x9b, 0x79, 0xf3, 0xde, 0xbc, 0x99, 0x07, 0xb7, 0x52, 0x9c, 0x14,
	0xaf, 0x0e, 0x77, 0xe7, 0x7b, 0x5a, 0xd4, 0x47, 0x1a, 0xdf, 0x1d, 0xa8, 0x62, 0xdf, 0xe9, 0x8f,
	0xa8, 0x96, 0x42, 0x09, 0xbc, 0xe1, 0x8d, 0xe2, 0x22, 0xea, 0x59, 0xd1, 0xc0, 0x9a, 0xcd, 0xc1,
	0x7d, 0xea, 0x88, 0xf8, 0x11, 0xdc, 0x86, 0xc9, 0x33, 0x93, 0x8d, 0x6f, 0x04, 0x56, 0x38, 0x8e,
	0x6f, 0xa3, 0xbf, 0x8b, 0xa2, 0x44, 0xd3, 0xc8, 0x40, 0x9f, 0xfd, 0xb2, 0xc0, 0xe9, 0x30, 0xfc,
	0x00, 0x76, 0x2d, 0xa4, 0xf2, 0x8d, 0xc0, 0x08, 0xc7, 0xf1, 0xeb, 0x7f, 0x6f, 0x88, 0x36, 0x42,
	0x2a, 0xa2, 0x05, 0xe8, 0x83, 0x9b, 0x8b, 0x92, 0xf2, 0xaa, 0xf1, 0xcd, 0xc0, 0x0a, 0x47, 0x64,
	0x68, 0xf1, 0x13, 0x58, 0xaa, 0x68, 0x7c, 0x4b, 0x6f, 0x7c, 0xfb, 0x9f, 0x8d, 0xdb, 0x65, 0xb2,
	0xae, 0x15, 0x17, 0x55, 0x43, 0x5a, 0xd5, 0xf4, 0x2b, 0xd8, 0xad, 0x09, 0xde, 0x80, 0x53, 0x9d,
	0xca, 0x3d, 0x93, 0xfa, 0xb2, 0x6b, 0xd2, 0x77, 0x38, 0x85, 0x2b, 0x1d, 0x50, 0x26, 0x0a, 0xdf,
	0x0c, 0x8c, 0x70, 0x44, 0xfe, 0xf4, 0x88, 0x60, 0x57, 0xb4, 0x64, 0xda, 0x79, 0x44, 0x74, 0x3d,
	0xfd, 0x69, 0x02, 0x5c, 0x3c, 0xf0, 0x0d, 0x4c, 0x8e, 0x4a, 0xd5, 0x4d, 0x2a, 0x59, 0xce, 0x25,
	0xcb, 0xba, 0xc7, 0xaf, 0xc8, 0xb5, 0x46, 0x49, 0x0f, 0xe2, 0x23, 0xd8, 0xa5, 0xc8, 0x99, 0x76,
	0x98, 0xc4, 0xf7, 0xcf, 0xfe, 0xa1, 0x2d, 0x5b, 0x21, 0xd1, 0x72, 0x7c, 0x07, 0xd8, 0x45, 0x9e,
	0x66, 0x4c, 0x2a, 0xfe, 0x8d, 0x67, 0x54, 0x0d, 0xe7, 0xbd, 0xea, 0x26, 0xf3, 0xcb, 0x00, 0x43,
	0xf0, 0xb2, 0x82, 0xb3, 0x4a, 0xa5, 0x19, 0x4d, 0xf7, 0xa7, 0x2a, 0x2f, 0x98, 0x6f, 0x6b, 0xf2,
	0xa4, 0xc3, 0xe7, 0xf4, 0xb3, 0x46, 0x67, 0x31, 0xb8, 0xbd, 0x13, 0xbe, 0x84, 0xf1, 0xe6, 0x21,
	0x49, 0xb6, 0x0b, 0xb2, 0xde, 0x3d, 0x2d, 0xbc, 0x17, 0x08, 0xe0, 0x24, 0x5f, 0x56, 0x9b, 0xe5,
	0xa3, 0x67, 0xb4, 0xf5, 0x6a, 0xb7, 0xdd, 0x3d, 0x2c, 0x3d, 0x73, 0xef, 0xe8, 0x9c, 0xde, 0xff,
	0x0e, 0x00, 0x00, 0xff, 0xff, 0x29, 0xe0, 0x45, 0x26, 0x7a, 0x02, 0x00, 0x00,
}
