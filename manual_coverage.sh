#!/bin/bash

echo "=== Manual Coverage Analysis ==="

if [ ! -f "coverage_new.out" ]; then
    echo "Coverage file not found"
    exit 1
fi

# Count total lines and covered lines
total_lines=$(grep -v "mode: set" coverage_new.out | wc -l)
covered_lines=$(grep -v "mode: set" coverage_new.out | awk '$3 > 0' | wc -l)

echo "Total statements: $total_lines"
echo "Covered statements: $covered_lines"

# Calculate percentage using bc if available, otherwise use awk
if command -v bc >/dev/null 2>&1; then
    percentage=$(echo "scale=2; $covered_lines * 100 / $total_lines" | bc)
else
    percentage=$(awk "BEGIN {printf \"%.2f\", $covered_lines * 100 / $total_lines}")
fi

echo "Coverage: $percentage%"

# Check if we meet the 97% target
if awk "BEGIN {exit !($percentage >= 97)}"; then
    echo "✅ Coverage target of 97% achieved!"
else
    needed=$(awk "BEGIN {printf \"%.2f\", 97 - $percentage}")
    echo "❌ Coverage is $percentage%, need $needed% more to reach 97%"
fi

echo ""
echo "=== Per-file breakdown ==="
echo "File coverage analysis:"

# Group by file and calculate per-file coverage
awk '
NR > 1 {
    # Extract filename from the first field
    split($1, parts, ":")
    filename = parts[1]
    gsub("github.com/khicago/repoll/", "", filename)
    
    total[filename]++
    if ($3 > 0) covered[filename]++
}
END {
    for (file in total) {
        if (total[file] > 0) {
            pct = (covered[file] / total[file]) * 100
            printf "%-30s %3d/%3d (%6.2f%%)\n", file, covered[file], total[file], pct
        }
    }
}' coverage_new.out | sort 