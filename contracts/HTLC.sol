// Copyright (c) 2018 Harry Roberts. All Rights Reserved.
// SPDX-License-Identifier: LGPL-3.0+

pragma solidity ^0.4.23;

contract HTLC
{
  event OnDeposit( address indexed receiver, bytes32 image, uint256 expiry );

  event OnRefund( bytes32 indexed image );

  event OnWithdraw( bytes32 indexed image, bytes32 preimage );

  enum ExchangeState
  {
    Invalid,    // Default state, invalid
    Deposited,
    Withdrawn,
    Refunded,
    Expired
  }

  struct Exchange
  {
    address sender;
    address receiver;
    uint256 amount;
    uint256 expiry;
    ExchangeState state;
  }

  mapping (bytes32 => Exchange) public exchanges;


  function GetState ( bytes32 in_image )
    public view returns (ExchangeState)
  {
    Exchange storage exch = exchanges[in_image];

    if( exch.state == ExchangeState.Invalid )
    {
      return ExchangeState.Invalid;
    }

    if( exch.expiry < block.timestamp )
    {
      return ExchangeState.Expired;
    }

    return exch.state;
  }


  function Deposit ( address in_receiver, bytes32 in_image, uint256 in_expiry )
    public payable
  {
    require( exchanges[in_image].state == ExchangeState.Invalid,
             "Duplicate exchange" );

    require( in_receiver != address(0x0),
             "Invalid receiver address" );

    require( in_expiry > block.timestamp,
             "Expiry not in future" );

    exchanges[in_image] = Exchange(
      msg.sender,
      in_receiver,
      msg.value,
      in_expiry,
      ExchangeState.Deposited
    );

    emit OnDeposit( in_receiver, in_image, in_expiry );
  }


  function Withdraw ( bytes32 in_image, bytes32 in_preimage )
    public
  {
    Exchange storage exch = exchanges[in_image];

    require( exch.state == ExchangeState.Deposited,
             "Unknown exchange, or invalid state" );

    require( exch.receiver == msg.sender,
             "Only receiver can Withdraw" );

    require( block.timestamp <= exch.expiry,
              "Exchange expired" );

    require( sha256(abi.encodePacked(in_preimage)) == in_image,
             "Bad preimage" );

    exch.state = ExchangeState.Withdrawn;

    msg.sender.transfer( exch.amount );

    emit OnWithdraw( in_image, in_preimage );
  }


  function Refund ( bytes32 in_image )
    public
  {
    Exchange storage exch = exchanges[in_image];

    require( exch.sender == msg.sender,
             "Only depositor can refund" );

    require( block.timestamp > exch.expiry,
             "Exchange not expired, cannot refund" );

    require( exch.state == ExchangeState.Deposited,
             "Unknown exchange, or invalid state" );

    exch.state = ExchangeState.Refunded;

    exch.sender.transfer( exch.amount );

    emit OnRefund( in_image );
  }
}