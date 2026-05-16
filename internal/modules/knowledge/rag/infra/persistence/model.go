package persistence

import (
	"time"

	"lattice-coding/internal/common/db"
)

type RAGSourcePO struct {
	db.BasePO
	SourceKey   string `gorm:"column:source_key;type:varchar(120);not null;uniqueIndex"`
	Name        string `gorm:"column:name;type:varchar(200);not null;index"`
	Type        string `gorm:"column:type;type:varchar(50);not null;index"`
	URI         string `gorm:"column:uri;type:varchar(1000);not null;default:''"`
	Description string `gorm:"column:description;type:text"`
	Owner       string `gorm:"column:owner;type:varchar(100);not null;default:'';index"`
	Status      string `gorm:"column:status;type:varchar(20);not null;default:'active';index"`
	Config      string `gorm:"column:config;type:json"`
	Metadata    string `gorm:"column:metadata;type:json"`
}

func (RAGSourcePO) TableName() string {
	return "rag_source"
}

type RAGDocumentPO struct {
	db.BasePO
	DocumentKey string `gorm:"column:document_key;type:varchar(120);not null;uniqueIndex"`
	SourceID    uint64 `gorm:"column:source_id;not null;index"`
	Title       string `gorm:"column:title;type:varchar(500);not null;index"`
	Type        string `gorm:"column:type;type:varchar(50);not null;default:'unknown';index"`
	URI         string `gorm:"column:uri;type:varchar(1000);not null;default:''"`
	Author      string `gorm:"column:author;type:varchar(200);not null;default:''"`
	Version     string `gorm:"column:version;type:varchar(100);not null;default:''"`
	Summary     string `gorm:"column:summary;type:text"`
	ContentHash string `gorm:"column:content_hash;type:varchar(128);not null;default:'';index"`
	Status      string `gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
	ParserName  string `gorm:"column:parser_name;type:varchar(80);not null;default:''"`
	ChunkCount  int    `gorm:"column:chunk_count;not null;default:0"`
	Metadata    string `gorm:"column:metadata;type:json"`
}

func (RAGDocumentPO) TableName() string {
	return "rag_document"
}

type RAGParsedDocumentPO struct {
	db.BasePO
	DocumentID  uint64 `gorm:"column:document_id;not null;index"`
	ParserName  string `gorm:"column:parser_name;type:varchar(80);not null;default:'';index"`
	ContentText string `gorm:"column:content_text;type:longtext"`
	Structure   string `gorm:"column:structure;type:json"`
	Sections    string `gorm:"column:sections;type:json"`
	Status      string `gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
	Error       string `gorm:"column:error;type:text"`
	Metadata    string `gorm:"column:metadata;type:json"`
}

func (RAGParsedDocumentPO) TableName() string {
	return "rag_parsed_document"
}

type RAGChunkPO struct {
	db.BasePO
	ChunkKey      string `gorm:"column:chunk_key;type:varchar(120);not null;uniqueIndex"`
	DocumentID    uint64 `gorm:"column:document_id;not null;index"`
	SourceID      uint64 `gorm:"column:source_id;not null;index"`
	Ordinal       int    `gorm:"column:ordinal;not null;default:0;index"`
	Title         string `gorm:"column:title;type:varchar(500);not null;default:''"`
	Content       string `gorm:"column:content;type:longtext"`
	ContentHash   string `gorm:"column:content_hash;type:varchar(128);not null;default:'';index"`
	Location      string `gorm:"column:location;type:varchar(500);not null;default:''"`
	TokenEstimate int    `gorm:"column:token_estimate;not null;default:0"`
	Strategy      string `gorm:"column:strategy;type:varchar(50);not null;default:'';index"`
	Metadata      string `gorm:"column:metadata;type:json"`
}

func (RAGChunkPO) TableName() string {
	return "rag_chunk"
}

