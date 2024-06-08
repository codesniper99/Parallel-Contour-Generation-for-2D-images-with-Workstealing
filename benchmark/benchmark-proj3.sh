#!/bin/bash
#
#SBATCH --mail-user=avaid@uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj3_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.input_f'input_fileN.stderr
#SBATCH --chdir=/home/avaid/Desktop/parallelProgramming/project-3/proj3/benchmark
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=4:00:00


module load golang/1.19
# Your command here
cd ..
types=("3000" "5000")
modes=("s" "p" "chunk")
values=("2" "4" "6" "8" "12")
runs=("1" "2" "3" "4" "5")

echo "Starting loop"

for mode in "${modes[@]}" 
do
    for type in "${types[@]}" 
    do 
        if [ "$mode" == "s" ] 
        then
            echo "s,$type,$mode" 
            for run in "${runs[@]}"
            do
                go run runner.go "$mode" "$type"
                
            done
            echo "-----"

        else 
            for value in "${values[@]}"
            do
                echo "$mode,$type,$mode,$value" 
                for run in "${runs[@]}"
                do
                    go run runner.go "$mode" "$type" "$value"
                    
                done
                echo "-----"
            done
        fi
    done
done
cd benchmark
python3 plot_graph.py $SLURM_JOB_ID
