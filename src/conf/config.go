package conf

import (
    "runtime"
    "flag"
)

func SetCPUMax() {
    num := runtime.NumCPU()
    t := runtime.GOMAXPROCS(num)
}

type Config struct {
    freq    *int
    maxiter *int
    cost    *float64
    eta     *float64
    cpu     *int
}

func GetConfig() *Config {
    c := new(Config)
    c.freq = flag.Int("f", 3, "freq limit")
    c.maxiter = flag.Int("m", 100000, "max cycle count")
    f.cost = flag.Float64("c", 1.0, "set float for cost parameter(default 1.0)")
    f.eta = flag.Float64("e", 0.0001, "set float for termination criterion(default 0.0001)")
    f.cpu = flag.Int("CPU", -1, "Set cpu numbers used for CRF, default system max")

    flag.Parse()
    if (*f.Cpu == -1) {
        SetCPUMax()
    }else {
        num := runtime.NumCPU()
        if *f.Cpu < num {
            runtime.GOMAXPROCS(*f.Cpu)
        }else {
            runtime.GOMAXPROCS(num)
        }
    }
    return c
}

func (c *Config)GetFreq() int {
    return *c.Freq
}
func (c *Config)GetMaxLimit() int{
    return *c.maxiter
}
func (c *Config)GetCost() float64{
    return *c.cost
}
func (c *Config)GetEta() float64 {
    return *c.eta
}

