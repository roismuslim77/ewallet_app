@extends('layouts.app')

@section('content')
<div class="container mt-4">
    <h2>Top-Up Saldo</h2>

    @if(session('message'))
        <div class="alert alert-success">
            {{ session('message') }}
        </div>
    @endif

    <form method="POST" action="/topup">
        @csrf
        <div class="mb-3">
            <label for="amount" class="form-label">Jumlah Top-Up</label>
            <input type="number" class="form-control" id="amount" name="amount" required>
        </div>

        <div class="mb-3">
            <label for="va" class="form-label">Pilih Virtual Account</label>
            <select class="form-control" id="va" name="va" required>
                @foreach($virtualAccounts as $key => $value)
                    <option value="{{ $key }}">{{ $value }}</option>
                @endforeach
            </select>
        </div>

        <button type="submit" class="btn btn-primary">Top-Up</button>
    </form>
</div>
@endsection
