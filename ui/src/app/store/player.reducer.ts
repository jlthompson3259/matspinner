import { createReducer, on } from '@ngrx/store';
import { PlayerActions } from './player.actions';

export const initialState: ReadonlyMap<number, string> = new Map();

export const playerReducer = createReducer(
  initialState,
  on(PlayerActions.addPlayerSuccess, (state, action) => ({
    ...state,
    [action.player.id]: action.player.name,
  })),
  on(PlayerActions.getAllPlayersSuccess, (state, action) =>
    action.players.reduce(
      (s, p) => ({
        ...s,
        [p.id]: p.name,
      }),
      initialState
    )
  ),
  on(PlayerActions.updatePlayerSuccess, (state, action) => ({
    ...state,
    [action.player.id]: action.player.name,
  }))
);
