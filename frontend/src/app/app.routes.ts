import { Routes } from '@angular/router';
import { LoginComponent } from './pages/login/login.component';
import { RegisterComponent } from './pages/register/register.component';
import { MainLayout } from './layout/main-layout/main-layout';
import { Dashboard } from './pages/dashboard/dashboard';
import { Holdings } from './pages/holdings/holdings';
import { Watchlist } from './pages/watchlist/watchlist';
import { authGuard } from './guards/auth-guard';

export const routes: Routes = [
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    {
        path: 'app',
        component: MainLayout,
        canActivate: [authGuard],
        children: [
            { path: 'dashboard', component: Dashboard },
            { path: 'holdings', component: Holdings },
            { path: 'watchlist', component: Watchlist },
            { path: '', redirectTo: 'dashboard', pathMatch: 'full' }
        ]
    },
    { path: '', redirectTo: '/login', pathMatch: 'full' },
    { path: '**', redirectTo: '/login' }
];
