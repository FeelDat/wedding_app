import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {

  constructor() { }

  pages = [
    {name:"Список гостей",url:"guests"},
    {name:"Рассадка",url:"disposition"},
    {name:"Пригласительное",url:"#"}
  ]
  ngOnInit(): void {
    console.log(this.pages)
  }

}
