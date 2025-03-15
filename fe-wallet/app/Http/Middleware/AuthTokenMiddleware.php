<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Session;

class AuthTokenMiddleware
{
    public function handle(Request $request, Closure $next)
    {
        if (!Session::has('token')) {
            return redirect('/login')->withErrors(['error' => 'Silakan login terlebih dahulu!']);
        }

        return $next($request);
    }
}
