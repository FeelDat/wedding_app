import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
// import { Observable } from 'rxjs/Observable';

// import 'rxjs/add/operator/map';
// import 'rxjs/add/operator/toPromise';

import { Worker } from './worker.model';

@Injectable({
  providedIn: 'root'
})
export class WorkerService {
  selectedWorker: Worker;
  workers: Worker[];

  readonly baseUrl = 'http://127.0.0.1:18000/guests';

  constructor( private http : HttpClient ) { }

  postWorker( worker: Worker ){
    return this.http.post(this.baseUrl, worker);
  }

  getWorkerList(){
    return this.http.get(this.baseUrl);
  }

  putWorker( worker: Worker ){
    return this.http.put(this.baseUrl + `/${worker._id}`, worker);
  }

  deleteWorker( _id: string ){
    return this.http.delete(this.baseUrl + `/${_id}`);
  }

  dropDisposition(){
    return this.http.delete('http://127.0.0.1:18000/disposition')
  }
}
