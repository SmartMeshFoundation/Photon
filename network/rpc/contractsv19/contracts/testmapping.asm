
======= testmapping.sol:TestMap =======
EVM assembly:
    /* "testmapping.sol":26:165  contract TestMap{... */
  mstore(0x40, 0x80)
  callvalue
    /* "--CODEGEN--":8:17   */
  dup1
    /* "--CODEGEN--":5:7   */
  iszero
  tag_1
  jumpi
    /* "--CODEGEN--":30:31   */
  0x0
    /* "--CODEGEN--":27:28   */
  dup1
    /* "--CODEGEN--":20:32   */
  revert
    /* "--CODEGEN--":5:7   */
tag_1:
    /* "testmapping.sol":26:165  contract TestMap{... */
  pop
  dataSize(sub_0)
  dup1
  dataOffset(sub_0)
  0x0
  codecopy
  0x0
  return
stop

sub_0: assembly {
        /* "testmapping.sol":26:165  contract TestMap{... */
      mstore(0x40, 0x80)
      jumpi(tag_1, lt(calldatasize, 0x4))
      calldataload(0x0)
      0x100000000000000000000000000000000000000000000000000000000
      swap1
      div
      0xffffffff
      and
      dup1
      0xd25ff42e
      eq
      tag_2
      jumpi
      dup1
      0xe5949b5d
      eq
      tag_3
      jumpi
    tag_1:
      0x0
      dup1
      revert
        /* "testmapping.sol":99:163  function TestSet() external{... */
    tag_2:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_4
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_4:
        /* "testmapping.sol":99:163  function TestSet() external{... */
      pop
      tag_5
      jump(tag_6)
    tag_5:
      stop
        /* "testmapping.sol":49:92  mapping(uint256 => uint256) public channels */
    tag_3:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_7
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_7:
        /* "testmapping.sol":49:92  mapping(uint256 => uint256) public channels */
      pop
      tag_8
      0x4
      dup1
      calldatasize
      sub
      dup2
      add
      swap1
      dup1
      dup1
      calldataload
      swap1
      0x20
      add
      swap1
      swap3
      swap2
      swap1
      pop
      pop
      pop
      jump(tag_9)
    tag_8:
      mload(0x40)
      dup1
      dup3
      dup2
      mstore
      0x20
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      return
        /* "testmapping.sol":99:163  function TestSet() external{... */
    tag_6:
        /* "testmapping.sol":151:155  0x77 */
      0x77
        /* "testmapping.sol":136:144  channels */
      0x0
        /* "testmapping.sol":136:150  channels[0x39] */
      dup1
        /* "testmapping.sol":145:149  0x39 */
      0x39
        /* "testmapping.sol":136:150  channels[0x39] */
      dup2
      mstore
      0x20
      add
      swap1
      dup2
      mstore
      0x20
      add
      0x0
      keccak256
        /* "testmapping.sol":136:155  channels[0x39]=0x77 */
      dup2
      swap1
      sstore
      pop
        /* "testmapping.sol":99:163  function TestSet() external{... */
      jump	// out
        /* "testmapping.sol":49:92  mapping(uint256 => uint256) public channels */
    tag_9:
      mstore(0x20, 0x0)
      dup1
      0x0
      mstore
      keccak256(0x0, 0x40)
      0x0
      swap2
      pop
      swap1
      pop
      sload
      dup2
      jump	// out

    auxdata: 0xa165627a7a72305820aca27ac3400377ea146d8422cc218e059231a5a274fccc5fd9e9b68ebc3f20f90029
}

