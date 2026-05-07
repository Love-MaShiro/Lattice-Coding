export interface ApiResponse<T = any> {
  code: number | string
  message: string
  data: T
}

export interface PageResult<T = any> {
  data: T[]
  total: number
  page: number
  size: number
}

export interface PageQuery {
  page?: number
  size?: number
  keyword?: string
}
