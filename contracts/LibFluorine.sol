pragma solidity ^0.4.18;

library LibFluorine {
	function HashTx (uint256 block_no, address tx_from, address tx_to, uint256 tx_value, bytes32 tx_input_hash )
	    internal pure returns (bytes32)
	{
		return keccak256(block_no, tx_from, tx_to, tx_value, tx_input_hash);
	}

	function HashEvent (uint256 block_no, address contract_addr, bytes32 topic, bytes32 data_hash )
	    internal pure returns (bytes32)
	{
		return keccak256(block_no, contract_addr, topic, data_hash);
	}
}

