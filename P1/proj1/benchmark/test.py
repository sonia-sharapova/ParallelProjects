import matplotlib.pyplot as plt
data_dirs = ["small", "mixture", "big"]
threads = [1, 2, 4, 6, 8, 12]

results = {'small': [6.66, 6.19, 5.63, 5.58, 6.39, 6.38], 
           'big':[110.96, 85.9, 85.05, 81.27, 78.41, 78.16],
           'mixture':[47.54, 40.03, 38.46, 38.04, 35.15, 34.76]}
'''
results = {'small': [19.63, 8.3, 5.44, 4.12, 4.24, 3.69], 
           'big':[217.97, 112.45, 86.38, 72.23, 67.75, 50.92],
           'mixture':[92.58, 60.33, 34.59, 36.9, 32.65, 34.2]}
'''

base_times = {data_dir: results[data_dir][0] for data_dir in results}  # Base time for 1 thread

# Plot for each data directory
for data_dir, times in results.items():
    speedups = [base_times[data_dir] / t for t in times]  # Calculate speedup
    plt.plot(threads, speedups, label=data_dir)

# Formatting the graph
#plt.title("Speedup for Parallelizing Image Files")
plt.title("Speedup for Parallelizing Images with Slicing")
plt.xlabel("Number of Threads")
plt.ylabel("Speedup")
plt.xticks(threads)
plt.legend(title="Image Size")
plt.grid()
plt.savefig("test-slice.png")
plt.show()