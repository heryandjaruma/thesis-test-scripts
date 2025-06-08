import os
import glob
import json
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

reports_dir = "/Users/heryandjaruma/Library/Mobile Documents/com~apple~CloudDocs/BINUS/TestScripts/reports-combined"
json_files = glob.glob(os.path.join(reports_dir, '*.json'))

TESTCASE="WF-R4"

all_latencies_data = []

for file_path in json_files:
    if TESTCASE in file_path:
        with open(file_path, 'r') as file:
            data = json.load(file)
            print(f"Loaded data from: {file_path}")
            
            # Extract iteration number from filename
            iteration = file_path.split('_')[-2]  # Gets the iteration number (1, 2, or 3)
            
            # Extract latencies data
            latencies_records = []
            
            for i, record in enumerate(data):
                # Latencies data - convert from nanoseconds to milliseconds
                latency_row = {
                    'iteration': iteration,
                    'second': i + 1,  # Time in seconds (1-50)
                    'mean': record['latencies']['mean'] / 1_000_000  # Convert ns to ms
                }
                latencies_records.append(latency_row)
            
            all_latencies_data.extend(latencies_records)

# Create DataFrame
latencies_df = pd.DataFrame(all_latencies_data)

# Display sample data
print("\nLatencies DataFrame:")
print(latencies_df.head())
print(f"Shape: {latencies_df.shape}")

# Create single plot for mean latencies only
fig, ax = plt.subplots(1, 1, figsize=(12, 8))

# Define similar color shades (different shades of blue) with transparency
colors = ['#1f77b4', '#5fa3d3', '#9ecae1']  # Light to dark blue shades
iteration_labels = ['Iteration 1', 'Iteration 2', 'Iteration 3']

# Plot mean latencies for each iteration with transparency and smaller dots
for i, iteration in enumerate(['1', '2', '3']):
    iteration_data = latencies_df[latencies_df['iteration'] == iteration]
    if not iteration_data.empty:
        ax.plot(iteration_data['second'], iteration_data['mean'], 
                color=colors[i], label=iteration_labels[i], 
                linewidth=1, marker='o', markersize=2, alpha=0.7)

ax.set_title('Mean Latencies Over Time by Iteration', fontsize=16, fontweight='bold')
ax.set_xlabel('Time (seconds)', fontsize=12)
ax.set_ylabel('Mean Latency (ms)', fontsize=12)
ax.legend(fontsize=11)
ax.grid(True, alpha=0.3)
ax.set_xlim(1, 50)

plt.tight_layout()
plt.show()