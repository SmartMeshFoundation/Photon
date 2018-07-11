pragma solidity ^0.4.0;

contract C {
    event Log0(bytes32  secrethash);
    event Log1(bytes32 b1,bytes32 bxx,bytes32 b2,bytes32 b3,bytes32 b4,bytes32 b5);
    event Log2(bytes32 b1,bytes32 bxx,bytes32 b2,bytes32 b3,bytes32 b4);
    function f0() public {
        // The next line creates a type error because uint[3] memory
        // cannot be converted to uint[] memory.
       emit Log0(0);
    }
    function f1() public{
         emit Log1(1,2,3,4,5,6);
    }
      function f2() public{
         emit Log2(1,2,3,4,5);
    }
}