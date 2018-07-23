package params

//NameTokenNetworkCreated event TokenNetworkCreated(address indexed token_address, address indexed token_network_address);
const NameTokenNetworkCreated = "TokenNetworkCreated"

//NameChannelOpened new channel event of token network
const NameChannelOpened = "ChannelOpened"

//NameChannelNewDeposit deposit event of token network
const NameChannelNewDeposit = "ChannelNewDeposit"

//NameChannelWithdraw withdraw event of token network
const NameChannelWithdraw = "ChannelWithdraw"

//NameChannelClosed event ChannelClosed(bytes32 indexed channel_identifier, address indexed closing_participant);
const NameChannelClosed = "ChannelClosed"

//NameChannelPunished punish event of token network
const NameChannelPunished = "ChannelPunished"

//NameChannelUnlocked unlock event of token network
const NameChannelUnlocked = "ChannelUnlocked"

//NameBalanceProofUpdated  update balance proof event of token network
const NameBalanceProofUpdated = "BalanceProofUpdated"

//NameChannelSettled  settle channel event of token network
const NameChannelSettled = "ChannelSettled"

//NameChannelCooperativeSettled represents channel cooperatively settled
const NameChannelCooperativeSettled = "ChannelCooperativeSettled"

//NameSecretRevealed name from contract
const NameSecretRevealed = "SecretRevealed"

//name of Monitoring Service

//NameNewDeposit event NewDeposit(address indexed receiver, uint amount);
const NameNewDeposit = "NewDeposit"

//NameNewBalanceProofReceived event of monitoring service
const NameNewBalanceProofReceived = "NewBalanceProofReceived"

//NameRewardClaimed event RewardClaimed(address indexed ms_address, uint amount, bytes32 indexed reward_identifier);
const NameRewardClaimed = "RewardClaimed"

//NameWithdrawn event Withdrawn(address indexed account, uint amount);
const NameWithdrawn = "Withdrawn"
