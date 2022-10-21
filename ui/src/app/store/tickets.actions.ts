import { createAction, props } from '@ngrx/store';
import { Tickets } from '../models/tickets';

export const TicketsActions = {
  getTickets: createAction(
    '[Tickets API] Get Tickets',
    props<{ ids: number[] }>()
  ),
  getTicketsSuccess: createAction(
    '[Tickets API] Get Tickets Success',
    props<{ tickets: Tickets[] }>()
  ),
  getTicketsFailure: createAction(
    '[Tickets API] Get Tickets Failure',
    props<{ error: string }>()
  ),

  incrementTickets: createAction(
    '[Tickets API] Increment Tickets',
    props<{ ids: number[] }>()
  ),
  incrementTicketsSuccess: createAction(
    '[Tickets API] Increment Tickets Success',
    props<{ tickets: Tickets[] }>()
  ),
  incrementTicketsFailure: createAction(
    '[Tickets API] Increment Tickets Failure',
    props<{ error: string }>()
  ),

  setTickets: createAction(
    '[Tickets API] Set Tickets',
    props<{ tickets: Tickets[] }>()
  ),
  setTicketsSuccess: createAction(
    '[Tickets API] Set Tickets Success',
    props<{ tickets: Tickets[] }>()
  ),
  setTicketsFailure: createAction(
    '[Tickets API] Set Tickets Failure',
    props<{ error: string }>()
  ),
};
