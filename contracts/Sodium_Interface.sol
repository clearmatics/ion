pragma solidity ^0.4.18;

interface Sodium_Interface
{
	function Verify( uint256 block_no, uint256 leaf_hash, uint256[] proof )
		public view returns (bool);
}
