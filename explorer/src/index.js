const express = require('express');
const cors = require('cors');
const app = express();
const port = process.env.PORT || 3001;

app.use(cors());
app.use(express.json());

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// Mock blockchain data endpoints
app.get('/api/blocks/latest', (req, res) => {
  res.json({
    height: 1000,
    hash: '0x' + Math.random().toString(16).substr(2, 64),
    time: new Date().toISOString(),
    txs: Math.floor(Math.random() * 10)
  });
});

app.get('/api/validators', (req, res) => {
  const validators = [];
  for (let i = 0; i < 21; i++) {
    validators.push({
      address: '0x' + Math.random().toString(16).substr(2, 40),
      moniker: `Validator-${i}`,
      voting_power: Math.floor(Math.random() * 1000000),
      commission: '0.10'
    });
  }
  res.json(validators);
});

app.get('/api/transactions/recent', (req, res) => {
  const txs = [];
  for (let i = 0; i < 10; i++) {
    txs.push({
      hash: '0x' + Math.random().toString(16).substr(2, 64),
      height: 1000 - i,
      type: 'transfer',
      amount: Math.floor(Math.random() * 1000) + ' NAMO',
      fee: '0.001 NAMO',
      status: 'success'
    });
  }
  res.json(txs);
});

app.listen(port, () => {
  console.log(`Explorer backend listening at http://localhost:${port}`);
});