# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: api.proto
# Protobuf Python Version: 4.25.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\tapi.proto\x12\x03\x61pi\"-\n\x1cTopLongMemesHighlightRequest\x12\r\n\x05limit\x18\x01 \x01(\x05\"\'\n\x16SearchHighlightRequest\x12\r\n\x05query\x18\x01 \x01(\t\"&\n\x15MonthHighlightRequest\x12\r\n\x05month\x18\x01 \x01(\x05\"\x17\n\x15\x45mptyHighlightRequest\"!\n\x11HighlightResponse\x12\x0c\n\x04text\x18\x01 \x01(\t2\xbb\x02\n\x10RequesterService\x12L\n\x0fGetTopLongMemes\x12!.api.TopLongMemesHighlightRequest\x1a\x16.api.HighlightResponse\x12M\n\x16SearchMemesBySubstring\x12\x1b.api.SearchHighlightRequest\x1a\x16.api.HighlightResponse\x12\x45\n\x0fGetMemesByMonth\x12\x1a.api.MonthHighlightRequest\x1a\x16.api.HighlightResponse\x12\x43\n\rGetRandomMeme\x12\x1a.api.EmptyHighlightRequest\x1a\x16.api.HighlightResponseB\x0cZ\nrequester/b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'api_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  _globals['DESCRIPTOR']._options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\nrequester/'
  _globals['_TOPLONGMEMESHIGHLIGHTREQUEST']._serialized_start=18
  _globals['_TOPLONGMEMESHIGHLIGHTREQUEST']._serialized_end=63
  _globals['_SEARCHHIGHLIGHTREQUEST']._serialized_start=65
  _globals['_SEARCHHIGHLIGHTREQUEST']._serialized_end=104
  _globals['_MONTHHIGHLIGHTREQUEST']._serialized_start=106
  _globals['_MONTHHIGHLIGHTREQUEST']._serialized_end=144
  _globals['_EMPTYHIGHLIGHTREQUEST']._serialized_start=146
  _globals['_EMPTYHIGHLIGHTREQUEST']._serialized_end=169
  _globals['_HIGHLIGHTRESPONSE']._serialized_start=171
  _globals['_HIGHLIGHTRESPONSE']._serialized_end=204
  _globals['_REQUESTERSERVICE']._serialized_start=207
  _globals['_REQUESTERSERVICE']._serialized_end=522
# @@protoc_insertion_point(module_scope)
