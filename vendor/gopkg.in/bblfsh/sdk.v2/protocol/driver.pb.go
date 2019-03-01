// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: driver.proto

/*
	Package protocol is a generated protocol buffer package.

	It is generated from these files:
		driver.proto

	It has these top-level messages:
		ParseRequest
		ParseResponse
		ParseError
*/
package protocol

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Mode int32

const (
	// DefaultMode selects the transformation mode that is considered to produce UAST of the best quality.
	Mode_DefaultMode Mode = 0
	// Native disables any UAST transformations and emits a native language AST as returned by the parser.
	Mode_Native Mode = 1
	// Preprocessed runs only basic transformation over native AST (normalize positional info, type fields).
	Mode_Preprocessed Mode = 2
	// Annotated UAST is based on native AST, but provides role annotations for nodes.
	Mode_Annotated Mode = 4
	// Semantic UAST normalizes native AST nodes to a unified structure where possible.
	Mode_Semantic Mode = 8
)

var Mode_name = map[int32]string{
	0: "DEFAULT_MODE",
	1: "NATIVE",
	2: "PREPROCESSED",
	4: "ANNOTATED",
	8: "SEMANTIC",
}
var Mode_value = map[string]int32{
	"DEFAULT_MODE": 0,
	"NATIVE":       1,
	"PREPROCESSED": 2,
	"ANNOTATED":    4,
	"SEMANTIC":     8,
}

func (x Mode) String() string {
	return proto.EnumName(Mode_name, int32(x))
}
func (Mode) EnumDescriptor() ([]byte, []int) { return fileDescriptorDriver, []int{0} }

// ParseRequest is a request to parse a file and get its UAST.
type ParseRequest struct {
	// Content stores the content of a source file. Required.
	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	// Language can be set optionally to disable automatic language detection.
	Language string `protobuf:"bytes,2,opt,name=language,proto3" json:"language,omitempty"`
	// Filename can be set optionally to assist automatic language detection.
	Filename string `protobuf:"bytes,3,opt,name=filename,proto3" json:"filename,omitempty"`
	// Mode sets a transformation pipeline used for UAST.
	Mode Mode `protobuf:"varint,4,opt,name=mode,proto3,enum=gopkg.in.bblfsh.sdk.v2.protocol.Mode" json:"mode,omitempty"`
}

func (m *ParseRequest) Reset()                    { *m = ParseRequest{} }
func (m *ParseRequest) String() string            { return proto.CompactTextString(m) }
func (*ParseRequest) ProtoMessage()               {}
func (*ParseRequest) Descriptor() ([]byte, []int) { return fileDescriptorDriver, []int{0} }

// ParseResponse is the reply to ParseRequest.
type ParseResponse struct {
	// UAST is a binary encoding of the resulting UAST.
	Uast []byte `protobuf:"bytes,1,opt,name=uast,proto3" json:"uast,omitempty"`
	// Language that was automatically detected.
	Language string `protobuf:"bytes,2,opt,name=language,proto3" json:"language,omitempty"`
	// Errors is a list of parsing errors.
	// Only set if parser was able to return a response. Otherwise gRPC error codes are used.
	Errors []*ParseError `protobuf:"bytes,3,rep,name=errors" json:"errors,omitempty"`
}

func (m *ParseResponse) Reset()                    { *m = ParseResponse{} }
func (m *ParseResponse) String() string            { return proto.CompactTextString(m) }
func (*ParseResponse) ProtoMessage()               {}
func (*ParseResponse) Descriptor() ([]byte, []int) { return fileDescriptorDriver, []int{1} }

type ParseError struct {
	// Text is an error message.
	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
}

func (m *ParseError) Reset()                    { *m = ParseError{} }
func (m *ParseError) String() string            { return proto.CompactTextString(m) }
func (*ParseError) ProtoMessage()               {}
func (*ParseError) Descriptor() ([]byte, []int) { return fileDescriptorDriver, []int{2} }

