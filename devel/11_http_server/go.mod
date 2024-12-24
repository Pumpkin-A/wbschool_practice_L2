module main

replace mymain => ./mymain

replace config => ./mymain/config

go 1.16

require mymain v0.0.0-00010101000000-000000000000 // indirect
