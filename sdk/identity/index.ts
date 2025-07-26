/**
 * DeshChain Identity SDK
 * 
 * A comprehensive SDK for interacting with DeshChain's Identity Module
 * Supporting W3C DIDs, Verifiable Credentials, and Zero-Knowledge Proofs
 */

import { SigningStargateClient, StargateClient } from '@cosmjs/stargate';
import { EncodeObject, Registry } from '@cosmjs/proto-signing';
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing';
import { Tendermint34Client } from '@cosmjs/tendermint-rpc';

// Identity types
export interface DID {
  did: string;
  controller: string;
  verificationMethod: VerificationMethod[];
  authentication: string[];
  service?: Service[];
  created: Date;
  updated: Date;
}

export interface VerificationMethod {
  id: string;
  type: string;
  controller: string;
  publicKeyMultibase?: string;
  blockchainAccountId?: string;
}

export interface Service {
  id: string;
  type: string;
  serviceEndpoint: string;
}

export interface VerifiableCredential {
  context: string[];
  id: string;
  type: string[];
  issuer: string;
  issuanceDate: Date;
  expirationDate?: Date;
  credentialSubject: any;
  proof?: Proof;
  status?: CredentialStatus;
}

export interface Proof {
  type: string;
  created: Date;
  verificationMethod: string;
  proofPurpose: string;
  proofValue: string;
}

export interface CredentialStatus {
  type: string;
  reason?: string;
}

export interface IdentityOptions {
  rpcEndpoint: string;
  chainId: string;
  prefix?: string;
}

export interface BiometricCredential {
  type: 'FINGERPRINT' | 'FACE' | 'IRIS' | 'VOICE' | 'PALM';
  templateHash: string;
  deviceId: string;
  score?: number;
}

export interface ZKProofRequest {
  type: string;
  statement: string;
  credentials: string[];
  options?: ZKProofOptions;
}

export interface ZKProofOptions {
  anonymitySet?: number;
  expiryMinutes?: number;
}

export class DeshChainIdentitySDK {
  private client: StargateClient | null = null;
  private signingClient: SigningStargateClient | null = null;
  private options: IdentityOptions;
  private registry: Registry;

  constructor(options: IdentityOptions) {
    this.options = {
      prefix: 'deshchain',
      ...options
    };
    this.registry = this.createRegistry();
  }

  /**
   * Initialize the SDK with a connection to the blockchain
   */
  async connect(): Promise<void> {
    const tendermintClient = await Tendermint34Client.connect(this.options.rpcEndpoint);
    this.client = await StargateClient.create(tendermintClient);
  }

