@extends('layouts.app')
@section('content')
<div class="container">
    <h2>Pay</h2>
    <form method="POST" action="/pay">
        @csrf
        <div class="mb-3">
            <label for="amount" class="form-label">Amount</label>
            <input type="number" class="form-control" id="amount" name="amount" required>
        </div>
        <button type="submit" class="btn btn-danger">Pay</button>
    </form>
</div>
@endsection
