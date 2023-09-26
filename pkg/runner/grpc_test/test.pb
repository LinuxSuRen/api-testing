
Ã

pkg/runner/grpc_test/test.protogrpctest"
Empty"&

HelloReply
message (	Rmessage"C
StreamMessage
MsgID (RMsgID
	ExpectLen (R	ExpectLen"D
StreamMessageRepeated+
data (2.grpctest.StreamMessageRdata"«
	BasicType
Int32 (RInt32
Int64 (RInt64
Uint32 (RUint32
Uint64 (RUint64
Float32 (RFloat32
Float64 (RFloat64
String (	RString
Bool (RBool"√
AdvancedType

Int32Array (R
Int32Array

Int64Array (R
Int64Array 
Uint32Array (RUint32Array 
Uint64Array (RUint64Array"
Float32Array (RFloat32Array"
Float64Array (RFloat64Array 
StringArray (	RStringArray
	BoolArray (R	BoolArrayO
HelloReplyMap	 (2).grpctest.AdvancedType.HelloReplyMapEntryRHelloReplyMapV
HelloReplyMapEntry
key (	Rkey*
value (2.grpctest.HelloReplyRvalue:82ê
Main.
Unary.grpctest.Empty.grpctest.HelloReplyJ
ClientStream.grpctest.StreamMessage.grpctest.StreamMessageRepeated(J
ServerStream.grpctest.StreamMessageRepeated.grpctest.StreamMessage0A
	BidStream.grpctest.StreamMessage.grpctest.StreamMessage(09
TestBasicType.grpctest.BasicType.grpctest.BasicTypeB
TestAdvancedType.grpctest.AdvancedType.grpctest.AdvancedTypeB8Z6github.com/linuxsuren/api-testing/pkg/runner/grpc_testbproto3