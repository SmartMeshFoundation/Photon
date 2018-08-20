#!/bin/sh
abigen --sol TokenNetworkRegistry.sol --pkg contracts --out TokenNetworkRegistry.go -solc ./solc
abigen -solc ./solc --sol test/tokens/HumanStandardToken.sol --pkg tokenstandard --out test/tokens/tokenstandard/HumanStandardToken.go
abigen -solc ./solc --sol test/tokens/CustomToken.sol --pkg tokencustom --out test/tokens/tokencustom/CustomToken.go
abigen -solc ./solc --sol test/tokens/HumanERC223Token.sol --pkg tokenerc223 --out test/tokens/tokenerc223/HumanERC223Token.go
abigen -solc ./solc --sol test/tokens/HumanERC223ApproveToken.sol --pkg tokenerc223approve --out test/tokens/tokenerc223approve/HumanERC223ApproveToken.go
abigen -solc ./solc --sol test/tokens/HumanEtherToken.sol --pkg tokenether --out test/tokens/tokenether/HumanEtherToken.go
