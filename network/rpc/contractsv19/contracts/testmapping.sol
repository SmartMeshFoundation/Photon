pragma solidity ^0.4.23;

contract TestMap{
     mapping(uint256 => uint256) public channels;
     function TestSet() external{
        channels[0x39]=0x77;
     }
}