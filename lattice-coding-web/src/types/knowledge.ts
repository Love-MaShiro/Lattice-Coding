export interface KnowledgeDocument {
  id: number
  title: string
  content: string
  tags: string[]
  createdAt: string
  updatedAt: string
}

export interface KnowledgeForm {
  title: string
  content: string
  tags: string[]
}