type RAGEmbeddingPO struct {
	db.BasePO
	ChunkID    uint64 `gorm:"column:chunk_id;not null;index"`
	Model      string `gorm:"column:model;type:varchar(120);not null;index"`
	Dimension  int    `gorm:"column:dimension;not null;default:0"`
	VectorRef  string `gorm:"column:vector_ref;type:varchar(200);not null;default:'';index"`
	VectorHash string `gorm:"column:vector_hash;type:varchar(128);not null;default:'';index"`
	Status     string `gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
	Error      string `gorm:"column:error;type:text"`
	Metadata   string `gorm:"column:metadata;type:json"`
}

func (RAGEmbeddingPO) TableName() string {
	return "rag_embedding"
}

type RAGKeywordIndexPO struct {
	db.BasePO
	ChunkID    uint64     `gorm:"column:chunk_id;not null;index"`
	IndexName  string     `gorm:"column:index_name;type:varchar(120);not null;index"`
	DocumentID uint64     `gorm:"column:document_id;not null;index"`
	SourceID   uint64     `gorm:"column:source_id;not null;index"`
	Status     string     `gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
	IndexedAt  *time.Time `gorm:"column:indexed_at"`
	Error      string     `gorm:"column:error;type:text"`
	Metadata   string     `gorm:"column:metadata;type:json"`
}

func (RAGKeywordIndexPO) TableName() string {
	return "rag_keyword_index"
}

type RAGRetrievalTracePO struct {
	db.BasePO
	Query        string `gorm:"column:query;type:text;not null"`
	Route        string `gorm:"column:route;type:varchar(50);not null;index"`
	VectorCount  int    `gorm:"column:vector_count;not null;default:0"`
	KeywordCount int    `gorm:"column:keyword_count;not null;default:0"`
	HybridCount  int    `gorm:"column:hybrid_count;not null;default:0"`
	LatencyMs    int64  `gorm:"column:latency_ms;not null;default:0"`
	Metadata     string `gorm:"column:metadata;type:json"`
}

func (RAGRetrievalTracePO) TableName() string {
	return "rag_retrieval_trace"
}

type RAGEvidencePO struct {
	db.BasePO
	EvidenceKey   string  `gorm:"column:evidence_key;type:varchar(120);not null;index"`
	Query         string  `gorm:"column:query;type:text;not null"`
	ChunkID       uint64  `gorm:"column:chunk_id;not null;index"`
	DocumentID    uint64  `gorm:"column:document_id;not null;index"`
	SourceID      uint64  `gorm:"column:source_id;not null;index"`
	Channel       string  `gorm:"column:channel;type:varchar(50);not null;index"`
	Title         string  `gorm:"column:title;type:varchar(500);not null;default:''"`
	Content       string  `gorm:"column:content;type:longtext"`
	Location      string  `gorm:"column:location;type:varchar(500);not null;default:''"`
	Score         float64 `gorm:"column:score;not null;default:0;index"`
	RerankScore   float64 `gorm:"column:rerank_score;not null;default:0;index"`
	TokenEstimate int     `gorm:"column:token_estimate;not null;default:0"`
	CitationURI   string  `gorm:"column:citation_uri;type:varchar(1000);not null;default:''"`
	Metadata      string  `gorm:"column:metadata;type:json"`
}

func (RAGEvidencePO) TableName() string {
	return "rag_evidence"
}

type RAGRerankResultPO struct {
	db.BasePO
	Query       string `gorm:"column:query;type:text;not null"`
	Model       string `gorm:"column:model;type:varchar(120);not null;default:'';index"`
	InputCount  int    `gorm:"column:input_count;not null;default:0"`
	OutputCount int    `gorm:"column:output_count;not null;default:0"`
	Items       string `gorm:"column:items;type:json"`
	Metadata    string `gorm:"column:metadata;type:json"`
}

func (RAGRerankResultPO) TableName() string {
	return "rag_rerank_result"
}

type RAGEvidenceGateResultPO struct {
	db.BasePO
	Query           string  `gorm:"column:query;type:text;not null"`
	Passed          bool    `gorm:"column:passed;not null;default:false;index"`
	Score           float64 `gorm:"column:score;not null;default:0;index"`
	Reason          string  `gorm:"column:reason;type:text"`
	MissingAspects  string  `gorm:"column:missing_aspects;type:json"`
	WeakEvidenceIDs string  `gorm:"column:weak_evidence_ids;type:json"`
	Contradictions  string  `gorm:"column:contradictions;type:json"`
	Suggestions     string  `gorm:"column:suggestions;type:json"`
	Status          string  `gorm:"column:status;type:varchar(20);not null;default:'review';index"`
	Metadata        string  `gorm:"column:metadata;type:json"`
}

func (RAGEvidenceGateResultPO) TableName() string {
	return "rag_evidence_gate_result"
}
