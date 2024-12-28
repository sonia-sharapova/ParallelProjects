#!/bin/bash

#SBATCH --mail-user=sharapova@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=benchmark_full
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/sharapova/autumn2024/Parallel/project-3-s-sharapova/proj3/benchmark
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=50:00

# Load required modules
module load golang/1.19

# Create directories for output if they don't exist
mkdir -p ./slurm/out
mkdir -p ./benchmark_results

# Set environment variables
export GOMAXPROCS=$SLURM_CPUS_PER_TASK
INPUT_DIR="../smallerData"  # Update this path

echo "Starting benchmark run at $(date)"
echo "Running on host: $(hostname)"
echo "Number of CPUs allocated: $SLURM_CPUS_PER_TASK"
echo "Input directory: $INPUT_DIR"

# Run the benchmark program
echo "Running benchmarks..."
go run ../cmd/benchmark/main.go -input "$INPUT_DIR" -maxworkers $SLURM_CPUS_PER_TASK

# If Python plotting script exists, generate plots
if [ -f "../scripts/plotting.py" ]; then
    echo "Generating plots..."
    python ../scripts/plotting.py
    
    # Check if plot generation was successful
    if [ -f "./benchmark_results/speedup.png" ]; then
        echo "Plot generated successfully"
    else
        echo "Warning: Plot generation may have failed"
    fi
else
    echo "Warning: Plotting script not found at ../scripts/plotting.py"
fi

echo "Benchmark completed at $(date)"

# Print summary of results location
echo "Results can be found in:"
echo "- CSV data: ./benchmark_results/results.csv"
echo "- Speedup plot: ./benchmark_results/speedup.png"

