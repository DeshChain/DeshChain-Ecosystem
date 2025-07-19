#!/usr/bin/env node
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

const fs = require('fs');
const path = require('path');
const glob = require('glob');

// License headers to check
const APACHE_LICENSE_PATTERN = /Copyright \d{4} DeshChain Foundation[\s\S]*Apache License, Version 2\.0/;
const CC_LICENSE_PATTERN = /Creative Commons Attribution-NonCommercial-ShareAlike 4\.0/;

// File patterns to check
const SOURCE_PATTERNS = [
  '**/*.go',
  '**/*.js',
  '**/*.ts',
  '**/*.jsx',
  '**/*.tsx',
  '**/*.py',
  '**/*.sh',
  '**/Makefile*',
  '**/Dockerfile*'
];

const CULTURAL_PATTERNS = [
  'cultural-data/**/*.json',
  'cultural-data/**/*.md'
];

const EXCLUDE_PATTERNS = [
  '**/node_modules/**',
  '**/vendor/**',
  '**/.git/**',
  '**/build/**',
  '**/dist/**',
  '**/*.min.js',
  '**/*.min.css'
];

// Check if a file has the appropriate license header
function checkLicenseHeader(filePath, pattern) {
  try {
    const content = fs.readFileSync(filePath, 'utf8');
    const firstLines = content.split('\n').slice(0, 20).join('\n');
    return pattern.test(firstLines);
  } catch (error) {
    console.error(`Error reading file ${filePath}:`, error.message);
    return false;
  }
}

// Get all files matching patterns
function getFiles(patterns, excludePatterns) {
  const files = new Set();
  
  patterns.forEach(pattern => {
    const matches = glob.sync(pattern, {
      ignore: excludePatterns,
      nodir: true
    });
    matches.forEach(file => files.add(file));
  });
  
  return Array.from(files);
}

// Main validation function
function validateLicenses() {
  console.log('🔍 DeshChain License Validation Tool\n');
  
  let errors = 0;
  let warnings = 0;
  
  // Check source code files for Apache license
  console.log('📋 Checking source files for Apache 2.0 license headers...\n');
  const sourceFiles = getFiles(SOURCE_PATTERNS, EXCLUDE_PATTERNS);
  
  sourceFiles.forEach(file => {
    if (!checkLicenseHeader(file, APACHE_LICENSE_PATTERN)) {
      console.log(`❌ Missing Apache license: ${file}`);
      errors++;
    }
  });
  
  if (sourceFiles.length > 0 && errors === 0) {
    console.log(`✅ All ${sourceFiles.length} source files have proper Apache 2.0 headers\n`);
  }
  
  // Check cultural data for CC license
  console.log('🎭 Checking cultural data for CC BY-NC-SA 4.0 license...\n');
  const culturalFiles = getFiles(CULTURAL_PATTERNS, EXCLUDE_PATTERNS);
  
  let culturalErrors = 0;
  culturalFiles.forEach(file => {
    if (!checkLicenseHeader(file, CC_LICENSE_PATTERN)) {
      console.log(`❌ Missing CC license: ${file}`);
      culturalErrors++;
    }
  });
  
  if (culturalFiles.length > 0 && culturalErrors === 0) {
    console.log(`✅ All ${culturalFiles.length} cultural files have proper CC BY-NC-SA 4.0 headers\n`);
  }
  
  errors += culturalErrors;
  
  // Check for required license files
  console.log('📄 Checking required license files...\n');
  const requiredFiles = [
    'LICENSE',
    'LICENSE-CULTURAL',
    'NOTICE',
    'cultural-data/LICENSE'
  ];
  
  requiredFiles.forEach(file => {
    if (fs.existsSync(file)) {
      console.log(`✅ Found: ${file}`);
    } else {
      console.log(`❌ Missing: ${file}`);
      errors++;
    }
  });
  
  // Summary
  console.log('\n📊 Summary:\n');
  console.log(`Total source files checked: ${sourceFiles.length}`);
  console.log(`Total cultural files checked: ${culturalFiles.length}`);
  console.log(`Errors: ${errors}`);
  console.log(`Warnings: ${warnings}`);
  
  if (errors > 0) {
    console.log('\n❌ License validation failed!');
    console.log('Run `npm run add-licenses` to automatically add missing headers.');
    process.exit(1);
  } else {
    console.log('\n✅ All license checks passed!');
    console.log('DeshChain maintains proper dual licensing compliance.');
  }
}

// Run validation
if (require.main === module) {
  validateLicenses();
}

module.exports = { checkLicenseHeader, validateLicenses };