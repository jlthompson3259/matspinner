import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Tickets, TicketsResponse } from '../models/tickets';
import { map, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class TicketsService {
  constructor(private http: HttpClient) {}

  public getTickets(ids: number[]): Observable<Tickets[]> {
    return this.http
      .get<TicketsResponse>('/tickets', {
        params: new HttpParams().append('ids', ids.join(',')),
      })
      .pipe(map((response) => response.tickets));
  }

  public incrementTickets(ids: number[]): Observable<Tickets[]> {
    return this.http
      .post<TicketsResponse>('/tickets/increment', {
        ids: ids.join(','),
      })
      .pipe(map((response) => response.tickets));
  }

  public setTickets(tickets: Tickets[]): Observable<Tickets[]> {
    return this.http
      .put<TicketsResponse>('/tickets/increment', {
        tickets: tickets,
      })
      .pipe(map((response) => response.tickets));
  }
}
