pragma solidity ^0.4.24;

import "smartraiden/TokenNetwork.sol";

contract TokenNetworkInternalsTest is TokenNetwork {
    constructor (address _token_address, address _secret_registry, uint256 _chain_id)
    TokenNetwork(_token_address, _secret_registry, _chain_id)
    public
    {

    }

    function get_max_safe_uint256() pure public returns (uint256) {
        return uint256(0 - 1);
    }

    function computeMerkleRootPublic(bytes32 lockhash, bytes merkle_proof)
    pure
    public
    returns (bytes32)
    {
        return computeMerkleRoot(lockhash, merkle_proof);
    }

    function recoverAddressFromCooperativeSettleSignaturePublic(
        bytes32 channel_identifier,
        address participant1,
        uint256 participant1_balance,
        address participant2,
        uint256 participant2_balance,
        uint64 open_blocknumber,
        bytes signature
    )
    view
    public
    returns (address signature_address)
    {
        return recoverAddressFromCooperativeSettleSignature(
            channel_identifier,
            participant1,
            participant1_balance,
            participant2,
            participant2_balance,
            open_blocknumber,
            signature
        );
    }
}
