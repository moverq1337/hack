from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder

_sym_db = _symbol_database.Default()

DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\tnlp.proto\x12\x02pb\"\x1c\n\x0cParseRequest\x12\x0c\n\x04text\x18\x01 \x01(\t\"$\n\rParseResponse\x12\x13\n\x0bparsed_data\x18\x01 \x01(\t\"9\n\x0cMatchRequest\x12\x13\n\x0bresume_text\x18\x01 \x01(\t\x12\x14\n\x0cvacancy_text\x18\x02 \x01(\t\"\x1e\n\rMatchResponse\x12\r\n\x05score\x18\x01 \x01(\x02\x32{\n\nNLPService\x12\x32\n\x0bParseResume\x12\x10.pb.ParseRequest\x1a\x11.pb.ParseResponse\x12\x39\n\x12MatchResumeVacancy\x12\x10.pb.MatchRequest\x1a\x11.pb.MatchResponseB+Z)github.com/moverq1337/VTBHack/internal/pbb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'nlp_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z)github.com/moverq1337/VTBHack/internal/pb'
  _globals['_PARSEREQUEST']._serialized_start=17
  _globals['_PARSEREQUEST']._serialized_end=45
  _globals['_PARSERESPONSE']._serialized_start=47
  _globals['_PARSERESPONSE']._serialized_end=83
  _globals['_MATCHREQUEST']._serialized_start=85
  _globals['_MATCHREQUEST']._serialized_end=142
  _globals['_MATCHRESPONSE']._serialized_start=144
  _globals['_MATCHRESPONSE']._serialized_end=174
  _globals['_NLPSERVICE']._serialized_start=176
  _globals['_NLPSERVICE']._serialized_end=299
