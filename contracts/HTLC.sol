// Copyright (c) 2018 Harry Roberts. All Rights Reserved.
// SPDX-License-Identifier: LGPL-3.0+

pragma solidity ^0.4.23;

contract HTLC
{
  event OnDeposit( address indexed receiver, bytes32 image, uint256 expiry );

  event OnRefund( bytes32 indexed exchange_id );

  event OnWithdraw( bytes32 indexed exchange_id );

  struct Exchange
  {
    address sender;
    address receiver;
    uint256 amount;
    bytes32 image;
    uint expiry;
  }

  mapping (bytes32 => Exchange) public exchanges;

  function Deposit ( address receiver, bytes32 image, uint256 expiry )
    public payable
  {
    bytes32 exch_id = sha256( abi.encodePacked(receiver, image, expiry) );

    require( exchanges[exch_id].sender == address(0x0), "Duplicate exchange" );
    require( receiver != address(0x0), "Invalid receiver address" );
    require( expiry > block.timestamp, "Expiry not in future" );

    exchanges[exch_id] = Exchange(
      msg.sender,
      receiver,
      msg.value,
      image,
      expiry
    );

    emit OnDeposit( receiver, image, expiry );
  }


  function Withdraw ( bytes32 in_id, bytes32 in_preimage )
    public
  {
    Exchange storage exch = exchanges[in_id];
    require( exch.sender != 0x0, "Non existant exchange" );
    require( exch.receiver == msg.sender, "Only receiver can Withdraw" );
    require( block.timestamp < exch.expiry, "Exchange expired" );
    require( sha256(abi.encodePacked(in_preimage)) == exch.image, "Bad preimage" );

    msg.sender.transfer( exch.amount );

    delete exchanges[in_id];

    emit OnWithdraw( in_id );
  }


  function Refund ( bytes32 in_id )
    public
  {
    Exchange storage exch = exchanges[in_id];
    require( exch.sender == msg.sender, "Only depositor can refund" );
    require( block.timestamp > exch.expiry, "Exchange not expired, cannot refund" );

    exch.sender.transfer( exch.amount );

    delete exchanges[in_id];

    emit OnRefund( in_id );
  }
}