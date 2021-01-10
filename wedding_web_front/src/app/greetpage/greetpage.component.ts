import { Component, OnInit } from '@angular/core';

import * as $ from 'jquery';

@Component({
  selector: 'app-greetpage',
  templateUrl: './greetpage.component.html',
  styleUrls: ['./greetpage.component.css']
})
export class GreetpageComponent implements OnInit {
  static counter: number = 0;

  constructor() { }

  author:string;

  ngOnInit() {
    this.author = 'Copyright © 2019 || designed by @2L';

    $('.section1').hover(function() {
        $('.choose_section').toggleClass('one-is-active');
    });

    $('.section2').hover(function() {
        $('.choose_section').toggleClass('two-is-active');
    });

    $('.scroll').click(function() {
        GreetpageComponent.counter++;
        if(GreetpageComponent.counter % 2 === 0){
          $('html, body').animate({scrollTop: $(".header").offset().top}, 800);
        }
        else{
          $('html, body').animate({scrollTop: $(".choose_section").offset().top}, 800);
        }
      // tslint:disable-next-line:triple-equals
        if (  $( this ).css( 'transform' ) == 'none' ){
          $(this).css("transform","rotate(-180deg)");
        }
        else {
          $(this).css("transform","" );
        }
    });
  }

}
