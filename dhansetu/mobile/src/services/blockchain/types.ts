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

import { GeneratedType } from '@cosmjs/proto-signing';

// Create a base message class that implements GeneratedType
export class BaseMessage implements GeneratedType {
  static readonly typeUrl: string;
  static readonly encode: (message: any) => Uint8Array;
  static readonly decode: (reader: any) => any;
  static readonly fromJSON: (object: any) => any;
  static readonly fromPartial: (object: any) => any;
  static readonly toJSON: (message: any) => unknown;
}

// Export the base type
export type DeshChainMessage = typeof BaseMessage;