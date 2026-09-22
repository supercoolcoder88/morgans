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
    fn empty() -> MemoryStats {
        MemoryStats {
            total_memory: 0,
            available_memory: 0,
            total_swap: 0,
            free_swap: 0,
            cached: 0,
        }
    }

    fn parse(&mut self, key: &str, value: &str) {
        let number: u64 = value.split_whitespace().next().unwrap().parse().unwrap();
        match key {
            "MemTotal" => self.total_memory = number,
            "MemAvailable" => self.available_memory = number,
            "SwapTotal" => self.total_swap = number,
            "SwapFree" => self.free_swap = number,
            "Cached" => self.cached = number,
            _ => (),
        }
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
    fn empty() -> CPUStats {
        CPUStats {
            idle: 0,
            io_wait: 0,
            user: 0,
            nice: 0,
            system: 0,
            irq: 0,
            softirq: 0,
        }
    }

    fn parse(&mut self, line: &str) {
        let mut values = line.split_whitespace();
        values.next(); // skip "cpu"

        let fields = [
            &mut self.user,
            &mut self.nice,
            &mut self.system,
            &mut self.idle,
            &mut self.io_wait,
            &mut self.irq,
            &mut self.softirq,
        ];

        for (field, value) in fields.into_iter().zip(values) {
            if let Ok(value) = value.parse::<u64>() {
                *field = value;
            }
        }
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let memory_info_file = File::open("/proc/meminfo")?;

    let reader = BufReader::new(memory_info_file);
    let mut memory_stats = MemoryStats::empty();

    for line in reader.lines() {
        let line = line?;

        if let Some((key, value)) = line.split_once(":") {
            memory_stats.parse(key, value);
        }
    }

    println!(
        "total_memory: {}, available_memory: {}, total_swap: {}, free_swap: {}, cached: {}",
        memory_stats.total_memory,
        memory_stats.available_memory,
        memory_stats.total_swap,
        memory_stats.free_swap,
        memory_stats.cached,
    );

    let cpu_info_file = File::open("/proc/stat")?;
    let mut cpu_reader = BufReader::new(cpu_info_file);
    let mut cpu_stats = CPUStats::empty();
    let mut line = String::new();

    cpu_reader.read_line(&mut line)?;
    cpu_stats.parse(line.as_str());
    println!("user: {}, nice: {}", cpu_stats.user, cpu_stats.idle,);
    Ok(())
}
