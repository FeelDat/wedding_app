import {Component, OnDestroy, OnInit} from '@angular/core';
import {CdkDragDrop, moveItemInArray, transferArrayItem} from '@angular/cdk/drag-drop';
import {WorkerService} from '../workers/shared/worker.service';
import {Worker} from '../workers/shared/worker.model';

import {Subscription} from 'rxjs';
import {AddserviceService} from '../workers/addworker/addservice.service';

@Component({
  selector: 'app-disposition',
  templateUrl: './disposition.component.html',
  styleUrls: ['./disposition.component.css']
})
export class DispositionComponent implements OnInit, OnDestroy  {
  zeroTable: Worker[] = [];
  firstTable: Worker[] = [];
  secondTable: Worker[] = [];
  public isOpened = false;
  subscription: Subscription;

  constructor(private workerService: WorkerService, private addWorker: AddserviceService) {
    this.subscription = this.addWorker.getMessage().subscribe(message => {
      if (message) {
        this.workerService.getWorkerList().subscribe((res) => {
          this.workerService.workers = res as Worker[];
        });
        this.getGuestsList();
      }
    });
  }
  getGuestsList(){
    console.log(this.tCount)
    // this.tCount = this.tCount+1;
    let num = 0;
    this.tables = [];
    this.tables = Array.from(Array(this.tCount), () => new Array(0));
    this.hidden = Array(this.tCount).fill(true);

    this.workerService.getWorkerList().subscribe((res) => {
      this.workerService.workers = res as Worker[];
      while (num <= this.tCount) {
        this.tables[num] = this.workerService.workers.filter(x=>parseInt(x.disposition) == num)
        num += 1;
      }
    });
    // console.log(this.tables)
    // this.tCount = 0;
  };
  onEdit( worker: Worker){
    this.workerService.selectedWorker = worker;
  }

  drop(event: CdkDragDrop<Worker[]>, table: string) {
    if (event.previousContainer === event.container) {
      moveItemInArray(event.container.data, event.previousIndex, event.currentIndex);
      this.entered(event.container.data as Worker[], table)
    } else {
      transferArrayItem(event.previousContainer.data,
        event.container.data,
        event.previousIndex,
        event.currentIndex);
      this.entered(event.container.data as Worker[], table)
    }
  }
  entered(guests: Worker[], table: string){
    for(const i of guests) {
      if(i.disposition !== table) {
        i.disposition = table;
        this.workerService.putWorker(i).subscribe((res) => {

        });
      }
    }
  }

  addWorkerShow(){
    this.isOpened = !this.isOpened;
  }

  dropDisposition() {
    this.workerService.dropDisposition().subscribe((res) => {
      this.getGuestsList();
    });
  }

  ngOnInit() {
    this.getGuestsList()
  }
  ngOnDestroy() {
    // unsubscribe to ensure no memory leaks
    this.subscription.unsubscribe();
  }
}
