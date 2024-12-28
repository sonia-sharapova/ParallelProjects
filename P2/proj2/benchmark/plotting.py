import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

# Load the data from the benchmark file
benchmark_file = "benchmark_results.csv"
data = pd.read_csv(benchmark_file)

# Calculate speedup: sequential_time / parallel_time
# Handle division by zero or invalid parallel_time
data['speedup'] = data.apply(
    lambda row: row['sequential_time'] / row['parallel_time'] if row['parallel_time'] > 0 else np.nan,
    axis=1
)

# Plot speedups for each test size
test_sizes = data['test_size'].unique()

plt.figure(figsize=(12, 6))
for test_size in test_sizes:
    subset = data[data['test_size'] == test_size]
    cleaned_speedup = subset.dropna(subset=['speedup'])  # Remove rows with NaN in speedup
    plt.plot(cleaned_speedup['threads'], cleaned_speedup['speedup'], marker='o', label=test_size)

# Add labels and legend
plt.xlabel("Number of Threads")
plt.ylabel("Speedup")
plt.title("Speedup vs Threads for Different Test Sizes")
plt.legend(title="Test Size")
plt.grid(True)

# Save the figure
output_file = "speedup_graph.png"
plt.savefig(output_file, dpi=300, bbox_inches="tight")  # Save as a high-resolution image

# Show the plot
plt.show()

print(f"Plot saved as {output_file}")
