pragma solidity ^0.4.24;

import "../IonCompatible.sol";

contract BlockStore is IonCompatible {
    bytes32[] public registeredChains;

    mapping (bytes32 => bool) public m_chains;

    modifier onlyIon() {
        require(msg.sender == address(ion), "Block does not exist for chain");
        _;
    }

    /*
    * onlyRegisteredChains
    * param: _id (bytes32) Unique id of chain supplied to function
    *
    * Modifier that checks if the provided chain id has been registered to this contract
    */
    modifier onlyRegisteredChains(bytes32 _chainId) {
        require(m_chains[_chainId], "Chain is not registered");
        _;
    }

    /*
    * Constructor
    * param: id (bytes32) Unique id to identify this chain that the contract is being deployed to.
    *
    * Supplied with a unique id to identify this chain to others that may interoperate with it.
    * The deployer must assert that the id is indeed public and that it is not already being used
    * by another chain
    */
    constructor(address _ionAddr) IonCompatible(_ionAddr) public {}

    /*
    * addChain
    * param: id (bytes32) Unique id of another chain to interoperate with
    *
    * Supplied with an id of another chain, checks if this id already exists in the known set of ids
    * and adds it to the list of known m_chains.
    *
    *Should be called by the validation registerChain() function
    */
    function addChain(bytes32 _chainId) onlyIon public returns (bool) {
        require( _chainId != ion.chainId(), "Cannot add this chain id to chain register" );
        require(!m_chains[_chainId], "Chain already exists" );

        m_chains[_chainId] = true;
        registeredChains.push(_chainId);

        return true;
    }

    function addBlock(bytes32 _chainId, bytes32 _blockHash, bytes _blockBlob) onlyIon onlyRegisteredChains(_chainId);
}
