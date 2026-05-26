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
