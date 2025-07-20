/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import { BaseMessage } from './types';

// Placeholder exports for custom DeshChain messages
// These would be properly generated from protobuf definitions

export const MsgSend = BaseMessage;
export const MsgDelegate = BaseMessage;
export const MsgUndelegate = BaseMessage;
export const MsgWithdrawDelegatorReward = BaseMessage;

// NAMO token specific messages
export const MsgTransferNAMO = BaseMessage;
export const MsgBurnNAMO = BaseMessage;
export const MsgClaimVesting = BaseMessage;

// Money Order DEX messages
export const MsgCreateMoneyOrder = BaseMessage;
export const MsgAcceptMoneyOrder = BaseMessage;
export const MsgCancelMoneyOrder = BaseMessage;

// Sikkebaaz launchpad messages
export const MsgCreateTokenLaunch = BaseMessage;
export const MsgParticipateInLaunch = BaseMessage;
export const MsgClaimTokens = BaseMessage;
export const MsgInitiateCommunityVeto = BaseMessage;
export const MsgVoteOnVeto = BaseMessage;

// DhanSetu messages
export const MsgCreateDhanPata = BaseMessage;
export const MsgUpdateMitraProfile = BaseMessage;
export const MsgSendViaVirtualAddress = BaseMessage;

// Gram Suraksha messages
export const MsgEnrollInSuraksha = BaseMessage;
export const MsgContributeToPool = BaseMessage;
export const MsgClaimMaturity = BaseMessage;

// Cultural messages
export const MsgSubmitQuote = BaseMessage;
export const MsgVoteOnQuote = BaseMessage;
export const MsgClaimFestivalBonus = BaseMessage;