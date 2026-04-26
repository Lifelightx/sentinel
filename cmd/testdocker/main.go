package main

import (
	"fmt"
	
	"log"
	"sentinel/internal/dockerstats"
)

func main(){
	data, err := dockerstats.Collect()
	if err != nil{
		log.Fatal(err)
	}
	for _, c := range data{
		fmt.Printf("%s | CPU %.2f%% | MEM %.2fMB | %s\n",
			c.Name,
			c.CPUPercentage,
			c.MemoryMB,
			c.Health,
		)
	}

}