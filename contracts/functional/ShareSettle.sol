// Copyright (c) 2016-2019 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "../storage/FabricStore.sol";
import "../libraries/ERC223Compatible.sol";
import "../libraries/RLP.sol";

/*
    This contract serves as a example of how to perform DvP using the Ion framework contracts, between
    Ethereum and Fabric networks.
*/

contract ShareSettle is ERC223ReceivingContract {
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLP for bytes;

    FabricStore blockStore;
    ERC223 sharesettle_currency;

    uint256 ionlock_balance;

    event State(uint blockNo, uint txNo, string value);
    event ShareTransfer(address _sender, address _receiver, string _org, uint256 _amount, uint256 _price, bytes32 ref);
    event ShareTrade(address _sender, address _receiver, string _org, uint256 _amount, uint256 _price, bytes32 ref);

    // Logs the unique references that have been submitted
    mapping(bytes32 => bool) public m_opened_trades;
    mapping(bytes32 => bool) public m_settled_trades;

    // Stores data associated with three trade stages above
    mapping(bytes32 => Trade) public m_trades;

    // Structure of the trade agreement
    struct Trade {
        string  org;
        address send;
        address recv;
		uint256 amount;
		uint256 price;
		uint256 total;
        bytes32 ref;
	}

    constructor(ERC223 _currency, address _storeAddr) public {
        assert(address(_currency)!=0);

        sharesettle_currency = _currency;
        blockStore = FabricStore(_storeAddr);
    }

/*
========================================================================================================================

    Payment Functions

========================================================================================================================
*/

    /**
    * Creates a new trade agreement for a cross-chain DvP payment between two counterparties
    *
    * @param _org Name of shares being traded
    * @param _recv Intended recipient
    * @param _amount Amount of tokens msg.sender will pay
    * @param _price Amount of tokens msg.sender will receive
    * @param _ref A reference to escrow the funds
    */
    function initiateTrade(string _org, address _recv, uint256 _amount, uint256 _price, bytes32 _ref) public {
        // Assert no trade has been initiated with this reference
        assert(!m_opened_trades[_ref]);
        m_opened_trades[_ref]=true;

        // Instantiate the trade agreement
        m_trades[_ref].send = msg.sender;
        m_trades[_ref].recv = _recv;
        m_trades[_ref].org = _org;
        m_trades[_ref].amount = _amount;
        m_trades[_ref].price = _price;
        m_trades[_ref].total = _amount*_price;        
    }
    
    /**
    * When ERC223 tokens are sent to this contract it escrows them and emits an event if msg sender is the lead counterparty
    * else it sends transfers them directly to the lead counterparty if sender is the follow couterparty.
    *
    * @param _from Who sent us the token
    * @param _value Amount of tokens
    * @param _ref Arbitrary data, to be used as the payment reference
    */
    function tokenFallback(address _from, uint _value, bytes32 _ref) public {
        assert(msg.sender==address(sharesettle_currency));
        assert(_value>0 && _value==m_trades[_ref].total);
        assert((ionlock_balance+_value)>ionlock_balance);
        assert(_from==m_trades[_ref].send || _from==m_trades[_ref].recv);
        assert(_value==m_trades[_ref].total);

        ionlock_balance += _value;

        emit ShareTransfer(
            m_trades[_ref].send,
            m_trades[_ref].recv,
            m_trades[_ref].org,
            m_trades[_ref].amount,
            m_trades[_ref].price,
            _ref
        );

    }

/*
========================================================================================================================

    Verification Function

========================================================================================================================
*/

    /**
    * When called transfers the escrowed token for a specific trade agreement to designated receiver
    * if the trade has been fulfilled on the Fabric ledger
    *
    * @param _chainId Identifier of the fabric chain on which the trigger has occured
    * @param _channelId specific channel of Fabric chain in which tx should have occured
    * @param _key Key in ledger that contains details of tx
    */
    function retrieveAndExecute(bytes32 _chainId, string _channelId, string _key) public {
        uint blockVersion;
        uint txVersion;
        string memory value;
        
        (blockVersion, txVersion, value) = blockStore.getState(_chainId, _channelId, _key);
        
        // Retrieve ledger details
        RLP.RLPItem[] memory rlpValue = bytes(value).toRLPItem().toList();
        RLP.RLPItem[] memory trade = rlpValue[0].toBytes().toRLPItem().toList();

        // Level deeper
        RLP.RLPItem[] memory next = trade[0].toBytes().toRLPItem().toList();

        // Retrieve trade agreement with unique reference
        bytes32 tradeRef = keccak256(abi.encodePacked(next[0].toAscii()));

        bool result = verifyTradeDetails(tradeRef, next[1].toBytes());

        if (result) {
            // Transfer funds to recipient
            sharesettle_currency.transfer(m_trades[tradeRef].recv, m_trades[tradeRef].total, tradeRef);
            emit ShareTransfer(
                m_trades[tradeRef].send,
                m_trades[tradeRef].recv,
                m_trades[tradeRef].org,
                m_trades[tradeRef].amount,
                m_trades[tradeRef].price,
                tradeRef
            );
            
            // Mark trade as settled
            m_settled_trades[tradeRef] = true;
        }

    }

    function verifyTradeDetails(bytes32 tradeRef, bytes tradeDetails) internal returns (bool) {
        // First verify trade exists as agreement
        assert(m_opened_trades[tradeRef]);
        Trade memory trade = m_trades[tradeRef];

        // Decode RLP encoded transfer details
        RLP.RLPItem[] memory tradeTx = tradeDetails.toRLPItem().toList();

        // Verify trade details match
        assert(keccak256(trade.org)==keccak256(tradeTx[0].toAscii()));
        assert(trade.recv==tradeTx[1].toAddress());
        assert(trade.send==tradeTx[2].toAddress());
        assert(trade.amount==tradeTx[3].toUint());
        assert(trade.price==tradeTx[4].toUint());

        return true;
    }

}