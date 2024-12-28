#!/bin/bash
#
#SBATCH --mail-user=sharapova@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj1_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/sharapova/autumn2024/Parallel/project-1-s-sharapova/proj1/benchmark
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=10:00


module load golang/1.19

# run python all go files through python to produce graph with results  
python slice_plot.py
python parfiles_plot.py

