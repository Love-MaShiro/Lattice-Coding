package persistence

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&RAGSourcePO{},
		&RAGDocumentPO{},
		&RAGParsedDocumentPO{},
		&RAGChunkPO{},
		&RAGEmbeddingPO{},
		&RAGKeywordIndexPO{},
		&RAGRetrievalTracePO{},
		&RAGEvidencePO{},
		&RAGRerankResultPO{},
		&RAGEvidenceGateResultPO{},
	)
}
