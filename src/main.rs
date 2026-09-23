use std::{
    fs::File,
    io::{BufRead, BufReader},
};

struct MemoryStats {
    total_memory: u64,
    available_memory: u64,
    total_swap: u64,
    free_swap: u64,
    cached: u64,
}

impl MemoryStats {
    pub fn fetch() -> Result<MemoryStats, Box<dyn std::error::Error>> {
        let memory_info_file = File::open("/proc/meminfo")?;

        let reader = BufReader::new(memory_info_file);
        let mut memory_stats = MemoryStats {
            total_memory: 0,
            available_memory: 0,
            total_swap: 0,
            free_swap: 0,
            cached: 0,
        };

        for line in reader.lines() {
            let line = line?;

            if let Some((key, value)) = line.split_once(":") {
                let number: u64 = value.split_whitespace().next().unwrap().parse().unwrap();
                match key {
                    "MemTotal" => memory_stats.total_memory = number,
                    "MemAvailable" => memory_stats.available_memory = number,
                    "SwapTotal" => memory_stats.total_swap = number,
                    "SwapFree" => memory_stats.free_swap = number,
                    "Cached" => memory_stats.cached = number,
                    _ => (),
                }
            }
        }

        assert!(
            memory_stats.total_memory > 0,
            "total memory available not parsed correct"
        );

        Ok(memory_stats)
    }
}

struct CPUStats {
    idle: u64,
    io_wait: u64,
    user: u64,
    nice: u64,
    system: u64,
    irq: u64,
    softirq: u64,
}

impl CPUStats {
    pub fn fetch() -> Result<CPUStats, Box<dyn std::error::Error>> {
        let cpu_info_file = File::open("/proc/stat")?;
        let mut cpu_reader = BufReader::new(cpu_info_file);
        let mut cpu_stats = CPUStats {
            idle: 0,
            io_wait: 0,
            user: 0,
            nice: 0,
            system: 0,
            irq: 0,
            softirq: 0,
        };
        let mut line = String::new();

        cpu_reader.read_line(&mut line)?;

        let mut values = line.split_whitespace();
        values.next();

        let fields = [
            &mut cpu_stats.user,
            &mut cpu_stats.nice,
            &mut cpu_stats.system,
            &mut cpu_stats.idle,
            &mut cpu_stats.io_wait,
            &mut cpu_stats.irq,
            &mut cpu_stats.softirq,
        ];

        for (field, value) in fields.into_iter().zip(values) {
            if let Ok(value) = value.parse::<u64>() {
                *field = value;
            }
        }

        assert!(cpu_stats.user > 0, "CPU stats not parsed correctly");

        Ok(cpu_stats)
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mem_stats = MemoryStats::fetch()?;
    println!("total: {}", mem_stats.total_memory);
    let cpu_stats = CPUStats::fetch()?;
    println!("user: {}, nice: {}", cpu_stats.user, cpu_stats.nice);
    Ok(())
}
