# vmtools
Tool to generate valid yaml files.

Currently, you can generate yaml to add additional users to an arbitrary number of virtual machines.

Usernames can come from stdin, a file specified by the -input paramater, or (as a last ditch effort) from a file called users in the directory of the binary. Input is prioritized as stdin(highest), -input paramater, users file (lowest). When using `-ip`, you must input a comma separated list. If you would prefer a space separated list, add your IP addresses at the end, after all other commandline flags/options.

Example usage:

`echo johndoe janedoe | ./adduser -ip 10.90.9.9,192.168.1.4 `

will generate:
```
---
additional_users:
  - username: johndoe
    vm_ip: 10.90.9.9
  - username: janedoe
    vm_ip: 10.90.9.9
  - username: johndoe
    vm_ip: 192.168.1.4
  - username: janedoe
    vm_ip: 192.168.1.4
```
Example with the IPs at the end
```
echo johndoe | ./adduser -indent 4  10.90.9.9 192.168.1.4
```

You can set the username input file directy using -input:
`./adduser -input list_of_users -ip 10.90.9.8,192.168.1.4`

An entry will be made for each user contained in the user file, for each of the ip addresses you've supplied.

The default output is stdout, but you can overide this by specifying a `-output $filename` argument. Or you can pipe the output of this into other commands (all errors print to stderr). Additionally, you can set the size of indentation using -indent $(number value), although due to a limitation of the underlying library, this must be at least 2. Anything below 2 will be set to 2. 

You can also specify a header using the `-header` flag and point it to a file that contains what you would like to be at the top of your yaml output. Without this file, it goes to the default behaviour of yaml: `---`. 
