/**
 * Basic Usage Example for DeshChain Identity SDK
 * 
 * This example demonstrates:
 * - Creating an identity
 * - Issuing credentials
 * - Verifying credentials
 * - Selective disclosure
 * - Zero-knowledge proofs
 */

import { 
  createIdentitySigningClient,
  DeshChainIdentitySDK,
  VerifiableCredential
} from '../index';

// Configuration
const RPC_ENDPOINT = process.env.DESHCHAIN_RPC || 'https://rpc.testnet.deshchain.com';
const CHAIN_ID = process.env.DESHCHAIN_CHAIN_ID || 'deshchain-testnet-1';
const MNEMONIC = process.env.DESHCHAIN_MNEMONIC || 'your test mnemonic here';

async function main() {
  console.log('🚀 DeshChain Identity SDK - Basic Usage Example\n');

  try {
    // 1. Initialize SDK
    console.log('1️⃣ Initializing SDK...');
    const { sdk, address } = await createIdentitySigningClient({
      rpcEndpoint: RPC_ENDPOINT,
      chainId: CHAIN_ID
    }, MNEMONIC);
    console.log(`✅ Connected with address: ${address}\n`);

    // 2. Create Identity
    console.log('2️⃣ Creating identity...');
    const publicKey = 'base64-encoded-public-key'; // In production, generate proper key
    const did = await sdk.createIdentity(address, publicKey, {
      name: 'Test User',
      type: 'individual',
      country: 'IN'
    });
    console.log(`✅ Created DID: ${did}\n`);

    // 3. Issue KYC Credential
    console.log('3️⃣ Issuing KYC credential...');
    const kycCredentialId = await sdk.issueCredential(
      address, // Issuer address
      did,     // Holder DID
      ['KYCCredential', 'IdentityCredential'],
      {
        fullName: 'Test User',
        dateOfBirth: '1990-01-01',
        nationality: 'Indian',
        kycLevel: 'standard',
        verifiedAt: new Date().toISOString(),
        riskRating: 'low',
        address: {
          street: '123 Main Street',
          city: 'Mumbai',
          state: 'Maharashtra',
          country: 'IN',
          postalCode: '400001'
        }
      },
      new Date(Date.now() + 180 * 24 * 60 * 60 * 1000) // 6 months expiry
    );
    console.log(`✅ Issued KYC credential: ${kycCredentialId}\n`);

    // 4. Issue Education Credential
    console.log('4️⃣ Issuing education credential...');
    const eduCredentialId = await sdk.issueCredential(
      address,
      did,
      ['EducationCredential', 'DegreeCredential'],
      {
        degree: 'Bachelor of Technology',
        field: 'Computer Science',
        university: 'Indian Institute of Technology',
        graduationYear: 2012,
        grade: 'First Class with Distinction',
        rollNumber: 'CS/2008/001'
      }
    );
    console.log(`✅ Issued education credential: ${eduCredentialId}\n`);

    // 5. Register Biometric
    console.log('5️⃣ Registering biometric credential...');
    const biometricCredId = await sdk.registerBiometric(address, {
      type: 'FINGERPRINT',
      templateHash: 'sha256-hash-of-biometric-template',
      deviceId: 'test-device-001'
    });
    console.log(`✅ Registered biometric: ${biometricCredId}\n`);

    // 6. Query Credentials
    console.log('6️⃣ Querying credentials...');
    const credentials = await sdk.getCredentialsBySubject(did);
    console.log(`✅ Found ${credentials.length} credentials for ${did}`);
    credentials.forEach((cred, index) => {
      console.log(`   ${index + 1}. ${cred.type.join(', ')} - ID: ${cred.id}`);
    });
    console.log();

    // 7. Verify Credential
    console.log('7️⃣ Verifying KYC credential...');
    const kycCredential = credentials.find(c => c.type.includes('KYCCredential'));
    if (kycCredential) {
      const isValid = await sdk.verifyCredential(kycCredential);
      console.log(`✅ KYC credential is ${isValid ? 'VALID' : 'INVALID'}\n`);
    }

    // 8. Selective Disclosure
    console.log('8️⃣ Presenting credentials with selective disclosure...');
    const verifierDID = 'did:desh:verifier123';
    const presentationId = await sdk.presentCredential(
      address,
      [kycCredentialId, eduCredentialId],
      verifierDID,
      {
        [kycCredentialId]: ['fullName', 'kycLevel', 'verifiedAt'],
        [eduCredentialId]: ['degree', 'university']
      }
    );
    console.log(`✅ Created presentation: ${presentationId}`);
    console.log('   Revealed only: name, KYC level, verification date, degree, and university\n');

    // 9. Zero-Knowledge Proof - Age Verification
    console.log('9️⃣ Creating zero-knowledge proof for age verification...');
    const ageProofId = await sdk.createZKProof(address, {
      type: 'age-range',
      statement: 'age >= 18',
      credentials: [kycCredentialId],
      options: {
        anonymitySet: 10,
        expiryMinutes: 60
      }
    });
    console.log(`✅ Created age proof: ${ageProofId}`);
    console.log('   Proves: User is 18 or older without revealing actual date of birth\n');

    // 10. Biometric Authentication
    console.log('🔟 Testing biometric authentication...');
    const isAuthenticated = await sdk.authenticateBiometric(address, {
      type: 'FINGERPRINT',
      templateHash: 'sha256-hash-of-biometric-template',
      deviceId: 'test-device-001'
    });
    console.log(`✅ Biometric authentication: ${isAuthenticated ? 'SUCCESS' : 'FAILED'}\n`);

    // 11. Update Privacy Settings
    console.log('1️⃣1️⃣ Updating privacy settings...');
    await sdk.updatePrivacySettings(address, {
      disclosureLevel: 'minimal',
      requireConsent: true,
      allowAnonymous: true
    });
    console.log('✅ Updated privacy settings to minimal disclosure with consent requirement\n');

    // Summary
    console.log('📊 Summary:');
    console.log(`   - Created identity: ${did}`);
    console.log(`   - Issued ${credentials.length} credentials`);
    console.log('   - Demonstrated selective disclosure');
    console.log('   - Created zero-knowledge proof');
    console.log('   - Registered and tested biometric authentication');
    console.log('   - Updated privacy settings');
    console.log('\n✨ Basic usage example completed successfully!');

  } catch (error) {
    console.error('❌ Error:', error);
    process.exit(1);
  }
}

// Helper function to demonstrate credential structure
function printCredential(credential: VerifiableCredential) {
  console.log('\n📄 Credential Structure:');
  console.log(JSON.stringify({
    context: credential.context,
    id: credential.id,
    type: credential.type,
    issuer: credential.issuer,
    issuanceDate: credential.issuanceDate,
    expirationDate: credential.expirationDate,
    credentialSubject: credential.credentialSubject,
    proof: credential.proof ? {
      type: credential.proof.type,
      created: credential.proof.created,
      verificationMethod: credential.proof.verificationMethod
    } : undefined
  }, null, 2));
}

// Run the example
main().catch(console.error);