pragma solidity ^0.4.18;

import "./Merkle.sol";
import "./Sodium_Interface.sol";

contract Sodium is Sodium_Interface
{
	mapping(uint256 => uint256) internal m_roots;
	uint256 LatestBlock;
	address Owner;


	function Sodium ()
		public
	{
		Owner = msg.sender;
		LatestBlock = block.number - (block.number % GroupSize());
	}


	function Destroy ()
		public
	{
		require( msg.sender == Owner );
		selfdestruct( msg.sender );
	}


	function GroupSize ()
		public pure
		returns (uint256)
	{
		// Kovan, 5 second block times
		// One merkle root per minute
		return 12;
	}


	function NextBlock ()
		public view
		returns (uint256)
	{
		return LatestBlock + GroupSize();
	}


	function GetMerkleRoot (uint256 block_no)
		public view
		returns (uint256)
	{
		return m_roots[ block_no - (block_no % GroupSize()) ];
	}


	function Update( uint256 start_block, uint256[] merkle_roots )
		public
	{
		bool success = false;
		uint256 latest_block;

		require( start_block == NextBlock() );

		require( msg.sender == Owner );

		for( uint256 i = 0; i < merkle_roots.length; i++ )
		{
	    // XXX: prevent overflow
			uint256 block_no = start_block + (i * GroupSize());

			if( m_roots[block_no] == 0 )
			{
				success = true;

				m_roots[block_no] = merkle_roots[i];
			}

			latest_block = block_no;
		}

		require( success );

		LatestBlock = latest_block;
	}


	function Verify( uint256 block_no, uint256 leaf_hash, uint256[] proof )
		public view
		returns (bool)
	{
		uint256 merkle_root = GetMerkleRoot(block_no);

		require( uint256(merkle_root) != 0 );

		return Merkle.Verify( merkle_root, leaf_hash, proof );
	}
}
