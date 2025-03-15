<?php
namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Session;

class WalletController extends Controller
{
    public function __construct()
    {
        $this->middleware('auth.token')->except(['login', 'authenticate', 'registerForm', 'register']);
    }

    public function index() {
        return view('welcome');
    }

    public function login() {
        return view('auth.login');
    }

    public function authenticate(Request $request) {
        $response = Http::post(env('API_URL') . '/auth/login', $request->all());
        if ($response->successful()) {
            Session::put('token', $response->json()['data']['token']);
            return redirect('/')->with('message', 'Login berhasil!');
        }
        return back()->withErrors(['error' => 'Login gagal!']);
    }

    public function registerForm() {
        return view('auth.register');
    }

    public function register(Request $request) {
        $request = $request->all();
        $request['birth_date'] =  $request['birth_date']."T00:00:00Z";
        $request['identity_photo_link'] = "test.jpg";

        $response = Http::post(env('API_URL') . '/auth/register', $request);
        if ($response->successful()) {
            return redirect('/login')->with('message', 'Registrasi berhasil! Silakan login.');
        }
        return back()->withErrors(['error' => 'Registrasi gagal!']);
    }

    public function logout() {
        Http::withHeaders(['Authorization' => Session::get('token')])->post(env('API_URL') . '/logout');
        Session::forget('token');
        return redirect('/login')->with('message', 'Logout berhasil!');
    }

    public function topupForm() {
        $virtualAccounts = [
            'bca' => 'BCA Virtual Account',
            'bni' => 'BNI Virtual Account',
            'bri' => 'BRI Virtual Account'
        ];
        return view('wallet.topup', compact('virtualAccounts'));
    }

    public function topup(Request $request) {
        $request = $request->all();
        $payload = [
            "amount_topup" => (float)$request['amount'],
            "amount_service"  => 0,
            "payment_type"  => "bank_transfer",
            "bank_name"  => $request['va'],
            "bank_code"  => $request['va']
        ];

        Http::withHeaders(['Authorization' => Session::get('token')])->post(env('API_URL') . '/wallet/topup', $payload);
        return redirect('/topup')->with('message', 'Top-up berhasil!');
    }

    public function payForm() {
        return view('wallet.pay');
    }

    public function pay(Request $request) {
        Http::withHeaders(['Authorization' => Session::get('token')])->post(env('API_URL') . '/wallet/pay', $request->all());
        return redirect('/pay')->with('message', 'Transaksi berhasil!');
    }

    public function transactions() {
        $transactions = Http::withHeaders(['Authorization' => Session::get('token')])->get(env('API_URL') . '/wallet/history')->json()['data'] ?? [];
        return view('wallet.transactions', compact('transactions'));
    }
}
