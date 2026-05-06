package application

type ProviderPageQuery struct {
	Page     int
	PageSize int
	Keyword  string
}

type ModelConfigPageQuery struct {
	Page       int
	PageSize   int
	ProviderID uint64
	Keyword    string
}
