export interface ShortenRequest {
  url: string;
}

export interface ShortenResponse {
  short_url: string;
}

export interface StatsResponse {
  count: number;
  original_url: string;
}

export interface ErrorResponse {
  error: string;
}
