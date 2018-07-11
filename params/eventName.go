package params

//NameTokenNetworkCreated event TokenNetworkCreated(address indexed token_address, address indexed token_network_address);
const NameTokenNetworkCreated = "TokenNetworkCreated"

//NameChannelOpened event ChannelOpened(
//        bytes32 indexed channel_identifier,
//        address indexed participant1,
//        address indexed participant2,
//        uint256 settle_timeout
//    );
const NameChannelOpened = "ChannelOpened"

//NameChannelNewDeposit event ChannelNewDeposit(
//        bytes32 indexed channel_identifier,
//        address indexed participant,
//        uint256 total_deposit
//    );
const NameChannelNewDeposit = "ChannelNewDeposit"

//NameChannelWithdraw event ChannelWithdraw(
//        bytes32 indexed channel_identifier,
//        address indexed participant, uint256 total_withdraw
//    );
const NameChannelWithdraw = "ChannelWithdraw"

//NameChannelClosed event ChannelClosed(bytes32 indexed channel_identifier, address indexed closing_participant);
const NameChannelClosed = "ChannelClosed"

//NameChannelUnlocked event ChannelUnlocked(
//        bytes32 indexed channel_identifier,
//        address indexed participant,
//        uint256 unlocked_amount,
//        uint256 returned_tokens
//    );
const NameChannelUnlocked = "ChannelUnlocked"

//NameBalanceProofUpdated event NonClosingBalanceProofUpdated(
//        bytes32 indexed channel_identifier,
//        address indexed closing_participant
//    );
const NameBalanceProofUpdated = "BalanceProofUpdated"

//NameChannelSettled event ChannelSettled(
//        bytes32 indexed channel_identifier,
//        uint256 participant1_amount,
//        uint256 participant2_amount
//    );
const NameChannelSettled = "ChannelSettled"

//NameChannelCooperativeSettled represents channel cooperatively settled
const NameChannelCooperativeSettled = "ChannelCooperativeSettled"

//NameSecretRevealed name from contract
const NameSecretRevealed = "SecretRevealed"

//name of Monitoring Service

//NameNewDeposit event NewDeposit(address indexed receiver, uint amount);
const NameNewDeposit = "NewDeposit"

/*
NameNewBalanceProofReceived event NewBalanceProofReceived(
        uint256 reward_amount,
        uint256 indexed nonce,
        address indexed ms_address,
        address indexed raiden_node_address
    );
*/
const NameNewBalanceProofReceived = "NewBalanceProofReceived"

//NameRewardClaimed event RewardClaimed(address indexed ms_address, uint amount, bytes32 indexed reward_identifier);
const NameRewardClaimed = "RewardClaimed"

//NameWithdrawn event Withdrawn(address indexed account, uint amount);
const NameWithdrawn = "Withdrawn"
