import { createAction, props } from '@ngrx/store';
import { Player } from '../models/player';

export const PlayerActions = {
  addPlayer: createAction(
    '[Player API] Add New Player',
    props<{ name: string }>()
  ),
  addPlayerSuccess: createAction(
    '[Player API] Add New Player Success',
    props<{ player: Player }>()
  ),
  addPlayerFailure: createAction(
    '[Player API] Add New Player Failure',
    props<{ error: string }>()
  ),

  getAllPlayers: createAction('[Player API] Get All Players'),
  getAllPlayersSuccess: createAction(
    '[Player API] Get All Players Success',
    props<{ players: Player[] }>()
  ),
  getAllPlayersFailure: createAction(
    '[Player API] Get All Players Failure',
    props<{ error: string }>()
  ),

  updatePlayer: createAction(
    '[Player API] Update Player',
    props<{ player: Player }>()
  ),
  updatePlayerSuccess: createAction(
    '[Player API] Update Player Success',
    props<{ player: Player }>()
  ),
  updatePlayerFailure: createAction(
    '[Player API] Update Player Failure',
    props<{ error: string }>()
  ),
};
