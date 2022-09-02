#!/usr/bin/expect

log_file /home/oonya/quagga-ospf-autocost/expect.log
set PW "zebra"
set addr [lindex $argv 0]
set interface [lindex $argv 1]
set cost [lindex $argv 2]

set timeout 5
spawn env LANG=C /usr/bin/telnet $addr 2604
expect "Password:"
send "${PW}\n"

expect "ospfd>"
send "en\n"
expect "ospfd#"

send "configure t\n"
expect "ospfd(config)#"

send "interface ${interface}\n"
expect "ospfd(config-if)"

send "ip ospf cost ${cost}\n"
expect "ospfd(config-if)"

send "end\n"
expect "ospfd#"

exit 0
