import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';
import { SpinResult, SpinResultResponse } from '../models/spin';

@Injectable({
  providedIn: 'root',
})
export class SpinService {
  constructor(private http: HttpClient) {}

  public spin(participantIds: number[]): Observable<SpinResult> {
    return this.http
      .post<SpinResultResponse>('/spin', {
        participantIds: participantIds,
      })
      .pipe(map((response) => response.result));
  }

  public getLastSpin(): Observable<SpinResult> {
    return this.http
      .get<SpinResultResponse>('/get-last-spin')
      .pipe(map((response) => response.result));
  }
}
