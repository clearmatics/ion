pragma solidity ^0.4.0;
import "./Ion.sol";

contract IonCompatible {
    Ion ion;

    constructor(address ionAddr) public {
        ion = Ion(ionAddr);
    }
}