pragma solidity ^0.4.24;

/// @title an interface for Token.
/// @notice it contains various utility functions that will be used by a token.
interface Token {
    /*name,decimals and symbol are optional*/
    /*string public name;                   //fancy name: eg Simon Bucks
    uint8 public decimals;                //How many decimals to show.
    string public symbol;                 //An identifier: eg SBX*/
    /// @return total amount of tokens
    function totalSupply() external view returns (uint256 supply);

    /// @param  _owner The address from which the balance will be retrieved
    /// @return The balance
    function balanceOf(address _owner) external view returns (uint256 balance);

    /// @notice send `_value` token to `_to` from `msg.sender`
    /// @param  _to The address of the recipient
    /// @param  _value The amount of token to be transferred
    /// @return Whether the transfer was successful or not
    function transfer(address _to, uint256 _value) external returns (bool success);

    /// @notice send `_value` token to `_to` from `_from` on the condition it is approved by `_from`
    /// @param  _from The address of the sender
    /// @param  _to The address of the recipient
    /// @param  _value The amount of token to be transferred
    /// @return Whether the transfer was successful or not
    function transferFrom(address _from, address _to, uint256 _value) external returns (bool success);

    /// @notice `msg.sender` approves `_spender` to spend `_value` tokens
    /// @param  _spender The address of the account able to transfer the tokens
    /// @param  _value The amount of wei to be approved for transfer
    /// @return Whether the approval was successful or not
    function approve(address _spender, uint256 _value) external returns (bool success);

    /// @param  _owner The address of the account owning tokens
    /// @param  _spender The address of the account able to transfer the tokens
    /// @return Amount of remaining tokens allowed to spent
    function allowance(address _owner, address _spender) external view returns (uint256 remaining);


    /** ERC20 Extension */
    function approveAndCall(address _spender, uint256 _amount, bytes _extraData) external returns (bool success);


    /** ERC223 Non-Standard */
    function transfer(address to, uint256 value, bytes data) external;

    event Transfer(address indexed from, address indexed to, uint256 value, bytes indexed data);
    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
}
