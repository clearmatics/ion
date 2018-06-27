// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "./ECVerify.sol";

contract Validation {
	address Owner;
	address[] validators;

	bytes32 prevBlockHash;
	bytes32 parentBlockHash;

	struct BlockHeader {
		bytes32 prevBlockHash;
	}

	mapping (bytes32 => BlockHeader) public m_blockheaders;
	mapping (address => bool) m_validators;

	event broadcastSig(address owner);
	event broadcastHashData(bytes header, bytes parentHash, bytes rootHash);
	event broadcastHash(bytes32 blockHash);

	/*
	*	@param _validators			list of validators at block 0
	*/
	constructor (address[] _validators, bytes32 genHash) public {
		Owner = msg.sender;
		for (uint i = 0; i < _validators.length; i++) {
			validators.push(_validators[i]);
			m_validators[_validators[i]] = true;
    	}

		prevBlockHash = genHash;	
	}

	/*
	* Returns the validators array
	*/
	function GetValidators() public view returns (address[] _validators) {
		return validators;
	}
	
	/*
	* Returns the latest block submitted
	*/
	function LatestBlock() public view returns (bytes32 _latestBlock) {
		return prevBlockHash;
	}

	/*
	* @param header  			header rlp encoded, with extraData signatures removed
	* @param prefixHeader		the new prefix for the signed hash header
	* @param prefixExtraData	the new prefix for the extraData field
	*/
	function ValidateBlock(bytes header, bytes prefixHeader, bytes prefixExtraData) public {
		uint256 length = header.length;
		bytes32 blockHash = keccak256(header);

		emit broadcastHash(blockHash);

		bytes memory headerStart 	= new bytes(length - 141);
		bytes memory extraData 		= new bytes(31);
		bytes memory extraDataSig	= new bytes(65);
		bytes memory headerEnd 		= new bytes(42);

		// Extract the start of the header and replace the length
		extractData(headerStart, header, 0, headerStart.length);
		assembly {
           let ret := staticcall(3000, 4, add(prefixHeader, 32), 2, add(headerStart, 33), 2)
	    }

		// Extract the real extra data and create the signed hash
		extractData(extraData, header, length-140, extraData.length);
		assembly {
			let ret := staticcall(3000, 4, add(prefixExtraData, 32), 1, add(extraData, 32), 1)
		}

		// Extract the end of the header
		extractData(headerEnd, header, length-42, headerEnd.length);
		bytes memory newHeader = mergeHash(headerStart, extraData, headerEnd);

		bytes32 hashData = keccak256(newHeader);

		// Extract the signature of the hash create above
		extractData(extraDataSig, header, length-107, extraDataSig.length);

		address sig_addr = ECVerify.ecrecovery(hashData, extraDataSig);
		require(m_validators[sig_addr]==true, "Signer not a validator!");

		prevBlockHash = blockHash;

		emit broadcastSig(sig_addr);

	}

	function mergeHash(bytes headerStart, bytes extraData, bytes headerEnd) internal view returns (bytes output) {
		// Get the lengths sorted because they're needed later...
		uint256 headerStartLength = headerStart.length;
		uint256 extraDataLength = extraData.length;
		uint256 extraDataStart = headerStartLength + 32;
		uint256 headerEndLength = headerEnd.length;
		uint256 headerEndStart = extraDataLength + headerStartLength + 32 + 2;
		uint256 newLength = headerStartLength + extraDataLength + headerEndLength + 2; // extra two is for the prefix
		bytes memory header = new bytes(newLength);


		// Add in the first part of the header
		assembly {
			let ret := staticcall(3000, 4, add(headerStart, 32), headerStartLength, add(header, 32), headerStartLength)
		}
		assembly {
			let ret := staticcall(3000, 4, add(extraData, 32), extraDataLength, add(header, extraDataStart), extraDataLength)
		}
		assembly {
			let ret := staticcall(3000, 4, add(headerEnd, 32), headerEndLength, add(header, headerEndStart), headerEndLength)
		}

		output = header;
	}

	/*
	* @param data	memory allocation for the data you need to extract
	* @param sig    array from which the data should be extracted
	* @param start  index which the data starts within the byte array
	* @param length total length of the data to be extracted
	*/
	function extractData(bytes data, bytes input, uint start, uint length) private pure {
		for (uint i=0; i<length; i++) {
			data[i] = input[start+i];
		}
	}

}
