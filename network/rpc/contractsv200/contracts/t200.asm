
======= ECVerify.sol:ECVerify =======
EVM assembly:
    /* "ECVerify.sol":26:1042  library ECVerify {... */
  dataSize(sub_0)
  dataOffset(sub_0)
    /* "--CODEGEN--":132:134   */
  0xb
    /* "--CODEGEN--":166:173   */
  dup3
    /* "--CODEGEN--":155:164   */
  dup3
    /* "--CODEGEN--":146:153   */
  dup3
    /* "--CODEGEN--":137:174   */
  codecopy
    /* "--CODEGEN--":252:259   */
  dup1
    /* "--CODEGEN--":246:260   */
  mload
    /* "--CODEGEN--":243:244   */
  0x0
    /* "--CODEGEN--":238:261   */
  byte
    /* "--CODEGEN--":232:236   */
  0x73
    /* "--CODEGEN--":229:262   */
  eq
    /* "--CODEGEN--":270:271   */
  0x0
    /* "--CODEGEN--":265:285   */
  dup2
  eq
  tag_2
  jumpi
    /* "--CODEGEN--":222:285   */
  jump(tag_1)
    /* "--CODEGEN--":265:285   */
tag_2:
    /* "--CODEGEN--":274:283   */
  invalid
    /* "--CODEGEN--":222:285   */
tag_1:
  pop
    /* "--CODEGEN--":298:307   */
  address
    /* "--CODEGEN--":295:296   */
  0x0
    /* "--CODEGEN--":288:308   */
  mstore
    /* "--CODEGEN--":328:332   */
  0x73
    /* "--CODEGEN--":319:326   */
  dup2
    /* "--CODEGEN--":311:333   */
  mstore8
    /* "--CODEGEN--":352:359   */
  dup3
    /* "--CODEGEN--":343:350   */
  dup2
    /* "--CODEGEN--":336:360   */
  return
stop

sub_0: assembly {
        /* "ECVerify.sol":26:1042  library ECVerify {... */
      eq(address, deployTimeAddress())
      mstore(0x40, 0x80)
      0x0
      dup1
      revert

    auxdata: 0xa165627a7a72305820109ad7eb8e45a0929082f078c8c53b1c0aec6a37f3dea6acd8ebb70b8b28c4e20029
}


======= SecretRegistry.sol:SecretRegistry =======
EVM assembly:
    /* "SecretRegistry.sol":26:1180  contract SecretRegistry {... */
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
    /* "SecretRegistry.sol":26:1180  contract SecretRegistry {... */
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
        /* "SecretRegistry.sol":26:1180  contract SecretRegistry {... */
      mstore(0x40, 0x80)
      jumpi(tag_1, lt(calldatasize, 0x4))
      and(div(calldataload(0x0), 0x100000000000000000000000000000000000000000000000000000000), 0xffffffff)
      0x12ad8bfc
      dup2
      eq
      tag_2
      jumpi
      dup1
      0x97340309
      eq
      tag_3
      jumpi
      dup1
      0xb32c65c8
      eq
      tag_4
      jumpi
      dup1
      0xc1f62946
      eq
      tag_5
      jumpi
    tag_1:
      0x0
      dup1
      revert
        /* "SecretRegistry.sol":674:1031  function registerSecret(bytes32 secret) public returns (bool) {... */
    tag_2:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_6
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_6:
      pop
        /* "SecretRegistry.sol":674:1031  function registerSecret(bytes32 secret) public returns (bool) {... */
      tag_7
      calldataload(0x4)
      jump(tag_8)
    tag_7:
      0x40
      dup1
      mload
      swap2
      iszero
      iszero
      dup3
      mstore
      mload
      swap1
      dup2
      swap1
      sub
      0x20
      add
      swap1
      return
        /* "SecretRegistry.sol":220:274  mapping(bytes32 => uint256) public secrethash_to_block */
    tag_3:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_9
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_9:
      pop
        /* "SecretRegistry.sol":220:274  mapping(bytes32 => uint256) public secrethash_to_block */
      tag_10
      calldataload(0x4)
      jump(tag_11)
    tag_10:
      0x40
      dup1
      mload
      swap2
      dup3
      mstore
      mload
      swap1
      dup2
      swap1
      sub
      0x20
      add
      swap1
      return
        /* "SecretRegistry.sol":97:146  string constant public contract_version = "0.3._" */
    tag_4:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_12
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_12:
        /* "SecretRegistry.sol":97:146  string constant public contract_version = "0.3._" */
      pop
      tag_13
      jump(tag_14)
    tag_13:
      0x40
      dup1
      mload
      0x20
      dup1
      dup3
      mstore
      dup4
      mload
      dup2
      dup4
      add
      mstore
      dup4
      mload
      swap2
      swap3
      dup4
      swap3
      swap1
      dup4
      add
      swap2
      dup6
      add
      swap1
      dup1
      dup4
      dup4
      0x0
        /* "--CODEGEN--":8:108   */
    tag_15:
        /* "--CODEGEN--":33:36   */
      dup4
        /* "--CODEGEN--":30:31   */
      dup2
        /* "--CODEGEN--":27:37   */
      lt
        /* "--CODEGEN--":8:108   */
      iszero
      tag_16
      jumpi
        /* "--CODEGEN--":90:101   */
      dup2
      dup2
      add
        /* "--CODEGEN--":84:102   */
      mload
        /* "--CODEGEN--":71:82   */
      dup4
      dup3
      add
        /* "--CODEGEN--":64:103   */
      mstore
        /* "--CODEGEN--":52:54   */
      0x20
        /* "--CODEGEN--":45:55   */
      add
        /* "--CODEGEN--":8:108   */
      jump(tag_15)
    tag_16:
        /* "--CODEGEN--":12:26   */
      pop
        /* "SecretRegistry.sol":97:146  string constant public contract_version = "0.3._" */
      pop
      pop
      pop
      swap1
      pop
      swap1
      dup2
      add
      swap1
      0x1f
      and
      dup1
      iszero
      tag_18
      jumpi
      dup1
      dup3
      sub
      dup1
      mload
      0x1
      dup4
      0x20
      sub
      0x100
      exp
      sub
      not
      and
      dup2
      mstore
      0x20
      add
      swap2
      pop
    tag_18:
      pop
      swap3
      pop
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      return
        /* "SecretRegistry.sol":1037:1178  function getSecretRevealBlockHeight(bytes32 secrethash) public view returns (uint256) {... */
    tag_5:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_19
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_19:
      pop
        /* "SecretRegistry.sol":1037:1178  function getSecretRevealBlockHeight(bytes32 secrethash) public view returns (uint256) {... */
      tag_10
      calldataload(0x4)
      jump(tag_21)
        /* "SecretRegistry.sol":674:1031  function registerSecret(bytes32 secret) public returns (bool) {... */
    tag_8:
        /* "SecretRegistry.sol":777:801  abi.encodePacked(secret) */
      0x40
      dup1
      mload
      0x20
      dup1
      dup3
      add
      dup5
      swap1
      mstore
      dup3
      mload
        /* "--CODEGEN--":26:47   */
      dup1
      dup4
      sub
        /* "--CODEGEN--":22:54   */
      dup3
      add
        /* "--CODEGEN--":6:55   */
      dup2
      mstore
        /* "SecretRegistry.sol":777:801  abi.encodePacked(secret) */
      swap2
      dup4
      add
      swap3
      dup4
      swap1
      mstore
        /* "SecretRegistry.sol":767:802  keccak256(abi.encodePacked(secret)) */
      dup2
      mload
        /* "SecretRegistry.sol":730:734  bool */
      0x0
      swap4
      dup5
      swap4
        /* "SecretRegistry.sol":777:801  abi.encodePacked(secret) */
      swap3
      swap1
      swap2
      dup3
      swap2
        /* "SecretRegistry.sol":767:802  keccak256(abi.encodePacked(secret)) */
      dup5
      add
      swap1
      dup1
        /* "SecretRegistry.sol":777:801  abi.encodePacked(secret) */
      dup4
        /* "SecretRegistry.sol":767:802  keccak256(abi.encodePacked(secret)) */
      dup4
        /* "--CODEGEN--":36:189   */
    tag_23:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_24
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_23)
    tag_24:
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":344:354   */
      dup2
      mload
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      swap4
      swap1
      swap4
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
      add
        /* "--CODEGEN--":311:320   */
      dup1
      not
        /* "--CODEGEN--":295:321   */
      swap1
      swap2
      and
        /* "--CODEGEN--":340:361   */
      swap3
      and
        /* "--CODEGEN--":377:397   */
      swap2
      swap1
      swap2
      or
        /* "--CODEGEN--":365:398   */
      swap1
      mstore
        /* "SecretRegistry.sol":767:802  keccak256(abi.encodePacked(secret)) */
      mload(0x40)
      swap3
      add
      dup3
      swap1
      sub
      swap1
      swap2
      keccak256
      swap4
      pop
      pop
        /* "SecretRegistry.sol":816:829  secret == 0x0 */
      dup5
      iszero
      swap2
      pop
      dup2
      swap1
      pop
        /* "SecretRegistry.sol":816:868  secret == 0x0 || secrethash_to_block[secrethash] > 0 */
      tag_26
      jumpi
      pop
        /* "SecretRegistry.sol":867:868  0 */
      0x0
        /* "SecretRegistry.sol":833:864  secrethash_to_block[secrethash] */
      dup2
      dup2
      mstore
      0x20
      dup2
      swap1
      mstore
      0x40
      dup2
      keccak256
      sload
        /* "SecretRegistry.sol":833:868  secrethash_to_block[secrethash] > 0 */
      gt
        /* "SecretRegistry.sol":816:868  secret == 0x0 || secrethash_to_block[secrethash] > 0 */
    tag_26:
        /* "SecretRegistry.sol":812:907  if (secret == 0x0 || secrethash_to_block[secrethash] > 0) {... */
      iszero
      tag_27
      jumpi
        /* "SecretRegistry.sol":891:896  false */
      0x0
        /* "SecretRegistry.sol":884:896  return false */
      swap2
      pop
      jump(tag_22)
        /* "SecretRegistry.sol":812:907  if (secret == 0x0 || secrethash_to_block[secrethash] > 0) {... */
    tag_27:
        /* "SecretRegistry.sol":916:935  secrethash_to_block */
      0x0
        /* "SecretRegistry.sol":916:947  secrethash_to_block[secrethash] */
      dup2
      dup2
      mstore
      0x20
      dup2
      swap1
      mstore
      0x40
      dup1
      dup3
      keccak256
        /* "SecretRegistry.sol":950:962  block.number */
      number
        /* "SecretRegistry.sol":916:962  secrethash_to_block[secrethash] = block.number */
      swap1
      sstore
        /* "SecretRegistry.sol":977:1003  SecretRevealed(secrethash) */
      mload
        /* "SecretRegistry.sol":936:946  secrethash */
      dup3
      swap2
        /* "SecretRegistry.sol":977:1003  SecretRevealed(secrethash) */
      0x9b7ddc883342824bd7ddbff103e7a69f8f2e60b96c075cd1b8b8b9713ecc75a4
      swap2
      log2
        /* "SecretRegistry.sol":1020:1024  true */
      0x1
        /* "SecretRegistry.sol":1013:1024  return true */
      swap2
      pop
        /* "SecretRegistry.sol":674:1031  function registerSecret(bytes32 secret) public returns (bool) {... */
    tag_22:
      pop
      swap2
      swap1
      pop
      jump	// out
        /* "SecretRegistry.sol":220:274  mapping(bytes32 => uint256) public secrethash_to_block */
    tag_11:
      0x0
      0x20
      dup2
      swap1
      mstore
      swap1
      dup2
      mstore
      0x40
      swap1
      keccak256
      sload
      dup2
      jump	// out
        /* "SecretRegistry.sol":97:146  string constant public contract_version = "0.3._" */
    tag_14:
      0x40
      dup1
      mload
      dup1
      dup3
      add
      swap1
      swap2
      mstore
      0x5
      dup2
      mstore
      0x302e332e5f000000000000000000000000000000000000000000000000000000
      0x20
      dup3
      add
      mstore
      dup2
      jump	// out
        /* "SecretRegistry.sol":1037:1178  function getSecretRevealBlockHeight(bytes32 secrethash) public view returns (uint256) {... */
    tag_21:
        /* "SecretRegistry.sol":1114:1121  uint256 */
      0x0
        /* "SecretRegistry.sol":1140:1171  secrethash_to_block[secrethash] */
      swap1
      dup2
      mstore
      0x20
      dup2
      swap1
      mstore
      0x40
      swap1
      keccak256
      sload
      swap1
        /* "SecretRegistry.sol":1037:1178  function getSecretRevealBlockHeight(bytes32 secrethash) public view returns (uint256) {... */
      jump	// out

    auxdata: 0xa165627a7a723058208b6cb9059ca52b0f9a206c84719926729a355d21589d1fa4cbcd45cb97f73f750029
}


======= Token.sol:Token =======
EVM assembly:


======= TokenNetwork200.sol:TokenNetwork =======
EVM assembly:
    /* "TokenNetwork200.sol":375:35794  contract TokenNetwork is Utils {... */
  mstore(0x40, 0x80)
    /* "TokenNetwork200.sol":4066:4618  constructor(address _token_address, address _secret_registry, uint256 _chain_id)... */
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
    /* "TokenNetwork200.sol":4066:4618  constructor(address _token_address, address _secret_registry, uint256 _chain_id)... */
  pop
  mload(0x40)
  0x60
  dup1
  bytecodeSize
  dup4
  codecopy
  dup2
  add
  0x40
  swap1
  dup2
  mstore
  dup2
  mload
  0x20
  dup4
  add
  mload
  swap2
  swap1
  swap3
  add
  mload
  sub(exp(0x2, 0xa0), 0x1)
    /* "TokenNetwork200.sol":4180:4201  _token_address != 0x0 */
  dup4
  and
  iszero
  iszero
    /* "TokenNetwork200.sol":4172:4202  require(_token_address != 0x0) */
  tag_4
  jumpi
  0x0
  dup1
  revert
tag_4:
  sub(exp(0x2, 0xa0), 0x1)
    /* "TokenNetwork200.sol":4220:4243  _secret_registry != 0x0 */
  dup3
  and
  iszero
  iszero
    /* "TokenNetwork200.sol":4212:4244  require(_secret_registry != 0x0) */
  tag_5
  jumpi
  0x0
  dup1
  revert
tag_5:
    /* "TokenNetwork200.sol":4274:4275  0 */
  0x0
    /* "TokenNetwork200.sol":4262:4275  _chain_id > 0 */
  dup2
  gt
    /* "TokenNetwork200.sol":4254:4276  require(_chain_id > 0) */
  tag_6
  jumpi
  0x0
  dup1
  revert
tag_6:
    /* "TokenNetwork200.sol":4294:4324  contractExists(_token_address) */
  tag_7
    /* "TokenNetwork200.sol":4309:4323  _token_address */
  dup4
    /* "TokenNetwork200.sol":4294:4308  contractExists */
  0x100000000
  tag_8
  dup2
  mul
    /* "TokenNetwork200.sol":4294:4324  contractExists(_token_address) */
  div
  jump	// in
tag_7:
    /* "TokenNetwork200.sol":4286:4325  require(contractExists(_token_address)) */
  iszero
  iszero
  tag_9
  jumpi
  0x0
  dup1
  revert
tag_9:
    /* "TokenNetwork200.sol":4343:4375  contractExists(_secret_registry) */
  tag_10
    /* "TokenNetwork200.sol":4358:4374  _secret_registry */
  dup3
    /* "TokenNetwork200.sol":4343:4357  contractExists */
  0x100000000
  tag_8
  dup2
  mul
    /* "TokenNetwork200.sol":4343:4375  contractExists(_secret_registry) */
  div
  jump	// in
tag_10:
    /* "TokenNetwork200.sol":4335:4376  require(contractExists(_secret_registry)) */
  iszero
  iszero
  tag_11
  jumpi
  0x0
  dup1
  revert
tag_11:
    /* "TokenNetwork200.sol":4387:4392  token */
  0x0
    /* "TokenNetwork200.sol":4387:4416  token = Token(_token_address) */
  dup1
  sload
  sub(exp(0x2, 0xa0), 0x1)
  dup1
  dup7
  and
  not(sub(exp(0x2, 0xa0), 0x1))
  swap3
  dup4
  and
  or
  dup1
  dup5
  sstore
  0x1
    /* "TokenNetwork200.sol":4427:4477  secret_registry = SecretRegistry(_secret_registry) */
  dup1
  sload
  dup8
  dup5
  and
  swap5
  and
  swap4
  swap1
  swap4
  or
  swap1
  swap3
  sstore
    /* "TokenNetwork200.sol":4487:4495  chain_id */
  0x2
    /* "TokenNetwork200.sol":4487:4507  chain_id = _chain_id */
  dup5
  swap1
  sstore
    /* "TokenNetwork200.sol":4587:4606  token.totalSupply() */
  0x40
  dup1
  mload
  0x18160ddd00000000000000000000000000000000000000000000000000000000
  dup2
  mstore
  swap1
  mload
    /* "TokenNetwork200.sol":4587:4592  token */
  swap3
  swap1
  swap2
  and
  swap2
    /* "TokenNetwork200.sol":4587:4604  token.totalSupply */
  0x18160ddd
  swap2
    /* "TokenNetwork200.sol":4587:4606  token.totalSupply() */
  0x4
  dup1
  dup3
  add
  swap3
  0x20
  swap3
  swap1
  swap2
  swap1
  dup3
  swap1
  sub
  add
  dup2
    /* "TokenNetwork200.sol":4387:4392  token */
  dup8
    /* "TokenNetwork200.sol":4587:4592  token */
  dup8
    /* "TokenNetwork200.sol":4587:4606  token.totalSupply() */
  dup1
  extcodesize
  iszero
    /* "--CODEGEN--":5:7   */
  dup1
  iszero
  tag_12
  jumpi
    /* "--CODEGEN--":30:31   */
  0x0
    /* "--CODEGEN--":27:28   */
  dup1
    /* "--CODEGEN--":20:32   */
  revert
    /* "--CODEGEN--":5:7   */
tag_12:
    /* "TokenNetwork200.sol":4587:4606  token.totalSupply() */
  pop
  gas
  call
  iszero
    /* "--CODEGEN--":8:17   */
  dup1
    /* "--CODEGEN--":5:7   */
  iszero
  tag_13
  jumpi
    /* "--CODEGEN--":45:61   */
  returndatasize
    /* "--CODEGEN--":42:43   */
  0x0
    /* "--CODEGEN--":39:40   */
  dup1
    /* "--CODEGEN--":24:62   */
  returndatacopy
    /* "--CODEGEN--":77:93   */
  returndatasize
    /* "--CODEGEN--":74:75   */
  0x0
    /* "--CODEGEN--":67:94   */
  revert
    /* "--CODEGEN--":5:7   */
tag_13:
    /* "TokenNetwork200.sol":4587:4606  token.totalSupply() */
  pop
  pop
  pop
  pop
  mload(0x40)
  returndatasize
    /* "--CODEGEN--":13:15   */
  0x20
    /* "--CODEGEN--":8:11   */
  dup2
    /* "--CODEGEN--":5:16   */
  lt
    /* "--CODEGEN--":2:4   */
  iszero
  tag_14
  jumpi
    /* "--CODEGEN--":29:30   */
  0x0
    /* "--CODEGEN--":26:27   */
  dup1
    /* "--CODEGEN--":19:31   */
  revert
    /* "--CODEGEN--":2:4   */
tag_14:
  pop
    /* "TokenNetwork200.sol":4587:4606  token.totalSupply() */
  mload
    /* "TokenNetwork200.sol":4587:4610  token.totalSupply() > 0 */
  gt
    /* "TokenNetwork200.sol":4579:4611  require(token.totalSupply() > 0) */
  tag_15
  jumpi
  0x0
  dup1
  revert
tag_15:
    /* "TokenNetwork200.sol":4066:4618  constructor(address _token_address, address _secret_registry, uint256 _chain_id)... */
  pop
  pop
  pop
    /* "TokenNetwork200.sol":375:35794  contract TokenNetwork is Utils {... */
  jump(tag_16)
    /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
tag_8:
    /* "Utils.sol":367:371  bool */
  0x0
    /* "Utils.sol":434:463  extcodesize(contract_address) */
  swap1
  extcodesize
    /* "Utils.sol":490:498  size > 0 */
  gt
  swap1
    /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
  jump	// out
    /* "TokenNetwork200.sol":375:35794  contract TokenNetwork is Utils {... */
tag_16:
  dataSize(sub_0)
  dup1
  dataOffset(sub_0)
  0x0
  codecopy
  0x0
  return
stop

