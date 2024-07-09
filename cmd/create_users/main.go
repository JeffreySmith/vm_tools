/*BSD 3-Clause License

Copyright (c) 2024, Jeffrey Smith

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/JeffreySmith/vmtools"
	"io"
	"os"
	"strings"
)

func main() {
	var OutputBuffer io.Writer = os.Stdout
	var InputBuffer io.Reader
	var header string
	var ips []string
	stdin := os.Stdin
	f, err := stdin.Stat()

	ip := flag.String("ip", "", "Comma separated list of ip addresses.")
	output := flag.String("output", "", "Output file for generated yaml.")
	input := flag.String("input", "", "Input file for user names.")
	header_path := flag.String("header", "", "Path to a file containing your yaml file header (optional).")
	indentation_level := flag.Int("indent", 2, "Set the indentation level. Must be > 2")
	flag.Parse()

	rest := flag.Args()

	if len(rest) > 0 {
		ips = rest
	}
	if len(*ip) == 0 && len(rest) == 0 {
		fmt.Fprintf(os.Stderr, "You must supply at least 1 ip address\n\n")
		fmt.Fprintf(os.Stderr, "Pass them either as a comma separated list after '-ip'\n")
		fmt.Fprintf(os.Stderr, "or as a space separated list at the end of your arguments.\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	} else {
		if len(rest) > 0 && len(*ip) > 0 {
			ips = append(strings.Split(*ip, ","), ips...)
			
		} else if len(rest) == 0 && len(*ip) > 0{
			ips = strings.Split(*ip,",")
		} 
	}

	if f.Size() > 0 {
		InputBuffer = os.Stdin
	} else if len(*input) > 0 {
		var err error
		InputBuffer, err = os.Open(*input)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		InputBuffer, err = os.Open("users")
		if err != nil {
			fmt.Fprintf(os.Stderr, "You must suppy an input through stdin, a supplied file, or a 'users' file in this directory\n")
			os.Exit(1)
		}
	}

	if len(*output) > 0 {
		var err error
		OutputBuffer, err = os.Create(*output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if len(*header_path) > 0 {
		_, err := os.Stat(*header_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read file %v: %v\n", *header_path, err)
			os.Exit(1)
		}
		f, err := os.ReadFile(*header_path)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %v: \n", err)
			os.Exit(1)
		}
		header = string(f)
	} else {
		header = "---"
	}
	config := vmtools.NewConfig(vmtools.WithOutput(OutputBuffer),
		vmtools.WithInput(InputBuffer),
		vmtools.WithHeader(header),
		vmtools.SetIndent(*indentation_level),
	)

	config.GetUsers(ips)

	_, err = config.GenerateYaml()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.WriteYaml()
}
