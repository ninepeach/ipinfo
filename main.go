package main

import (
    "bufio"
    "bytes"
    "fmt"
    "net"
    "os"
    "sort"
    "strings"
)

func main() {
    // Open the file containing the CIDRs
    file, err := os.Open("test")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    // Read the CIDRs from the file into a slice
    var cidrs []*net.IPNet
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        _, cidr, err := net.ParseCIDR(strings.TrimSpace(scanner.Text()))
        if err != nil {
            fmt.Println("Error parsing CIDR:", err)
            continue
        }

        cidrs = append(cidrs, cidr)
    }
    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    // Merge the CIDRs
    merged := mergeCIDRs(cidrs)

    // Print the merged CIDRs
    fmt.Println("Merged CIDRs:")
    for _, cidr := range merged {
        fmt.Println(cidr.String())
    }
}

func mergeCIDRs(cidrs []*net.IPNet) []*net.IPNet {
    // Sort the CIDRs by their IP addresses
    sort.Slice(cidrs, func(i, j int) bool {
        return bytes.Compare(cidrs[i].IP, cidrs[j].IP) < 0
    })

    // Initialize a slice to hold the merged CIDRs
    var merged []*net.IPNet

    // Iterate through the sorted CIDRs and merge adjacent ones
    for i := 0; i < len(cidrs); i++ {
        // Check if the current CIDR overlaps with the previous one
        if i > 0 && cidrs[i-1].Contains(cidrs[i].IP) {
            // Merge the current CIDR with the previous one
            tmpCidr : = mergeIPNets(merged[len(merged)-1], cidrs[i])
            if tmpCidr!=nil {
                merged[len(merged)-1] = mergeIPNets(merged[len(merged)-1], cidrs[i])
            }
        } else {
            // Add the current CIDR to the merged slice
            merged = append(merged, cidrs[i])
        }
    }

    return merged
}


func mergeIPNets(a, b *net.IPNet) *net.IPNet {
    // Find the common prefix length of the two IP networks
    aLen, _ := a.Mask.Size()
    bLen, _ := b.Mask.Size()
    prefixLen := aLen
    if bLen < aLen {
        prefixLen = bLen
    }
    // Compare the common prefix bits to make sure they match
    aPrefix := a.IP.To4()[:prefixLen/8]
    bPrefix := b.IP.To4()[:prefixLen/8]
    if bytes.Compare(aPrefix, bPrefix) != 0 || prefixLen%8 != 0 {
        return nil
    }
    // Merge the two IP networks
    var merged net.IPNet
    merged.IP = a.IP.Mask(net.CIDRMask(prefixLen, 32))
    merged.Mask = net.IPNet{
        Mask: net.IPv4Mask(
            net.CIDRMask(prefixLen, 32)[0],
            net.CIDRMask(prefixLen, 32)[1],
            net.CIDRMask(prefixLen, 32)[2],
            net.CIDRMask(prefixLen, 32)[3],
        ),
    }.Mask
    return &merged
}
