pragma solidity ^0.4.18;

import "./Sodium_Interface.sol";
import "./ECVerify.sol";
import "./ERC223Compatible.sol";

/*

B can only deposit if he has proof of A's deposit
Either A or B can cancel B's deposit
A can only cancel A's deposit if she has proof that B's deposit cancelled
B can only withdraw A's deposit if he has proof that A withdrew B's deposit


keccak256("Deposit(address,hash256,uint256)")
*/

contract Fluoride is ERC223ReceivingContract
{
	enum State
	{
		Invalid,	// TODO: verify that State.Invalid == 0s
		StartedA,
		StartedB,
		DepositedA,
		DepositedB
	}

	struct Data
	{
		State state;
		address owner;
		address token;
		uint256 amount;

		// This will be the address that the merkle proof must be supplied from
		address counterparty_contract;

		// Transfer to counterparty upon withdraw, on same chain as owner
		address counterparty;

		uint expire;
	}


	event OnDeposit( bytes32 trade_id );

	event OnCancel( bytes32 trade_id );

	event OnWithdraw( bytes32 trade_id );

	event OnTimeout( bytes32 trade_id );


	Sodium_Interface internal m_sodium;

	mapping( bytes32 => Data ) m_exchanges;


	function Fluoride( Sodium_Interface sodium_address )
		public
	{
		m_sodium = sodium_address;
	}

  function bytesToBytes32(bytes b, uint offset) private pure returns (bytes32) {
    bytes32 out;

    for (uint i = 0; i < 32; i++) {
      out |= bytes32(b[offset + i] & 0xFF) >> (i * 8);
    }
    return out;
  }

  event OnTokenTransfer(address _from, uint _value, bytes _data, bytes32 trade_id, string side, address owner);
	function tokenFallback(address _from, uint _value, bytes _data)
		public
	{
		// Load _data bytes into trade_id
		bytes32 trade_id;
		require( _data.length == 32 );
    /*
		assembly {
			//let ret := staticcall(3000, 4, add(_data, 32), 32, trade_id, 32)
			//switch ret case 0 { invalid }
      //let x := mload(0x40)
      //mstore(x,)
      //trade_id := x
			//let ret := staticcall(3000, 4, x, 32, trade_id, 32)
			//switch ret case 0 { invalid }
      //codecopy(trade_id,_data,32)
      calldatacopy(trade_id,calldata(_data),32)
		}
    */
    trade_id = bytesToBytes32(_data,0);

		// Load exchange, must be in correct state
		Data storage trade = m_exchanges[trade_id];

		State state = trade.state;
		require( state == State.StartedA || state == State.StartedB );

		// XXX: is it necessary to check _from when tx.origin is also checked?
		// Is it necessary to check either?
		require( trade.owner == tx.origin );
		require( _from == trade.owner );

		// The value *must* be checked
		require( _value == trade.amount );
		// As must the ERC-223 token
		require( trade.token == msg.sender );

		// Transition from Started to Deposited after recieving the funds
		if( state == State.StartedA ) {
			trade.state = State.DepositedA;
		}
		else if( state == State.StartedB ) {
			trade.state = State.DepositedB;
		}
    OnTokenTransfer( _from, _value, _data, trade_id, "OUT",trade.owner);
	}


	function getcodesize(address _addr)
		internal constant
		returns (uint256 outsz)
	{
        assembly {
            outsz := extcodesize(_addr)
        }
    }


	function VerifyTradeAgreement( bytes32 a_hash, bytes a_sig, bytes32 b_hash, bytes b_sig, bytes c_sig )
		internal constant
		returns (bytes32, address, address)
	{
		// Initiator sends signed message to Counterparty
		var a_addr = ECVerify.ecrecovery(a_hash, a_sig);

		// Counterparty confirms signed message from Initiator
		// b_hash includes fingerprint of all info
		var ab_hash = keccak256(a_hash, a_addr, b_hash);
		var b_addr = ECVerify.ecrecovery(ab_hash, b_sig);

		// Closer accepts Counterparty offer
		var abc_hash = keccak256(ab_hash, b_addr);
		var c_addr = ECVerify.ecrecovery(abc_hash, c_sig);

		var trade_id = keccak256(abc_hash, c_sig);

		require( c_addr == a_addr );

		// Initiator and Closer must be the same
		return (trade_id, a_addr, b_addr);
	}


	function Start_OnAbyA( address a_contract, uint a_expire, address a_token, uint256 a_amount, bytes a_sig, address b_contract, bytes32 b_state, bytes b_sig, bytes c_sig )
		public returns (bytes32)
	{
		require( a_contract == address(this) );
		require( getcodesize(a_token) > 0 );

		bytes32 trade_id;
		address a_addr;
		address b_addr;
		(trade_id, a_addr, b_addr) = VerifyTradeAgreement(
			keccak256(a_contract, keccak256(a_expire, a_token, a_amount)),
			a_sig,
			keccak256(b_contract, b_state),
			b_sig,
			c_sig );

		m_exchanges[trade_id] = Data({
			state: State.StartedA,
			owner: msg.sender,
			token: a_token,
			amount: a_amount,
			expire: a_expire,
			counterparty_contract: b_contract,
			counterparty: b_addr
		});

		OnDeposit(trade_id);
	}


	/**
	* Deposit by Bob on Bob's chain
	*
	* He must have verified Alice has deposited first, but this isn't 100% necessary as they can
	* simultaneously perform their deposits if they trust each other.
	*/
	function Start_OnBbyB( address a_contract, bytes32 a_state, bytes a_sig, address b_contract, uint256 b_expire, address b_token, uint256 b_amount, bytes b_sig, bytes c_sig )
		public returns (bytes32)
	{
		require( b_contract == address(this) );
		require( getcodesize(b_token) > 0 );

		bytes32 trade_id;
		address a_addr;
		address b_addr;
		(trade_id, a_addr, b_addr) = VerifyTradeAgreement(
			keccak256(a_contract, a_state),
			a_sig,
			keccak256(b_contract, keccak256(b_expire, b_token, b_amount)),
			b_sig,
			c_sig );

		m_exchanges[trade_id] = Data({
			state: State.StartedB,
			owner: msg.sender,
			token: ERC223(b_token),
			amount: b_amount,
			expire: b_expire,
			counterparty_contract: a_contract,
			counterparty: a_addr
		});

		OnDeposit(trade_id);
	}


	/**
	* Bob cancels his own deposit on his own chain
	*/
	function Cancel( bytes32 trade_id )
		public
	{
		Data storage trade = m_exchanges[trade_id];

		// XXX: This can only be cancelled on A chain if she supplies proof!

		require( trade.state == State.StartedB );
		require( trade.expire <= block.timestamp );
		require( trade.owner == msg.sender || trade.counterparty == msg.sender );

		delete m_exchanges[trade_id];

		OnCancel( trade_id );

		// If funds have been deposited, return them to owner
		State state = trade.state;
		if( state == State.DepositedB || state == State.DepositedA )
		{
			ERC223(trade.token).transfer(trade.owner, trade.amount);
		}
	}


  // TODO: THE WITHDRAW FUNCTION IS NOT FINiSHED!
  // IT DOES NOT  VERIFY IF ANY OF THE SIGNATURES WERE CREATED OR ANYTHING
  // THE FIRST REQUIRE CHECKS IF A HAS STARTED, IT SHOULD CHECK IF A DEPOSITED
	/**
	* Withdraw by Bob on Alice's chain
	*
	* This can be triggered by anybody that can supply the proof
	*/
	function Withdraw( bytes32 trade_id, uint256 block_no, uint256[] proof )
		public
	{
		Data storage trade = m_exchanges[trade_id];
		//require( trade.state == State.StartedA ); // TODO: Does this check need to exist?!
		require( trade.state == State.DepositedB || trade.state == State.DepositedA );
		require( trade.expire > block.timestamp );

		// TODO: create hash of event to expect from other contract
		// TODO: insert the Topic name between contract and trade id
		var expect_event = keccak256(trade.counterparty_contract, /* topic id, unknown yet */ trade_id);

		require( m_sodium.Verify( block_no, uint256(expect_event), proof ) );

		//delete m_exchanges[trade_id]; // TODO: delete failing!

		OnWithdraw(trade_id);

		ERC223(trade.token).transfer(trade.counterparty, trade.amount);
	}
}
