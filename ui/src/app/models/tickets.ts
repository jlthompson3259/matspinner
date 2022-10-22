export interface Tickets {
  id: number;
  tickets: number;
}

export interface TicketsResponse {
  tickets: Tickets[];
  error: string;
}