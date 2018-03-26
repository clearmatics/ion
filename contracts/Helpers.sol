pragma solidity ^0.4.18;

/**
* Keccak256 Generator for use in testing...
*/
contract Helpers
{

	function a_Keccak256 (string a_input)
		pure public returns (bytes32)
	{
		var a_hash  = keccak256(a_input);

		return a_hash;
	}
}
