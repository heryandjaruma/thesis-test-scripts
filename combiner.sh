#!/bin/bash

iteration=3
app="WF"
testcase="U6"
REPORTS_DIR="/Users/heryandjaruma/Library/Mobile Documents/com~apple~CloudDocs/BINUS/TestScripts/reports/"
OUTPUT_FILE="reports-combined/combined_${app}-${testcase}_${iteration}_reports.json"

echo "[" > "$OUTPUT_FILE"


first=true
for i in {1..50}; do
    file="${REPORTS_DIR}load_test_report_${app}-${testcase}_${iteration}_${i}.json"
    if [ -f "$file" ]; then
        if [ "$first" = true ]; then
            first=false
        else
            echo "," >> "$OUTPUT_FILE"
        fi
        cat "$file" >> "$OUTPUT_FILE"
    fi
done

echo "]" >> "$OUTPUT_FILE"

echo "Combined ${app}-${testcase}_${iteration} reports into $OUTPUT_FILE"