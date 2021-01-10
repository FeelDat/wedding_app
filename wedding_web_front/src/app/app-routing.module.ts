import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { GreetpageComponent } from './greetpage/greetpage.component';
import { WorkersComponent } from './workers/workers.component';
import { DispositionComponent } from './disposition/disposition.component'
import {LoginComponent} from './login/login.component';
import {AuthGuard} from './_helpers/auth.guard';
import { WelcomeComponent } from './welcome/welcome.component';


const routes: Routes = [
    { path: '', pathMatch: 'full', component: GreetpageComponent },
    { path: 'guests', component: WorkersComponent, canActivate: [AuthGuard] },
    { path: 'disposition', component: DispositionComponent, canActivate: [AuthGuard] },
    { path: 'login', component: LoginComponent },
    { path: 'welcome', component: WelcomeComponent },
    { path: '**', redirectTo: '' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
