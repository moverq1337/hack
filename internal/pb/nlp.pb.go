package pb

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ParseRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Text          string                 `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ParseRequest) Reset() {
	*x = ParseRequest{}
	mi := &file_proto_nlp_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ParseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseRequest) ProtoMessage() {}

func (x *ParseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_nlp_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ParseRequest) Descriptor() ([]byte, []int) {
	return file_proto_nlp_proto_rawDescGZIP(), []int{0}
}

func (x *ParseRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type ParseResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ParsedData    string                 `protobuf:"bytes,1,opt,name=parsed_data,json=parsedData,proto3" json:"parsed_data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ParseResponse) Reset() {
	*x = ParseResponse{}
	mi := &file_proto_nlp_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ParseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseResponse) ProtoMessage() {}

func (x *ParseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_nlp_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ParseResponse) Descriptor() ([]byte, []int) {
	return file_proto_nlp_proto_rawDescGZIP(), []int{1}
}

func (x *ParseResponse) GetParsedData() string {
	if x != nil {
		return x.ParsedData
	}
	return ""
}

type MatchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ResumeText    string                 `protobuf:"bytes,1,opt,name=resume_text,json=resumeText,proto3" json:"resume_text,omitempty"`
	VacancyText   string                 `protobuf:"bytes,2,opt,name=vacancy_text,json=vacancyText,proto3" json:"vacancy_text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MatchRequest) Reset() {
	*x = MatchRequest{}
	mi := &file_proto_nlp_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MatchRequest) ProtoMessage() {}

func (x *MatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_nlp_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MatchRequest) Descriptor() ([]byte, []int) {
	return file_proto_nlp_proto_rawDescGZIP(), []int{2}
}

func (x *MatchRequest) GetResumeText() string {
	if x != nil {
		return x.ResumeText
	}
	return ""
}

func (x *MatchRequest) GetVacancyText() string {
	if x != nil {
		return x.VacancyText
	}
	return ""
}

type MatchResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Score         float32                `protobuf:"fixed32,1,opt,name=score,proto3" json:"score,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MatchResponse) Reset() {
	*x = MatchResponse{}
	mi := &file_proto_nlp_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MatchResponse) ProtoMessage() {}

func (x *MatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_nlp_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MatchResponse) Descriptor() ([]byte, []int) {
	return file_proto_nlp_proto_rawDescGZIP(), []int{3}
}

func (x *MatchResponse) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

var File_proto_nlp_proto protoreflect.FileDescriptor

const file_proto_nlp_proto_rawDesc = "" +
	"\n" +
	"\x0fproto/nlp.proto\x12\x02pb\"\"\n" +
	"\fParseRequest\x12\x12\n" +
	"\x04text\x18\x01 \x01(\tR\x04text\"0\n" +
	"\rParseResponse\x12\x1f\n" +
	"\vparsed_data\x18\x01 \x01(\tR\n" +
	"parsedData\"R\n" +
	"\fMatchRequest\x12\x1f\n" +
	"\vresume_text\x18\x01 \x01(\tR\n" +
	"resumeText\x12!\n" +
	"\fvacancy_text\x18\x02 \x01(\tR\vvacancyText\"%\n" +
	"\rMatchResponse\x12\x14\n" +
	"\x05score\x18\x01 \x01(\x02R\x05score2\x7f\n" +
	"\n" +
	"NLPService\x124\n" +
	"\vParseResume\x12\x10.pb.ParseRequest\x1a\x11.pb.ParseResponse\"\x00\x12;\n" +
	"\x12MatchResumeVacancy\x12\x10.pb.MatchRequest\x1a\x11.pb.MatchResponse\"\x00B+Z)github.com/moverq1337/VTBHack/internal/pbb\x06proto3"

var (
	file_proto_nlp_proto_rawDescOnce sync.Once
	file_proto_nlp_proto_rawDescData []byte
)

func file_proto_nlp_proto_rawDescGZIP() []byte {
	file_proto_nlp_proto_rawDescOnce.Do(func() {
		file_proto_nlp_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_nlp_proto_rawDesc), len(file_proto_nlp_proto_rawDesc)))
	})
	return file_proto_nlp_proto_rawDescData
}

var file_proto_nlp_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_nlp_proto_goTypes = []any{
	(*ParseRequest)(nil),  
	(*ParseResponse)(nil), 
	(*MatchRequest)(nil),  
	(*MatchResponse)(nil), 
}
var file_proto_nlp_proto_depIdxs = []int32{
	0, 
	2, 
	1, 
	3, 
	2, 
	0, 
	0, 
	0, 
	0, 
}

func init() { file_proto_nlp_proto_init() }
func file_proto_nlp_proto_init() {
	if File_proto_nlp_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_nlp_proto_rawDesc), len(file_proto_nlp_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_nlp_proto_goTypes,
		DependencyIndexes: file_proto_nlp_proto_depIdxs,
		MessageInfos:      file_proto_nlp_proto_msgTypes,
	}.Build()
	File_proto_nlp_proto = out.File
	file_proto_nlp_proto_goTypes = nil
	file_proto_nlp_proto_depIdxs = nil
}
