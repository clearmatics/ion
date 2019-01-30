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
    ERC223 ionlock_currency;

    uint256 ionlock_balance;

    event State(uint blockNo, uint txNo, string value);
    event StringEvent(string value);
    event BytesEvent(bytes value);
    event Bytes32Event(bytes32 value);
    event AddressEvent(address value);
    event UintEvent(uint value);
    event BoolEvent(bool value);

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

        ionlock_currency = _currency;
        blockStore = FabricStore(_storeAddr);
    }

/*
========================================================================================================================

    Payment Functions

========================================================================================================================
*/

    /**
    * Creates a new trade agreement for a cross-chain PvP payment between two counterparties
    *
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
        assert(msg.sender==address(ionlock_currency));
        assert(_value>0 && _value==m_trades[_ref].total);
        assert((ionlock_balance+_value)>ionlock_balance);
        assert(_from==m_trades[_ref].send || _from==m_trades[_ref].recv);
        assert(_value==m_trades[_ref].total);

        ionlock_balance += _value;

        // if (_from==m_trades[_ref].send) {emit StringEvent(
        //     emit IonTransfer(
        //         m_trades[_ref].send, 
        //         m_trades[_ref].recv,
        //         address(this),
        //         m_trades[_ref].amount,
        //         m_trades[_ref].price,
        //         _ref
        //     );
        // } 
        // if (_from==m_trades[_ref].recv) {
        //     // Assert receiver has enough funds to perform transfer
        //     // assert(ionlock_currency.balanceOf(msg.sender)>=m_trades[_ref].valueRecv);

        //     // settle(_ref);
        // }

    }

    function execute(uint _blockNo, uint _txNo, string _value) internal {
        emit State(_blockNo, _txNo, _value);
    }

/*
========================================================================================================================

    Verification Function

========================================================================================================================
*/

    function retrieveAndExecute(bytes32 _chainId, string _channelId, string _key) public {
        uint blockVersion;
        uint txVersion;
        string memory value;
        
        (blockVersion, txVersion, value) = blockStore.getState(_chainId, _channelId, _key);
        
        // Retrieve ledger details
        RLP.RLPItem[] memory rlpValue = bytes(value).toRLPItem().toList();
        RLP.RLPItem[] memory trade = rlpValue[0].toBytes().toRLPItem().toList();
        emit BytesEvent(trade[0].toBytes());

        // Level deeper
        RLP.RLPItem[] memory next = trade[0].toBytes().toRLPItem().toList();

        // Retrieve trade agreement with unique reference
        bytes32 tradeRef = keccak256(abi.encodePacked(next[0].toAscii()));
        emit Bytes32Event(tradeRef);

        bool result = verifyTradeDetails(tradeRef, next[1].toBytes());

        emit BoolEvent(result);
        // emit test2(next[1].toBytes());
        // emit test2(next[2].toBytes());

        // retrieveBalance(next[2].toBytes());
        // emit test(value);
        

        // RLP.RLPItem[] memory next = item1.toRLPItem().toList();
        // bytes memory item2 = next[0].toBytes();
        // // emit test2(item2);

        // RLP.RLPItem[] memory next1 = item2.toRLPItem().toList();
        // bytes memory item3 = next1[0].toBytes();
        // // emit test2(item3);
        // bytes memory item4 = next1[1].toBytes();
        // // emit test2(item4);

        // bytes memory item5 = next1[2].toBytes();
        // emit test2(item5);

        // RLP.RLPItem[] memory transaction2 = item5.toRLPItem().toList();
        // bytes memory item6 = transaction2[0].toBytes();
        // emit test2(item6);

        // item4 = next1[2].toBytes();
        // RLP.RLPItem[] memory next2 = item4.toRLPItem().toList();
        // bytes memory item4 = next2[0].toBytes();
        // emit test2(item4);

        // bytes memory item2 = item1.toBytes();
        // emit test2(item2);
        // bytes memory item3 = transaction[2].toBytes();
        // emit test2(item3);
        // emit test2(key2);

    }

    function verifyTradeDetails(bytes32 tradeRef, bytes tradeDetails) internal returns (bool) {
        // First verify trade exists as agreement
        assert(m_opened_trades[tradeRef]);
        Trade memory trade = m_trades[tradeRef];

        // Decode RLP encoded transfer details
        RLP.RLPItem[] memory tradeTx = tradeDetails.toRLPItem().toList();

        // Verify trade details match
        assert(keccak256(trade.org)==keccak256(tradeTx[0].toAscii()));
        emit AddressEvent(tradeTx[2].toAddress());
        // assert(trade.recv==tradeTx[2].toAddress());

        return true;
    }

    function retrieveBalance(bytes data) internal {
        RLP.RLPItem[] memory ledger = data.toRLPItem().toList();
        RLP.RLPItem[] memory balance = ledger[0].toBytes().toRLPItem().toList();
        emit AddressEvent(balance[0].toAddress());
        emit UintEvent(balance[1].toUint());

    }

    // Convert an hexadecimal string to raw bytes
    // function fromHex(string s) public pure returns (bytes) {
    //     bytes memory ss = bytes(s);
    //     require(ss.length%2 == 0); // length must be even
    //     bytes memory r = new bytes(ss.length/2);
    //     for (uint i=0; i<ss.length/2; ++i) {
    //         r[i] = byte(fromHexChar(uint(ss[2*i])) * 16 +
    //                     fromHexChar(uint(ss[2*i+1])));
    //     }
    //     return r;
    // }

}