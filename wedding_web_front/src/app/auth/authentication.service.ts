import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import {User} from '../_helpers/user';

@Injectable({ providedIn: 'root' })
export class AuthenticationService {
  private currentUserSubject: BehaviorSubject<User>;
  public currentUser: Observable<User>;

  constructor(private http: HttpClient) {
    this.currentUserSubject = new BehaviorSubject<User>(JSON.parse(localStorage.getItem('currentUser')));
    this.currentUser = this.currentUserSubject.asObservable();
  }

  public get currentUserValue(): User {
    return this.currentUserSubject.value;
  }

  login(username: string, password: string) {
    return this.http.post<any>(`http://127.0.0.1:18000/login`+`/${username}`+`/${password}`, username)
      .pipe(map(user => {
        // store user details and jwt token in local storage to keep user logged in between page refreshes
        localStorage.setItem('currentUser', JSON.stringify(user));
        this.currentUserSubject.next(user);
        console.log(user)
        return user;
      }));
  }
  readonly baseUrl = 'http://127.0.0.1:18000/logout';

  logout() {
     return this.http.post<any>(this.baseUrl,null).subscribe(res => { 
      localStorage.removeItem('currentUser');
      this.currentUserSubject.next(null);
    }, error => {
      console.log("Error", error);
    });
  }
}
