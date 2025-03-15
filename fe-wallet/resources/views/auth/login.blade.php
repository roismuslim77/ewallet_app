@extends('layouts.app')
@section('content')
<div class="container">
    <h2>Login</h2>
    <form method="POST" action="/login">
        @csrf
        <div class="mb-3">
            <label for="email" class="form-label">Email</label>
            <input type="email" class="form-control" id="email" name="email" required>
        </div>
        <div class="mb-3">
            <label for="password" class="form-label">Password</label>
            <input type="password" class="form-control" id="password" name="password" required>
        </div>
        <button type="submit" class="btn btn-primary">Login</button>
        <a href="/register">Buat Akun baru</button>
    </form>
</div>
@endsection
