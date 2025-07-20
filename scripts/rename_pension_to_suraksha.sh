#!/bin/bash

# DeshChain Pension to Suraksha Rename Script
# This script renames all instances of "pension" to "suraksha" for regulatory compliance

echo "🔄 Starting Pension to Suraksha rename process..."

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Backup function
backup_file() {
    if [ -f "$1" ]; then
        cp "$1" "$1.backup_$(date +%Y%m%d_%H%M%S)"
    fi
}

# Function to rename files
rename_files() {
    echo -e "${YELLOW}📁 Renaming files containing 'pension'...${NC}"
    
    # Find all files with pension in the name
    find . -type f -name "*pension*" ! -path "./scripts/*" ! -path "./.git/*" | while read file; do
        newfile=$(echo "$file" | sed 's/pension/suraksha/g' | sed 's/Pension/Suraksha/g')
        if [ "$file" != "$newfile" ]; then
            echo "  Renaming: $file -> $newfile"
            mkdir -p "$(dirname "$newfile")"
            mv "$file" "$newfile"
        fi
    done
}

# Function to update file contents
update_file_contents() {
    echo -e "${YELLOW}📝 Updating file contents...${NC}"
    
    # Define file patterns to update
    patterns=(
        "*.go"
        "*.proto"
        "*.md"
        "*.yaml"
        "*.yml"
        "*.json"
        "*.ts"
        "*.tsx"
        "*.js"
        "*.jsx"
    )
    
    for pattern in "${patterns[@]}"; do
        echo "  Processing $pattern files..."
        find . -type f -name "$pattern" ! -path "./scripts/*" ! -path "./.git/*" ! -path "./node_modules/*" | while read file; do
            if grep -q -i "pension" "$file"; then
                backup_file "$file"
                
                # Perform case-sensitive replacements
                sed -i.tmp 's/GramPension/GramSuraksha/g' "$file"
                sed -i.tmp 's/grampension/gramsuraksha/g' "$file"
                sed -i.tmp 's/Pension/Suraksha/g' "$file"
                sed -i.tmp 's/pension/suraksha/g' "$file"
                sed -i.tmp 's/PENSION/SURAKSHA/g' "$file"
                
                # Clean up temp files
                rm -f "$file.tmp"
                
                echo "    ✓ Updated: $file"
            fi
        done
    done
}

# Function to update module paths
update_module_paths() {
    echo -e "${YELLOW}🔧 Updating module paths...${NC}"
    
    # Update Go imports
    find . -type f -name "*.go" ! -path "./.git/*" | while read file; do
        if grep -q "github.com/deshchain/deshchain/x/grampension" "$file"; then
            sed -i 's|github.com/deshchain/deshchain/x/grampension|github.com/deshchain/deshchain/x/gramsuraksha|g' "$file"
            echo "    ✓ Updated imports in: $file"
        fi
    done
}

# Function to rename directories
rename_directories() {
    echo -e "${YELLOW}📂 Renaming directories...${NC}"
    
    # Rename the main module directory
    if [ -d "x/grampension" ]; then
        echo "  Renaming x/grampension -> x/gramsuraksha"
        mv x/grampension x/gramsuraksha
    fi
    
    # Rename any other pension directories
    find . -type d -name "*pension*" ! -path "./.git/*" | sort -r | while read dir; do
        newdir=$(echo "$dir" | sed 's/pension/suraksha/g' | sed 's/Pension/Suraksha/g')
        if [ "$dir" != "$newdir" ] && [ -d "$dir" ]; then
            echo "  Renaming: $dir -> $newdir"
            mv "$dir" "$newdir"
        fi
    done
}

# Function to update proto files specifically
update_proto_files() {
    echo -e "${YELLOW}🔗 Updating protobuf definitions...${NC}"
    
    # Update proto package names
    find . -type f -name "*.proto" ! -path "./.git/*" | while read file; do
        if grep -q "grampension" "$file"; then
            sed -i 's/package deshchain.grampension/package deshchain.gramsuraksha/g' "$file"
            sed -i 's/option go_package = ".*grampension.*"/option go_package = "github.com\/deshchain\/deshchain\/x\/gramsuraksha\/types"/g' "$file"
            echo "    ✓ Updated proto: $file"
        fi
    done
}

# Function to update specific terminology
update_terminology() {
    echo -e "${YELLOW}💬 Updating specific terminology...${NC}"
    
    # Update specific pension-related terms
    find . -type f \( -name "*.go" -o -name "*.md" -o -name "*.proto" \) ! -path "./.git/*" | while read file; do
        # Backup before terminology updates
        backup_file "$file"
        
        # Update specific terms
        sed -i 's/pension scheme/suraksha pool/gi' "$file"
        sed -i 's/pension plan/suraksha plan/gi' "$file"
        sed -i 's/pensioner/suraksha member/gi' "$file"
        sed -i 's/pension fund/suraksha fund/gi' "$file"
        sed -i 's/pension benefits/suraksha benefits/gi' "$file"
        sed -i 's/retirement pension/retirement suraksha/gi' "$file"
        
        # Update guarantee language for compliance
        sed -i 's/guaranteed returns/target returns/gi' "$file"
        sed -i 's/guaranteed 50%/target 50%/gi' "$file"
    done
}

# Function to generate summary report
generate_report() {
    echo -e "${YELLOW}📊 Generating summary report...${NC}"
    
    report_file="pension_to_suraksha_report_$(date +%Y%m%d_%H%M%S).txt"
    
    {
        echo "Pension to Suraksha Rename Report"
        echo "================================="
        echo "Date: $(date)"
        echo ""
        echo "Files renamed:"
        find . -name "*.backup_*" -type f | wc -l
        echo ""
        echo "Backup files created:"
        find . -name "*.backup_*" -type f
        echo ""
        echo "Remaining 'pension' references (please review manually):"
        grep -r -i "pension" . --exclude-dir=.git --exclude-dir=scripts --exclude="*.backup_*" | head -20
    } > "$report_file"
    
    echo -e "${GREEN}✅ Report saved to: $report_file${NC}"
}

# Main execution
main() {
    echo "This script will rename all 'pension' references to 'suraksha'"
    echo "Backup files will be created for all modified files"
    echo ""
    read -p "Continue? (y/n): " -n 1 -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Execute rename operations
        rename_directories
        rename_files
        update_file_contents
        update_module_paths
        update_proto_files
        update_terminology
        
        # Generate report
        generate_report
        
        echo -e "${GREEN}✅ Pension to Suraksha rename completed!${NC}"
        echo -e "${YELLOW}⚠️  Please review the report and test the changes${NC}"
    else
        echo -e "${RED}❌ Operation cancelled${NC}"
    fi
}

# Run main function
main