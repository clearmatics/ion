// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.18;

import "./ERC223Compatible.sol";
import "./IonCompatible.sol";


contract IonLock is ERC223ReceivingContract, IonCompatible{
    uint256 m_balance;

    ERC223 m_currency;

    IonLinkInterface m_ion;

    // Logs the unique references that have been submitted
    mapping(bytes32 => bool) m_withdraws;

    // Keeps reference to latest block that transfer was performed on
    uint256 public LatestBlock;

    constructor( ERC223 currency, IonLinkInterface ion ) public {
        require( address(currency) != 0 );

        m_ion = ion;

        m_currency = currency;

        LatestBlock = block.number;
    }


    /**
    * When ERC223 tokens are sent to this contract it mints the
    * equivalent value of its own currency and gives it to the
    * sender using Ion.
    *
    * @param _from Who sent us the token
    * @param _value Amount of tokens
    * @param _data Arbitrary data, to be used as the payment reference
    */
    function tokenFallback(address _from, uint _value, bytes _data) public {
        require( msg.sender == address(m_currency) );

        require( _value > 0 );

        require( (m_balance + _value) > m_balance );

        m_balance += _value;

        bytes32 ref = keccak256(_data);

        emit IonMint( _value, ref );

        emit IonTransfer( _from, address(this), _value, ref, _data );

        LatestBlock = block.number;
    }


    /**
    * When given proof that a payment exists in a proof uploaded to IonLink
    * it will allow the sender to withdraw tokens of the specified value.
    *
    * @param _value Amount of token to withdraw
    * @param _ref Payment reference hash
    * @param _block_id IonLink block ID
    * @param _proof Merkle proof
    */
    function Withdraw( uint256 _value, bytes32 _ref, uint256 _block_id, uint256[] _proof ) public {
        require( false == m_withdraws[_ref] );

        // Definition of leaf structure
        uint256 leaf_hash = uint256(keccak256(msg.sender, m_currency, address(this), _value, _ref));

        require( m_ion.Verify(_block_id, leaf_hash, _proof) );

        m_withdraws[_ref] = true;

        require( (m_balance - _value) < _value );

        m_balance -= _value;

        m_currency.transfer(msg.sender, _value);

        emit IonWithdraw(msg.sender, m_currency, _value, _ref);
    }
}
