export interface Player {
  id: number;
  name: string;
}

export interface PlayerResponse {
  player: Player,
  error: string
}

export interface PlayersResponse {
  players: Player[],
  error: string
}