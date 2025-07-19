#!/bin/bash

# Script to add Apache 2.0 license headers to source files
# Copyright 2024 DeshChain Foundation

LICENSE_HEADER_GO='/*
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

'

LICENSE_HEADER_JS='/*
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

'

LICENSE_HEADER_PY='# Copyright 2024 DeshChain Foundation
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

'

LICENSE_HEADER_SH='#!/bin/bash
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

'

# Function to check if file already has license header
has_license_header() {
    local file="$1"
    head -5 "$file" | grep -q "Copyright.*DeshChain Foundation"
}

# Function to add license header to Go files
add_go_license() {
    local file="$1"
    if ! has_license_header "$file"; then
        echo "Adding license header to: $file"
        # Create temporary file with license header and original content
        {
            echo -n "$LICENSE_HEADER_GO"
            cat "$file"
        } > "$file.tmp"
        mv "$file.tmp" "$file"
    else
        echo "License header already exists in: $file"
    fi
}

# Function to add license header to JS/TS files
add_js_license() {
    local file="$1"
    if ! has_license_header "$file"; then
        echo "Adding license header to: $file"
        {
            echo -n "$LICENSE_HEADER_JS"
            cat "$file"
        } > "$file.tmp"
        mv "$file.tmp" "$file"
    else
        echo "License header already exists in: $file"
    fi
}

# Function to add license header to Python files
add_py_license() {
    local file="$1"
    if ! has_license_header "$file"; then
        echo "Adding license header to: $file"
        # Handle shebang line if present
        if head -1 "$file" | grep -q "^#!"; then
            {
                head -1 "$file"
                echo -n "$LICENSE_HEADER_PY"
                tail -n +2 "$file"
            } > "$file.tmp"
        else
            {
                echo -n "$LICENSE_HEADER_PY"
                cat "$file"
            } > "$file.tmp"
        fi
        mv "$file.tmp" "$file"
    else
        echo "License header already exists in: $file"
    fi
}

# Function to add license header to shell scripts
add_sh_license() {
    local file="$1"
    if ! has_license_header "$file"; then
        echo "Adding license header to: $file"
        # Replace existing shebang with our licensed version
        if head -1 "$file" | grep -q "^#!"; then
            {
                echo -n "$LICENSE_HEADER_SH"
                tail -n +2 "$file"
            } > "$file.tmp"
        else
            {
                echo -n "$LICENSE_HEADER_SH"
                cat "$file"
            } > "$file.tmp"
        fi
        mv "$file.tmp" "$file"
        chmod +x "$file"
    else
        echo "License header already exists in: $file"
    fi
}

# Main execution
echo "Adding Apache 2.0 license headers to source files..."

# Process Go files
echo "Processing Go files..."
find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" | while read -r file; do
    add_go_license "$file"
done

# Process JavaScript and TypeScript files
echo "Processing JavaScript/TypeScript files..."
find . \( -name "*.js" -o -name "*.ts" -o -name "*.jsx" -o -name "*.tsx" \) \
    -not -path "./node_modules/*" -not -path "./vendor/*" -not -path "./.git/*" | while read -r file; do
    add_js_license "$file"
done

# Process Python files
echo "Processing Python files..."
find . -name "*.py" -not -path "./vendor/*" -not -path "./.git/*" | while read -r file; do
    add_py_license "$file"
done

# Process Shell scripts
echo "Processing Shell scripts..."
find . -name "*.sh" -not -path "./vendor/*" -not -path "./.git/*" | while read -r file; do
    add_sh_license "$file"
done

echo "License header addition completed!"