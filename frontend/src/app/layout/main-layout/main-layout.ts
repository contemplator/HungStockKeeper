import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { PopoverModule } from 'primeng/popover';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive, ButtonModule, PopoverModule],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss',
})
export class MainLayout implements OnInit {
  userEmail: string = '';

  constructor(private router: Router, private authService: AuthService) {}

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
