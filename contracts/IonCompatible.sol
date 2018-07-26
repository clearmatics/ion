pragma solidity ^0.4.0;
import "./Ion.sol";

contract IonCompatible {
    /*  The Ion contract that proofs would be made to. Ensure that prior to verification attempts that the relevant
        blocks have been submitted to the Ion contract. */
    Ion internal ion;

    constructor(address ionAddr) public {
        ion = Ion(ionAddr);
    }
}
