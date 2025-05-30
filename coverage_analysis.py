#!/usr/bin/env python3

import sys
import os

def analyze_coverage():
    if not os.path.exists('coverage_new.out'):
        print("Coverage file not found")
        return
    
    file_stats = {}
    total_statements = 0
    covered_statements = 0
    
    with open('coverage_new.out', 'r') as f:
        lines = f.readlines()
    
    # Skip the first line (mode: set)
    for line in lines[1:]:
        line = line.strip()
        if not line:
            continue
        
        parts = line.split()
        if len(parts) < 3:
            continue
        
        # Parse: filename:start.col,end.col statements covered
        file_part = parts[0]
        covered = int(parts[2])
        
        # Extract filename
        colon_index = file_part.rfind(':')
        if colon_index == -1:
            continue
        
        filename = file_part[:colon_index]
        filename = filename.replace('github.com/khicago/repoll/', '')
        
        if filename not in file_stats:
            file_stats[filename] = {'total': 0, 'covered': 0}
        
        file_stats[filename]['total'] += 1
        total_statements += 1
        
        if covered > 0:
            file_stats[filename]['covered'] += 1
            covered_statements += 1
    
    # Calculate overall coverage
    overall_percent = (covered_statements / total_statements * 100) if total_statements > 0 else 0
    
    print("=== Coverage Analysis ===")
    print(f"Total statements: {total_statements}")
    print(f"Covered statements: {covered_statements}")
    print(f"Overall coverage: {overall_percent:.2f}%")
    print()
    
    # Per-file breakdown
    print("=== Per-file Coverage ===")
    for filename in sorted(file_stats.keys()):
        stats = file_stats[filename]
        if stats['total'] > 0:
            percent = (stats['covered'] / stats['total']) * 100
            print(f"{filename:<30} {stats['covered']:3d}/{stats['total']:3d} ({percent:6.2f}%)")
    
    print("=" * 60)
    print(f"{'TOTAL':<30} {covered_statements:3d}/{total_statements:3d} ({overall_percent:6.2f}%)")
    
    if overall_percent >= 97.0:
        print("✅ Coverage target of 97% achieved!")
    else:
        needed = 97.0 - overall_percent
        print(f"❌ Coverage is {overall_percent:.2f}%, need {needed:.2f}% more to reach 97%")

if __name__ == "__main__":
    analyze_coverage() 