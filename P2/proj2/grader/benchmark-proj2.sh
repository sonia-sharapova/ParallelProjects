#!/bin/bash
#
#SBATCH --mail-user=sharapova@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj2_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/sharapova/autumn2024/Parallel/project-2-s-sharapova/proj2/grader
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=5:00


module load golang/1.19
go test proj2/twitter -v -count=1
