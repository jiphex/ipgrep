// ipgrep - scan STDIN for IP(v4 or v6) addresses that match the
//          CIDR addresses specified as arguments.
//
// Usage: ipgrep 1.2.3.4/24 2.3.4.5/12 2001:db0:1::30/64 ... < listofips

package main

import (
  "strings"
  "net"
  "fmt"
  "os"
  "bufio"
  "flag"
  "sync"
)

func str2net(s string) (ipm *net.IPNet) {
  if !strings.Contains(s, "/") {
    s += "/32"
  }
  _,ipm,err := net.ParseCIDR(s)
  if err != nil {
    ipm = nil
  }
  return
}

type linelast struct {
  line string
  last bool
}

func ipMatches(target *net.IPNet, queue <-chan linelast, matches chan<- net.IP, wg *sync.WaitGroup) {
  // fmt.Printf("Net matching %s\n", target)
  for {
    linelast := <-queue
    line := linelast.line
    lscan := bufio.NewScanner(strings.NewReader(line))
    lscan.Split(bufio.ScanWords)
    for lscan.Scan() {
      ip := net.ParseIP(lscan.Text())
      // fmt.Printf("Test: [%s]\n", ip)
      // fmt.Println(ip)
      if ip != nil && (*target).Contains(ip) {
        matches <- ip
      }
    }
    if linelast.last {
      wg.Done()
    }
  }
}

func main() {
  verbose := flag.Bool("verbose", false, "Show errors/stuff")
  flag.Parse()
  s := make([]chan linelast, 0)
  matches := make(chan net.IP)
  wait := &sync.WaitGroup{}
  for _,i := range flag.Args() {
    m := str2net(i)
    if m != nil {
      recp := make(chan linelast)
      wait.Add(1)
      go ipMatches(m, recp, matches, wait)
      s = append(s,recp)
    } else {
      if *verbose {
        fmt.Fprintf(os.Stderr, "Not an IP/CIDR: %s\n", i)
      }
    }
  }
  scanner := bufio.NewScanner(os.Stdin)
  go func() {
    m := make(map[string]int)
    for {
      match := (<- matches).String()
      if m[match]>0 {
        if *verbose {
          fmt.Printf("For the %d'th time: %s\n", m[match]+1, match)
        }
      } else {
        fmt.Println(match) 
      }
      m[match]+=1
    }
  }()
  for scanner.Scan() {
    line := scanner.Text()
    for _,i := range s {
      i <- linelast{line,false}
    }
  }
  for _,i := range s {
    i <- linelast{"",true}
  }
  wait.Wait()
}
