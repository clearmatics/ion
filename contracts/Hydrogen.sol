pragma solidity ^0.4.18;

import "./Ion.sol";
import "./ERC223Compatible.sol";

contract Hydrogen is Ion, ERC223ReceivingContract {

    mapping (bytes32 => TradeAgreement) internal m_trades;
    mapping (address => bytes32[]) internal m_trades_initiated;
    mapping (bytes32 => bool) internal m_trade_deposited;

    ERC223 internal currency;

    /*
    InitiateTradeAgreement
        Creates a new trade agreement between a sender and recipient of a token with withdraw and refund hash
        Hashes are pre-agreed between the parties. The first depositor must choose the refund reference and the
        second depositor must choose the withdraw reference. Both parties then must share the hash of their chosen
        reference and construct trade agreements on each chain with this function.

    Arguments:
        _token: ERC223 Token address being traded
        _recipient: address of recipient of trade
        _value: amount of token being traded
        _withdrawHash: keccak256 hash of withdraw reference
        _refundHash: keccak256 hash of refund reference

    Returns:
        trade_id: keccak256 hash of the sender, recipient, token address, value, withdraw hash and refund hash

    Events:
        OnTradeInitiated
    */
    function InitiateTradeAgreement(
        address _token,
        address _initiator,
        address _counterparty,
        uint256 _value,
        bytes32 _withdrawHash,
        bytes32 _refundHash)
        public returns (bytes32)
    {
        require( _value > 0 );
        bytes32 tradeId = keccak256(_initiator, _counterparty, _token, _value, _withdrawHash, _refundHash);

        TradeAgreement storage trade = m_trades[tradeId];
        require( trade.value == 0, "Identical trade agreement already exists" );
        trade.initiator = _initiator;
        trade.counterparty = _counterparty;
        trade.token = _token;
        trade.value = _value;
        trade.withdrawHash = _withdrawHash;
        trade.refundHash = _refundHash;

        if ( msg.sender == _initiator ) {
            trade.fundsUnlocked = false;
        } else if ( msg.sender = _counterparty ) {
            trade.fundsUnlocked = true;
        } else {
            revert( "Trade initiated by non-participant" );
        }

        bytes32[] storage tradesInitiatedBySender = m_trades_initiated[msg.sender];
        tradesInitiatedBySender.push(tradeId);

        emit OnTradeInitiated(tradeId, msg.sender);

        return tradeId;
    }

    /*
    tokenFallback
        Fallback function when ERC223 tokens are transferred to this contract. Reverts the payment if a valid trade
        agreement was not found.

    Arguments:
        _from: address of sender of token
        _value: amount of token sent
        _data: byte array that should contain the trade_id in the first 32 bytes

    Events:
        OnDeposit
    */
    function tokenFallback(address _from, uint _value, bytes _data) public {
        require( _data.length == 32, "Data length not 32" );

        bytes32 trade_id = bytesToBytes32(_data, 0);

        TradeAgreement storage trade = m_trades[trade_id];

        //TODO: Return the funds to sender if trade agreement is not found
        require( trade.token == msg.sender, "Trade sender not matched" );
        require( trade.value == _value, "Trade value not matched" );

        m_trade_deposited[trade_id] = true;

        emit OnDeposit(trade_id, _from);
    }

    function getTradesWithAddress(address _owner) public view returns (bytes32[]) {
        return m_trades_initiated[_owner];
    }

    function UnlockFunds(bytes32 _tradeId) public {
        TradeAgreement storage trade = m_trades[_tradeId];
        require( trade.value > 0, "Trade agreement does not exist" );
        require( msg.sender == trade.initiator );

        trade.fundsUnlocked = true;

        emit OnFundsUnlocked(_tradeId);
    }

    function CheckCounterpartyDeposit(bytes32 _tradeId, address _counterparty) returns (bool) {
        TradeAgreement storage trade = m_trades[_tradeId];
        require( trade.value > 0, "Trade agreement does not exist" );
        require( msg.sender == trade.initiator, "Caller is not Initiator" );
        require( _counterparty == trade.counterparty, "Supplied counterparty address is not counterparty in trade agreement" );

        return m_trade_deposited[_tradeId];
    }

    /*
    Withdraw
        Given a valid trade_id and withdraw reference, will verify the trade agreement and credit the caller with
        funds if they are the designated recipient of the funds under the trade agreement.

    Arguments:
        _trade_id: The trade agreement hash identifying the trade
        _withdrawRef: The pre-image of the withdraw hash used to form the trade agreement

    Events:
        OnWithdraw
    */
    function Withdraw(bytes32 _tradeId, bytes _withdrawRef) public {
        TradeAgreement storage trade = m_trades[_tradeId];
        require( trade.value > 0, "Trade agreement does not exist" );
        require( msg.sender == trade.recipient );
        require( trade.fundsUnlocked );

        bytes32 withdrawHash = keccak256(_withdrawRef);
        require( trade.withdrawHash == withdrawHash, "Withdraw reference provided is incorrect" );

        ERC223 token = ERC223(trade.token);
        token.transfer(trade.recipient, trade.value);

        emit OnWithdraw(_tradeId, msg.sender, _withdrawRef);
    }

    /*
    Refund
        Given a valid trade_id and refund reference, will verify the trade agreement and return the funds to the caller
        if they are the original depositor of the funds under the trade agreement.

    Arguments:
        _trade_id: The trade agreement hash identifying the trade
        _refundRef: The pre-image of the refund hash used to form the trade agreement

    Events:
        OnRefund
    */
    function Refund(bytes32 _tradeId, bytes _refundRef) public {
        TradeAgreement storage trade = m_trades[_tradeId];
        require( trade.value > 0, "Trade agreement does not exist" );
        require( msg.sender == trade.owner );

        bytes32 refundHash = keccak256(_refundRef);
        require( trade.refundHash == refundHash, "Refund reference provided is incorrect" );

        ERC223 token = ERC223(trade.token);
        token.transfer(trade.owner, trade.value);

        emit OnRefund(_tradeId, msg.sender, _refundRef);
    }

    /*
------------------------------------------------------------------------------------------------------------------------
    Helper functions
------------------------------------------------------------------------------------------------------------------------
    */

    // This assembly doesn't work
    function bytesToBytes32(bytes _data) private view returns (bytes32) {
        bytes32 trade_id;
        // Assembly for byte copying is cheaper than byte array iteration
        assembly {
            let ret := staticcall(3000, 4, add(_data, 32), 32, trade_id, 32)
            switch ret case 0 { invalid }
        }
        return trade_id;
    }

    // Takes a byte array and prunes the first 32 bytes into a bytes32 variable
    function bytesToBytes32(bytes b, uint offset) private pure returns (bytes32) {
        bytes32 out;

        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[offset + i] & 0xFF) >> (i * 8);
        }
        return out;
    }

}
