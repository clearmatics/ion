// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "./ERC223.sol";

contract ERC223ReceivingContract {
    function tokenFallback(address _from, uint _value, bytes32 _ref) public;
}
