pragma solidity ^0.4.23;

import "./Token.sol";
import "./Utils.sol";
import "./ECVerify.sol";
import "./SecretRegistry.sol";

/*
相比17改动:
主要是把 Participant 中的 initialized 去掉,
都是通过getChannelIdentifier 来获取 channel id, 保证初始化过.
*/
contract TokenNetwork is Utils {

    /*
     *  Data structures
     */

    string constant public contract_version = "0.3._";

    // Instance of the token used as digital currency by the channels
    Token public token;
    // Instance of SecretRegistry used for storing secrets revealed in a mediating transfer.
    SecretRegistry public secret_registry;
    // Chain ID as specified by EIP155 used in balance proof signatures to avoid replay attacks
    uint256 public chain_id;
    // Channel identifier is a uint256, incremented after each new channel
    mapping(uint256 => Channel) public channels;
    mapping(bytes32 => uint256) public  openedchannels;
    // Start from 1 instead of 0, otherwise the first channel will have an additional
    // 15000 gas cost than the rest
    uint256 public last_channel_index = 0;

    struct Participant {
        // Total amount of token transferred to this smart contract
        uint256 deposit;
        // The latest known merkle root of the pending hash-time locks, used to
        // validate the withdrawn proofs.
        bytes32 locksroot;
        // The latest known transferred_amount from this node to the other
        // participant, used to compute the net balance on settlement.
        uint256 transferred_amount;
        // Monotonically increasing counter of the off-chain transfers
        // Value used to order transfers and only accept the latest on calls to
        // update
        uint256 nonce;
        /*
        对方通过密码解锁的数量总和.
        */
        uint256 unlocked_amount;
    }

    struct Channel {
        // After the channel has been uncooperatively closed, this value represents the
        // block number after which settleChannel can be called.
        uint256 settle_block_number;
        address closing_participant;
        // Channel state
        // 1 = open, 2 = closed
        // 0 = non-existent or settled
        uint8 state;
        mapping(address => Participant) participants;
    }

    /*
*  Events
*/
    event ChannelOpened(
        uint256 indexed channel_identifier,
        address indexed participant1,
        address indexed participant2,
        uint256 settle_timeout
    );

    event ChannelNewDeposit(
        uint256 indexed channel_identifier,
        address indexed participant,
        uint256 total_deposit
    );

    event ChannelClosed(uint256 indexed channel_identifier, address indexed closing_participant);
    event ChannelUnlocked(
        uint256 channel_identifier,
        address payer_participant,
        bytes32 locskroot, //解锁的 locksroot
        uint256 transferred_amount
    );

    event NonClosingBalanceProofUpdated(
        uint256 indexed channel_identifier,
        address indexed closing_participant
    );

    event ChannelSettled(
        uint256 indexed channel_identifier,
        uint256 participant1_amount,
        uint256 participant2_amount
    );
    event Channelwithdraw(
        uint256 channel_identifier,
        uint256 participant1_deposit,
        uint256 participant2_deposit,
        uint256 participant1_withdraw,
        uint256 participant2_withdraw
    );
    /*
 * Modifiers
 */

    modifier isOpen(uint256 channel_identifier) {
        require(channels[channel_identifier].state == 1);
        _;
    }

    modifier isClosed(uint256 channel_identifier) {
        require(channels[channel_identifier].state == 2);
        _;
    }

    modifier settleTimeoutValid(uint256 timeout) {
        require(timeout >= 6 && timeout <= 2700000);
        _;
    }
    /*
   *  Constructor
   */

    constructor(address _token_address, address _secret_registry, uint256 _chain_id)
    public
    {
        require(_token_address != 0x0);
        require(_secret_registry != 0x0);
        require(_chain_id > 0);
        require(contractExists(_token_address));
        require(contractExists(_secret_registry));

        token = Token(_token_address);

        secret_registry = SecretRegistry(_secret_registry);
        chain_id = _chain_id;

        // Make sure the contract is indeed a token contract
        require(token.totalSupply() > 0);
    }

    /*
     *  Public functions
     */


    function openChannel(address participant1, address participant2, uint256 settle_timeout)
    settleTimeoutValid(settle_timeout)
    public
    {
        uint256 channel_identifier;
        bytes32 channel_hash;
        require(participant1 != 0x0);
        require(participant2 != 0x0);
        require(participant1 != participant2);
        last_channel_index += 1;
        channel_identifier = last_channel_index;
        Channel storage channel = channels[channel_identifier];

        // Store channel information
        channel.settle_block_number = settle_timeout;
        // Mark channel as opened
        channel.state = 1;
        channel_hash = getChannelHash(participant1, participant2);
        //不允许有两个节点重复重建通道.
        require(openedchannels[channel_hash] == 0);
        openedchannels[channel_hash] = channel_identifier;
        emit ChannelOpened(channel_identifier, participant1, participant2, settle_timeout);
    }
/*
必须在通道 open 状态调用,可以重复调用多次,任何人都可以调用.
total_deposit 是为了防止重放.
*/
    function setTotalDeposit(address participant,address partner, uint256 total_deposit)
    public
    {
        require(total_deposit > 0);
        uint256 added_deposit;
        uint256 channel_identifier;
        channel_identifier=getChannelIdentifier(participant,partner);
        Channel storage channel = channels[channel_identifier];
        Participant storage participant_state = channel.participants[participant];
        require(participant_state.deposit < total_deposit);
        // Calculate the actual amount of tokens that will be transferred
        added_deposit = total_deposit - participant_state.deposit;

        // Update the participant's channel deposit
        participant_state.deposit += added_deposit;


        // Do the transfer
        require(token.transferFrom(msg.sender, address(this), added_deposit));
        require(channel.state==1);
        emit ChannelNewDeposit(channel_identifier, participant, participant_state.deposit);
        //如果 token 可能的 totalSupply 大于 uint256,说明这个 token 分文不值,分文不值的 token 发生什么都无所谓.
        //require(participant_state.deposit >= added_deposit);
        //防止溢出,有必要么?我是想不到原因.
        //require(channel_deposit >= participant_state.deposit);
        //require(channel_deposit >= partner_state.deposit);
    }
/*
任何人都可以调用,调用一次相当于新创建了通道,所以无法重放攻击
*/
    function withDraw(
        address participant1,
        address participant2,
        uint256 participant1_deposit,
        uint256 participant2_deposit,
        uint256 participant1_withdraw,
        uint256 participant2_withdraw,
        bytes participant1_signature,
        bytes participant2_signature
    )
    public
    {
        uint256 total_deposit;
        uint256 channel_identifier;
        channel_identifier=getChannelIdentifier(participant1,participant2);
        Channel storage channel = channels[channel_identifier];
        require(channel.state == 1);

        //验证发送方签名
        bytes32 message_hash = keccak256(abi.encodePacked(
                participant1,
                participant1_deposit,
                participant2,
                participant2_deposit,
                participant1_withdraw,
                channel_identifier,
                address(this),
                chain_id
            ));
        require(participant1 == ECVerify.ecverify(message_hash, participant1_signature));
        //验证接收方签名
        message_hash = keccak256(abi.encodePacked(
                participant1,
                participant1_deposit,
                participant2,
                participant2_deposit,
                participant1_withdraw,
                participant2_withdraw,
                channel_identifier,
                address(this),
                chain_id
            ));
        require(participant2 == ECVerify.ecverify(message_hash, participant2_signature));
        Participant storage participant1_state = channel.participants[participant1];
        Participant storage participant2_state = channel.participants[participant2];
        //The sum of the provided deposit must be equal to the total available deposit
        total_deposit = participant1_state.deposit + participant2_state.deposit;
        require(participant1_deposit <= total_deposit);
        require(participant2_deposit <= total_deposit);
        require((participant1_deposit + participant2_deposit) == total_deposit);

        // Do the token transfers
        if (participant1_withdraw > 0) {
            require(token.transfer(participant1, participant1_withdraw));
        }
        if (participant2_withdraw > 0) {
            require(token.transfer(participant2, participant2_withdraw));
        }
        require(participant1_withdraw <= participant1_deposit);
        require(participant2_withdraw <= participant2_deposit);
        participant1_deposit = participant1_deposit - participant1_withdraw;
        participant2_deposit = participant2_deposit - participant2_withdraw;
        openNew(participant1,participant2,participant1_deposit,participant2_deposit,channel.settle_block_number);

        delete channel.participants[participant1];
        delete channel.participants[participant2];
        delete channels[channel_identifier];
        openedchannels[message_hash]=last_channel_index;

        emit Channelwithdraw(channel_identifier, participant1_deposit, participant2_deposit, participant1_withdraw, participant2_withdraw);

    }

    function openNew(address participant1,address participant2,uint256 participant1_deposit,uint256 participant2_deposit,uint256 settle_timeout) internal {
        last_channel_index+=1;
        Channel storage channel=channels[last_channel_index];
        channel.state=1;
        channel.settle_block_number=settle_timeout;
        Participant storage participant1_state=channel.participants[participant1];
        participant1_state.deposit=participant1_deposit;
        Participant storage participant2_state=channel.participants[participant2];
        participant2_state.deposit=participant2_deposit;
    }
    /*
    只能是通道参与方调用,只能调用一次,必须是在通道打开状态调用.
    */
    function closeChannel(
        address partner,
        uint256 transferred_amount,
        bytes32 locksroot,
        uint256 nonce,
        bytes32 additional_hash,
        bytes signature
    )
    public
    {
        uint256 channel_identifier; 
        address recovered_partner_address;
        channel_identifier=getChannelIdentifier(msg.sender,partner);
        Channel storage channel = channels[channel_identifier];
        require(channel.state==1);
        // Mark the channel as closed and mark the closing participant
        channel.state = 2;
        channel.closing_participant = msg.sender;
        // This is the block number at which the channel can be settled.
        channel.settle_block_number += uint256(block.number);
        // Nonce 0 means that the closer never received a transfer, therefore never received a
        // balance proof, or he is intentionally not providing the latest transfer, in which case
        // the closing party is going to lose the tokens that were transferred to him.
        if (nonce > 0) {
            recovered_partner_address = recoverAddressFromBalanceProof(
                channel_identifier,
                transferred_amount,
                locksroot,
                nonce,
                additional_hash,
                signature
            );
            updateBalanceProofData(channel_identifier, recovered_partner_address, nonce, locksroot, transferred_amount);
            require(partner==recovered_partner_address);
        }
        emit ChannelClosed(channel_identifier, msg.sender);
    }
    /*
    任何人都可以调用,可以调用多次,只要在有效期内.
    */
    function updateNonClosingBalanceProof(
        address closing_participant,
        address non_closing_participant,
        bytes32 locksroot,
        uint256 transferred_amount,
        uint256 nonce,
        bytes32 additional_hash,
        bytes closing_signature,
        bytes non_closing_signature
    )
    public
    {
        uint256 channel_identifier;
        channel_identifier=getChannelIdentifier(closing_participant,non_closing_participant);
        Channel storage channel = channels[channel_identifier];
        require(channel.state==2);
        require(channel.settle_block_number >= block.number);
        require(nonce > 0);
        // 防止closing 冒充 non_closing_participant,冒充了会得利么?相当于可以重复更新了,再次提交新的证据?
        // 所以,如果允许 closing 一方再次提交证据,是不是可以把channel.closing_participant删除掉呢?
        require(channel.closing_participant==closing_participant);

        require(non_closing_participant == recoverAddressFromBalanceProofUpdateMessage(
            channel_identifier,
            transferred_amount,
            locksroot,
            nonce,
            additional_hash,
            closing_signature,
            non_closing_signature
        ));
        require(closing_participant == recoverAddressFromBalanceProof(
            channel_identifier,
            transferred_amount,
            locksroot,
            nonce,
            additional_hash,
            closing_signature
        ));
        // Update the balance proof data for the closing_participant
        updateBalanceProofData(channel_identifier, closing_participant, nonce, locksroot, transferred_amount);
        emit NonClosingBalanceProofUpdated(channel_identifier, closing_participant);
    }
    /*
    任何人都可以调用,要在有效期内调用.通道状态必须是关闭,
    目前测试200个 lock 都是可以正常工作的.
    */
    function unlock(
        address participant,
        address partner,
        bytes merkle_tree_leaves
    ) 
    public
    {
        uint256 channel_identifier;
        require(merkle_tree_leaves.length > 0);
        channel_identifier=getChannelIdentifier(participant,partner);
        bytes32 computed_locksroot;
        uint256 unlocked_amount;
        Channel storage channel = channels[channel_identifier];
        require(channel.settle_block_number >= block.number);
        require(channel.state==2);
        Participant storage participant_state = channel.participants[participant];
        //如果为0,就没必要 unlock 了,
        require(participant_state.locksroot != 0);
        // Calculate the locksroot for the pending transfers and the amount of tokens
        // corresponding to the locked transfers with secrets revealed on chain.
        (computed_locksroot, unlocked_amount) = getMerkleRootAndUnlockedAmount(merkle_tree_leaves);
        require(computed_locksroot == participant_state.locksroot);

        /*
        因为每次都是设置unlocked amount, 所以也不用担心重放攻击,重放也没有多余的收益.

        考虑到惩罚机制能否生效的问题
        假设 A第一 次给 B 交易, B声明放弃此笔交易,但是回来看到 secret registry 有对应的密码就想到链上提交证据来获取收益.
        因为 unlock 几乎可以在有效期内反复调用,这就会导致惩罚的时候设置participant_state.unlocked_amount 没有任何意义.
        但是如果设置transfered_amount 为0,也没有任何惩罚意义,因为本来就是0,
        将deposit 设置为0意义也不大,因为有可能 deposit 本来就是0
        所以这笔钱注定要丢失?
        */
        participant_state.unlocked_amount = unlocked_amount;
        emit ChannelUnlocked(channel_identifier, participant,computed_locksroot, unlocked_amount);
    }
    /*

        /// @notice punish partner unlock a obsolete lock which he has annouced to abandon .
        // Anyone can call punishObsoleteUnlock  on behalf of a channel participant.
        /// @param channel_identifier The channel identifier - mapping key used for `channels`.
        /// @param beneficiary Address of the participant who owes the locked tokens.
        /// //@param expiration_block Block height at which the lock expires.
        /// @param locked_amount Amount of tokens that the locked transfer values.
        /// @param hashlock hash of a preimage used in a HTL Transfer
        /// @param additional_hash Computed from the message. Used for message authentication.
        /// @param signature signature of partner who has annouced to abandon this transfer,whether or not he knows the password.
        */

    function punishObsoleteUnlock(
        address beneficiary,
        address cheater,
        bytes32 lockhash,
        bytes32 additional_hash,
        bytes signature,
        bytes merkle_proof)
    public
    {
        uint256 channel_identifier;
        bytes32 locksroot;
        channel_identifier=getChannelIdentifier(beneficiary,cheater);
        Channel storage channel = channels[channel_identifier];
        require(channel.state==2);
        Participant storage beneficiary_state = channel.participants[beneficiary];

        // Check that the partner is a channel participant.
        // An empty locksroot means there are no pending locks
        require(beneficiary_state.locksroot != 0);
        /*
        the cheater provides his signature of lockhash to annouce that he has already abandon this transfer.
        */
        require(cheater == recoverAddressFromUnlockProof(
            channel_identifier,
            lockhash,
            additional_hash,
            signature
        ));
        Participant storage cheater_state = channel.participants[cheater];

        /*
        证明这个 lockhash 包含在受益方的 locksroot 中,既可以说明cheater 用了旧的证明,他声明放弃了的锁.
        */
        locksroot=computeMerkleRoot(lockhash,merkle_proof);
        require(beneficiary_state.locksroot==locksroot);
        /*
        punish the cheater.
        */
        beneficiary_state.transferred_amount = 0;
        beneficiary_state.unlocked_amount=0;
        beneficiary_state.deposit=beneficiary_state.deposit+cheater_state.deposit;
        cheater_state.deposit=0;
        
    }

    /*
    任何人都可以调用,只能调用一次
    */
    function settleChannel(
        address participant1,
        address participant2
    )
    public
    {
        uint256 participant1_amount;
        uint256 participant2_amount;
        uint256 participant1_transfer_amount;
        uint256 participant2_transfer_amount;
        uint256 total_deposit;
        uint256 channel_identifier;
        bytes32 channel_hash;
        channel_hash = getChannelHash(participant1, participant2);
        channel_identifier = openedchannels[channel_hash];
        Channel storage channel = channels[channel_identifier];
        // Channel must be closed
        require(channel.state == 2);

        // Settlement window must be over
        require(channel.settle_block_number < block.number);

        Participant storage participant1_state = channel.participants[participant1];
        Participant storage participant2_state = channel.participants[participant2];
        /*
        这两个 require 真有必要么? 只要 channel state 非0,肯定是 iniitilized 啊
        */
        total_deposit = participant1_state.deposit + participant2_state.deposit;
        participant1_transfer_amount=participant1_state.transferred_amount+participant1_state.unlocked_amount;
        participant2_transfer_amount=participant2_state.transferred_amount+participant2_state.unlocked_amount;

        participant1_amount = (
        participant1_state.deposit
        + participant2_transfer_amount
        - participant1_transfer_amount
        );
        // There are 2 cases that require attention here:
        // case1. If participant1 does NOT provide a balance proof or provides an old balance proof
        // case2. If participant2 does NOT provide a balance proof or provides an old balance proof
        // The issue is that we need to react differently in both cases. However, both cases have
        // an end result of participant1_amount > total_available_deposit. Therefore:

        // case1: participant2_transferred_amount can be [0, real_participant2_transferred_amount)
        // This can trigger an underflow -> participant1_amount > total_available_deposit
        // We need to make participant1_amount = 0 in this case, otherwise it can be
        // an attack vector. participant1 must lose all/some of its tokens if it does not
        // provide a valid balance proof.
        if (
            (participant1_state.deposit + participant2_transfer_amount) < participant1_transfer_amount
        ) {
            participant1_amount = 0;
        }

        // case2: participant1_transferred_amount can be [0, real_participant1_transferred_amount)
        // This means participant1_amount > total_available_deposit.
        // We need to limit participant1_amount to total_available_deposit. It is fine if
        // participant1 gets all the available tokens if participant2 has not provided a
        // valid balance proof.
        participant1_amount = min(participant1_amount, total_deposit);
        // At this point `participant1_amount` is between [0,total_deposit], so this is safe.
        participant2_amount = total_deposit - participant1_amount;
        // participant1_amount is the amount of tokens that participant1 will receive
        // participant2_amount is the amount of tokens that participant2 will receive
        // Remove the channel data from storage
        delete channel.participants[participant1];
        delete channel.participants[participant2];
        delete channels[channel_identifier];
        delete openedchannels[channel_hash];


        // Do the actual token transfers
        if (participant1_amount > 0) {
            require(token.transfer(participant1, participant1_amount));
        }

        if (participant2_amount > 0) {
            require(token.transfer(participant2, participant2_amount));
        }

        emit ChannelSettled(
            channel_identifier,
            participant1_amount,
            participant2_amount
        );
    }

    /*
    任何人都可以调用,只能调用一次.
    */
    function cooperativeSettle(
        address participant1_address,
        uint256 participant1_balance,
        address participant2_address,
        uint256 participant2_balance,
        bytes participant1_signature,
        bytes participant2_signature
    )
    public
    {
        address participant;
        uint256 total_available_deposit;
        bytes32 channel_hash;
        uint256 channel_identifier;
        channel_hash = getChannelHash(participant1_address, participant2_address);
        channel_identifier = openedchannels[channel_hash];
        Channel storage channel = channels[channel_identifier];

        participant = recoverAddressFromCooperativeSettleSignature(
            channel_identifier,
            participant1_address,
            participant1_balance,
            participant2_address,
            participant2_balance,
            participant1_signature
        );
        require(participant1_address == participant);
        participant = recoverAddressFromCooperativeSettleSignature(
            channel_identifier,
            participant1_address,
            participant1_balance,
            participant2_address,
            participant2_balance,
            participant2_signature
        );
        require(participant2_address == participant);

        Participant storage participant1_state = channel.participants[participant1_address];
        Participant storage participant2_state = channel.participants[participant2_address];
        // The channel must be open
        require(channel.state == 1);


        total_available_deposit = participant1_state.deposit + participant2_state.deposit;

        // Remove channel data from storage before doing the token transfers
        delete channel.participants[participant1_address];
        delete channel.participants[participant2_address];
        delete channels[channel_identifier];
        delete openedchannels[channel_hash];
        // Do the token transfers
        if (participant1_balance > 0) {
            require(token.transfer(participant1_address, participant1_balance));
        }

        if (participant2_balance > 0) {
            require(token.transfer(participant2_address, participant2_balance));
        }


        // The sum of the provided balances must be equal to the total available deposit
        require(total_available_deposit == (participant1_balance + participant2_balance));
        require(total_available_deposit >= participant1_balance);
        require(total_available_deposit >= participant2_balance);
        emit ChannelSettled(channel_identifier, participant1_balance, participant2_balance);
    }

    function getChannelIdentifier(address participant1,address participant2) view internal returns (uint256){
        bytes32 channel_hash;
        if (participant1 < participant2) {
            channel_hash= keccak256(abi.encodePacked(participant1, participant2));
        } else {
            channel_hash=keccak256(abi.encodePacked(participant2, participant1));
        }
        return openedchannels[channel_hash];
    }

    function getChannelHash(address participant, address partner)
    pure
    internal
    returns (bytes32)
    {
        // Lexicographic order of the channel addresses
        // This limits the number of channels that can be opened between two nodes to 1.
        if (participant < partner) {
            return keccak256(abi.encodePacked(participant, partner));
        } else {
            return keccak256(abi.encodePacked(partner, participant));
        }
    }


    function updateBalanceProofData(
        uint256 channel_identifier,
        address participant,
        uint256 nonce,
        bytes32 locksroot,
        uint256 transferred_amount
    )
    internal
    {
        Channel storage channel = channels[channel_identifier];
        Participant storage participant_state = channel.participants[participant];

        // Multiple calls to updateNonClosingBalanceProof can be made and we need to store
        // the last known balance proof data
        require(nonce > participant_state.nonce);
        participant_state.nonce = nonce;
        participant_state.locksroot = locksroot;
        participant_state.transferred_amount = transferred_amount;
    }

    function getChannelInfo(address participant1,address participant2)
    view
    external
    returns (uint256,uint256, address, uint8)
    {

        uint256 channel_identifier;
        bytes32 channel_hash;
        channel_hash=getChannelHash(participant1,participant2);
        channel_identifier=openedchannels[channel_hash]; 
        Channel storage channel = channels[channel_identifier];

        return (
        channel_identifier,
        channel.settle_block_number,
        channel.closing_participant,
        channel.state
        );
    }

    function getChannelParticipantInfo( address participant,address partner)
    view
    external
    returns (uint256, bytes32, uint256, uint256,uint256)
    {

        uint256 channel_identifier=getChannelIdentifier(participant,partner);
        Channel storage channel=channels[channel_identifier];
        Participant storage participant_state = channel.participants[participant];

        return (
        participant_state.deposit,
        participant_state.locksroot,
        participant_state.transferred_amount,
        participant_state.nonce,
        participant_state.unlocked_amount
        );
    }



    /*
     * Internal Functions
     */


    function recoverAddressFromBalanceProof(
        uint256 channel_identifier,
        uint256 transferred_amount,
        bytes32 locksroot,
        uint256 nonce,
        bytes32 additional_hash,
        bytes signature
    )
    view
    internal
    returns (address signature_address)
    {
        bytes32 message_hash = keccak256(abi.encodePacked(
                transferred_amount,
                locksroot,
                nonce,
                additional_hash,
                channel_identifier,
                address(this),
                chain_id
            ));

        signature_address = ECVerify.ecverify(message_hash, signature);
    }


    function recoverAddressFromBalanceProofUpdateMessage(
        uint256 channel_identifier,
        uint256 transferred_amount,
        bytes32 locksroot,
        uint256 nonce,
        bytes32 additional_hash,
        bytes closing_signature,
        bytes non_closing_signature
    )
    view
    internal
    returns (address signature_address)
    {
        bytes32 message_hash = keccak256(abi.encodePacked(
                transferred_amount,
                locksroot,
                nonce,
                additional_hash,
                channel_identifier,
                address(this),
                chain_id,
                closing_signature
            ));

        signature_address = ECVerify.ecverify(message_hash, non_closing_signature);
    }

    function recoverAddressFromCooperativeSettleSignature(
        uint256 channel_identifier,
        address participant1,
        uint256 participant1_balance,
        address participant2,
        uint256 participant2_balance,
        bytes signature
    )
    view
    internal
    returns (address signature_address)
    {
        bytes32 message_hash = keccak256(abi.encodePacked(
                participant1,
                participant1_balance,
                participant2,
                participant2_balance,
                channel_identifier,
                address(this),
                chain_id
            ));

        signature_address = ECVerify.ecverify(message_hash, signature);
    }


    function getMerkleRootAndUnlockedAmount(bytes merkle_tree_leaves)
    view
    internal
    returns (bytes32, uint256)
    {
        uint256 length = merkle_tree_leaves.length;

        // each merkle_tree lock component has this form:
        // (locked_amount || expiration_block || secrethash) = 3 * 32 bytes
        require(length % 96 == 0);

        uint256 i;
        uint256 total_unlocked_amount;
        uint256 unlocked_amount;
        bytes32 lockhash;
        bytes32 merkle_root;

        bytes32[] memory merkle_layer = new bytes32[](length / 96 + 1);

        for (i = 32; i < length; i += 96) {
            (lockhash, unlocked_amount) = getLockDataFromMerkleTree(merkle_tree_leaves, i);
            total_unlocked_amount += unlocked_amount;
            merkle_layer[i / 96] = lockhash;
        }

        length /= 96;

        while (length > 1) {
            if (length % 2 != 0) {
                merkle_layer[length] = merkle_layer[length - 1];
                length += 1;
            }

            for (i = 0; i < length - 1; i += 2) {
                if (merkle_layer[i] == merkle_layer[i + 1]) {
                    lockhash = merkle_layer[i];
                } else if (merkle_layer[i] < merkle_layer[i + 1]) {
                    lockhash = keccak256(abi.encodePacked(merkle_layer[i], merkle_layer[i + 1]));
                } else {
                    lockhash = keccak256(abi.encodePacked(merkle_layer[i + 1], merkle_layer[i]));
                }
                merkle_layer[i / 2] = lockhash;
            }
            length = i / 2;
        }

        merkle_root = merkle_layer[0];

        return (merkle_root, total_unlocked_amount);
    }

    function getLockDataFromMerkleTree(bytes merkle_tree_leaves, uint256 offset)
    view
    internal
    returns (bytes32, uint256)
    {
        uint256 expiration_block;
        uint256 locked_amount;
        uint256 reveal_block;
        bytes32 secrethash;
        bytes32 lockhash;

        if (merkle_tree_leaves.length <= offset) {
            return (lockhash, 0);
        }

        assembly {
            expiration_block := mload(add(merkle_tree_leaves, offset))
            locked_amount := mload(add(merkle_tree_leaves, add(offset, 32)))
            secrethash := mload(add(merkle_tree_leaves, add(offset, 64)))
        }

        // Calculate the lockhash for computing the merkle root
        lockhash = keccak256(abi.encodePacked(expiration_block, locked_amount, secrethash));

        // Check if the lock's secret was revealed in the SecretRegistry
        // The secret must have been revealed in the SecretRegistry contract before the lock's
        // expiration_block in order for the hash time lock transfer to be successful.
        reveal_block = secret_registry.getSecretRevealBlockHeight(secrethash);
        if (reveal_block == 0 || expiration_block <= reveal_block) {
            locked_amount = 0;
        }

        return (lockhash, locked_amount);
    }

    function computeMerkleRoot(bytes32 lockhash, bytes merkle_proof)
    pure
    internal
    returns (bytes32)
    {
        require(merkle_proof.length % 32 == 0);

        uint256 i;
        bytes32 el;

        for (i = 32; i <= merkle_proof.length; i += 32) {
            assembly {
                el := mload(add(merkle_proof, i))
            }

            if (lockhash < el) {
                lockhash = keccak256(abi.encodePacked(lockhash, el));
            } else {
                lockhash = keccak256(abi.encodePacked(el, lockhash));
            }
        }

        return lockhash;
    }

    function recoverAddressFromUnlockProof(
        uint256 channel_identifier,
        bytes32 lockhash,
        bytes32 additional_hash,
        bytes signature
    )
    view
    internal
    returns (address signature_address)
    {
        bytes32 message_hash = keccak256(abi.encodePacked(
                lockhash,
                channel_identifier,
                address(this),
                chain_id,
                additional_hash
            ));

        signature_address = ECVerify.ecverify(message_hash, signature);
    }

    function min(uint256 a, uint256 b) pure internal returns (uint256)
    {
        return a > b ? b : a;
    }

    function max(uint256 a, uint256 b) pure internal returns (uint256)
    {
        return a > b ? a : b;
    }

}