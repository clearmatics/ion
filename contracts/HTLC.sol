// Copyright (c) 2018 Harry Roberts. All Rights Reserved.
// SPDX-License-Identifier: LGPL-3.0+

pragma solidity 0.4.24;


contract HTLC {

    event OnDeposit( bytes32 exchGUID, address indexed receiver, bytes32 secretHashed, uint256 expiry );

    event OnRefund( bytes32 indexed exchGUID );

    event OnWithdraw( bytes32 indexed exchGUID, bytes32 secret );

    enum ExchangeState {
        Invalid,    // Default state, invalid
        Deposited,
        Withdrawn,
        Refunded,
        Expired
    }

    struct Exchange {
        bytes32 secretHashed;
        address sender;
        address receiver;
        uint256 amount;
        uint256 expiry;
        ExchangeState state;
    }

    mapping (bytes32 => Exchange) public exchanges;

    function GetExchange ( bytes32 inExchGUID )
        internal view returns (Exchange storage)
    {
        Exchange storage exch = exchanges[inExchGUID];
        require( exch.state != ExchangeState.Invalid );
        return exch;
    }

    function GetSender ( bytes32 inExchGUID )
        public view returns (address)
    {
        return GetExchange(inExchGUID).sender;
    }

    function GetReceiver ( bytes32 inExchGUID )
        public view returns (address)
    {
        return GetExchange(inExchGUID).receiver;
    }

    function GetSecretHashed ( bytes32 inExchGUID )
        public view returns (bytes32)
    {
        return GetExchange(inExchGUID).secretHashed;
    }

    function GetExpiry ( bytes32 inExchGUID )
        public view returns (uint256)
    {
        return GetExchange(inExchGUID).expiry;
    }

    function GetAmount ( bytes32 inExchGUID )
        public view returns (uint256)
    {
        return GetExchange(inExchGUID).amount;
    }

    function GetState ( bytes32 inExchGUID )
        public view returns (ExchangeState)
    {
        Exchange storage exch = exchanges[inExchGUID];

        if (exch.state == ExchangeState.Invalid) {
            return ExchangeState.Invalid;
        }

        if (exch.expiry < block.timestamp) {
            return ExchangeState.Expired;
        }

        return exch.state;
    }

    function Deposit ( address inReceiver, bytes32 inSecretHashed, uint256 inExpiry )
        public payable returns (bytes32)
    {
        // GUID must be predictable
        bytes32 exchGUID = sha256(abi.encodePacked(inReceiver, inSecretHashed));

        require(exchanges[exchGUID].state == ExchangeState.Invalid, "Duplicate exchange");

        require(inReceiver != address(0x0), "Invalid receiver address");

        require(inExpiry > block.timestamp, "Expiry not in future");

        exchanges[exchGUID] = Exchange(
            inSecretHashed,
            msg.sender,
            inReceiver,
            msg.value,
            inExpiry,
            ExchangeState.Deposited
        );

        emit OnDeposit( exchGUID, inReceiver, inSecretHashed, inExpiry);

        return exchGUID;
    }

    function Withdraw ( bytes32 inExchGUID, bytes32 inSecret )
        public
    {
        Exchange storage exch = exchanges[inExchGUID];

        require(exch.state == ExchangeState.Deposited, "Unknown exchange, or invalid state");

        require(exch.receiver == msg.sender, "Only receiver can Withdraw");

        require(block.timestamp <= exch.expiry, "Exchange expired");

        require(sha256(abi.encodePacked(inSecret)) == exch.secretHashed, "Bad secret");

        exch.state = ExchangeState.Withdrawn;

        msg.sender.transfer(exch.amount);

        emit OnWithdraw(inExchGUID, inSecret);
    }

    function Refund ( bytes32 inExchGUID )
        public
    {
        Exchange storage exch = exchanges[inExchGUID];

        require(exch.sender == msg.sender, "Only depositor can refund");

        require(block.timestamp > exch.expiry, "Exchange not expired, cannot refund");

        require(exch.state == ExchangeState.Deposited, "Unknown exchange, or invalid state");

        exch.state = ExchangeState.Refunded;

        exch.sender.transfer(exch.amount);

        emit OnRefund(inExchGUID);
    }
}
