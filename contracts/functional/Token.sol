// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "../libraries/ERC223Compatible.sol";
import '../libraries/SafeMath.sol';

/**
 * @title Reference implementation of the ERC223 standard token.
 */
contract Token is ERC223 {
    using SafeMath for uint256;

    mapping(address => uint256) balances; // List of user balances.

    event AccountTransfer();
    event test(address to);

    constructor() public {
        totalSupply = 0;
    }


    function mint(uint256 _value) public {
        totalSupply = totalSupply.add(_value);

        balances[msg.sender] = balances[msg.sender].add(_value);
    }


    function burn(uint256 _value) public {
        balances[msg.sender] = balances[msg.sender].sub(_value);

        totalSupply = totalSupply.sub(_value);
    }

    function metadataTransfer(address _to, uint256 _value, bytes32 _ref) public {
        transfer(_to, _value, _ref);
    }

    function rawTransfer(address _to, uint256 _value) public {
        transfer(_to, _value);
    }

    /**
     * @dev Transfer the specified amount of tokens to the specified address.
     *      Invokes the `tokenFallback` function if the recipient is a contract.
     *      The token transfer fails if the recipient is a contract
     *      but does not implement the `tokenFallback` function
     *      or the fallback function to receive funds.
     *
     * @param _to    Receiver address.
     * @param _value Amount of tokens that will be transferred.
     * @param _ref  Transaction metadata.
     */
    function transfer(address _to, uint256 _value, bytes32 _ref) public {
        balances[msg.sender] = balances[msg.sender].sub(_value);
        balances[_to] = balances[_to].add(_value);

        // Standard function transfer similar to ERC20 transfer with no _ref .
        // Added due to backwards compatibility reasons .
        uint256 codeLength;

        assembly {
            // Retrieve the size of the code on target address, this needs assembly .
            codeLength := extcodesize(_to)
        }

        if(codeLength>0) {
            ERC223ReceivingContract receiver = ERC223ReceivingContract(_to);
            receiver.tokenFallback(msg.sender, _value, _ref);
        } else {
            emit AccountTransfer();
        }

        emit Transfer(msg.sender, _to, _value, _ref);

    }


    /**
     * @dev Transfer the specified amount of tokens to the specified address.
     *      This function works the same with the previous one
     *      but doesn't contain `_data` param.
     *      Added due to backwards compatibility reasons.
     *
     * @param _to    Receiver address.
     * @param _value Amount of tokens that will be transferred.
     */
    function transfer(address _to, uint256 _value) public {
        bytes32 empty;
        return transfer(_to, _value, empty);
    }


    /**
     * @dev Returns balance of the `_owner`.
     *
     * @param _owner   The address whose balance will be returned.
     * @return balance Balance of the `_owner`.
     */
    function balanceOf(address _owner) public view returns (uint256) {
        return balances[_owner];
    }
}