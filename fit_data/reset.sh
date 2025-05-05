rm data.db
printf "set fit true\nset fit:users 0\nwrite\n" | ./key_value_db data.db > /dev/null