func init() {
	proto.RegisterType((*ParseRequest)(nil), "gopkg.in.bblfsh.sdk.v2.protocol.ParseRequest")
	proto.RegisterType((*ParseResponse)(nil), "gopkg.in.bblfsh.sdk.v2.protocol.ParseResponse")
	proto.RegisterType((*ParseError)(nil), "gopkg.in.bblfsh.sdk.v2.protocol.ParseError")
	proto.RegisterEnum("gopkg.in.bblfsh.sdk.v2.protocol.Mode", Mode_name, Mode_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Driver service

type DriverClient interface {
	// Parse returns an UAST for a given source file.
	Parse(ctx context.Context, in *ParseRequest, opts ...grpc.CallOption) (*ParseResponse, error)
}

type driverClient struct {
	cc *grpc.ClientConn
}

func NewDriverClient(cc *grpc.ClientConn) DriverClient {
	return &driverClient{cc}
}

func (c *driverClient) Parse(ctx context.Context, in *ParseRequest, opts ...grpc.CallOption) (*ParseResponse, error) {
	out := new(ParseResponse)
	err := grpc.Invoke(ctx, "/gopkg.in.bblfsh.sdk.v2.protocol.Driver/Parse", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Driver service

type DriverServer interface {
	// Parse returns an UAST for a given source file.
	Parse(context.Context, *ParseRequest) (*ParseResponse, error)
}

func RegisterDriverServer(s *grpc.Server, srv DriverServer) {
	s.RegisterService(&_Driver_serviceDesc, srv)
}

func _Driver_Parse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Parse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gopkg.in.bblfsh.sdk.v2.protocol.Driver/Parse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Parse(ctx, req.(*ParseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Driver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gopkg.in.bblfsh.sdk.v2.protocol.Driver",
	HandlerType: (*DriverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Parse",
			Handler:    _Driver_Parse_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "driver.proto",
}

func (m *ParseRequest) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ParseRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Content) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Content)))
		i += copy(dAtA[i:], m.Content)
	}
	if len(m.Language) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Language)))
		i += copy(dAtA[i:], m.Language)
	}
	if len(m.Filename) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Filename)))
		i += copy(dAtA[i:], m.Filename)
	}
	if m.Mode != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintDriver(dAtA, i, uint64(m.Mode))
	}
	return i, nil
}

func (m *ParseResponse) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ParseResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Uast) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Uast)))
		i += copy(dAtA[i:], m.Uast)
	}
	if len(m.Language) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Language)))
		i += copy(dAtA[i:], m.Language)
	}
	if len(m.Errors) > 0 {
		for _, msg := range m.Errors {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintDriver(dAtA, i, uint64(msg.ProtoSize()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *ParseError) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ParseError) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Text) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDriver(dAtA, i, uint64(len(m.Text)))
		i += copy(dAtA[i:], m.Text)
	}
	return i, nil
}

