<?php
use Illuminate\Support\Facades\Route;
use App\Http\Controllers\WalletController;
use App\Http\Controllers\HomeController;

Route::get('/', [HomeController::class, 'index']);
Route::get('/login', [WalletController::class, 'login'])->name('login');
Route::post('/login', [WalletController::class, 'authenticate']);
Route::get('/logout', [WalletController::class, 'logout'])->name('logout');

Route::get('/register', [WalletController::class, 'registerForm'])->name('register');
Route::post('/register', [WalletController::class, 'register']);

Route::get('/topup', [WalletController::class, 'topupForm']);
Route::post('/topup', [WalletController::class, 'topup']);
Route::get('/pay', [WalletController::class, 'payForm']);
Route::post('/pay', [WalletController::class, 'pay']);
Route::get('/transactions', [WalletController::class, 'transactions']);
