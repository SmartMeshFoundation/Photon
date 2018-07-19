package params

//NameTokenNetworkCreated event TokenNetworkCreated(address indexed token_address, address indexed token_network_address);
const NameTokenNetworkCreated = "TokenNetworkCreated"

const NameChannelOpened = "ChannelOpened"

const NameChannelNewDeposit = "ChannelNewDeposit"

const NameChannelWithdraw = "ChannelWithdraw"

//NameChannelClosed event ChannelClosed(bytes32 indexed channel_identifier, address indexed closing_participant);
const NameChannelClosed = "ChannelClosed"

const NameChannelPunished = "ChannelPunished"

const NameChannelUnlocked = "ChannelUnlocked"

const NameBalanceProofUpdated = "BalanceProofUpdated"

const NameChannelSettled = "ChannelSettled"

//NameChannelCooperativeSettled represents channel cooperatively settled
const NameChannelCooperativeSettled = "ChannelCooperativeSettled"

//NameSecretRevealed name from contract
const NameSecretRevealed = "SecretRevealed"

//name of Monitoring Service

//NameNewDeposit event NewDeposit(address indexed receiver, uint amount);
const NameNewDeposit = "NewDeposit"

const NameNewBalanceProofReceived = "NewBalanceProofReceived"

//NameRewardClaimed event RewardClaimed(address indexed ms_address, uint amount, bytes32 indexed reward_identifier);
const NameRewardClaimed = "RewardClaimed"

//NameWithdrawn event Withdrawn(address indexed account, uint amount);
const NameWithdrawn = "Withdrawn"