sub_0: assembly {
        /* "TokenNetwork200.sol":375:35794  contract TokenNetwork is Utils {... */
      mstore(0x40, 0x80)
      jumpi(tag_1, lt(calldatasize, 0x4))
      and(div(calldataload(0x0), 0x100000000000000000000000000000000000000000000000000000000), 0xffffffff)
      0x14e5e4de
      dup2
      eq
      tag_2
      jumpi
      dup1
      0x24d73a93
      eq
      tag_3
      jumpi
      dup1
      0x2506f1b3
      eq
      tag_4
      jumpi
      dup1
      0x3af973b1
      eq
      tag_5
      jumpi
      dup1
      0x70ba2a76
      eq
      tag_6
      jumpi
      dup1
      0x7709bc78
      eq
      tag_7
      jumpi
      dup1
      0x7a7ebd7b
      eq
      tag_8
      jumpi
      dup1
      0x7ed74ad9
      eq
      tag_9
      jumpi
      dup1
      0x8568536a
      eq
      tag_10
      jumpi
      dup1
      0x862ceb1a
      eq
      tag_11
      jumpi
      dup1
      0x9375cff2
      eq
      tag_12
      jumpi
      dup1
      0xac133709
      eq
      tag_13
      jumpi
      dup1
      0xaef91441
      eq
      tag_14
      jumpi
      dup1
      0xb32c65c8
      eq
      tag_15
      jumpi
      dup1
      0xc10fd1bb
      eq
      tag_16
      jumpi
      dup1
      0xdaea0ba8
      eq
      tag_17
      jumpi
      dup1
      0xf56451eb
      eq
      tag_18
      jumpi
      dup1
      0xf94c9e13
      eq
      tag_19
      jumpi
      dup1
      0xfc0c546a
      eq
      tag_20
      jumpi
    tag_1:
      0x0
      dup1
      revert
        /* "TokenNetwork200.sol":7168:10171  function withDraw(... */
    tag_2:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_21
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_21:
      pop
        /* "TokenNetwork200.sol":7168:10171  function withDraw(... */
      0x40
      dup1
      mload
      0x20
      0x1f
      calldataload(0xc4)
      0x4
      dup2
      dup2
      add
      calldataload
      swap3
      dup4
      add
      dup5
      swap1
      div
      dup5
      mul
      dup6
      add
      dup5
      add
      swap1
      swap6
      mstore
      dup2
      dup5
      mstore
      tag_22
      swap5
      0xffffffffffffffffffffffffffffffffffffffff
      dup2
      calldataload
      dup2
      and
      swap6
      0x24
      dup1
      calldataload
      swap1
      swap3
      and
      swap6
      calldataload(0x44)
      swap6
      calldataload(0x64)
      swap6
      calldataload(0x84)
      swap6
      calldataload(0xa4)
      swap6
      calldatasize
      swap6
      swap2
      swap5
      0xe4
      swap5
      swap3
      swap4
      swap1
      swap2
      add
      swap2
      swap1
      dup2
      swap1
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      pop
      0x40
      dup1
      mload
      0x20
      0x1f
      dup10
      calldataload
      dup12
      add
      dup1
      calldataload
      swap2
      dup3
      add
      dup4
      swap1
      div
      dup4
      mul
      dup5
      add
      dup4
      add
      swap1
      swap5
      mstore
      dup1
      dup4
      mstore
      swap8
      swap11
      swap10
      swap9
      dup2
      add
      swap8
      swap2
      swap7
      pop
      swap2
      dup3
      add
      swap5
      pop
      swap3
      pop
      dup3
      swap2
      pop
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      swap5
      swap8
      pop
      tag_23
      swap7
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump
    tag_22:
      stop
        /* "TokenNetwork200.sol":806:843  SecretRegistry public secret_registry */
    tag_3:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_24
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_24:
        /* "TokenNetwork200.sol":806:843  SecretRegistry public secret_registry */
      pop
      tag_25
      jump(tag_26)
    tag_25:
      0x40
      dup1
      mload
      0xffffffffffffffffffffffffffffffffffffffff
      swap1
      swap3
      and
      dup3
      mstore
      mload
      swap1
      dup2
      swap1
      sub
      0x20
      add
      swap1
      return
        /* "TokenNetwork200.sol":19476:23609  function settleChannel(... */
    tag_4:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_27
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_27:
      pop
        /* "TokenNetwork200.sol":19476:23609  function settleChannel(... */
      tag_22
      0xffffffffffffffffffffffffffffffffffffffff
      calldataload(0x4)
      dup2
      and
      swap1
      calldataload(0x24)
      and
      calldataload(0x44)
      calldataload(0x64)
      calldataload(0x84)
      calldataload(0xa4)
      calldataload(0xc4)
      calldataload(0xe4)
      jump(tag_29)
        /* "TokenNetwork200.sol":996:1019  uint256 public chain_id */
    tag_5:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_30
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_30:
        /* "TokenNetwork200.sol":996:1019  uint256 public chain_id */
      pop
      tag_31
      jump(tag_32)
    tag_31:
      0x40
      dup1
      mload
      swap2
      dup3
      mstore
      mload
      swap1
      dup2
      swap1
      sub
      0x20
      add
      swap1
      return
        /* "TokenNetwork200.sol":12122:13929  function updateBalanceProof(... */
    tag_6:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_33
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_33:
      pop
        /* "TokenNetwork200.sol":12122:13929  function updateBalanceProof(... */
      0x40
      dup1
      mload
      0x20
      0x1f
      calldataload(0x124)
      0x4
      dup2
      dup2
      add
      calldataload
      swap3
      dup4
      add
      dup5
      swap1
      div
      dup5
      mul
      dup6
      add
      dup5
      add
      swap1
      swap6
      mstore
      dup2
      dup5
      mstore
      tag_22
      swap5
      0xffffffffffffffffffffffffffffffffffffffff
      dup2
      calldataload
      dup2
      and
      swap6
      0x24
      dup1
      calldataload
      swap1
      swap3
      and
      swap6
      calldataload(0x44)
      swap6
      calldataload(0x64)
      swap6
      calldataload(0x84)
      swap6
      calldataload(0xa4)
      swap6
      calldataload(0xc4)
      swap6
      calldataload(0xe4)
      swap6
      calldataload(0x104)
      swap6
      calldatasize
      swap6
      0x144
      swap5
      add
      swap2
      dup2
      swap1
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      pop
      0x40
      dup1
      mload
      0x20
      0x1f
      dup10
      calldataload
      dup12
      add
      dup1
      calldataload
      swap2
      dup3
      add
      dup4
      swap1
      div
      dup4
      mul
      dup5
      add
      dup4
      add
      swap1
      swap5
      mstore
      dup1
      dup4
      mstore
      swap8
      swap11
      swap10
      swap9
      dup2
      add
      swap8
      swap2
      swap7
      pop
      swap2
      dup3
      add
      swap5
      pop
      swap3
      pop
      dup3
      swap2
      pop
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      swap5
      swap8
      pop
      tag_35
      swap7
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump
        /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
    tag_7:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_36
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_36:
      pop
        /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
      tag_37
      and(calldataload(0x4), 0xffffffffffffffffffffffffffffffffffffffff)
      jump(tag_38)
    tag_37:
      0x40
      dup1
      mload
      swap2
      iszero
      iszero
      dup3
      mstore
      mload
      swap1
      dup2
      swap1
      sub
      0x20
      add
      swap1
      return
        /* "TokenNetwork200.sol":1100:1143  mapping(bytes32 => Channel) public channels */
    tag_8:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_39
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_39:
      pop
        /* "TokenNetwork200.sol":1100:1143  mapping(bytes32 => Channel) public channels */
      tag_40
      calldataload(0x4)
      jump(tag_41)
    tag_40:
      0x40
      dup1
      mload
      0xffffffffffffffff
      swap5
      dup6
      and
      dup2
      mstore
      swap3
      swap1
      swap4
      and
      0x20
      dup4
      add
      mstore
      0xff
      and
      dup2
      dup4
      add
      mstore
      swap1
      mload
      swap1
      dup2
      swap1
      sub
      0x60
      add
      swap1
      return
        /* "TokenNetwork200.sol":508:612  bytes32 constant public invalid_balance_hash=keccak256(uint64(0xffffffffffffffff),bytes32(0),uint256(0)) */
    tag_9:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_42
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_42:
        /* "TokenNetwork200.sol":508:612  bytes32 constant public invalid_balance_hash=keccak256(uint64(0xffffffffffffffff),bytes32(0),uint256(0)) */
      pop
      tag_31
      jump(tag_44)
        /* "TokenNetwork200.sol":23678:26342  function cooperativeSettle(... */
    tag_10:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_45
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_45:
      pop
        /* "TokenNetwork200.sol":23678:26342  function cooperativeSettle(... */
      0x40
      dup1
      mload
      0x20
      0x1f
      calldataload(0x84)
      0x4
      dup2
      dup2
      add
      calldataload
      swap3
      dup4
      add
      dup5
      swap1
      div
      dup5
      mul
      dup6
      add
      dup5
      add
      swap1
      swap6
      mstore
      dup2
      dup5
      mstore
      tag_22
      swap5
      0xffffffffffffffffffffffffffffffffffffffff
      dup2
      calldataload
      dup2
      and
      swap6
      0x24
      dup1
      calldataload
      swap7
      calldataload(0x44)
      swap1
      swap4
      and
      swap6
      calldataload(0x64)
      swap6
      calldatasize
      swap6
      swap5
      0xa4
      swap5
      swap4
      swap2
      swap1
      swap2
      add
      swap2
      swap1
      dup2
      swap1
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      pop
      0x40
      dup1
      mload
      0x20
      0x1f
      dup10
      calldataload
      dup12
      add
      dup1
      calldataload
      swap2
      dup3
      add
      dup4
      swap1
      div
      dup4
      mul
      dup5
      add
      dup4
      add
      swap1
      swap5
      mstore
      dup1
      dup4
      mstore
      swap8
      swap11
      swap10
      swap9
      dup2
      add
      swap8
      swap2
      swap7
      pop
      swap2
      dup3
      add
      swap5
      pop
      swap3
      pop
      dup3
      swap2
      pop
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      swap5
      swap8
      pop
      tag_47
      swap7
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump
        /* "TokenNetwork200.sol":17485:19408  function punishObsoleteUnlock(... */
    tag_11:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_48
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_48:
      pop
        /* "TokenNetwork200.sol":17485:19408  function punishObsoleteUnlock(... */
      0x40
      dup1
      mload
      0x20
      0x1f
      calldataload(0xc4)
      0x4
      dup2
      dup2
      add
      calldataload
      swap3
      dup4
      add
      dup5
      swap1
      div
      dup5
      mul
      dup6
      add
      dup5
      add
      swap1
      swap6
      mstore
      dup2
      dup5
      mstore
      tag_22
      swap5
      0xffffffffffffffffffffffffffffffffffffffff
      dup2
      calldataload
      dup2
      and
      swap6
      0x24
      dup1
      calldataload
      swap1
      swap3
      and
      swap6
      calldataload(0x44)
      swap6
      calldataload(0x64)
      swap6
      calldataload(0x84)
      swap6
      calldataload(0xa4)
      swap6
      calldatasize
      swap6
      swap2
      swap5
      0xe4
      swap5
      swap3
      swap4
      swap1
      swap2
      add
      swap2
      swap1
      dup2
      swap1
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      pop
      0x40
      dup1
      mload
      0x20
      0x1f
      dup10
      calldataload
      dup12
      add
      dup1
      calldataload
      swap2
      dup3
      add
      dup4
      swap1
      div
      dup4
      mul
      dup5
      add
      dup4
      add
      swap1
      swap5
      mstore
      dup1
      dup4
      mstore
      swap8
      swap11
      swap10
      swap9
      dup2
      add
      swap8
      swap2
      swap7
      pop
      swap2
      dup3
      add
      swap5
      pop
      swap3
      pop
      dup3
      swap2
      pop
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      swap5
      swap8
      pop
      tag_50
      swap7
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump
        /* "TokenNetwork200.sol":849:894  uint64 constant public punish_block_number=10 */
    tag_12:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_51
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_51:
        /* "TokenNetwork200.sol":849:894  uint64 constant public punish_block_number=10 */
      pop
      tag_52
      jump(tag_53)
    tag_52:
      0x40
      dup1
      mload
      0xffffffffffffffff
      swap1
      swap3
      and
      dup3
      mstore
      mload
      swap1
      dup2
      swap1
      sub
      0x20
      add
      swap1
      return
        /* "TokenNetwork200.sol":28519:28983  function getChannelParticipantInfo( address participant,address partner)... */
    tag_13:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_54
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_54:
      pop
        /* "TokenNetwork200.sol":28519:28983  function getChannelParticipantInfo( address participant,address partner)... */
      tag_55
      0xffffffffffffffffffffffffffffffffffffffff
      calldataload(0x4)
      dup2
      and
      swap1
      calldataload(0x24)
      and
      jump(tag_56)
    tag_55:
      0x40
      dup1
      mload
      swap3
      dup4
      mstore
      0x20
      dup4
      add
      swap2
      swap1
      swap2
      mstore
      dup1
      mload
      swap2
      dup3
      swap1
      sub
      add
      swap1
      return
        /* "TokenNetwork200.sol":4666:5550  function openChannel(address participant1, address participant2, uint64 settle_timeout)... */
    tag_14:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_57
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_57:
      pop
        /* "TokenNetwork200.sol":4666:5550  function openChannel(address participant1, address participant2, uint64 settle_timeout)... */
      tag_22
      0xffffffffffffffffffffffffffffffffffffffff
      calldataload(0x4)
      dup2
      and
      swap1
      calldataload(0x24)
      and
      and(calldataload(0x44), 0xffffffffffffffff)
      jump(tag_59)
        /* "TokenNetwork200.sol":453:502  string constant public contract_version = "0.3._" */
    tag_15:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_60
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_60:
        /* "TokenNetwork200.sol":453:502  string constant public contract_version = "0.3._" */
      pop
      tag_61
      jump(tag_62)
    tag_61:
      0x40
      dup1
      mload
      0x20
      dup1
      dup3
      mstore
      dup4
      mload
      dup2
      dup4
      add
      mstore
      dup4
      mload
      swap2
      swap3
      dup4
      swap3
      swap1
      dup4
      add
      swap2
      dup6
      add
      swap1
      dup1
      dup4
      dup4
      0x0
        /* "--CODEGEN--":8:108   */
    tag_63:
        /* "--CODEGEN--":33:36   */
      dup4
        /* "--CODEGEN--":30:31   */
      dup2
        /* "--CODEGEN--":27:37   */
      lt
        /* "--CODEGEN--":8:108   */
      iszero
      tag_64
      jumpi
        /* "--CODEGEN--":90:101   */
      dup2
      dup2
      add
        /* "--CODEGEN--":84:102   */
      mload
        /* "--CODEGEN--":71:82   */
      dup4
      dup3
      add
        /* "--CODEGEN--":64:103   */
      mstore
        /* "--CODEGEN--":52:54   */
      0x20
        /* "--CODEGEN--":45:55   */
      add
        /* "--CODEGEN--":8:108   */
      jump(tag_63)
    tag_64:
        /* "--CODEGEN--":12:26   */
      pop
        /* "TokenNetwork200.sol":453:502  string constant public contract_version = "0.3._" */
      pop
      pop
      pop
      swap1
      pop
      swap1
      dup2
      add
      swap1
      0x1f
      and
      dup1
      iszero
      tag_66
      jumpi
      dup1
      dup3
      sub
      dup1
      mload
      0x1
      dup4
      0x20
      sub
      0x100
      exp
      sub
      not
      and
      dup2
      mstore
      0x20
      add
      swap2
      pop
    tag_66:
      pop
      swap3
      pop
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      return
        /* "TokenNetwork200.sol":5699:7055  function setTotalDeposit(address participant,address partner, uint256 total_deposit)... */
    tag_16:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_67
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_67:
      pop
        /* "TokenNetwork200.sol":5699:7055  function setTotalDeposit(address participant,address partner, uint256 total_deposit)... */
      tag_22
      0xffffffffffffffffffffffffffffffffffffffff
      calldataload(0x4)
      dup2
      and
      swap1
      calldataload(0x24)
      and
      calldataload(0x44)
      jump(tag_69)
        /* "TokenNetwork200.sol":10282:11927  function closeChannel(... */
    tag_17:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_70
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_70:
      pop
        /* "TokenNetwork200.sol":10282:11927  function closeChannel(... */
      0x40
      dup1
      mload
      0x20
      0x4
      calldataload(0xa4)
      dup2
      dup2
      add
      calldataload
      0x1f
      dup2
      add
      dup5
      swap1
      div
      dup5
      mul
      dup6
      add
      dup5
      add
      swap1
      swap6
      mstore
      dup5
      dup5
      mstore
      tag_22
      swap5
      dup3
      calldataload
      0xffffffffffffffffffffffffffffffffffffffff
      and
      swap5
      0x24
      dup1
      calldataload
      swap6
      calldataload(0x44)
      swap6
      calldataload(0x64)
      swap6
      calldataload(0x84)
      swap6
      calldatasize
      swap6
      swap3
      swap5
      0xc4
      swap5
      swap3
      add
      swap2
      dup2
      swap1
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      swap5
      swap8
      pop
      tag_72
      swap7
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump
        /* "TokenNetwork200.sol":14088:16317  function unlock(... */
    tag_18:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_73
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_73:
      pop
        /* "TokenNetwork200.sol":14088:16317  function unlock(... */
      0x40
      dup1
      mload
      0x20
      0x4
      calldataload(0xa4)
      dup2
      dup2
      add
      calldataload
      0x1f
      dup2
      add
      dup5
      swap1
      div
      dup5
      mul
      dup6
      add
      dup5
      add
      swap1
      swap6
      mstore
      dup5
      dup5
      mstore
      tag_22
      swap5
      dup3
      calldataload
      0xffffffffffffffffffffffffffffffffffffffff
      swap1
      dup2
      and
      swap6
      0x24
      dup1
      calldataload
      swap1
      swap3
      and
      swap6
      calldataload(0x44)
      swap6
      calldataload(0x64)
      swap6
      calldataload(0x84)
      swap6
      calldatasize
      swap6
      swap3
      swap5
      0xc4
      swap5
      swap1
      swap4
      swap1
      swap3
      add
      swap2
      dup2
      swap1
      dup5
      add
      dup4
      dup3
      dup1
      dup3
      dup5
      calldatacopy
      pop
      swap5
      swap8
      pop
      tag_75
      swap7
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump
        /* "TokenNetwork200.sol":28042:28513  function getChannelInfo(address participant1,address participant2)... */
    tag_19:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_76
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_76:
      pop
        /* "TokenNetwork200.sol":28042:28513  function getChannelInfo(address participant1,address participant2)... */
      tag_77
      0xffffffffffffffffffffffffffffffffffffffff
      calldataload(0x4)
      dup2
      and
      swap1
      calldataload(0x24)
      and
      jump(tag_78)
    tag_77:
      0x40
      dup1
      mload
      swap5
      dup6
      mstore
      0xffffffffffffffff
      swap4
      dup5
      and
      0x20
      dup7
      add
      mstore
      swap2
      swap1
      swap3
      and
      dup4
      dup3
      add
      mstore
      0xff
      swap1
      swap2
      and
      0x60
      dup4
      add
      mstore
      mload
      swap1
      dup2
      swap1
      sub
      0x80
      add
      swap1
      return
        /* "TokenNetwork200.sol":689:707  Token public token */
    tag_20:
      callvalue
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_79
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_79:
        /* "TokenNetwork200.sol":689:707  Token public token */
      pop
      tag_25
      jump(tag_81)
        /* "TokenNetwork200.sol":7168:10171  function withDraw(... */
    tag_23:
        /* "TokenNetwork200.sol":7507:7528  uint256 total_deposit */
      0x0
        /* "TokenNetwork200.sol":7538:7564  bytes32 channel_identifier */
      dup1
        /* "TokenNetwork200.sol":7650:7673  Channel storage channel */
      0x0
        /* "TokenNetwork200.sol":7783:7803  bytes32 message_hash */
      dup1
        /* "TokenNetwork200.sol":8794:8832  Participant storage participant1_state */
      0x0
        /* "TokenNetwork200.sol":8879:8917  Participant storage participant2_state */
      dup1
        /* "TokenNetwork200.sol":7593:7640  getChannelIdentifier(participant1,participant2) */
      tag_83
        /* "TokenNetwork200.sol":7614:7626  participant1 */
      dup15
        /* "TokenNetwork200.sol":7627:7639  participant2 */
      dup15
        /* "TokenNetwork200.sol":7593:7613  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":7593:7640  getChannelIdentifier(participant1,participant2) */
      jump	// in
    tag_83:
        /* "TokenNetwork200.sol":7676:7704  channels[channel_identifier] */
      0x0
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":7676:7684  channels */
      0x3
        /* "TokenNetwork200.sol":7676:7704  channels[channel_identifier] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":7722:7735  channel.state */
      dup1
      sload
        /* "TokenNetwork200.sol":7676:7704  channels[channel_identifier] */
      swap2
      swap7
      pop
      swap5
      pop
        /* "TokenNetwork200.sol":7722:7735  channel.state */
      0x100000000000000000000000000000000
      swap1
      div
      0xff
      and
        /* "TokenNetwork200.sol":7739:7740  1 */
      0x1
        /* "TokenNetwork200.sol":7722:7740  channel.state == 1 */
      eq
        /* "TokenNetwork200.sol":7714:7741  require(channel.state == 1) */
      tag_85
      jumpi
      0x0
      dup1
      revert
    tag_85:
        /* "TokenNetwork200.sol":7850:7862  participant1 */
      dup14
        /* "TokenNetwork200.sol":7880:7900  participant1_deposit */
      dup13
        /* "TokenNetwork200.sol":7918:7930  participant2 */
      dup15
        /* "TokenNetwork200.sol":7948:7968  participant2_deposit */
      dup14
        /* "TokenNetwork200.sol":7986:8007  participant1_withdraw */
      dup14
        /* "TokenNetwork200.sol":8025:8043  channel_identifier */
      dup10
        /* "TokenNetwork200.sol":8061:8068  channel */
      dup10
        /* "TokenNetwork200.sol":8061:8085  channel.open_blocknumber */
      0x0
      add
      0x8
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffff
      and
        /* "TokenNetwork200.sol":8111:8115  this */
      address
        /* "TokenNetwork200.sol":8134:8142  chain_id */
      sload(0x2)
        /* "TokenNetwork200.sol":7816:8156  abi.encodePacked(... */
      add(0x20, mload(0x40))
      dup1
      dup11
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      dup10
      dup2
      mstore
      0x20
      add
      dup9
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      dup8
      dup2
      mstore
      0x20
      add
      dup7
      dup2
      mstore
      0x20
      add
      dup6
      not(0x0)
      and
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup5
      0xffffffffffffffff
      and
      0xffffffffffffffff
      and
      0x1000000000000000000000000000000000000000000000000
      mul
      dup2
      mstore
      0x8
      add
      dup4
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      dup3
      dup2
      mstore
      0x20
      add
      swap10
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      mload(0x40)
        /* "--CODEGEN--":49:53   */
      0x20
        /* "--CODEGEN--":39:46   */
      dup2
        /* "--CODEGEN--":30:37   */
      dup4
        /* "--CODEGEN--":26:47   */
      sub
        /* "--CODEGEN--":22:54   */
      sub
        /* "--CODEGEN--":13:20   */
      dup2
        /* "--CODEGEN--":6:55   */
      mstore
        /* "TokenNetwork200.sol":7816:8156  abi.encodePacked(... */
      swap1
      0x40
      mstore
        /* "TokenNetwork200.sol":7806:8157  keccak256(abi.encodePacked(... */
      mload(0x40)
      dup1
      dup3
      dup1
      mload
      swap1
      0x20
      add
      swap1
      dup1
      dup4
      dup4
        /* "--CODEGEN--":36:189   */
    tag_86:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_87
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_86)
    tag_87:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":7806:8157  keccak256(abi.encodePacked(... */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":7783:8157  bytes32 message_hash = keccak256(abi.encodePacked(... */
      swap3
      pop
        /* "TokenNetwork200.sol":8191:8246  ECVerify.ecverify(message_hash, participant1_signature) */
      tag_89
        /* "TokenNetwork200.sol":8209:8221  message_hash */
      dup4
        /* "TokenNetwork200.sol":8223:8245  participant1_signature */
      dup10
        /* "TokenNetwork200.sol":8191:8208  ECVerify.ecverify */
      tag_90
        /* "TokenNetwork200.sol":8191:8246  ECVerify.ecverify(message_hash, participant1_signature) */
      jump	// in
    tag_89:
        /* "TokenNetwork200.sol":8175:8246  participant1 == ECVerify.ecverify(message_hash, participant1_signature) */
      0xffffffffffffffffffffffffffffffffffffffff
      dup16
      dup2
      and
      swap2
      and
      eq
        /* "TokenNetwork200.sol":8167:8247  require(participant1 == ECVerify.ecverify(message_hash, participant1_signature)) */
      tag_91
      jumpi
      0x0
      dup1
      revert
    tag_91:
        /* "TokenNetwork200.sol":8348:8360  participant1 */
      dup14
        /* "TokenNetwork200.sol":8378:8398  participant1_deposit */
      dup13
        /* "TokenNetwork200.sol":8416:8428  participant2 */
      dup15
        /* "TokenNetwork200.sol":8446:8466  participant2_deposit */
      dup14
        /* "TokenNetwork200.sol":8484:8505  participant1_withdraw */
      dup14
        /* "TokenNetwork200.sol":8523:8544  participant2_withdraw */
      dup14
        /* "TokenNetwork200.sol":8562:8580  channel_identifier */
      dup11
        /* "TokenNetwork200.sol":8598:8605  channel */
      dup11
        /* "TokenNetwork200.sol":8598:8622  channel.open_blocknumber */
      0x0
      add
      0x8
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffff
      and
        /* "TokenNetwork200.sol":8648:8652  this */
      address
        /* "TokenNetwork200.sol":8671:8679  chain_id */
      sload(0x2)
        /* "TokenNetwork200.sol":8314:8693  abi.encodePacked(... */
      add(0x20, mload(0x40))
      dup1
      dup12
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      dup11
      dup2
      mstore
      0x20
      add
      dup10
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      dup9
      dup2
      mstore
      0x20
      add
      dup8
      dup2
      mstore
      0x20
      add
      dup7
      dup2
      mstore
      0x20
      add
      dup6
      not(0x0)
      and
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup5
      0xffffffffffffffff
      and
      0xffffffffffffffff
      and
      0x1000000000000000000000000000000000000000000000000
      mul
      dup2
      mstore
      0x8
      add
      dup4
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      dup3
      dup2
      mstore
      0x20
      add
      swap11
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      mload(0x40)
        /* "--CODEGEN--":49:53   */
      0x20
        /* "--CODEGEN--":39:46   */
      dup2
        /* "--CODEGEN--":30:37   */
      dup4
        /* "--CODEGEN--":26:47   */
      sub
        /* "--CODEGEN--":22:54   */
      sub
        /* "--CODEGEN--":13:20   */
      dup2
        /* "--CODEGEN--":6:55   */
      mstore
        /* "TokenNetwork200.sol":8314:8693  abi.encodePacked(... */
      swap1
      0x40
      mstore
        /* "TokenNetwork200.sol":8304:8694  keccak256(abi.encodePacked(... */
      mload(0x40)
      dup1
      dup3
      dup1
      mload
      swap1
      0x20
      add
      swap1
      dup1
      dup4
      dup4
        /* "--CODEGEN--":36:189   */
    tag_92:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_93
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_92)
    tag_93:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":8304:8694  keccak256(abi.encodePacked(... */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":8289:8694  message_hash = keccak256(abi.encodePacked(... */
      swap3
      pop
        /* "TokenNetwork200.sol":8728:8783  ECVerify.ecverify(message_hash, participant2_signature) */
      tag_95
        /* "TokenNetwork200.sol":8746:8758  message_hash */
      dup4
        /* "TokenNetwork200.sol":8760:8782  participant2_signature */
      dup9
        /* "TokenNetwork200.sol":8728:8745  ECVerify.ecverify */
      tag_90
        /* "TokenNetwork200.sol":8728:8783  ECVerify.ecverify(message_hash, participant2_signature) */
      jump	// in
    tag_95:
        /* "TokenNetwork200.sol":8712:8783  participant2 == ECVerify.ecverify(message_hash, participant2_signature) */
      0xffffffffffffffffffffffffffffffffffffffff
      dup15
      dup2
      and
      swap2
      and
      eq
        /* "TokenNetwork200.sol":8704:8784  require(participant2 == ECVerify.ecverify(message_hash, participant2_signature)) */
      tag_96
      jumpi
      0x0
      dup1
      revert
    tag_96:
      pop
      pop
        /* "TokenNetwork200.sol":8835:8869  channel.participants[participant1] */
      0xffffffffffffffffffffffffffffffffffffffff
      dup1
      dup14
      and
      0x0
      swap1
      dup2
      mstore
        /* "TokenNetwork200.sol":8835:8855  channel.participants */
      0x1
      dup5
      add
        /* "TokenNetwork200.sol":8835:8869  channel.participants[participant1] */
      0x20
      mstore
      0x40
      dup1
      dup3
      keccak256
        /* "TokenNetwork200.sol":8920:8954  channel.participants[participant2] */
      swap3
      dup15
      and
      dup3
      mstore
      swap1
      keccak256
        /* "TokenNetwork200.sol":9096:9122  participant2_state.deposit */
      dup1
      sload
        /* "TokenNetwork200.sol":9067:9093  participant1_state.deposit */
      dup3
      sload
        /* "TokenNetwork200.sol":9067:9122  participant1_state.deposit + participant2_state.deposit */
      add
      swap6
      pop
        /* "TokenNetwork200.sol":9140:9177  participant1_deposit <= total_deposit */
      dup6
      dup13
      gt
      iszero
        /* "TokenNetwork200.sol":9132:9178  require(participant1_deposit <= total_deposit) */
      tag_97
      jumpi
      0x0
      dup1
      revert
    tag_97:
        /* "TokenNetwork200.sol":9196:9233  participant2_deposit <= total_deposit */
      dup6
      dup12
      gt
      iszero
        /* "TokenNetwork200.sol":9188:9234  require(participant2_deposit <= total_deposit) */
      tag_98
      jumpi
      0x0
      dup1
      revert
    tag_98:
        /* "TokenNetwork200.sol":9253:9296  participant1_deposit + participant2_deposit */
      dup12
      dup12
      add
        /* "TokenNetwork200.sol":9252:9314  (participant1_deposit + participant2_deposit) == total_deposit */
      dup7
      eq
        /* "TokenNetwork200.sol":9244:9315  require((participant1_deposit + participant2_deposit) == total_deposit) */
      tag_99
      jumpi
      0x0
      dup1
      revert
    tag_99:
        /* "TokenNetwork200.sol":9388:9389  0 */
      0x0
        /* "TokenNetwork200.sol":9364:9385  participant1_withdraw */
      dup11
        /* "TokenNetwork200.sol":9364:9389  participant1_withdraw > 0 */
      gt
        /* "TokenNetwork200.sol":9360:9476  if (participant1_withdraw > 0) {... */
      iszero
      tag_104
      jumpi
        /* "TokenNetwork200.sol":9413:9418  token */
      0x0
      dup1
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffffffffffffffffffffffffffff
      and
        /* "TokenNetwork200.sol":9413:9427  token.transfer */
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xa9059cbb
        /* "TokenNetwork200.sol":9428:9440  participant1 */
      dup16
        /* "TokenNetwork200.sol":9442:9463  participant1_withdraw */
      dup13
        /* "TokenNetwork200.sol":9413:9464  token.transfer(participant1, participant1_withdraw) */
      mload(0x40)
      dup4
      0xffffffff
      and
      0x100000000000000000000000000000000000000000000000000000000
      mul
      dup2
      mstore
      0x4
      add
      dup1
      dup4
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      dup2
      mstore
      0x20
      add
      dup3
      dup2
      mstore
      0x20
      add
      swap3
      pop
      pop
      pop
      0x20
      mload(0x40)
      dup1
      dup4
      sub
      dup2
      0x0
      dup8
      dup1
      extcodesize
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_101
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_101:
        /* "TokenNetwork200.sol":9413:9464  token.transfer(participant1, participant1_withdraw) */
      pop
      gas
      call
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_102
      jumpi
        /* "--CODEGEN--":45:61   */
      returndatasize
        /* "--CODEGEN--":42:43   */
      0x0
        /* "--CODEGEN--":39:40   */
      dup1
        /* "--CODEGEN--":24:62   */
      returndatacopy
        /* "--CODEGEN--":77:93   */
      returndatasize
        /* "--CODEGEN--":74:75   */
      0x0
        /* "--CODEGEN--":67:94   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_102:
        /* "TokenNetwork200.sol":9413:9464  token.transfer(participant1, participant1_withdraw) */
      pop
      pop
      pop
      pop
      mload(0x40)
      returndatasize
        /* "--CODEGEN--":13:15   */
      0x20
        /* "--CODEGEN--":8:11   */
      dup2
        /* "--CODEGEN--":5:16   */
      lt
        /* "--CODEGEN--":2:4   */
      iszero
      tag_103
      jumpi
        /* "--CODEGEN--":29:30   */
      0x0
        /* "--CODEGEN--":26:27   */
      dup1
        /* "--CODEGEN--":19:31   */
      revert
        /* "--CODEGEN--":2:4   */
    tag_103:
      pop
        /* "TokenNetwork200.sol":9413:9464  token.transfer(participant1, participant1_withdraw) */
      mload
        /* "TokenNetwork200.sol":9405:9465  require(token.transfer(participant1, participant1_withdraw)) */
      iszero
      iszero
      tag_104
      jumpi
      0x0
      dup1
      revert
    tag_104:
        /* "TokenNetwork200.sol":9513:9514  0 */
      0x0
        /* "TokenNetwork200.sol":9489:9510  participant2_withdraw */
      dup10
        /* "TokenNetwork200.sol":9489:9514  participant2_withdraw > 0 */
      gt
        /* "TokenNetwork200.sol":9485:9601  if (participant2_withdraw > 0) {... */
      iszero
      tag_109
      jumpi
        /* "TokenNetwork200.sol":9538:9543  token */
      0x0
      dup1
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffffffffffffffffffffffffffff
      and
        /* "TokenNetwork200.sol":9538:9552  token.transfer */
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xa9059cbb
        /* "TokenNetwork200.sol":9553:9565  participant2 */
      dup15
        /* "TokenNetwork200.sol":9567:9588  participant2_withdraw */
      dup12
        /* "TokenNetwork200.sol":9538:9589  token.transfer(participant2, participant2_withdraw) */
      mload(0x40)
      dup4
      0xffffffff
      and
      0x100000000000000000000000000000000000000000000000000000000
      mul
      dup2
      mstore
      0x4
      add
      dup1
      dup4
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      dup2
      mstore
      0x20
      add
      dup3
      dup2
      mstore
      0x20
      add
      swap3
      pop
      pop
      pop
      0x20
      mload(0x40)
      dup1
      dup4
      sub
      dup2
      0x0
      dup8
      dup1
      extcodesize
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_106
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_106:
        /* "TokenNetwork200.sol":9538:9589  token.transfer(participant2, participant2_withdraw) */
      pop
      gas
      call
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_107
      jumpi
        /* "--CODEGEN--":45:61   */
      returndatasize
        /* "--CODEGEN--":42:43   */
      0x0
        /* "--CODEGEN--":39:40   */
      dup1
        /* "--CODEGEN--":24:62   */
      returndatacopy
        /* "--CODEGEN--":77:93   */
      returndatasize
        /* "--CODEGEN--":74:75   */
      0x0
        /* "--CODEGEN--":67:94   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_107:
        /* "TokenNetwork200.sol":9538:9589  token.transfer(participant2, participant2_withdraw) */
      pop
      pop
      pop
      pop
      mload(0x40)
      returndatasize
        /* "--CODEGEN--":13:15   */
      0x20
        /* "--CODEGEN--":8:11   */
      dup2
        /* "--CODEGEN--":5:16   */
      lt
        /* "--CODEGEN--":2:4   */
      iszero
      tag_108
      jumpi
        /* "--CODEGEN--":29:30   */
      0x0
        /* "--CODEGEN--":26:27   */
      dup1
        /* "--CODEGEN--":19:31   */
      revert
        /* "--CODEGEN--":2:4   */
    tag_108:
      pop
        /* "TokenNetwork200.sol":9538:9589  token.transfer(participant2, participant2_withdraw) */
      mload
        /* "TokenNetwork200.sol":9530:9590  require(token.transfer(participant2, participant2_withdraw)) */
      iszero
      iszero
      tag_109
      jumpi
      0x0
      dup1
      revert
    tag_109:
        /* "TokenNetwork200.sol":9618:9663  participant1_withdraw <= participant1_deposit */
      dup12
      dup11
      gt
      iszero
        /* "TokenNetwork200.sol":9610:9664  require(participant1_withdraw <= participant1_deposit) */
      tag_110
      jumpi
      0x0
      dup1
      revert
    tag_110:
        /* "TokenNetwork200.sol":9682:9727  participant2_withdraw <= participant2_deposit */
      dup11
      dup10
      gt
      iszero
        /* "TokenNetwork200.sol":9674:9728  require(participant2_withdraw <= participant2_deposit) */
      tag_111
      jumpi
      0x0
      dup1
      revert
    tag_111:
        /* "TokenNetwork200.sol":9767:9811  participant1_deposit - participant1_withdraw */
      dup10
      dup13
      sub
        /* "TokenNetwork200.sol":9738:9811  participant1_state.deposit = participant1_deposit - participant1_withdraw */
      dup3
      sstore
        /* "TokenNetwork200.sol":9850:9894  participant2_deposit - participant2_withdraw */
      dup9
      dup12
      sub
        /* "TokenNetwork200.sol":9821:9894  participant2_state.deposit = participant2_deposit - participant2_withdraw */
      dup2
      sstore
        /* "TokenNetwork200.sol":9977:10022  channel.open_blocknumber=uint64(block.number) */
      dup4
      sload
      0xffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff
      and
      0x10000000000000000
        /* "TokenNetwork200.sol":10009:10021  block.number */
      number
        /* "TokenNetwork200.sol":9977:10022  channel.open_blocknumber=uint64(block.number) */
      0xffffffffffffffff
      and
      mul
      or
      dup5
      sstore
        /* "TokenNetwork200.sol":10038:10163  Channelwithdraw(channel_identifier, participant1_deposit, participant2_deposit, participant1_withdraw, participant2_withdraw) */
      0x40
      dup1
      mload
      dup7
      dup2
      mstore
      0x20
      dup2
      add
      dup15
      swap1
      mstore
      dup1
      dup3
      add
      dup14
      swap1
      mstore
      0x60
      dup2
      add
      dup13
      swap1
      mstore
      0x80
      dup2
      add
      dup12
      swap1
      mstore
      swap1
      mload
      0xddcd9a7ecf9971deae217332ae11d54f4a7da25e414fb1401fa9fc63c143c264
      swap2
      0xa0
      swap1
      dup3
      swap1
      sub
      add
      swap1
      log1
        /* "TokenNetwork200.sol":7168:10171  function withDraw(... */
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":806:843  SecretRegistry public secret_registry */
    tag_26:
      and(0xffffffffffffffffffffffffffffffffffffffff, sload(0x1))
      dup2
      jump	// out
        /* "TokenNetwork200.sol":19476:23609  function settleChannel(... */
    tag_29:
        /* "TokenNetwork200.sol":19840:19867  uint256 participant1_amount */
      0x0
        /* "TokenNetwork200.sol":19877:19898  uint256 total_deposit */
      dup1
        /* "TokenNetwork200.sol":19908:19934  bytes32 channel_identifier */
      0x0
        /* "TokenNetwork200.sol":20020:20043  Channel storage channel */
      dup1
        /* "TokenNetwork200.sol":20307:20345  Participant storage participant1_state */
      0x0
        /* "TokenNetwork200.sol":20392:20430  Participant storage participant2_state */
      dup1
        /* "TokenNetwork200.sol":19963:20010  getChannelIdentifier(participant1,participant2) */
      tag_113
        /* "TokenNetwork200.sol":19984:19996  participant1 */
      dup15
        /* "TokenNetwork200.sol":19997:20009  participant2 */
      dup15
        /* "TokenNetwork200.sol":19963:19983  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":19963:20010  getChannelIdentifier(participant1,participant2) */
      jump	// in
    tag_113:
        /* "TokenNetwork200.sol":20046:20074  channels[channel_identifier] */
      0x0
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":20046:20054  channels */
      0x3
        /* "TokenNetwork200.sol":20046:20074  channels[channel_identifier] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":20126:20139  channel.state */
      dup1
      sload
        /* "TokenNetwork200.sol":20046:20074  channels[channel_identifier] */
      swap2
      swap6
      pop
      swap4
      pop
        /* "TokenNetwork200.sol":20126:20139  channel.state */
      0x100000000000000000000000000000000
      swap1
      div
      0xff
      and
        /* "TokenNetwork200.sol":20143:20144  2 */
      0x2
        /* "TokenNetwork200.sol":20126:20144  channel.state == 2 */
      eq
        /* "TokenNetwork200.sol":20118:20145  require(channel.state == 2) */
      tag_114
      jumpi
      0x0
      dup1
      revert
    tag_114:
        /* "TokenNetwork200.sol":20206:20233  channel.settle_block_number */
      dup3
      sload
        /* "TokenNetwork200.sol":20256:20268  block.number */
      number
        /* "TokenNetwork200.sol":20206:20233  channel.settle_block_number */
      0xffffffffffffffff
      swap2
      dup3
      and
        /* "TokenNetwork200.sol":892:894  10 */
      0xa
        /* "TokenNetwork200.sol":20206:20253  channel.settle_block_number+punish_block_number */
      add
        /* "TokenNetwork200.sol":20206:20268  channel.settle_block_number+punish_block_number < block.number */
      swap1
      swap2
      and
      lt
        /* "TokenNetwork200.sol":20198:20269  require(channel.settle_block_number+punish_block_number < block.number) */
      tag_115
      jumpi
      0x0
      dup1
      revert
    tag_115:
      pop
      pop
        /* "TokenNetwork200.sol":20348:20382  channel.participants[participant1] */
      0xffffffffffffffffffffffffffffffffffffffff
      dup1
      dup14
      and
      0x0
      swap1
      dup2
      mstore
        /* "TokenNetwork200.sol":20348:20368  channel.participants */
      0x1
      dup4
      add
        /* "TokenNetwork200.sol":20348:20382  channel.participants[participant1] */
      0x20
      mstore
      0x40
      dup1
      dup3
      keccak256
        /* "TokenNetwork200.sol":20433:20467  channel.participants[participant2] */
      swap3
      dup15
      and
      dup3
      mstore
      swap1
      keccak256
        /* "TokenNetwork200.sol":20582:20673  calceBalanceHash(participant1_nonce,participant1_locksroot,participant1_transferred_amount) */
      tag_116
        /* "TokenNetwork200.sol":20599:20617  participant1_nonce */
      dup9
        /* "TokenNetwork200.sol":20618:20640  participant1_locksroot */
      dup12
        /* "TokenNetwork200.sol":20641:20672  participant1_transferred_amount */
      dup15
        /* "TokenNetwork200.sol":20582:20598  calceBalanceHash */
      tag_117
        /* "TokenNetwork200.sol":20582:20673  calceBalanceHash(participant1_nonce,participant1_locksroot,participant1_transferred_amount) */
      jump	// in
    tag_116:
        /* "TokenNetwork200.sol":20549:20580  participant1_state.balance_hash */
      0x1
      dup4
      add
      sload
        /* "TokenNetwork200.sol":20549:20673  participant1_state.balance_hash==calceBalanceHash(participant1_nonce,participant1_locksroot,participant1_transferred_amount) */
      eq
        /* "TokenNetwork200.sol":20541:20674  require(participant1_state.balance_hash==calceBalanceHash(participant1_nonce,participant1_locksroot,participant1_transferred_amount)) */
      tag_118
      jumpi
      0x0
      dup1
      revert
    tag_118:
        /* "TokenNetwork200.sol":20725:20816  calceBalanceHash(participant2_nonce,participant2_locksroot,participant2_transferred_amount) */
      tag_119
        /* "TokenNetwork200.sol":20742:20760  participant2_nonce */
      dup8
        /* "TokenNetwork200.sol":20761:20783  participant2_locksroot */
      dup11
        /* "TokenNetwork200.sol":20784:20815  participant2_transferred_amount */
      dup14
        /* "TokenNetwork200.sol":20725:20741  calceBalanceHash */
      tag_117
        /* "TokenNetwork200.sol":20725:20816  calceBalanceHash(participant2_nonce,participant2_locksroot,participant2_transferred_amount) */
      jump	// in
    tag_119:
        /* "TokenNetwork200.sol":20692:20723  participant2_state.balance_hash */
      0x1
      dup3
      add
      sload
        /* "TokenNetwork200.sol":20692:20816  participant2_state.balance_hash==calceBalanceHash(participant2_nonce,participant2_locksroot,participant2_transferred_amount) */
      eq
        /* "TokenNetwork200.sol":20684:20817  require(participant2_state.balance_hash==calceBalanceHash(participant2_nonce,participant2_locksroot,participant2_transferred_amount)) */
      tag_120
      jumpi
      0x0
      dup1
      revert
    tag_120:
        /* "TokenNetwork200.sol":20873:20899  participant2_state.deposit */
      dup1
      sload
        /* "TokenNetwork200.sol":20844:20870  participant1_state.deposit */
      dup3
      sload
        /* "TokenNetwork200.sol":20942:21010  participant1_state.deposit... */
      dup13
      dup2
      add
        /* "TokenNetwork200.sol":20942:21052  participant1_state.deposit... */
      dup15
      dup2
      sub
      swap9
      pop
        /* "TokenNetwork200.sol":20844:20899  participant1_state.deposit + participant2_state.deposit */
      swap2
      add
      swap6
      pop
        /* "TokenNetwork200.sol":21938:22034  (participant1_state.deposit + participant2_transferred_amount) < participant1_transferred_amount */
      dup13
      gt
        /* "TokenNetwork200.sol":21921:22093  if (... */
      iszero
      tag_121
      jumpi
        /* "TokenNetwork200.sol":22081:22082  0 */
      0x0
        /* "TokenNetwork200.sol":22059:22082  participant1_amount = 0 */
      swap6
      pop
        /* "TokenNetwork200.sol":21921:22093  if (... */
    tag_121:
        /* "TokenNetwork200.sol":22504:22543  min(participant1_amount, total_deposit) */
      tag_122
        /* "TokenNetwork200.sol":22508:22527  participant1_amount */
      dup7
        /* "TokenNetwork200.sol":22529:22542  total_deposit */
      dup7
        /* "TokenNetwork200.sol":22504:22507  min */
      tag_123
        /* "TokenNetwork200.sol":22504:22543  min(participant1_amount, total_deposit) */
      jump	// in
    tag_122:
        /* "TokenNetwork200.sol":22482:22543  participant1_amount = min(participant1_amount, total_deposit) */
      swap6
      pop
        /* "TokenNetwork200.sol":22758:22777  participant1_amount */
      dup6
        /* "TokenNetwork200.sol":22742:22755  total_deposit */
      dup6
        /* "TokenNetwork200.sol":22742:22777  total_deposit - participant1_amount */
      sub
        /* "TokenNetwork200.sol":22708:22777  participant2_transferred_amount = total_deposit - participant1_amount */
      swap11
      pop
        /* "TokenNetwork200.sol":23014:23021  channel */
      dup3
        /* "TokenNetwork200.sol":23014:23034  channel.participants */
      0x1
      add
        /* "TokenNetwork200.sol":23014:23048  channel.participants[participant1] */
      0x0
        /* "TokenNetwork200.sol":23035:23047  participant1 */
      dup16
        /* "TokenNetwork200.sol":23014:23048  channel.participants[participant1] */
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
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
      0x0
        /* "TokenNetwork200.sol":23007:23048  delete channel.participants[participant1] */
      dup1
      dup3
      add
      0x0
      swap1
      sstore
      0x1
      dup3
      add
      0x0
      swap1
      sstore
      pop
      pop
        /* "TokenNetwork200.sol":23065:23072  channel */
      dup3
        /* "TokenNetwork200.sol":23065:23085  channel.participants */
      0x1
      add
        /* "TokenNetwork200.sol":23065:23099  channel.participants[participant2] */
      0x0
        /* "TokenNetwork200.sol":23086:23098  participant2 */
      dup15
        /* "TokenNetwork200.sol":23065:23099  channel.participants[participant2] */
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
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
      0x0
        /* "TokenNetwork200.sol":23058:23099  delete channel.participants[participant2] */
      dup1
      dup3
      add
      0x0
      swap1
      sstore
      0x1
      dup3
      add
      0x0
      swap1
      sstore
      pop
      pop
        /* "TokenNetwork200.sol":23116:23124  channels */
      0x3
        /* "TokenNetwork200.sol":23116:23144  channels[channel_identifier] */
      0x0
        /* "TokenNetwork200.sol":23125:23143  channel_identifier */
      dup6
        /* "TokenNetwork200.sol":23116:23144  channels[channel_identifier] */
      not(0x0)
      and
      not(0x0)
      and
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
      0x0
        /* "TokenNetwork200.sol":23109:23144  delete channels[channel_identifier] */
      dup1
      dup3
      add
      exp(0x100, 0x0)
      dup2
      sload
      swap1
      0xffffffffffffffff
      mul
      not
      and
      swap1
      sstore
      0x0
      dup3
      add
      exp(0x100, 0x8)
      dup2
      sload
      swap1
      0xffffffffffffffff
      mul
      not
      and
      swap1
      sstore
      0x0
      dup3
      add
      exp(0x100, 0x10)
      dup2
      sload
      swap1
      0xff
      mul
      not
      and
      swap1
      sstore
      pop
      pop
        /* "TokenNetwork200.sol":23221:23222  0 */
      0x0
        /* "TokenNetwork200.sol":23199:23218  participant1_amount */
      dup7
        /* "TokenNetwork200.sol":23199:23222  participant1_amount > 0 */
      gt
        /* "TokenNetwork200.sol":23195:23307  if (participant1_amount > 0) {... */
      iszero
      tag_128
      jumpi
        /* "TokenNetwork200.sol":23246:23251  token */
      0x0
      dup1
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffffffffffffffffffffffffffff
      and
        /* "TokenNetwork200.sol":23246:23260  token.transfer */
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xa9059cbb
        /* "TokenNetwork200.sol":23261:23273  participant1 */
      dup16
        /* "TokenNetwork200.sol":23275:23294  participant1_amount */
      dup9
        /* "TokenNetwork200.sol":23246:23295  token.transfer(participant1, participant1_amount) */
      mload(0x40)
      dup4
      0xffffffff
      and
      0x100000000000000000000000000000000000000000000000000000000
      mul
      dup2
      mstore
      0x4
      add
      dup1
      dup4
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      dup2
      mstore
      0x20
      add
      dup3
      dup2
      mstore
      0x20
      add
      swap3
      pop
      pop
      pop
      0x20
      mload(0x40)
      dup1
      dup4
      sub
      dup2
      0x0
      dup8
      dup1
      extcodesize
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_125
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_125:
        /* "TokenNetwork200.sol":23246:23295  token.transfer(participant1, participant1_amount) */
      pop
      gas
      call
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_126
      jumpi
        /* "--CODEGEN--":45:61   */
      returndatasize
        /* "--CODEGEN--":42:43   */
      0x0
        /* "--CODEGEN--":39:40   */
      dup1
        /* "--CODEGEN--":24:62   */
      returndatacopy
        /* "--CODEGEN--":77:93   */
      returndatasize
        /* "--CODEGEN--":74:75   */
      0x0
        /* "--CODEGEN--":67:94   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_126:
        /* "TokenNetwork200.sol":23246:23295  token.transfer(participant1, participant1_amount) */
      pop
      pop
      pop
      pop
      mload(0x40)
      returndatasize
        /* "--CODEGEN--":13:15   */
      0x20
        /* "--CODEGEN--":8:11   */
      dup2
        /* "--CODEGEN--":5:16   */
      lt
        /* "--CODEGEN--":2:4   */
      iszero
      tag_127
      jumpi
        /* "--CODEGEN--":29:30   */
      0x0
        /* "--CODEGEN--":26:27   */
      dup1
        /* "--CODEGEN--":19:31   */
      revert
        /* "--CODEGEN--":2:4   */
    tag_127:
      pop
        /* "TokenNetwork200.sol":23246:23295  token.transfer(participant1, participant1_amount) */
      mload
        /* "TokenNetwork200.sol":23238:23296  require(token.transfer(participant1, participant1_amount)) */
      iszero
      iszero
      tag_128
      jumpi
      0x0
      dup1
      revert
    tag_128:
        /* "TokenNetwork200.sol":23355:23356  0 */
      0x0
        /* "TokenNetwork200.sol":23321:23352  participant2_transferred_amount */
      dup12
        /* "TokenNetwork200.sol":23321:23356  participant2_transferred_amount > 0 */
      gt
        /* "TokenNetwork200.sol":23317:23453  if (participant2_transferred_amount > 0) {... */
      iszero
      tag_133
      jumpi
        /* "TokenNetwork200.sol":23380:23385  token */
      0x0
      dup1
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffffffffffffffffffffffffffff
      and
        /* "TokenNetwork200.sol":23380:23394  token.transfer */
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xa9059cbb
        /* "TokenNetwork200.sol":23395:23407  participant2 */
      dup15
        /* "TokenNetwork200.sol":23409:23440  participant2_transferred_amount */
      dup14
        /* "TokenNetwork200.sol":23380:23441  token.transfer(participant2, participant2_transferred_amount) */
      mload(0x40)
      dup4
      0xffffffff
      and
      0x100000000000000000000000000000000000000000000000000000000
      mul
      dup2
      mstore
      0x4
      add
      dup1
      dup4
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      dup2
      mstore
      0x20
      add
      dup3
      dup2
      mstore
      0x20
      add
      swap3
      pop
      pop
      pop
      0x20
      mload(0x40)
      dup1
      dup4
      sub
      dup2
      0x0
      dup8
      dup1
      extcodesize
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_130
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_130:
        /* "TokenNetwork200.sol":23380:23441  token.transfer(participant2, participant2_transferred_amount) */
      pop
      gas
      call
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_131
      jumpi
        /* "--CODEGEN--":45:61   */
      returndatasize
        /* "--CODEGEN--":42:43   */
      0x0
        /* "--CODEGEN--":39:40   */
      dup1
        /* "--CODEGEN--":24:62   */
      returndatacopy
        /* "--CODEGEN--":77:93   */
      returndatasize
        /* "--CODEGEN--":74:75   */
      0x0
        /* "--CODEGEN--":67:94   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_131:
        /* "TokenNetwork200.sol":23380:23441  token.transfer(participant2, participant2_transferred_amount) */
      pop
      pop
      pop
      pop
      mload(0x40)
      returndatasize
        /* "--CODEGEN--":13:15   */
      0x20
        /* "--CODEGEN--":8:11   */
      dup2
        /* "--CODEGEN--":5:16   */
      lt
        /* "--CODEGEN--":2:4   */
      iszero
      tag_132
      jumpi
        /* "--CODEGEN--":29:30   */
      0x0
        /* "--CODEGEN--":26:27   */
      dup1
        /* "--CODEGEN--":19:31   */
      revert
        /* "--CODEGEN--":2:4   */
    tag_132:
      pop
        /* "TokenNetwork200.sol":23380:23441  token.transfer(participant2, participant2_transferred_amount) */
      mload
        /* "TokenNetwork200.sol":23372:23442  require(token.transfer(participant2, participant2_transferred_amount)) */
      iszero
      iszero
      tag_133
      jumpi
      0x0
      dup1
      revert
    tag_133:
        /* "TokenNetwork200.sol":23468:23602  ChannelSettled(... */
      0x40
      dup1
      mload
      dup8
      dup2
      mstore
      0x20
      dup2
      add
      dup14
      swap1
      mstore
      dup2
      mload
        /* "TokenNetwork200.sol":23496:23514  channel_identifier */
      dup7
      swap3
        /* "TokenNetwork200.sol":23468:23602  ChannelSettled(... */
      0xf94fb5c0628a82dc90648e8dc5e983f632633b0d26603d64e8cc042ca0790aa4
      swap3
      dup3
      swap1
      sub
      add
      swap1
      log2
        /* "TokenNetwork200.sol":19476:23609  function settleChannel(... */
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":996:1019  uint256 public chain_id */
    tag_32:
      sload(0x2)
      dup2
      jump	// out
        /* "TokenNetwork200.sol":12122:13929  function updateBalanceProof(... */
    tag_35:
        /* "TokenNetwork200.sol":12522:12548  bytes32 channel_identifier */
      0x0
        /* "TokenNetwork200.sol":12629:12652  Channel storage channel */
      dup1
        /* "TokenNetwork200.sol":12577:12619  getChannelIdentifier(participant, partner) */
      tag_135
        /* "TokenNetwork200.sol":12598:12609  participant */
      dup14
        /* "TokenNetwork200.sol":12611:12618  partner */
      dup14
        /* "TokenNetwork200.sol":12577:12597  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":12577:12619  getChannelIdentifier(participant, partner) */
      jump	// in
    tag_135:
        /* "TokenNetwork200.sol":12655:12683  channels[channel_identifier] */
      0x0
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":12655:12663  channels */
      0x3
        /* "TokenNetwork200.sol":12655:12683  channels[channel_identifier] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":12701:12714  channel.state */
      dup1
      sload
        /* "TokenNetwork200.sol":12655:12683  channels[channel_identifier] */
      swap2
      swap4
      pop
      swap2
      pop
        /* "TokenNetwork200.sol":12701:12714  channel.state */
      0x100000000000000000000000000000000
      swap1
      div
      0xff
      and
        /* "TokenNetwork200.sol":12716:12717  2 */
      0x2
        /* "TokenNetwork200.sol":12701:12717  channel.state==2 */
      eq
        /* "TokenNetwork200.sol":12693:12718  require(channel.state==2) */
      tag_136
      jumpi
      0x0
      dup1
      revert
    tag_136:
        /* "TokenNetwork200.sol":12736:12763  channel.settle_block_number */
      dup1
      sload
        /* "TokenNetwork200.sol":12767:12779  block.number */
      number
        /* "TokenNetwork200.sol":12736:12763  channel.settle_block_number */
      0xffffffffffffffff
      swap1
      swap2
      and
        /* "TokenNetwork200.sol":12736:12779  channel.settle_block_number >= block.number */
      lt
      iszero
        /* "TokenNetwork200.sol":12728:12780  require(channel.settle_block_number >= block.number) */
      tag_137
      jumpi
      0x0
      dup1
      revert
    tag_137:
        /* "TokenNetwork200.sol":12806:12807  0 */
      0x0
        /* "TokenNetwork200.sol":12798:12807  nonce > 0 */
      dup10
      gt
        /* "TokenNetwork200.sol":12790:12808  require(nonce > 0) */
      tag_138
      jumpi
      0x0
      dup1
      revert
    tag_138:
        /* "TokenNetwork200.sol":12839:13131  recoverAddressFromBalanceProofUpdateMessage(... */
      tag_139
        /* "TokenNetwork200.sol":12896:12914  channel_identifier */
      dup3
        /* "TokenNetwork200.sol":12928:12946  transferred_amount */
      dup13
        /* "TokenNetwork200.sol":12960:12969  locksroot */
      dup13
        /* "TokenNetwork200.sol":12983:12988  nonce */
      dup13
        /* "TokenNetwork200.sol":13002:13009  channel */
      dup6
        /* "TokenNetwork200.sol":13002:13026  channel.open_blocknumber */
      0x0
      add
      0x8
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffff
      and
        /* "TokenNetwork200.sol":13040:13055  additional_hash */
      dup11
        /* "TokenNetwork200.sol":13069:13090  participant_signature */
      dup11
        /* "TokenNetwork200.sol":13104:13121  partner_signature */
      dup11
        /* "TokenNetwork200.sol":12839:12882  recoverAddressFromBalanceProofUpdateMessage */
      tag_140
        /* "TokenNetwork200.sol":12839:13131  recoverAddressFromBalanceProofUpdateMessage(... */
      jump	// in
    tag_139:
        /* "TokenNetwork200.sol":12828:13131  partner == recoverAddressFromBalanceProofUpdateMessage(... */
      0xffffffffffffffffffffffffffffffffffffffff
      dup14
      dup2
      and
      swap2
      and
      eq
        /* "TokenNetwork200.sol":12820:13132  require(partner == recoverAddressFromBalanceProofUpdateMessage(... */
      tag_141
      jumpi
      0x0
      dup1
      revert
    tag_141:
        /* "TokenNetwork200.sol":13165:13413  recoverAddressFromBalanceProof(... */
      tag_142
        /* "TokenNetwork200.sol":13209:13227  channel_identifier */
      dup3
        /* "TokenNetwork200.sol":13241:13259  transferred_amount */
      dup13
        /* "TokenNetwork200.sol":13273:13282  locksroot */
      dup13
        /* "TokenNetwork200.sol":13296:13301  nonce */
      dup13
        /* "TokenNetwork200.sol":13315:13322  channel */
      dup6
        /* "TokenNetwork200.sol":13315:13339  channel.open_blocknumber */
      0x0
      add
      0x8
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffff
      and
        /* "TokenNetwork200.sol":13353:13368  additional_hash */
      dup11
        /* "TokenNetwork200.sol":13382:13403  participant_signature */
      dup11
        /* "TokenNetwork200.sol":13165:13195  recoverAddressFromBalanceProof */
      tag_143
        /* "TokenNetwork200.sol":13165:13413  recoverAddressFromBalanceProof(... */
      jump	// in
    tag_142:
        /* "TokenNetwork200.sol":13150:13413  participant == recoverAddressFromBalanceProof(... */
      0xffffffffffffffffffffffffffffffffffffffff
      dup15
      dup2
      and
      swap2
      and
      eq
        /* "TokenNetwork200.sol":13142:13414  require(participant == recoverAddressFromBalanceProof(... */
      tag_144
      jumpi
      0x0
      dup1
      revert
    tag_144:
        /* "TokenNetwork200.sol":13493:13605  verifyBalanceHashIsValid( channel_identifier, participant, old_transferred_amount,old_locksroot,old_nonce,nonce) */
      tag_145
        /* "TokenNetwork200.sol":13519:13537  channel_identifier */
      dup3
        /* "TokenNetwork200.sol":13539:13550  participant */
      dup15
        /* "TokenNetwork200.sol":13552:13574  old_transferred_amount */
      dup11
        /* "TokenNetwork200.sol":13575:13588  old_locksroot */
      dup11
        /* "TokenNetwork200.sol":13589:13598  old_nonce */
      dup11
        /* "TokenNetwork200.sol":13599:13604  nonce */
      dup15
        /* "TokenNetwork200.sol":13493:13517  verifyBalanceHashIsValid */
      tag_146
        /* "TokenNetwork200.sol":13493:13605  verifyBalanceHashIsValid( channel_identifier, participant, old_transferred_amount,old_locksroot,old_nonce,nonce) */
      jump	// in
    tag_145:
        /* "TokenNetwork200.sol":13615:13699  updateBalanceHash(channel_identifier,participant,nonce,locksroot,transferred_amount) */
      tag_147
        /* "TokenNetwork200.sol":13633:13651  channel_identifier */
      dup3
        /* "TokenNetwork200.sol":13652:13663  participant */
      dup15
        /* "TokenNetwork200.sol":13664:13669  nonce */
      dup12
        /* "TokenNetwork200.sol":13670:13679  locksroot */
      dup14
        /* "TokenNetwork200.sol":13680:13698  transferred_amount */
      dup16
        /* "TokenNetwork200.sol":13615:13632  updateBalanceHash */
      tag_148
        /* "TokenNetwork200.sol":13615:13699  updateBalanceHash(channel_identifier,participant,nonce,locksroot,transferred_amount) */
      jump	// in
    tag_147:
        /* "TokenNetwork200.sol":13835:13922  BalanceProofUpdated(channel_identifier, participant,nonce,locksroot,transferred_amount) */
      0x40
      dup1
      mload
      0xffffffffffffffffffffffffffffffffffffffff
      dup16
      and
      dup2
      mstore
      0x20
      dup2
      add
      dup12
      swap1
      mstore
      dup1
      dup3
      add
      dup13
      swap1
      mstore
      0x60
      dup2
      add
      dup14
      swap1
      mstore
      swap1
      mload
        /* "TokenNetwork200.sol":13855:13873  channel_identifier */
      dup4
      swap2
        /* "TokenNetwork200.sol":13835:13922  BalanceProofUpdated(channel_identifier, participant,nonce,locksroot,transferred_amount) */
      0x426b08fab60dee0cfad0cca6c15f9c5cf729b718857f18e03880482a81120487
      swap2
      swap1
      dup2
      swap1
      sub
      0x80
      add
      swap1
      log2
        /* "TokenNetwork200.sol":12122:13929  function updateBalanceProof(... */
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
    tag_38:
        /* "Utils.sol":367:371  bool */
      0x0
        /* "Utils.sol":434:463  extcodesize(contract_address) */
      swap1
      extcodesize
        /* "Utils.sol":490:498  size > 0 */
      gt
      swap1
        /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
      jump	// out
        /* "TokenNetwork200.sol":1100:1143  mapping(bytes32 => Channel) public channels */
    tag_41:
      mstore(0x20, 0x3)
      0x0
      swap1
      dup2
      mstore
      0x40
      swap1
      keccak256
      sload
      0xffffffffffffffff
      dup1
      dup3
      and
      swap2
      0x10000000000000000
      dup2
      div
      swap1
      swap2
      and
      swap1
      0x100000000000000000000000000000000
      swap1
      div
      0xff
      and
      dup4
      jump	// out
        /* "TokenNetwork200.sol":508:612  bytes32 constant public invalid_balance_hash=keccak256(uint64(0xffffffffffffffff),bytes32(0),uint256(0)) */
    tag_44:
        /* "TokenNetwork200.sol":553:612  keccak256(uint64(0xffffffffffffffff),bytes32(0),uint256(0)) */
      0x40
      dup1
      mload
      0xffffffffffffffff000000000000000000000000000000000000000000000000
      dup2
      mstore
        /* "TokenNetwork200.sol":598:599  0 */
      0x0
        /* "TokenNetwork200.sol":553:612  keccak256(uint64(0xffffffffffffffff),bytes32(0),uint256(0)) */
      0x8
      dup3
      add
      dup2
      swap1
      mstore
      0x28
      dup3
      add
      mstore
      swap1
      mload
      swap1
      dup2
      swap1
      sub
      0x48
      add
      swap1
      keccak256
        /* "TokenNetwork200.sol":508:612  bytes32 constant public invalid_balance_hash=keccak256(uint64(0xffffffffffffffff),bytes32(0),uint256(0)) */
      dup2
      jump	// out
        /* "TokenNetwork200.sol":23678:26342  function cooperativeSettle(... */
    tag_47:
        /* "TokenNetwork200.sol":23964:23983  address participant */
      0x0
        /* "TokenNetwork200.sol":23993:24024  uint256 total_available_deposit */
      dup1
        /* "TokenNetwork200.sol":24034:24060  bytes32 channel_identifier */
      0x0
        /* "TokenNetwork200.sol":24070:24093  uint64 open_blocknumber */
      dup1
        /* "TokenNetwork200.sol":24198:24221  Channel storage channel */
      0x0
        /* "TokenNetwork200.sol":25120:25158  Participant storage participant1_state */
      dup1
        /* "TokenNetwork200.sol":25213:25251  Participant storage participant2_state */
      0x0
        /* "TokenNetwork200.sol":24124:24188  getChannelIdentifier(participant1_address, participant2_address) */
      tag_151
        /* "TokenNetwork200.sol":24145:24165  participant1_address */
      dup14
        /* "TokenNetwork200.sol":24167:24187  participant2_address */
      dup13
        /* "TokenNetwork200.sol":24124:24144  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":24124:24188  getChannelIdentifier(participant1_address, participant2_address) */
      jump	// in
    tag_151:
        /* "TokenNetwork200.sol":24224:24252  channels[channel_identifier] */
      0x0
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":24224:24232  channels */
      0x3
        /* "TokenNetwork200.sol":24224:24252  channels[channel_identifier] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":24306:24319  channel.state */
      dup1
      sload
        /* "TokenNetwork200.sol":24224:24252  channels[channel_identifier] */
      swap2
      swap7
      pop
      swap4
      pop
        /* "TokenNetwork200.sol":24306:24319  channel.state */
      0x100000000000000000000000000000000
      swap1
      div
      0xff
      and
        /* "TokenNetwork200.sol":24323:24324  1 */
      0x1
        /* "TokenNetwork200.sol":24306:24324  channel.state == 1 */
      eq
        /* "TokenNetwork200.sol":24298:24325  require(channel.state == 1) */
      tag_152
      jumpi
      0x0
      dup1
      revert
    tag_152:
        /* "TokenNetwork200.sol":24353:24377  channel.open_blocknumber */
      dup3
      sload
      0x10000000000000000
      swap1
      div
      0xffffffffffffffff
      and
      swap4
      pop
        /* "TokenNetwork200.sol":24401:24689  recoverAddressFromCooperativeSettleSignature(... */
      tag_153
        /* "TokenNetwork200.sol":24459:24477  channel_identifier */
      dup6
        /* "TokenNetwork200.sol":24491:24511  participant1_address */
      dup15
        /* "TokenNetwork200.sol":24525:24545  participant1_balance */
      dup15
        /* "TokenNetwork200.sol":24559:24579  participant2_address */
      dup15
        /* "TokenNetwork200.sol":24593:24613  participant2_balance */
      dup15
        /* "TokenNetwork200.sol":24353:24377  channel.open_blocknumber */
      dup10
        /* "TokenNetwork200.sol":24657:24679  participant1_signature */
      dup16
        /* "TokenNetwork200.sol":24401:24445  recoverAddressFromCooperativeSettleSignature */
      tag_154
        /* "TokenNetwork200.sol":24401:24689  recoverAddressFromCooperativeSettleSignature(... */
      jump	// in
    tag_153:
        /* "TokenNetwork200.sol":24387:24689  participant = recoverAddressFromCooperativeSettleSignature(... */
      swap7
      pop
        /* "TokenNetwork200.sol":24707:24742  participant1_address == participant */
      0xffffffffffffffffffffffffffffffffffffffff
      dup14
      dup2
      and
      swap1
      dup9
      and
      eq
        /* "TokenNetwork200.sol":24699:24743  require(participant1_address == participant) */
      tag_155
      jumpi
      0x0
      dup1
      revert
    tag_155:
        /* "TokenNetwork200.sol":24767:25055  recoverAddressFromCooperativeSettleSignature(... */
      tag_156
        /* "TokenNetwork200.sol":24825:24843  channel_identifier */
      dup6
        /* "TokenNetwork200.sol":24857:24877  participant1_address */
      dup15
        /* "TokenNetwork200.sol":24891:24911  participant1_balance */
      dup15
        /* "TokenNetwork200.sol":24925:24945  participant2_address */
      dup15
        /* "TokenNetwork200.sol":24959:24979  participant2_balance */
      dup15
        /* "TokenNetwork200.sol":24993:25009  open_blocknumber */
      dup10
        /* "TokenNetwork200.sol":25023:25045  participant2_signature */
      dup15
        /* "TokenNetwork200.sol":24767:24811  recoverAddressFromCooperativeSettleSignature */
      tag_154
        /* "TokenNetwork200.sol":24767:25055  recoverAddressFromCooperativeSettleSignature(... */
      jump	// in
    tag_156:
        /* "TokenNetwork200.sol":24753:25055  participant = recoverAddressFromCooperativeSettleSignature(... */
      swap7
      pop
        /* "TokenNetwork200.sol":25073:25108  participant2_address == participant */
      0xffffffffffffffffffffffffffffffffffffffff
      dup12
      dup2
      and
      swap1
      dup9
      and
      eq
        /* "TokenNetwork200.sol":25065:25109  require(participant2_address == participant) */
      tag_157
      jumpi
      0x0
      dup1
      revert
    tag_157:
      pop
      pop
        /* "TokenNetwork200.sol":25161:25203  channel.participants[participant1_address] */
      0xffffffffffffffffffffffffffffffffffffffff
      dup1
      dup13
      and
      0x0
      swap1
      dup2
      mstore
        /* "TokenNetwork200.sol":25161:25181  channel.participants */
      0x1
      dup1
      dup5
      add
        /* "TokenNetwork200.sol":25161:25203  channel.participants[participant1_address] */
      0x20
      swap1
      dup2
      mstore
      0x40
      dup1
      dup5
      keccak256
        /* "TokenNetwork200.sol":25254:25296  channel.participants[participant2_address] */
      swap5
      dup15
      and
      dup5
      mstore
      dup1
      dup5
      keccak256
        /* "TokenNetwork200.sol":25364:25390  participant2_state.deposit */
      dup1
      sload
        /* "TokenNetwork200.sol":25335:25361  participant1_state.deposit */
      dup7
      sload
        /* "TokenNetwork200.sol":25478:25527  delete channel.participants[participant1_address] */
      dup7
      dup9
      sstore
      dup8
      dup7
      add
      dup8
      swap1
      sstore
        /* "TokenNetwork200.sol":25537:25586  delete channel.participants[participant2_address] */
      dup7
      dup4
      sstore
      swap5
      dup3
      add
      dup7
      swap1
      sstore
        /* "TokenNetwork200.sol":25603:25631  channels[channel_identifier] */
      dup10
      dup7
      mstore
        /* "TokenNetwork200.sol":25603:25611  channels */
      0x3
        /* "TokenNetwork200.sol":25603:25631  channels[channel_identifier] */
      swap1
      swap4
      mstore
      swap1
      dup5
      keccak256
        /* "TokenNetwork200.sol":25596:25631  delete channels[channel_identifier] */
      dup1
      sload
      0xffffffffffffffffffffffffffffff0000000000000000000000000000000000
      and
      swap1
      sstore
        /* "TokenNetwork200.sol":25335:25390  participant1_state.deposit + participant2_state.deposit */
      swap2
      add
      swap7
      pop
        /* "TokenNetwork200.sol":25254:25296  channel.participants[participant2_address] */
      swap1
        /* "TokenNetwork200.sol":25679:25703  participant1_balance > 0 */
      dup13
      gt
        /* "TokenNetwork200.sol":25675:25797  if (participant1_balance > 0) {... */
      iszero
      tag_162
      jumpi
        /* "TokenNetwork200.sol":25727:25732  token */
      0x0
      dup1
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffffffffffffffffffffffffffff
      and
        /* "TokenNetwork200.sol":25727:25741  token.transfer */
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xa9059cbb
        /* "TokenNetwork200.sol":25742:25762  participant1_address */
      dup15
        /* "TokenNetwork200.sol":25764:25784  participant1_balance */
      dup15
        /* "TokenNetwork200.sol":25727:25785  token.transfer(participant1_address, participant1_balance) */
      mload(0x40)
      dup4
      0xffffffff
      and
      0x100000000000000000000000000000000000000000000000000000000
      mul
      dup2
      mstore
      0x4
      add
      dup1
      dup4
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      dup2
      mstore
      0x20
      add
      dup3
      dup2
      mstore
      0x20
      add
      swap3
      pop
      pop
      pop
      0x20
      mload(0x40)
      dup1
      dup4
      sub
      dup2
      0x0
      dup8
      dup1
      extcodesize
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_159
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_159:
        /* "TokenNetwork200.sol":25727:25785  token.transfer(participant1_address, participant1_balance) */
      pop
      gas
      call
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_160
      jumpi
        /* "--CODEGEN--":45:61   */
      returndatasize
        /* "--CODEGEN--":42:43   */
      0x0
        /* "--CODEGEN--":39:40   */
      dup1
        /* "--CODEGEN--":24:62   */
      returndatacopy
        /* "--CODEGEN--":77:93   */
      returndatasize
        /* "--CODEGEN--":74:75   */
      0x0
        /* "--CODEGEN--":67:94   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_160:
        /* "TokenNetwork200.sol":25727:25785  token.transfer(participant1_address, participant1_balance) */
      pop
      pop
      pop
      pop
      mload(0x40)
      returndatasize
        /* "--CODEGEN--":13:15   */
      0x20
        /* "--CODEGEN--":8:11   */
      dup2
        /* "--CODEGEN--":5:16   */
      lt
        /* "--CODEGEN--":2:4   */
      iszero
      tag_161
      jumpi
        /* "--CODEGEN--":29:30   */
      0x0
        /* "--CODEGEN--":26:27   */
      dup1
        /* "--CODEGEN--":19:31   */
      revert
        /* "--CODEGEN--":2:4   */
    tag_161:
      pop
        /* "TokenNetwork200.sol":25727:25785  token.transfer(participant1_address, participant1_balance) */
      mload
        /* "TokenNetwork200.sol":25719:25786  require(token.transfer(participant1_address, participant1_balance)) */
      iszero
      iszero
      tag_162
      jumpi
      0x0
      dup1
      revert
    tag_162:
        /* "TokenNetwork200.sol":25834:25835  0 */
      0x0
        /* "TokenNetwork200.sol":25811:25831  participant2_balance */
      dup11
        /* "TokenNetwork200.sol":25811:25835  participant2_balance > 0 */
      gt
        /* "TokenNetwork200.sol":25807:25929  if (participant2_balance > 0) {... */
      iszero
      tag_167
      jumpi
        /* "TokenNetwork200.sol":25859:25864  token */
      0x0
      dup1
      sload
        /* "TokenNetwork200.sol":25859:25917  token.transfer(participant2_address, participant2_balance) */
      0x40
      dup1
      mload
      0xa9059cbb00000000000000000000000000000000000000000000000000000000
      dup2
      mstore
        /* "TokenNetwork200.sol":25859:25864  token */
      0xffffffffffffffffffffffffffffffffffffffff
        /* "TokenNetwork200.sol":25859:25917  token.transfer(participant2_address, participant2_balance) */
      dup16
      dup2
      and
      0x4
      dup4
      add
      mstore
      0x24
      dup3
      add
      dup16
      swap1
      mstore
      swap2
      mload
        /* "TokenNetwork200.sol":25859:25864  token */
      swap2
      swap1
      swap3
      and
      swap3
        /* "TokenNetwork200.sol":25859:25873  token.transfer */
      0xa9059cbb
      swap3
        /* "TokenNetwork200.sol":25859:25917  token.transfer(participant2_address, participant2_balance) */
      0x44
      dup1
      dup3
      add
      swap4
      0x20
      swap4
      swap1
      swap3
      dup4
      swap1
      sub
      swap1
      swap2
      add
      swap1
      dup3
      swap1
        /* "TokenNetwork200.sol":25859:25864  token */
      dup8
        /* "TokenNetwork200.sol":25859:25917  token.transfer(participant2_address, participant2_balance) */
      dup1
      extcodesize
      iszero
        /* "--CODEGEN--":5:7   */
      dup1
      iszero
      tag_164
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_164:
        /* "TokenNetwork200.sol":25859:25917  token.transfer(participant2_address, participant2_balance) */
      pop
      gas
      call
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_165
      jumpi
        /* "--CODEGEN--":45:61   */
      returndatasize
        /* "--CODEGEN--":42:43   */
      0x0
        /* "--CODEGEN--":39:40   */
      dup1
        /* "--CODEGEN--":24:62   */
      returndatacopy
        /* "--CODEGEN--":77:93   */
      returndatasize
        /* "--CODEGEN--":74:75   */
      0x0
        /* "--CODEGEN--":67:94   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_165:
        /* "TokenNetwork200.sol":25859:25917  token.transfer(participant2_address, participant2_balance) */
      pop
      pop
      pop
      pop
      mload(0x40)
      returndatasize
        /* "--CODEGEN--":13:15   */
      0x20
        /* "--CODEGEN--":8:11   */
      dup2
        /* "--CODEGEN--":5:16   */
      lt
        /* "--CODEGEN--":2:4   */
      iszero
      tag_166
      jumpi
        /* "--CODEGEN--":29:30   */
      0x0
        /* "--CODEGEN--":26:27   */
      dup1
        /* "--CODEGEN--":19:31   */
      revert
        /* "--CODEGEN--":2:4   */
    tag_166:
      pop
        /* "TokenNetwork200.sol":25859:25917  token.transfer(participant2_address, participant2_balance) */
      mload
        /* "TokenNetwork200.sol":25851:25918  require(token.transfer(participant2_address, participant2_balance)) */
      iszero
      iszero
      tag_167
      jumpi
      0x0
      dup1
      revert
    tag_167:
        /* "TokenNetwork200.sol":26065:26108  participant1_balance + participant2_balance */
      dup12
      dup11
      add
        /* "TokenNetwork200.sol":26037:26109  total_available_deposit == (participant1_balance + participant2_balance) */
      dup7
      eq
        /* "TokenNetwork200.sol":26029:26110  require(total_available_deposit == (participant1_balance + participant2_balance)) */
      tag_168
      jumpi
      0x0
      dup1
      revert
    tag_168:
        /* "TokenNetwork200.sol":26128:26175  total_available_deposit >= participant1_balance */
      dup12
      dup7
      lt
      iszero
        /* "TokenNetwork200.sol":26120:26176  require(total_available_deposit >= participant1_balance) */
      tag_169
      jumpi
      0x0
      dup1
      revert
    tag_169:
        /* "TokenNetwork200.sol":26194:26241  total_available_deposit >= participant2_balance */
      dup10
      dup7
      lt
      iszero
        /* "TokenNetwork200.sol":26186:26242  require(total_available_deposit >= participant2_balance) */
      tag_170
      jumpi
      0x0
      dup1
      revert
    tag_170:
        /* "TokenNetwork200.sol":26257:26335  ChannelSettled(channel_identifier, participant1_balance, participant2_balance) */
      0x40
      dup1
      mload
      dup14
      dup2
      mstore
      0x20
      dup2
      add
      dup13
      swap1
      mstore
      dup2
      mload
        /* "TokenNetwork200.sol":26272:26290  channel_identifier */
      dup8
      swap3
        /* "TokenNetwork200.sol":26257:26335  ChannelSettled(channel_identifier, participant1_balance, participant2_balance) */
      0xf94fb5c0628a82dc90648e8dc5e983f632633b0d26603d64e8cc042ca0790aa4
      swap3
      dup3
      swap1
      sub
      add
      swap1
      log2
        /* "TokenNetwork200.sol":23678:26342  function cooperativeSettle(... */
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":17485:19408  function punishObsoleteUnlock(... */
    tag_50:
        /* "TokenNetwork200.sol":17790:17816  bytes32 channel_identifier */
      0x0
        /* "TokenNetwork200.sol":17826:17843  bytes32 locksroot */
      dup1
        /* "TokenNetwork200.sol":17853:17873  bytes32 balance_hash */
      0x0
        /* "TokenNetwork200.sol":17953:17976  Channel storage channel */
      dup1
        /* "TokenNetwork200.sol":18052:18089  Participant storage beneficiary_state */
      0x0
        /* "TokenNetwork200.sol":18781:18814  Participant storage cheater_state */
      dup1
        /* "TokenNetwork200.sol":17902:17943  getChannelIdentifier(beneficiary,cheater) */
      tag_172
        /* "TokenNetwork200.sol":17923:17934  beneficiary */
      dup15
        /* "TokenNetwork200.sol":17935:17942  cheater */
      dup15
        /* "TokenNetwork200.sol":17902:17922  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":17902:17943  getChannelIdentifier(beneficiary,cheater) */
      jump	// in
    tag_172:
        /* "TokenNetwork200.sol":17979:18007  channels[channel_identifier] */
      0x0
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":17979:17987  channels */
      0x3
        /* "TokenNetwork200.sol":17979:18007  channels[channel_identifier] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":18025:18038  channel.state */
      dup1
      sload
        /* "TokenNetwork200.sol":17979:18007  channels[channel_identifier] */
      swap2
      swap8
      pop
      swap4
      pop
        /* "TokenNetwork200.sol":18025:18038  channel.state */
      0x100000000000000000000000000000000
      swap1
      div
      0xff
      and
        /* "TokenNetwork200.sol":18040:18041  2 */
      0x2
        /* "TokenNetwork200.sol":18025:18041  channel.state==2 */
      eq
        /* "TokenNetwork200.sol":18017:18042  require(channel.state==2) */
      tag_173
      jumpi
      0x0
      dup1
      revert
    tag_173:
        /* "TokenNetwork200.sol":18092:18125  channel.participants[beneficiary] */
      0xffffffffffffffffffffffffffffffffffffffff
      dup15
      and
      0x0
      swap1
      dup2
      mstore
        /* "TokenNetwork200.sol":18092:18112  channel.participants */
      0x1
      dup1
      dup6
      add
        /* "TokenNetwork200.sol":18092:18125  channel.participants[beneficiary] */
      0x20
      mstore
      0x40
      swap1
      swap2
      keccak256
        /* "TokenNetwork200.sol":18148:18178  beneficiary_state.balance_hash */
      swap1
      dup2
      add
      sload
      swap5
      pop
        /* "TokenNetwork200.sol":18092:18125  channel.participants[beneficiary] */
      swap2
      pop
        /* "TokenNetwork200.sol":18320:18337  balance_hash != 0 */
      dup4
      iszero
      iszero
        /* "TokenNetwork200.sol":18312:18338  require(balance_hash != 0) */
      tag_174
      jumpi
      0x0
      dup1
      revert
    tag_174:
        /* "TokenNetwork200.sol":18683:18707  channel.open_blocknumber */
      dup3
      sload
        /* "TokenNetwork200.sol":18579:18770  recoverAddressFromUnlockProof(... */
      tag_175
      swap1
        /* "TokenNetwork200.sol":18622:18640  channel_identifier */
      dup8
      swap1
        /* "TokenNetwork200.sol":18654:18662  lockhash */
      dup15
      swap1
        /* "TokenNetwork200.sol":18683:18707  channel.open_blocknumber */
      0x10000000000000000
      swap1
      div
      0xffffffffffffffff
      and
        /* "TokenNetwork200.sol":18722:18737  additional_hash */
      dup13
        /* "TokenNetwork200.sol":18751:18760  signature */
      dup13
        /* "TokenNetwork200.sol":18579:18608  recoverAddressFromUnlockProof */
      tag_176
        /* "TokenNetwork200.sol":18579:18770  recoverAddressFromUnlockProof(... */
      jump	// in
    tag_175:
        /* "TokenNetwork200.sol":18568:18770  cheater == recoverAddressFromUnlockProof(... */
      0xffffffffffffffffffffffffffffffffffffffff
      dup15
      dup2
      and
      swap2
      and
      eq
        /* "TokenNetwork200.sol":18560:18771  require(cheater == recoverAddressFromUnlockProof(... */
      tag_177
      jumpi
      0x0
      dup1
      revert
    tag_177:
      pop
        /* "TokenNetwork200.sol":18817:18846  channel.participants[cheater] */
      0xffffffffffffffffffffffffffffffffffffffff
      dup13
      and
      0x0
      swap1
      dup2
      mstore
        /* "TokenNetwork200.sol":18817:18837  channel.participants */
      0x1
      dup4
      add
        /* "TokenNetwork200.sol":18817:18846  channel.participants[cheater] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":19023:19063  computeMerkleRoot(lockhash,merkle_proof) */
      tag_178
        /* "TokenNetwork200.sol":19041:19049  lockhash */
      dup13
        /* "TokenNetwork200.sol":19050:19062  merkle_proof */
      dup9
        /* "TokenNetwork200.sol":19023:19040  computeMerkleRoot */
      tag_179
        /* "TokenNetwork200.sol":19023:19063  computeMerkleRoot(lockhash,merkle_proof) */
      jump	// in
    tag_178:
        /* "TokenNetwork200.sol":19013:19063  locksroot=computeMerkleRoot(lockhash,merkle_proof) */
      swap5
      pop
        /* "TokenNetwork200.sol":19095:19171  calceBalanceHash(beneficiary_nonce,locksroot,beneficiary_transferred_amount) */
      tag_180
        /* "TokenNetwork200.sol":19112:19129  beneficiary_nonce */
      dup11
        /* "TokenNetwork200.sol":19130:19139  locksroot */
      dup7
        /* "TokenNetwork200.sol":19140:19170  beneficiary_transferred_amount */
      dup14
        /* "TokenNetwork200.sol":19095:19111  calceBalanceHash */
      tag_117
        /* "TokenNetwork200.sol":19095:19171  calceBalanceHash(beneficiary_nonce,locksroot,beneficiary_transferred_amount) */
      jump	// in
    tag_180:
        /* "TokenNetwork200.sol":19081:19171  balance_hash==calceBalanceHash(beneficiary_nonce,locksroot,beneficiary_transferred_amount) */
      dup5
      eq
        /* "TokenNetwork200.sol":19073:19172  require(balance_hash==calceBalanceHash(beneficiary_nonce,locksroot,beneficiary_transferred_amount)) */
      tag_181
      jumpi
      0x0
      dup1
      revert
    tag_181:
        /* "TokenNetwork200.sol":553:612  keccak256(uint64(0xffffffffffffffff),bytes32(0),uint256(0)) */
      0x40
      dup1
      mload
      0xffffffffffffffff000000000000000000000000000000000000000000000000
      dup2
      mstore
        /* "TokenNetwork200.sol":598:599  0 */
      0x0
        /* "TokenNetwork200.sol":553:612  keccak256(uint64(0xffffffffffffffff),bytes32(0),uint256(0)) */
      0x8
      dup3
      add
      dup2
      swap1
      mstore
      0x28
      dup3
      add
      dup2
      swap1
      mstore
      swap2
      mload
      swap1
      dup2
      swap1
      sub
      0x48
      add
      swap1
      keccak256
        /* "TokenNetwork200.sol":590:600  bytes32(0) */
      0x1
        /* "TokenNetwork200.sol":19232:19262  beneficiary_state.balance_hash */
      dup5
      add
        /* "TokenNetwork200.sol":19232:19285  beneficiary_state.balance_hash = invalid_balance_hash */
      sstore
        /* "TokenNetwork200.sol":19347:19368  cheater_state.deposit */
      dup2
      sload
        /* "TokenNetwork200.sol":19321:19346  beneficiary_state.deposit */
      dup4
      sload
        /* "TokenNetwork200.sol":19321:19368  beneficiary_state.deposit+cheater_state.deposit */
      add
        /* "TokenNetwork200.sol":19295:19368  beneficiary_state.deposit=beneficiary_state.deposit+cheater_state.deposit */
      swap1
      swap3
      sstore
        /* "TokenNetwork200.sol":19378:19401  cheater_state.deposit=0 */
      sstore
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
        /* "TokenNetwork200.sol":17485:19408  function punishObsoleteUnlock(... */
      jump	// out
        /* "TokenNetwork200.sol":849:894  uint64 constant public punish_block_number=10 */
    tag_53:
        /* "TokenNetwork200.sol":892:894  10 */
      0xa
        /* "TokenNetwork200.sol":849:894  uint64 constant public punish_block_number=10 */
      dup2
      jump	// out
        /* "TokenNetwork200.sol":28519:28983  function getChannelParticipantInfo( address participant,address partner)... */
    tag_56:
        /* "TokenNetwork200.sol":28627:28634  uint256 */
      0x0
        /* "TokenNetwork200.sol":28636:28643  bytes32 */
      dup1
        /* "TokenNetwork200.sol":28660:28686  bytes32 channel_identifier */
      0x0
        /* "TokenNetwork200.sol":28738:28761  Channel storage channel */
      dup1
        /* "TokenNetwork200.sol":28800:28837  Participant storage participant_state */
      0x0
        /* "TokenNetwork200.sol":28687:28728  getChannelIdentifier(participant,partner) */
      tag_183
        /* "TokenNetwork200.sol":28708:28719  participant */
      dup8
        /* "TokenNetwork200.sol":28720:28727  partner */
      dup8
        /* "TokenNetwork200.sol":28687:28707  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":28687:28728  getChannelIdentifier(participant,partner) */
      jump	// in
    tag_183:
        /* "TokenNetwork200.sol":28762:28790  channels[channel_identifier] */
      0x0
      swap1
      dup2
      mstore
        /* "TokenNetwork200.sol":28762:28770  channels */
      0x3
        /* "TokenNetwork200.sol":28762:28790  channels[channel_identifier] */
      0x20
      swap1
      dup2
      mstore
      0x40
      dup1
      dup4
      keccak256
        /* "TokenNetwork200.sol":28840:28873  channel.participants[participant] */
      0xffffffffffffffffffffffffffffffffffffffff
      swap11
      swap1
      swap11
      and
      dup4
      mstore
        /* "TokenNetwork200.sol":28840:28860  channel.participants */
      0x1
      swap10
      dup11
      add
        /* "TokenNetwork200.sol":28840:28873  channel.participants[participant] */
      swap1
      swap2
      mstore
      swap1
      keccak256
        /* "TokenNetwork200.sol":28901:28926  participant_state.deposit */
      dup1
      sload
        /* "TokenNetwork200.sol":28936:28966  participant_state.balance_hash */
      swap8
      add
      sload
        /* "TokenNetwork200.sol":28901:28926  participant_state.deposit */
      swap7
      swap8
        /* "TokenNetwork200.sol":28519:28983  function getChannelParticipantInfo( address participant,address partner)... */
      swap6
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":4666:5550  function openChannel(address participant1, address participant2, uint64 settle_timeout)... */
    tag_59:
        /* "TokenNetwork200.sol":4818:4844  bytes32 channel_identifier */
      0x0
        /* "TokenNetwork200.sol":5053:5076  Channel storage channel */
      dup1
        /* "TokenNetwork200.sol":4777:4791  settle_timeout */
      dup3
        /* "TokenNetwork200.sol":3987:3988  6 */
      0x6
        /* "TokenNetwork200.sol":3976:3983  timeout */
      dup2
        /* "TokenNetwork200.sol":3976:3988  timeout >= 6 */
      0xffffffffffffffff
      and
      lt
      iszero
        /* "TokenNetwork200.sol":3976:4010  timeout >= 6 && timeout <= 2700000 */
      dup1
      iszero
      tag_185
      jumpi
      pop
        /* "TokenNetwork200.sol":4003:4010  2700000 */
      0x2932e0
        /* "TokenNetwork200.sol":3992:3999  timeout */
      dup2
        /* "TokenNetwork200.sol":3992:4010  timeout <= 2700000 */
      0xffffffffffffffff
      and
      gt
      iszero
        /* "TokenNetwork200.sol":3976:4010  timeout >= 6 && timeout <= 2700000 */
    tag_185:
        /* "TokenNetwork200.sol":3968:4011  require(timeout >= 6 && timeout <= 2700000) */
      iszero
      iszero
      tag_186
      jumpi
      0x0
      dup1
      revert
    tag_186:
        /* "TokenNetwork200.sol":4862:4881  participant1 != 0x0 */
      0xffffffffffffffffffffffffffffffffffffffff
      dup7
      and
      iszero
      iszero
        /* "TokenNetwork200.sol":4854:4882  require(participant1 != 0x0) */
      tag_188
      jumpi
      0x0
      dup1
      revert
    tag_188:
        /* "TokenNetwork200.sol":4900:4919  participant2 != 0x0 */
      0xffffffffffffffffffffffffffffffffffffffff
      dup6
      and
      iszero
      iszero
        /* "TokenNetwork200.sol":4892:4920  require(participant2 != 0x0) */
      tag_189
      jumpi
      0x0
      dup1
      revert
    tag_189:
        /* "TokenNetwork200.sol":4938:4966  participant1 != participant2 */
      0xffffffffffffffffffffffffffffffffffffffff
      dup7
      dup2
      and
      swap1
      dup7
      and
      eq
      iszero
        /* "TokenNetwork200.sol":4930:4967  require(participant1 != participant2) */
      tag_190
      jumpi
      0x0
      dup1
      revert
    tag_190:
        /* "TokenNetwork200.sol":4996:5043  getChannelIdentifier(participant1,participant2) */
      tag_191
        /* "TokenNetwork200.sol":5017:5029  participant1 */
      dup7
        /* "TokenNetwork200.sol":5030:5042  participant2 */
      dup7
        /* "TokenNetwork200.sol":4996:5016  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":4996:5043  getChannelIdentifier(participant1,participant2) */
      jump	// in
    tag_191:
        /* "TokenNetwork200.sol":5079:5107  channels[channel_identifier] */
      0x0
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":5079:5087  channels */
      0x3
        /* "TokenNetwork200.sol":5079:5107  channels[channel_identifier] */
      0x20
      swap1
      dup2
      mstore
      0x40
      swap2
      dup3
      swap1
      keccak256
        /* "TokenNetwork200.sol":5290:5334  channel.settle_block_number = settle_timeout */
      dup1
      sload
        /* "TokenNetwork200.sol":5433:5450  channel.state = 1 */
      0xffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff
        /* "TokenNetwork200.sol":5376:5388  block.number */
      number
        /* "TokenNetwork200.sol":5290:5334  channel.settle_block_number = settle_timeout */
      0xffffffffffffffff
        /* "TokenNetwork200.sol":5344:5389  channel.open_blocknumber=uint64(block.number) */
      swap1
      dup2
      and
      0x10000000000000000
      mul
      0xffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff
        /* "TokenNetwork200.sol":5290:5334  channel.settle_block_number = settle_timeout */
      swap2
      dup13
      and
      0xffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000
      swap1
      swap5
      and
      dup5
      or
        /* "TokenNetwork200.sol":5344:5389  channel.open_blocknumber=uint64(block.number) */
      swap2
      swap1
      swap2
      and
      or
        /* "TokenNetwork200.sol":5433:5450  channel.state = 1 */
      and
      0x100000000000000000000000000000000
      or
      dup3
      sstore
        /* "TokenNetwork200.sol":5466:5543  ChannelOpened(channel_identifier, participant1, participant2, settle_timeout) */
      dup4
      mload
      0xffffffffffffffffffffffffffffffffffffffff
      dup1
      dup14
      and
      dup3
      mstore
      dup12
      and
      swap4
      dup2
      add
      swap4
      swap1
      swap4
      mstore
      dup3
      dup5
      add
      mstore
      swap2
      mload
        /* "TokenNetwork200.sol":4977:5043  channel_identifier=getChannelIdentifier(participant1,participant2) */
      swap3
      swap6
      pop
        /* "TokenNetwork200.sol":5079:5107  channels[channel_identifier] */
      swap1
      swap4
      pop
        /* "TokenNetwork200.sol":4977:5043  channel_identifier=getChannelIdentifier(participant1,participant2) */
      dup5
      swap2
        /* "TokenNetwork200.sol":5466:5543  ChannelOpened(channel_identifier, participant1, participant2, settle_timeout) */
      0x448d27f1fe12f92a2070111296e68fd6ef0a01c0e05bf5819eda0dbcf267bf3d
      swap2
      dup2
      swap1
      sub
      0x60
      add
      swap1
      log2
        /* "TokenNetwork200.sol":4666:5550  function openChannel(address participant1, address participant2, uint64 settle_timeout)... */
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":453:502  string constant public contract_version = "0.3._" */
    tag_62:
      0x40
      dup1
      mload
      dup1
      dup3
      add
      swap1
      swap2
      mstore
      0x5
      dup2
      mstore
      0x302e332e5f000000000000000000000000000000000000000000000000000000
      0x20
      dup3
      add
      mstore
      dup2
      jump	// out
        /* "TokenNetwork200.sol":5699:7055  function setTotalDeposit(address participant,address partner, uint256 total_deposit)... */
    tag_69:
        /* "TokenNetwork200.sol":5845:5866  uint256 added_deposit */
      0x0
      dup1
      dup1
      dup1
        /* "TokenNetwork200.sol":5817:5834  total_deposit > 0 */
      dup1
      dup6
      gt
        /* "TokenNetwork200.sol":5809:5835  require(total_deposit > 0) */
      tag_193
      jumpi
      0x0
      dup1
      revert
    tag_193:
        /* "TokenNetwork200.sol":5931:5972  getChannelIdentifier(participant,partner) */
      tag_194
        /* "TokenNetwork200.sol":5952:5963  participant */
      dup8
        /* "TokenNetwork200.sol":5964:5971  partner */
      dup8
        /* "TokenNetwork200.sol":5931:5951  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":5931:5972  getChannelIdentifier(participant,partner) */
      jump	// in
    tag_194:
        /* "TokenNetwork200.sol":6008:6036  channels[channel_identifier] */
      0x0
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":6008:6016  channels */
      0x3
        /* "TokenNetwork200.sol":6008:6036  channels[channel_identifier] */
      0x20
      swap1
      dup2
      mstore
      0x40
      dup1
      dup4
      keccak256
        /* "TokenNetwork200.sol":6086:6119  channel.participants[participant] */
      0xffffffffffffffffffffffffffffffffffffffff
      dup13
      and
      dup5
      mstore
        /* "TokenNetwork200.sol":6086:6106  channel.participants */
      0x1
      dup2
      add
        /* "TokenNetwork200.sol":6086:6119  channel.participants[participant] */
      swap1
      swap3
      mstore
      swap1
      swap2
      keccak256
        /* "TokenNetwork200.sol":6137:6162  participant_state.deposit */
      dup1
      sload
        /* "TokenNetwork200.sol":5912:5972  channel_identifier=getChannelIdentifier(participant,partner) */
      swap3
      swap6
      pop
        /* "TokenNetwork200.sol":6008:6036  channels[channel_identifier] */
      swap1
      swap4
      pop
        /* "TokenNetwork200.sol":6086:6119  channel.participants[participant] */
      swap2
      pop
        /* "TokenNetwork200.sol":6137:6178  participant_state.deposit < total_deposit */
      dup6
      gt
        /* "TokenNetwork200.sol":6129:6179  require(participant_state.deposit < total_deposit) */
      tag_195
      jumpi
      0x0
      dup1
      revert
    tag_195:
        /* "TokenNetwork200.sol":6295:6320  participant_state.deposit */
      dup1
      sload
        /* "TokenNetwork200.sol":6279:6320  total_deposit - participant_state.deposit */
      dup1
      dup7
      sub
        /* "TokenNetwork200.sol":6383:6425  participant_state.deposit += added_deposit */
      swap1
      dup2
      add
      dup3
      sstore
        /* "TokenNetwork200.sol":6295:6320  participant_state.deposit */
      0x0
        /* "TokenNetwork200.sol":6472:6477  token */
      dup1
      sload
        /* "TokenNetwork200.sol":6472:6532  token.transferFrom(msg.sender, address(this), added_deposit) */
      0x40
      dup1
      mload
      0x23b872dd00000000000000000000000000000000000000000000000000000000
      dup2
      mstore
        /* "TokenNetwork200.sol":6491:6501  msg.sender */
      caller
        /* "TokenNetwork200.sol":6472:6532  token.transferFrom(msg.sender, address(this), added_deposit) */
      0x4
      dup3
      add
      mstore
        /* "TokenNetwork200.sol":6511:6515  this */
      address
        /* "TokenNetwork200.sol":6472:6532  token.transferFrom(msg.sender, address(this), added_deposit) */
      0x24
      dup3
      add
      mstore
      0x44
      dup2
      add
      dup6
      swap1
      mstore
      swap1
      mload
        /* "TokenNetwork200.sol":6279:6320  total_deposit - participant_state.deposit */
      swap4
      swap8
      pop
        /* "TokenNetwork200.sol":6472:6477  token */
      0xffffffffffffffffffffffffffffffffffffffff
      swap1
      swap2
      and
      swap3
        /* "TokenNetwork200.sol":6472:6490  token.transferFrom */
      0x23b872dd
      swap3
        /* "TokenNetwork200.sol":6472:6532  token.transferFrom(msg.sender, address(this), added_deposit) */
      0x64
      dup1
      dup5
      add
      swap4
      0x20
      swap4
      swap3
      swap1
      dup4
      swap1
      sub
      swap1
      swap2
      add
      swap1
      dup3
      swap1
        /* "TokenNetwork200.sol":6472:6477  token */
      dup8
        /* "TokenNetwork200.sol":6472:6532  token.transferFrom(msg.sender, address(this), added_deposit) */
      dup1
      extcodesize
      iszero
        /* "--CODEGEN--":5:7   */
      dup1
      iszero
      tag_196
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_196:
        /* "TokenNetwork200.sol":6472:6532  token.transferFrom(msg.sender, address(this), added_deposit) */
      pop
      gas
      call
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_197
      jumpi
        /* "--CODEGEN--":45:61   */
      returndatasize
        /* "--CODEGEN--":42:43   */
      0x0
        /* "--CODEGEN--":39:40   */
      dup1
        /* "--CODEGEN--":24:62   */
      returndatacopy
        /* "--CODEGEN--":77:93   */
      returndatasize
        /* "--CODEGEN--":74:75   */
      0x0
        /* "--CODEGEN--":67:94   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_197:
        /* "TokenNetwork200.sol":6472:6532  token.transferFrom(msg.sender, address(this), added_deposit) */
      pop
      pop
      pop
      pop
      mload(0x40)
      returndatasize
        /* "--CODEGEN--":13:15   */
      0x20
        /* "--CODEGEN--":8:11   */
      dup2
        /* "--CODEGEN--":5:16   */
      lt
        /* "--CODEGEN--":2:4   */
      iszero
      tag_198
      jumpi
        /* "--CODEGEN--":29:30   */
      0x0
        /* "--CODEGEN--":26:27   */
      dup1
        /* "--CODEGEN--":19:31   */
      revert
        /* "--CODEGEN--":2:4   */
    tag_198:
      pop
        /* "TokenNetwork200.sol":6472:6532  token.transferFrom(msg.sender, address(this), added_deposit) */
      mload
        /* "TokenNetwork200.sol":6464:6533  require(token.transferFrom(msg.sender, address(this), added_deposit)) */
      iszero
      iszero
      tag_199
      jumpi
      0x0
      dup1
      revert
    tag_199:
        /* "TokenNetwork200.sol":6551:6564  channel.state */
      dup2
      sload
      0x100000000000000000000000000000000
      swap1
      div
      0xff
      and
        /* "TokenNetwork200.sol":6566:6567  1 */
      0x1
        /* "TokenNetwork200.sol":6551:6567  channel.state==1 */
      eq
        /* "TokenNetwork200.sol":6543:6568  require(channel.state==1) */
      tag_200
      jumpi
      0x0
      dup1
      revert
    tag_200:
        /* "TokenNetwork200.sol":6634:6659  participant_state.deposit */
      dup1
      sload
        /* "TokenNetwork200.sol":6583:6660  ChannelNewDeposit(channel_identifier, participant, participant_state.deposit) */
      0x40
      dup1
      mload
      0xffffffffffffffffffffffffffffffffffffffff
      dup11
      and
      dup2
      mstore
      0x20
      dup2
      add
      swap3
      swap1
      swap3
      mstore
      dup1
      mload
        /* "TokenNetwork200.sol":6601:6619  channel_identifier */
      dup6
      swap3
        /* "TokenNetwork200.sol":6583:6660  ChannelNewDeposit(channel_identifier, participant, participant_state.deposit) */
      0x346e981e2bfa2366dc2307a8f1fa24779830a01121b1275fe565c6b98bb4d34
      swap3
      swap1
      dup3
      swap1
      sub
      add
      swap1
      log2
        /* "TokenNetwork200.sol":5699:7055  function setTotalDeposit(address participant,address partner, uint256 total_deposit)... */
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":10282:11927  function closeChannel(... */
    tag_72:
        /* "TokenNetwork200.sol":10504:10530  bytes32 channel_identifier */
      0x0
        /* "TokenNetwork200.sol":10540:10573  address recovered_partner_address */
      dup1
        /* "TokenNetwork200.sol":10652:10675  Channel storage channel */
      0x0
        /* "TokenNetwork200.sol":11292:11325  Participant storage partner_state */
      dup1
        /* "TokenNetwork200.sol":10602:10642  getChannelIdentifier(msg.sender,partner) */
      tag_202
        /* "TokenNetwork200.sol":10623:10633  msg.sender */
      caller
        /* "TokenNetwork200.sol":10634:10641  partner */
      dup12
        /* "TokenNetwork200.sol":10602:10622  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":10602:10642  getChannelIdentifier(msg.sender,partner) */
      jump	// in
    tag_202:
        /* "TokenNetwork200.sol":10678:10706  channels[channel_identifier] */
      0x0
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":10678:10686  channels */
      0x3
        /* "TokenNetwork200.sol":10678:10706  channels[channel_identifier] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":10724:10737  channel.state */
      dup1
      sload
        /* "TokenNetwork200.sol":10678:10706  channels[channel_identifier] */
      swap2
      swap6
      pop
      swap3
      pop
        /* "TokenNetwork200.sol":10724:10737  channel.state */
      0x100000000000000000000000000000000
      swap1
      div
      0xff
      and
        /* "TokenNetwork200.sol":10739:10740  1 */
      0x1
        /* "TokenNetwork200.sol":10724:10740  channel.state==1 */
      eq
        /* "TokenNetwork200.sol":10716:10741  require(channel.state==1) */
      tag_203
      jumpi
      0x0
      dup1
      revert
    tag_203:
        /* "TokenNetwork200.sol":10822:10839  channel.state = 2 */
      dup2
      sload
        /* "TokenNetwork200.sol":10922:10973  channel.settle_block_number += uint64(block.number) */
      0xffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000
        /* "TokenNetwork200.sol":10822:10839  channel.state = 2 */
      0xffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff
      swap1
      swap2
      and
      0x200000000000000000000000000000000
      or
        /* "TokenNetwork200.sol":10922:10973  channel.settle_block_number += uint64(block.number) */
      swap1
      dup2
      and
        /* "TokenNetwork200.sol":10960:10972  block.number */
      number
        /* "TokenNetwork200.sol":10922:10973  channel.settle_block_number += uint64(block.number) */
      0xffffffffffffffff
      swap3
      dup4
      and
      add
      swap1
      swap2
      and
      or
      dup3
      sstore
      0x0
        /* "TokenNetwork200.sol":11267:11276  nonce > 0 */
      dup8
      gt
        /* "TokenNetwork200.sol":11263:11826  if (nonce > 0) {... */
      iszero
      tag_204
      jumpi
      pop
        /* "TokenNetwork200.sol":11326:11355  channel.participants[partner] */
      0xffffffffffffffffffffffffffffffffffffffff
      dup10
      and
      0x0
      swap1
      dup2
      mstore
        /* "TokenNetwork200.sol":11326:11346  channel.participants */
      0x1
      dup3
      add
        /* "TokenNetwork200.sol":11326:11355  channel.participants[partner] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":11567:11591  channel.open_blocknumber */
      dup2
      sload
        /* "TokenNetwork200.sol":11397:11665  recoverAddressFromBalanceProof(... */
      tag_205
      swap1
        /* "TokenNetwork200.sol":11445:11463  channel_identifier */
      dup6
      swap1
        /* "TokenNetwork200.sol":11481:11499  transferred_amount */
      dup12
      swap1
        /* "TokenNetwork200.sol":11517:11526  locksroot */
      dup12
      swap1
        /* "TokenNetwork200.sol":11544:11549  nonce */
      dup12
      swap1
        /* "TokenNetwork200.sol":11567:11591  channel.open_blocknumber */
      0x10000000000000000
      swap1
      div
      0xffffffffffffffff
      and
        /* "TokenNetwork200.sol":11609:11624  additional_hash */
      dup12
        /* "TokenNetwork200.sol":11642:11651  signature */
      dup12
        /* "TokenNetwork200.sol":11397:11427  recoverAddressFromBalanceProof */
      tag_143
        /* "TokenNetwork200.sol":11397:11665  recoverAddressFromBalanceProof(... */
      jump	// in
    tag_205:
        /* "TokenNetwork200.sol":11369:11665  recovered_partner_address = recoverAddressFromBalanceProof(... */
      swap3
      pop
        /* "TokenNetwork200.sol":11687:11721  partner==recovered_partner_address */
      0xffffffffffffffffffffffffffffffffffffffff
      dup11
      dup2
      and
      swap1
      dup5
      and
      eq
        /* "TokenNetwork200.sol":11679:11722  require(partner==recovered_partner_address) */
      tag_206
      jumpi
      0x0
      dup1
      revert
    tag_206:
        /* "TokenNetwork200.sol":11763:11815  calceBalanceHash(nonce,locksroot,transferred_amount) */
      tag_207
        /* "TokenNetwork200.sol":11780:11785  nonce */
      dup8
        /* "TokenNetwork200.sol":11786:11795  locksroot */
      dup10
        /* "TokenNetwork200.sol":11796:11814  transferred_amount */
      dup12
        /* "TokenNetwork200.sol":11763:11779  calceBalanceHash */
      tag_117
        /* "TokenNetwork200.sol":11763:11815  calceBalanceHash(nonce,locksroot,transferred_amount) */
      jump	// in
    tag_207:
        /* "TokenNetwork200.sol":11736:11762  partner_state.balance_hash */
      0x1
      dup3
      add
        /* "TokenNetwork200.sol":11736:11815  partner_state.balance_hash=calceBalanceHash(nonce,locksroot,transferred_amount) */
      sstore
        /* "TokenNetwork200.sol":11263:11826  if (nonce > 0) {... */
    tag_204:
        /* "TokenNetwork200.sol":11840:11920  ChannelClosed(channel_identifier, msg.sender,nonce,locksroot,transferred_amount) */
      0x40
      dup1
      mload
        /* "TokenNetwork200.sol":11874:11884  msg.sender */
      caller
        /* "TokenNetwork200.sol":11840:11920  ChannelClosed(channel_identifier, msg.sender,nonce,locksroot,transferred_amount) */
      dup2
      mstore
      0x20
      dup2
      add
      dup10
      swap1
      mstore
      dup1
      dup3
      add
      dup11
      swap1
      mstore
      0x60
      dup2
      add
      dup12
      swap1
      mstore
      swap1
      mload
        /* "TokenNetwork200.sol":11854:11872  channel_identifier */
      dup6
      swap2
        /* "TokenNetwork200.sol":11840:11920  ChannelClosed(channel_identifier, msg.sender,nonce,locksroot,transferred_amount) */
      0x939a8c03f193dee253805a941b9d1eb72504a6e61c270cc17a24fe3403d86a21
      swap2
      swap1
      dup2
      swap1
      sub
      0x80
      add
      swap1
      log2
        /* "TokenNetwork200.sol":10282:11927  function closeChannel(... */
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":14088:16317  function unlock(... */
    tag_75:
        /* "TokenNetwork200.sol":14310:14336  bytes32 channel_identifier */
      0x0
        /* "TokenNetwork200.sol":14346:14368  bytes32 locksroot_hash */
      dup1
        /* "TokenNetwork200.sol":14378:14404  bytes32 computed_locksroot */
      0x0
        /* "TokenNetwork200.sol":14414:14437  uint256 unlocked_amount */
      dup1
        /* "TokenNetwork200.sol":14447:14467  bytes32 balance_hash */
      0x0
        /* "TokenNetwork200.sol":14595:14618  Channel storage channel */
      dup1
        /* "TokenNetwork200.sol":14659:14696  Participant storage participant_state */
      0x0
        /* "TokenNetwork200.sol":14513:14514  0 */
      dup1
        /* "TokenNetwork200.sol":14485:14503  merkle_tree_leaves */
      dup9
        /* "TokenNetwork200.sol":14485:14510  merkle_tree_leaves.length */
      mload
        /* "TokenNetwork200.sol":14485:14514  merkle_tree_leaves.length > 0 */
      gt
        /* "TokenNetwork200.sol":14477:14515  require(merkle_tree_leaves.length > 0) */
      iszero
      iszero
      tag_209
      jumpi
      0x0
      dup1
      revert
    tag_209:
        /* "TokenNetwork200.sol":14544:14585  getChannelIdentifier(participant,partner) */
      tag_210
        /* "TokenNetwork200.sol":14565:14576  participant */
      dup14
        /* "TokenNetwork200.sol":14577:14584  partner */
      dup14
        /* "TokenNetwork200.sol":14544:14564  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":14544:14585  getChannelIdentifier(participant,partner) */
      jump	// in
    tag_210:
        /* "TokenNetwork200.sol":14525:14585  channel_identifier=getChannelIdentifier(participant,partner) */
      swap7
      pop
        /* "TokenNetwork200.sol":14621:14629  channels */
      0x3
        /* "TokenNetwork200.sol":14621:14649  channels[channel_identifier] */
      0x0
        /* "TokenNetwork200.sol":14630:14648  channel_identifier */
      dup9
        /* "TokenNetwork200.sol":14621:14649  channels[channel_identifier] */
      not(0x0)
      and
      not(0x0)
      and
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
        /* "TokenNetwork200.sol":14595:14649  Channel storage channel = channels[channel_identifier] */
      swap2
      pop
        /* "TokenNetwork200.sol":14699:14706  channel */
      dup2
        /* "TokenNetwork200.sol":14699:14719  channel.participants */
      0x1
      add
        /* "TokenNetwork200.sol":14699:14732  channel.participants[participant] */
      0x0
        /* "TokenNetwork200.sol":14720:14731  participant */
      dup15
        /* "TokenNetwork200.sol":14699:14732  channel.participants[participant] */
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
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
        /* "TokenNetwork200.sol":14659:14732  Participant storage participant_state = channel.participants[participant] */
      swap1
      pop
        /* "TokenNetwork200.sol":14781:14793  block.number */
      number
        /* "TokenNetwork200.sol":14750:14757  channel */
      dup3
        /* "TokenNetwork200.sol":14750:14777  channel.settle_block_number */
      0x0
      add
      0x0
      swap1
      sload
      swap1
      0x100
      exp
      swap1
      div
      0xffffffffffffffff
      and
        /* "TokenNetwork200.sol":14750:14793  channel.settle_block_number >= block.number */
      0xffffffffffffffff
      and
      lt
      iszero
        /* "TokenNetwork200.sol":14742:14794  require(channel.settle_block_number >= block.number) */
      iszero
      iszero
      tag_211
      jumpi
      0x0
      dup1
      revert
    tag_211:
        /* "TokenNetwork200.sol":14812:14825  channel.state */
      dup2
      sload
      0x100000000000000000000000000000000
      swap1
      div
      0xff
      and
        /* "TokenNetwork200.sol":14827:14828  2 */
      0x2
        /* "TokenNetwork200.sol":14812:14828  channel.state==2 */
      eq
        /* "TokenNetwork200.sol":14804:14829  require(channel.state==2) */
      tag_212
      jumpi
      0x0
      dup1
      revert
    tag_212:
        /* "TokenNetwork200.sol":14852:14882  participant_state.balance_hash */
      0x1
      dup2
      add
      sload
      swap3
      pop
        /* "TokenNetwork200.sol":14946:14960  locksroot != 0 */
      dup10
      iszero
      iszero
        /* "TokenNetwork200.sol":14938:14961  require(locksroot != 0) */
      tag_213
      jumpi
      0x0
      dup1
      revert
    tag_213:
        /* "TokenNetwork200.sol":15103:15155  keccak256(balance_hash,locksroot,channel_identifier) */
      0x40
      dup1
      mload
      dup5
      dup2
      mstore
      0x20
      dup1
      dup3
      add
      dup14
      swap1
      mstore
      dup2
      dup4
      add
      dup11
      swap1
      mstore
      dup3
      mload
      swap2
      dup3
      swap1
      sub
      0x60
      add
      swap1
      swap2
      keccak256
      0x0
        /* "TokenNetwork200.sol":15258:15292  unlocked_locksroot[locksroot_hash] */
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":15258:15276  unlocked_locksroot */
      0x4
        /* "TokenNetwork200.sol":15258:15292  unlocked_locksroot[locksroot_hash] */
      swap1
      swap3
      mstore
      swap2
      swap1
      keccak256
      sload
        /* "TokenNetwork200.sol":15103:15155  keccak256(balance_hash,locksroot,channel_identifier) */
      swap1
      swap7
      pop
        /* "TokenNetwork200.sol":15258:15292  unlocked_locksroot[locksroot_hash] */
      0xff
      and
        /* "TokenNetwork200.sol":15258:15299  unlocked_locksroot[locksroot_hash]==false */
      iszero
        /* "TokenNetwork200.sol":15250:15300  require(unlocked_locksroot[locksroot_hash]==false) */
      tag_214
      jumpi
      0x0
      dup1
      revert
    tag_214:
        /* "TokenNetwork200.sol":15310:15344  unlocked_locksroot[locksroot_hash] */
      0x0
      dup7
      dup2
      mstore
        /* "TokenNetwork200.sol":15310:15328  unlocked_locksroot */
      0x4
        /* "TokenNetwork200.sol":15310:15344  unlocked_locksroot[locksroot_hash] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":15310:15349  unlocked_locksroot[locksroot_hash]=true */
      dup1
      sload
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00
      and
        /* "TokenNetwork200.sol":15345:15349  true */
      0x1
        /* "TokenNetwork200.sol":15310:15349  unlocked_locksroot[locksroot_hash]=true */
      or
      swap1
      sstore
        /* "TokenNetwork200.sol":15566:15616  getMerkleRootAndUnlockedAmount(merkle_tree_leaves) */
      tag_215
        /* "TokenNetwork200.sol":15597:15615  merkle_tree_leaves */
      dup9
        /* "TokenNetwork200.sol":15566:15596  getMerkleRootAndUnlockedAmount */
      tag_216
        /* "TokenNetwork200.sol":15566:15616  getMerkleRootAndUnlockedAmount(merkle_tree_leaves) */
      jump	// in
    tag_215:
        /* "TokenNetwork200.sol":15526:15616  (computed_locksroot, unlocked_amount) = getMerkleRootAndUnlockedAmount(merkle_tree_leaves) */
      swap1
      swap6
      pop
      swap4
      pop
        /* "TokenNetwork200.sol":15650:15651  0 */
      0x0
        /* "TokenNetwork200.sol":15634:15651  unlocked_amount>0 */
      dup5
      gt
        /* "TokenNetwork200.sol":15626:15652  require(unlocked_amount>0) */
      tag_217
      jumpi
      0x0
      dup1
      revert
    tag_217:
        /* "TokenNetwork200.sol":15686:15717  computed_locksroot == locksroot */
      dup5
      dup11
      eq
        /* "TokenNetwork200.sol":15678:15718  require(computed_locksroot == locksroot) */
      tag_218
      jumpi
      0x0
      dup1
      revert
    tag_218:
        /* "TokenNetwork200.sol":15750:15803  calceBalanceHash(nonce,locksroot,transferered_amount) */
      tag_219
        /* "TokenNetwork200.sol":15767:15772  nonce */
      dup10
        /* "TokenNetwork200.sol":15773:15782  locksroot */
      dup12
        /* "TokenNetwork200.sol":15783:15802  transferered_amount */
      dup14
        /* "TokenNetwork200.sol":15750:15766  calceBalanceHash */
      tag_117
        /* "TokenNetwork200.sol":15750:15803  calceBalanceHash(nonce,locksroot,transferered_amount) */
      jump	// in
    tag_219:
        /* "TokenNetwork200.sol":15736:15803  balance_hash==calceBalanceHash(nonce,locksroot,transferered_amount) */
      dup4
      eq
        /* "TokenNetwork200.sol":15728:15804  require(balance_hash==calceBalanceHash(nonce,locksroot,transferered_amount)) */
      tag_220
      jumpi
      0x0
      dup1
      revert
    tag_220:
        /* "TokenNetwork200.sol":15998:16036  transferered_amount += unlocked_amount */
      swap10
      dup4
      add
      swap10
        /* "TokenNetwork200.sol":16148:16201  calceBalanceHash(nonce,locksroot,transferered_amount) */
      tag_221
        /* "TokenNetwork200.sol":16165:16170  nonce */
      dup10
        /* "TokenNetwork200.sol":16171:16180  locksroot */
      dup12
        /* "TokenNetwork200.sol":15998:16036  transferered_amount += unlocked_amount */
      dup14
        /* "TokenNetwork200.sol":16148:16164  calceBalanceHash */
      tag_117
        /* "TokenNetwork200.sol":16148:16201  calceBalanceHash(nonce,locksroot,transferered_amount) */
      jump	// in
    tag_221:
        /* "TokenNetwork200.sol":16117:16147  participant_state.balance_hash */
      0x1
      dup3
      add
        /* "TokenNetwork200.sol":16117:16201  participant_state.balance_hash=calceBalanceHash(nonce,locksroot,transferered_amount) */
      sstore
        /* "TokenNetwork200.sol":16216:16310  ChannelUnlocked(channel_identifier, participant,nonce,computed_locksroot, transferered_amount) */
      0x40
      dup1
      mload
      0xffffffffffffffffffffffffffffffffffffffff
      dup16
      and
      dup2
      mstore
      0x20
      dup2
      add
      dup12
      swap1
      mstore
      dup1
      dup3
      add
      dup8
      swap1
      mstore
      0x60
      dup2
      add
      dup14
      swap1
      mstore
      swap1
      mload
        /* "TokenNetwork200.sol":16232:16250  channel_identifier */
      dup9
      swap2
        /* "TokenNetwork200.sol":16216:16310  ChannelUnlocked(channel_identifier, participant,nonce,computed_locksroot, transferered_amount) */
      0x27e50253fbeac1c4afe1bc9bcb631bfdf66273ef9514547087b59ea681d9b584
      swap2
      swap1
      dup2
      swap1
      sub
      0x80
      add
      swap1
      log2
        /* "TokenNetwork200.sol":14088:16317  function unlock(... */
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":28042:28513  function getChannelInfo(address participant1,address participant2)... */
    tag_78:
        /* "TokenNetwork200.sol":28144:28151  bytes32 */
      0x0
        /* "TokenNetwork200.sol":28152:28158  uint64 */
      dup1
        /* "TokenNetwork200.sol":28159:28165  uint64 */
      0x0
        /* "TokenNetwork200.sol":28168:28173  uint8 */
      dup1
        /* "TokenNetwork200.sol":28190:28216  bytes32 channel_identifier */
      0x0
        /* "TokenNetwork200.sol":28302:28325  Channel storage channel */
      dup1
        /* "TokenNetwork200.sol":28245:28292  getChannelIdentifier(participant1,participant2) */
      tag_223
        /* "TokenNetwork200.sol":28266:28278  participant1 */
      dup9
        /* "TokenNetwork200.sol":28279:28291  participant2 */
      dup9
        /* "TokenNetwork200.sol":28245:28265  getChannelIdentifier */
      tag_84
        /* "TokenNetwork200.sol":28245:28292  getChannelIdentifier(participant1,participant2) */
      jump	// in
    tag_223:
        /* "TokenNetwork200.sol":28328:28356  channels[channel_identifier] */
      0x0
      dup2
      dup2
      mstore
        /* "TokenNetwork200.sol":28328:28336  channels */
      0x3
        /* "TokenNetwork200.sol":28328:28356  channels[channel_identifier] */
      0x20
      mstore
      0x40
      swap1
      keccak256
        /* "TokenNetwork200.sol":28412:28439  channel.settle_block_number */
      sload
        /* "TokenNetwork200.sol":28328:28356  channels[channel_identifier] */
      swap1
      swap10
        /* "TokenNetwork200.sol":28412:28439  channel.settle_block_number */
      0xffffffffffffffff
      dup1
      dup4
      and
      swap11
      pop
        /* "TokenNetwork200.sol":28449:28473  channel.open_blocknumber */
      0x10000000000000000
      dup4
      div
      and
      swap9
      pop
        /* "TokenNetwork200.sol":28483:28496  channel.state */
      0x100000000000000000000000000000000
      swap1
      swap2
      div
      0xff
      and
      swap7
      pop
        /* "TokenNetwork200.sol":28042:28513  function getChannelInfo(address participant1,address participant2)... */
      swap5
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":689:707  Token public token */
    tag_81:
      and(0xffffffffffffffffffffffffffffffffffffffff, sload(0x0))
      dup2
      jump	// out
        /* "TokenNetwork200.sol":26348:26681  function getChannelIdentifier(address participant1,address participant2) pure internal returns (bytes32){... */
    tag_84:
        /* "TokenNetwork200.sol":26444:26451  bytes32 */
      0x0
        /* "TokenNetwork200.sol":26481:26493  participant2 */
      dup2
        /* "TokenNetwork200.sol":26466:26493  participant1 < participant2 */
      0xffffffffffffffffffffffffffffffffffffffff
      and
        /* "TokenNetwork200.sol":26466:26478  participant1 */
      dup4
        /* "TokenNetwork200.sol":26466:26493  participant1 < participant2 */
      0xffffffffffffffffffffffffffffffffffffffff
      and
      lt
        /* "TokenNetwork200.sol":26462:26675  if (participant1 < participant2) {... */
      iszero
      tag_225
      jumpi
        /* "TokenNetwork200.sol":26543:26555  participant1 */
      dup3
        /* "TokenNetwork200.sol":26557:26569  participant2 */
      dup3
        /* "TokenNetwork200.sol":26526:26570  abi.encodePacked(participant1, participant2) */
      add(0x20, mload(0x40))
      dup1
      dup4
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      dup3
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      swap3
      pop
      pop
      pop
      mload(0x40)
        /* "--CODEGEN--":49:53   */
      0x20
        /* "--CODEGEN--":39:46   */
      dup2
        /* "--CODEGEN--":30:37   */
      dup4
        /* "--CODEGEN--":26:47   */
      sub
        /* "--CODEGEN--":22:54   */
      sub
        /* "--CODEGEN--":13:20   */
      dup2
        /* "--CODEGEN--":6:55   */
      mstore
        /* "TokenNetwork200.sol":26526:26570  abi.encodePacked(participant1, participant2) */
      swap1
      0x40
      mstore
        /* "TokenNetwork200.sol":26516:26571  keccak256(abi.encodePacked(participant1, participant2)) */
      mload(0x40)
      dup1
      dup3
      dup1
      mload
      swap1
      0x20
      add
      swap1
      dup1
      dup4
      dup4
        /* "--CODEGEN--":36:189   */
    tag_226:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_227
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_226)
    tag_227:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":26516:26571  keccak256(abi.encodePacked(participant1, participant2)) */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":26509:26571  return keccak256(abi.encodePacked(participant1, participant2)) */
      swap1
      pop
      jump(tag_229)
        /* "TokenNetwork200.sol":26462:26675  if (participant1 < participant2) {... */
    tag_225:
        /* "TokenNetwork200.sol":26636:26648  participant2 */
      dup2
        /* "TokenNetwork200.sol":26650:26662  participant1 */
      dup4
        /* "TokenNetwork200.sol":26619:26663  abi.encodePacked(participant2, participant1) */
      add(0x20, mload(0x40))
      dup1
      dup4
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      dup3
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      swap3
      pop
      pop
      pop
      mload(0x40)
        /* "--CODEGEN--":49:53   */
      0x20
        /* "--CODEGEN--":39:46   */
      dup2
        /* "--CODEGEN--":30:37   */
      dup4
        /* "--CODEGEN--":26:47   */
      sub
        /* "--CODEGEN--":22:54   */
      sub
        /* "--CODEGEN--":13:20   */
      dup2
        /* "--CODEGEN--":6:55   */
      mstore
        /* "TokenNetwork200.sol":26619:26663  abi.encodePacked(participant2, participant1) */
      swap1
      0x40
      mstore
        /* "TokenNetwork200.sol":26609:26664  keccak256(abi.encodePacked(participant2, participant1)) */
      mload(0x40)
      dup1
      dup3
      dup1
      mload
      swap1
      0x20
      add
      swap1
      dup1
      dup4
      dup4
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_227
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_226)
        /* "TokenNetwork200.sol":26462:26675  if (participant1 < participant2) {... */
    tag_229:
        /* "TokenNetwork200.sol":26348:26681  function getChannelIdentifier(address participant1,address participant2) pure internal returns (bytes32){... */
      swap3
      swap2
      pop
      pop
      jump	// out
        /* "ECVerify.sol":50:1040  function ecverify(bytes32 hash, bytes signature)... */
    tag_90:
        /* "ECVerify.sol":146:171  address signature_address */
      0x0
        /* "ECVerify.sol":229:238  bytes32 r */
      dup1
        /* "ECVerify.sol":248:257  bytes32 s */
      0x0
        /* "ECVerify.sol":267:274  uint8 v */
      dup1
        /* "ECVerify.sol":195:204  signature */
      dup5
        /* "ECVerify.sol":195:211  signature.length */
      mload
        /* "ECVerify.sol":215:217  65 */
      0x41
        /* "ECVerify.sol":195:217  signature.length == 65 */
      eq
        /* "ECVerify.sol":187:218  require(signature.length == 65) */
      iszero
      iszero
      tag_234
      jumpi
      0x0
      dup1
      revert
    tag_234:
      pop
      pop
      pop
        /* "ECVerify.sol":492:494  32 */
      0x20
        /* "ECVerify.sol":477:495  add(signature, 32) */
      dup3
      add
        /* "ECVerify.sol":471:496  mload(add(signature, 32)) */
      mload
        /* "ECVerify.sol":535:537  64 */
      0x40
        /* "ECVerify.sol":520:538  add(signature, 64) */
      dup4
      add
        /* "ECVerify.sol":514:539  mload(add(signature, 64)) */
      mload
        /* "ECVerify.sol":668:670  96 */
      0x60
        /* "ECVerify.sol":653:671  add(signature, 96) */
      dup5
      add
        /* "ECVerify.sol":647:672  mload(add(signature, 96)) */
      mload
        /* "ECVerify.sol":644:645  0 */
      0x0
        /* "ECVerify.sol":639:673  byte(0, mload(add(signature, 96))) */
      byte
        /* "ECVerify.sol":783:785  27 */
      0x1b
        /* "ECVerify.sol":779:785  v < 27 */
      0xff
      dup3
      and
      lt
        /* "ECVerify.sol":775:819  if (v < 27) {... */
      iszero
      tag_235
      jumpi
        /* "ECVerify.sol":806:808  27 */
      0x1b
        /* "ECVerify.sol":801:808  v += 27 */
      add
        /* "ECVerify.sol":775:819  if (v < 27) {... */
    tag_235:
        /* "ECVerify.sol":837:838  v */
      dup1
        /* "ECVerify.sol":837:844  v == 27 */
      0xff
      and
        /* "ECVerify.sol":842:844  27 */
      0x1b
        /* "ECVerify.sol":837:844  v == 27 */
      eq
        /* "ECVerify.sol":837:855  v == 27 || v == 28 */
      dup1
      tag_236
      jumpi
      pop
        /* "ECVerify.sol":848:849  v */
      dup1
        /* "ECVerify.sol":848:855  v == 28 */
      0xff
      and
        /* "ECVerify.sol":853:855  28 */
      0x1c
        /* "ECVerify.sol":848:855  v == 28 */
      eq
        /* "ECVerify.sol":837:855  v == 27 || v == 28 */
    tag_236:
        /* "ECVerify.sol":829:856  require(v == 27 || v == 28) */
      iszero
      iszero
      tag_237
      jumpi
      0x0
      dup1
      revert
    tag_237:
        /* "ECVerify.sol":887:911  ecrecover(hash, v, r, s) */
      0x40
      dup1
      mload
      0x0
      dup1
      dup3
      mstore
      0x20
      dup1
      dup4
      add
      dup1
      dup6
      mstore
      dup11
      swap1
      mstore
      0xff
      dup6
      and
      dup4
      dup6
      add
      mstore
      0x60
      dup4
      add
      dup8
      swap1
      mstore
      0x80
      dup4
      add
      dup7
      swap1
      mstore
      swap3
      mload
      0x1
      swap4
      0xa0
      dup1
      dup6
      add
      swap5
      swap2
      swap4
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      dup5
      add
      swap4
      swap3
      dup4
      swap1
      sub
      swap1
      swap2
      add
      swap2
      swap1
      dup7
      gas
      call
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_238
      jumpi
        /* "--CODEGEN--":45:61   */
      returndatasize
        /* "--CODEGEN--":42:43   */
      0x0
        /* "--CODEGEN--":39:40   */
      dup1
        /* "--CODEGEN--":24:62   */
      returndatacopy
        /* "--CODEGEN--":77:93   */
      returndatasize
        /* "--CODEGEN--":74:75   */
      0x0
        /* "--CODEGEN--":67:94   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_238:
      pop
      pop
        /* "ECVerify.sol":887:911  ecrecover(hash, v, r, s) */
      mload(add(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0, mload(0x40)))
      swap5
      pop
      pop
        /* "ECVerify.sol":973:997  signature_address != 0x0 */
      0xffffffffffffffffffffffffffffffffffffffff
      dup5
      and
      iszero
      iszero
        /* "ECVerify.sol":965:998  require(signature_address != 0x0) */
      tag_239
      jumpi
      0x0
      dup1
      revert
    tag_239:
        /* "ECVerify.sol":50:1040  function ecverify(bytes32 hash, bytes signature)... */
      pop
      pop
      pop
      swap3
      swap2
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":27755:28036  function calceBalanceHash(uint256 nonce,bytes32 locksroot,uint256 transferred_amount) pure internal returns(bytes32){... */
    tag_117:
        /* "TokenNetwork200.sol":27863:27870  bytes32 */
      0x0
        /* "TokenNetwork200.sol":27885:27893  nonce==0 */
      dup4
      iszero
        /* "TokenNetwork200.sol":27885:27909  nonce==0 && locksroot==0 */
      dup1
      iszero
      tag_241
      jumpi
      pop
        /* "TokenNetwork200.sol":27897:27909  locksroot==0 */
      dup3
      iszero
        /* "TokenNetwork200.sol":27885:27909  nonce==0 && locksroot==0 */
    tag_241:
        /* "TokenNetwork200.sol":27885:27934  nonce==0 && locksroot==0 && transferred_amount==0 */
      dup1
      iszero
      tag_242
      jumpi
      pop
        /* "TokenNetwork200.sol":27913:27934  transferred_amount==0 */
      dup2
      iszero
        /* "TokenNetwork200.sol":27885:27934  nonce==0 && locksroot==0 && transferred_amount==0 */
    tag_242:
        /* "TokenNetwork200.sol":27881:27968  if( nonce==0 && locksroot==0 && transferred_amount==0){... */
      iszero
      tag_243
      jumpi
      pop
        /* "TokenNetwork200.sol":27956:27957  0 */
      0x0
        /* "TokenNetwork200.sol":27949:27957  return 0 */
      jump(tag_240)
        /* "TokenNetwork200.sol":27881:27968  if( nonce==0 && locksroot==0 && transferred_amount==0){... */
    tag_243:
      pop
        /* "TokenNetwork200.sol":27984:28029  keccak256(nonce,locksroot,transferred_amount) */
      0x40
      dup1
      mload
      dup5
      dup2
      mstore
      0x20
      dup2
      add
      dup5
      swap1
      mstore
      dup1
      dup3
      add
      dup4
      swap1
      mstore
      swap1
      mload
      swap1
      dup2
      swap1
      sub
      0x60
      add
      swap1
      keccak256
        /* "TokenNetwork200.sol":27755:28036  function calceBalanceHash(uint256 nonce,bytes32 locksroot,uint256 transferred_amount) pure internal returns(bytes32){... */
    tag_240:
      swap4
      swap3
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":35569:35677  function min(uint256 a, uint256 b) pure internal returns (uint256)... */
    tag_123:
        /* "TokenNetwork200.sol":35627:35634  uint256 */
      0x0
        /* "TokenNetwork200.sol":35661:35662  b */
      dup2
        /* "TokenNetwork200.sol":35657:35658  a */
      dup4
        /* "TokenNetwork200.sol":35657:35662  a > b */
      gt
        /* "TokenNetwork200.sol":35657:35670  a > b ? b : a */
      tag_245
      jumpi
        /* "TokenNetwork200.sol":35669:35670  a */
      dup3
        /* "TokenNetwork200.sol":35657:35670  a > b ? b : a */
      jump(tag_240)
    tag_245:
      pop
        /* "TokenNetwork200.sol":35665:35666  b */
      swap2
        /* "TokenNetwork200.sol":35650:35670  return a > b ? b : a */
      swap1
      pop
        /* "TokenNetwork200.sol":35569:35677  function min(uint256 a, uint256 b) pure internal returns (uint256)... */
      jump	// out
        /* "TokenNetwork200.sol":29766:30596  function recoverAddressFromBalanceProofUpdateMessage(... */
    tag_140:
        /* "TokenNetwork200.sol":30118:30143  address signature_address */
      0x0
        /* "TokenNetwork200.sol":30159:30179  bytes32 message_hash */
      dup1
        /* "TokenNetwork200.sol":30226:30244  transferred_amount */
      dup9
        /* "TokenNetwork200.sol":30262:30271  locksroot */
      dup9
        /* "TokenNetwork200.sol":30289:30294  nonce */
      dup9
        /* "TokenNetwork200.sol":30312:30327  additional_hash */
      dup8
        /* "TokenNetwork200.sol":30345:30363  channel_identifier */
      dup14
        /* "TokenNetwork200.sol":30381:30397  open_blocknumber */
      dup11
        /* "TokenNetwork200.sol":30423:30427  this */
      address
        /* "TokenNetwork200.sol":30446:30454  chain_id */
      sload(0x2)
        /* "TokenNetwork200.sol":30472:30489  closing_signature */
      dup12
        /* "TokenNetwork200.sol":30192:30503  abi.encodePacked(... */
      add(0x20, mload(0x40))
      dup1
      dup11
      dup2
      mstore
      0x20
      add
      dup10
      not(0x0)
      and
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup9
      dup2
      mstore
      0x20
      add
      dup8
      not(0x0)
      and
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup7
      not(0x0)
      and
      not(0x0)
      and
      dup2
      mstore
      0x20
      add
      dup6
      0xffffffffffffffff
      and
      0xffffffffffffffff
      and
      0x1000000000000000000000000000000000000000000000000
      mul
      dup2
      mstore
      0x8
      add
      dup5
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0xffffffffffffffffffffffffffffffffffffffff
      and
      0x1000000000000000000000000
      mul
      dup2
      mstore
      0x14
      add
      dup4
      dup2
      mstore
      0x20
      add
      dup3
      dup1
      mload
      swap1
      0x20
      add
      swap1
      dup1
      dup4
      dup4
        /* "--CODEGEN--":36:189   */
    tag_248:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_249
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_248)
    tag_249:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":30192:30503  abi.encodePacked(... */
      pop
      pop
      pop
      swap1
      pop
      add
      swap10
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      mload(0x40)
        /* "--CODEGEN--":49:53   */
      0x20
        /* "--CODEGEN--":39:46   */
      dup2
        /* "--CODEGEN--":30:37   */
      dup4
        /* "--CODEGEN--":26:47   */
      sub
        /* "--CODEGEN--":22:54   */
      sub
        /* "--CODEGEN--":13:20   */
      dup2
        /* "--CODEGEN--":6:55   */
      mstore
        /* "TokenNetwork200.sol":30192:30503  abi.encodePacked(... */
      swap1
      0x40
      mstore
        /* "TokenNetwork200.sol":30182:30504  keccak256(abi.encodePacked(... */
      mload(0x40)
      dup1
      dup3
      dup1
      mload
      swap1
      0x20
      add
      swap1
      dup1
      dup4
      dup4
        /* "--CODEGEN--":36:189   */
    tag_251:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_252
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_251)
    tag_252:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":30182:30504  keccak256(abi.encodePacked(... */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":30159:30504  bytes32 message_hash = keccak256(abi.encodePacked(... */
      swap1
      pop
        /* "TokenNetwork200.sol":30535:30589  ECVerify.ecverify(message_hash, non_closing_signature) */
      tag_254
        /* "TokenNetwork200.sol":30553:30565  message_hash */
      dup2
        /* "TokenNetwork200.sol":30567:30588  non_closing_signature */
      dup5
        /* "TokenNetwork200.sol":30535:30552  ECVerify.ecverify */
      tag_90
        /* "TokenNetwork200.sol":30535:30589  ECVerify.ecverify(message_hash, non_closing_signature) */
      jump	// in
    tag_254:
        /* "TokenNetwork200.sol":30515:30589  signature_address = ECVerify.ecverify(message_hash, non_closing_signature) */
      swap11
        /* "TokenNetwork200.sol":29766:30596  function recoverAddressFromBalanceProofUpdateMessage(... */
      swap10
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":29034:29759  function recoverAddressFromBalanceProof(... */
    tag_143:
        /* "TokenNetwork200.sol":29656:29664  chain_id */
      sload(0x2)
        /* "TokenNetwork200.sol":29402:29678  abi.encodePacked(... */
      0x40
      dup1
      mload
      0x20
      dup1
      dup3
      add
      dup11
      swap1
      mstore
      dup2
      dup4
      add
      dup10
      swap1
      mstore
      0x60
      dup3
      add
      dup9
      swap1
      mstore
      0x80
      dup3
      add
      dup7
      swap1
      mstore
      0xa0
      dup3
      add
      dup12
      swap1
      mstore
      0x1000000000000000000000000000000000000000000000000
      0xffffffffffffffff
      dup9
      and
      mul
      0xc0
      dup4
      add
      mstore
      0x1000000000000000000000000
        /* "TokenNetwork200.sol":29633:29637  this */
      address
        /* "TokenNetwork200.sol":29402:29678  abi.encodePacked(... */
      mul
      0xc8
      dup4
      add
      mstore
      0xdc
      dup1
      dup4
      add
      swap5
      swap1
      swap5
      mstore
      dup3
      mload
        /* "--CODEGEN--":26:47   */
      dup1
      dup4
      sub
        /* "--CODEGEN--":22:54   */
      swap1
      swap5
      add
        /* "--CODEGEN--":6:55   */
      dup5
      mstore
        /* "TokenNetwork200.sol":29402:29678  abi.encodePacked(... */
      0xfc
      swap1
      swap2
      add
      swap2
      dup3
      swap1
      mstore
        /* "TokenNetwork200.sol":29392:29679  keccak256(abi.encodePacked(... */
      dup3
      mload
        /* "TokenNetwork200.sol":29328:29353  address signature_address */
      0x0
      swap4
      dup5
      swap4
        /* "TokenNetwork200.sol":29402:29678  abi.encodePacked(... */
      swap1
      swap3
      swap1
      swap2
      dup3
      swap2
        /* "TokenNetwork200.sol":29392:29679  keccak256(abi.encodePacked(... */
      dup5
      add
      swap1
      dup1
        /* "TokenNetwork200.sol":29402:29678  abi.encodePacked(... */
      dup4
        /* "TokenNetwork200.sol":29392:29679  keccak256(abi.encodePacked(... */
      dup4
        /* "--CODEGEN--":36:189   */
    tag_256:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_257
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_256)
    tag_257:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":29392:29679  keccak256(abi.encodePacked(... */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":29369:29679  bytes32 message_hash = keccak256(abi.encodePacked(... */
      swap1
      pop
        /* "TokenNetwork200.sol":29710:29752  ECVerify.ecverify(message_hash, signature) */
      tag_259
        /* "TokenNetwork200.sol":29728:29740  message_hash */
      dup2
        /* "TokenNetwork200.sol":29742:29751  signature */
      dup5
        /* "TokenNetwork200.sol":29710:29727  ECVerify.ecverify */
      tag_90
        /* "TokenNetwork200.sol":29710:29752  ECVerify.ecverify(message_hash, signature) */
      jump	// in
    tag_259:
        /* "TokenNetwork200.sol":29690:29752  signature_address = ECVerify.ecverify(message_hash, signature) */
      swap10
        /* "TokenNetwork200.sol":29034:29759  function recoverAddressFromBalanceProof(... */
      swap9
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":26688:27263  function verifyBalanceHashIsValid(... */
    tag_146:
        /* "TokenNetwork200.sol":26939:26962  Channel storage channel */
      0x0
        /* "TokenNetwork200.sol":26965:26993  channels[channel_identifier] */
      dup7
      dup2
      mstore
        /* "TokenNetwork200.sol":26965:26973  channels */
      0x3
        /* "TokenNetwork200.sol":26965:26993  channels[channel_identifier] */
      0x20
      swap1
      dup2
      mstore
      0x40
      dup1
      dup4
      keccak256
        /* "TokenNetwork200.sol":27043:27076  channel.participants[participant] */
      0xffffffffffffffffffffffffffffffffffffffff
      dup10
      and
      dup5
      mstore
        /* "TokenNetwork200.sol":27043:27063  channel.participants */
      0x1
      dup2
      add
        /* "TokenNetwork200.sol":27043:27076  channel.participants[participant] */
      swap1
      swap3
      mstore
      dup3
      keccak256
        /* "TokenNetwork200.sol":26965:26993  channels[channel_identifier] */
      swap1
      swap2
        /* "TokenNetwork200.sol":27107:27159  calceBalanceHash(nonce,locksroot,transferred_amount) */
      tag_261
        /* "TokenNetwork200.sol":27124:27129  nonce */
      dup6
        /* "TokenNetwork200.sol":27130:27139  locksroot */
      dup8
        /* "TokenNetwork200.sol":27140:27158  transferred_amount */
      dup10
        /* "TokenNetwork200.sol":27107:27123  calceBalanceHash */
      tag_117
        /* "TokenNetwork200.sol":27107:27159  calceBalanceHash(nonce,locksroot,transferred_amount) */
      jump	// in
    tag_261:
        /* "TokenNetwork200.sol":27177:27207  participant_state.balance_hash */
      0x1
      dup4
      add
      sload
        /* "TokenNetwork200.sol":27086:27159  bytes32 balance_hash=calceBalanceHash(nonce,locksroot,transferred_amount) */
      swap1
      swap2
      pop
        /* "TokenNetwork200.sol":27177:27221  participant_state.balance_hash==balance_hash */
      dup2
      eq
        /* "TokenNetwork200.sol":27169:27222  require(participant_state.balance_hash==balance_hash) */
      tag_262
      jumpi
      0x0
      dup1
      revert
    tag_262:
        /* "TokenNetwork200.sol":27240:27255  new_nonce>nonce */
      dup5
      dup5
      gt
        /* "TokenNetwork200.sol":27232:27256  require(new_nonce>nonce) */
      tag_263
      jumpi
      0x0
      dup1
      revert
    tag_263:
        /* "TokenNetwork200.sol":26688:27263  function verifyBalanceHashIsValid(... */
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":27268:27750  function updateBalanceHash(        bytes32 channel_identifier,... */
    tag_148:
        /* "TokenNetwork200.sol":27470:27493  Channel storage channel */
      0x0
        /* "TokenNetwork200.sol":27496:27524  channels[channel_identifier] */
      dup6
      dup2
      mstore
        /* "TokenNetwork200.sol":27496:27504  channels */
      0x3
        /* "TokenNetwork200.sol":27496:27524  channels[channel_identifier] */
      0x20
      swap1
      dup2
      mstore
      0x40
      dup1
      dup4
      keccak256
        /* "TokenNetwork200.sol":27574:27607  channel.participants[participant] */
      0xffffffffffffffffffffffffffffffffffffffff
      dup9
      and
      dup5
      mstore
        /* "TokenNetwork200.sol":27574:27594  channel.participants */
      0x1
      dup2
      add
        /* "TokenNetwork200.sol":27574:27607  channel.participants[participant] */
      swap1
      swap3
      mstore
      dup3
      keccak256
        /* "TokenNetwork200.sol":27496:27524  channels[channel_identifier] */
      swap1
      swap2
        /* "TokenNetwork200.sol":27638:27690  calceBalanceHash(nonce,locksroot,transferred_amount) */
      tag_265
        /* "TokenNetwork200.sol":27655:27660  nonce */
      dup5
        /* "TokenNetwork200.sol":27661:27670  locksroot */
      dup7
        /* "TokenNetwork200.sol":27671:27689  transferred_amount */
      dup9
        /* "TokenNetwork200.sol":27638:27654  calceBalanceHash */
      tag_117
        /* "TokenNetwork200.sol":27638:27690  calceBalanceHash(nonce,locksroot,transferred_amount) */
      jump	// in
    tag_265:
        /* "TokenNetwork200.sol":27700:27730  participant_state.balance_hash */
      0x1
      swap1
      swap3
      add
        /* "TokenNetwork200.sol":27700:27743  participant_state.balance_hash=balance_hash */
      swap2
      swap1
      swap2
      sstore
      pop
      pop
      pop
      pop
      pop
      pop
      pop
        /* "TokenNetwork200.sol":27268:27750  function updateBalanceHash(        bytes32 channel_identifier,... */
      jump	// out
        /* "TokenNetwork200.sol":30602:31375  function recoverAddressFromCooperativeSettleSignature(... */
    tag_154:
        /* "TokenNetwork200.sol":31272:31280  chain_id */
      sload(0x2)
        /* "TokenNetwork200.sol":31001:31294  abi.encodePacked(... */
      0x40
      dup1
      mload
      0x1000000000000000000000000
      0xffffffffffffffffffffffffffffffffffffffff
      dup1
      dup12
      and
      dup3
      mul
      0x20
      dup1
      dup6
      add
      swap2
      swap1
      swap2
      mstore
      0x34
      dup5
      add
      dup12
      swap1
      mstore
      swap1
      dup10
      and
      dup3
      mul
      0x54
      dup5
      add
      mstore
      0x68
      dup4
      add
      dup9
      swap1
      mstore
      0x88
      dup4
      add
      dup13
      swap1
      mstore
      0x1000000000000000000000000000000000000000000000000
      0xffffffffffffffff
      dup9
      and
      mul
      0xa8
      dup5
      add
      mstore
        /* "TokenNetwork200.sol":31249:31253  this */
      address
        /* "TokenNetwork200.sol":31001:31294  abi.encodePacked(... */
      swap2
      swap1
      swap2
      mul
      0xb0
      dup4
      add
      mstore
      0xc4
      dup1
      dup4
      add
      swap5
      swap1
      swap5
      mstore
      dup3
      mload
        /* "--CODEGEN--":26:47   */
      dup1
      dup4
      sub
        /* "--CODEGEN--":22:54   */
      swap1
      swap5
      add
        /* "--CODEGEN--":6:55   */
      dup5
      mstore
        /* "TokenNetwork200.sol":31001:31294  abi.encodePacked(... */
      0xe4
      swap1
      swap2
      add
      swap2
      dup3
      swap1
      mstore
        /* "TokenNetwork200.sol":30991:31295  keccak256(abi.encodePacked(... */
      dup3
      mload
        /* "TokenNetwork200.sol":30927:30952  address signature_address */
      0x0
      swap4
      dup5
      swap4
        /* "TokenNetwork200.sol":31001:31294  abi.encodePacked(... */
      swap1
      swap3
      swap1
      swap2
      dup3
      swap2
        /* "TokenNetwork200.sol":30991:31295  keccak256(abi.encodePacked(... */
      dup5
      add
      swap1
      dup1
        /* "TokenNetwork200.sol":31001:31294  abi.encodePacked(... */
      dup4
        /* "TokenNetwork200.sol":30991:31295  keccak256(abi.encodePacked(... */
      dup4
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_257
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_256)
        /* "TokenNetwork200.sol":34959:35563  function recoverAddressFromUnlockProof(... */
    tag_176:
        /* "TokenNetwork200.sol":35427:35435  chain_id */
      sload(0x2)
        /* "TokenNetwork200.sol":35266:35482  abi.encodePacked(... */
      0x40
      dup1
      mload
      0x20
      dup1
      dup3
      add
      dup9
      swap1
      mstore
      dup2
      dup4
      add
      dup10
      swap1
      mstore
      0x1000000000000000000000000000000000000000000000000
      0xffffffffffffffff
      dup9
      and
      mul
      0x60
      dup4
      add
      mstore
      0x1000000000000000000000000
        /* "TokenNetwork200.sol":35404:35408  this */
      address
        /* "TokenNetwork200.sol":35266:35482  abi.encodePacked(... */
      mul
      0x68
      dup4
      add
      mstore
      0x7c
      dup3
      add
      swap4
      swap1
      swap4
      mstore
      0x9c
      dup1
      dup3
      add
      dup7
      swap1
      mstore
      dup3
      mload
        /* "--CODEGEN--":26:47   */
      dup1
      dup4
      sub
        /* "--CODEGEN--":22:54   */
      swap1
      swap2
      add
        /* "--CODEGEN--":6:55   */
      dup2
      mstore
        /* "TokenNetwork200.sol":35266:35482  abi.encodePacked(... */
      0xbc
      swap1
      swap2
      add
      swap2
      dup3
      swap1
      mstore
        /* "TokenNetwork200.sol":35256:35483  keccak256(abi.encodePacked(... */
      dup1
      mload
        /* "TokenNetwork200.sol":35192:35217  address signature_address */
      0x0
      swap4
      dup5
      swap4
        /* "TokenNetwork200.sol":35266:35482  abi.encodePacked(... */
      swap2
      dup3
      swap2
        /* "TokenNetwork200.sol":35256:35483  keccak256(abi.encodePacked(... */
      dup5
      add
      swap1
      dup1
        /* "TokenNetwork200.sol":35266:35482  abi.encodePacked(... */
      dup4
        /* "TokenNetwork200.sol":35256:35483  keccak256(abi.encodePacked(... */
      dup4
        /* "--CODEGEN--":36:189   */
    tag_272:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_273
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_272)
    tag_273:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":35256:35483  keccak256(abi.encodePacked(... */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":35233:35483  bytes32 message_hash = keccak256(abi.encodePacked(... */
      swap1
      pop
        /* "TokenNetwork200.sol":35514:35556  ECVerify.ecverify(message_hash, signature) */
      tag_275
        /* "TokenNetwork200.sol":35532:35544  message_hash */
      dup2
        /* "TokenNetwork200.sol":35546:35555  signature */
      dup5
        /* "TokenNetwork200.sol":35514:35531  ECVerify.ecverify */
      tag_90
        /* "TokenNetwork200.sol":35514:35556  ECVerify.ecverify(message_hash, signature) */
      jump	// in
    tag_275:
        /* "TokenNetwork200.sol":35494:35556  signature_address = ECVerify.ecverify(message_hash, signature) */
      swap8
        /* "TokenNetwork200.sol":34959:35563  function recoverAddressFromUnlockProof(... */
      swap7
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":34354:34953  function computeMerkleRoot(bytes32 lockhash, bytes merkle_proof)... */
    tag_179:
        /* "TokenNetwork200.sol":34454:34461  bytes32 */
      0x0
        /* "TokenNetwork200.sol":34526:34535  uint256 i */
      dup1
        /* "TokenNetwork200.sol":34545:34555  bytes32 el */
      0x0
        /* "TokenNetwork200.sol":34507:34509  32 */
      0x20
        /* "TokenNetwork200.sol":34485:34497  merkle_proof */
      dup5
        /* "TokenNetwork200.sol":34485:34504  merkle_proof.length */
      mload
        /* "TokenNetwork200.sol":34485:34509  merkle_proof.length % 32 */
      dup2
      iszero
      iszero
      tag_277
      jumpi
      invalid
    tag_277:
      mod
        /* "TokenNetwork200.sol":34485:34514  merkle_proof.length % 32 == 0 */
      iszero
        /* "TokenNetwork200.sol":34477:34515  require(merkle_proof.length % 32 == 0) */
      tag_278
      jumpi
      0x0
      dup1
      revert
    tag_278:
        /* "TokenNetwork200.sol":34575:34577  32 */
      0x20
        /* "TokenNetwork200.sol":34571:34577  i = 32 */
      swap2
      pop
        /* "TokenNetwork200.sol":34566:34921  for (i = 32; i <= merkle_proof.length; i += 32) {... */
    tag_279:
        /* "TokenNetwork200.sol":34584:34603  merkle_proof.length */
      dup4
      mload
        /* "TokenNetwork200.sol":34579:34603  i <= merkle_proof.length */
      dup3
      gt
        /* "TokenNetwork200.sol":34566:34921  for (i = 32; i <= merkle_proof.length; i += 32) {... */
      tag_280
      jumpi
      pop
        /* "TokenNetwork200.sol":34667:34687  add(merkle_proof, i) */
      dup3
      dup2
      add
        /* "TokenNetwork200.sol":34661:34688  mload(add(merkle_proof, i)) */
      mload
        /* "TokenNetwork200.sol":34720:34733  lockhash < el */
      dup1
      dup6
      lt
        /* "TokenNetwork200.sol":34716:34911  if (lockhash < el) {... */
      iszero
      tag_282
      jumpi
        /* "TokenNetwork200.sol":34774:34804  abi.encodePacked(lockhash, el) */
      0x40
      dup1
      mload
      0x20
      dup1
      dup3
      add
      dup9
      swap1
      mstore
      dup2
      dup4
      add
      dup5
      swap1
      mstore
      dup3
      mload
        /* "--CODEGEN--":26:47   */
      dup1
      dup4
      sub
        /* "--CODEGEN--":22:54   */
      dup5
      add
        /* "--CODEGEN--":6:55   */
      dup2
      mstore
        /* "TokenNetwork200.sol":34774:34804  abi.encodePacked(lockhash, el) */
      0x60
      swap1
      swap3
      add
      swap3
      dup4
      swap1
      mstore
        /* "TokenNetwork200.sol":34764:34805  keccak256(abi.encodePacked(lockhash, el)) */
      dup2
      mload
        /* "TokenNetwork200.sol":34774:34804  abi.encodePacked(lockhash, el) */
      swap2
      swap3
      swap2
      dup3
      swap2
        /* "TokenNetwork200.sol":34764:34805  keccak256(abi.encodePacked(lockhash, el)) */
      dup5
      add
      swap1
      dup1
        /* "TokenNetwork200.sol":34774:34804  abi.encodePacked(lockhash, el) */
      dup4
        /* "TokenNetwork200.sol":34764:34805  keccak256(abi.encodePacked(lockhash, el)) */
      dup4
        /* "--CODEGEN--":36:189   */
    tag_283:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_284
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_283)
    tag_284:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":34764:34805  keccak256(abi.encodePacked(lockhash, el)) */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":34753:34805  lockhash = keccak256(abi.encodePacked(lockhash, el)) */
      swap5
      pop
        /* "TokenNetwork200.sol":34716:34911  if (lockhash < el) {... */
      jump(tag_286)
    tag_282:
        /* "TokenNetwork200.sol":34865:34895  abi.encodePacked(el, lockhash) */
      0x40
      dup1
      mload
      0x20
      dup1
      dup3
      add
      dup5
      swap1
      mstore
      dup2
      dup4
      add
      dup9
      swap1
      mstore
      dup3
      mload
        /* "--CODEGEN--":26:47   */
      dup1
      dup4
      sub
        /* "--CODEGEN--":22:54   */
      dup5
      add
        /* "--CODEGEN--":6:55   */
      dup2
      mstore
        /* "TokenNetwork200.sol":34865:34895  abi.encodePacked(el, lockhash) */
      0x60
      swap1
      swap3
      add
      swap3
      dup4
      swap1
      mstore
        /* "TokenNetwork200.sol":34855:34896  keccak256(abi.encodePacked(el, lockhash)) */
      dup2
      mload
        /* "TokenNetwork200.sol":34865:34895  abi.encodePacked(el, lockhash) */
      swap2
      swap3
      swap2
      dup3
      swap2
        /* "TokenNetwork200.sol":34855:34896  keccak256(abi.encodePacked(el, lockhash)) */
      dup5
      add
      swap1
      dup1
        /* "TokenNetwork200.sol":34865:34895  abi.encodePacked(el, lockhash) */
      dup4
        /* "TokenNetwork200.sol":34855:34896  keccak256(abi.encodePacked(el, lockhash)) */
      dup4
        /* "--CODEGEN--":36:189   */
    tag_287:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_288
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_287)
    tag_288:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":34855:34896  keccak256(abi.encodePacked(el, lockhash)) */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":34844:34896  lockhash = keccak256(abi.encodePacked(el, lockhash)) */
      swap5
      pop
        /* "TokenNetwork200.sol":34716:34911  if (lockhash < el) {... */
    tag_286:
        /* "TokenNetwork200.sol":34610:34612  32 */
      0x20
        /* "TokenNetwork200.sol":34605:34612  i += 32 */
      dup3
      add
      swap2
      pop
        /* "TokenNetwork200.sol":34566:34921  for (i = 32; i <= merkle_proof.length; i += 32) {... */
      jump(tag_279)
    tag_280:
      pop
        /* "TokenNetwork200.sol":34938:34946  lockhash */
      swap3
      swap4
        /* "TokenNetwork200.sol":34354:34953  function computeMerkleRoot(bytes32 lockhash, bytes merkle_proof)... */
      swap3
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":31382:33058  function getMerkleRootAndUnlockedAmount(bytes merkle_tree_leaves)... */
    tag_216:
        /* "TokenNetwork200.sol":31532:31557  merkle_tree_leaves.length */
      dup1
      mload
        /* "TokenNetwork200.sol":31483:31490  bytes32 */
      0x0
      swap1
      dup2
      swap1
      dup2
      dup1
      dup1
      dup1
      dup1
        /* "TokenNetwork200.sol":31885:31914  bytes32[] memory merkle_layer */
      0x60
      dup1
        /* "TokenNetwork200.sol":31532:31557  merkle_tree_leaves.length */
      dup8
        /* "TokenNetwork200.sol":31710:31721  length % 96 */
      mod
        /* "TokenNetwork200.sol":31710:31726  length % 96 == 0 */
      iszero
        /* "TokenNetwork200.sol":31702:31727  require(length % 96 == 0) */
      tag_292
      jumpi
      0x0
      dup1
      revert
    tag_292:
        /* "TokenNetwork200.sol":31940:31942  96 */
      0x60
        /* "TokenNetwork200.sol":31931:31937  length */
      dup8
        /* "TokenNetwork200.sol":31931:31942  length / 96 */
      div
        /* "TokenNetwork200.sol":31945:31946  1 */
      0x1
        /* "TokenNetwork200.sol":31931:31946  length / 96 + 1 */
      add
        /* "TokenNetwork200.sol":31917:31947  new bytes32[](length / 96 + 1) */
      mload(0x40)
      swap1
      dup1
      dup3
      mstore
      dup1
      0x20
      mul
      0x20
      add
      dup3
      add
      0x40
      mstore
      dup1
      iszero
      tag_294
      jumpi
      dup2
      0x20
      add
        /* "--CODEGEN--":29:31   */
      0x20
        /* "--CODEGEN--":21:27   */
      dup3
        /* "--CODEGEN--":17:32   */
      mul
        /* "--CODEGEN--":117:121   */
      dup1
        /* "--CODEGEN--":105:115   */
      codesize
        /* "--CODEGEN--":97:103   */
      dup4
        /* "--CODEGEN--":88:122   */
      codecopy
        /* "--CODEGEN--":136:153   */
      add
      swap1
      pop
        /* "TokenNetwork200.sol":31917:31947  new bytes32[](length / 96 + 1) */
    tag_294:
      pop
        /* "TokenNetwork200.sol":31885:31947  bytes32[] memory merkle_layer = new bytes32[](length / 96 + 1) */
      swap1
      pop
        /* "TokenNetwork200.sol":31967:31969  32 */
      0x20
        /* "TokenNetwork200.sol":31963:31969  i = 32 */
      swap6
      pop
        /* "TokenNetwork200.sol":31958:32194  for (i = 32; i < length; i += 96) {... */
    tag_295:
        /* "TokenNetwork200.sol":31975:31981  length */
      dup7
        /* "TokenNetwork200.sol":31971:31972  i */
      dup7
        /* "TokenNetwork200.sol":31971:31981  i < length */
      lt
        /* "TokenNetwork200.sol":31958:32194  for (i = 32; i < length; i += 96) {... */
      iszero
      tag_296
      jumpi
        /* "TokenNetwork200.sol":32036:32084  getLockDataFromMerkleTree(merkle_tree_leaves, i) */
      tag_298
        /* "TokenNetwork200.sol":32062:32080  merkle_tree_leaves */
      dup11
        /* "TokenNetwork200.sol":32082:32083  i */
      dup8
        /* "TokenNetwork200.sol":32036:32061  getLockDataFromMerkleTree */
      tag_299
        /* "TokenNetwork200.sol":32036:32084  getLockDataFromMerkleTree(merkle_tree_leaves, i) */
      jump	// in
    tag_298:
        /* "TokenNetwork200.sol":32098:32138  total_unlocked_amount += unlocked_amount */
      swap6
      dup7
      add
      swap6
        /* "TokenNetwork200.sol":32006:32084  (lockhash, unlocked_amount) = getLockDataFromMerkleTree(merkle_tree_leaves, i) */
      swap5
      pop
      swap3
      pop
      dup3
        /* "TokenNetwork200.sol":32152:32164  merkle_layer */
      dup2
        /* "TokenNetwork200.sol":32169:32171  96 */
      0x60
        /* "TokenNetwork200.sol":32165:32166  i */
      dup9
        /* "TokenNetwork200.sol":32165:32171  i / 96 */
      div
        /* "TokenNetwork200.sol":32152:32172  merkle_layer[i / 96] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_301
      jumpi
      invalid
    tag_301:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      add
        /* "TokenNetwork200.sol":32152:32183  merkle_layer[i / 96] = lockhash */
      mstore
        /* "TokenNetwork200.sol":31988:31990  96 */
      0x60
        /* "TokenNetwork200.sol":31983:31990  i += 96 */
      swap6
      swap1
      swap6
      add
      swap5
        /* "TokenNetwork200.sol":31958:32194  for (i = 32; i < length; i += 96) {... */
      jump(tag_295)
    tag_296:
        /* "TokenNetwork200.sol":32214:32216  96 */
      0x60
        /* "TokenNetwork200.sol":32204:32216  length /= 96 */
      dup8
      div
      swap7
      pop
        /* "TokenNetwork200.sol":32227:32958  while (length > 1) {... */
    tag_303:
        /* "TokenNetwork200.sol":32243:32244  1 */
      0x1
        /* "TokenNetwork200.sol":32234:32240  length */
      dup8
        /* "TokenNetwork200.sol":32234:32244  length > 1 */
      gt
        /* "TokenNetwork200.sol":32227:32958  while (length > 1) {... */
      iszero
      tag_304
      jumpi
        /* "TokenNetwork200.sol":32273:32274  2 */
      0x2
        /* "TokenNetwork200.sol":32264:32270  length */
      dup8
        /* "TokenNetwork200.sol":32264:32274  length % 2 */
      mod
        /* "TokenNetwork200.sol":32264:32279  length % 2 != 0 */
      iszero
        /* "TokenNetwork200.sol":32260:32390  if (length % 2 != 0) {... */
      tag_306
      jumpi
        /* "TokenNetwork200.sol":32322:32334  merkle_layer */
      dup1
        /* "TokenNetwork200.sol":32344:32345  1 */
      0x1
        /* "TokenNetwork200.sol":32335:32341  length */
      dup9
        /* "TokenNetwork200.sol":32335:32345  length - 1 */
      sub
        /* "TokenNetwork200.sol":32322:32346  merkle_layer[length - 1] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_307
      jumpi
      invalid
    tag_307:
      swap1
      0x20
      add
      swap1
      0x20
      mul
      add
      mload
        /* "TokenNetwork200.sol":32299:32311  merkle_layer */
      dup2
        /* "TokenNetwork200.sol":32312:32318  length */
      dup9
        /* "TokenNetwork200.sol":32299:32319  merkle_layer[length] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_308
      jumpi
      invalid
    tag_308:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      add
        /* "TokenNetwork200.sol":32299:32346  merkle_layer[length] = merkle_layer[length - 1] */
      mstore
        /* "TokenNetwork200.sol":32374:32375  1 */
      0x1
        /* "TokenNetwork200.sol":32364:32375  length += 1 */
      swap7
      swap1
      swap7
      add
      swap6
        /* "TokenNetwork200.sol":32260:32390  if (length % 2 != 0) {... */
    tag_306:
        /* "TokenNetwork200.sol":32413:32414  0 */
      0x0
        /* "TokenNetwork200.sol":32409:32414  i = 0 */
      swap6
      pop
        /* "TokenNetwork200.sol":32404:32920  for (i = 0; i < length - 1; i += 2) {... */
    tag_309:
        /* "TokenNetwork200.sol":32429:32430  1 */
      0x1
        /* "TokenNetwork200.sol":32420:32426  length */
      dup8
        /* "TokenNetwork200.sol":32420:32430  length - 1 */
      sub
        /* "TokenNetwork200.sol":32416:32417  i */
      dup7
        /* "TokenNetwork200.sol":32416:32430  i < length - 1 */
      lt
        /* "TokenNetwork200.sol":32404:32920  for (i = 0; i < length - 1; i += 2) {... */
      iszero
      tag_310
      jumpi
        /* "TokenNetwork200.sol":32481:32493  merkle_layer */
      dup1
        /* "TokenNetwork200.sol":32494:32495  i */
      dup7
        /* "TokenNetwork200.sol":32498:32499  1 */
      0x1
        /* "TokenNetwork200.sol":32494:32499  i + 1 */
      add
        /* "TokenNetwork200.sol":32481:32500  merkle_layer[i + 1] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_312
      jumpi
      invalid
    tag_312:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      add
      mload
        /* "TokenNetwork200.sol":32462:32477  merkle_layer[i] */
      dup2
      mload
        /* "TokenNetwork200.sol":32462:32474  merkle_layer */
      dup3
      swap1
        /* "TokenNetwork200.sol":32475:32476  i */
      dup9
      swap1
        /* "TokenNetwork200.sol":32462:32477  merkle_layer[i] */
      dup2
      lt
      tag_313
      jumpi
      invalid
    tag_313:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      add
      mload
        /* "TokenNetwork200.sol":32462:32500  merkle_layer[i] == merkle_layer[i + 1] */
      eq
        /* "TokenNetwork200.sol":32458:32858  if (merkle_layer[i] == merkle_layer[i + 1]) {... */
      iszero
      tag_314
      jumpi
        /* "TokenNetwork200.sol":32535:32547  merkle_layer */
      dup1
        /* "TokenNetwork200.sol":32548:32549  i */
      dup7
        /* "TokenNetwork200.sol":32535:32550  merkle_layer[i] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_315
      jumpi
      invalid
    tag_315:
      swap1
      0x20
      add
      swap1
      0x20
      mul
      add
      mload
        /* "TokenNetwork200.sol":32524:32550  lockhash = merkle_layer[i] */
      swap3
      pop
        /* "TokenNetwork200.sol":32458:32858  if (merkle_layer[i] == merkle_layer[i + 1]) {... */
      jump(tag_325)
    tag_314:
        /* "TokenNetwork200.sol":32597:32609  merkle_layer */
      dup1
        /* "TokenNetwork200.sol":32610:32611  i */
      dup7
        /* "TokenNetwork200.sol":32614:32615  1 */
      0x1
        /* "TokenNetwork200.sol":32610:32615  i + 1 */
      add
        /* "TokenNetwork200.sol":32597:32616  merkle_layer[i + 1] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_317
      jumpi
      invalid
    tag_317:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      add
      mload
        /* "TokenNetwork200.sol":32579:32594  merkle_layer[i] */
      dup2
      mload
        /* "TokenNetwork200.sol":32579:32591  merkle_layer */
      dup3
      swap1
        /* "TokenNetwork200.sol":32592:32593  i */
      dup9
      swap1
        /* "TokenNetwork200.sol":32579:32594  merkle_layer[i] */
      dup2
      lt
      tag_318
      jumpi
      invalid
    tag_318:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      add
      mload
        /* "TokenNetwork200.sol":32579:32616  merkle_layer[i] < merkle_layer[i + 1] */
      lt
        /* "TokenNetwork200.sol":32575:32858  if (merkle_layer[i] < merkle_layer[i + 1]) {... */
      iszero
      tag_319
      jumpi
        /* "TokenNetwork200.sol":32678:32690  merkle_layer */
      dup1
        /* "TokenNetwork200.sol":32691:32692  i */
      dup7
        /* "TokenNetwork200.sol":32678:32693  merkle_layer[i] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_320
      jumpi
      invalid
    tag_320:
      swap1
      0x20
      add
      swap1
      0x20
      mul
      add
      mload
        /* "TokenNetwork200.sol":32695:32707  merkle_layer */
      dup2
        /* "TokenNetwork200.sol":32708:32709  i */
      dup8
        /* "TokenNetwork200.sol":32712:32713  1 */
      0x1
        /* "TokenNetwork200.sol":32708:32713  i + 1 */
      add
        /* "TokenNetwork200.sol":32695:32714  merkle_layer[i + 1] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_321
      jumpi
      invalid
    tag_321:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      dup2
      add
      mload
        /* "TokenNetwork200.sol":32661:32715  abi.encodePacked(merkle_layer[i], merkle_layer[i + 1]) */
      0x40
      dup1
      mload
      dup1
      dup5
      add
      swap5
      swap1
      swap5
      mstore
      dup4
      dup2
      add
      swap2
      swap1
      swap2
      mstore
      dup1
      mload
        /* "--CODEGEN--":26:47   */
      dup1
      dup5
      sub
        /* "--CODEGEN--":22:54   */
      dup3
      add
        /* "--CODEGEN--":6:55   */
      dup2
      mstore
        /* "TokenNetwork200.sol":32661:32715  abi.encodePacked(merkle_layer[i], merkle_layer[i + 1]) */
      0x60
      swap1
      swap4
      add
      swap1
      dup2
      swap1
      mstore
        /* "TokenNetwork200.sol":32651:32716  keccak256(abi.encodePacked(merkle_layer[i], merkle_layer[i + 1])) */
      dup3
      mload
        /* "TokenNetwork200.sol":32661:32715  abi.encodePacked(merkle_layer[i], merkle_layer[i + 1]) */
      swap1
      swap2
      dup3
      swap2
        /* "TokenNetwork200.sol":32651:32716  keccak256(abi.encodePacked(merkle_layer[i], merkle_layer[i + 1])) */
      swap1
      dup5
      add
      swap1
      dup1
        /* "TokenNetwork200.sol":32661:32715  abi.encodePacked(merkle_layer[i], merkle_layer[i + 1]) */
      dup4
        /* "TokenNetwork200.sol":32651:32716  keccak256(abi.encodePacked(merkle_layer[i], merkle_layer[i + 1])) */
      dup4
        /* "--CODEGEN--":36:189   */
    tag_322:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_323
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_322)
    tag_323:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":32651:32716  keccak256(abi.encodePacked(merkle_layer[i], merkle_layer[i + 1])) */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":32640:32716  lockhash = keccak256(abi.encodePacked(merkle_layer[i], merkle_layer[i + 1])) */
      swap3
      pop
        /* "TokenNetwork200.sol":32575:32858  if (merkle_layer[i] < merkle_layer[i + 1]) {... */
      jump(tag_325)
    tag_319:
        /* "TokenNetwork200.sol":32801:32813  merkle_layer */
      dup1
        /* "TokenNetwork200.sol":32814:32815  i */
      dup7
        /* "TokenNetwork200.sol":32818:32819  1 */
      0x1
        /* "TokenNetwork200.sol":32814:32819  i + 1 */
      add
        /* "TokenNetwork200.sol":32801:32820  merkle_layer[i + 1] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_326
      jumpi
      invalid
    tag_326:
      swap1
      0x20
      add
      swap1
      0x20
      mul
      add
      mload
        /* "TokenNetwork200.sol":32822:32834  merkle_layer */
      dup2
        /* "TokenNetwork200.sol":32835:32836  i */
      dup8
        /* "TokenNetwork200.sol":32822:32837  merkle_layer[i] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_327
      jumpi
      invalid
    tag_327:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      dup2
      add
      mload
        /* "TokenNetwork200.sol":32784:32838  abi.encodePacked(merkle_layer[i + 1], merkle_layer[i]) */
      0x40
      dup1
      mload
      dup1
      dup5
      add
      swap5
      swap1
      swap5
      mstore
      dup4
      dup2
      add
      swap2
      swap1
      swap2
      mstore
      dup1
      mload
        /* "--CODEGEN--":26:47   */
      dup1
      dup5
      sub
        /* "--CODEGEN--":22:54   */
      dup3
      add
        /* "--CODEGEN--":6:55   */
      dup2
      mstore
        /* "TokenNetwork200.sol":32784:32838  abi.encodePacked(merkle_layer[i + 1], merkle_layer[i]) */
      0x60
      swap1
      swap4
      add
      swap1
      dup2
      swap1
      mstore
        /* "TokenNetwork200.sol":32774:32839  keccak256(abi.encodePacked(merkle_layer[i + 1], merkle_layer[i])) */
      dup3
      mload
        /* "TokenNetwork200.sol":32784:32838  abi.encodePacked(merkle_layer[i + 1], merkle_layer[i]) */
      swap1
      swap2
      dup3
      swap2
        /* "TokenNetwork200.sol":32774:32839  keccak256(abi.encodePacked(merkle_layer[i + 1], merkle_layer[i])) */
      swap1
      dup5
      add
      swap1
      dup1
        /* "TokenNetwork200.sol":32784:32838  abi.encodePacked(merkle_layer[i + 1], merkle_layer[i]) */
      dup4
        /* "TokenNetwork200.sol":32774:32839  keccak256(abi.encodePacked(merkle_layer[i + 1], merkle_layer[i])) */
      dup4
        /* "--CODEGEN--":36:189   */
    tag_328:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_329
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_328)
    tag_329:
        /* "--CODEGEN--":274:275   */
      0x1
        /* "--CODEGEN--":267:270   */
      dup4
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      sub
        /* "--CODEGEN--":315:319   */
      dup1
        /* "--CODEGEN--":311:320   */
      not
        /* "--CODEGEN--":305:308   */
      dup3
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":295:321   */
      and
        /* "--CODEGEN--":356:360   */
      dup2
        /* "--CODEGEN--":350:353   */
      dup5
        /* "--CODEGEN--":344:354   */
      mload
        /* "--CODEGEN--":340:361   */
      and
        /* "--CODEGEN--":389:396   */
      dup1
        /* "--CODEGEN--":380:387   */
      dup3
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":372:375   */
      dup6
        /* "--CODEGEN--":365:398   */
      mstore
        /* "--CODEGEN--":3:402   */
      pop
      pop
      pop
        /* "TokenNetwork200.sol":32774:32839  keccak256(abi.encodePacked(merkle_layer[i + 1], merkle_layer[i])) */
      pop
      pop
      pop
      swap1
      pop
      add
      swap2
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      keccak256
        /* "TokenNetwork200.sol":32763:32839  lockhash = keccak256(abi.encodePacked(merkle_layer[i + 1], merkle_layer[i])) */
      swap3
      pop
        /* "TokenNetwork200.sol":32575:32858  if (merkle_layer[i] < merkle_layer[i + 1]) {... */
    tag_325:
        /* "TokenNetwork200.sol":32897:32905  lockhash */
      dup3
        /* "TokenNetwork200.sol":32875:32887  merkle_layer */
      dup2
        /* "TokenNetwork200.sol":32892:32893  2 */
      0x2
        /* "TokenNetwork200.sol":32888:32889  i */
      dup9
        /* "TokenNetwork200.sol":32888:32893  i / 2 */
      div
        /* "TokenNetwork200.sol":32875:32894  merkle_layer[i / 2] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_332
      jumpi
      invalid
    tag_332:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      add
        /* "TokenNetwork200.sol":32875:32905  merkle_layer[i / 2] = lockhash */
      mstore
        /* "TokenNetwork200.sol":32437:32438  2 */
      0x2
        /* "TokenNetwork200.sol":32432:32438  i += 2 */
      swap6
      swap1
      swap6
      add
      swap5
        /* "TokenNetwork200.sol":32404:32920  for (i = 0; i < length - 1; i += 2) {... */
      jump(tag_309)
    tag_310:
        /* "TokenNetwork200.sol":32946:32947  2 */
      0x2
        /* "TokenNetwork200.sol":32942:32943  i */
      dup7
        /* "TokenNetwork200.sol":32942:32947  i / 2 */
      div
        /* "TokenNetwork200.sol":32933:32947  length = i / 2 */
      swap7
      pop
        /* "TokenNetwork200.sol":32227:32958  while (length > 1) {... */
      jump(tag_303)
    tag_304:
        /* "TokenNetwork200.sol":32982:32994  merkle_layer */
      dup1
        /* "TokenNetwork200.sol":32995:32996  0 */
      0x0
        /* "TokenNetwork200.sol":32982:32997  merkle_layer[0] */
      dup2
      mload
      dup2
      lt
      iszero
      iszero
      tag_334
      jumpi
      invalid
    tag_334:
      0x20
      swap1
      dup2
      mul
      swap1
      swap2
      add
      add
      mload
      swap11
        /* "TokenNetwork200.sol":33029:33050  total_unlocked_amount */
      swap5
      swap10
      pop
        /* "TokenNetwork200.sol":31382:33058  function getMerkleRootAndUnlockedAmount(bytes merkle_tree_leaves)... */
      swap4
      swap8
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      pop
      jump	// out
        /* "TokenNetwork200.sol":33064:34348  function getLockDataFromMerkleTree(bytes merkle_tree_leaves, uint256 offset)... */
    tag_299:
        /* "TokenNetwork200.sol":33176:33183  bytes32 */
      0x0
        /* "TokenNetwork200.sol":33185:33192  uint256 */
      dup1
        /* "TokenNetwork200.sol":33208:33232  uint256 expiration_block */
      0x0
        /* "TokenNetwork200.sol":33242:33263  uint256 locked_amount */
      dup1
        /* "TokenNetwork200.sol":33273:33293  uint256 reveal_block */
      0x0
        /* "TokenNetwork200.sol":33303:33321  bytes32 secrethash */
      dup1
        /* "TokenNetwork200.sol":33331:33347  bytes32 lockhash */
      0x0
        /* "TokenNetwork200.sol":33391:33397  offset */
      dup8
        /* "TokenNetwork200.sol":33362:33380  merkle_tree_leaves */
      dup10
        /* "TokenNetwork200.sol":33362:33387  merkle_tree_leaves.length */
      mload
        /* "TokenNetwork200.sol":33362:33397  merkle_tree_leaves.length <= offset */
      gt
      iszero
        /* "TokenNetwork200.sol":33358:33444  if (merkle_tree_leaves.length <= offset) {... */
      iszero
      tag_336
      jumpi
        /* "TokenNetwork200.sol":33421:33429  lockhash */
      swap6
      pop
        /* "TokenNetwork200.sol":33431:33432  0 */
      0x0
      swap5
      pop
        /* "TokenNetwork200.sol":33421:33429  lockhash */
      dup6
        /* "TokenNetwork200.sol":33413:33433  return (lockhash, 0) */
      jump(tag_335)
        /* "TokenNetwork200.sol":33358:33444  if (merkle_tree_leaves.length <= offset) {... */
    tag_336:
        /* "TokenNetwork200.sol":33503:33534  add(merkle_tree_leaves, offset) */
      dup9
      dup9
      add
        /* "TokenNetwork200.sol":33497:33535  mload(add(merkle_tree_leaves, offset)) */
      dup1
      mload
        /* "TokenNetwork200.sol":33607:33609  32 */
      0x20
        /* "TokenNetwork200.sol":33571:33611  add(merkle_tree_leaves, add(offset, 32)) */
      dup1
      dup4
      add
        /* "TokenNetwork200.sol":33565:33612  mload(add(merkle_tree_leaves, add(offset, 32))) */
      mload
        /* "TokenNetwork200.sol":33681:33683  64 */
      0x40
        /* "TokenNetwork200.sol":33645:33685  add(merkle_tree_leaves, add(offset, 64)) */
      swap4
      dup5
      add
        /* "TokenNetwork200.sol":33639:33686  mload(add(merkle_tree_leaves, add(offset, 64))) */
      mload
        /* "TokenNetwork200.sol":33791:33852  abi.encodePacked(expiration_block, locked_amount, secrethash) */
      dup5
      mload
      dup1
      dup5
      add
      dup6
      swap1
      mstore
      dup1
      dup7
      add
      dup4
      swap1
      mstore
      0x60
      dup1
      dup3
      add
      dup4
      swap1
      mstore
      dup7
      mload
        /* "--CODEGEN--":26:47   */
      dup1
      dup4
      sub
        /* "--CODEGEN--":22:54   */
      swap1
      swap2
      add
        /* "--CODEGEN--":6:55   */
      dup2
      mstore
        /* "TokenNetwork200.sol":33791:33852  abi.encodePacked(expiration_block, locked_amount, secrethash) */
      0x80
      swap1
      swap2
      add
      swap6
      dup7
      swap1
      mstore
        /* "TokenNetwork200.sol":33781:33853  keccak256(abi.encodePacked(expiration_block, locked_amount, secrethash)) */
      dup1
      mload
        /* "TokenNetwork200.sol":33497:33535  mload(add(merkle_tree_leaves, offset)) */
      swap5
      swap11
      pop
        /* "TokenNetwork200.sol":33565:33612  mload(add(merkle_tree_leaves, add(offset, 32))) */
      swap2
      swap9
      pop
        /* "TokenNetwork200.sol":33639:33686  mload(add(merkle_tree_leaves, add(offset, 64))) */
      swap6
      pop
        /* "TokenNetwork200.sol":33791:33852  abi.encodePacked(expiration_block, locked_amount, secrethash) */
      swap3
      swap2
      dup3
      swap2
        /* "TokenNetwork200.sol":33781:33853  keccak256(abi.encodePacked(expiration_block, locked_amount, secrethash)) */
      dup5
      add
      swap1
      dup1
        /* "TokenNetwork200.sol":33791:33852  abi.encodePacked(expiration_block, locked_amount, secrethash) */
      dup4
        /* "TokenNetwork200.sol":33781:33853  keccak256(abi.encodePacked(expiration_block, locked_amount, secrethash)) */
      dup4
        /* "--CODEGEN--":36:189   */
    tag_337:
        /* "--CODEGEN--":66:68   */
      0x20
        /* "--CODEGEN--":58:69   */
      dup4
      lt
        /* "--CODEGEN--":36:189   */
      tag_338
      jumpi
        /* "--CODEGEN--":176:186   */
      dup1
      mload
        /* "--CODEGEN--":164:187   */
      dup3
      mstore
        /* "--CODEGEN--":139:151   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0
      swap1
      swap3
      add
      swap2
        /* "--CODEGEN--":98:100   */
      0x20
        /* "--CODEGEN--":89:101   */
      swap2
      dup3
      add
      swap2
        /* "--CODEGEN--":114:126   */
      add
        /* "--CODEGEN--":36:189   */
      jump(tag_337)
    tag_338:
        /* "--CODEGEN--":299:309   */
      mload
        /* "--CODEGEN--":344:354   */
      dup2
      mload
        /* "--CODEGEN--":263:265   */
      0x20
        /* "--CODEGEN--":259:271   */
      swap4
      dup5
      sub
        /* "--CODEGEN--":254:257   */
      0x100
        /* "--CODEGEN--":250:272   */
      exp
        /* "--CODEGEN--":246:276   */
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
      add
        /* "--CODEGEN--":311:320   */
      dup1
      not
        /* "--CODEGEN--":295:321   */
      swap1
      swap3
      and
        /* "--CODEGEN--":340:361   */
      swap2
      and
        /* "--CODEGEN--":377:397   */
      or
        /* "--CODEGEN--":365:398   */
      swap1
      mstore
        /* "TokenNetwork200.sol":33781:33853  keccak256(abi.encodePacked(expiration_block, locked_amount, secrethash)) */
      0x40
      dup1
      mload
      swap3
      swap1
      swap5
      add
      dup3
      swap1
      sub
      dup3
      keccak256
        /* "--CODEGEN--":274:275   */
      0x1
        /* "TokenNetwork200.sol":34134:34149  secret_registry */
      sload
        /* "TokenNetwork200.sol":34134:34188  secret_registry.getSecretRevealBlockHeight(secrethash) */
      0xc1f6294600000000000000000000000000000000000000000000000000000000
      dup5
      mstore
      0x4
      dup5
      add
      dup11
      swap1
      mstore
      swap5
      mload
        /* "TokenNetwork200.sol":33781:33853  keccak256(abi.encodePacked(expiration_block, locked_amount, secrethash)) */
      swap1
      swap8
      pop
        /* "TokenNetwork200.sol":34134:34149  secret_registry */
      0xffffffffffffffffffffffffffffffffffffffff
      swap1
      swap5
      and
      swap6
      pop
        /* "TokenNetwork200.sol":34134:34176  secret_registry.getSecretRevealBlockHeight */
      0xc1f62946
      swap5
      pop
        /* "TokenNetwork200.sol":34134:34188  secret_registry.getSecretRevealBlockHeight(secrethash) */
      0x24
      dup1
      dup4
      add
      swap5
        /* "--CODEGEN--":263:265   */
      swap2
      swap4
      pop
        /* "TokenNetwork200.sol":34134:34188  secret_registry.getSecretRevealBlockHeight(secrethash) */
      swap1
      swap2
      dup3
      swap1
      sub
      add
      dup2
      0x0
        /* "TokenNetwork200.sol":34134:34149  secret_registry */
      dup8
        /* "TokenNetwork200.sol":34134:34188  secret_registry.getSecretRevealBlockHeight(secrethash) */
      dup1
      extcodesize
      iszero
        /* "--CODEGEN--":5:7   */
      dup1
      iszero
      tag_340
      jumpi
        /* "--CODEGEN--":30:31   */
      0x0
        /* "--CODEGEN--":27:28   */
      dup1
        /* "--CODEGEN--":20:32   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_340:
        /* "TokenNetwork200.sol":34134:34188  secret_registry.getSecretRevealBlockHeight(secrethash) */
      pop
      gas
      call
      iszero
        /* "--CODEGEN--":8:17   */
      dup1
        /* "--CODEGEN--":5:7   */
      iszero
      tag_341
      jumpi
        /* "--CODEGEN--":45:61   */
      returndatasize
        /* "--CODEGEN--":42:43   */
      0x0
        /* "--CODEGEN--":39:40   */
      dup1
        /* "--CODEGEN--":24:62   */
      returndatacopy
        /* "--CODEGEN--":77:93   */
      returndatasize
        /* "--CODEGEN--":74:75   */
      0x0
        /* "--CODEGEN--":67:94   */
      revert
        /* "--CODEGEN--":5:7   */
    tag_341:
        /* "TokenNetwork200.sol":34134:34188  secret_registry.getSecretRevealBlockHeight(secrethash) */
      pop
      pop
      pop
      pop
      mload(0x40)
      returndatasize
        /* "--CODEGEN--":13:15   */
      0x20
        /* "--CODEGEN--":8:11   */
      dup2
        /* "--CODEGEN--":5:16   */
      lt
        /* "--CODEGEN--":2:4   */
      iszero
      tag_342
      jumpi
        /* "--CODEGEN--":29:30   */
      0x0
        /* "--CODEGEN--":26:27   */
      dup1
        /* "--CODEGEN--":19:31   */
      revert
        /* "--CODEGEN--":2:4   */
    tag_342:
      pop
        /* "TokenNetwork200.sol":34134:34188  secret_registry.getSecretRevealBlockHeight(secrethash) */
      mload
      swap3
      pop
        /* "TokenNetwork200.sol":34202:34219  reveal_block == 0 */
      dup3
      iszero
      dup1
        /* "TokenNetwork200.sol":34202:34255  reveal_block == 0 || expiration_block <= reveal_block */
      tag_343
      jumpi
      pop
        /* "TokenNetwork200.sol":34243:34255  reveal_block */
      dup3
        /* "TokenNetwork200.sol":34223:34239  expiration_block */
      dup6
        /* "TokenNetwork200.sol":34223:34255  expiration_block <= reveal_block */
      gt
      iszero
        /* "TokenNetwork200.sol":34202:34255  reveal_block == 0 || expiration_block <= reveal_block */
    tag_343:
        /* "TokenNetwork200.sol":34198:34299  if (reveal_block == 0 || expiration_block <= reveal_block) {... */
      iszero
      tag_344
      jumpi
        /* "TokenNetwork200.sol":34287:34288  0 */
      0x0
        /* "TokenNetwork200.sol":34271:34288  locked_amount = 0 */
      swap4
      pop
        /* "TokenNetwork200.sol":34198:34299  if (reveal_block == 0 || expiration_block <= reveal_block) {... */
    tag_344:
        /* "TokenNetwork200.sol":34317:34325  lockhash */
      dup1
        /* "TokenNetwork200.sol":34327:34340  locked_amount */
      dup5
        /* "TokenNetwork200.sol":34309:34341  return (lockhash, locked_amount) */
      swap7
      pop
      swap7
      pop
        /* "TokenNetwork200.sol":33064:34348  function getLockDataFromMerkleTree(bytes merkle_tree_leaves, uint256 offset)... */
    tag_335:
      pop
      pop
      pop
      pop
      pop
      swap3
      pop
      swap3
      swap1
      pop
      jump	// out

    auxdata: 0xa165627a7a7230582051528a1cf43bd82f7b42336f06d11503f20e8237f740221c349c24136c8d8cf20029
}


