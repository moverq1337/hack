import grpc
import warnings

import nlp_pb2 as nlp__pb2

GRPC_GENERATED_VERSION = '1.64.1'
GRPC_VERSION = grpc.__version__
EXPECTED_ERROR_RELEASE = '1.65.0'
SCHEDULED_RELEASE_DATE = 'June 25, 2024'
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    warnings.warn(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in nlp_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
        + f' This warning will become an error in {EXPECTED_ERROR_RELEASE},'
        + f' scheduled for release on {SCHEDULED_RELEASE_DATE}.',
        RuntimeWarning
    )

class NLPServiceStub(object):
    
    def __init__(self, channel):
        
        self.ParseResume = channel.unary_unary(
                '/pb.NLPService/ParseResume',
                request_serializer=nlp__pb2.ParseRequest.SerializeToString,
                response_deserializer=nlp__pb2.ParseResponse.FromString,
                _registered_method=True)
        self.MatchResumeVacancy = channel.unary_unary(
                '/pb.NLPService/MatchResumeVacancy',
                request_serializer=nlp__pb2.MatchRequest.SerializeToString,
                response_deserializer=nlp__pb2.MatchResponse.FromString,
                _registered_method=True)

class NLPServiceServicer(object):
    
    def ParseResume(self, request, context):
        
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def MatchResumeVacancy(self, request, context):
        
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

def add_NLPServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'ParseResume': grpc.unary_unary_rpc_method_handler(
                    servicer.ParseResume,
                    request_deserializer=nlp__pb2.ParseRequest.FromString,
                    response_serializer=nlp__pb2.ParseResponse.SerializeToString,
            ),
            'MatchResumeVacancy': grpc.unary_unary_rpc_method_handler(
                    servicer.MatchResumeVacancy,
                    request_deserializer=nlp__pb2.MatchRequest.FromString,
                    response_serializer=nlp__pb2.MatchResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'pb.NLPService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('pb.NLPService', rpc_method_handlers)

class NLPService(object):
    
    @staticmethod
    def ParseResume(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/pb.NLPService/ParseResume',
            nlp__pb2.ParseRequest.SerializeToString,
            nlp__pb2.ParseResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def MatchResumeVacancy(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/pb.NLPService/MatchResumeVacancy',
            nlp__pb2.MatchRequest.SerializeToString,
            nlp__pb2.MatchResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
