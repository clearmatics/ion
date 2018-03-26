pragma solidity ^0.4.18;

/**
* Hash Time Lock Contract
*/
contract HTLC
{

	event OnDeposit( uint256 lock_id, address receiver );
	event OnClaim( uint256 lock_id, bytes32 preimage, uint256 timeout );
	event OnRecover( address receiver, address lockreceiver, address msgsender );
	event OnVerification( address verified );
	event Test( uint256 lock_id, uint256 timeout );

	uint256 constant MAX_TIMEOUT = 60 * 60 * 24;

	struct LockState {
		uint256 timeout;
		bytes32 hash;
		address receiver;
		address owner;
		uint256 value;
	}
	mapping(uint256 => LockState) internal m_locks;

	uint256 internal m_ctr;

  modifier fundsSent() {
      require(msg.value > 0);
      _;
  }
  modifier futureTimelock(uint _time) {
      // only requirement is the timelock time is after the last blocktime (now).
      // probably want something a bit further in the future then this.
      // but this is still a useful sanity check:
      require(_time > now);
      _;
  }

	/**
	* Alice deposits her funds (the Input, or side A)
	*
	* These can only be withdrawn by the intended `receiver` if they provide
	* the `preimage` which results in `hash`.
	*
	* On side `B` alice reveals her secret `preimage`
	*
	* If the timeout has been reached Alice (the `owner`) can request a Refund.
	*/
	function Deposit (uint256 timeout, bytes32 hash, address receiver)
		external
		payable
		fundsSent
		futureTimelock(timeout)
		returns (uint256)
	{
		require( timeout >= now );
		require( timeout < (now + MAX_TIMEOUT) );
		require( msg.value > 0 );

		var lock_id = m_ctr;
		m_ctr += 1;
		m_locks[lock_id] = LockState(timeout, hash, receiver, msg.sender, msg.value);

		OnDeposit(lock_id, receiver);

		return lock_id;
	}


	function Claim (uint256 lock_id, bytes32 preimage, uint8 v, bytes32 r, bytes32 s)
		public
	{
		/* bytes32 result = sha256("randomhash"); */
		LockState storage lock = m_locks[lock_id];
		require( lock.timeout >= now );
		require( sha256(preimage) == lock.hash );

    // Only receiver can provide the preimage to withdraw
		bytes memory prefix = "\x19Ethereum Signed Message:\n32";
		bytes32 prefixedHash = keccak256(prefix, preimage);
		address receiver = ecrecover(prefixedHash, v, r, s);
		OnVerification(receiver);
		require( receiver == lock.receiver );

		var value = lock.value;
		delete m_locks[lock_id];

		msg.sender.transfer(value);

		OnClaim(lock_id, preimage, lock.timeout);
	}


	function Refund (uint256 lock_id,  bytes32 preimage, uint8 v, bytes32 r, bytes32 s)
		public
	{
		LockState storage lock = m_locks[lock_id];
		require( lock.timeout < now );
		/* Test(lock.timeout, now); */

		// Lock owner must authorise msg.sender to retrieve refund
		bytes memory prefix = "\x19Ethereum Signed Message:\n32";
		bytes32 prefixedHash = keccak256(prefix, preimage);
		address owner = ecrecover(prefixedHash, v, r, s);
		OnVerification(owner);
		require( owner == lock.owner );

		var value = lock.value;
		delete m_locks[lock_id];

		msg.sender.transfer(value);
	}


	function Verify(bytes32 hash, uint8 v, bytes32 r, bytes32 s)
	public
	{
		bytes memory prefix = "\x19Ethereum Signed Message:\n32";
		bytes32 prefixedHash = keccak256(prefix, hash);
		address receiver = ecrecover(prefixedHash, v, r, s);

		OnVerification(receiver);
	}

}