======= Utils.sol:Utils =======
EVM assembly:
    /* "Utils.sol":26:507  contract Utils {... */
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
    /* "Utils.sol":26:507  contract Utils {... */
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
        /* "Utils.sol":26:507  contract Utils {... */
      mstore(0x40, 0x80)
      jumpi(tag_1, lt(calldatasize, 0x4))
      and(div(calldataload(0x0), 0x100000000000000000000000000000000000000000000000000000000), 0xffffffff)
      0x7709bc78
      dup2
      eq
      tag_2
      jumpi
      dup1
      0xb32c65c8
      eq
      tag_3
      jumpi
    tag_1:
      0x0
      dup1
      revert
        /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
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
      pop
        /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
      tag_5
      and(calldataload(0x4), 0xffffffffffffffffffffffffffffffffffffffff)
      jump(tag_6)
    tag_5:
      0x40
      dup1
      mload
      swap2
      iszero
      iszero
      dup3
      mstore
      mload
      swap1
      dup2
      swap1
      sub
      0x20
      add
      swap1
      return
        /* "Utils.sol":47:96  string constant public contract_version = "0.3._" */
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
        /* "Utils.sol":47:96  string constant public contract_version = "0.3._" */
      pop
      tag_8
      jump(tag_9)
    tag_8:
      0x40
      dup1
      mload
      0x20
      dup1
      dup3
      mstore
      dup4
      mload
      dup2
      dup4
      add
      mstore
      dup4
      mload
      swap2
      swap3
      dup4
      swap3
      swap1
      dup4
      add
      swap2
      dup6
      add
      swap1
      dup1
      dup4
      dup4
      0x0
        /* "--CODEGEN--":8:108   */
    tag_10:
        /* "--CODEGEN--":33:36   */
      dup4
        /* "--CODEGEN--":30:31   */
      dup2
        /* "--CODEGEN--":27:37   */
      lt
        /* "--CODEGEN--":8:108   */
      iszero
      tag_11
      jumpi
        /* "--CODEGEN--":90:101   */
      dup2
      dup2
      add
        /* "--CODEGEN--":84:102   */
      mload
        /* "--CODEGEN--":71:82   */
      dup4
      dup3
      add
        /* "--CODEGEN--":64:103   */
      mstore
        /* "--CODEGEN--":52:54   */
      0x20
        /* "--CODEGEN--":45:55   */
      add
        /* "--CODEGEN--":8:108   */
      jump(tag_10)
    tag_11:
        /* "--CODEGEN--":12:26   */
      pop
        /* "Utils.sol":47:96  string constant public contract_version = "0.3._" */
      pop
      pop
      pop
      swap1
      pop
      swap1
      dup2
      add
      swap1
      0x1f
      and
      dup1
      iszero
      tag_13
      jumpi
      dup1
      dup3
      sub
      dup1
      mload
      0x1
      dup4
      0x20
      sub
      0x100
      exp
      sub
      not
      and
      dup2
      mstore
      0x20
      add
      swap2
      pop
    tag_13:
      pop
      swap3
      pop
      pop
      pop
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      return
        /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
    tag_6:
        /* "Utils.sol":367:371  bool */
      0x0
        /* "Utils.sol":434:463  extcodesize(contract_address) */
      swap1
      extcodesize
        /* "Utils.sol":490:498  size > 0 */
      gt
      swap1
        /* "Utils.sol":296:505  function contractExists(address contract_address) public view returns (bool) {... */
      jump	// out
        /* "Utils.sol":47:96  string constant public contract_version = "0.3._" */
    tag_9:
      0x40
      dup1
      mload
      dup1
      dup3
      add
      swap1
      swap2
      mstore
      0x5
      dup2
      mstore
      0x302e332e5f000000000000000000000000000000000000000000000000000000
      0x20
      dup3
      add
      mstore
      dup2
      jump	// out

    auxdata: 0xa165627a7a723058207bd06324e5db36a2af32694c0c3430898df2b9e3072ec0f6cadc5308df1f0a490029
}

