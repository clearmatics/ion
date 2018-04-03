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

  /**
  * Supplies a sequence of merkle roots which create a hash-chain
  *
  *   hash = H(hash, root)
  */
	function Update( uint256 _new_block_root )
		public
	{
		require( msg.sender == Owner);
		uint256 prev_hash = LatestBlock;
		uint256 new_block_hash;

		if (prev_hash != 0) {
			new_block_hash = uint256(keccak256(prev_hash, _new_block_root));
		} else {
			new_block_hash = uint256(keccak256(_new_block_root));
		}

		IonBlock storage blk = GetBlock(new_block_hash);

		blk.root = _new_block_root;

		// Record state at time of block creation
		blk.prev = prev_hash;
		blk.time = block.timestamp;

		LatestBlock = new_block_hash;
	}


	function Verify( uint256 block_id, uint256 leaf_hash, uint256[] proof )
		public view
		returns (bool)
	{
		return Merkle.Verify( GetRoot(block_id), leaf_hash, proof );
	}
}
