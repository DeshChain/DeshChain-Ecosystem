const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const rateLimit = require('express-rate-limit');
const NodeCache = require('node-cache');
const winston = require('winston');
const { DirectSecp256k1HdWallet } = require('@cosmjs/proto-signing');
const { SigningStargateClient, GasPrice } = require('@cosmjs/stargate');

// Initialize logger
const logger = winston.createLogger({
  level: 'info',
  format: winston.format.json(),
  transports: [
    new winston.transports.Console({
      format: winston.format.simple(),
    }),
  ],
});

// Initialize Express app
const app = express();

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json());

// Initialize cache (TTL in seconds)
const cache = new NodeCache({ stdTTL: parseInt(process.env.COOLDOWN_TIME || '3600') });

// Rate limiting
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 5, // Limit each IP to 5 requests per windowMs
  message: 'Too many requests from this IP, please try again later.',
});

app.use('/api/faucet', limiter);

// Faucet configuration
const config = {
  rpcUrl: process.env.RPC_URL || 'http://localhost:26657',
  chainId: process.env.CHAIN_ID || 'deshchain-testnet-1',
  mnemonic: process.env.FAUCET_MNEMONIC || 'witness effort dose make crucial vote nature glove observe dilemma alpha invite lady wage fall shaft stock melody birth check refuse emotion fiscal cruise',
  dripAmount: process.env.DRIP_AMOUNT || '1000000000unamo',
  cooldownTime: parseInt(process.env.COOLDOWN_TIME || '3600'),
  prefix: 'desh',
};

// Initialize wallet and client
let wallet;
let client;
let faucetAddress;

async function initializeFaucet() {
  try {
    // Create wallet from mnemonic
    wallet = await DirectSecp256k1HdWallet.fromMnemonic(config.mnemonic, {
      prefix: config.prefix,
    });

    // Get faucet address
    const [firstAccount] = await wallet.getAccounts();
    faucetAddress = firstAccount.address;
    logger.info(`Faucet address: ${faucetAddress}`);

    // Connect to blockchain
    client = await SigningStargateClient.connectWithSigner(
      config.rpcUrl,
      wallet,
      {
        gasPrice: GasPrice.fromString('0.025unamo'),
      }
    );

    logger.info('Connected to DeshChain');

    // Check faucet balance
    const balance = await client.getAllBalances(faucetAddress);
    logger.info(`Faucet balance: ${JSON.stringify(balance)}`);
  } catch (error) {
    logger.error('Failed to initialize faucet:', error);
    setTimeout(initializeFaucet, 5000);
  }
}

// API Routes
app.get('/api/faucet/info', async (req, res) => {
  try {
    const balance = await client.getAllBalances(faucetAddress);
    const chainId = await client.getChainId();
    const height = await client.getHeight();

    res.json({
      faucetAddress,
      balance,
      chainId,
      height,
      dripAmount: config.dripAmount,
      cooldownTime: config.cooldownTime,
    });
  } catch (error) {
    res.status(500).json({ error: 'Failed to get faucet info' });
  }
});

app.post('/api/faucet/request', async (req, res) => {
  try {
    const { address } = req.body;

    // Validate address
    if (!address || !address.startsWith(config.prefix)) {
      return res.status(400).json({ error: 'Invalid address' });
    }

    // Check cooldown
    const lastRequest = cache.get(address);
    if (lastRequest) {
      const remainingTime = config.cooldownTime - Math.floor((Date.now() - lastRequest) / 1000);
      return res.status(429).json({
        error: 'Address is in cooldown period',
        remainingTime,
      });
    }

    // Check IP cooldown
    const ip = req.ip || req.connection.remoteAddress;
    const ipLastRequest = cache.get(`ip:${ip}`);
    if (ipLastRequest) {
      const remainingTime = config.cooldownTime - Math.floor((Date.now() - ipLastRequest) / 1000);
      return res.status(429).json({
        error: 'IP address is in cooldown period',
        remainingTime,
      });
    }

    // Parse drip amount
    const [amount, denom] = config.dripAmount.match(/(\d+)(\w+)/).slice(1);

    // Send tokens
    const result = await client.sendTokens(
      faucetAddress,
      address,
      [{ amount, denom }],
      {
        amount: [{ amount: '5000', denom: 'unamo' }],
        gas: '200000',
      },
      'DeshChain Testnet Faucet'
    );

    // Set cooldown
    cache.set(address, Date.now());
    cache.set(`ip:${ip}`, Date.now());

    // Log transaction
    logger.info(`Sent ${config.dripAmount} to ${address}. Tx: ${result.transactionHash}`);

    res.json({
      success: true,
      transactionHash: result.transactionHash,
      amount: config.dripAmount,
      recipient: address,
    });
  } catch (error) {
    logger.error('Faucet request error:', error);
    res.status(500).json({ error: 'Failed to process faucet request' });
  }
});

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

