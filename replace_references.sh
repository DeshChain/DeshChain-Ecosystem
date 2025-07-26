#!/bin/bash

# Replace deshchain.bharat with deshchain.com
find . -type f \( -name "*.md" -o -name "*.go" -o -name "*.ts" -o -name "*.js" -o -name "*.html" -o -name "*.json" \) \
  -not -path "./venv/*" \
  -not -path "./node_modules/*" \
  -not -path "./.git/*" \
  -exec sed -i 's/deshchain\.bharat/deshchain.com/g' {} +

echo "Replaced all instances of deshchain.bharat with deshchain.com"

# List of files to check for founder references
echo "Files that need manual review for founder references:"
echo "- README.md (NAMO token tribute section)"
echo "- x/cultural/types/cultural_data.go (Historical references)"
echo "- x/nft/types/pradhan_sevak.go (Complete NFT tribute file)"
echo "- x/nft/genesis_nft.go (NFT minting references)"