  /**
   * Initialize signing client with wallet
   */
  async connectWithSigner(mnemonic: string): Promise<string> {
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
      prefix: this.options.prefix
    });
    
    const [account] = await wallet.getAccounts();
    const tendermintClient = await Tendermint34Client.connect(this.options.rpcEndpoint);
    
    this.signingClient = await SigningStargateClient.createWithSigner(
      tendermintClient,
      wallet,
      { registry: this.registry }
    );

    return account.address;
  }

  /**
   * Create a new identity with DID
   */
  async createIdentity(
    address: string,
    publicKey: string,
    metadata?: Record<string, string>
  ): Promise<string> {
    if (!this.signingClient) {
      throw new Error('Signing client not initialized. Call connectWithSigner first.');
    }

    const msg: EncodeObject = {
      typeUrl: '/deshchain.identity.MsgCreateIdentity',
      value: {
        creator: address,
        publicKey,
        metadata
      }
    };

    const result = await this.signingClient.signAndBroadcast(
      address,
      [msg],
      'auto',
      'Create identity'
    );

    if (result.code !== 0) {
      throw new Error(`Failed to create identity: ${result.rawLog}`);
    }

    // Extract DID from result
    const did = this.extractDIDFromResult(result);
    return did;
  }

  /**
   * Register a DID document
   */
  async registerDID(
    address: string,
    didDocument: Partial<DID>
  ): Promise<void> {
    if (!this.signingClient) {
      throw new Error('Signing client not initialized');
    }

    const msg: EncodeObject = {
      typeUrl: '/deshchain.identity.MsgRegisterDID',
      value: {
        creator: address,
        didDocument
      }
    };

    const result = await this.signingClient.signAndBroadcast(
      address,
      [msg],
      'auto',
      'Register DID'
    );

    if (result.code !== 0) {
      throw new Error(`Failed to register DID: ${result.rawLog}`);
    }
  }

  /**
   * Issue a verifiable credential
   */
  async issueCredential(
    issuerAddress: string,
    holderDID: string,
    credentialType: string[],
    credentialSubject: any,
    expirationDate?: Date
  ): Promise<string> {
    if (!this.signingClient) {
      throw new Error('Signing client not initialized');
    }

    const credential: Partial<VerifiableCredential> = {
      context: [
        'https://www.w3.org/2018/credentials/v1',
        'https://deshchain.com/contexts/v1'
      ],
      type: ['VerifiableCredential', ...credentialType],
      issuer: `did:desh:${issuerAddress}`,
      issuanceDate: new Date(),
      expirationDate,
      credentialSubject: {
        id: holderDID,
        ...credentialSubject
      }
    };

    const msg: EncodeObject = {
      typeUrl: '/deshchain.identity.MsgIssueCredential',
      value: {
        issuer: issuerAddress,
        holder: holderDID,
        credential
      }
    };

    const result = await this.signingClient.signAndBroadcast(
      issuerAddress,
      [msg],
      'auto',
      'Issue credential'
    );

    if (result.code !== 0) {
      throw new Error(`Failed to issue credential: ${result.rawLog}`);
    }

    return this.extractCredentialIDFromResult(result);
  }

  /**
   * Verify a credential
   */
  async verifyCredential(credential: VerifiableCredential): Promise<boolean> {
    if (!this.client) {
      throw new Error('Client not initialized');
    }

    try {
      const response = await this.client.queryContractSmart(
        this.getIdentityModuleAddress(),
        {
          verify_credential: {
            credential
          }
        }
      );

      return response.valid === true;
    } catch (error) {
      console.error('Credential verification failed:', error);
      return false;
    }
  }

  /**
   * Present credentials with selective disclosure
   */
  async presentCredential(
    holderAddress: string,
    credentialIds: string[],
    verifierDID: string,
    revealedClaims: Record<string, string[]>
  ): Promise<string> {
    if (!this.signingClient) {
      throw new Error('Signing client not initialized');
    }

    const msg: EncodeObject = {
      typeUrl: '/deshchain.identity.MsgPresentCredential',
      value: {
        holder: holderAddress,
        credentialIds,
        verifier: verifierDID,
        revealedClaims
      }
    };

    const result = await this.signingClient.signAndBroadcast(
      holderAddress,
      [msg],
      'auto',
      'Present credential'
    );

    if (result.code !== 0) {
      throw new Error(`Failed to present credential: ${result.rawLog}`);
    }

    return this.extractPresentationIDFromResult(result);
  }

  /**
   * Create a zero-knowledge proof
   */
  async createZKProof(
    proverAddress: string,
    proofRequest: ZKProofRequest
  ): Promise<string> {
    if (!this.signingClient) {
      throw new Error('Signing client not initialized');
    }

    const msg: EncodeObject = {
      typeUrl: '/deshchain.identity.MsgCreateZKProof',
      value: {
        prover: proverAddress,
        proofType: proofRequest.type,
        statement: proofRequest.statement,
        credentials: proofRequest.credentials,
        options: proofRequest.options
      }
    };

    const result = await this.signingClient.signAndBroadcast(
      proverAddress,
      [msg],
      'auto',
      'Create ZK proof'
    );

    if (result.code !== 0) {
      throw new Error(`Failed to create ZK proof: ${result.rawLog}`);
    }

    return this.extractProofIDFromResult(result);
  }

  /**
   * Register biometric credential
   */
  async registerBiometric(
    userAddress: string,
    biometric: BiometricCredential
  ): Promise<string> {
    if (!this.signingClient) {
      throw new Error('Signing client not initialized');
    }

    const credentialSubject = {
      biometricType: biometric.type,
      templateHash: biometric.templateHash,
      deviceId: biometric.deviceId
    };

    return this.issueCredential(
      userAddress,
      `did:desh:${userAddress}`,
      ['BiometricCredential'],
      credentialSubject,
      new Date(Date.now() + 2 * 365 * 24 * 60 * 60 * 1000) // 2 years
    );
  }

  /**
   * Authenticate with biometric
   */
  async authenticateBiometric(
    userAddress: string,
    biometric: BiometricCredential
  ): Promise<boolean> {
    if (!this.client) {
      throw new Error('Client not initialized');
    }

    const response = await this.client.queryContractSmart(
      this.getIdentityModuleAddress(),
      {
        authenticate_biometric: {
          user: userAddress,
          biometric_type: biometric.type,
          template_hash: biometric.templateHash,
          device_id: biometric.deviceId
        }
      }
    );

    return response.success === true && (response.score || 0) >= 0.85;
  }

  /**
   * Link Aadhaar with consent
   */
  async linkAadhaar(
    userAddress: string,
    aadhaarHash: string,
    consentArtefact: string
  ): Promise<void> {
    if (!this.signingClient) {
      throw new Error('Signing client not initialized');
    }

    const msg: EncodeObject = {
      typeUrl: '/deshchain.identity.MsgLinkAadhaar',
      value: {
        user: userAddress,
        aadhaarHash,
        consentArtefact
      }
    };

    const result = await this.signingClient.signAndBroadcast(
      userAddress,
      [msg],
      'auto',
      'Link Aadhaar'
    );

    if (result.code !== 0) {
      throw new Error(`Failed to link Aadhaar: ${result.rawLog}`);
    }
  }

  /**
   * Query identity by DID
   */
  async getIdentity(did: string): Promise<DID | null> {
    if (!this.client) {
      throw new Error('Client not initialized');
    }

    try {
      const response = await this.client.queryContractSmart(
        this.getIdentityModuleAddress(),
        {
          get_identity: { did }
        }
      );

      return response as DID;
    } catch (error) {
      console.error('Failed to get identity:', error);
      return null;
    }
  }

  /**
   * Query credentials by subject
   */
  async getCredentialsBySubject(subjectDID: string): Promise<VerifiableCredential[]> {
    if (!this.client) {
      throw new Error('Client not initialized');
    }

    try {
      const response = await this.client.queryContractSmart(
        this.getIdentityModuleAddress(),
        {
          get_credentials_by_subject: { subject: subjectDID }
        }
      );

      return response.credentials || [];
    } catch (error) {
      console.error('Failed to get credentials:', error);
      return [];
    }
  }

  /**
   * Update privacy settings
   */
  async updatePrivacySettings(
    userAddress: string,
    settings: {
      disclosureLevel: 'minimal' | 'standard' | 'full';
      requireConsent: boolean;
      allowAnonymous: boolean;
    }
  ): Promise<void> {
    if (!this.signingClient) {
      throw new Error('Signing client not initialized');
    }

    const msg: EncodeObject = {
      typeUrl: '/deshchain.identity.MsgUpdatePrivacySettings',
      value: {
        user: userAddress,
        settings
      }
    };

    const result = await this.signingClient.signAndBroadcast(
      userAddress,
      [msg],
      'auto',
      'Update privacy settings'
    );

    if (result.code !== 0) {
      throw new Error(`Failed to update privacy settings: ${result.rawLog}`);
    }
  }

  // Helper methods
  
  private createRegistry(): Registry {
    // Create custom registry with identity module types
    const registry = new Registry();
    
    // Register identity module types
    // This would include all the message types from the identity module
    
    return registry;
  }

  private getIdentityModuleAddress(): string {
    // Return the identity module address based on chain
    return 'deshchain1identity...'; // Placeholder
  }

  private extractDIDFromResult(result: any): string {
    // Extract DID from transaction result
    const events = result.events || [];
    for (const event of events) {
      if (event.type === 'identity_created') {
        const didAttr = event.attributes.find((attr: any) => attr.key === 'did');
        if (didAttr) {
          return didAttr.value;
        }
      }
    }
    throw new Error('DID not found in result');
  }

  private extractCredentialIDFromResult(result: any): string {
    // Extract credential ID from result
    const events = result.events || [];
    for (const event of events) {
      if (event.type === 'credential_issued') {
        const idAttr = event.attributes.find((attr: any) => attr.key === 'credential_id');
        if (idAttr) {
          return idAttr.value;
        }
      }
    }
    throw new Error('Credential ID not found in result');
  }

  private extractPresentationIDFromResult(result: any): string {
    // Extract presentation ID from result
    const events = result.events || [];
    for (const event of events) {
      if (event.type === 'credential_presented') {
        const idAttr = event.attributes.find((attr: any) => attr.key === 'presentation_id');
        if (idAttr) {
          return idAttr.value;
        }
      }
    }
    throw new Error('Presentation ID not found in result');
  }

  private extractProofIDFromResult(result: any): string {
    // Extract proof ID from result
    const events = result.events || [];
    for (const event of events) {
      if (event.type === 'zk_proof_created') {
        const idAttr = event.attributes.find((attr: any) => attr.key === 'proof_id');
        if (idAttr) {
          return idAttr.value;
        }
      }
    }
    throw new Error('Proof ID not found in result');
  }
}

// Export convenience functions

export async function createIdentityClient(options: IdentityOptions): Promise<DeshChainIdentitySDK> {
  const sdk = new DeshChainIdentitySDK(options);
  await sdk.connect();
  return sdk;
}

export async function createIdentitySigningClient(
  options: IdentityOptions,
  mnemonic: string
): Promise<{ sdk: DeshChainIdentitySDK; address: string }> {
  const sdk = new DeshChainIdentitySDK(options);
  const address = await sdk.connectWithSigner(mnemonic);
  return { sdk, address };
}