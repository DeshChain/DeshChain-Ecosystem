#!/bin/bash
# Copyright 2024 DeshChain Foundation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

echo "🔍 DeshChain License Validation Tool"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
TOTAL_FILES=0
MISSING_LICENSE=0
ERRORS=0

# Check for Apache license in Go files
echo "📋 Checking Go files for Apache 2.0 license headers..."
GO_FILES=$(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" | wc -l)
GO_MISSING=$(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -exec sh -c 'head -5 "$1" | grep -q "Copyright.*DeshChain Foundation" || echo "$1"' _ {} \; | wc -l)

if [ $GO_MISSING -eq 0 ]; then
    echo -e "${GREEN}✅ All $GO_FILES Go files have proper Apache 2.0 headers${NC}"
else
    echo -e "${RED}❌ $GO_MISSING Go files missing Apache license headers${NC}"
    MISSING_LICENSE=$((MISSING_LICENSE + GO_MISSING))
fi
echo ""

# Check for Apache license in JavaScript/TypeScript files
echo "📋 Checking JavaScript/TypeScript files for Apache 2.0 license headers..."
JS_FILES=$(find . \( -name "*.js" -o -name "*.ts" -o -name "*.jsx" -o -name "*.tsx" \) -not -path "./node_modules/*" -not -path "./vendor/*" -not -path "./.git/*" | wc -l)

if [ $JS_FILES -gt 0 ]; then
    JS_MISSING=$(find . \( -name "*.js" -o -name "*.ts" -o -name "*.jsx" -o -name "*.tsx" \) -not -path "./node_modules/*" -not -path "./vendor/*" -not -path "./.git/*" -exec sh -c 'head -5 "$1" | grep -q "Copyright.*DeshChain Foundation" || echo "$1"' _ {} \; | wc -l)
    
    if [ $JS_MISSING -eq 0 ]; then
        echo -e "${GREEN}✅ All $JS_FILES JavaScript/TypeScript files have proper Apache 2.0 headers${NC}"
    else
        echo -e "${RED}❌ $JS_MISSING JavaScript/TypeScript files missing Apache license headers${NC}"
        MISSING_LICENSE=$((MISSING_LICENSE + JS_MISSING))
    fi
else
    echo -e "${YELLOW}ℹ️  No JavaScript/TypeScript files found${NC}"
fi
echo ""

# Check for Apache license in Python files
echo "📋 Checking Python files for Apache 2.0 license headers..."
PY_FILES=$(find . -name "*.py" -not -path "./vendor/*" -not -path "./.git/*" | wc -l)

if [ $PY_FILES -gt 0 ]; then
    PY_MISSING=$(find . -name "*.py" -not -path "./vendor/*" -not -path "./.git/*" -exec sh -c 'head -5 "$1" | grep -q "Copyright.*DeshChain Foundation" || echo "$1"' _ {} \; | wc -l)
    
    if [ $PY_MISSING -eq 0 ]; then
        echo -e "${GREEN}✅ All $PY_FILES Python files have proper Apache 2.0 headers${NC}"
    else
        echo -e "${RED}❌ $PY_MISSING Python files missing Apache license headers${NC}"
        MISSING_LICENSE=$((MISSING_LICENSE + PY_MISSING))
    fi
else
    echo -e "${YELLOW}ℹ️  No Python files found${NC}"
fi
echo ""

# Check for Apache license in Shell scripts
echo "📋 Checking Shell scripts for Apache 2.0 license headers..."
SH_FILES=$(find . -name "*.sh" -not -path "./vendor/*" -not -path "./.git/*" | wc -l)

if [ $SH_FILES -gt 0 ]; then
    SH_MISSING=$(find . -name "*.sh" -not -path "./vendor/*" -not -path "./.git/*" -exec sh -c 'head -5 "$1" | grep -q "Copyright.*DeshChain Foundation" || echo "$1"' _ {} \; | wc -l)
    
    if [ $SH_MISSING -eq 0 ]; then
        echo -e "${GREEN}✅ All $SH_FILES Shell scripts have proper Apache 2.0 headers${NC}"
    else
        echo -e "${RED}❌ $SH_MISSING Shell scripts missing Apache license headers${NC}"
        MISSING_LICENSE=$((MISSING_LICENSE + SH_MISSING))
    fi
else
    echo -e "${YELLOW}ℹ️  No Shell scripts found${NC}"
fi
echo ""

# Check required license files
echo "📄 Checking required license files..."
REQUIRED_FILES=("LICENSE" "LICENSE-CULTURAL" "NOTICE" "cultural-data/LICENSE")

for file in "${REQUIRED_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✅ Found: $file${NC}"
    else
        echo -e "${RED}❌ Missing: $file${NC}"
        ERRORS=$((ERRORS + 1))
    fi
done
echo ""

# Check cultural data in Go source files
echo "🎭 Checking cultural data in source code..."
CULTURAL_GO_FILE="x/cultural/types/cultural_data.go"

if [ -f "$CULTURAL_GO_FILE" ]; then
    if head -20 "$CULTURAL_GO_FILE" | grep -q "Copyright.*DeshChain Foundation"; then
        echo -e "${GREEN}✅ Cultural data in Go source has Apache 2.0 license${NC}"
        echo -e "${GREEN}✅ Cultural content embedded in source code (quotes, events, wisdom)${NC}"
    else
        echo -e "${RED}❌ Cultural data file missing license header${NC}"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo -e "${RED}❌ Cultural data file not found: $CULTURAL_GO_FILE${NC}"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# Calculate totals
TOTAL_FILES=$((GO_FILES + JS_FILES + PY_FILES + SH_FILES))
TOTAL_ERRORS=$((ERRORS + MISSING_LICENSE))

# Summary
echo "📊 Summary:"
echo "============================"
echo "Total source files checked: $TOTAL_FILES"
echo "Files missing license headers: $MISSING_LICENSE"
echo "Other errors: $ERRORS"
echo "Total errors: $TOTAL_ERRORS"
echo ""

if [ $TOTAL_ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ All license checks passed!${NC}"
    echo "DeshChain maintains proper dual licensing compliance."
    echo ""
    echo "📋 License Structure:"
    echo "- Source Code: Apache 2.0"
    echo "- Cultural Data: Embedded in source with Apache 2.0"
    echo "- External Cultural Content: CC BY-NC-SA 4.0"
    exit 0
else
    echo -e "${RED}❌ License validation failed!${NC}"
    echo "Run 'scripts/add-license-headers.sh' to automatically add missing headers."
    exit 1
fi