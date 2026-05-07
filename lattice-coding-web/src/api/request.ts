import axios, {
  AxiosError,
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  InternalAxiosRequestConfig
} from 'axios'
import { ElMessage } from 'element-plus'
import type { ApiResponse } from '@/types/api'

const request: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 60000
})

request.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError): Promise<never> => Promise.reject(error)
)

request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>): any => {
    const res = response.data
    const code = String(res?.code ?? '').trim()

    if (!['0', '200', '10000'].includes(code)) {
      ElMessage.error(res?.message || '请求失败')
      return Promise.reject(response)
    }

    return res.data
  },
  (error: AxiosError): Promise<never> => {
    if (error.response) {
      const status = error.response.status
      let message = '请求失败'

      if (status === 401) {
        message = '未授权，请重新登录'
      } else if (status === 403) {
        message = '拒绝访问'
      } else if (status === 404) {
        message = '请求资源不存在'
      } else if (status >= 500) {
        message = '服务器错误'
      } else if (error.code === 'ECONNABORTED') {
        message = '请求超时'
      } else if (error.message) {
        message = error.message
      }

      ElMessage.error(message)
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }

    return Promise.reject(error)
  }
)

export function get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return request.get(url, config)
}

export function post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return request.post(url, data, config)
}

export function put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return request.put(url, data, config)
}

export function del<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return request.delete(url, config)
}

export default request
