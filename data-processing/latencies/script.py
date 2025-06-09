import os
import glob
import json
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

reports_dir = "/Users/heryandjaruma/Library/Mobile Documents/com~apple~CloudDocs/BINUS/TestScripts/reports-combined"
json_files = glob.glob(os.path.join(reports_dir, '*.json'))

CASE = 6
OPERATION = "D"
TESTCASE1=f"WF-{OPERATION}{CASE}"
TESTCASE2=f"VT-{OPERATION}{CASE}"

all_latencies_data = []

# Process both test cases
for testcase in [TESTCASE1, TESTCASE2]:
    for file_path in json_files:
        if testcase in file_path:
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
                        'testcase': testcase,
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

# Calculate mean across all iterations for each test case
mean_by_testcase = latencies_df.groupby(['testcase', 'second'])['mean'].mean().reset_index()

# Create single plot for mean latencies comparison
fig, ax = plt.subplots(1, 1, figsize=(12, 8))

# Define better color groups with stronger distinction
# WF-R4: Blue family, VT-R4: Orange family
colors = {
    TESTCASE1: ['#1f77b4', '#4a90e2', '#7bb3f0'],  # Blue shades (dark to light)
    TESTCASE2: ['#ff7f0e', '#ff9f40', '#ffbf73']   # Orange shades (dark to light)
}

# Mean line colors (darker versions)
mean_colors = {
    TESTCASE1: '#0d4f8c',  # Darker blue
    TESTCASE2: '#cc5500'   # Darker orange
}

iteration_labels = ['Iterasi 1', 'Iterasi 2', 'Iterasi 3']

# Plot mean latencies for each test case and iteration
for testcase in [TESTCASE1, TESTCASE2]:
    for i, iteration in enumerate(['1', '2', '3']):
        iteration_data = latencies_df[(latencies_df['testcase'] == testcase) & 
                                    (latencies_df['iteration'] == iteration)]
        if not iteration_data.empty:
            label = f"{testcase} - {iteration_labels[i]}"
            ax.plot(iteration_data['second'], iteration_data['mean'], 
                    color=colors[testcase][i], label=label, 
                    linewidth=1, alpha=0.8)

# Plot mean lines for each test case
for testcase in [TESTCASE1, TESTCASE2]:
    testcase_mean_data = mean_by_testcase[mean_by_testcase['testcase'] == testcase]
    if not testcase_mean_data.empty:
        ax.plot(testcase_mean_data['second'], testcase_mean_data['mean'], 
                color=mean_colors[testcase], label=f"{testcase} - Rata-rata", 
                linewidth=1.5, linestyle='--', alpha=1)

ax.set_title(f'Perbandingan Rata-rata Latensi: {TESTCASE1} dan {TESTCASE2}', fontsize=16, fontweight='bold')
ax.set_xlabel('Detik Ke-', fontsize=12)
ax.set_ylabel('Rata-rata Latensi (ms)', fontsize=12)
ax.legend(fontsize=10, loc='upper right')
ax.grid(True, alpha=0.3)
ax.set_xlim(1, 50)



# Export results in the requested format
results_dir = "results"
os.makedirs(results_dir, exist_ok=True)

# Create filename in the format TESTCASE1_vs_TESTCASE2_latencies
filename_base = f"{TESTCASE1}_vs_{TESTCASE2}_latencies"
plt.tight_layout()

# Save the plot
plt.savefig(f"{results_dir}/{filename_base}.png", dpi=300, bbox_inches='tight')

plt.show()