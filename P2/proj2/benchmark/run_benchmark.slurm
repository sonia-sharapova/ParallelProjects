#!/bin/bash
#
#SBATCH --mail-user=sharapova@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=benchmark_full
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/sharapova/autumn2024/Parallel/project-2-s-sharapova/proj2/benchmark
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=50:00

# load Go module
module load golang/1.19


# Set up environment variables
BENCHMARK_PROGRAM="benchmark.go"
RESULTS_FILE="benchmark_results.csv"

# Test sizes and thread counts
TEST_SIZES=("xsmall" "small" "medium" "large" "xlarge")
THREAD_COUNTS=(2 4 6 8 12)
RUNS=5  # Number of runs for averaging

# Initialize the results CSV
echo "test_size,threads,sequential_time,parallel_time" > $RESULTS_FILE

# Function to calculate the average time from multiple runs
calculate_average() {
    local times=("$@")
    local sum=0
    for time in "${times[@]}"; do
        sum=$(echo "$sum + $time" | bc)
    done
    echo "scale=3; $sum / ${#times[@]}" | bc
}

# Run sequential benchmarks
declare -A SEQUENTIAL_TIMES
for test_size in "${TEST_SIZES[@]}"; do
    echo "Running sequential benchmark for test size: $test_size"
    seq_times=()
    for ((i=1; i<=RUNS; i++)); do
        time=$(go run $BENCHMARK_PROGRAM s $test_size)
        seq_times+=("$time")
    done
    SEQ_AVG=$(calculate_average "${seq_times[@]}")
    SEQUENTIAL_TIMES[$test_size]=$SEQ_AVG
    echo "$test_size,1,$SEQ_AVG," >> $RESULTS_FILE
done

# Run parallel benchmarks
for test_size in "${TEST_SIZES[@]}"; do
    for threads in "${THREAD_COUNTS[@]}"; do
        echo "Running parallel benchmark for test size: $test_size with $threads threads"
        par_times=()
        for ((i=1; i<=RUNS; i++)); do
            time=$(go run $BENCHMARK_PROGRAM p $test_size $threads)
            par_times+=("$time")
        done
        PAR_AVG=$(calculate_average "${par_times[@]}")
        SEQ_TIME=${SEQUENTIAL_TIMES[$test_size]}  # Fetch the sequential time
        echo "$test_size,$threads,$SEQ_TIME,$PAR_AVG" >> $RESULTS_FILE
    done
done

echo "Benchmarking completed. Results stored in $RESULTS_FILE."