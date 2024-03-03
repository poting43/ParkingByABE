def parse_log_file(file_path):
    stage_times = {}
    
    with open(file_path, 'r') as file:
        lines = file.readlines()
        
        for line in lines:
            parts = line.strip().split(' : ')
            if len(parts) == 2:
                stage, time = parts
                stage_times[stage] = float(time[:-1])  # Remove 's' and convert to float
    
    return stage_times

def calculate_average_times(log_files):
    num_files = len(log_files)
    total_times = {stage: 0.0 for stage in ['Setup', 'Key extration for Bob', 'Encryption by Alice', 'Decryption by Bob']}
    
    for log_file in log_files:
        stage_times = parse_log_file(log_file)
        for stage, time in stage_times.items():
            total_times[stage] += time
    
    average_times = {stage: total_time / num_files for stage, total_time in total_times.items()}
    return average_times

if __name__ == "__main__":
    num_logs = 10
    log_files = [f"output{i}.txt" for i in range(1, num_logs + 1)]
    
    average_times = calculate_average_times(log_files)
    
    for stage, avg_time in average_times.items():
        print(f"Average {stage} time: {avg_time:.6f}s")
