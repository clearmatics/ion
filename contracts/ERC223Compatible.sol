pragma solidity ^0.4.18;

import "./ERC223.sol";

contract ERC223ReceivingContract
{
    function tokenFallback(address _from, uint _value, bytes _data) public;
}
