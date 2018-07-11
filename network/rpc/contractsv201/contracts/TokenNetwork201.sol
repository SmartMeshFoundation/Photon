pragma solidity ^0.4.23;

import "./Token.sol";
import "./Utils.sol";
import "./ECVerify.sol";
import "./SecretRegistry.sol";

/*
主要是在上一版基础之上
1. 运行 openanddeposit
2. updatebalanceproof 两个版本
3. event 提供更详细的nonce,transferamount,locksroot
*/
contract TokenNetwork is Utils {

    /*
     *  Data structures
     */

    string constant public contract_version = "0.3._";
    bytes32 constant public invalid_balance_hash=keccak256(uint64(0xffffffffffffffff),bytes32(0),uint256(0));

    // Instance of the token used as digital currency by the channels
    Token public token;
    // Instance of SecretRegistry used for storing secrets revealed in a mediating transfer.
    SecretRegistry public secret_registry;
    uint64 constant public punish_block_number=10;
    // Chain ID as specified by EIP155 used in balance proof signatures to avoid replay attacks
    uint256 public chain_id;
    // Channel identifier is a uint256, incremented after each new channel
    mapping(bytes32 => Channel) public channels;
    /*
    这个 locksroot 是否已经解锁过了,防止重放攻击
    */
    mapping(bytes32=>bool) unlocked_locksroot;
    // Start from 1 instead of 0, otherwise the first channel will have an additional
    // 15000 gas cost than the rest
    struct Participant {
        // Total amount of token transferred to this smart contract
        uint256 deposit;
        // The latest known merkle root of the pending hash-time locks, used to
        // validate the withdrawn proofs.
        /*
        nonce,locksroot,transferered_amount 的 hash
        */
        bytes32 balance_hash;
    }

    struct Channel {
        // After the channel has been uncooperatively closed, this value represents the
        // block number after which settleChannel can be called.
        uint64 settle_block_number;
        uint64 open_blocknumber;
        /*
        不关心是谁关闭了通道,关闭通道一方也可以再次更新证据,只要他能提供更新的 nonce 就可以了
        */
        //address closing_participant;
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
        bytes32 indexed channel_identifier,
        address  participant1,
        address  participant2,
        uint256 settle_timeout
    );

    event ChannelNewDeposit(
        bytes32 indexed channel_identifier,
        address  participant,
        uint256 total_deposit
    );

    event ChannelClosed(bytes32 indexed channel_identifier, address  closing_participant,uint256 nonce,bytes32 locksroot,uint256 transferred_amount);
    event ChannelUnlocked(
        bytes32 indexed channel_identifier,
        address payer_participant,
        uint256 nonce,
        bytes32 locskroot, //解锁的 locksroot
        uint256 transferred_amount
    );

    event BalanceProofUpdated(
        bytes32 indexed channel_identifier,
        address participant,
        uint256 nonce,
        bytes32 locksroot,
        uint256 transferred_amount
    );

    event ChannelSettled(
        bytes32 indexed channel_identifier,
        uint256 participant1_amount,
        uint256 participant2_amount
    );
    event Channelwithdraw(
        bytes32 channel_identifier,
        uint256 participant1_deposit,
        uint256 participant2_deposit,
        uint256 participant1_withdraw,
        uint256 participant2_withdraw
    );
    /*
 * Modifiers
 */

    modifier isOpen(bytes32 channel_identifier) {
        require(channels[channel_identifier].state == 1);
        _;
    }

    modifier isClosed(bytes32 channel_identifier) {
        require(channels[channel_identifier].state == 2);
        _;
    }

    modifier settleTimeoutValid(uint64 timeout) {
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


    function openChannel(address participant1, address participant2, uint64 settle_timeout)
    settleTimeoutValid(settle_timeout)
    public
    {
        bytes32 channel_identifier;
        require(participant1 != 0x0);
        require(participant2 != 0x0);
        require(participant1 != participant2);
        channel_identifier=getChannelIdentifier(participant1,participant2);
        Channel storage channel = channels[channel_identifier];
        /*
        保证channel没有被创建过
        */
        require(channel.state==0);
        /*
        可以考虑通过汇编指令将三句话合成一句,节省 gas, 说不定编译器已经这么做了.
        */
        // Store channel information
        channel.settle_block_number = settle_timeout;
        channel.open_blocknumber=uint64(block.number);
        // Mark channel as opened
        channel.state = 1;

        emit ChannelOpened(channel_identifier, participant1, participant2, settle_timeout);
    }
    function openChannelWithDeposit(address participant, address partner, uint64 settle_timeout,uint256 deposit)
    settleTimeoutValid(settle_timeout)
    public
    {
        bytes32 channel_identifier;
        require(participant != 0x0);
        require(partner != 0x0);
        require(participant != partner);
        channel_identifier=getChannelIdentifier(participant, partner);
        Channel storage channel = channels[channel_identifier];
        Participant storage participant_state = channel.participants[participant];
        /*
        保证channel没有被创建过
        */
        require(channel.state==0);

        // Store channel information
        channel.settle_block_number = settle_timeout;
        channel.open_blocknumber=uint64(block.number);
        // Mark channel as opened
        channel.state = 1;
        require(token.transferFrom(msg.sender, address(this), deposit));
        participant_state.deposit=deposit;
        emit ChannelOpened(channel_identifier, participant, partner, settle_timeout);
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
        uint256 current_deposit;
        bytes32 channel_identifier;
        channel_identifier=getChannelIdentifier(participant,partner);
        Channel storage channel = channels[channel_identifier];
        Participant storage participant_state = channel.participants[participant];
        current_deposit=participant_state.deposit;
        require(current_deposit < total_deposit);
        // Calculate the actual amount of tokens that will be transferred
        added_deposit = total_deposit - current_deposit;

        // Update the participant's channel deposit
        participant_state.deposit = current_deposit+added_deposit;


        // Do the transfer
        require(token.transferFrom(msg.sender, address(this), added_deposit));
        require(channel.state==1);
        emit ChannelNewDeposit(channel_identifier, participant, total_deposit);
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
        uint256 participant1_deposit,
        uint256 participant1_withdraw,
        address participant2,
        uint256 participant2_deposit,
        uint256 participant2_withdraw,
        bytes participant1_signature,
        bytes participant2_signature
    )
    public
    {
        uint256 total_deposit;
        bytes32 channel_identifier;
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
                channel.open_blocknumber,
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
                channel.open_blocknumber,
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
        participant1_state.deposit = participant1_deposit - participant1_withdraw;
        participant2_state.deposit = participant2_deposit - participant2_withdraw;
        //相当于 通道 settle 有新开了.老的签名都作废了.
        channel.open_blocknumber=uint64(block.number);

        emit Channelwithdraw(channel_identifier, participant1_deposit, participant2_deposit, participant1_withdraw, participant2_withdraw);

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
        bytes32 channel_identifier;
        address recovered_partner_address;
        channel_identifier=getChannelIdentifier(msg.sender,partner);
        Channel storage channel = channels[channel_identifier];
        require(channel.state==1);
        // Mark the channel as closed and mark the closing participant
        channel.state = 2;
        // This is the block number at which the channel can be settled.
        channel.settle_block_number += uint64(block.number);
        // Nonce 0 means that the closer never received a transfer, therefore never received a
        // balance proof, or he is intentionally not providing the latest transfer, in which case
        // the closing party is going to lose the tokens that were transferred to him.
        if (nonce > 0) {
            Participant storage partner_state=channel.participants[partner];
            recovered_partner_address = recoverAddressFromBalanceProof(
                channel_identifier,
                transferred_amount,
                locksroot,
                nonce,
                channel.open_blocknumber,
                additional_hash,
                signature
            );
            require(partner==recovered_partner_address);
            partner_state.balance_hash=calceBalanceHash(transferred_amount,locksroot,nonce);
        }
        emit ChannelClosed(channel_identifier, msg.sender,nonce,locksroot,transferred_amount);
    }
    /*
    任何人都可以调用,可以调用多次,只要在有效期内.
    包括 closing 方和非 close 方都可以反复调用在,只要能够提供更新的 nonce 即可.
    */
    function updateBalanceProofDelegate(
        address participant,
        address partner,
        uint256 transferred_amount,
        bytes32 locksroot,
        uint256 nonce,
        uint256 old_transferred_amount,
        bytes32 old_locksroot,
        uint256 old_nonce,
        bytes32 additional_hash,
        bytes participant_signature,
        bytes partner_signature
    )
    public
    {
        bytes32 channel_identifier;
        channel_identifier=getChannelIdentifier(participant, partner);
        Channel storage channel = channels[channel_identifier];
        require(channel.state==2);
        require(channel.settle_block_number >= block.number);
        require(nonce > 0);


        require(partner == recoverAddressFromBalanceProofUpdateMessage(
            channel_identifier,
            transferred_amount,
            locksroot,
            nonce,
            channel.open_blocknumber,
            additional_hash,
            participant_signature,
            partner_signature
        ));
        require(participant == recoverAddressFromBalanceProof(
            channel_identifier,
            transferred_amount,
            locksroot,
            nonce,
            channel.open_blocknumber,
            additional_hash,
            participant_signature
        ));
        // Update the balance proof data for the closing_participant
        verifyBalanceHashIsValid( channel_identifier, participant, old_transferred_amount,old_locksroot,old_nonce,nonce);
        updateBalanceHash(channel_identifier,participant,transferred_amount,locksroot,nonce);
        //todo 如何修复 stack too deep 的问题呢?最后一个参数 transferred_amount 多余,但是又很有用.
        emit BalanceProofUpdated(channel_identifier, participant,nonce,locksroot,transferred_amount);
    }
    function updateBalanceProof(
        address participant,
        uint256 transferred_amount,
        bytes32 locksroot,
        uint256 nonce,
        uint256 old_transferred_amount,
        bytes32 old_locksroot,
        uint256 old_nonce,
        bytes32 additional_hash,
        bytes participant_signature
    )
    public
    {
        bytes32 channel_identifier;
        channel_identifier=getChannelIdentifier(participant, msg.sender);
        Channel storage channel = channels[channel_identifier];
        require(channel.state==2);
        require(channel.settle_block_number >= block.number);
        require(nonce > 0);

        require(participant == recoverAddressFromBalanceProof(
            channel_identifier,
            transferred_amount,
            locksroot,
            nonce,
            channel.open_blocknumber,
            additional_hash,
            participant_signature
        ));
        // Update the balance proof data for the closing_participant
        verifyBalanceHashIsValid( channel_identifier, participant, old_transferred_amount,old_locksroot,old_nonce,nonce);
        updateBalanceHash(channel_identifier,participant,transferred_amount,locksroot,nonce);
        //todo 如何修复 stack too deep 的问题呢?最后一个参数 transferred_amount 多余,但是又很有用.
        emit BalanceProofUpdated(channel_identifier, participant,nonce,locksroot,transferred_amount);
    }
    /*
    任何人都可以调用,要在有效期内调用.通道状态必须是关闭,
    目前测试200个 lock 都是可以正常工作的.
    */
    function unlock(
        address participant,
        address partner,
        uint256 transferered_amount,
        bytes32 locksroot,
        uint256 nonce,
        bytes merkle_tree_leaves
    )
    public
    {
        bytes32 channel_identifier;
        bytes32 locksroot_hash;
        bytes32 computed_locksroot;
        uint256 unlocked_amount;
        bytes32 balance_hash;
        require(merkle_tree_leaves.length > 0);
        channel_identifier=getChannelIdentifier(participant,partner);
        Channel storage channel = channels[channel_identifier];
        Participant storage participant_state = channel.participants[participant];
        require(channel.settle_block_number >= block.number);
        require(channel.state==2);
        balance_hash=participant_state.balance_hash;
        //如果为0,就没必要 unlock 了,
        require(locksroot != 0);
        //最好加上 channel_identifier,否则比如 A-B-C 交易, C 如果 unlock 的话, B 就无法 unlock 了.
        locksroot_hash=keccak256(balance_hash,locksroot,channel_identifier);
        //这个 locksroot 还没有 unlock 过,不允许反复 unlock, 重放攻击
        require(unlocked_locksroot[locksroot_hash]==false);
        unlocked_locksroot[locksroot_hash]=true;
        // Calculate the locksroot for the pending transfers and the amount of tokens
        // corresponding to the locked transfers with secrets revealed on chain.
        (computed_locksroot, unlocked_amount) = getMerkleRootAndUnlockedAmount(merkle_tree_leaves);
        require(unlocked_amount>0); //有必要么?
        require(computed_locksroot == locksroot);
        require(balance_hash==calceBalanceHash(transferered_amount,locksroot,nonce));
        /*
        会不会溢出呢? 两人持续交易?
        正常来说,不会,
        但是如果是恶意的会溢出,但是溢出对于 partner 也没好处啊
        */
        transferered_amount += unlocked_amount;
        /*
       注意transferered_amount已经更新了,
        */
        participant_state.balance_hash=calceBalanceHash(transferered_amount,locksroot,nonce);
        emit ChannelUnlocked(channel_identifier, participant,nonce,computed_locksroot, transferered_amount);
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

    /*
    给 punish 一方留出了专门的 punishBlock 时间,punish 一方可以选择在那个时候提交证据,也可以在这之前.
    如果能够提供 old beneficiary_transferred_amount,也就是加 unlock 之前的,可以将unlocked_locksroot中的删除,从而再节省 gas, 不过意义好像并不大.
    */
    function punishObsoleteUnlock(
        address beneficiary,
        address cheater,
        bytes32 lockhash,
        uint256 beneficiary_transferred_amount,
        uint256 beneficiary_nonce,
        bytes32 additional_hash,
        bytes signature,
        bytes merkle_proof)
    public
    {
        bytes32 channel_identifier;
        bytes32 locksroot;
        bytes32 balance_hash;
        channel_identifier=getChannelIdentifier(beneficiary,cheater);
        Channel storage channel = channels[channel_identifier];
        require(channel.state==2);
        Participant storage beneficiary_state = channel.participants[beneficiary];
        balance_hash=beneficiary_state.balance_hash;

        // Check that the partner is a channel participant.
        // An empty locksroot means there are no pending locks
        require(balance_hash != 0);
        /*
        the cheater provides his signature of lockhash to annouce that he has already abandon this transfer.
        需要加上 chainid,tokennetworkAddress不?有其他重放的可能不?
        */
        require(cheater == recoverAddressFromUnlockProof(
            channel_identifier,
            lockhash,
            uint64(channel.open_blocknumber),
            additional_hash,
            signature
        ));
        Participant storage cheater_state = channel.participants[cheater];

        /*
        证明这个 lockhash 包含在受益方的 locksroot 中,既可以说明cheater 用了旧的证明,他声明放弃了的锁.
        */
        locksroot=computeMerkleRoot(lockhash,merkle_proof);
        require(balance_hash==calceBalanceHash(beneficiary_transferred_amount,locksroot,beneficiary_nonce));
        /*
        punish the cheater.
        */
        beneficiary_state.balance_hash = invalid_balance_hash;
        beneficiary_state.deposit=beneficiary_state.deposit+cheater_state.deposit;
        cheater_state.deposit=0;
    }

    /*
    任何人都可以调用,只能调用一次
    */
    function settleChannel(
        address participant1,
        uint256 participant1_transferred_amount,
        bytes32 participant1_locksroot,
        uint256 participant1_nonce,
        address participant2,
        uint256 participant2_transferred_amount,
        bytes32 participant2_locksroot,
        uint256 participant2_nonce
    )
    public
    {
        uint256 participant1_amount;
        uint256 total_deposit;
        bytes32 channel_identifier;
        channel_identifier=getChannelIdentifier(participant1,participant2);
        Channel storage channel = channels[channel_identifier];
        // Channel must be closed
        require(channel.state == 2);

        // Settlement window must be over
        require(channel.settle_block_number+punish_block_number < block.number); //300是给 punish 用的,

        Participant storage participant1_state = channel.participants[participant1];
        Participant storage participant2_state = channel.participants[participant2];
        /*
        验证提供的参数是有效的
        */
        require(participant1_state.balance_hash==calceBalanceHash(participant1_transferred_amount,participant1_locksroot,participant1_nonce));
        require(participant2_state.balance_hash==calceBalanceHash(participant2_transferred_amount,participant2_locksroot,participant2_nonce));

        total_deposit = participant1_state.deposit + participant2_state.deposit;

        participant1_amount = (
        participant1_state.deposit
        + participant2_transferred_amount
        - participant1_transferred_amount
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
            (participant1_state.deposit + participant2_transferred_amount) < participant1_transferred_amount
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
        //变量复用是因为局部变量不能超过16个
        participant2_transferred_amount = total_deposit - participant1_amount;
        // participant1_amount is the amount of tokens that participant1 will receive
        // participant2_amount is the amount of tokens that participant2 will receive
        // Remove the channel data from storage
        delete channel.participants[participant1];
        delete channel.participants[participant2];
        delete channels[channel_identifier];
        // Do the actual token transfers
        if (participant1_amount > 0) {
            require(token.transfer(participant1, participant1_amount));
        }

        if (participant2_transferred_amount > 0) {
            require(token.transfer(participant2, participant2_transferred_amount));
        }

        emit ChannelSettled(
            channel_identifier,
            participant1_amount,
            participant2_transferred_amount
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
        bytes32 channel_identifier;
        uint64 open_blocknumber;
        channel_identifier = getChannelIdentifier(participant1_address, participant2_address);
        Channel storage channel = channels[channel_identifier];
        // The channel must be open
        require(channel.state == 1);

        open_blocknumber=channel.open_blocknumber;
        participant = recoverAddressFromCooperativeSettleSignature(
            channel_identifier,
            participant1_address,
            participant1_balance,
            participant2_address,
            participant2_balance,
            open_blocknumber,
            participant1_signature
        );
        require(participant1_address == participant);
        participant = recoverAddressFromCooperativeSettleSignature(
            channel_identifier,
            participant1_address,
            participant1_balance,
            participant2_address,
            participant2_balance,
            open_blocknumber,
            participant2_signature
        );
        require(participant2_address == participant);

        Participant storage participant1_state = channel.participants[participant1_address];
        Participant storage participant2_state = channel.participants[participant2_address];



        total_available_deposit = participant1_state.deposit + participant2_state.deposit;

        // Remove channel data from storage before doing the token transfers
        delete channel.participants[participant1_address];
        delete channel.participants[participant2_address];
        delete channels[channel_identifier];
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

    function getChannelIdentifier(address participant1,address participant2) pure internal returns (bytes32){
        if (participant1 < participant2) {
            return keccak256(abi.encodePacked(participant1, participant2));
        } else {
            return keccak256(abi.encodePacked(participant2, participant1));
        }
    }


    function verifyBalanceHashIsValid(
        bytes32 channel_identifier,
        address participant,
        uint256 transferred_amount,
        bytes32 locksroot,
        uint256 nonce,
        uint256 new_nonce

    )
    view internal
    {
        Channel storage channel = channels[channel_identifier];
        Participant storage participant_state = channel.participants[participant];
        bytes32 balance_hash=calceBalanceHash(transferred_amount,locksroot,nonce);
        require(participant_state.balance_hash==balance_hash);
        require(new_nonce>nonce);
    }
    function updateBalanceHash(
        bytes32 channel_identifier,
        address participant,
        uint256 transferred_amount,
        bytes32 locksroot,
        uint256 nonce)  internal
    {
        Channel storage channel = channels[channel_identifier];
        Participant storage participant_state = channel.participants[participant];
        bytes32 balance_hash=calceBalanceHash(transferred_amount,locksroot,nonce);
        participant_state.balance_hash=balance_hash;
    }
    function calceBalanceHash(uint256 transferred_amount,bytes32 locksroot,uint256 nonce ) pure internal returns(bytes32){
        if( nonce==0 && locksroot==0 && transferred_amount==0){
            return 0;
        }
        return keccak256(nonce,locksroot,transferred_amount);
    }

    function getChannelInfo(address participant1,address participant2)
    view
    external
    returns (bytes32,uint64,uint64 , uint8)
    {

        bytes32 channel_identifier;
        channel_identifier=getChannelIdentifier(participant1,participant2);
        Channel storage channel = channels[channel_identifier];

        return (
        channel_identifier,
        channel.settle_block_number,
        channel.open_blocknumber,
        channel.state
        );
    }

    function getChannelParticipantInfo( address participant,address partner)
    view
    external
    returns (uint256, bytes32)
    {

        bytes32 channel_identifier=getChannelIdentifier(participant,partner);
        Channel storage channel=channels[channel_identifier];
        Participant storage participant_state = channel.participants[participant];

        return (
        participant_state.deposit,
        participant_state.balance_hash
        );
    }



    /*
     * Internal Functions
     */


    function recoverAddressFromBalanceProof(
        bytes32 channel_identifier,
        uint256 transferred_amount,
        bytes32 locksroot,
        uint256 nonce,
        uint64 open_blocknumber,
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
                open_blocknumber,
                address(this),
                chain_id
            ));

        signature_address = ECVerify.ecverify(message_hash, signature);
    }


    function recoverAddressFromBalanceProofUpdateMessage(
        bytes32 channel_identifier,
        uint256 transferred_amount,
        bytes32 locksroot,
        uint256 nonce,
        uint64 open_blocknumber,
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
                open_blocknumber,
                address(this),
                chain_id,
                closing_signature
            ));

        signature_address = ECVerify.ecverify(message_hash, non_closing_signature);
    }

    function recoverAddressFromCooperativeSettleSignature(
        bytes32 channel_identifier,
        address participant1,
        uint256 participant1_balance,
        address participant2,
        uint256 participant2_balance,
        uint64 open_blocknumber,
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
                open_blocknumber,
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
        bytes32 channel_identifier,
        bytes32 lockhash,
        uint64 open_blocknumber,
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
                open_blocknumber,
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