export interface Stats {
  messagesSent: number;
  messagesPending: number;
  messagesFailed: number;
  devicesOnline: number;
  devicesTotal: number;
}

export interface Me {
  login: string;
}

export interface LoginRequest {
  login: string;
  password: string;
}

export interface LoginResponse {
  login: string;
}

export interface Device {
  id: string;
  name: string;
  lastSeen: string;
  isOnline: boolean;
  createdAt: string;
}

export interface MessageListItem {
  id: string;
  deviceId: string;
  state: string;
  recipients: RecipientState[];
  textPreview: string;
  createdAt: string;
}

export interface MessageDetail {
  id: string;
  deviceId: string;
  state: string;
  isHashed: boolean;
  isEncrypted: boolean;
  recipients: RecipientState[];
  states: Record<string, string>;
  textMessage?: { text: string };
  dataMessage?: { data: string; port: number };
  hashedMessage?: { hash: string };
}

export interface RecipientState {
  phoneNumber: string;
  state: string;
  error?: string;
}

export interface SendMessageRequest {
  phoneNumbers: string[];
  text: string;
  simNumber?: number;
  priority?: number;
}

export interface ListMessagesResponse {
  items: MessageListItem[];
  total: number;
}
