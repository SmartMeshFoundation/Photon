
======= testlog.sol:C =======
EVM assembly:
    /* "testlog.sol":25:535  contract C {... */
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
    /* "testlog.sol":25:535  contract C {... */
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
        /* "testlog.sol":25:535  contract C {... */
      mstore(0x40, 0x80)
      jumpi(tag_1, lt(calldatasize, 0x4))
      calldataload(0x0)
      0x100000000000000000000000000000000000000000000000000000000
      swap1
      div
      0xffffffff
      and
      dup1
      0x9942ec6f
      eq
      tag_2
      jumpi
      dup1
      0xa5850475
      eq
      tag_3
      jumpi
      dup1
      0xc27fc305
      eq
      tag_4
      jumpi
    tag_1:
      0x0
      dup1
      revert
        /* "testlog.sol":475:533  function f2() public{... */
    tag_2:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_5
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_5:
        /* "testlog.sol":475:533  function f2() public{... */
      pop
      tag_6
      jump(tag_7)
    tag_6:
      stop
        /* "testlog.sol":236:403  function f0() public {... */
    tag_3:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_8
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_8:
        /* "testlog.sol":236:403  function f0() public {... */
      pop
      tag_9
      jump(tag_10)
    tag_9:
      stop
        /* "testlog.sol":408:468  function f1() public{... */
    tag_4:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_11
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_11:
        /* "testlog.sol":408:468  function f1() public{... */
      pop
      tag_12
      jump(tag_13)
    tag_12:
      stop
        /* "testlog.sol":475:533  function f2() public{... */
    tag_7:
        /* "testlog.sol":511:526  Log2(1,2,3,4,5) */
      0x3e5fe01be41c5c64c156c4321fa3572869a25c15bb8af3ad6d024d1eb7beeec2
        /* "testlog.sol":516:517  1 */
      0x1
        /* "testlog.sol":518:519  2 */
      0x2
        /* "testlog.sol":520:521  3 */
      0x3
        /* "testlog.sol":522:523  4 */
      0x4
        /* "testlog.sol":524:525  5 */
      0x5
        /* "testlog.sol":511:526  Log2(1,2,3,4,5) */
      mload(0x40)
      dup1
      dup7
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup6
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup5
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup4
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup3
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      swap6
      pop
      pop
      pop
      pop
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      log1
        /* "testlog.sol":475:533  function f2() public{... */
      jump	// out
        /* "testlog.sol":236:403  function f0() public {... */
    tag_10:
        /* "testlog.sol":389:396  Log0(0) */
      0xff4395dddaafcca59f9c36e11fd2e7e13c0360fa63a7431f94e11cf15cdd3203
        /* "testlog.sol":394:395  0 */
      0x0
        /* "testlog.sol":389:396  Log0(0) */
      mload(0x40)
      dup1
      dup3
      0x1
      mul
      not(0x0)
      and
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
      log1
        /* "testlog.sol":236:403  function f0() public {... */
      jump	// out
        /* "testlog.sol":408:468  function f1() public{... */
    tag_13:
        /* "testlog.sol":444:461  Log1(1,2,3,4,5,6) */
      0x1f24f32e32a7a27bff917b1bb401cd069d64b15d6ec3d517d76e161fab7ea7df
        /* "testlog.sol":449:450  1 */
      0x1
        /* "testlog.sol":451:452  2 */
      0x2
        /* "testlog.sol":453:454  3 */
      0x3
        /* "testlog.sol":455:456  4 */
      0x4
        /* "testlog.sol":457:458  5 */
      0x5
        /* "testlog.sol":459:460  6 */
      0x6
        /* "testlog.sol":444:461  Log1(1,2,3,4,5,6) */
      mload(0x40)
      dup1
      dup8
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup7
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup6
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup5
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup4
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup3
      0x1
      mul
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      swap7
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      log1
        /* "testlog.sol":408:468  function f1() public{... */
      jump	// out

    auxdata: 0xa165627a7a72305820fcf51c6a721742a5620b141cd1518a100a7d1ac9aa210deceeb7b065474979a00029
}

