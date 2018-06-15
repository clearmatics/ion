// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.18;

contract ERC223ReceivingContract {
    function tokenFallback(address _from, uint _value, bytes _data) public;
}

contract ERC223 {
    uint public totalSupply;

    function balanceOf(address who) constant public returns (uint);

    function transfer(address to, uint value) public;

    function transfer(address to, uint value, bytes data) public;

    event Transfer(address indexed from, address indexed to, uint value, bytes data);
}
