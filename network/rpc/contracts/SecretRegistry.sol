pragma solidity ^0.4.23;

/// @title contract for registering secret to unlock transferred transactions.
contract SecretRegistry {

    /** A string literal denoting the version of this contract. */
    string constant public contract_version = "0.3._";

    // secrethash => block number at which the secret was revealed
    mapping(bytes32 => uint256) public secrethash_to_block;


    /** a secret is registered */
    event SecretRevealed(bytes32 indexed secret);

    /// @notice Registers a hash time lock secret and saves the block number.
    /// This allows the lock to be unlocked after the expiration block.
    /// @param secret The secret used to lock the hash time lock.
    /// @return true if secret was registered, false if the secret was already registered.
    function registerSecret(bytes32 secret) public returns (bool) {
        bytes32 secrethash = sha256(abi.encodePacked(secret));
        //secret already registered
        if (secret == 0x0 || secrethash_to_block[secrethash] > 0) {
            revert();
        }
        secrethash_to_block[secrethash] = block.number;
        emit SecretRevealed(secret);
        return true;
    }

    /// @notice function to return the block height corresponding to secret hash.
    /// @param secrethash                   a secret hash to be revealed.
    /// @return a 256-bit unsigned integer  block number in which this secrect was revealed.
    function getSecretRevealBlockHeight(bytes32 secrethash) public view returns (uint256) {
        return secrethash_to_block[secrethash];
    }
}
