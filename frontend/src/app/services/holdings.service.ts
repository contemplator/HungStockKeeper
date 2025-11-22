import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Holding {
  ID?: number;
  symbol: string;
  quantity: number;
  cost_basis: number;
  purchase_date: Date;
  note?: string;
  CreatedAt?: Date;
  UpdatedAt?: Date;
}

@Injectable({
  providedIn: 'root'
})
export class HoldingsService {
  private apiUrl = '/api/holdings';

  constructor(private http: HttpClient) { }

  getHoldings(): Observable<Holding[]> {
    return this.http.get<Holding[]>(this.apiUrl);
  }

  getHolding(id: number): Observable<Holding> {
    return this.http.get<Holding>(`${this.apiUrl}/${id}`);
  }

  createHolding(holding: Holding): Observable<Holding> {
    return this.http.post<Holding>(this.apiUrl, holding);
  }

  updateHolding(id: number, holding: Partial<Holding>): Observable<Holding> {
    return this.http.put<Holding>(`${this.apiUrl}/${id}`, holding);
  }

  deleteHolding(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }
}
