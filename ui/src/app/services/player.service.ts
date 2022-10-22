import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { Player, PlayerResponse, PlayersResponse } from '../models/player';

@Injectable({
  providedIn: 'root',
})
export class PlayerService {
  constructor(private http: HttpClient) {}

  public getAllPlayers(): Observable<Player[]> {
    return this.http
      .get<PlayersResponse>('/players')
      .pipe(map((response) => response.players));
  }

  public addNewPlayer(name: string): Observable<Player> {
    return this.http
      .post<PlayerResponse>('/players', {
        name: name
      })
      .pipe(map((response) => response.player));
  }

  public updatePlayer(player: Player): Observable<Player> {
    return this.http
      .put<PlayerResponse>('/players', {
        player: player
      })
      .pipe(map((response) => response.player));
  }
}
