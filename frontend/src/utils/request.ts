import type { RequestOptions } from '@@/plugin-request/request';
import type { RequestConfig } from '@umijs/max';
import { message, notification } from 'antd';
import { getToken, setToken, getRefreshToken, removeToken } from './auth';
import { refreshToken } from '@/services/backend/user';

// 错误处理方案： 错误类型
enum ErrorShowType {
  SILENT = 0,
  WARN_MESSAGE = 1,
  ERROR_MESSAGE = 2,
  NOTIFICATION = 3,
  REDIRECT = 9,
}

// 与后端约定的响应数据格式
interface IResponse<T = any> {
  success: boolean;
  data: T;
  errorCode?: number;
  errorMessage?: string;
  errorShowType?: ErrorShowType;
}

// 请求队列项类型
interface RequestQueueItem {
  resolve: (value?: string | null) => void;
  reject: (reason?: any) => void;
}

// Token 刷新状态管理
let isRefreshing = false;
let failedQueue: RequestQueueItem[] = [];

/**
 * 处理 Token 刷新
 * @returns 新的 access token
 */
const handleTokenRefresh = async (): Promise<string> => {
  try {
    const refreshTokenValue = getRefreshToken();
    if (!refreshTokenValue) {
      throw new Error('No refresh token available');
    }

    const response = await refreshToken({ refreshToken: refreshTokenValue });
    if (response.success) {
      setToken(response.data.accessToken, response.data.refreshToken);
      return response.data.accessToken;
    }
    throw new Error(response.errorMessage || '刷新 token 失败');
  } catch (error) {
    removeToken();
    window.location.href = '/user/login';
    throw error;
  }
};

/**
 * 处理请求队列
 * @param error 错误对象
 * @param token 新的 token 或 null
 */
const processQueue = (error: any, token: string | null = null): void => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

/**
 * 处理业务错误
 * @param errorInfo 错误信息
 */
const handleBusinessError = (errorInfo: IResponse): void => {
  const { errorCode, errorMessage, errorShowType } = errorInfo;

  // 按错误码处理
  const errorCodeHandlers: Record<number, () => void> = {
    400: () => notification.error({ message: '请求错误', description: errorMessage }),
    401: () => {}, // 401 在拦截器中特殊处理
    403: () => notification.error({ message: '权限不足', description: '请联系管理员' }),
    404: () => notification.error({ message: '资源不存在', description: '请联系管理员' }),
    500: () => notification.error({ message: '服务器错误', description: '请联系管理员' }),
  };

  if (errorCode && errorCodeHandlers[errorCode]) {
    errorCodeHandlers[errorCode]();
    return;
  }

  // 按错误展示类型处理
  switch (errorShowType) {
    case ErrorShowType.SILENT:
      break;
    case ErrorShowType.WARN_MESSAGE:
      message.warning(errorMessage);
      break;
    case ErrorShowType.ERROR_MESSAGE:
      message.error(errorMessage);
      break;
    case ErrorShowType.NOTIFICATION:
      notification.open({ description: errorMessage, message: errorCode?.toString() });
      break;
    case ErrorShowType.REDIRECT:
      // TODO: 实现重定向逻辑
      break;
    default:
      message.error(errorMessage || '未知错误');
  }
};

/**
 * @name 错误处理
 * pro 自带的错误处理， 可以在这里做自己的改动
 * @doc https://umijs.org/docs/max/request#配置
 */
/**
 * 错误配置
 */
export const errorConfig: RequestConfig = {
  // 错误处理：umi@3 的错误处理方案
  errorConfig: {
    // 错误抛出
    errorThrower: (response: unknown) => {
      const { success, errorCode, errorMessage, errorShowType, data } = response as IResponse;
      if (!success) {
        const error = new Error(errorMessage);
        error.name = 'BizError';
        error.info = { errorCode, errorMessage, errorShowType, data };
        throw error;
      }
    },

    // 错误接收及处理
    errorHandler: async (error: any, opts: any) => {
      if (opts?.skipErrorHandler) throw error;

      // 业务错误
      if (error.name === 'BizError') {
        const errorInfo = error.info as IResponse;
        if (errorInfo) {
          // 401 特殊处理 - token 刷新逻辑
          if (errorInfo.errorCode === 401) {
            const originalRequest = error.config;

            if (!isRefreshing) {
              isRefreshing = true;
              try {
                const newToken = await handleTokenRefresh();
                processQueue(null, newToken);

                // 重试原始请求
                originalRequest.headers.Authorization = `Bearer ${newToken}`;
                return fetch(originalRequest);
              } catch (err) {
                processQueue(err, null);
                return Promise.reject(err);
              } finally {
                isRefreshing = false;
              }
            }

            // 如果正在刷新 token，将请求加入队列
            return new Promise((resolve, reject) => {
              failedQueue.push({ resolve, reject });
            });
          }

          handleBusinessError(errorInfo);
        }
        return;
      }

      if (error.response) {
        // 请求成功发出且服务器响应了状态码，但状态码超出 2xx 范围
        message.error(`请求错误，状态码: ${error.response.status}`);
      } else if (error.request) {
        // 请求已发出但没有收到响应
        message.error('请求无响应，请重试');
      } else {
        // 请求设置出错
        message.error('请求错误，请重试');
      }
    },
  },

  // 请求拦截器
  requestInterceptors: [
    (config: RequestOptions) => {
      // 在浏览器环境中添加 token
      if (typeof window !== 'undefined') {
        const token = getToken();
        if (token) {
          config.headers = {
            ...config.headers,
            Authorization: `Bearer ${token}`,
          };
        }
      }
      return config;
    },
  ],
};
