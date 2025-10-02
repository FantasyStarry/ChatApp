// API客户端配置
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// 通用请求函数
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
}

// 用户相关API
export const userApi = {
  // 用户登录
  login: (credentials: { username: string; password: string }) =>
    apiRequest<{ token: string; user: unknown }>('/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    }),

  // 获取当前用户信息
  getCurrentUser: () =>
    apiRequest<unknown>('/profile'),

  // 用户登出
  logout: () =>
    apiRequest<unknown>('/logout', {
      method: 'POST',
    }),
};

// 聊天室相关API
export const chatroomApi = {
  // 获取聊天室列表
  getChatrooms: () =>
    apiRequest<unknown[]>('/chatrooms'),

  // 创建聊天室
  createChatroom: (chatroomData: { name: string; description?: string }) =>
    apiRequest<unknown>('/chatrooms', {
      method: 'POST',
      body: JSON.stringify(chatroomData),
    }),

  // 获取聊天室详情
  getChatroom: (id: string) =>
    apiRequest<unknown>(`/chatrooms/${id}`),
};

// 消息相关API
export const messageApi = {
  // 获取聊天室消息
  getMessages: (chatroomId: string, params?: { limit?: number; offset?: number }) => {
    const queryParams = new URLSearchParams();
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());
    
    const queryString = queryParams.toString();
    return apiRequest<unknown[]>(`/chatrooms/${chatroomId}/messages${queryString ? `?${queryString}` : ''}`);
  },
};

// WebSocket连接管理
export class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectInterval = 3000;
  private currentChatroomId: string | null = null;

  constructor(
    private onMessage: (data: unknown) => void,
    private onError?: (error: Event) => void,
    private onClose?: (event: CloseEvent) => void
  ) {}

  connect(chatroomId: string, token?: string) {
    const wsBaseUrl = process.env.NEXT_PUBLIC_WS_BASE_URL || 'ws://localhost:8080/ws';
    const url = token ? `${wsBaseUrl}/${chatroomId}?token=${token}` : `${wsBaseUrl}/${chatroomId}`;
    
    this.currentChatroomId = chatroomId;

    try {
      this.ws = new WebSocket(url);

      this.ws.onopen = () => {
        console.log(`WebSocket connected to chatroom ${chatroomId}`);
        this.reconnectAttempts = 0;
      };

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          this.onMessage(data);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.onError?.(error);
      };

      this.ws.onclose = (event) => {
        console.log('WebSocket closed:', event.code, event.reason);
        this.onClose?.(event);
        this.attemptReconnect();
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts && this.currentChatroomId) {
      this.reconnectAttempts++;
      console.log(`Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
      
      setTimeout(() => {
        this.connect(this.currentChatroomId!);
      }, this.reconnectInterval);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }

  sendMessage(data: unknown) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  getConnectionState() {
    return this.ws?.readyState;
  }
}
