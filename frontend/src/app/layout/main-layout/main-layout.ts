import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink, RouterLinkActive, RouterOutlet, NavigationEnd } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { PopoverModule } from 'primeng/popover';
import { AuthService } from '../../services/auth.service';
import { filter } from 'rxjs/operators';

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive, ButtonModule, PopoverModule],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss',
})
export class MainLayout implements OnInit {
  userEmail: string = '';
  pageTitle: string = 'Hung Stock Keeper';

  constructor(private router: Router, private authService: AuthService) {
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe((event: any) => {
      this.updateTitle(event.url);
    });
  }

  ngOnInit() {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        this.userEmail = user.email;
      } catch (e) {
        console.error('Error parsing user data', e);
      }
    }
    this.updateTitle(this.router.url);
  }

  updateTitle(url: string) {
    if (url.includes('/app/holdings')) {
      this.pageTitle = '庫存管理';
    } else if (url.includes('/app/dashboard')) {
      this.pageTitle = 'Dashboard';
    } else if (url.includes('/app/watchlist')) {
      this.pageTitle = 'Watchlist';
    } else {
      this.pageTitle = 'Hung Stock Keeper';
    }
  }

  logout() {
    this.authService.logout().subscribe({
      next: () => {
        localStorage.removeItem('user');
        this.router.navigate(['/login']);
      },
      error: (err) => {
        console.error('Logout failed', err);
        // Force logout on client side even if server fails
        localStorage.removeItem('user');
        this.router.navigate(['/login']);
      }
    });
  }
}
