#!/usr/bin/expect

log_file /var/log/expect.log
set PW "zebra"
set interface [lindex $argv 0]
set cost [lindex $argv 1]

set timeout 5
spawn env LANG=C /usr/bin/telnet localhost 2604
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
