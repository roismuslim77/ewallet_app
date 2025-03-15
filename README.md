# Wallet APP with laravel and GoLang
Golang Architecture with DDD (Domain Data Driven) and integrate with Payment Gateway Midtrans


## Project setup & RUN

```bash
-   run backend service first in /ewallet
-   next run website in /fe-wallet
```


## Endpoint API
```bash
- POST /auth/login: Untuk authentication. 
- POST /auth/register: Untuk register user access. 

- POST /wallet/topup: Untuk top-up saldo. 
- POST /wallet/webhook: Untuk mendapatkan status dari midtrans via webhook. 
- POST /wallet/pay: Untuk melakukan transaksi transfer ke pengguna lain. 
- GET /wallet/history: Untuk melihat riwayat transaksi. 
```