func encodeFixed64Driver(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Driver(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintDriver(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ParseRequest) ProtoSize() (n int) {
	var l int
	_ = l
	l = len(m.Content)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	l = len(m.Language)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	l = len(m.Filename)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	if m.Mode != 0 {
		n += 1 + sovDriver(uint64(m.Mode))
	}
	return n
}

func (m *ParseResponse) ProtoSize() (n int) {
	var l int
	_ = l
	l = len(m.Uast)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	l = len(m.Language)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	if len(m.Errors) > 0 {
		for _, e := range m.Errors {
			l = e.ProtoSize()
			n += 1 + l + sovDriver(uint64(l))
		}
	}
	return n
}

func (m *ParseError) ProtoSize() (n int) {
	var l int
	_ = l
	l = len(m.Text)
	if l > 0 {
		n += 1 + l + sovDriver(uint64(l))
	}
	return n
}

func sovDriver(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozDriver(x uint64) (n int) {
	return sovDriver(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ParseRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDriver
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ParseRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ParseRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Content", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Content = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Language", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Language = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Filename", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Filename = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Mode", wireType)
			}
			m.Mode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Mode |= (Mode(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDriver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDriver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ParseResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDriver
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ParseResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ParseResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uast", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uast = append(m.Uast[:0], dAtA[iNdEx:postIndex]...)
			if m.Uast == nil {
				m.Uast = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Language", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Language = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Errors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Errors = append(m.Errors, &ParseError{})
			if err := m.Errors[len(m.Errors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDriver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDriver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ParseError) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDriver
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ParseError: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ParseError: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Text", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDriver
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Text = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDriver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDriver
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipDriver(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDriver
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDriver
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthDriver
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowDriver
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipDriver(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthDriver = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDriver   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("driver.proto", fileDescriptorDriver) }

var fileDescriptorDriver = []byte{
	// 457 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xcd, 0x8e, 0x93, 0x50,
	0x14, 0xee, 0x9d, 0x56, 0x6c, 0xef, 0x30, 0x4a, 0xee, 0xc2, 0x10, 0x62, 0x10, 0x9b, 0x98, 0x34,
	0x9a, 0x61, 0x92, 0xba, 0x72, 0x89, 0x05, 0x93, 0x49, 0x2c, 0x6d, 0x28, 0xba, 0x70, 0x63, 0xf8,
	0x39, 0x30, 0x64, 0x28, 0x17, 0xb9, 0x97, 0xc6, 0x07, 0x70, 0xc5, 0x2b, 0x18, 0xa2, 0x8f, 0x33,
	0x4b, 0x1f, 0x41, 0xeb, 0x8b, 0x18, 0x2e, 0xad, 0xee, 0x9c, 0xd9, 0x7d, 0xdf, 0xfd, 0xce, 0x97,
	0xf3, 0x9d, 0x73, 0x0f, 0x96, 0xe3, 0x2a, 0xdb, 0x41, 0x65, 0x96, 0x15, 0xe5, 0x94, 0x3c, 0x49,
	0x69, 0x79, 0x9d, 0x9a, 0x59, 0x61, 0x86, 0x61, 0x9e, 0xb0, 0x2b, 0x93, 0xc5, 0xd7, 0xe6, 0x6e,
	0xde, 0xab, 0x11, 0xcd, 0xb5, 0xf3, 0x34, 0xe3, 0x57, 0x75, 0x68, 0x46, 0x74, 0x7b, 0x91, 0xd2,
	0x94, 0x5e, 0x08, 0x25, 0xac, 0x13, 0xc1, 0x04, 0x11, 0xa8, 0x77, 0x4c, 0xbf, 0x22, 0x2c, 0xaf,
	0x83, 0x8a, 0x81, 0x07, 0x9f, 0x6a, 0x60, 0x9c, 0xa8, 0xf8, 0x7e, 0x44, 0x0b, 0x0e, 0x05, 0x57,
	0x91, 0x81, 0x66, 0x13, 0xef, 0x48, 0x89, 0x86, 0xc7, 0x79, 0x50, 0xa4, 0x75, 0x90, 0x82, 0x7a,
	0x22, 0xa4, 0xbf, 0xbc, 0xd3, 0x92, 0x2c, 0x87, 0x22, 0xd8, 0x82, 0x3a, 0xec, 0xb5, 0x23, 0x27,
	0xaf, 0xf0, 0x68, 0x4b, 0x63, 0x50, 0x47, 0x06, 0x9a, 0x3d, 0x98, 0x3f, 0x33, 0x6f, 0x99, 0xc0,
	0x5c, 0xd2, 0x18, 0x3c, 0x61, 0x99, 0x7e, 0x41, 0xf8, 0xec, 0x90, 0x8e, 0x95, 0xb4, 0x60, 0x40,
	0x08, 0x1e, 0xd5, 0x01, 0xeb, 0xb3, 0xc9, 0x9e, 0xc0, 0xff, 0x0d, 0xb6, 0xc0, 0x12, 0x54, 0x15,
	0xad, 0x98, 0x3a, 0x34, 0x86, 0xb3, 0xd3, 0xf9, 0x8b, 0x5b, 0xdb, 0x8b, 0x7e, 0x4e, 0xe7, 0xf1,
	0x0e, 0xd6, 0xa9, 0x81, 0xf1, 0xbf, 0xd7, 0x2e, 0x02, 0x87, 0xcf, 0xc7, 0xf5, 0x08, 0xfc, 0xfc,
	0x1b, 0xc2, 0xa3, 0x2e, 0x37, 0x79, 0x8a, 0x65, 0xdb, 0x79, 0x63, 0xbd, 0x7b, 0xeb, 0x7f, 0x5c,
	0xae, 0x6c, 0x47, 0x19, 0x68, 0x0f, 0x9b, 0xd6, 0x38, 0xb5, 0x21, 0x09, 0xea, 0x9c, 0x8b, 0x92,
	0x47, 0x58, 0x72, 0x2d, 0xff, 0xf2, 0xbd, 0xa3, 0x20, 0x0d, 0x37, 0xad, 0x21, 0xb9, 0x01, 0xcf,
	0x76, 0x40, 0xa6, 0x58, 0x5e, 0x7b, 0xce, 0xda, 0x5b, 0x2d, 0x9c, 0xcd, 0xc6, 0xb1, 0x95, 0x13,
	0x4d, 0x69, 0x5a, 0x43, 0x5e, 0x57, 0x50, 0x56, 0x34, 0x02, 0xc6, 0x20, 0x26, 0x8f, 0xf1, 0xc4,
	0x72, 0xdd, 0x95, 0x6f, 0xf9, 0x8e, 0xad, 0x8c, 0xb4, 0xb3, 0xa6, 0x35, 0x26, 0x56, 0x51, 0x50,
	0x1e, 0x70, 0x88, 0xbb, 0x45, 0x6c, 0x9c, 0xa5, 0xe5, 0xfa, 0x97, 0x0b, 0x65, 0xac, 0xc9, 0x4d,
	0x6b, 0x8c, 0x37, 0xb0, 0x0d, 0x0a, 0x9e, 0x45, 0xf3, 0x12, 0x4b, 0xb6, 0x38, 0x24, 0x92, 0xe0,
	0x7b, 0x62, 0x1a, 0x72, 0x7e, 0xb7, 0x5d, 0x1c, 0x2e, 0x43, 0x33, 0xef, 0x5a, 0xde, 0x7f, 0xd5,
	0x6b, 0xfd, 0xe6, 0x97, 0x3e, 0xb8, 0xd9, 0xeb, 0xe8, 0xc7, 0x5e, 0x47, 0x3f, 0xf7, 0xfa, 0xe0,
	0xfb, 0x6f, 0x1d, 0x7d, 0x18, 0x1f, 0xab, 0x43, 0x49, 0xa0, 0x97, 0x7f, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x64, 0x59, 0x48, 0x92, 0xe1, 0x02, 0x00, 0x00,
}
