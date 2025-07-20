#!/bin/bash

# Targeted Pension to Suraksha Rename Script
# This script performs specific renames for the pension modules

echo "🔄 Starting targeted Pension to Suraksha rename..."

# Step 1: Rename directories
echo "📂 Renaming directories..."

# Rename grampension to gramsuraksha
if [ -d "x/grampension" ]; then
    echo "  Renaming x/grampension -> x/gramsuraksha"
    mv x/grampension x/gramsuraksha
fi

# Rename urbanpension to urbansuraksha
if [ -d "x/urbanpension" ]; then
    echo "  Renaming x/urbanpension -> x/urbansuraksha"
    mv x/urbanpension x/urbansuraksha
fi

# Rename proto directories
if [ -d "proto/deshchain/grampension" ]; then
    echo "  Renaming proto/deshchain/grampension -> proto/deshchain/gramsuraksha"
    mv proto/deshchain/grampension proto/deshchain/gramsuraksha
fi

# Rename Batua screens
if [ -d "batua/mobile/lib/ui/screens/pension" ]; then
    echo "  Renaming batua/.../pension -> batua/.../suraksha"
    mv batua/mobile/lib/ui/screens/pension batua/mobile/lib/ui/screens/suraksha
fi

# Step 2: Rename files
echo "📁 Renaming files..."

# Rename specific files
files_to_rename=(
    "x/gramsuraksha/types/pension_scheme.go:x/gramsuraksha/types/suraksha_scheme.go"
    "x/moneyorder/keeper/pension_liquidity.go:x/moneyorder/keeper/suraksha_liquidity.go"
    "x/moneyorder/keeper/pension_liquidity_test.go:x/moneyorder/keeper/suraksha_liquidity_test.go"
    "x/urbansuraksha/keeper/urban_pension_keeper.go:x/urbansuraksha/keeper/urban_suraksha_keeper.go"
    "x/urbansuraksha/types/urban_pension_scheme.go:x/urbansuraksha/types/urban_suraksha_scheme.go"
    "proto/deshchain/gramsuraksha/v1/grampension.proto:proto/deshchain/gramsuraksha/v1/gramsuraksha.proto"
    "batua/mobile/lib/ui/screens/suraksha/pension_scheme_screen.dart:batua/mobile/lib/ui/screens/suraksha/suraksha_scheme_screen.dart"
)

for rename in "${files_to_rename[@]}"; do
    IFS=':' read -r oldfile newfile <<< "$rename"
    if [ -f "$oldfile" ]; then
        echo "  Renaming: $oldfile -> $newfile"
        mv "$oldfile" "$newfile"
    fi
done

# Step 3: Update imports and references
echo "📝 Updating imports and references..."

# Update Go imports
echo "  Updating Go imports..."
find . -type f -name "*.go" ! -path "./.git/*" -exec sed -i.bak \
    -e 's|github.com/deshchain/deshchain/x/grampension|github.com/deshchain/deshchain/x/gramsuraksha|g' \
    -e 's|github.com/deshchain/deshchain/x/urbanpension|github.com/deshchain/deshchain/x/urbansuraksha|g' \
    {} \;

# Update proto imports
echo "  Updating proto imports..."
find . -type f -name "*.proto" ! -path "./.git/*" -exec sed -i.bak \
    -e 's|deshchain.grampension|deshchain.gramsuraksha|g' \
    -e 's|deshchain/grampension|deshchain/gramsuraksha|g' \
    {} \;

# Update type names and references
echo "  Updating type names..."
find . -type f \( -name "*.go" -o -name "*.proto" \) ! -path "./.git/*" -exec sed -i.bak \
    -e 's/PensionScheme/SurakshaScheme/g' \
    -e 's/pensionScheme/surakshaScheme/g' \
    -e 's/pension_scheme/suraksha_scheme/g' \
    -e 's/PensionParticipant/SurakshaParticipant/g' \
    -e 's/pensionParticipant/surakshaParticipant/g' \
    -e 's/pension_participant/suraksha_participant/g' \
    -e 's/PensionContribution/SurakshaContribution/g' \
    -e 's/pensionContribution/surakshaContribution/g' \
    -e 's/pension_contribution/suraksha_contribution/g' \
    -e 's/PensionMaturity/SurakshaMaturity/g' \
    -e 's/pensionMaturity/surakshaMaturity/g' \
    -e 's/pension_maturity/suraksha_maturity/g' \
    -e 's/PensionReserve/SurakshaReserve/g' \
    -e 's/pensionReserve/surakshaReserve/g' \
    -e 's/pension_reserve/suraksha_reserve/g' \
    -e 's/GetPension/GetSuraksha/g' \
    -e 's/SetPension/SetSuraksha/g' \
    -e 's/QueryPension/QuerySuraksha/g' \
    -e 's/MsgPension/MsgSuraksha/g' \
    {} \;

# Update Flutter/Dart files
echo "  Updating Flutter files..."
find batua -type f -name "*.dart" ! -path "./.git/*" -exec sed -i.bak \
    -e 's/PensionSchemeScreen/SurakshaSchemeScreen/g' \
    -e 's/pension_scheme_screen/suraksha_scheme_screen/g' \
    -e 's/Gram Pension Scheme/Gram Suraksha Pool/g' \
    -e 's/pension/suraksha/g' \
    {} \;

# Clean up backup files
echo "🧹 Cleaning up backup files..."
find . -name "*.bak" -delete

echo "✅ Targeted rename complete!"
echo ""
echo "⚠️  Please run the following commands to verify:"
echo "  1. go mod tidy"
echo "  2. make proto-gen (if you have proto generation)"
echo "  3. go test ./..."