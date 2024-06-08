import matplotlib.pyplot as plt
import sys
import json

job_id = sys.argv[1]
job_id = "./slurm/out/" + job_id +".slurm1.stdout"
with open(job_id, 'r') as file:

    threads = [2,4,6,8,12]
    types = [ "3000", "5000"]

    sequential_dict = {
                "3000": {},
                "5000": {}
    }
    parallel_dict = {
                "3000": {},
                "5000": {}
    }
    chunk_dict = {
                "3000": {},
                "5000": {}
    }
    chunk_sum = 0
    par_sum = 0
    seq_sum = 0
    thread_t = 1
    fail = 0
    last = 'n'
    for line in file:
        parts = line.split(',')  # Split each line into columns
        parts = [e.strip('\n') for e in parts]
        print(parts)
        if parts[0] == "Starting loop":
            print("yes")
            continue
        
        
        if parts[0] == '-----':
            if last == 's':
                if fail <5 :
                    seq_sum = seq_sum/(5-fail)
                    print("Average " ,seq_sum)
                    sequential_dict[mode] = seq_sum
                else:
                    seq_sum = 149
                    print("Average " ,seq_sum)
                    sequential_dict[mode] = seq_sum
            elif last == 'p':
                if fail<5:
                    par_sum = par_sum/5
                    print("Average " ,par_sum)
                    parallel_dict[mode].update({thread_t : par_sum})
                else:
                    par_sum = 149
                    print("Average " ,par_sum)
                    parallel_dict[mode].update({thread_t : par_sum})
            elif last == 'chunk':
                if fail<5:
                    chunk_sum = chunk_sum/5
                    print("Average " ,chunk_sum)
                    chunk_dict[mode].update({thread_t : chunk_sum})
                else:
                    chunk_sum = 149
                    print("Average " ,chunk_sum)
                    chunk_dict[mode].update({thread_t : chunk_sum})


            continue

        elif parts[0] == 's':
            mode = parts[1]
            seq_sum = 0
            last = 's'
            fail = 0
            continue
        
        elif parts[0] == 'p':       
            mode = parts[1]
            par_sum = 0
            last = 'p'
            fail = 0
            thread_t = int(parts[3])
            continue
        elif parts[0] == 'chunk':       
            mode = parts[1]
            chunk_sum = 0
            last = 'chunk'
            fail = 0
            thread_t = int(parts[3])
            continue
        elif last == 's':
            if parts[0] == "exit status 1":
                fail+=1
            else:
                seq_sum += float(parts[0])
        elif last == 'p':
            if parts[0] == "exit status 1":
                fail+=1
            else:
                par_sum += float(parts[0])
        elif last == 'chunk':
            if parts[0] == "exit status 1":
                fail+=1
            else:
                chunk_sum += float(parts[0])
            

        
            
pretty_dict = json.dumps(sequential_dict, indent=4)
print(pretty_dict)       

pretty_dict = json.dumps(parallel_dict, indent=4)
print(pretty_dict)

pretty_dict = json.dumps(chunk_dict, indent=4)
print("CHONKY")
print(pretty_dict)
print("CHUNKS")
# Create the plot

for typ in types:
    speedup_arr = []
    for t in threads:
        speedup_arr.append(sequential_dict[typ]/parallel_dict[typ][t])
    print(speedup_arr)
    plt.plot(threads,speedup_arr, label=typ)
    plt.xlabel('Threads')
    plt.ylabel('Speedup')
    plt.legend()
plt.title('Speedup Line Plot WorkStealing')

file_name = "Final_WorkSteal_"+str(sys.argv[1])+'.png'
plt.savefig(file_name, dpi=300, bbox_inches='tight')
plt.close()

for typ in types:
    speedup_arr = []
    for t in threads:
        speedup_arr.append(sequential_dict[typ]/chunk_dict[typ][t])
    print(speedup_arr)
    plt.plot(threads,speedup_arr, label=typ)
    plt.xlabel('Threads')
    plt.ylabel('Speedup')
    plt.legend()
plt.title('Speedup Line Basic Parallel Implement Plot')

file_name = "Final_Chunk_"+str(sys.argv[1])+'.png'
plt.savefig(file_name, dpi=300, bbox_inches='tight')
plt.close()
