import { Component, OnInit } from '@angular/core';
import {WorkerService} from '../shared/worker.service';
import {NgForm} from '@angular/forms';
import {Worker} from '../shared/worker.model';
import {Observable, Subject} from 'rxjs';
import {AddserviceService} from './addservice.service';

@Component({
  selector: 'app-addworker',
  templateUrl: './addworker.component.html',
  styleUrls: ['./addworker.component.css']
})
export class AddworkerComponent implements OnInit {

  constructor(private workerService: WorkerService, private addService: AddserviceService) { }

  onAddWorker(form: NgForm){
    if( form.invalid ){
      return;
    }
    if(form.value._id === ''){
      this.workerService.postWorker(form.value).subscribe((res) => {
        this.resetForm();
      });
    }
    else{
      console.log(form.value);
      this.workerService.putWorker(form.value).subscribe((res) => {
        this.resetForm();
      });
    }
    this.addService.sendMessage('new guest')
  }

  resetForm(form?: NgForm){
    if(form){
      form.reset();
    }
    this.workerService.selectedWorker = {
      _id: '',
      name: '',
      number: '',
      disposition: '',
    };
    this.workerService.getWorkerList().subscribe((res) => {
      this.workerService.workers = res as Worker[];
    });
  }

  onDelete( worker: Worker ){
    if(confirm('Delete this record?') === true){
        this.workerService.deleteWorker(worker._id).subscribe((res) => {
            // this.refreshWorkerList();
        });
    }
  }

  ngOnInit() {
    this.resetForm();
  }

}
