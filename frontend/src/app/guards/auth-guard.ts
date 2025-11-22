import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';

export const authGuard: CanActivateFn = (route, state) => {
  const router = inject(Router);
  
  // 由於我們使用 HttpOnly Cookie，前端 JS 無法直接讀取 Cookie 來判斷 JWT 是否存在。
  // 因此，我們檢查 localStorage 中的 user 資料作為前端路由判斷的依據。
  // 真正的安全性由後端 API 驗證 Cookie 來把關。
  const user = localStorage.getItem('user');

  if (user) {
    return true;
  }

  // 未登入，導向登入頁
  router.navigate(['/login']);
  return false;
};
