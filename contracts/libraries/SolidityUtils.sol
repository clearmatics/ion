// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

// A library of funky data manipulation stuff
library SolUtils {
	/*
    * @description  copies 32 bytes from input into the output
	* @param output	memory allocation for the data you need to extract
	* @param input  array from which the data should be extracted
	* @param buf	index which the data starts within the byte array needs to have 32 bytes appended
	*/
	function BytesToBytes32(bytes input, uint256 buf) internal pure returns (bytes32 output) {
		buf = buf + 32;
        assembly {
			output := mload(add(input, buf))
		}
	}

	/*
    * @description  copies 20 bytes from input into the output
	* @param output	memory allocation for the data you need to extract
	* @param input  array from which the data should be extracted
	* @param buf	index which the data starts within the byte array needs to have 32 bytes appended
	*/
	function BytesToBytes20(bytes input, uint256 buf) internal pure returns (bytes20) {
        bytes20 output;

        for (uint i = 0; i < 20; i++) {
            output |= bytes20(input[buf + i] & 0xFF) >> (i * 8);
        }
        return output;
    }

/*
    * @description  copies 20 bytes from input into the output returning an address
	* @param output	memory allocation for the data you need to extract
	* @param input  array from which the data should be extracted
	* @param buf	index which the data starts within the byte array needs to have 32 bytes appended
	*/
	function BytesToAddress(bytes input, uint256 buf) internal pure returns (address output) {
		buf = buf + 20;
		assembly {
			output := mload(add(input, buf))
		} 
	}

	/*
    * @description  copies output.length bytes from the input into the output
	* @param output	memory allocation for the data you need to extract
	* @param input  array from which the data should be extracted
	* @param buf	index which the data starts within the byte array
	*/
	function BytesToBytes(bytes output, bytes input, uint256 buf) constant internal {
		uint256 outputLength = output.length;
		buf = buf + 32; // Append 32 as we need to point past the variable type definition
		assembly {
           let ret := staticcall(3000, 4, add(input, buf), outputLength, add(output, 32), outputLength)
	    }
	}

	function UintToString(uint _i) internal pure returns (string memory _uintAsString) {
		if (_i == 0) {
			return "0";
		}
		uint j = _i;
		uint len;
		while (j != 0) {
			len++;
			j /= 10;
		}
		bytes memory bstr = new bytes(len);
		uint k = len - 1;
		while (_i != 0) {
			bstr[k--] = byte(uint8(48 + _i % 10));
			_i /= 10;
		}
		return string(bstr);
	}

	function BoolToString(bool _b) internal pure returns (string memory) {
		if (_b)
            return "true";
        else
            return "false";
	}

}
