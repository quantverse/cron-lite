package main

import (
	"os"
	"os/exec"
	"log"
	"flag"
	"io/ioutil"
	"path/filepath"

	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"

)

type ScheduledJob struct {
	Command string	`yaml:"command"`
	Args	string  `yaml:"args"`
	Rule	string  `yaml:"rule"`
}

func (sj ScheduledJob) Run() {
	log.Print(sj)
	cmd := exec.Command(sj.Command, sj.Args);
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Print("command failed to start")
		log.Print(err)
	}
}

const version = "0.1.0"

func main() {
	fnamePtr := flag.String("f", "", "YAML file name with rules")
	flag.Parse()

	if *fnamePtr=="" {
		log.Fatal("Please specify filename with -f")
	}


	log.Printf("cron-lite v%s by Quantverse", version)
	log.Printf("loading config from %s", *fnamePtr)

	filename, _ := filepath.Abs(*fnamePtr)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	m := make([]ScheduledJob, 0)

	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		panic(err)
	}
	cnt := 0
	c := cron.New()
	for _, item := range m {
		log.Printf("adding new job: [%s %s]: %s\n", item.Command, item.Args, item.Rule)
		c.AddJob(item.Rule, item);
		cnt++;
	}
	log.Printf("%d jobs added\n", cnt)
	if cnt > 0 {
		c.Start()
		log.Print("cron engine started")
		select{}
	}

}
