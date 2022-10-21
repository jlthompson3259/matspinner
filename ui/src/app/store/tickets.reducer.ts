import { createReducer, on } from '@ngrx/store';
import { TicketsActions } from './tickets.actions';

export const initialState: ReadonlyMap<number, string> = new Map();

export const ticketReducer = createReducer(
  initialState,
  on(TicketsActions.getTicketsSuccess, (state, action) =>
    action.tickets.reduce(
      (s, t) => ({
        ...s,
        [t.id]: t.tickets,
      }),
      state
    )
  ),
  on(TicketsActions.setTicketsSuccess, (state, action) =>
    action.tickets.reduce(
      (s, t) => ({
        ...s,
        [t.id]: t.tickets,
      }),
      state
    )
  ),
  on(TicketsActions.incrementTicketsSuccess, (state, action) =>
    action.tickets.reduce(
      (s, t) => ({
        ...s,
        [t.id]: t.tickets,
      }),
      state
    )
  )
);