// Homepage
app.get('/', (req, res) => {
  res.send(`
    <!DOCTYPE html>
    <html>
    <head>
      <title>DeshChain Testnet Faucet</title>
      <style>
        body {
          font-family: Arial, sans-serif;
          max-width: 600px;
          margin: 50px auto;
          padding: 20px;
          background: linear-gradient(135deg, #FF9933 0%, #FFFFFF 50%, #138808 100%);
          min-height: 100vh;
        }
        .container {
          background: white;
          padding: 30px;
          border-radius: 10px;
          box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        h1 {
          color: #FF9933;
          text-align: center;
        }
        input, button {
          width: 100%;
          padding: 10px;
          margin: 10px 0;
          border: 1px solid #ddd;
          border-radius: 5px;
        }
        button {
          background: #138808;
          color: white;
          cursor: pointer;
          font-weight: bold;
        }
        button:hover {
          background: #0f6606;
        }
        .message {
          padding: 10px;
          margin: 10px 0;
          border-radius: 5px;
          text-align: center;
        }
        .success {
          background: #d4edda;
          color: #155724;
          border: 1px solid #c3e6cb;
        }
        .error {
          background: #f8d7da;
          color: #721c24;
          border: 1px solid #f5c6cb;
        }
        .info {
          background: #d1ecf1;
          color: #0c5460;
          border: 1px solid #bee5eb;
        }
      </style>
    </head>
    <body>
      <div class="container">
        <h1>🚰 DeshChain Testnet Faucet</h1>
        <p>Get free testnet NAMO tokens for testing DeshChain features.</p>
        
        <input type="text" id="address" placeholder="Enter your desh1... address" />
        <button onclick="requestTokens()">Request Tokens</button>
        
        <div id="message"></div>
        
        <div class="info">
          <p><strong>Drip Amount:</strong> ${config.dripAmount}</p>
          <p><strong>Cooldown:</strong> ${config.cooldownTime} seconds</p>
        </div>
      </div>
      
      <script>
        async function requestTokens() {
          const address = document.getElementById('address').value;
          const messageDiv = document.getElementById('message');
          
          if (!address || !address.startsWith('desh')) {
            messageDiv.className = 'message error';
            messageDiv.innerHTML = 'Please enter a valid desh1... address';
            return;
          }
          
          try {
            const response = await fetch('/api/faucet/request', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ address })
            });
            
            const data = await response.json();
            
            if (response.ok) {
              messageDiv.className = 'message success';
              messageDiv.innerHTML = \`✅ Success! Sent \${data.amount} to \${data.recipient}<br>
                <a href="http://localhost:3000/tx/\${data.transactionHash}" target="_blank">View Transaction</a>\`;
            } else {
              messageDiv.className = 'message error';
              messageDiv.innerHTML = \`❌ Error: \${data.error}\`;
              if (data.remainingTime) {
                messageDiv.innerHTML += \`<br>Please wait \${data.remainingTime} seconds\`;
              }
            }
          } catch (error) {
            messageDiv.className = 'message error';
            messageDiv.innerHTML = '❌ Network error. Please try again.';
          }
        }
      </script>
    </body>
    </html>
  `);
});

// Start server
const PORT = process.env.PORT || 4000;

async function start() {
  await initializeFaucet();
  
  app.listen(PORT, () => {
    logger.info(`Faucet service running on port ${PORT}`);
  });
}

start().catch((error) => {
  logger.error('Failed to start:', error);
  process.exit(1);
});