import type { Device, ListMessagesResponse, LoginRequest, LoginResponse, Me, MessageDetail, SendMessageRequest, Stats } from './types';

const BASE = '/api/v1';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  if (!res.ok) {
    throw new Error(`API error: ${res.status}`);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return res.json();
}

export function login(data: LoginRequest) {
  return request<LoginResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function logout() {
  return request<void>('/auth/logout', { method: 'POST' });
}

export function me() {
  return request<Me>('/auth/me');
}

export function stats() {
  return request<Stats>('/stats');
}

export function listMessages(params?: string) {
  return request<ListMessagesResponse>(`/messages${params ? '?' + params : ''}`);
}

export function getMessage(id: string) {
  return request<MessageDetail>(`/messages/${encodeURIComponent(id)}`);
}

export function sendMessage(data: SendMessageRequest) {
  return request<MessageDetail>('/messages', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function listDevices() {
  return request<Device[]>('/devices');
}

export function deleteDevice(id: string) {
  return request<void>(`/devices/${encodeURIComponent(id)}`, { method: 'DELETE' });
}
