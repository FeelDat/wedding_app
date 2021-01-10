import { Component, OnInit } from '@angular/core';
import {AuthenticationService} from '../auth/authentication.service';
import { NgForm, NgModel } from '@angular/forms';
import { WorkerService } from './shared/worker.service';
import { Worker } from './shared/worker.model';
import { Router } from '@angular/router';

@Component({
  selector: 'app-workers',
  templateUrl: './workers.component.html',
  styleUrls: ['./workers.component.css']
})
export class WorkersComponent implements OnInit {

  constructor(private workerService: WorkerService,
    private router: Router,
    private authenticationService: AuthenticationService
    ) { }

  refreshWorkerList(){
    this.workerService.getWorkerList().subscribe((res) => {
      this.workerService.workers = res as Worker[];
    });
  }

  onEdit( worker: Worker){
      this.workerService.selectedWorker = worker;
    }

  // onDelete( worker: Worker ){
  //   if(confirm('Delete this record?') === true){
  //       this.workerService.deleteWorker(worker._id).subscribe((res) => {
  //           this.refreshWorkerList();
  //       });
  //   }
  // }

  ngOnInit() {
    this.refreshWorkerList();
    // this.resetForm();
  }

  logout(){
    this.authenticationService.logout();
    location.reload()
  }

}
