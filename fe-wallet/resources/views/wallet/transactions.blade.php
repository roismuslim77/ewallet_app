@extends('layouts.app')
@section('content')
<div class="container">
    <h2>Transaction History</h2>
    <ul class="list-group">
        @foreach($transactions as $transaction)
            <li class="list-group-item">{{ $transaction['description'] }} [{{$transaction['no_acc']}}]</li>
        @endforeach
    </ul>
</div>
@endsection
