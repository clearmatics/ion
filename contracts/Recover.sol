// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "./libraries/ECVerify.sol";

contract Recover {
	address Owner;

	event broadcastSig(address owner);
	event broadcastHashData(bytes header, bytes parentHash, bytes rootHash);
	/* event test(bytes start, bytes data); */
	event broadcastHash(bytes32 blockHash);
	event test(bytes header);

	constructor () public {
		Owner = msg.sender;
	}

	/*
	 * @param data  			data that has been signed
	 * @param sig    			signature of data
	 */
	function VerifyHash(bytes32 data, bytes sig) public {
		address sig_addr = ECVerify.ecrecovery(data, sig);

		emit broadcastSig(sig_addr);
	}

	/*
	* @param header  			header rlp encoded, with extraData signatures removed
	* @param sig    			extraData signatures
	*/
	function VerifyBlock(bytes header, bytes sig) public {
		bytes32 hashData = keccak256(header);
		address sig_addr = ECVerify.ecrecovery(hashData, sig);

		bytes memory parentHash = new bytes(32);
		bytes memory rootHash = new bytes(32);

		// get parentHash and rootHash
		extractData(parentHash, header, 4, 32);
		extractData(rootHash, header, 91, 32);

		emit broadcastHashData(header, parentHash, rootHash);
		emit broadcastSig(sig_addr);
	}

	/*
	* @param header  					header rlp encoded, with extraData signatures removed
	* @param prefixHeader			the new prefix for the signed hash header
	* @param prefixExtraData	the new prefix for the extraData field
	*/
	function ExtractHash(bytes header, bytes prefixHeader, bytes prefixExtraData) public {
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
	* @param data	  			memory allocation for the data you need to extract
	* @param sig    			array from which the data should be extracted
	* @param start   			index which the data starts within the byte array
	* @param length  			total length of the data to be extracted
	*/
	function extractData(bytes data, bytes input, uint start, uint length) private pure {
		for (uint i=0; i<length; i++) {
			data[i] = input[start+i];
		}
	}

}
