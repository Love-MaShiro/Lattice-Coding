package persistence

import packdomain "lattice-coding/internal/modules/knowledge/pack/domain"

func packToPO(pack *packdomain.KnowledgePack, po *KnowledgePackPO) {
	po.ID = pack.ID
	po.PackKey = pack.PackKey
	po.Query = pack.Query
	po.Intent = pack.Intent
	po.Route = pack.Route
	po.Status = string(pack.Status)
	po.TokenEstimate = pack.TokenEstimate
	po.MaxTokens = pack.MaxTokens
	po.PromptContext = pack.PromptContext
	po.Warnings = pack.Warnings
	po.Options = pack.Options
	po.Meta = pack.Meta
}

func poToPack(po *KnowledgePackPO) *packdomain.KnowledgePack {
	return &packdomain.KnowledgePack{
		ID:            po.ID,
		PackKey:       po.PackKey,
		Query:         po.Query,
		Intent:        po.Intent,
		Route:         po.Route,
		Status:        packdomain.PackStatus(po.Status),
		TokenEstimate: po.TokenEstimate,
		MaxTokens:     po.MaxTokens,
		PromptContext: po.PromptContext,
		Warnings:      po.Warnings,
		Options:       po.Options,
		Meta:          po.Meta,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
	}
}

func itemToPO(packID uint64, item *packdomain.KnowledgeItem, po *KnowledgeItemPO) {
	po.ID = item.ID
	po.PackID = packID
	po.ItemKey = item.ItemKey
	po.SourceKind = string(item.SourceKind)
	po.SourceID = item.SourceID
	po.SourceType = item.SourceType
	po.Title = item.Title
	po.Content = item.Content
	po.Location = item.Location
	po.Score = item.Score
	po.TokenEstimate = item.TokenEstimate
	po.CitationKey = item.CitationKey
	po.Metadata = item.Metadata
	po.SortOrder = item.SortOrder
}

func poToItem(po *KnowledgeItemPO) packdomain.KnowledgeItem {
	return packdomain.KnowledgeItem{
		ID:            po.ID,
		PackID:        po.PackID,
		ItemKey:       po.ItemKey,
		SourceKind:    packdomain.PackSourceKind(po.SourceKind),
		SourceID:      po.SourceID,
		SourceType:    po.SourceType,
		Title:         po.Title,
		Content:       po.Content,
		Location:      po.Location,
		Score:         po.Score,
		TokenEstimate: po.TokenEstimate,
		CitationKey:   po.CitationKey,
		Metadata:      po.Metadata,
		SortOrder:     po.SortOrder,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
	}
}

func citationToPO(packID uint64, citation *packdomain.KnowledgeCitation, po *KnowledgeCitationPO) {
	po.ID = citation.ID
	po.PackID = packID
	po.CitationKey = citation.CitationKey
	po.SourceKind = string(citation.SourceKind)
	po.SourceID = citation.SourceID
	po.Title = citation.Title
	po.Location = citation.Location
	po.URI = citation.URI
	po.Score = citation.Score
	po.Metadata = citation.Metadata
	po.SortOrder = citation.SortOrder
}

func poToCitation(po *KnowledgeCitationPO) packdomain.KnowledgeCitation {
	return packdomain.KnowledgeCitation{
		ID:          po.ID,
		PackID:      po.PackID,
		CitationKey: po.CitationKey,
		SourceKind:  packdomain.PackSourceKind(po.SourceKind),
		SourceID:    po.SourceID,
		Title:       po.Title,
		Location:    po.Location,
		URI:         po.URI,
		Score:       po.Score,
		Metadata:    po.Metadata,
		SortOrder:   po.SortOrder,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}
