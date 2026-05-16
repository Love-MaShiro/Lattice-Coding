package domain

type SourceType string

const (
	SourceTypeFile     SourceType = "file"
	SourceTypeWeb      SourceType = "web"
	SourceTypeCodeRepo SourceType = "code_repo"
	SourceTypeDatabase SourceType = "database"
	SourceTypeManual   SourceType = "manual"
	SourceTypeDataset  SourceType = "dataset"
)

type SourceStatus string

const (
	SourceStatusActive   SourceStatus = "active"
	SourceStatusDisabled SourceStatus = "disabled"
)

type DocumentType string

const (
	DocumentTypeMarkdown DocumentType = "markdown"
	DocumentTypeHTML     DocumentType = "html"
	DocumentTypePDF      DocumentType = "pdf"
	DocumentTypeText     DocumentType = "text"
	DocumentTypeCode     DocumentType = "code"
	DocumentTypeUnknown  DocumentType = "unknown"
)

type DocumentStatus string

const (
	DocumentStatusPending DocumentStatus = "pending"
	DocumentStatusParsed  DocumentStatus = "parsed"
	DocumentStatusIndexed DocumentStatus = "indexed"
	DocumentStatusFailed  DocumentStatus = "failed"
)

type ParseStatus string

const (
	ParseStatusPending ParseStatus = "pending"
	ParseStatusParsed  ParseStatus = "parsed"
	ParseStatusFailed  ParseStatus = "failed"
)

type ChunkStrategy string

const (
	ChunkStrategyToken   ChunkStrategy = "token"
	ChunkStrategySection ChunkStrategy = "section"
	ChunkStrategyCode    ChunkStrategy = "code"
)

type EmbeddingStatus string

const (
	EmbeddingStatusPending EmbeddingStatus = "pending"
	EmbeddingStatusReady   EmbeddingStatus = "ready"
	EmbeddingStatusFailed  EmbeddingStatus = "failed"
)

type KeywordIndexStatus string

const (
	KeywordIndexStatusPending KeywordIndexStatus = "pending"
	KeywordIndexStatusReady   KeywordIndexStatus = "ready"
	KeywordIndexStatusFailed  KeywordIndexStatus = "failed"
)

type RetrievalRoute string

const (
	RetrievalRouteVector  RetrievalRoute = "vector"
	RetrievalRouteKeyword RetrievalRoute = "keyword"
	RetrievalRouteHybrid  RetrievalRoute = "hybrid"
)

type RetrievalChannel string

const (
	RetrievalChannelVector  RetrievalChannel = "vector"
	RetrievalChannelKeyword RetrievalChannel = "keyword"
)

type EvidenceGateStatus string

const (
	EvidenceGateStatusPassed EvidenceGateStatus = "passed"
	EvidenceGateStatusFailed EvidenceGateStatus = "failed"
	EvidenceGateStatusReview EvidenceGateStatus = "review"
)
