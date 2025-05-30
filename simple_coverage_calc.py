#!/usr/bin/env python3

# Read the coverage file and calculate coverage manually
with open('coverage_new.out', 'r') as f:
    lines = f.readlines()

total_lines = len(lines) - 1  # Exclude the first line "mode: set"
covered_lines = 0

print("=== Simple Coverage Calculation ===")
print(f"Total lines in coverage file: {len(lines)}")
print(f"Data lines (excluding mode line): {total_lines}")

# Count covered lines (where the last number > 0)
for line in lines[1:]:  # Skip first line
    parts = line.strip().split()
    if len(parts) >= 3:
        coverage_count = int(parts[2])
        if coverage_count > 0:
            covered_lines += 1

coverage_percent = (covered_lines / total_lines) * 100 if total_lines > 0 else 0

print(f"Covered lines: {covered_lines}")
print(f"Coverage percentage: {coverage_percent:.2f}%")

if coverage_percent >= 97.0:
    print("✅ Coverage target of 97% achieved!")
else:
    needed = 97.0 - coverage_percent
    print(f"❌ Coverage is {coverage_percent:.2f}%, need {needed:.2f}% more to reach 97%")

# Show some sample lines for verification
print("\n=== Sample coverage lines ===")
for i, line in enumerate(lines[1:6]):  # Show first 5 data lines
    parts = line.strip().split()
    if len(parts) >= 3:
        print(f"Line {i+1}: {parts[0]} -> {parts[2]} ({'covered' if int(parts[2]) > 0 else 'not covered'})") 