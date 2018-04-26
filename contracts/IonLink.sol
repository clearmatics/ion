// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.18;

import "./Merkle.sol";
import "./IonCompatible.sol";

contract IonLink is IonLinkInterface
{
	struct IonBlock
	{
	    uint256 root;
	    uint256 prev;
	    uint256 time;
	}

	mapping(uint256 => IonBlock) internal m_blocks;

	uint256 public LatestBlock;

	address Owner;

	event IonLinkUpdated();

	function IonLink ( uint256 genesis )
		public
	{
		Owner = msg.sender;
		LatestBlock = genesis;
	}


	function Destroy ()
		public
	{
		require( msg.sender == Owner );

		selfdestruct( msg.sender );
	}


	function GetBlock( uint256 block_id )
	    internal view returns (IonBlock storage)
	{
	    IonBlock storage blk = m_blocks[block_id];

	    return blk;
	}


    function GetTime( uint256 block_id )
	    public view returns (uint256)
	{
	    return GetBlock(block_id).time;
	}


    function GetPrevious( uint256 block_id )
	    public view returns (uint256)
	{
	    return GetBlock(block_id).prev;
	}


	function GetRoot( uint256 block_id )
	    public view returns (uint256)
	{
	    return GetBlock(block_id).root;
	}

	function GetLatestBlock()
	    public view returns (uint256)
	{
	    return LatestBlock;
	}

	/**
	* Supplies a sequence of merkle roots which create a hash-chain
	*
	*   hash = H(hash, root)
	*/
	function Update( uint256[] in_state )
		public
	{
		require( in_state.length > 1 );

		uint256 prev_hash = LatestBlock;

		for( uint256 i = 0; i < in_state.length; i++ )
		{
			uint256 block_hash = uint256(keccak256(prev_hash, in_state[i]));

			IonBlock storage blk = m_blocks[block_hash];

			blk.root = in_state[i];

			// Record state at time of block creation
			blk.prev = prev_hash;
			blk.time = block.timestamp;

			prev_hash = block_hash;
		}

		LatestBlock = prev_hash;

		emit IonLinkUpdated();
	}


	function Verify( uint256 block_id, uint256 leaf_hash, uint256[] proof )
		public view
		returns (bool)
	{
		return Merkle.Verify( GetRoot(block_id), leaf_hash, proof );
	}
}
