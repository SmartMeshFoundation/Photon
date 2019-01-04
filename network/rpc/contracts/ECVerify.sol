pragma solidity ^0.4.24;

/// @title ECVerify Library
/// @notice A reusable library for verifying whether a message is from a specific participant.
library ECVerify {

    /// @notice It is to verify whether provided hash is corresponding to an address.
    /// @dev signature has to be 65-byte long, dividing into two 32-bytes and 1 byte structure.
    /// @param hash                 a 32-byte hash value to be verified
    /// @param signature            a dynamic sized bytes of signature
    /// @return signature_address   a 20-byte signature address.
    function ecverify(bytes32 hash, bytes memory signature)
        internal
        pure
        returns (address signature_address)
    {
        require(signature.length == 65);

        bytes32 r;
        bytes32 s;
        uint8 v;

        // The signature format is a compact form of:
        //   {bytes32 r}{bytes32 s}{uint8 v}
        // Compact means, uint8 is not padded to 32 bytes.
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))

            // Here we are loading the last 32 bytes, including 31 bytes of 's'.
            v := byte(0, mload(add(signature, 96)))
        }

        // Version of signature should be 27 or 28, but 0 and 1 are also possible
        if (v < 27) {
            v += 27;
        }

        require(v == 27 || v == 28);

        signature_address = ecrecover(hash, v, r, s);

        // ecrecover returns zero on error
        require(signature_address != 0x0);

        return signature_address;
    }
}
