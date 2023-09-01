package logic

import (
	"github.com/ShyunnY/cruise/pkg/reader"
	"github.com/ShyunnY/cruise/pkg/storage"
)

// ServiceCtx
// service context include store,reader components for service need
type ServiceCtx struct {
	Store  storage.Storage
	Reader reader.Reader
}
