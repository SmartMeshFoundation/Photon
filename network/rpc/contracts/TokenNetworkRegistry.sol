pragma solidity ^0.4.23;

import "./Utils.sol";
import "./Token.sol";
import "./TokenNetwork.sol";

/// @title contract to register a TokenNetwork.
/// @notice
contract TokenNetworkRegistry is Utils {

    string constant public contract_version = "0.3._";
    address public secret_registry_address;
    uint256 public chain_id;

    // Token address => TokenNetwork address
    mapping(address => address) public token_to_token_networks;

    event TokenNetworkCreated(address indexed token_address, address indexed token_network_address);

    /// @notice constructor for this contract.
    /// @dev    _chain_id must be greater than 0, and _secret_registry_address must exist.
    /// @param  _secrect_registry_address    an address of contract for secret registry.
    /// @param  _chain_id                    a 256-bit unsigned integer
    constructor(address _secret_registry_address, uint256 _chain_id) public {
        require(_chain_id > 0);
        require(_secret_registry_address != 0x0);
        require(contractExists(_secret_registry_address));
        secret_registry_address = _secret_registry_address;
        chain_id = _chain_id;
    }

    /// @notice function to create a ERC20 token network.
    /// @param  _token_address          the place that tokens are from.
    /// @return token_network_address   the address of a token network.
    function createERC20TokenNetwork(address _token_address)
        external
        returns (address token_network_address)
    {
        require(token_to_token_networks[_token_address] == 0x0);

        // Token contract checks are in the corresponding TokenNetwork contract
        token_network_address = new TokenNetwork(
            _token_address,
            secret_registry_address,
            chain_id
        );

        token_to_token_networks[_token_address] = token_network_address;
        emit TokenNetworkCreated(_token_address, token_network_address);

        return token_network_address;
    }
}
