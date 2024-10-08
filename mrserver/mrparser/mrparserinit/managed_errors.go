package mrparserinit

import (
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// ManagedHttpErrors - comment func.
func ManagedHttpErrors() []mrinit.EnrichedError {
	return []mrinit.EnrichedError{
		mrinit.WrapProto(mrparser.ErrHttpRequestFileSizeMin),
		mrinit.WrapProto(mrparser.ErrHttpRequestFileSizeMax),
		mrinit.WrapProto(mrparser.ErrHttpRequestFileExtension),
		mrinit.WrapProto(mrparser.ErrHttpRequestFileTotalSizeMax),
		mrinit.WrapProto(mrparser.ErrHttpRequestFileContentType),
		mrinit.WrapProto(mrparser.ErrHttpRequestFileUnsupportedType),
		mrinit.WrapProto(mrparser.ErrHttpRequestImageWidthMax),
		mrinit.WrapProto(mrparser.ErrHttpRequestImageHeightMax),
	}
}
