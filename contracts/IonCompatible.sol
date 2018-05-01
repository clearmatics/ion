// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.18;

contract IonLinkInterface
{

    function Verify(uint256 block_id, uint256 leaf_hash, uint256[] proof)
		public view returns (bool);
}


contract IonCompatible
{
    event IonTransfer(address _recipient, address _currency, uint256 indexed value, bytes32 indexed ref, bytes data);

    event IonWithdraw(address _recipient, address _currency, uint256 indexed value, bytes32 indexed ref);

    event IonMint(uint256 value, bytes32 indexed ref);

    event IonBurn(uint256 value, bytes32 indexed ref);
}
