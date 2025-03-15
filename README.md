# Wallet APP with laravel and GoLang
Golang Architecture with DDD (Domain Data Driven) and integrate with Payment Gateway Midtrans


## Project setup & RUN

```bash
-   run backend service first in /ewallet
-   next run website in /fe-wallet
-   php artisan serve
```


## Endpoint API
POST /auth/login: Untuk authentication. 
POST /auth/register: Untuk register user access. 

POST /walllet/topup: Untuk top-up saldo. 
POST /walllet/webhook: Untuk mendapatkan status dari midtrans via webhook. 
POST /walllet/pay: Untuk melakukan transaksi transfer ke pengguna lain. 
GET /walllet/history: Untuk melihat riwayat transaksi. 