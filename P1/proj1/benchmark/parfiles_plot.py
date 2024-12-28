import matplotlib.pyplot as plt
import subprocess
import time


def run_test(data_dir, threads):
    start_time = time.time()

    # Run executable
    result = subprocess.run(["go", "run", "../editor/editor.go", data_dir, "parfiles", str(threads)], capture_output=True, text=True)

    # Access the output of the Go program
    print(result.stdout)

    end_time = time.time()
    runtime = end_time - start_time
    return runtime


data_dirs = ["small", "mixture", "big"]
threads = [1, 2, 4, 6, 8, 12]
results = {data_dir: [] for data_dir in data_dirs}

for data_dir in data_dirs:
    for t in threads:
        print(f"running {data_dir} with {t} threads")
        runtime = run_test(data_dir, t)
        results[data_dir].append(runtime)
        print(f"Data Dir: {data_dir}, Threads: {t}, Runtime: {runtime:.2f} seconds")

# proceed to plotting

print(results)

base_times = {data_dir: results[data_dir][0] for data_dir in results}  # Base time for 1 thread

# Plot for each data directory
for data_dir, times in results.items():
    speedups = [base_times[data_dir] / t for t in times]  # Calculate speedup
    plt.plot(threads, speedups, label=data_dir)

# Formatting the graph
plt.title("Speedup for Multiple Images in Parallel")
plt.xlabel("Number of Threads")
plt.ylabel("Speedup")
plt.xticks(threads)
plt.legend(title="File Size")
plt.grid()
plt.savefig("speedup-images.png")
plt.show()

